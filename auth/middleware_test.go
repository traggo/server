package auth_test

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/traggo/server/auth"
	"github.com/traggo/server/model"
	"github.com/traggo/server/test"
)

func TestMiddleware_noAuthentication_noAuthentication(t *testing.T) {
	db := test.InMemoryDB(t)
	defer db.Close()
	spy := &requestSpy{}

	auth.Middleware(db.DB)(spy).ServeHTTP(httptest.NewRecorder(), httptest.NewRequest("GET", "/test", nil))

	ctx := spy.req.Context()
	assert.Nil(t, auth.GetUser(ctx))
	assert.Nil(t, auth.GetDevice(ctx))
}

func TestMiddleware_query_notExistingToken_noAuthentication(t *testing.T) {
	db := test.InMemoryDB(t)
	defer db.Close()
	spy := &requestSpy{}

	auth.Middleware(db.DB)(spy).ServeHTTP(httptest.NewRecorder(), httptest.NewRequest("GET", "/test?token=abc", nil))

	ctx := spy.req.Context()
	assert.Nil(t, auth.GetUser(ctx))
	assert.Nil(t, auth.GetDevice(ctx))
}

func TestMiddleware_query_validToken_authenticates(t *testing.T) {
	db := test.InMemoryDB(t)
	defer db.Close()
	builder := db.User(1)
	user := builder.User
	device := builder.NewDevice(2, "abc", "test")
	spy := &requestSpy{}

	auth.Middleware(db.DB)(spy).ServeHTTP(httptest.NewRecorder(), httptest.NewRequest("GET", "/test?token=abc", nil))

	ctx := spy.req.Context()
	assert.Equal(t, &user, auth.GetUser(ctx))
	assert.Equal(t, &device, auth.GetDevice(ctx))
}

func TestMiddleware_header_notExistingToken_noAuthentication(t *testing.T) {
	db := test.InMemoryDB(t)
	defer db.Close()
	spy := &requestSpy{}

	request := httptest.NewRequest("GET", "/test?token=abc", nil)
	request.Header["Authorization"] = []string{"traggo abc"}
	auth.Middleware(db.DB)(spy).ServeHTTP(httptest.NewRecorder(), request)

	ctx := spy.req.Context()
	assert.Nil(t, auth.GetUser(ctx))
	assert.Nil(t, auth.GetDevice(ctx))
}

func TestMiddleware_header_validToken_authenticates(t *testing.T) {
	db := test.InMemoryDB(t)
	defer db.Close()
	builder := db.User(1)
	user := builder.User
	device := builder.NewDevice(2, "abc", "test")
	spy := &requestSpy{}

	request := httptest.NewRequest("GET", "/test", nil)
	request.Header.Set("Authorization", "traggo abc")
	auth.Middleware(db.DB)(spy).ServeHTTP(httptest.NewRecorder(), request)

	ctx := spy.req.Context()
	assert.Equal(t, &user, auth.GetUser(ctx))
	assert.Equal(t, &device, auth.GetDevice(ctx))
}

func TestMiddleware_header_validToken_invalidAuthenticationType_noAuthentication(t *testing.T) {
	db := test.InMemoryDB(t)
	defer db.Close()
	db.User(1).NewDevice(2, "abc", "test")
	spy := &requestSpy{}

	request := httptest.NewRequest("GET", "/test", nil)
	request.Header.Set("Authorization", "basic abc")

	auth.Middleware(db.DB)(spy).ServeHTTP(httptest.NewRecorder(), request)

	ctx := spy.req.Context()
	assert.Nil(t, auth.GetUser(ctx))
	assert.Nil(t, auth.GetDevice(ctx))
}

