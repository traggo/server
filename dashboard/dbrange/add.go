package dbrange

import (
	"context"

	"github.com/traggo/server/auth"
	"github.com/traggo/server/dashboard/convert"
	"github.com/traggo/server/dashboard/util"
	"github.com/traggo/server/generated/gqlmodel"
)

// AddDashboardRange adds a dashboard range.
func (r *ResolverForRange) AddDashboardRange(ctx context.Context, dashboardID int, rangeArg gqlmodel.InputNamedDateRange) (*gqlmodel.NamedDateRange, error) {
	userID := auth.GetUser(ctx).ID
	if _, err := util.FindDashboard(r.DB, userID, dashboardID); err != nil {
		return nil, err
	}

	rangeToAdd, err := convert.ToInternalDashboardRange(rangeArg)
	if err != nil {
		return nil, err
	}

	rangeToAdd.DashboardID = dashboardID

	save := r.DB.Create(&rangeToAdd)
	return convert.ToExternalDashboardRange(rangeToAdd), save.Error
}
