package user

import (
	"testing"

	"github.com/jinzhu/gorm"
	"github.com/stretchr/testify/require"
	"github.com/traggo/server/schema"
)

var (
	ponyPW       = []byte{1}
	unicornPW    = []byte{2}
	fakePassword = func(pw string, strength int) []byte {
		if pw == "pony" {
			return ponyPW
		} else if pw == "unicorn" {
			return unicornPW
		}
		panic("unknown pw")
	}
)

func assertUserExist(t *testing.T, db *gorm.DB, expected schema.User) {
	foundUser := new(schema.User)
	find := db.Find(foundUser, expected.ID)
	require.Nil(t, find.Error)
	require.NotNil(t, foundUser)
	require.Equal(t, expected, *foundUser)
}

func assertUserCount(t *testing.T, db *gorm.DB, expected int) {
	count := new(int)
	db.Model(new(schema.User)).Count(count)
	require.Equal(t, expected, *count)
}
