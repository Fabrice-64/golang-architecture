# Commands to start the server and test it.
1. `go build server.go`
2. `./server`
3. Authorization Request : `http://localhost:9096/authorize?client_id=000000&response_type=code`
4. Grant Token Request: `http://localhost:9096/token?grant_type=client_credentials&client_id=000000&client_secret=999999&scope=read`

