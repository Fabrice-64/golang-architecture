package main

import (
	"encoding/base64"
	"fmt"
)

// curl -u user:pass -v google.com
// Will produce this as outcome, user and pass in base64
//Authorization: Basic dXNlcjpwYXNz

func main() {
	p := base64.StdEncoding.EncodeToString([]byte("user:pass"))
	fmt.Println("Encoded String: ", p)
}

// the output will be the same string.
