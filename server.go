package main

import (
	"net/http"

	mgo "gopkg.in/mgo.v2"

	"github.com/gorilla/mux"
	"github.com/unrolled/render"
	"github.com/urfave/negroni"
)

type key int

const (
	currentUserKey key = iota
)

type handlerConfig struct {
	renderer     *render.Render
	mongoSession *mgo.Session
}

func newHandlerConfig() *handlerConfig {
	return &handlerConfig{
		mongoSession: newMongoSession(),
		renderer: render.New(render.Options{
			Layout: "layout",
			RenderPartialsWithoutPrefix: true,
		}),
	}
}

func newServer() *negroni.Negroni {
	n := negroni.Classic()
	mx := mux.NewRouter()

	handlerConfig := newHandlerConfig()
	authConfig := newAuthConfig()

	// This will serve files under /assets/<filename>
	mx.PathPrefix("/assets/").Handler(http.StripPrefix("/assets/", http.FileServer(http.Dir("./assets/"))))

	mx.HandleFunc("/", loginViewHandler(handlerConfig)).Methods("GET")
	mx.HandleFunc("/logout", logoutHandler(handlerConfig)).Methods("GET")
	mx.Handle("/callback", callbackHandler(authConfig)).Methods("GET")

	isAuthenticatedHandler := isAuthenticatedHandler(authConfig)
	mx.Handle("/dashboard", secure(isAuthenticatedHandler, dashboardHandler(handlerConfig))).Methods("GET")
	mx.Handle("/taxonomy", secure(isAuthenticatedHandler, taxonomyListHandler(handlerConfig))).Methods("GET")

	n.UseHandler(mx)

	return n
}

func secure(handlerWithNext negroni.HandlerFunc, handlerFunc http.HandlerFunc) *negroni.Negroni {
	return negroni.New(
		negroni.HandlerFunc(handlerWithNext),
		negroni.Wrap(handlerFunc),
	)
}
