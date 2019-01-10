package device

import (
	"context"

	"github.com/jinzhu/copier"
	"github.com/traggo/server/auth"
	"github.com/traggo/server/generated/gqlmodel"
)

// RemoveCurrentDevice removes the current authenticated device
func (r *ResolverForDevice) RemoveCurrentDevice(ctx context.Context) (gqlmodel.Device, error) {
	device := auth.GetDevice(ctx)
	remove := r.DB.Delete(device)
	gqlDevice := &gqlmodel.Device{}
	copier.Copy(gqlDevice, device)
	return *gqlDevice, remove.Error
}
