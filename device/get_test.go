package device

import (
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/traggo/server/generated/gqlmodel"
	"github.com/traggo/server/model"
	"github.com/traggo/server/test"
	"github.com/traggo/server/test/fake"
)

func TestGQL_Devices(t *testing.T) {
	db := test.InMemoryDB(t)
	defer db.Close()
	db.Create(&model.Device{
		ID:        1,
		Name:      "Android",
		Token:     "abc",
		UserID:    1,
		CreatedAt: test.Time("2009-06-30T18:30:00Z"),
		ActiveAt:  test.Time("2018-06-30T18:30:00Z"),
		ExpiresAt: test.Time("2022-06-30T18:30:00Z"),
	})
	db.Create(&model.Device{
		ID:        2,
		Name:      "Browser",
		Token:     "abcd",
		UserID:    2,
		CreatedAt: test.Time("2004-06-30T18:30:00Z"),
		ActiveAt:  test.Time("2015-06-30T18:30:00Z"),
		ExpiresAt: test.Time("2026-06-30T18:30:00Z"),
	})

	resolver := ResolverForDevice{DB: db.DB}
	devices, err := resolver.Devices(fake.User(1))

	require.Nil(t, err)
	expected := []gqlmodel.Device{
		{
			ID:        1,
			Name:      "Android",
			CreatedAt: test.ModelTime("2009-06-30T18:30:00Z"),
			ActiveAt:  test.ModelTime("2018-06-30T18:30:00Z"),
			ExpiresAt: test.ModelTime("2022-06-30T18:30:00Z"),
		},
	}
	require.Equal(t, expected, devices)
}
