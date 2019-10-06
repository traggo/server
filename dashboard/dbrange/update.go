package dbrange

import (
	"context"

	"github.com/traggo/server/auth"
	"github.com/traggo/server/dashboard/convert"
	"github.com/traggo/server/dashboard/util"
	"github.com/traggo/server/generated/gqlmodel"
)

// UpdateDashboardRange updates a dashboard range.
func (r *ResolverForRange) UpdateDashboardRange(ctx context.Context, rangeID int, rangeArg gqlmodel.InputNamedDateRange) (*gqlmodel.NamedDateRange, error) {
	userID := auth.GetUser(ctx).ID

	dashboardRange, err := util.FindDashboardRange(r.DB, rangeID)
	if err != nil {
		return nil, err
	}

	if _, err := util.FindDashboard(r.DB, userID, dashboardRange.DashboardID); err != nil {
		return nil, err
	}

	rangeToUpdate, err := convert.ToInternalDashboardRange(rangeArg)
	if err != nil {
		return nil, err
	}

	rangeToUpdate.DashboardID = dashboardRange.DashboardID
	rangeToUpdate.ID = dashboardRange.ID

	save := r.DB.Save(rangeToUpdate)
	return convert.ToExternalDashboardRange(rangeToUpdate), save.Error
}
