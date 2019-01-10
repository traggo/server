package device

import (
	"context"
	"fmt"

	"github.com/jinzhu/copier"
	"github.com/traggo/server/auth"
	"github.com/traggo/server/generated/gqlmodel"
	"github.com/traggo/server/model"
)

// RemoveDevice removes a device
func (r *ResolverForDevice) RemoveDevice(ctx context.Context, id int) (*gqlmodel.Device, error) {
	device := model.Device{ID: id}
	if r.DB.Where(&model.Device{UserID: auth.GetUser(ctx).ID}).Find(&device).RecordNotFound() {
		return nil, fmt.Errorf("device not found")
	}

	remove := r.DB.Delete(&device)
	gqlDevice := &gqlmodel.Device{}
	copier.Copy(gqlDevice, &device)
	return gqlDevice, remove.Error
}
