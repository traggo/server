package main

import (
	"net/http"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"github.com/traggo/server/database"
	"github.com/traggo/server/gql"
	"github.com/traggo/server/logger"
)

func main() {
	// TODO configurable
	logger.Init(zerolog.DebugLevel)
	db, err := database.New("sqlite3", "file::memory:?mode=memory&cache=shared")
	if err != nil {
		log.Fatal().Err(err).Msg("Failed to connect to the database")
	}

	http.Handle("/graphql", gql.Handler(db, 10))
	http.ListenAndServe(":3030", nil) // TODO configurable port
}
