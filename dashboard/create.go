package dashboard

import (
	"context"

	"github.com/traggo/server/dashboard/convert"

	"github.com/traggo/server/auth"
	"github.com/traggo/server/generated/gqlmodel"
	"github.com/traggo/server/model"
)

// CreateDashboard creates a dashboard.
func (r *ResolverForDashboard) CreateDashboard(ctx context.Context, name string) (*gqlmodel.Dashboard, error) {
	userID := auth.GetUser(ctx).ID
	dashboard := model.Dashboard{
		UserID: userID,
		Name:   name,
	}

	create := r.DB.Create(&dashboard)
	if create.Error != nil {
		return &gqlmodel.Dashboard{}, create.Error
	}

	return convert.ToExternalDashboard(dashboard)
}
