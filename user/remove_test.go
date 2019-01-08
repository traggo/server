package user

import (
	"context"
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/traggo/server/test"
)

func TestGQL_RemoveUser_succeeds_removesUser(t *testing.T) {
	db := test.InMemoryDB(t)
	defer db.Close()
	db.NewUserPass(1, "jmattheis", unicornPW, true)

	resolver := ResolverForUser{DB: db.DB, PassStrength: 4}
	_, err := resolver.RemoveUser(context.Background(), 1)

	require.Nil(t, err)
	assertUserCount(t, db, 0)
}

func TestGQL_RemoveUser_fails_notExistingUser(t *testing.T) {
	db := test.InMemoryDB(t)
	defer db.Close()
	db.NewUserPass(1, "jmattheis", unicornPW, true)

	resolver := ResolverForUser{DB: db.DB, PassStrength: 4}
	_, err := resolver.RemoveUser(context.Background(), 3)

	require.EqualError(t, err, "user with id 3 does not exist")
	assertUserCount(t, db, 1)
}
