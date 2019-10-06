package entry

import (
	"context"

	"github.com/traggo/server/dashboard/convert"
	"github.com/traggo/server/dashboard/util"

	"github.com/traggo/server/auth"
	"github.com/traggo/server/generated/gqlmodel"
	"github.com/traggo/server/model"
)

// RemoveDashboardEntry removes a dashboard entry.
func (r *ResolverForEntry) RemoveDashboardEntry(ctx context.Context, id int) (*gqlmodel.DashboardEntry, error) {

	userID := auth.GetUser(ctx).ID

	entry, err := util.FindDashboardEntry(r.DB, id)
	if err != nil {
		return nil, err
	}

	if _, err := util.FindDashboard(r.DB, userID, entry.DashboardID); err != nil {
		return nil, err
	}

	remove := r.DB.Delete(&model.DashboardEntry{}, id)
	if remove.Error != nil {
		return &gqlmodel.DashboardEntry{}, remove.Error
	}

	return convert.ToExternalEntry(entry)
}
