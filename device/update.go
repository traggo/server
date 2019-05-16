package device

import (
	"context"
	"errors"

	"github.com/jinzhu/copier"
	"github.com/traggo/server/auth"
	"github.com/traggo/server/generated/gqlmodel"
	"github.com/traggo/server/model"
)

// UpdateDevice updates a device.
func (r *ResolverForDevice) UpdateDevice(ctx context.Context, id int, name string, expiresAt model.Time) (*gqlmodel.Device, error) {
	device := new(model.Device)
	if r.DB.Where("user_id = ?", auth.GetUser(ctx).ID).Find(device, id).RecordNotFound() {
		return nil, errors.New("device not found")
	}

	device.Name = name
	device.ExpiresAt = expiresAt.UTC()
	update := r.DB.Save(device)
	gqlDevice := &gqlmodel.Device{}
	copier.Copy(gqlDevice, device)
	return gqlDevice, update.Error
}
