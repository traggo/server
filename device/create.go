package device

import (
	"context"
	"errors"
	"time"

	"github.com/jinzhu/copier"
	"github.com/jinzhu/gorm"
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

	errPasswordLoginDisabled = errors.New("password login is disabled, use single sign-on")
)

// Login creates a device.
func (r *ResolverForDevice) Login(ctx context.Context, username string, pass string, deviceName string, deviceType gqlmodel.DeviceType, cookie bool) (*gqlmodel.Login, error) {
	if r.PasswordLoginDisabled {
		return nil, errPasswordLoginDisabled
	}

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

// CreateDevice creates a device.
func (r *ResolverForDevice) CreateDevice(ctx context.Context, deviceName string, deviceType gqlmodel.DeviceType) (*gqlmodel.Login, error) {

	user := auth.GetUser(ctx)

	return r.createDeviceInternal(ctx, user, deviceName, deviceType, false)
}

func (r *ResolverForDevice) createDeviceInternal(ctx context.Context, user *model.User, deviceName string, deviceType gqlmodel.DeviceType, cookie bool) (*gqlmodel.Login, error) {
	device, err := Create(r.DB, user, deviceName, model.DeviceType(deviceType))
	if err != nil {
		return nil, err
	}

	if cookie {
		auth.GetCreateSession(ctx)(device.Token, device.Type.Seconds())
	}

	gqlUser := &gqlmodel.User{}
	copier.Copy(gqlUser, user)
	gqlDevice := &gqlmodel.Device{}
	copier.Copy(gqlDevice, device)
	return &gqlmodel.Login{
		User:   gqlUser,
		Device: gqlDevice,
		Token:  device.Token}, nil
}

// Create persists a new device with a fresh token for the user.
func Create(db *gorm.DB, user *model.User, deviceName string, deviceType model.DeviceType) (*model.Device, error) {
	if err := deviceType.Valid(); err != nil {
		return nil, err
	}

	token := randToken(20)
	for !db.Where("token = ?", token).Find(new(model.Device)).RecordNotFound() {
		token = randToken(20)
	}

	now := timeNow()
	device := &model.Device{
		Token:     token,
		UserID:    user.ID,
		Name:      deviceName,
		Type:      deviceType,
		CreatedAt: now.UTC(),
		ActiveAt:  now.UTC(),
	}

	if err := db.Create(device).Error; err != nil {
		return nil, err
	}
	return device, nil
}
