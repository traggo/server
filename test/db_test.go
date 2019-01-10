package test_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/traggo/server/model"
	"github.com/traggo/server/test"
)

func TestInMemoryDB(t *testing.T) {
	db := test.InMemoryDB(t)
	defer db.Close()
	assert.NotNil(t, db)
}

func TestDatabase(t *testing.T) {
	db := test.InMemoryDB(t)
	defer db.Close()
	db.User(1).NewDevice(1, "lol", "test device")
	db.NewUser(2, "abc", true)
	db.NewUserPass(3, "xxx", []byte{5, 5}, true)

	expectedUsers := []model.User{
		{ID: 1, Name: "test1", Pass: []uint8{1, 2, 3}, Admin: false},
		{ID: 2, Name: "abc", Pass: []uint8{1, 2, 3}, Admin: true},
		{ID: 3, Name: "xxx", Pass: []uint8{5, 5}, Admin: true}}

	var users []model.User
	db.Find(&users)
	assert.Equal(t, expectedUsers, users)

	expectedDevices := []model.Device{{
		ID:        1,
		Token:     "lol",
		Name:      "test device",
		UserID:    1,
		ExpiresAt: test.Time("2009-06-30T18:30:00+02:00"),
		ActiveAt:  test.Time("2009-06-30T18:30:00+02:00"),
		CreatedAt: test.Time("2009-06-30T18:30:00+02:00")}}

	var devices []model.Device
	db.Find(&devices)
	assert.Equal(t, expectedDevices, devices)
}
