package device

import (
	"bytes"
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/traggo/server/auth"
	"github.com/traggo/server/generated/gqlmodel"
	"github.com/traggo/server/model"
	"github.com/traggo/server/test"
)

var (
	existingPassword = []byte{1}
)

func TestGQL_CreateDevice_fails_userDoesNotExist(t *testing.T) {
	db := test.InMemoryDB(t)
	defer db.Close()

	now := test.Time("2018-06-30T18:30:00+02:00")
	timeDispose := fakeTime(now)
	defer timeDispose()

	resolver := ResolverForDevice{DB: db.DB}
	login, err := resolver.CreateDevice(
		context.Background(),
		"jmattheis",
		"123",
		"test",
		model.Time(test.Time("2019-06-30T18:30:00+02:00")),

		false)

	assert.Nil(t, login)
	assert.Equal(t, errUserPassWrong, err)
	assertDeviceCount(t, db, 0)
}

func TestGQL_CreateDevice_fails_ExpireAtAlreadyExpired(t *testing.T) {
	db := test.InMemoryDB(t)
	defer db.Close()

	now := test.Time("2018-06-30T18:30:00+02:00")
	timeDispose := fakeTime(now)
	defer timeDispose()

	resolver := ResolverForDevice{DB: db.DB}
	login, err := resolver.CreateDevice(
		context.Background(),
		"jmattheis",
		"123",
		"test",
		model.Time(test.Time("2012-06-30T18:30:00+02:00")),

		false)

	assert.Nil(t, login)
	assert.EqualError(t, err, "expiresAt must be in the future")
	assertDeviceCount(t, db, 0)
}

func TestGQL_CreateDevice_fails_wrongPass(t *testing.T) {
	db := test.InMemoryDB(t)
	defer db.Close()
	db.Create(&model.User{
		Name: "jmattheis",
		ID:   1,
		Pass: existingPassword,
	})
	pwDispose := fakePassword()
	defer pwDispose()

	now := test.Time("2018-06-30T18:30:00+02:00")
	timeDispose := fakeTime(now)
	defer timeDispose()

	resolver := ResolverForDevice{DB: db.DB}
	login, err := resolver.CreateDevice(
		context.Background(),
		"jmattheis",
		"123",
		"test",
		model.Time(test.Time("2019-06-30T18:30:00+02:00")),

		false)

	assert.Nil(t, login)
	assert.Equal(t, errUserPassWrong, err)
	assertDeviceCount(t, db, 0)
}

func TestGQL_CreateDevice_succeeds(t *testing.T) {
	db := test.InMemoryDB(t)
	defer db.Close()
	db.Create(&model.User{
		Name: "jmattheis",
		ID:   1,
		Pass: existingPassword,
	})
	pwDispose := fakePassword()
	defer pwDispose()

	tokenDispose := fakeToken("firstToken")
	defer tokenDispose()

	now := test.Time("2018-06-30T18:30:00+02:00")
	timeDispose := fakeTime(now)
	defer timeDispose()

	resolver := ResolverForDevice{DB: db.DB}
	expireDate := test.Time("2019-06-30T18:30:00+02:00")
	login, err := resolver.CreateDevice(
		context.Background(),
		"jmattheis",
		"unicorn",
		"test",
		model.Time(expireDate),
		false)

	assert.Nil(t, err)

	expected := &gqlmodel.Login{
		Token: "firstToken",
		User:  gqlmodel.User{Admin: false, ID: 1, Name: "jmattheis"},
		Device: gqlmodel.Device{
			ID:        1,
			Name:      "test",
			ExpiresAt: model.Time(expireDate),
			CreatedAt: model.Time(now),
			ActiveAt:  model.Time(now),
		},
	}
	assert.Equal(t, expected, login)
	assertDeviceCount(t, db, 1)
}

func TestGQL_CreateDevice_setsCookie(t *testing.T) {
	db := test.InMemoryDB(t)
	defer db.Close()
	db.Create(&model.User{
		Name: "jmattheis",
		ID:   1,
		Pass: existingPassword,
	})
	pwDispose := fakePassword()
	defer pwDispose()

	tokenDispose := fakeToken("firstToken")
	defer tokenDispose()

	now := test.Time("2018-06-30T18:30:00+02:00")
	timeDispose := fakeTime(now)
	defer timeDispose()

	resolver := ResolverForDevice{DB: db.DB}
	expireDate := test.Time("2018-07-01T18:30:00+02:00")

	createSessionToken := ""
	createSessionAge := 0

	_, err := resolver.CreateDevice(
		auth.WithCreateSession(context.Background(), func(token string, age int) {
			createSessionToken = token
			createSessionAge = age
		}),
		"jmattheis",
		"unicorn",
		"test",
		model.Time(expireDate),
		true)

	assert.Nil(t, err)

	assert.Equal(t, "firstToken", createSessionToken)
	assert.Equal(t, 60*60*24, createSessionAge)
	assertDeviceCount(t, db, 1)
}

func TestGQL_CreateDevice_succeeds_withExistingToken(t *testing.T) {
	db := test.InMemoryDB(t)
	defer db.Close()
	db.Create(&model.User{
		Name: "jmattheis",
		ID:   1,
		Pass: existingPassword,
	})
	db.Create(&model.Device{
		Name:   "test",
		ID:     1,
		Token:  "firstToken",
		UserID: 1,
	})
	pwDispose := fakePassword()
	defer pwDispose()

	tokenDispose := fakeToken("firstToken", "secondToken")
	defer tokenDispose()

	now := test.Time("2018-06-30T18:30:00+02:00")
	timeDispose := fakeTime(now)
	defer timeDispose()

	resolver := ResolverForDevice{DB: db.DB}
	expireDate := test.Time("2022-06-30T18:30:00+02:00")
	login, err := resolver.CreateDevice(
		context.Background(),
		"jmattheis",
		"unicorn",
		"test",
		model.Time(expireDate),
		false)

	assert.Nil(t, err)

	assert.Equal(t, "secondToken", login.Token)
	assertDeviceCount(t, db, 2)
}

func fakePassword() func() {
	old := comparePassword
	comparePassword = func(hashedPassword, password []byte) bool {
		return bytes.Equal(password, []byte("unicorn")) && bytes.Equal(hashedPassword, existingPassword)
	}
	return func() {
		comparePassword = old
	}
}

func fakeTime(t time.Time) func() {
	old := timeNow
	timeNow = func() time.Time {
		return t
	}
	return func() {
		timeNow = old
	}
}

func fakeToken(token ...string) func() {
	old := randToken
	remaining := token
	randToken = func(n int) string {
		if len(remaining) > 0 {
			used := remaining[0]
			remaining = remaining[1:]
			return used
		}
		panic("oops no token specified")
	}
	return func() {
		randToken = old
	}
}
