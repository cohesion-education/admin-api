package config

import (
	"cloud.google.com/go/datastore"
	"github.com/unrolled/render"
)

type HandlerConfig struct {
	Renderer        *render.Render
	DatastoreClient *datastore.Client
}

func NewHandlerConfig(datastoreClient *datastore.Client) (*HandlerConfig, error) {
	return &HandlerConfig{
		DatastoreClient: datastoreClient,
		Renderer: render.New(render.Options{
			Layout: "layout",
			RenderPartialsWithoutPrefix: true,
		}),
	}, nil
}
