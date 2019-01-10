package graphql

import (
	"github.com/jinzhu/gorm"
	"github.com/traggo/server/device"
	"github.com/traggo/server/generated/gqlschema"
	"github.com/traggo/server/tag"
	"github.com/traggo/server/user"
)

// NewResolver combines all resolvers to a resolver root.
func NewResolver(db *gorm.DB, passStrength int) gqlschema.ResolverRoot {
	return &resolver{
		user.ResolverForUser{
			DB:           db,
			PassStrength: passStrength,
		},
		tag.ResolverForTag{
			DB: db,
		},
		device.ResolverForDevice{
			DB: db,
		},
	}
}

type resolver struct {
	user.ResolverForUser
	tag.ResolverForTag
	device.ResolverForDevice
}

func (r *resolver) RootMutation() gqlschema.RootMutationResolver {
	return r
}

func (r *resolver) RootQuery() gqlschema.RootQueryResolver {
	return r
}
