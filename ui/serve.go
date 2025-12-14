package ui

import (
	"embed"
	"fmt"
	"io"
	"io/fs"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/rs/zerolog/log"
)

//go:embed build
var uiDir embed.FS
var buildDir, _ = fs.Sub(uiDir, "build")

// Register registers the ui on the root path.
func Register(r *mux.Router) {
	r.Handle("/", serveFile("index.html", "text/html"))
	r.Handle("/index.html", serveFile("index.html", "text/html"))
	r.Handle("/manifest.json", serveFile("manifest.json", "application/json"))
	r.Handle("/service-worker.js", serveFile("service-worker.js", "text/javascript"))
	r.Handle("/asset-manifest.json", serveFile("asset-manifest.json", "application/json"))

	fileServer := http.FileServer(http.FS(buildDir))
	r.Handle("/static/{type}/{resource}", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Cache-Control", "public, max-age=31536000, immutable")
		fileServer.ServeHTTP(w, r)
	}))

	r.Handle("/favicon.ico", serveFile("favicon.ico", "image/x-icon"))
	for _, size := range []string{"16x16", "32x32", "192x192", "256x256"} {
		fileName := fmt.Sprintf("favicon-%s.png", size)
		r.Handle("/"+fileName, serveFile(fileName, "image/png"))
	}
}

func serveFile(name, contentType string) http.HandlerFunc {
	file, err := buildDir.Open(name)
	if err != nil {
		log.Panic().Err(err).Msgf("could not find %s", file)
	}
	defer file.Close()
	content, err := io.ReadAll(file)
	if err != nil {
		log.Panic().Err(err).Msgf("could not read %s", file)
	}

	return func(writer http.ResponseWriter, reg *http.Request) {
		writer.Header().Set("Content-Type", contentType)
		_, _ = writer.Write(content)
	}
}
