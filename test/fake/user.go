package fake

import (
	"context"

	"github.com/traggo/server/auth"
	"github.com/traggo/server/model"
)

// User create a context with a fake user.
func User(id int) context.Context {
	return UserWithPerm(id, true)
}

// UserWithPerm create a context with a fake user.
func UserWithPerm(id int, admin bool) context.Context {
	return auth.WithUser(context.Background(), &model.User{
		ID:    id,
		Name:  "fake",
		Admin: admin,
	})
}
