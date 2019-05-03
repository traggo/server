package ui

import (
	"net/http"

	"github.com/gobuffalo/packr/v2"
	"github.com/gorilla/mux"
)

var box = packr.New("UI", "build")

// Register registers the ui on the root path.
func Register(r *mux.Router) {
	r.Handle("/", serveFile("index.html", "text/html"))
	r.Handle("/index.html", serveFile("index.html", "text/html"))
	r.Handle("/manifest.json", serveFile("manifest.json", "application/json"))
	r.Handle("/service-worker.js", serveFile("service-worker.js", "text/javascript"))
	r.Handle("/assets-manifest.json", serveFile("asserts-manifest.json", "application/json"))
	r.Handle("/static/{type}/{resource}", http.FileServer(box))
}

func serveFile(name, contentType string) http.HandlerFunc {
	return func(writer http.ResponseWriter, reg *http.Request) {
		writer.Header().Set("Content-Type", contentType)
		content, err := box.Find(name)
		if err != nil {
			panic(err)
		}
		writer.Write(content)
	}
}
