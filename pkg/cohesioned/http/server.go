package http

import (
	"fmt"
	"log"
	"net/http"

	"github.com/cohesion-education/api/pkg/cohesioned/auth"
	"github.com/cohesion-education/api/pkg/cohesioned/config"
	"github.com/cohesion-education/api/pkg/cohesioned/profile"
	"github.com/cohesion-education/api/pkg/cohesioned/student"
	"github.com/cohesion-education/api/pkg/cohesioned/taxonomy"
	"github.com/cohesion-education/api/pkg/cohesioned/video"
	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
	"github.com/joho/godotenv"
	"github.com/unrolled/render"
	"github.com/urfave/negroni"
)

var (
	apiRenderer = render.New()
)

func Run(port string) {
	s := newServer()
	s.Run(":" + port)
}

func newServer() *negroni.Negroni {
	if err := godotenv.Load(); err != nil {
		fmt.Println("Failed to load .env file. Components will fallback to loading from VCAP_SERVICES or env vars")
	}

	authConfig, err := config.NewAuthConfig()
	if err != nil {
		log.Fatal(err)
	}

	awsConfig, err := config.NewAwsConfig()
	if err != nil {
		log.Fatal(err)
	}

	db, err := awsConfig.DialRDS()
	if err != nil {
		log.Fatal(err)
	}

	profileRepo := profile.NewAwsRepo(db)
	taxonomyRepo := taxonomy.NewAwsRepo(db)
	studentRepo := student.NewAwsRepo(db)
	videoRepo := video.NewAwsRepo(db, awsConfig)

	adminVideoService := video.NewService(videoRepo, awsConfig)

	n := negroni.Classic()
	mx := mux.NewRouter()
	mx.StrictSlash(true)

	//TODO - register 404
	// mx.NotFoundHandler = cohesioned.NotFoundViewHandler(homepageRenderer)

	//Public APIs
	// mx.Methods(http.MethodGet).Path("/api/homepage").Handler(homepage.HomepageHandler(apiRenderer, homepageRepo))
	mx.Methods(http.MethodGet).Path("/api/taxonomy").Handler(taxonomy.ListHandler(apiRenderer, taxonomyRepo))
	mx.Methods(http.MethodGet).Path("/api/taxonomy/{id:[0-9]+}/children").Handler(taxonomy.ListChildrenHandler(apiRenderer, taxonomyRepo))
	mx.Methods(http.MethodGet).Path("/api/taxonomy/recursive").Handler(taxonomy.RecursiveListHandler(apiRenderer, taxonomyRepo))
	mx.Methods(http.MethodGet).Path("/api/taxonomy/flatten").Handler(taxonomy.FlatListHandler(apiRenderer, taxonomyRepo))

	authMiddleware := negroni.New(
		negroni.HandlerFunc(auth.CheckJwt(apiRenderer, profileRepo, authConfig)),
	)

	//endpoints that require Admin priveleges
	requiresAdmin(http.MethodPost, "/api/taxonomy", taxonomy.AddHandler(apiRenderer, taxonomyRepo), mx, authMiddleware)
	requiresAdmin(http.MethodPut, "/api/taxonomy/{id:[0-9]+}", taxonomy.UpdateHandler(apiRenderer, taxonomyRepo), mx, authMiddleware)
	requiresAdmin(http.MethodGet, "/api/videos", video.ListHandler(apiRenderer, adminVideoService), mx, authMiddleware)
	requiresAdmin(http.MethodPost, "/api/video", video.AddHandler(apiRenderer, adminVideoService), mx, authMiddleware)
	requiresAdmin(http.MethodPost, "/api/video/upload/{id:[0-9]+}", video.UploadHandler(apiRenderer, adminVideoService), mx, authMiddleware)
	requiresAdmin(http.MethodDelete, "/api/video/{id:[0-9]+}", video.DeleteHandler(apiRenderer, adminVideoService), mx, authMiddleware)
	requiresAdmin(http.MethodPut, "/api/video/{id:[0-9]+}", video.UpdateHandler(apiRenderer, adminVideoService), mx, authMiddleware)

	//endpoints that only require Authentication
	requiresAuth(http.MethodGet, "/api/profile", profile.GetCurrentUserHandler(apiRenderer, profileRepo), mx, authMiddleware)
	requiresAuth(http.MethodPost, "/api/profile", profile.SaveHandler(apiRenderer, profileRepo), mx, authMiddleware)
	requiresAuth(http.MethodPut, "/api/profile", profile.UpdateHandler(apiRenderer, profileRepo), mx, authMiddleware)
	requiresAuth(http.MethodGet, "/api/profile/students", student.ListHandler(apiRenderer, studentRepo), mx, authMiddleware)
	requiresAuth(http.MethodPost, "/api/profile/students", student.SaveHandler(apiRenderer, studentRepo), mx, authMiddleware)
	requiresAuth(http.MethodPost, "/api/profile/preferences", profile.SavePreferencesHandler(apiRenderer, profileRepo), mx, authMiddleware)
	// requiresAuth(http.MethodGet, "/api/video/{id:[0-9]+}", video.GetByIDHandler(apiRenderer, videoRepo, awsConfig), mx, authMiddleware)

	allowedOrigins := handlers.AllowedOrigins([]string{"*"})
	allowedHeaders := handlers.AllowedHeaders([]string{"authorization", "content-type", "content-length"})
	allowedMethods := handlers.AllowedMethods([]string{"GET", "POST", "PUT", "DELETE", "HEAD"})
	n.UseHandler(handlers.CORS(allowedOrigins, allowedHeaders, allowedMethods)(mx))
	return n
}

func requiresAuth(method string, uri string, handler http.Handler, mx *mux.Router, authMiddleware *negroni.Negroni) {
	mx.Methods(method).Path(uri).Handler(authMiddleware.With(
		negroni.Wrap(handler),
	))
}

func requiresAdmin(method string, uri string, handler http.Handler, mx *mux.Router, authMiddleware *negroni.Negroni) {
	mx.Methods(method).Path(uri).Handler(authMiddleware.With(
		negroni.HandlerFunc(auth.IsAdmin(apiRenderer)),
		negroni.Wrap(handler),
	))
}
