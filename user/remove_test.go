package user

import (
	"context"
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/traggo/server/model"
	"github.com/traggo/server/test"
)

func TestGQL_RemoveUser_succeeds_removesUser(t *testing.T) {
	db := test.InMemoryDB(t)
	defer db.Close()
	db.Create(&model.User{
		Name:  "jmattheis",
		Pass:  unicornPW,
		ID:    1,
		Admin: true,
	})
	resolver := ResolverForUser{DB: db, PassStrength: 4}
	_, err := resolver.RemoveUser(context.Background(), 1)
	require.Nil(t, err)
	assertUserCount(t, db, 0)
}

func TestGQL_RemoveUser_fails_notExistingUser(t *testing.T) {
	db := test.InMemoryDB(t)
	defer db.Close()
	db.Create(&model.User{
		Name:  "jmattheis",
		Pass:  unicornPW,
		ID:    1,
		Admin: true,
	})
	resolver := ResolverForUser{DB: db, PassStrength: 4}
	_, err := resolver.RemoveUser(context.Background(), 3)
	require.EqualError(t, err, "user with id 3 does not exist")
	assertUserCount(t, db, 1)
}
