package device

import (
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/traggo/server/generated/gqlmodel"
	"github.com/traggo/server/model"
	"github.com/traggo/server/test"
	"github.com/traggo/server/test/fake"
)

func TestGQL_UpdateDevice_succeeds_updatesDevice(t *testing.T) {
	db := test.InMemoryDB(t)
	defer db.Close()
	db.Create(&model.User{
		Name:  "jmattheis",
		ID:    1,
		Admin: true,
	})
	db.Create(&model.Device{
		Name:      "old name",
		ID:        1,
		UserID:    1,
		CreatedAt: test.Time("2009-06-30T18:30:00+02:00"),
		ActiveAt:  test.Time("2018-06-30T18:30:00+02:00"),
		ExpiresAt: test.Time("2022-06-30T18:30:00+02:00"),
	})

	resolver := ResolverForDevice{DB: db.DB}
	device, err := resolver.UpdateDevice(fake.User(1), 1, "updated name")
	require.Nil(t, err)

	expected := &gqlmodel.Device{
		Name:      "updated name",
		ID:        1,
		CreatedAt: test.ModelTime("2009-06-30T18:30:00+02:00"),
		ActiveAt:  test.ModelTime("2018-06-30T18:30:00+02:00"),
		ExpiresAt: test.ModelTime("2022-06-30T18:30:00+02:00"),
	}
	require.Equal(t, expected, device)
	assertDeviceCount(t, db, 1)
	assertDeviceExist(t, db, model.Device{
		Name:      "updated name",
		ID:        1,
		UserID:    1,
		CreatedAt: test.Time("2009-06-30T18:30:00+02:00"),
		ActiveAt:  test.Time("2018-06-30T18:30:00+02:00"),
		ExpiresAt: test.Time("2022-06-30T18:30:00+02:00"),
	})
}

func TestGQL_UpdateDevice_fails_notExistingDevice(t *testing.T) {
	db := test.InMemoryDB(t)
	defer db.Close()
	db.Create(&model.User{
		Name:  "jmattheis",
		ID:    1,
		Admin: true,
	})
	resolver := ResolverForDevice{DB: db.DB}
	_, err := resolver.UpdateDevice(fake.User(1), 3, "tst")
	require.EqualError(t, err, "device not found")

	assertDeviceCount(t, db, 0)
}

func TestGQL_UpdateDevice_fails_noPermissions(t *testing.T) {
	db := test.InMemoryDB(t)
	defer db.Close()
	db.Create(&model.User{
		Name:  "jmattheis",
		ID:    1,
		Admin: true,
	})
	db.Create(&model.User{
		Name:  "broderp",
		ID:    2,
		Admin: true,
	})
	db.Create(&model.Device{
		Name:      "old name",
		ID:        66,
		UserID:    2,
		CreatedAt: test.Time("2009-06-30T18:30:00+02:00"),
		ActiveAt:  test.Time("2018-06-30T18:30:00+02:00"),
		ExpiresAt: test.Time("2022-06-30T18:30:00+02:00"),
	})
	resolver := ResolverForDevice{DB: db.DB}
	_, err := resolver.UpdateDevice(fake.User(1), 66, "tst")
	require.EqualError(t, err, "device not found")

	assertDeviceCount(t, db, 1)
}
