package cohesioned

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/cohesion-education/admin-api/pkg/admin"
	"github.com/cohesion-education/admin-api/pkg/auth"
	"github.com/cohesion-education/admin-api/pkg/config"
	"github.com/cohesion-education/admin-api/pkg/gcp"
	"github.com/cohesion-education/admin-api/pkg/taxonomy"
	"github.com/gorilla/mux"
	"github.com/joho/godotenv"
	"github.com/urfave/negroni"
)

func newServer() *negroni.Negroni {
	if err := godotenv.Load(); err != nil {
		fmt.Println("Failed to load .env file. Components will fallback to loading from VCAP_SERVICES or env vars")
	}

	n := negroni.Classic()
	mx := mux.NewRouter()

	ctx := context.TODO()
	gcpProjectID := os.Getenv("DATASTORE_PROJECT_ID")
	datastoreClient, err := gcp.NewDatastoreClient(ctx, gcpProjectID)
	if err != nil {
		log.Fatal(err)
	}

	handlerConfig, err := config.NewHandlerConfig(datastoreClient)
	if err != nil {
		log.Fatal(err)
	}

	authConfig, err := config.NewAuthConfig()
	if err != nil {
		log.Fatal(err)
	}

	// This will serve files under /assets/<filename>
	mx.PathPrefix("/assets/").Handler(http.StripPrefix("/assets/", http.FileServer(http.Dir("./assets/"))))

	mx.HandleFunc("/", auth.LoginViewHandler(handlerConfig)).Methods("GET")
	mx.HandleFunc("/logout", auth.LogoutHandler(handlerConfig)).Methods("GET")
	mx.Handle("/callback", auth.CallbackHandler(authConfig)).Methods("GET")

	isAuthenticatedHandler := auth.IsAuthenticatedHandler(authConfig)
	mx.Handle("/dashboard", secure(isAuthenticatedHandler, admin.DashboardViewHandler(handlerConfig))).Methods("GET")
	mx.Handle("/taxonomy", secure(isAuthenticatedHandler, taxonomy.ListHandler(handlerConfig))).Methods("GET")
	mx.Handle("/api/taxonomy", secure(isAuthenticatedHandler, taxonomy.AddHandler(handlerConfig))).Methods("POST")

	n.UseHandler(mx)

	return n
}

func secure(handlerWithNext negroni.HandlerFunc, handlerFunc http.HandlerFunc) *negroni.Negroni {
	return negroni.New(
		negroni.HandlerFunc(handlerWithNext),
		negroni.Wrap(handlerFunc),
	)
}
