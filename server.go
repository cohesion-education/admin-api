package main

import (
	"encoding/gob"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/gorilla/sessions"
	"github.com/unrolled/render"
	"github.com/urfave/negroni"
)

var (
	store sessions.Store
)

func newServer() *negroni.Negroni {
	store = sessions.NewCookieStore([]byte("todo-inject-me"))
	gob.Register(map[string]interface{}{})

	config := newAuth0Config()

	// mongoDialInfo := &mgo.DialInfo{
	// 	Addrs:    []string{"127.0.0.1"},
	// 	Database: "cohesion-education",
	// 	Username: "admin",
	// 	Password: "password",
	// 	Source:   "admin",
	// 	Timeout:  60 * time.Second,
	// }
	// session, err := mgo.DialWithInfo(mongoDialInfo)
	//
	// //session, err := mgo.Dial(" mongodb://admin:password@127.0.0.1:27017/cohesion-education?authSource=admin")
	// if err != nil {
	// 	log.Fatalf("Failed to connect to mongodb: %v", err)
	// }
	//
	// taxonomyRepo := &taxonomyRepository{session: session}

	r := render.New(render.Options{
		Layout: "layout",
		RenderPartialsWithoutPrefix: true,
	})
	n := negroni.Classic()
	mx := mux.NewRouter()

	// This will serve files under /assets/<filename>
	mx.PathPrefix("/assets/").Handler(http.StripPrefix("/assets/", http.FileServer(http.Dir("./templates/ubold/assets/"))))

	mx.HandleFunc("/", loginViewHandler(r, store)).Methods("GET")
	mx.HandleFunc("/logout", logoutHandler(r, store)).Methods("GET")
	mx.Handle("/callback", callbackHandler(config, store)).Methods("GET")
	mx.Handle("/dashboard", secure(dashboardHandler(r, store))).Methods("GET")
	mx.Handle("/users", secure(userListHandler(r, store))).Methods("GET")
	mx.Handle("/taxonomy", secure(taxonomyListHandler(r, nil))).Methods("GET")

	n.UseHandler(mx)

	return n
}

func secure(handlerFunc http.HandlerFunc) *negroni.Negroni {
	return negroni.New(
		negroni.HandlerFunc(isAuthenticatedHandler(store)),
		negroni.Wrap(handlerFunc),
	)
}
