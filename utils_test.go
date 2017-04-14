package main

import (
	"bytes"
	"net/http"

	"github.com/unrolled/render"
)

var (
	renderer = render.New(render.Options{
		Layout: "layout",
		RenderPartialsWithoutPrefix: true,
	})

	hc = &handlerConfig{
		renderer:     renderer,
		mongoSession: nil,
	}
)

func renderHTMLWithNoLayout(templateFileName string, data interface{}) []byte {
	return renderHTML(templateFileName, data, render.HTMLOptions{Layout: ""})
}

func renderHTML(templateFileName string, data interface{}, htmlOpt ...render.HTMLOptions) []byte {
	buffer := bytes.NewBuffer(make([]byte, 0))
	err := renderer.HTML(buffer, http.StatusOK, templateFileName, data, htmlOpt...)
	if err != nil {
		panic("Failed to render template " + templateFileName + "; error: " + err.Error())
	}

	return buffer.Bytes()
}
