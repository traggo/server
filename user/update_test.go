package user

import (
	"context"
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/traggo/server/generated/gqlmodel"
	"github.com/traggo/server/model"
	"github.com/traggo/server/test"
)

func TestGQL_UpdateUser_succeeds_updatesUser(t *testing.T) {
	createPassword = fakePassword
	db := test.InMemoryDB(t)
	defer db.Close()
	db.Create(&model.User{
		Name:  "jmattheis",
		Pass:  unicornPW,
		ID:    1,
		Admin: true,
	})

	resolver := ResolverForUser{DB: db, PassStrength: 4}
	user, err := resolver.UpdateUser(context.Background(), 1, "broder", pointer("pony"), false)
	require.Nil(t, err)

	expected := &gqlmodel.User{
		Name:  "broder",
		ID:    1,
		Admin: false,
	}
	require.Equal(t, expected, user)
	assertUserCount(t, db, 1)
	assertUserExist(t, db, model.User{
		Name:  "broder",
		ID:    1,
		Admin: false,
		Pass:  ponyPW,
	})
}

func TestGQL_UpdateUser_succeeds_preservesPassword(t *testing.T) {
	createPassword = fakePassword
	db := test.InMemoryDB(t)
	defer db.Close()
	db.Create(&model.User{
		Name:  "jmattheis",
		Pass:  unicornPW,
		ID:    1,
		Admin: true,
	})

	resolver := ResolverForUser{DB: db, PassStrength: 4}
	user, err := resolver.UpdateUser(context.Background(), 1, "broder", nil, false)
	require.Nil(t, err)

	expected := &gqlmodel.User{
		Name:  "broder",
		ID:    1,
		Admin: false,
	}
	require.Equal(t, expected, user)
	assertUserCount(t, db, 1)
	assertUserExist(t, db, model.User{
		Name:  "broder",
		ID:    1,
		Admin: false,
		Pass:  unicornPW,
	})
}

func TestGQL_UpdateUser_fails_notExistingUser(t *testing.T) {
	createPassword = fakePassword
	db := test.InMemoryDB(t)
	defer db.Close()

	resolver := ResolverForUser{DB: db, PassStrength: 4}
	_, err := resolver.UpdateUser(context.Background(), 1, "broder", nil, false)
	require.EqualError(t, err, "user with id 1 does not exist")

	assertUserCount(t, db, 0)
}

func pointer(s string) *string {
	return &s
}
