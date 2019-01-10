package device

import (
	"context"

	"github.com/jinzhu/copier"
	"github.com/traggo/server/auth"
	"github.com/traggo/server/generated/gqlmodel"
)

// CurrentDevice returns the current device.
func (r *ResolverForDevice) CurrentDevice(ctx context.Context) (*gqlmodel.Device, error) {
	device := auth.GetDevice(ctx)
	if device == nil {
		return nil, nil
	}
	var result gqlmodel.Device
	copier.Copy(&result, device)
	return &result, nil
}
