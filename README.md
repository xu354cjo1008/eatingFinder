# eatingFinder
Library for weather information crawl and parsing
## How to use:
go get github.com/xu354cjo1008/eatingFiner
### Set your api keys and http port in config/app.toml
for example:  
[development]  
apiHost = "localhost"  
apiPort = 8080  
googleApiKey = "abcdefgh"  
cwdApiKey = "12345678"  
###Run api server
./eatingFinder -mode api -port <port number>  
###Run web server
configure api server host name and port number  
./eatingFinder -mode web -port <port number>  
