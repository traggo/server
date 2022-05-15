package auth

import (
	"context"
	"errors"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/gorilla/mux"
	"github.com/jinzhu/gorm"
	"github.com/rs/zerolog/log"
	"github.com/traggo/server/model"
)

type key string

var (
	traggoSession         = "traggo"
	createSessionKey  key = "createSession"
	destroySessionKey key = "destroySession"
	userKey           key = "user"
	deviceKey         key = "device"
)

// Middleware is the auth middleware which sets user and device context parameters.
func Middleware(db *gorm.DB) mux.MiddlewareFunc {
	return func(handler http.Handler) http.Handler {
		return http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
			handler.ServeHTTP(writer, sessionCallbacks(reqisterUser(request, writer, db), writer))
		})
	}
}

func sessionCallbacks(request *http.Request, writer http.ResponseWriter) *http.Request {
	createSession := func(token string, maxAge int) {
		http.SetCookie(writer, &http.Cookie{
			Name:   traggoSession,
			Value:  token,
			MaxAge: maxAge,
		})
	}
	destroySession := func() {
		http.SetCookie(writer, &http.Cookie{
			Name:   traggoSession,
			Value:  "",
			MaxAge: -1,
		})
	}
	return request.WithContext(WithDestroySession(WithCreateSession(request.Context(), createSession), destroySession))
}

func reqisterUser(request *http.Request, writer http.ResponseWriter, db *gorm.DB) *http.Request {
	token, err := getToken(request)
	if err != nil {
		msg := "Issue fetching token"
		log.Info().Msg(msg)
		http.Error(writer, msg, 400)
		return request
	}

	device := &model.Device{}
	if db.Where("token = ?", token).Find(device).RecordNotFound() {
		msg := "No device with token found"
		log.Info().Str("token", token).Msg(msg)
		http.Error(writer, msg, 400)
		return request
	}

	user := &model.User{}
	if db.Find(user, device.UserID).RecordNotFound() {
		log.Panic().Int("userID", device.UserID).Int("deviceID", device.ID).Msg("User not found")
	}
	log.Info().Int("userid", device.UserID).Str("username", user.Name).Msg("User found")

	impersonate := request.Header.Get("X-Traggo-Impersonate")
	if impersonate != "" {
		if !user.Admin {
			msg := "Trying to impersonate without being admin"
			log.Info().Str("impersonate", impersonate).Msg(msg)
			http.Error(writer, msg, 403)
			return request
		}
		userID, err := strconv.Atoi(impersonate)
		if err != nil {
			msg := "Unable to parse impersonation header"
			log.Info().Str("impersonate", impersonate).Msg(msg)
			http.Error(writer, msg, 400)
			return request
		}
		impersonateUser := &model.User{}
		if db.Find(impersonateUser, userID).RecordNotFound() {
			msg := "Impersonation user not found"
			log.Info().Int("userid", userID).Msg(msg)
			http.Error(writer, msg, 400)
			return request
		}
		user = impersonateUser
		log.Info().Str("username", user.Name).Msg("Impersonation user")
	}

	if device.ActiveAt.Before(time.Now().Add(5 * -time.Minute)) {
		log.Debug().Int("deviceId", device.ID).Str("deviceName", device.Name).Msg("update device activeAt")
		device.ActiveAt = timeNow()
		db.Save(device)

		if cookie, err := request.Cookie(traggoSession); err == nil && cookie != nil {
			cookie.MaxAge = device.Type.Seconds()
			http.SetCookie(writer, cookie)
		}
	}

	return request.WithContext(WithUser(WithDevice(request.Context(), device), user))
}

func getToken(request *http.Request) (string, error) {
	if value := request.Header.Get("Authorization"); value != "" && strings.HasPrefix(value, "traggo ") {
		log.Debug().Msg("Using authorization header for authentication")
		key := strings.TrimPrefix(value, "traggo ")
		return key, nil
	}

	if cookie, err := request.Cookie("traggo"); err == nil && cookie != nil && cookie.Value != "" {
		log.Debug().Msg("Using cookie for authentication")
		return cookie.Value, nil
	}

	if token := request.FormValue("token"); token != "" {
		log.Debug().Msg("Using query parameter for authentication")
		return token, nil
	}

	return "", errors.New("no token found")
}

// WithDevice adds the authenticated device to the context.
func WithDevice(ctx context.Context, device *model.Device) context.Context {
	return context.WithValue(
		ctx,
		deviceKey,
		device)
}

// WithUser adds the authenticated user to the context.
func WithUser(ctx context.Context, user *model.User) context.Context {
	return context.WithValue(
		ctx,
		userKey,
		user)
}

// WithCreateSession adds the create session method
func WithCreateSession(ctx context.Context, f func(token string, age int)) context.Context {
	return context.WithValue(ctx, createSessionKey, f)
}

// WithDestroySession adds the destroy session method
func WithDestroySession(ctx context.Context, f func()) context.Context {
	return context.WithValue(ctx, destroySessionKey, f)
}

// GetUser returns the authenticated user or panics.
func GetUser(ctx context.Context) *model.User {
	if user, ok := ctx.Value(userKey).(*model.User); ok {
		return user
	}
	return nil
}

// GetDevice returns the authenticated device or nil.
func GetDevice(ctx context.Context) *model.Device {
	if device, ok := ctx.Value(deviceKey).(*model.Device); ok {
		return device
	}
	return nil
}

// GetCreateSession gets the create session callback
func GetCreateSession(ctx context.Context) func(token string, age int) {
	if f, ok := ctx.Value(createSessionKey).(func(string, int)); ok {
		return f
	}
	panic("create session must exist")
}

// GetDestroySession gets the destroy session callback
func GetDestroySession(ctx context.Context) func() {
	if f, ok := ctx.Value(destroySessionKey).(func()); ok {
		return f
	}
	panic("create session must exist")
}
