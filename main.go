package main

import (
	"os"
	"time"

	"github.com/gorilla/mux"
	"github.com/jinzhu/gorm"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"github.com/traggo/server/auth"
	"github.com/traggo/server/config"
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

	conf, errs := config.Get()
	logger.Init(conf.LogLevel.AsZeroLogLevel())

	exit := false
	for _, err := range errs {
		log.WithLevel(err.Level).Msg(err.Msg)
		exit = exit || err.Level == zerolog.FatalLevel || err.Level == zerolog.PanicLevel
	}
	if exit {
		os.Exit(1)
	}
	log.Debug().Interface("config", conf).Msg("Using")

	db := initDatabase(conf)

	stopCleanUp := initCleanUp(db)

	router := initRouter(db, conf)

	log.Info().Int("port", conf.Port).Msg("Start listening")
	if err := server.Start(router, conf.Port); err != nil {
		log.Fatal().Err(err).Msg("Server Error")
	}
	stopCleanUp <- true
}

func initRouter(db *gorm.DB, conf config.Config) *mux.Router {
	gqlHandler := graphql.Handler(
		"/graphql",
		graphql.NewResolver(db, conf.PassStrength),
		graphql.NewDirective())

	router := mux.NewRouter()
	router.Use(auth.Middleware(db))
	router.HandleFunc("/graphql", gqlHandler)
	return router
}

func initDatabase(conf config.Config) *gorm.DB {
	db, err := database.New(conf.DatabaseDialect, conf.DatabaseConnection)
	if err != nil {
		log.Fatal().Err(err).Msg("Failed to connect to the database")
	}
	c := new(int)
	if db.Model(new(model.User)).Count(c); *c == 0 {
		log.Info().Msg("Creating default user.")
		db.Create(&model.User{
			Name:  conf.DefaultUserName,
			Pass:  password.CreatePassword(conf.DefaultUserPass, conf.PassStrength),
			Admin: true})
	}

	return db
}

func initCleanUp(db *gorm.DB) chan<- bool {
	stopCleanUp := make(chan bool)
	go auth.CleanUp(db, time.Hour, stopCleanUp)
	return stopCleanUp
}