func TestMiddleware_cookie_notExistingToken_noAuthentication(t *testing.T) {
	db := test.InMemoryDB(t)
	defer db.Close()
	spy := &requestSpy{}

	request := httptest.NewRequest("GET", "/test?token=abc", nil)
	request.AddCookie(&http.Cookie{
		Name:  "traggo",
		Value: "abc",
	})
	auth.Middleware(db.DB)(spy).ServeHTTP(httptest.NewRecorder(), request)

	ctx := spy.req.Context()
	assert.Nil(t, auth.GetUser(ctx))
	assert.Nil(t, auth.GetDevice(ctx))
}

func TestMiddleware_cookie_validToken_noAuthentication(t *testing.T) {
	db := test.InMemoryDB(t)
	defer db.Close()
	builder := db.User(1)
	user := builder.User
	device := builder.NewDevice(2, "abc", "test")
	spy := &requestSpy{}

	request := httptest.NewRequest("GET", "/test?token=abc", nil)
	request.AddCookie(&http.Cookie{
		Name:  "traggo",
		Value: "abc",
	})
	auth.Middleware(db.DB)(spy).ServeHTTP(httptest.NewRecorder(), request)

	ctx := spy.req.Context()
	assert.Equal(t, &user, auth.GetUser(ctx))
	assert.Equal(t, &device, auth.GetDevice(ctx))
}

func TestMiddleware_createSession_setsCookie(t *testing.T) {
	db := test.InMemoryDB(t)
	defer db.Close()
	spy := &requestSpy{}
	recorder := httptest.NewRecorder()

	request := httptest.NewRequest("GET", "/test", nil)
	auth.Middleware(db.DB)(spy).ServeHTTP(recorder, request)

	ctx := spy.req.Context()
	auth.GetCreateSession(ctx)("new token", 60)

	cookieHeader := recorder.Header().Get("Set-Cookie")
	assert.Equal(t, `traggo="new token"; Max-Age=60`, cookieHeader)
}

func TestMiddleware_noCallbackExecuted_noCookieSet(t *testing.T) {
	db := test.InMemoryDB(t)
	defer db.Close()
	spy := &requestSpy{}
	recorder := httptest.NewRecorder()

	request := httptest.NewRequest("GET", "/test", nil)
	auth.Middleware(db.DB)(spy).ServeHTTP(recorder, request)

	cookieHeader := recorder.Header().Get("Set-Cookie")
	assert.Equal(t, "", cookieHeader)
}

func TestMiddleware_destroySession_destroysCookie(t *testing.T) {
	db := test.InMemoryDB(t)
	defer db.Close()
	spy := &requestSpy{}
	recorder := httptest.NewRecorder()

	request := httptest.NewRequest("GET", "/test", nil)
	auth.Middleware(db.DB)(spy).ServeHTTP(recorder, request)

	ctx := spy.req.Context()
	auth.GetDestroySession(ctx)()

	cookieHeader := recorder.Header().Get("Set-Cookie")
	assert.Equal(t, "traggo=; Max-Age=0", cookieHeader)
}

func TestGetCreateSession_panicsWhenMiddlewareWasNotExecuted(t *testing.T) {
	assert.Panics(t, func() {
		auth.GetCreateSession(context.Background())
	})
}

func TestGetDestroySession_panicsWhenMiddlewareWasNotExecuted(t *testing.T) {
	assert.Panics(t, func() {
		auth.GetDestroySession(context.Background())
	})
}

func TestMiddleware_panicsWhenDeviceExistButUserDoesNot(t *testing.T) {
	db := test.InMemoryDB(t)
	defer db.Close()
	db.Create(&model.Device{
		ID:     1,
		UserID: 5,
		Token:  "abc",
	})
	spy := &requestSpy{}

	assert.Panics(t, func() {
		auth.Middleware(db.DB)(spy).ServeHTTP(httptest.NewRecorder(), httptest.NewRequest("GET", "/test?token=abc", nil))
	})
}

type requestSpy struct {
	req *http.Request
}

func (s *requestSpy) ServeHTTP(_ http.ResponseWriter, req *http.Request) {
	s.req = req
}
