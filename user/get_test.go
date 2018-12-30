package user

import (
	"context"
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/traggo/server/generated/gqlmodel"
	"github.com/traggo/server/schema"
	"github.com/traggo/server/test"
)

func TestGQL_Users(t *testing.T) {
	db := test.InMemoryDB(t)
	defer db.Close()
	resolver := ResolverForUser{DB: db, PassStrength: 4}
	db.Create(&schema.User{
		Name:  "jmattheis",
		Pass:  unicornPW,
		ID:    1,
		Admin: true,
	})
	db.Create(&schema.User{
		Name:  "broderpeters",
		Pass:  ponyPW,
		ID:    2,
		Admin: false,
	})

	users, err := resolver.Users(context.Background())

	require.Nil(t, err)
	expected := []gqlmodel.User{
		{
			Name:  "jmattheis",
			ID:    1,
			Admin: true,
		},
		{
			Name:  "broderpeters",
			ID:    2,
			Admin: false,
		},
	}
	require.Equal(t, expected, users)
}
