package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/cohesion-education/admin-api/pkg/cohesioned"
	"github.com/cohesion-education/admin-api/pkg/cohesioned/auth"
	"github.com/cohesion-education/admin-api/pkg/cohesioned/config"
	"github.com/cohesion-education/admin-api/pkg/cohesioned/dashboard"
	"github.com/cohesion-education/admin-api/pkg/cohesioned/gcp"
	"github.com/cohesion-education/admin-api/pkg/cohesioned/homepage"
	"github.com/cohesion-education/admin-api/pkg/cohesioned/taxonomy"
	"github.com/cohesion-education/admin-api/pkg/cohesioned/video"
	"github.com/gorilla/mux"
	"github.com/joho/godotenv"
	"github.com/unrolled/render"
	"github.com/urfave/negroni"
)

var (
	apiRenderer = render.New()

	homepageRenderer = render.New(render.Options{
		Layout: "homepage/layout",
		RenderPartialsWithoutPrefix: true,
	})

	adminDashboardRenderer = render.New(render.Options{
		Layout: "admin-layout",
		RenderPartialsWithoutPrefix: true,
	})

	userDashboardRenderer = render.New(render.Options{
		Layout: "user-layout",
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

	homepageRepo := homepage.NewGCPDatastoreRepo(datastoreClient)
	taxonomyRepo := taxonomy.NewGCPDatastoreRepo(datastoreClient)
	videoRepo := video.NewGCPRepo(datastoreClient, storageClient, videoStorageBucketName)

	n := negroni.Classic()
	mx := mux.NewRouter()
	mx.StrictSlash(true)

	// Static Assets
	mx.PathPrefix("/assets/").Handler(http.StripPrefix("/assets/", http.FileServer(http.Dir("./assets/"))))

	mx.NotFoundHandler = cohesioned.NotFoundViewHandler(homepageRenderer)
	mx.Methods(http.MethodGet).Path("/401").Handler(cohesioned.UnauthorizedViewHandler(homepageRenderer))
	mx.Methods(http.MethodGet).Path("/403").Handler(cohesioned.ForbiddenViewHandler(homepageRenderer))
	mx.Methods(http.MethodGet).Path("/500").Handler(cohesioned.InternalServerErrorViewHandler(homepageRenderer))

	//Public Routes
	mx.Methods(http.MethodGet).Path("/").Handler(homepage.HomepageViewHandler(homepageRenderer, homepageRepo))
	mx.Methods(http.MethodGet).Path("/logout").Handler(auth.LogoutHandler(authConfig))
	mx.Methods(http.MethodGet).Path("/callback").Handler(auth.CallbackHandler(authConfig))
	mx.Methods(http.MethodGet).Path("/auth/config").Handler(auth.ConfigHandler(authConfig))

	//Public APIs
	mx.Methods(http.MethodGet).Path("/api/taxonomy/flatten").Handler(taxonomy.FlatListHandler(apiRenderer, taxonomyRepo))

	isAuthenticatedHandler := auth.IsAuthenticatedHandler(authConfig)
	authMiddleware := negroni.New(
		negroni.HandlerFunc(isAuthenticatedHandler),
	)

	//User Routes
	requiresAuth(http.MethodGet, "/dashboard", dashboard.UserViewHandler(userDashboardRenderer), mx, authMiddleware)

	//Admin Routes
	adminRouter := mux.NewRouter()
	adminRouter.Methods(http.MethodGet).Path("/admin/dashboard").Handler(dashboard.AdminViewHandler(adminDashboardRenderer))
	adminRouter.Methods(http.MethodGet).Path("/admin/homepage").Handler(homepage.FormViewHandler(adminDashboardRenderer, homepageRepo))
	adminRouter.Methods(http.MethodPost).Path("/admin/homepage").Handler(homepage.SaveHandler(adminDashboardRenderer, homepageRepo))
	adminRouter.Methods(http.MethodGet).Path("/admin/taxonomy").Handler(taxonomy.ListViewHandler(adminDashboardRenderer, taxonomyRepo))
	adminRouter.Methods(http.MethodGet).Path("/admin/video").Handler(video.ListViewHandler(adminDashboardRenderer, videoRepo))
	adminRouter.Methods(http.MethodGet).Path("/admin/video/{id:[0-9]+}").Handler(video.ShowViewHandler(adminDashboardRenderer, videoRepo))
	adminRouter.Methods(http.MethodGet).Path("/admin/video/add").Handler(video.FormViewHandler(adminDashboardRenderer, videoRepo))
	adminRouter.Methods(http.MethodGet).Path("/admin/video/edit/{id:[0-9]+}").Handler(video.FormViewHandler(adminDashboardRenderer, videoRepo))
	mx.PathPrefix("/admin").Handler(authMiddleware.With(
		negroni.HandlerFunc(auth.IsAdmin),
		negroni.Wrap(adminRouter),
	))

	//APIs that require Admin priveleges
	requiresAdmin(http.MethodPost, "/api/taxonomy", taxonomy.AddHandler(apiRenderer, taxonomyRepo), mx, authMiddleware)
	requiresAdmin(http.MethodPost, "/api/video", video.SaveHandler(apiRenderer, videoRepo), mx, authMiddleware)
	requiresAdmin(http.MethodPut, "/api/video", video.UpdateHandler(apiRenderer, videoRepo), mx, authMiddleware)

	//APIs that only require Authentication
	requiresAuth(http.MethodGet, "/api/taxonomy", taxonomy.ListHandler(apiRenderer, taxonomyRepo), mx, authMiddleware)
	requiresAuth(http.MethodGet, "/api/taxonomy/{id:[0-9]+}/children", taxonomy.ListChildrenHandler(apiRenderer, taxonomyRepo), mx, authMiddleware)
	requiresAuth(http.MethodGet, "/api/video/stream/{id:[0-9]+}", video.StreamHandler(apiRenderer, videoRepo, gcpConfig), mx, authMiddleware)

	n.UseHandler(mx)
	return n
}

func requiresAuth(method string, uri string, handler http.Handler, mx *mux.Router, authMiddleware *negroni.Negroni) {
	mx.Methods(method).Path(uri).Handler(authMiddleware.With(
		negroni.Wrap(handler),
	))
}

func requiresAdmin(method string, uri string, handler http.Handler, mx *mux.Router, authMiddleware *negroni.Negroni) {
	mx.Methods(method).Path(uri).Handler(authMiddleware.With(
		negroni.HandlerFunc(auth.IsAdmin),
		negroni.Wrap(handler),
	))
}
