package dashboard

import (
	"context"

	"github.com/traggo/server/dashboard/convert"

	"github.com/traggo/server/auth"
	"github.com/traggo/server/generated/gqlmodel"
	"github.com/traggo/server/model"
)

// Dashboards returns all dashboards.
func (r *ResolverForDashboard) Dashboards(ctx context.Context) ([]*gqlmodel.Dashboard, error) {
	userID := auth.GetUser(ctx).ID

	dashboards := []model.Dashboard{}

	q := r.DB
	q = q.Preload("Entries")
	q = q.Preload("Entries.TagFilters")
	q = q.Preload("Ranges")

	find := q.Where(&model.Dashboard{UserID: userID}).Find(&dashboards)

	if find.Error != nil {
		return nil, find.Error
	}

	return convert.ToExternalDashboards(dashboards)
}
