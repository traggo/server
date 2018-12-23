package database

import (
	"errors"
	"os"
	"testing"

	"github.com/rs/zerolog"
	"github.com/stretchr/testify/assert"
	"github.com/traggo/server/logger"
)

func TestMain(m *testing.M) {
	logger.Init(zerolog.WarnLevel)
	os.Exit(m.Run())
}

func TestInvalidDialect(t *testing.T) {
	_, err := New("asdf", "testdb.db")
	assert.NotNil(t, err)
}

func TestCreateSqliteFolder(t *testing.T) {
	// ensure path not exists
	os.RemoveAll("somepath")

	db, err := New("sqlite3", "somepath/testdb.db")
	assert.Nil(t, err)
	assert.DirExists(t, "somepath")
	db.Close()

	assert.Nil(t, os.RemoveAll("somepath"))
}

func TestWithAlreadyExistingSqliteFolder(t *testing.T) {
	// ensure path not exists
	os.RemoveAll("somepath")
	os.MkdirAll("somepath", 0777)

	db, err := New("sqlite3", "somepath/testdb.db")
	assert.Nil(t, err)
	assert.DirExists(t, "somepath")
	db.Close()

	assert.Nil(t, os.RemoveAll("somepath"))
}

func TestPanicsOnMkdirError(t *testing.T) {
	os.RemoveAll("somepath")
	mkdirAll = func(path string, perm os.FileMode) error {
		return errors.New("ERROR")
	}
	assert.Panics(t, func() {
		New("sqlite3", "somepath/test.db")
	})
}
