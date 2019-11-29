package database

import (
	"os"
	"path/filepath"

	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"    // enable the mysql dialect
	_ "github.com/jinzhu/gorm/dialects/postgres" // enable the postgres dialect
	_ "github.com/jinzhu/gorm/dialects/sqlite"   // enable the sqlite3 dialect
	"github.com/rs/zerolog/log"
	"github.com/traggo/server/logger"
	"github.com/traggo/server/model"
)

var mkdirAll = os.MkdirAll

// New creates a gorm instance.
func New(dialect, connection string) (*gorm.DB, error) {
	createDirectoryIfSqlite(dialect, connection)

	db, err := gorm.Open(dialect, connection)
	if err != nil {
		return nil, err
	}
	db.LogMode(true)
	db.SetLogger(&logger.DatabaseLogger{})

	// We normally don't need that much connections, so we limit them. F.ex. mysql complains about
	// "too many connections".
	db.DB().SetMaxOpenConns(10)

	if dialect == "sqlite3" {
		// We use the database connection inside the handlers from the http
		// framework, therefore concurrent access occurs. Sqlite cannot handle
		// concurrent writes, so we limit sqlite to one connection.
		// see https://github.com/mattn/go-sqlite3/issues/274
		db.DB().SetMaxOpenConns(1)
		db.Exec("PRAGMA foreign_keys = ON")
	}

	log.Debug().Msg("Auto migrating schema's")
	db.AutoMigrate(model.All()...)

	log.Debug().Msg("Database initialized")
	return db, nil
}

func createDirectoryIfSqlite(dialect string, connection string) {
	if dialect == "sqlite3" {
		if _, err := os.Stat(filepath.Dir(connection)); os.IsNotExist(err) {
			if err := mkdirAll(filepath.Dir(connection), 0777); err != nil {
				panic(err)
			}
		}
	}
}
