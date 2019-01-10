package device

import (
	"context"
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/traggo/server/generated/gqlmodel"
	"github.com/traggo/server/model"
	"github.com/traggo/server/test"
	"github.com/traggo/server/test/fake"
)

func TestGQL_CurrentDevice_withDevice(t *testing.T) {
	db := test.InMemoryDB(t)
	defer db.Close()

	device := &model.Device{
		ID:        2,
		Name:      "Browser",
		Token:     "abcd",
		UserID:    2,
		CreatedAt: test.Time("2004-06-30T18:30:00+02:00"),
		ActiveAt:  test.Time("2015-06-30T18:30:00+02:00"),
		ExpiresAt: test.Time("2026-06-30T18:30:00+02:00"),
	}

	resolver := ResolverForDevice{DB: db.DB}
	result, err := resolver.CurrentDevice(fake.Device(device))

	require.Nil(t, err)
	expected := &gqlmodel.Device{
		ID:        2,
		Name:      "Browser",
		CreatedAt: test.ModelTime("2004-06-30T18:30:00+02:00"),
		ActiveAt:  test.ModelTime("2015-06-30T18:30:00+02:00"),
		ExpiresAt: test.ModelTime("2026-06-30T18:30:00+02:00"),
	}

	require.Equal(t, expected, result)
}

func TestGQL_CurrentDevice_noDevice(t *testing.T) {
	db := test.InMemoryDB(t)
	defer db.Close()

	resolver := ResolverForDevice{DB: db.DB}
	result, err := resolver.CurrentDevice(context.Background())

	require.Nil(t, err)

	require.Nil(t, result)
}
