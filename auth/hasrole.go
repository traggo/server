package auth

import (
	"context"
	"errors"

	"github.com/99designs/gqlgen/graphql"
	"github.com/traggo/server/generated/gqlmodel"
)

// HasRole checks if the current user has sufficient permissions.
func HasRole() func(ctx context.Context, obj interface{}, next graphql.Resolver, role gqlmodel.Role) (res interface{}, err error) {
	return func(ctx context.Context, obj interface{}, next graphql.Resolver, role gqlmodel.Role) (interface{}, error) {
		user := GetUser(ctx)

		if user == nil {
			return nil, errors.New("you need to login")
		}

		if role == gqlmodel.RoleAdmin && !user.Admin {
			return nil, errors.New("permission denied")
		}

		return next(ctx)
	}
}
