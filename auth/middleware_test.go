package auth

import (
	"context"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/traggo/server/test"
)

func TestMiddleware_noAuthentication_noAuthentication(t *testing.T) {
	db := test.InMemoryDB(t)
	defer db.Close()
	spy := &requestSpy{}

	Middleware(db.DB)(spy).ServeHTTP(httptest.NewRecorder(), httptest.NewRequest("GET", "/test", nil))

	ctx := spy.req.Context()
	assert.Nil(t, GetUser(ctx))
	assert.Nil(t, GetDevice(ctx))
}

func TestMiddleware_query_notExistingToken_noAuthentication(t *testing.T) {
	db := test.InMemoryDB(t)
	defer db.Close()
	spy := &requestSpy{}

	Middleware(db.DB)(spy).ServeHTTP(httptest.NewRecorder(), httptest.NewRequest("GET", "/test?token=abc", nil))

	ctx := spy.req.Context()
	assert.Nil(t, GetUser(ctx))
	assert.Nil(t, GetDevice(ctx))
}

func TestMiddleware_query_validToken_authenticates(t *testing.T) {
	now := test.Time("2018-06-30T18:30:00Z")
	timeDispose := fakeTime(now)
	defer timeDispose()

	db := test.InMemoryDB(t)
	defer db.Close()
	builder := db.User(1)
	user := builder.User
	device := builder.NewDevice(2, "abc", "test")
	spy := &requestSpy{}

	Middleware(db.DB)(spy).ServeHTTP(httptest.NewRecorder(), httptest.NewRequest("GET", "/test?token=abc", nil))

	ctx := spy.req.Context()
	assert.Equal(t, &user, GetUser(ctx))

	device.ActiveAt = now
	assert.Equal(t, &device, GetDevice(ctx))
}

func TestMiddleware_header_notExistingToken_noAuthentication(t *testing.T) {
	db := test.InMemoryDB(t)
	defer db.Close()
	spy := &requestSpy{}

	request := httptest.NewRequest("GET", "/test?token=abc", nil)
	request.Header["Authorization"] = []string{"traggo abc"}
	Middleware(db.DB)(spy).ServeHTTP(httptest.NewRecorder(), request)

	ctx := spy.req.Context()
	assert.Nil(t, GetUser(ctx))
	assert.Nil(t, GetDevice(ctx))
}

func TestMiddleware_header_validToken_authenticates(t *testing.T) {
	now := test.Time("2018-06-30T18:30:00Z")
	timeDispose := fakeTime(now)
	defer timeDispose()

	db := test.InMemoryDB(t)
	defer db.Close()
	builder := db.User(1)
	user := builder.User
	device := builder.NewDevice(2, "abc", "test")
	spy := &requestSpy{}

	request := httptest.NewRequest("GET", "/test", nil)
	request.Header.Set("Authorization", "traggo abc")
	Middleware(db.DB)(spy).ServeHTTP(httptest.NewRecorder(), request)

	ctx := spy.req.Context()
	assert.Equal(t, &user, GetUser(ctx))

	device.ActiveAt = now
	assert.Equal(t, &device, GetDevice(ctx))
}

func TestMiddleware_header_validToken_invalidAuthenticationType_noAuthentication(t *testing.T) {
	db := test.InMemoryDB(t)
	defer db.Close()
	db.User(1).NewDevice(2, "abc", "test")
	spy := &requestSpy{}

	request := httptest.NewRequest("GET", "/test", nil)
	request.Header.Set("Authorization", "basic abc")

	Middleware(db.DB)(spy).ServeHTTP(httptest.NewRecorder(), request)

	ctx := spy.req.Context()
	assert.Nil(t, GetUser(ctx))
	assert.Nil(t, GetDevice(ctx))
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
	Middleware(db.DB)(spy).ServeHTTP(httptest.NewRecorder(), request)

	ctx := spy.req.Context()
	assert.Nil(t, GetUser(ctx))
	assert.Nil(t, GetDevice(ctx))
}

func TestMiddleware_cookie_validToken_authenticates(t *testing.T) {
	now := test.Time("2018-06-30T18:30:00Z")
	timeDispose := fakeTime(now)
	defer timeDispose()

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
	Middleware(db.DB)(spy).ServeHTTP(httptest.NewRecorder(), request)

	ctx := spy.req.Context()
	assert.Equal(t, &user, GetUser(ctx))

	device.ActiveAt = now
	assert.Equal(t, &device, GetDevice(ctx))
}

func TestMiddleware_createSession_setsCookie(t *testing.T) {
	db := test.InMemoryDB(t)
	defer db.Close()
	spy := &requestSpy{}
	recorder := httptest.NewRecorder()

	request := httptest.NewRequest("GET", "/test", nil)
	Middleware(db.DB)(spy).ServeHTTP(recorder, request)

	ctx := spy.req.Context()
	GetCreateSession(ctx)("new token", 60)

	cookieHeader := recorder.Header().Get("Set-Cookie")
	assert.Equal(t, `traggo="new token"; Max-Age=60`, cookieHeader)
}

func TestMiddleware_noCallbackExecuted_noCookieSet(t *testing.T) {
	db := test.InMemoryDB(t)
	defer db.Close()
	spy := &requestSpy{}
	recorder := httptest.NewRecorder()

	request := httptest.NewRequest("GET", "/test", nil)
	Middleware(db.DB)(spy).ServeHTTP(recorder, request)

	cookieHeader := recorder.Header().Get("Set-Cookie")
	assert.Equal(t, "", cookieHeader)
}

