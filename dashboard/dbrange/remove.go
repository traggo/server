package dbrange

import (
	"context"
	"errors"
	"strings"

	"github.com/traggo/server/auth"
	"github.com/traggo/server/dashboard/convert"
	"github.com/traggo/server/dashboard/util"
	"github.com/traggo/server/generated/gqlmodel"
	"github.com/traggo/server/model"
)

// RemoveDashboardRange removes a dashboard range.
func (r *ResolverForRange) RemoveDashboardRange(ctx context.Context, rangeID int) (*gqlmodel.NamedDateRange, error) {
	userID := auth.GetUser(ctx).ID
	dashboardRange, err := util.FindDashboardRange(r.DB, rangeID)
	if err != nil {
		return nil, err
	}

	if _, err := util.FindDashboard(r.DB, userID, dashboardRange.DashboardID); err != nil {
		return nil, err
	}

	entries := []model.DashboardEntry{}
	if err := r.DB.Where(model.DashboardEntry{RangeID: rangeID}).Find(&entries).Error; err != nil {
		return nil, err
	}

	names := []string{}
	for _, entry := range entries {
		names = append(names, entry.Title)
	}

	if len(entries) != 0 {
		return nil, errors.New("range is used in entries: " + strings.Join(names, ","))
	}

	remove := r.DB.Delete(&model.DashboardRange{}, rangeID)
	return convert.ToExternalDashboardRange(dashboardRange), remove.Error
}
