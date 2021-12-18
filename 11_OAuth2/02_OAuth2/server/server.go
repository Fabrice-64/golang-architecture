package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/http/httputil"
	"net/url"
	"os"
	"time"

	"github.com/go-oauth2/oauth2/v4/errors"
	"github.com/go-oauth2/oauth2/v4/generates"
	"github.com/go-oauth2/oauth2/v4/manage"
	"github.com/go-oauth2/oauth2/v4/models"
	"github.com/go-oauth2/oauth2/v4/server"
	"github.com/go-oauth2/oauth2/v4/store"
	"github.com/go-session/session"
)

var (
	dumpvar   bool
	idvar     string
	secretvar string
	domainvar string
	portvar   int
)

func init() {
	flag.BoolVar(&dumpvar, "d", true, "Dump requests and responses")
	flag.StringVar(&idvar, "i", "222222", "The client ID is being passed in")
	flag.StringVar(&secretvar, "str", "ABDCDEF", "The client secret is being passed in")
	flag.StringVar(&domainvar, "r", "http://localhost:9096", "The domain of the redirect url")
	flag.IntVar(&portvar, "p", 9096, "The base port for the server")
}

func main() {
	flag.Parse()
	if dumpvar {
		log.Println("Dump client requests")
	}
	manager := manage.NewDefaultManager()
	manager.SetAuthorizeCodeTokenCfg(manage.DefaultAuthorizeCodeTokenCfg)
	manager.MustTokenStorage(store.NewMemoryTokenStore())
	manager.MapAccessGenerate(generates.NewAccessGenerate())

	clientStore := store.NewClientStore()
	clientStore.Set(idvar, &models.Client{
		ID:     idvar,
		Secret: secretvar,
		Domain: domainvar,
	})
	manager.MapClientStorage(clientStore)

	// Create authorization server
	srv := server.NewServer(server.NewConfig(), manager)
	srv.SetPasswordAuthorizationHandler(func(un, pwd string) (userID string, err error) {
		if un == "test" && pwd == "test" {
			userID = "test"
		}
		return
	})
	srv.SetUserAuthorizationHandler(userAuthorizeHandler)
	srv.SetInternalErrorHandler(func(err error) (re *errors.Response) {
		log.Println("response error: ", re.Error.Error())
		return
	})
	srv.SetResponseErrorHandler(func(re *errors.Response) {
		log.Println(re.Error.Error())
	})
	http.HandleFunc("/login", loginHandler)
	http.HandleFunc("/auth", authHandler)

	http.HandleFunc("/oauth/authorize", func(w http.ResponseWriter, req *http.Request) {
		if dumpvar {
			dumpRequest(os.Stdout, "authorize", req)
		}
		store, err := session.Start(req.Context(), w, req)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		var form url.Values
		if v, ok := store.Get("ReturnUri"); ok {
			form = v.(url.Values)
		}
		req.Form = form
		store.Delete("ReturnUri")
		store.Save()
		err = srv.HandleAuthorizeRequest(w, req)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
		}
	})
	http.HandleFunc("oauth/token", func(w http.ResponseWriter, req *http.Request) {
		if dumpvar {
			dumpRequest(os.Stdout, "token", req) //ignore the error
		}
		err := srv.HandleTokenRequest(w, req)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
	})
	http.HandleFunc("/test", func(w http.ResponseWriter, req *http.Request) {
		if dumpvar {
			dumpRequest(os.Stdout, "test", req)
		}
		token, err := srv.ValidationBearerToken(req)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		data := map[string]interface{}{
			"expires_in": int64(token.GetAccessCreateAt().Add(token.GetAccessExpiresIn()).Sub(time.Now()).Seconds()),
			"client_id":  token.GetClientID(),
			"user_id":    token.GetUserID(),
		}
		e := json.NewEncoder(w)
		e.SetIndent("", " ")
		e.Encode(data)
	})
	log.Printf("Server is running at %d port.\n", portvar)
	log.Printf("Point your OAuth client Auth at %s:%d%s", "http://localhost", portvar, "/oauth/authorize")
	log.Printf("Point your OAuth client Token at %s:%d%s", "http://localhost", portvar, "/oauth/token")
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%d", portvar), nil))
}

func dumpRequest(w io.Writer, header string, req *http.Request) error {
	data, err := httputil.DumpRequest(req, true)
	if err != nil {
		return err
	}
	w.Write([]byte("\n" + header + ":\n"))
	w.Write(data)
	return nil
}

func userAuthorizeHandler(w http.ResponseWriter, req *http.Request) (userID string, err error) {
	if dumpvar {
		_ = dumpRequest(os.Stdout, "userAuthorizeHandler", req)
	}
	store, err := session.Start(req.Context(), w, req)
	if err != nil {
		return
	}
	uid, ok := store.Get("LoggedInUserID")
	if !ok {
		if req.Form == nil {
			req.ParseForm()
		}
		store.Set("Client URI", req.Form)
		store.Save()
		w.Header().Set("Location", "/login")
		w.WriteHeader(http.StatusFound)
		return
	}
	userID = uid.(string)
	store.Delete("LoggedInUserID")
	store.Save()
	return
}

func loginHandler(w http.ResponseWriter, req *http.Request) {
	if dumpvar {
		_ = dumpRequest(os.Stdout, "login", req) // ignore the error
	}
	store, err := session.Start(req.Context(), w, req)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	if req.Method == "POST" {
		if err := req.ParseForm(); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		store.Set("LoggedInUserID", req.Form.Get("username"))
		store.Save()
		w.Header().Set("Location", "/auth")
		w.WriteHeader(http.StatusFound)
		return
	}
	outputHTML(w, req, "static/login.html")
}

func outputHTML(w http.ResponseWriter, req *http.Request, filename string) {
	file, err := os.Open(filename)
	if err != nil {
		http.Error(w, "Internal Error", http.StatusInternalServerError)
		return
	}
	defer file.Close()
	fi, _ := file.Stat()
	http.ServeContent(w, req, file.Name(), fi.ModTime(), file)
}

func authHandler(w http.ResponseWriter, req *http.Request) {
	if dumpvar {
		_ = dumpRequest(os.Stdout, "auth", req) //ignore the error
	}
	store, err := session.Start(nil, w, req)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	if _, ok := store.Get("LoggedInUserID"); !ok {
		w.Header().Set("Location", "/login")
		w.WriteHeader(http.StatusFound)
		return
	}
	outputHTML(w, req, "static/auth.html")
}
