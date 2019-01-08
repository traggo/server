package user

import (
	"context"
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/traggo/server/generated/gqlmodel"
	"github.com/traggo/server/model"
	"github.com/traggo/server/test"
)

func TestGQL_CreateUser_succeeds_addsUser(t *testing.T) {
	createPassword = fakePassword
	db := test.InMemoryDB(t)
	defer db.Close()

	resolver := ResolverForUser{DB: db.DB, PassStrength: 4}
	user, err := resolver.CreateUser(context.Background(), "jmattheis", "unicorn", true)

	require.Nil(t, err)
	expected := &gqlmodel.User{
		Name:  "jmattheis",
		Admin: true,
		ID:    1,
	}
	require.Equal(t, expected, user)
	assertUserExist(t, db, model.User{
		Name:  "jmattheis",
		Pass:  unicornPW,
		ID:    1,
		Admin: true,
	})
	assertUserCount(t, db, 1)
}

func TestGQL_CreateUser_fails_userAlreadyExists(t *testing.T) {
	createPassword = fakePassword
	db := test.InMemoryDB(t)
	defer db.Close()
	db.Create(&model.User{
		Name:  "jmattheis",
		Pass:  unicornPW,
		ID:    1,
		Admin: true,
	})

	resolver := ResolverForUser{DB: db.DB, PassStrength: 4}
	_, err := resolver.CreateUser(context.Background(), "jmattheis", "unicorn", true)
	require.EqualError(t, err, "user with name 'jmattheis' does already exist")
	assertUserCount(t, db, 1)
}
