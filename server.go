package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/cohesion-education/admin-api/pkg/cohesioned/admin"
	"github.com/cohesion-education/admin-api/pkg/cohesioned/auth"
	"github.com/cohesion-education/admin-api/pkg/cohesioned/config"
	"github.com/cohesion-education/admin-api/pkg/cohesioned/gcp"
	"github.com/cohesion-education/admin-api/pkg/cohesioned/taxonomy"
	"github.com/cohesion-education/admin-api/pkg/cohesioned/video"
	"github.com/gorilla/mux"
	"github.com/joho/godotenv"
	"github.com/unrolled/render"
	"github.com/urfave/negroni"
)

var (
	renderer = render.New(render.Options{
		Layout: "layout",
		RenderPartialsWithoutPrefix: true,
	})
)

func newServer() *negroni.Negroni {
	if err := godotenv.Load(); err != nil {
		fmt.Println("Failed to load .env file. Components will fallback to loading from VCAP_SERVICES or env vars")
	}

	n := negroni.Classic()
	mx := mux.NewRouter()

	authConfig, err := config.NewAuthConfig()
	if err != nil {
		log.Fatal(err)
	}

	ctx := context.TODO()
	gcpProjectID := os.Getenv("DATASTORE_PROJECT_ID")
	datastoreClient, err := gcp.NewDatastoreClient(ctx, gcpProjectID)
	if err != nil {
		log.Fatalf("Failed to get datastore client %v", err)
	}

	videoStorageBucketName := os.Getenv("GCP_STORAGE_VIDEO_BUCKET")
	storageClient, err := gcp.NewStorageClient(ctx)
	if err != nil {
		log.Fatalf("Failed to get storage client %v", err)
	}

	taxonomyRepo := taxonomy.NewGCPDatastoreRepo(datastoreClient)
	videoRepo := video.NewGCPRepo(datastoreClient, storageClient, videoStorageBucketName)

	// This will serve files under /assets/<filename>
	mx.PathPrefix("/assets/").Handler(http.StripPrefix("/assets/", http.FileServer(http.Dir("./assets/"))))

	mx.HandleFunc("/", auth.LoginViewHandler(renderer)).Methods("GET")
	mx.HandleFunc("/logout", auth.LogoutHandler(renderer)).Methods("GET")
	mx.Handle("/callback", auth.CallbackHandler(authConfig)).Methods("GET")

	isAuthenticatedHandler := auth.IsAuthenticatedHandler(authConfig)
	mx.Handle("/admin/dashboard", secure(isAuthenticatedHandler, admin.DashboardViewHandler(renderer))).Methods("GET")
	mx.Handle("/admin/taxonomy", secure(isAuthenticatedHandler, taxonomy.ListHandler(renderer, taxonomyRepo))).Methods("GET")
	mx.Handle("/admin/video", secure(isAuthenticatedHandler, video.ListHandler(renderer, videoRepo))).Methods("GET")
	mx.Handle("/admin/video/{id:[0-9]+}", secure(isAuthenticatedHandler, video.ShowHandler(renderer, videoRepo))).Methods("GET")
	mx.Handle("/admin/video/add", secure(isAuthenticatedHandler, video.FormHandler(renderer, videoRepo))).Methods("GET")
	mx.Handle("/api/taxonomy", secure(isAuthenticatedHandler, taxonomy.AddHandler(renderer, taxonomyRepo))).Methods("POST")
	mx.Handle("/api/taxonomy", secure(isAuthenticatedHandler, taxonomy.ListJSONHandler(renderer, taxonomyRepo))).Methods("GET")
	mx.Handle("/api/taxonomy/{id:[0-9]+}/children", secure(isAuthenticatedHandler, taxonomy.ListChildrenHandler(renderer, taxonomyRepo))).Methods("GET")
	mx.Handle("/api/video", secure(isAuthenticatedHandler, video.SaveHandler(renderer, videoRepo))).Methods("POST")
	mx.Handle("/api/video/stream/{id:[0-9]+}", secure(isAuthenticatedHandler, video.StreamHandler(renderer, videoRepo))).Methods("GET")

	n.UseHandler(mx)

	return n
}

func secure(handlerWithNext negroni.HandlerFunc, handlerFunc http.HandlerFunc) *negroni.Negroni {
	return negroni.New(
		negroni.HandlerFunc(handlerWithNext),
		negroni.Wrap(handlerFunc),
	)
}
