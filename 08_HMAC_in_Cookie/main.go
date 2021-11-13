package main

import (
	"crypto/hmac"
	"crypto/sha256"
	"fmt"
	"io"
	"net/http"
	"strings"
)

func main() {
	http.HandleFunc("/", foo)
	http.HandleFunc("/submit", bar)
	http.ListenAndServe(":8080", nil)
}

func getCode(msg string) string {
	h := hmac.New(sha256.New, []byte("This is the secret sentence !"))
	h.Write([]byte(msg))
	return fmt.Sprintf("%x", h.Sum(nil))

}

func bar(w http.ResponseWriter, req *http.Request) {
	if req.Method != http.MethodPost {
		http.Redirect(w, req, "/", http.StatusSeeOther)
		return
	}
	email := req.FormValue("email")
	if email == "" {
		http.Redirect(w, req, "/", http.StatusSeeOther)
		return
	}
	code := getCode(email)
	c := http.Cookie{
		Name:  "crypto-course",
		Value: code + "|" + email,
	}
	http.SetCookie(w, &c)
	http.Redirect(w, req, "/", http.StatusSeeOther)

}

func foo(w http.ResponseWriter, req *http.Request) {
	c, err := req.Cookie("crypto-course")
	if err != nil {
		c = &http.Cookie{}
	}

	isEqual := true
	xs := strings.SplitN(c.Value, "|", 2)
	if len(xs) == 2 {
		cCode := xs[0]
		cEmail := xs[1]
		code := getCode(cEmail)
		isEqual = hmac.Equal([]byte(cCode), []byte(code))
	}
	msg := "Not Logged In"
	if isEqual {
		msg = "Logged In"
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
