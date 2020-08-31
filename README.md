# PkgCloudStats
Connects to Packagecloud.io API to deliver stats on packages and repos. This was created to help a OSS company track the number of specific repos and packages being installed because the built in stats that show do not make a clear distinction. 

## How to use
(make sure you have go already installed lol)
From command line - 
1. edit the JSON file to your needs. Keep the same structure. Just change the values. 
2. build main.go to make it an executable
```
go. build main.go
```
3. ```./main```

## Cmmnd Line Flag Options and Output
You have 2 flag at your disposale: 
1. -debug : if set to true, more specific log statements would be outputted. Default is false
2. -config : Path set to json config file.  Default set to "Pkgcloud-Counter-config.json" assumed to be in same directory as main.go

The number of downloads for each package will output to the console for each version shown respectivley. 
The number of installs will also output to console for each repo
