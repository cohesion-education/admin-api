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

	authConfig, err := config.NewAuthConfig()
	if err != nil {
		log.Fatal(err)
	}

	gcpKeyfileLocation := os.Getenv("GCP_KEYFILE_LOCATION")
	if len(gcpKeyfileLocation) == 0 {
		log.Fatal("Required env var GCP_KEYFILE_LOCATION not set")
	}

	videoStorageBucketName := os.Getenv("GCP_STORAGE_VIDEO_BUCKET")
	if len(videoStorageBucketName) == 0 {
		log.Fatal("Required env var GCP_STORAGE_VIDEO_BUCKET not set")
	}

	gcpConfig, err := gcp.NewConfig(gcpKeyfileLocation)
	if err != nil {
		log.Fatal(err)
	}

	ctx := context.TODO()
	datastoreClient, err := gcp.NewDatastoreClient(ctx, gcpConfig)
	if err != nil {
		log.Fatalf("Failed to get datastore client %v", err)
	}

	storageClient, err := gcp.NewStorageClient(ctx, gcpConfig)
	if err != nil {
		log.Fatalf("Failed to get storage client %v", err)
	}

	taxonomyRepo := taxonomy.NewGCPDatastoreRepo(datastoreClient)
	videoRepo := video.NewGCPRepo(datastoreClient, storageClient, videoStorageBucketName)

	n := negroni.Classic()
	mx := mux.NewRouter()

	// Static Assets
	mx.PathPrefix("/assets/").Handler(http.StripPrefix("/assets/", http.FileServer(http.Dir("./assets/"))))

	//Public Routes
	mx.Methods("GET").Path("/").Handler(auth.LoginViewHandler(renderer))
	mx.Methods("GET").Path("/admin/login").Handler(auth.LoginViewHandler(renderer))
	mx.Methods("GET").Path("/logout").Handler(auth.LogoutHandler(renderer))
	mx.Methods("GET").Path("/callback").Handler(auth.CallbackHandler(authConfig))
	mx.Methods("GET").Path("/auth/config").Handler(auth.ConfigHandler(authConfig))

	isAuthenticatedHandler := auth.IsAuthenticatedHandler(authConfig)
	commonMiddleware := negroni.New(
		negroni.HandlerFunc(isAuthenticatedHandler),
	)

	//Admin Routes
	adminSubRouter := mux.NewRouter()
	adminSubRouter.Methods("GET").Path("/admin/dashboard").Handler(admin.DashboardViewHandler(renderer))
	adminSubRouter.Methods("GET").Path("/admin/taxonomy").Handler(taxonomy.ListHandler(renderer, taxonomyRepo))
	adminSubRouter.Methods("GET").Path("/admin/video").Handler(video.ListHandler(renderer, videoRepo))
	adminSubRouter.Methods("GET").Path("/admin/video/{id:[0-9]+}").Handler(video.ShowHandler(renderer, videoRepo))
	adminSubRouter.Methods("GET").Path("/admin/video/add").Handler(video.FormHandler(renderer, videoRepo))
	mx.PathPrefix("/admin").Handler(commonMiddleware.With(
		negroni.HandlerFunc(auth.IsAdmin),
		negroni.Wrap(adminSubRouter),
	))

	//APIs that require Admin priveleges
	requireAdmin("POST", "/api/taxonomy", taxonomy.AddHandler(renderer, taxonomyRepo), mx, commonMiddleware)
	requireAdmin("GET", "/api/taxonomy", taxonomy.ListJSONHandler(renderer, taxonomyRepo), mx, commonMiddleware)
	requireAdmin("GET", "/api/taxonomy/{id:[0-9]+}/children", taxonomy.ListChildrenHandler(renderer, taxonomyRepo), mx, commonMiddleware)
	requireAdmin("POST", "/api/video", video.SaveHandler(renderer, videoRepo), mx, commonMiddleware)
	requireAdmin("GET", "/api/video/stream/{id:[0-9]+}", video.StreamHandler(renderer, videoRepo, gcpConfig), mx, commonMiddleware)

	n.UseHandler(mx)
	return n
}

func requireAdmin(method string, uri string, handler http.Handler, mx *mux.Router, commonMiddleware *negroni.Negroni) {
	mx.Methods(method).Path(uri).Handler(commonMiddleware.With(
		negroni.HandlerFunc(auth.IsAdmin),
		negroni.Wrap(handler),
	))
}
