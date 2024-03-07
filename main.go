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
	"github.com/traggo/server/ui"
	"github.com/traggo/server/user/password"
)

var (
	// BuildMode will be set via ldflags
	BuildMode = mode.Dev
	// BuildCommit will be set via ldflags
	BuildCommit = "-"
	// BuildVersion will be set via ldflags
	BuildVersion = "develop"
	// BuildDate will be set via ldflags
	BuildDate = "unknown"
)

func main() {
	mode.Set(BuildMode)

	version := model.Version{Commit: BuildCommit, BuildDate: BuildDate, Name: BuildVersion}

	conf, errs := config.Get()
	logger.Init(conf.LogLevel.AsZeroLogLevel())
	log.Info().Str("commit", BuildCommit).Str("buildDate", BuildDate).Str("buildMode", BuildMode).Str("version", BuildVersion).Msg("Traggo")

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

	router := initRouter(db, conf, version)

	log.Info().Int("port", conf.Port).Msg("Start listening")
	if err := server.Start(router, conf.Port); err != nil {
		log.Fatal().Err(err).Msg("Server Error")
	}
	stopCleanUp <- true
}

func initRouter(db *gorm.DB, conf config.Config, version model.Version) *mux.Router {
	gqlHandler := graphql.Handler(
		"/graphql",
		graphql.NewResolver(db, conf.PassStrength, version),
		graphql.NewDirective())

	router := mux.NewRouter()
	router.Use(auth.Middleware(db))
	router.HandleFunc("/graphql", gqlHandler)
	ui.Register(router)
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
	go auth.CleanUp(db, time.Minute*10, stopCleanUp)
	return stopCleanUp
}
