package main

import (
	"time"

	"github.com/gorilla/mux"
	"github.com/jinzhu/gorm"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"github.com/traggo/server/auth"
	"github.com/traggo/server/config/mode"
	"github.com/traggo/server/database"
	"github.com/traggo/server/graphql"
	"github.com/traggo/server/logger"
	"github.com/traggo/server/model"
	"github.com/traggo/server/server"
	"github.com/traggo/server/user/password"
)

var (
	// Mode the build mode
	Mode = mode.Dev
)

func main() {
	mode.Set(Mode)
	// TODO configurable
	logger.Init(zerolog.DebugLevel)
	db, err := database.New("sqlite3", "file::memory:?mode=memory&cache=shared")
	if err != nil {
		log.Fatal().Err(err).Msg("Failed to connect to the database")
	}
	db.Create(&model.User{Name: "admin", Pass: password.CreatePassword("admin", 10), Admin: true})

	stopCleanUp := make(chan bool)
	go auth.CleanUp(db, time.Hour, stopCleanUp)

	router := initRouter(db)

	port := 3030
	log.Info().Int("port", port).Msg("Start listening")
	if err := server.Start(router, port); err != nil {
		log.Fatal().Err(err).Msg("Server Error")
	}
	stopCleanUp <- true
}

func initRouter(db *gorm.DB) *mux.Router {
	gqlHandler := graphql.Handler(
		"/graphql",
		graphql.NewResolver(db, 10),
		graphql.NewDirective())

	router := mux.NewRouter()
	router.Use(auth.Middleware(db))
	router.HandleFunc("/graphql", gqlHandler)
	return router
}
