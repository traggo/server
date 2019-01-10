package device

import (
	"context"

	"github.com/jinzhu/copier"
	"github.com/traggo/server/auth"
	"github.com/traggo/server/generated/gqlmodel"
	"github.com/traggo/server/model"
)

// Devices returns all devices.
func (r *ResolverForDevice) Devices(ctx context.Context) ([]gqlmodel.Device, error) {
	user := auth.GetUser(ctx)
	var devices []model.Device
	find := r.DB.Where(&model.Device{UserID: user.ID}).Find(&devices)
	var result []gqlmodel.Device
	copier.Copy(&result, &devices)
	return result, find.Error
}
