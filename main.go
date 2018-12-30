package main

import (
	"net/http"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"github.com/traggo/server/database"
	"github.com/traggo/server/graphql"
	"github.com/traggo/server/logger"
	"github.com/traggo/server/server"
)

func main() {
	// TODO configurable
	logger.Init(zerolog.DebugLevel)
	db, err := database.New("sqlite3", "file::memory:?mode=memory&cache=shared")
	if err != nil {
		log.Fatal().Err(err).Msg("Failed to connect to the database")
	}

	gqlHandler := graphql.Handler("/graphql", graphql.NewResolver(db, 10))

	mux := http.NewServeMux()
	mux.HandleFunc("/graphql", gqlHandler)

	port := 3030
	log.Info().Int("port", port).Msg("Start listening")
	if err := server.Start(mux, port); err != nil {
		log.Fatal().Err(err).Msg("Server Error")
	}
}
