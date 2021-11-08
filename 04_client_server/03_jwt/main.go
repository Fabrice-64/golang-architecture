package main

import (
	"fmt"
	"time"

	"github.com/golang-jwt/jwt"
)

type UserClaims struct {
	jwt.StandardClaims
	sessionID int64
}

func (u *UserClaims) Valid() error {
	if !u.VerifyExpiresAt(time.Now().Unix(), true) {
		return fmt.Errorf("Token has Expired")
	}
	if u.sessionID == 0 {
		return fmt.Errorf("Invalid Session")
	}
	return nil
}

func main() {

}
