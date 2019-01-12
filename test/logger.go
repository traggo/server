package test

import (
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"github.com/stretchr/testify/assert"
	"github.com/traggo/server/logger"
)

// Logger the test logger with util methods
type Logger struct {
	old     zerolog.Logger
	entries []Entry
	t       assert.TestingT
}

// Entry a test logging entry
type Entry struct {
	Level   zerolog.Level
	Message string
}

// Dispose resets the logger
func (l *Logger) Dispose() {
	log.Logger = l.old
}

// NewLogger creates a new test logger
func NewLogger(t assert.TestingT) *Logger {
	logger := &Logger{t: t}
	log.Logger = zerolog.New(&noop{}).With().Timestamp().Logger().Hook(logger)
	return logger
}

// Run records log entries
func (l *Logger) Run(e *zerolog.Event, level zerolog.Level, message string) {
	l.entries = append(l.entries, Entry{Level: level, Message: message})
}

// AssertCount asserts the amount of recorded logging entries
func (l *Logger) AssertCount(count int) {
	assert.Len(l.t, l.entries, count)
}

// AssertEntryExists asserts that a logging entry exists
func (l *Logger) AssertEntryExists(entry Entry) {
	assert.Contains(l.t, l.entries, entry)
}

type noop struct {
}

func (*noop) Write(p []byte) (n int, err error) {
	return len(p), nil
}

// LogDebug enables debug log
func LogDebug() {
	logger.Init(zerolog.DebugLevel)
}
