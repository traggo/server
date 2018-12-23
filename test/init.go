package test

import (
	"github.com/rs/zerolog"
	"github.com/traggo/server/logger"
)

func init() {
	logger.Init(zerolog.WarnLevel)
}
