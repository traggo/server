package test_test

import (
	"testing"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"github.com/stretchr/testify/assert"
	"github.com/traggo/server/test"
)

func TestLogger_AssertCount_Succeeds(t *testing.T) {
	logger := test.NewLogger(t)
	defer logger.Dispose()

	log.Error().Msg("error")

	logger.AssertCount(1)
}

func TestLogger_AssertEntry_Succeeds(t *testing.T) {
	logger := test.NewLogger(t)
	defer logger.Dispose()

	log.Error().Msg("error")

	logger.AssertEntryExists(test.Entry{Message: "error", Level: zerolog.ErrorLevel})
}

func TestLogger_AssertCount_Fails(t *testing.T) {
	fake := &fakeTesting{}
	logger := test.NewLogger(fake)
	defer logger.Dispose()

	log.Error().Msg("error")

	logger.AssertCount(2)
	assert.True(t, fake.hasErrors)
}

func TestLogger_AssertEntry_Fails_wrongLevel(t *testing.T) {
	fake := &fakeTesting{}
	logger := test.NewLogger(fake)
	defer logger.Dispose()

	log.Error().Msg("error")

	logger.AssertEntryExists(test.Entry{Message: "error", Level: zerolog.InfoLevel})

	assert.True(t, fake.hasErrors)
}

func TestLogger_AssertEntry_Fails_wrongMessage(t *testing.T) {
	fake := &fakeTesting{}
	logger := test.NewLogger(fake)
	defer logger.Dispose()

	log.Error().Msg("error")

	logger.AssertEntryExists(test.Entry{Message: "info", Level: zerolog.ErrorLevel})

	assert.True(t, fake.hasErrors)
}
