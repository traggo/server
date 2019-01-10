package user

import (
	"context"
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/traggo/server/generated/gqlmodel"
	"github.com/traggo/server/test"
	"github.com/traggo/server/test/fake"
)

func TestGQL_CurrentUser_withUser(t *testing.T) {
	db := test.InMemoryDB(t)
	defer db.Close()

	resolver := ResolverForUser{DB: db.DB}
	result, err := resolver.CurrentUser(fake.UserWithPerm(2, true))

	require.Nil(t, err)
	expected := &gqlmodel.User{
		ID:    2,
		Name:  "fake",
		Admin: true,
	}

	require.Equal(t, expected, result)
}

func TestGQL_CurrentUser_noUser(t *testing.T) {
	db := test.InMemoryDB(t)
	defer db.Close()

	resolver := ResolverForUser{DB: db.DB}
	result, err := resolver.CurrentUser(context.Background())

	require.Nil(t, err)
	require.Nil(t, result)
}
