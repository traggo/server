package device

import (
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/traggo/server/model"
	"github.com/traggo/server/test"
	"github.com/traggo/server/test/fake"
)

func TestGQL_RemoveDevice_succeeds_removesDevice(t *testing.T) {
	db := test.InMemoryDB(t)
	defer db.Close()
	db.User(1)
	db.Create(&model.Device{
		ID:        55,
		Name:      "Android",
		Token:     "abc",
		UserID:    1,
		CreatedAt: test.Time("2009-06-30T18:30:00+02:00"),
		ActiveAt:  test.Time("2018-06-30T18:30:00+02:00"),
		ExpiresAt: test.Time("2022-06-30T18:30:00+02:00"),
	})
	resolver := ResolverForDevice{DB: db.DB}
	_, err := resolver.RemoveDevice(fake.User(1), 55)
	require.Nil(t, err)
	assertDeviceCount(t, db, 0)
}

func TestGQL_RemoveDevice_fails_notExistingDevice(t *testing.T) {
	db := test.InMemoryDB(t)
	defer db.Close()
	db.User(1)

	resolver := ResolverForDevice{DB: db.DB}
	_, err := resolver.RemoveDevice(fake.User(1), 55)
	require.EqualError(t, err, "device not found")
}

func TestGQL_RemoveDevice_fails_notPermission(t *testing.T) {
	db := test.InMemoryDB(t)
	defer db.Close()
	db.User(1)
	db.User(2)
	db.Create(&model.Device{
		ID:        55,
		Name:      "Android",
		Token:     "abc",
		UserID:    2,
		CreatedAt: test.Time("2009-06-30T18:30:00+02:00"),
		ActiveAt:  test.Time("2018-06-30T18:30:00+02:00"),
		ExpiresAt: test.Time("2022-06-30T18:30:00+02:00"),
	})
	resolver := ResolverForDevice{DB: db.DB}
	_, err := resolver.RemoveDevice(fake.User(1), 55)
	require.EqualError(t, err, "device not found")
}
