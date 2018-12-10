package test

import (
	"github.com/jinzhu/gorm"
	"github.com/stretchr/testify/assert"
	"github.com/traggo/server/database"
)

// InMemoryDB create a in memory database for testing.
func InMemoryDB(t assert.TestingT) *gorm.DB {
	db, err := database.New("sqlite3", "file::memory:?mode=memory&cache=shared")
	assert.Nil(t, err)
	return db
}
