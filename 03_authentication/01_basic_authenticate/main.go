package main

import (
	"encoding/base64"
	"fmt"
)

// curl -u user:pass -v google.com
// Will produce the following outcome, user:pass in base64.
// Therefore both username and password have to be used associated with :
//Authorization: Basic dXNlcjpwYXNz
// To be noticed: http sends back a line with "Authorization" although it is an authentication

func main() {
	p := base64.StdEncoding.EncodeToString([]byte("user:pass"))
	fmt.Println("Encoded String: ", p)
}

// the output will be the same string.
// base64 is sent at each exchange but is not secure. To use only with https.
