package device

import (
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/traggo/server/generated/gqlmodel"
	"github.com/traggo/server/model"
	"github.com/traggo/server/test"
	"github.com/traggo/server/test/fake"
)

func TestGQL_RemoveCurrentDevice_succeeds_removesDevice(t *testing.T) {
	db := test.InMemoryDB(t)
	defer db.Close()
	db.User(1)
	device := &model.Device{
		ID:        55,
		Name:      "Android",
		Token:     "abc",
		UserID:    1,
		CreatedAt: test.Time("2009-06-30T18:30:00+02:00"),
		ActiveAt:  test.Time("2018-06-30T18:30:00+02:00"),
		ExpiresAt: test.Time("2022-06-30T18:30:00+02:00"),
	}
	db.Create(device)
	resolver := ResolverForDevice{DB: db.DB}
	gqlDevice, err := resolver.RemoveCurrentDevice(fake.Device(device))
	require.Nil(t, err)

	expected := &gqlmodel.Device{
		ID:        55,
		Name:      "Android",
		CreatedAt: test.ModelTime("2009-06-30T18:30:00+02:00"),
		ActiveAt:  test.ModelTime("2018-06-30T18:30:00+02:00"),
		ExpiresAt: test.ModelTime("2022-06-30T18:30:00+02:00"),
	}

	require.Equal(t, expected, gqlDevice)
	assertDeviceCount(t, db, 0)
}
