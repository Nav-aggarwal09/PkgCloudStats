/*
 * $Header$
 * This program gets the download and install counts for a package in packagecloud.io
 * See https://packagecloud.io/docs/api
 * Note that the downloads API is for each file while the installs API for each repo
 *
 * You have to list all the repo, package, version details in the config file.
 *
 * AUTHOR:
 *      Arnav Aggarwal, arnav0908_at_gmail.com
 */

package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strings"
)

const DEFAULT_START_DATEZ = "19900101Z"
const DEFAULT_CONFIG = "Pkgcloud-Counter-config.json"

var DebugLvl bool
var API_TOKEN string

// structure of Json input file
type CounterHeader struct {
	API_Token string      `json:"API_TOKEN"`
	User      string      `json:"user"`
	UserRepo  string      `json:"userrepo"`
	Package   string      `json:"packagename"`
	Packages  []SensuPkg  `json:"packages"`
	Repos     []SensuRepo `json:"repositories"`
}

type SensuPkg struct {
	PkgVersion string   `json:"pkgversion"` // 5.20.6
	Release    string   `json:"release"`    //3425
	Arch       string   `json:"arch"`
	Distro     string   `json:"distro_p"`
	Version    []string `json:"version_p"`
	Downloads  []int    //for each version
}

type SensuRepo struct {
	Distro   string `json:"distro_r"`  // Ubuntu, el, etc
	Version  string `json:"version_r"` // 5, 6, focal, bionic
	Installs int
}

type CloudResponse struct {
	Value int `json:"value"`
}

func getDownloads(pkgptr *SensuPkg, pkgname, startdatez, user, repo string) {

	for i, ver := range pkgptr.Version {
		url := fmt.Sprintf("https://%v:@packagecloud.io/api/v1/repos/%v/%v/"+
			"package/rpm/%v/%v/%v/%v/%v/%v/stats/downloads/count.json?start_date=%v",
			API_TOKEN, user, repo, pkgptr.Distro, ver,
			pkgname, pkgptr.Arch, pkgptr.PkgVersion, pkgptr.Release, startdatez)
		response, err := http.Get(url)
		if err != nil {
			log.Fatal("Get request failed: ", err.Error())
		}

		responseData, err := ioutil.ReadAll(response.Body)
		if err != nil {
			log.Fatal("Error reading data: ", err)
		}
		var cr CloudResponse
		json.Unmarshal(responseData, &cr)
		if DebugLvl {
			log.Printf("pkgVer %v had %v downloads", pkgptr.PkgVersion, cr.Value)
		}

		if len(pkgptr.Downloads)-1 >= i {
			pkgptr.Downloads[i] = cr.Value
		} else {
			pkgptr.Downloads = append(pkgptr.Downloads, cr.Value)
		}

	}
	if DebugLvl {
		log.Printf("Finished download count for pkg ver %v", pkgptr.PkgVersion)
	}

} // end getDownloads()

func getInstalls(sensuptr *SensuRepo, startdatez, user, repo string) {

	// ubuntu = distro
	// version = 67878
	url := fmt.Sprintf("https://%v:@packagecloud.io/api/v1/repos/%v/%v/"+
		"/stats/installs/%v/%v/count.json?start_date=%v",
		API_TOKEN, user, repo, sensuptr.Distro, sensuptr.Version, startdatez)
	response, err := http.Get(url)

	if err != nil {
		log.Fatal("Get request failed: ", err.Error())
	}

	responseData, err := ioutil.ReadAll(response.Body)
	if err != nil {
		log.Fatal("Error reading data: ", err)
	}

	var cr CloudResponse
	json.Unmarshal(responseData, &cr)
	if DebugLvl {
		log.Printf("%v/%v had %v installs", sensuptr.Distro, sensuptr.Version, cr.Value)
	}

	sensuptr.Installs = cr.Value

} // end getinstalls()

func main() {

	/**
	 * Read and parse the JSON configuration file
	 */
	var config string

	filePtr := flag.String("config", DEFAULT_CONFIG, "path to config file")
	debugPtr := flag.Bool("debug", false, "Log level")
	flag.Parse()

	DebugLvl = *debugPtr // set log level
	log.Println("Set debug to ", DebugLvl)

	if strings.EqualFold(*filePtr, "ENV") {
		if os.Getenv("PKGCLOUD_API_TOKEN") == "" {
			log.Fatal("Cannot find Pkgcloud API token in enviornment")
		} else {
			config = os.Getenv("PKGCLOUD_API_TOKEN")
		}
	} else {
		config = *filePtr
	}

	log.Println("Set file to ", config)

	file, err := os.Open(config)
	if err != nil {
		log.Fatalf("Error opening file %v :", config, err)
	}

	defer file.Close()
	filebytes, err := ioutil.ReadAll(file)
	if err != nil {
		log.Fatal("Error reading file: ", err)
	}

	var parsedfile CounterHeader
	err = json.Unmarshal(filebytes, &parsedfile)
	if err != nil {
		log.Fatal("Could not parse Json: ", err)
	}
	// to test if json was read
	if parsedfile.API_Token == "" {
		log.Fatal("Cannot parse API_TOKEN in file")

	}
	if DebugLvl {
		log.Println("Parsed JSON file")
	}

	API_TOKEN = parsedfile.API_Token

	/**
	 * Get the DOWNLOADS count
	 */
	log.Println("Starting download count for given packages")
	for i := range parsedfile.Packages {
		getDownloads(&parsedfile.Packages[i], parsedfile.Package, DEFAULT_START_DATEZ,
			parsedfile.User, parsedfile.UserRepo)
		if DebugLvl {
			log.Printf("Package %v has %v downloads",
				parsedfile.Packages[i].PkgVersion, parsedfile.Packages[i].Downloads)
		}

	}

	/**
	 * Get the INSTALLS count
	 */
	log.Println("Starting install count for given repos")
	for i := range parsedfile.Repos {
		getInstalls(&parsedfile.Repos[i], DEFAULT_START_DATEZ, parsedfile.User, parsedfile.UserRepo)
	}

	/**
	 * Print results
	 */
	// downloads for packages
	fmt.Println("DOWNLOADS FOR ", parsedfile.Package)
	for _, pkg := range parsedfile.Packages {
		fmt.Println(pkg)
	}

	// installs for repos
	fmt.Println("INSTALLS")
	for _, repo := range parsedfile.Repos {
		fmt.Println(repo)
	}

} // end: main()
