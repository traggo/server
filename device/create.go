package device

import (
	"context"
	"errors"
	"time"

	"github.com/jinzhu/copier"
	"github.com/traggo/server/auth"
	"github.com/traggo/server/auth/rand"
	"github.com/traggo/server/generated/gqlmodel"
	"github.com/traggo/server/model"
	"github.com/traggo/server/user/password"
)

var (
	timeNow          = time.Now
	randToken        = rand.Token
	comparePassword  = password.ComparePassword
	errUserPassWrong = errors.New("username/password combination does not exist")
)

// Login creates a device.
func (r *ResolverForDevice) Login(ctx context.Context, username string, pass string, deviceName string, deviceType gqlmodel.DeviceType, cookie bool) (*gqlmodel.Login, error) {

	user := new(model.User)
	find := r.DB.Where("name = ?", username).Find(user)

	if find.RecordNotFound() {
		return nil, errUserPassWrong
	}

	if !comparePassword(user.Pass, []byte(pass)) {
		return nil, errUserPassWrong
	}

	return r.createDeviceInternal(ctx, user, deviceName, deviceType, cookie)
}


// Impersonate creates a device.
func (r *ResolverForDevice) Impersonate(ctx context.Context, username string, deviceName string, deviceType gqlmodel.DeviceType, cookie bool) (*gqlmodel.Login, error) {

	current_user := auth.GetUser(ctx)

    if current_user.Admin == false {
		return nil, errors.New("needs to be admin to impersonate")
    }

	user := new(model.User)
	find := r.DB.Where("name = ?", username).Find(user)
	if find.RecordNotFound() {
		return nil, errors.New("username does not exist")
	}

	return r.createDeviceInternal(ctx, user, deviceName, deviceType, cookie)
}

// CreateDevice creates a device.
func (r *ResolverForDevice) CreateDevice(ctx context.Context, deviceName string, deviceType gqlmodel.DeviceType) (*gqlmodel.Login, error) {

	user := auth.GetUser(ctx)

	return r.createDeviceInternal(ctx, user, deviceName, deviceType, false)
}

func (r *ResolverForDevice) createDeviceInternal(ctx context.Context, user *model.User, deviceName string, deviceType gqlmodel.DeviceType, cookie bool) (*gqlmodel.Login, error) {

	token := randToken(20)
	for !r.DB.Where("token = ?", token).Find(new(model.Device)).RecordNotFound() {
		token = randToken(20)
	}

	now := timeNow()
	device := &model.Device{
		Token:     token,
		UserID:    user.ID,
		Name:      deviceName,
		Type:      model.DeviceType(deviceType),
		CreatedAt: now.UTC(),
		ActiveAt:  now.UTC(),
	}

	if err := device.Type.Valid(); err != nil {
		return nil, err
	}

	if cookie {
		auth.GetCreateSession(ctx)(token, device.Type.Seconds())
	}

	create := r.DB.Create(device)

	gqlUser := &gqlmodel.User{}
	copier.Copy(gqlUser, user)
	gqlDevice := &gqlmodel.Device{}
	copier.Copy(gqlDevice, device)
	return &gqlmodel.Login{
		User:   gqlUser,
		Device: gqlDevice,
		Token:  token}, create.Error
}
