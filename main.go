package main

import (
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"github.com/traggo/server/database"
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

	port := 3030
	log.Info().Int("port", port).Msg("Start listening")
	if err := server.Start(db, 10, port); err != nil {
		log.Fatal().Err(err).Msg("Server Error")
	}
}
