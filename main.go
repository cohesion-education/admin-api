package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"

	jwtmiddleware "github.com/auth0/go-jwt-middleware"
	"github.com/cloudfoundry-community/go-cfenv"
	"github.com/codegangsta/negroni"
	jwt "github.com/dgrijalva/jwt-go"
	"github.com/gorilla/mux"
	"github.com/joho/godotenv"
)

func main() {
	var auth0ClientSecret string

	if appEnv, err := cfenv.Current(); err == nil {
		if auth0Service, err := appEnv.Services.WithName("admin-api-auth0"); err == nil {
			if clientSecret, ok := auth0Service.CredentialString("secret"); ok {
				auth0ClientSecret = clientSecret
			}
		}
	} else {
		if err := godotenv.Load(); err == nil {
			auth0ClientSecret = os.Getenv("AUTH0_CLIENT_SECRET")
		}
	}

	if len(auth0ClientSecret) == 0 {
		log.Fatal("Failed to load auth0 client secret from VCAP_SERVICES or from .env file")
	}

	port := os.Getenv("PORT")
	if len(port) == 0 {
		port = "3001"
	}

	startServer(port, auth0ClientSecret)
}

func startServer(port string, auth0ClientSecret string) {
	r := mux.NewRouter()

	jwtMiddleware := jwtmiddleware.New(jwtmiddleware.Options{
		ValidationKeyGetter: func(token *jwt.Token) (interface{}, error) {
			secret := []byte(auth0ClientSecret)
			return secret, nil
		},
	})

	r.HandleFunc("/ping", pingHandler)
	r.Handle("/secured/ping", negroni.New(
		negroni.HandlerFunc(jwtMiddleware.HandlerWithNext),
		negroni.Wrap(http.HandlerFunc(securedPingHandler)),
	))
	http.Handle("/", r)

	fmt.Print("Server started and listening on ", port)
	http.ListenAndServe(":"+port, nil)
}

type response struct {
	Text string `json:"text"`
}

func respondJSON(text string, w http.ResponseWriter) {
	r := response{text}

	jsonResponse, err := json.Marshal(r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(jsonResponse)
}

func pingHandler(w http.ResponseWriter, r *http.Request) {
	respondJSON("All good. You don't need to be authenticated to call this", w)
}

func securedPingHandler(w http.ResponseWriter, r *http.Request) {
	respondJSON("All good. You only get this message if you're authenticated", w)
}
