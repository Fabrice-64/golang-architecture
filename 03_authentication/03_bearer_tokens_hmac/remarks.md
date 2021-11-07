# Bearer Tokens & Hmac

## bearer tokens
added to http spec with OAUTH2
uses authorization header & keyword “bearer”
to prevent faked bearer tokens, use cryptographic “signing”
cryptographic signing is a way to prove that the value was created/validated by a certain person
HMAC
https://godoc.org/crypto/hmac
