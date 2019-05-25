package graphql

import (
	"context"
	"github.com/traggo/server/setting"

	"github.com/jinzhu/copier"
	"github.com/jinzhu/gorm"
	"github.com/traggo/server/device"
	"github.com/traggo/server/generated/gqlmodel"
	"github.com/traggo/server/generated/gqlschema"
	"github.com/traggo/server/model"
	"github.com/traggo/server/statistics"
	"github.com/traggo/server/tag"
	"github.com/traggo/server/timespan"
	"github.com/traggo/server/user"
)

// NewResolver combines all resolvers to a resolver root.
func NewResolver(db *gorm.DB, passStrength int, version model.Version) gqlschema.ResolverRoot {
	return &resolver{
		ResolverForUser: user.ResolverForUser{
			DB:           db,
			PassStrength: passStrength,
		},
		ResolverForTag: tag.ResolverForTag{
			DB: db,
		},
		ResolverForDevice: device.ResolverForDevice{
			DB: db,
		},
		ResolverForTimeSpan: timespan.ResolverForTimeSpan{
			DB: db,
		},
		ResolverForStatistics: statistics.ResolverForStatistics{
			DB: db,
		},
		ResolverForSettings: setting.ResolverForSettings{
			DB: db,
		},
		version: version,
	}
}

type resolver struct {
	user.ResolverForUser
	tag.ResolverForTag
	device.ResolverForDevice
	timespan.ResolverForTimeSpan
	statistics.ResolverForStatistics
	version model.Version
	setting.ResolverForSettings
}

func (r *resolver) RootMutation() gqlschema.RootMutationResolver {
	return r
}

func (r *resolver) RootQuery() gqlschema.RootQueryResolver {
	return r
}

func (r *resolver) Version(ctx context.Context) (*gqlmodel.Version, error) {
	gql := &gqlmodel.Version{}
	copier.Copy(gql, r.version)
	return gql, nil
}
