# PkgCloudStats
Connects to Packagecloud.io API to deliver stats on packages and repos. This was created to help a OSS company track the number of specific repos and packages being installed because the built in stats that show do not make a clear distinction. 

## How to use
(make sure you have go already installed lol)
From command line - 
1. edit the JSON file to your needs. Keep the same structure. Just change the values. *Note- You may set the API Token to "ENV" where the program will pull the token from the OS enviornment variable labeled "PKGCLOUD_API_TOKEN"* 
2. build main.go to make it an executable
```
go. build main.go
```
3. Finally, run ```./main```

## Cmmnd Line Flag Options and Output
You have 2 flag at your disposal: 
1. -debug : if set to true, more specific log statements would be outputted. Default is false
2. -config : Path set to json config file.  Default set to "Pkgcloud-Counter-config.json" assumed to be in same directory as main.go. 

The number of downloads for each package will output to the console for each version shown respectivley. 
The number of installs will also output to console for each repo
### Sample output
2020/08/31 12:47:12 Set debug to  false

2020/08/31 12:47:12 Set file to  Pkgcloud-Counter-config.json

2020/08/31 12:47:12 Starting download count for given packages

2020/08/31 12:47:15 Starting install count for given repos

DOWNLOADS FOR  sensu-go-backend

{6.0.0 3003 el [6 7 8] [173 1061 154]}

{5.21.0 14262 el [6 7 8] [311 1745 296]}

{5.20.2 12959 el [6 7 8] [121 819 123]}

INSTALLS

{el 6 155041}

{el 7 38130}

{el 8 908}
