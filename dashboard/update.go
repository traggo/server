package dashboard

import (
	"context"

	"github.com/traggo/server/dashboard/convert"

	"github.com/traggo/server/dashboard/util"

	"github.com/traggo/server/auth"
	"github.com/traggo/server/generated/gqlmodel"
	"github.com/traggo/server/model"
)

// UpdateDashboard updates a dashboard.
func (r *ResolverForDashboard) UpdateDashboard(ctx context.Context, id int, name string) (*gqlmodel.Dashboard, error) {
	userID := auth.GetUser(ctx).ID

	dashboard, err := util.FindDashboard(r.DB, userID, id)
	if err != nil {
		return nil, err
	}

	dashboard.Name = name

	save := r.DB.Save(dashboard)

	if save.Error != nil {
		return nil, save.Error
	}

	var entries []model.DashboardEntry
	if err := r.DB.Where(&model.DashboardEntry{DashboardID: dashboard.ID}).Find(&entries).Error; err != nil {
		return nil, err
	}
	dashboard.Entries = entries

	var ranges []model.DashboardRange
	if err := r.DB.Where(&model.DashboardRange{DashboardID: dashboard.ID}).Find(&ranges).Error; err != nil {
		return nil, err
	}
	dashboard.Ranges = ranges

	return convert.ToExternalDashboard(dashboard)
}
