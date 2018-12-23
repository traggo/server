package logger_test

import (
	"testing"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"github.com/traggo/server/logger"
)

func TestInit_LoggingWorks(t *testing.T) {
	logger.Init(zerolog.InfoLevel)
	log.Info().Msg("test logging")
}
