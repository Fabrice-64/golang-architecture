package main

import (
	"io"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/joho/godotenv"
)

type myClaims struct {
	jwt.StandardClaims
}

func main() {
	http.HandleFunc("/", index)
	http.HandleFunc("/submit", bar)
	http.ListenAndServe(":8080", nil)
}

func jwtCookie(email string) (string, error) {
	envVar := getEnvVar("SECRET_KEY")
	myKey := []byte(envVar)
	claims := &myClaims{
		jwt.StandardClaims{
			ExpiresAt: time.Now().Add(15 * time.Minute).Unix(),
			Issuer:    email,
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	signedString, err := token.SignedString(myKey)
	if err != nil {
		log.Fatalln("Error While Signing the Token: ", err)
		return "", err
	}
	return signedString, nil
}

func index(w http.ResponseWriter, req *http.Request) {
	c, err := req.Cookie("jwt-cookie")
	if err != nil {
		c = &http.Cookie{
			Name:  "jwt-cookie",
			Value: "NIL",
		}
	}
	// new
	msg := ""
	if c.Value == "NIL" {
		msg = "Not Logged In"
	} else {
		msg = "Logged In"
	}
	envVar := getEnvVar("SECRET_KEY")
	signedString := c.Value

	// sample token is expired.  override time so it parses as valid
	token, err := jwt.ParseWithClaims(signedString, &myClaims{}, func(token *jwt.Token) (interface{}, error) {
		return []byte(envVar), nil
	})
	if err != nil {
		msg = "Not Logged In or Token Invalid"
		log.Printf("Token Invalid : %v", err)
	}
	if err == nil && token.Valid {
		msg = "Logged In and Token is Valid"
	}

	html := `<!doctype html>
		<html>
		<head>
	<title>HMAC and Cookies</title>
	<meta name="description" content="HMAC and Cookies">
	<meta name="keywords" content="hmac cookies ">
	</head>
	<body>
		<h1>Form to fill</h1>
		<h2>` + msg + `</h2>
		<p>Cookie Value: ` + c.Value + `</p>
		<form action="/submit" method="post">
			<label for="email">E-Mail</label>
			<input type="email" name="email" id="email" placeholder="email">
			<input type="submit" value="Valider">
		</form>
	</body>
	</html>`
	io.WriteString(w, html)

}
func bar(w http.ResponseWriter, req *http.Request) {
	if req.Method != http.MethodPost {
		http.Redirect(w, req, "/", http.StatusSeeOther)
		log.Printf("Method was not Post: %x", http.StatusMethodNotAllowed)
		return
	}

	email := req.FormValue("email")
	if email == "" {
		http.Redirect(w, req, "/", http.StatusSeeOther)
	}
	log.Printf("Value: %s", email)
	signedString, err := jwtCookie(email)
	if err != nil {
		http.Error(w, "Could not get a JWT Token", http.StatusInternalServerError)
		return
	}
	//log.Println("Token Signé: ", signedString)
	c := http.Cookie{
		Name:  "jwt-cookie",
		Value: signedString,
	}
	http.SetCookie(w, &c)
	http.Redirect(w, req, "/", http.StatusSeeOther)

}

func getEnvVar(key string) string {
	err := godotenv.Load(".env")
	if err != nil {
		log.Panicln("could not get access to env var ", err)
	}
	return os.Getenv(key)
}
