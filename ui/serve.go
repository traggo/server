package ui

import (
	"net/http"

	"github.com/gobuffalo/packr/v2"
	"github.com/gorilla/mux"
)

// Register registers the ui on the root path.
func Register(r *mux.Router, box *packr.Box) {
	r.Handle("/", serveFile("index.html", "text/html", box))
	r.Handle("/index.html", serveFile("index.html", "text/html", box))
	r.Handle("/manifest.json", serveFile("manifest.json", "application/json", box))
	r.Handle("/service-worker.js", serveFile("service-worker.js", "text/javascript", box))
	r.Handle("/assets-manifest.json", serveFile("asserts-manifest.json", "application/json", box))
	r.Handle("/static/{type}/{resource}", http.FileServer(box))
}

func serveFile(name, contentType string, box *packr.Box) http.HandlerFunc {
	return func(writer http.ResponseWriter, reg *http.Request) {
		writer.Header().Set("Content-Type", contentType)
		content, err := box.Find(name)
		if err != nil {
			panic(err)
		}
		writer.Write(content)
	}
}