func TestMiddleware_destroySession_destroysCookie(t *testing.T) {
	db := test.InMemoryDB(t)
	defer db.Close()
	spy := &requestSpy{}
	recorder := httptest.NewRecorder()

	request := httptest.NewRequest("GET", "/test", nil)
	Middleware(db.DB)(spy).ServeHTTP(recorder, request)

	ctx := spy.req.Context()
	GetDestroySession(ctx)()

	cookieHeader := recorder.Header().Get("Set-Cookie")
	assert.Equal(t, "traggo=; Max-Age=0", cookieHeader)
}

func TestMiddleware_impersonate_no_admin(t *testing.T) {
	now := test.Time("2018-06-30T18:30:00Z")
	timeDispose := fakeTime(now)
	defer timeDispose()

	db := test.InMemoryDB(t)
	defer db.Close()
	builder := db.User(1)
	builder.NewDevice(2, "abc", "test")
	spy := &requestSpy{}
	recorder := httptest.NewRecorder()

	request := httptest.NewRequest("GET", "/test?token=abc", nil)
	request.Header.Set("X-Traggo-Impersonate", "empty")
	Middleware(db.DB)(spy).ServeHTTP(recorder, request)

	response := recorder.Result()
	assert.Equal(t, 403, response.StatusCode)

	bodyBytes, _ := io.ReadAll(response.Body)
	bodyString := string(bodyBytes)
	assert.Equal(t, "Trying to impersonate without being admin\n", bodyString)

	ctx := spy.req.Context()
	assert.Nil(t, GetUser(ctx))
	assert.Nil(t, GetDevice(ctx))
}

func TestMiddleware_impersonate_invalid_personate_header(t *testing.T) {
	now := test.Time("2018-06-30T18:30:00Z")
	timeDispose := fakeTime(now)
	defer timeDispose()

	db := test.InMemoryDB(t)
	defer db.Close()
	builder := db.User(1)
	user := builder.User
	user.Admin = true
	db.Save(user)
	builder.NewDevice(2, "abc", "test")
	spy := &requestSpy{}
	recorder := httptest.NewRecorder()

	request := httptest.NewRequest("GET", "/test?token=abc", nil)
	request.Header.Set("X-Traggo-Impersonate", "invalid")
	Middleware(db.DB)(spy).ServeHTTP(recorder, request)

	response := recorder.Result()
	assert.Equal(t, 400, response.StatusCode)

	bodyBytes, _ := io.ReadAll(response.Body)
	bodyString := string(bodyBytes)
	assert.Equal(t, "Unable to parse impersonation header\n", bodyString)

	ctx := spy.req.Context()
	assert.Nil(t, GetUser(ctx))
	assert.Nil(t, GetDevice(ctx))
}

func TestMiddleware_impersonate_non_existing_user(t *testing.T) {
	now := test.Time("2018-06-30T18:30:00Z")
	timeDispose := fakeTime(now)
	defer timeDispose()

	db := test.InMemoryDB(t)
	defer db.Close()
	builder := db.User(1)
	user := builder.User
	user.Admin = true
	db.Save(user)
	builder.NewDevice(2, "abc", "test")
	spy := &requestSpy{}
	recorder := httptest.NewRecorder()

	request := httptest.NewRequest("GET", "/test?token=abc", nil)
	request.Header.Set("X-Traggo-Impersonate", "42")
	Middleware(db.DB)(spy).ServeHTTP(recorder, request)

	response := recorder.Result()
	assert.Equal(t, 400, response.StatusCode)

	bodyBytes, _ := io.ReadAll(response.Body)
	bodyString := string(bodyBytes)
	assert.Equal(t, "Impersonation user not found\n", bodyString)

	ctx := spy.req.Context()
	assert.Nil(t, GetUser(ctx))
	assert.Nil(t, GetDevice(ctx))
}

func TestMiddleware_impersonate_happy(t *testing.T) {
	now := test.Time("2018-06-30T18:30:00Z")
	timeDispose := fakeTime(now)
	defer timeDispose()

	db := test.InMemoryDB(t)
	defer db.Close()
	admin_builder := db.User(1)
	admin_user := admin_builder.User
	admin_user.Admin = true
	db.Save(admin_user)
	admin_device := admin_builder.NewDevice(2, "abc", "test")

	builder := db.User(2)
	user := builder.User

	spy := &requestSpy{}
	recorder := httptest.NewRecorder()

	request := httptest.NewRequest("GET", "/test?token=abc", nil)
	request.Header.Set("X-Traggo-Impersonate", "2")
	Middleware(db.DB)(spy).ServeHTTP(recorder, request)

	response := recorder.Result()
	assert.Equal(t, 200, response.StatusCode)

	ctx := spy.req.Context()
	assert.Equal(t, &user, GetUser(ctx))

	admin_device.ActiveAt = now
	assert.Equal(t, &admin_device, GetDevice(ctx))
}

func TestGetCreateSession_panicsWhenMiddlewareWasNotExecuted(t *testing.T) {
	assert.Panics(t, func() {
		GetCreateSession(context.Background())
	})
}

func TestGetDestroySession_panicsWhenMiddlewareWasNotExecuted(t *testing.T) {
	assert.Panics(t, func() {
		GetDestroySession(context.Background())
	})
}

type requestSpy struct {
	req *http.Request
}

func (s *requestSpy) ServeHTTP(_ http.ResponseWriter, req *http.Request) {
	s.req = req
}
