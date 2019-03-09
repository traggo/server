package test

import (
	"fmt"

	"github.com/jinzhu/gorm"
	"github.com/stretchr/testify/assert"
	"github.com/traggo/server/database"
	"github.com/traggo/server/model"
)

// InMemoryDB create a in memory database for testing.
func InMemoryDB(t assert.TestingT) *Database {
	db, err := database.New("sqlite3", "file::memory:?mode=memory&cache=shared")
	assert.Nil(t, err)
	return &Database{DB: db}
}

// Database wraps the gorm.DB and provides helper methods
type Database struct {
	*gorm.DB
}

// User creates a user
func (d *Database) User(id int) *UserWithDatabase {
	user := d.NewUser(id, fmt.Sprint("test", id), false)
	return &UserWithDatabase{
		User: user,
		DB:   d.DB,
	}
}

// UserWithDatabase wraps gorm.DB and provides helper methods
type UserWithDatabase struct {
	User model.User
	*gorm.DB
}

// NewDevice creates a device.
func (d *UserWithDatabase) NewDevice(id int, token string, name string) model.Device {
	device := model.Device{
		ID:        id,
		Token:     token,
		UserID:    d.User.ID,
		Name:      name,
		ExpiresAt: Time("2009-06-30T18:30:00Z"),
		ActiveAt:  Time("2009-06-30T18:30:00Z"),
		CreatedAt: Time("2009-06-30T18:30:00Z"),
	}
	d.Create(&device)
	return device
}

// NewUser creates a user
func (d *Database) NewUser(id int, name string, admin bool) model.User {
	return d.NewUserPass(id, name, []byte{1, 2, 3}, admin)
}

// NewUserPass creates a user
func (d *Database) NewUserPass(id int, name string, pass []byte, admin bool) model.User {
	user := model.User{
		ID:    id,
		Name:  name,
		Pass:  pass,
		Admin: admin,
	}
	d.Create(&user)
	return user
}
