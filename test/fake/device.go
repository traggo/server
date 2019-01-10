package fake

import (
	"context"

	"github.com/traggo/server/auth"
	"github.com/traggo/server/model"
)

// Device creates a context with a fake device.
func Device(device *model.Device) context.Context {
	return auth.WithDevice(context.Background(), device)
}
