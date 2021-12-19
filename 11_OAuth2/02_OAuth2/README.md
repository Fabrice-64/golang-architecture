# This exercise combines client and server using OAuth2 protocol

Server and Client are to be run together.

# Server
## Configure
`$ cd ../server`
`$ go build server.go`
`$ ./server`

You don't need the client to launch it ;-

## Outcome
`2021/12/18 09:34:41 Dump client requests`
`2021/12/18 09:34:41 Server is running at 9096 port.`
`2021/12/18 09:34:41 Point your OAuth client Auth at http://localhost:9096/oauth/authorize`
`2021/12/18 09:34:41 Point your OAuth client Token at http://localhost:9096/oauth/token`

# Client
## Configure
`$ cd ../client`
`$ go build client.go`
`$ ./client`

## Outcome
The user can then start using the client & server:
`http:localhost:9094/try`
`http:localhost:9094/refresh`
`http:localhost:9094/pwd`
`http:localhost:9094/login`
