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

func TestUser(t *testing.T) {
	db := test.InMemoryDB(t)
	defer db.Close()
	db.User(1)
	db.NewUser(2, "abc", true)
	db.NewUserPass(3, "xxx", []byte{5, 5}, true)

	expected := []model.User{
		{ID: 1, Name: "test1", Pass: []uint8{1, 2, 3}, Admin: false},
		{ID: 2, Name: "abc", Pass: []uint8{1, 2, 3}, Admin: true},
		{ID: 3, Name: "xxx", Pass: []uint8{5, 5}, Admin: true}}

	var users []model.User
	db.Find(&users)
	assert.Equal(t, expected, users)
}
