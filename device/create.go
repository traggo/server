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

// CreateDevice creates a device.
func (r *ResolverForDevice) CreateDevice(ctx context.Context, username string, pass string, deviceName string, expiresAt model.Time, cookie bool) (*gqlmodel.Login, error) {
	if !expiresAt.Time().After(timeNow()) {
		return nil, errors.New("expiresAt must be in the future")
	}

	user := new(model.User)
	find := r.DB.Where("name = ?", username).Find(user)

	if find.RecordNotFound() {
		return nil, errUserPassWrong
	}

	if !comparePassword(user.Pass, []byte(pass)) {
		return nil, errUserPassWrong
	}

	token := randToken(20)
	for !r.DB.Where("token = ?", token).Find(new(model.Device)).RecordNotFound() {
		token = randToken(20)
	}

	device := &model.Device{
		Token:     token,
		UserID:    user.ID,
		Name:      deviceName,
		ExpiresAt: expiresAt.Time(),
		CreatedAt: timeNow(),
		ActiveAt:  timeNow(),
	}

	if cookie {
		age := int(expiresAt.Time().Sub(timeNow()).Seconds())
		auth.GetCreateSession(ctx)(token, age)
	}

	create := r.DB.Create(device)

	gqlUser := &gqlmodel.User{}
	copier.Copy(gqlUser, user)
	gqlDevice := &gqlmodel.Device{}
	copier.Copy(gqlDevice, device)
	return &gqlmodel.Login{
		User:   *gqlUser,
		Device: *gqlDevice,
		Token:  token}, create.Error
}
