package oidc

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"
	"time"

	"github.com/go-jose/go-jose/v4"
	"github.com/go-jose/go-jose/v4/jwt"
	"github.com/gorilla/mux"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/traggo/server/auth"
	"github.com/traggo/server/config"
	"github.com/traggo/server/model"
	"github.com/traggo/server/test"
)

const (
	clientID     = "traggo-test"
	clientSecret = "synthetic-secret"
)

// mockIDP is a minimal OIDC provider: discovery, jwks and token endpoint.
type mockIDP struct {
	srv    *httptest.Server
	key    *rsa.PrivateKey
	claims map[string]interface{}
	// overrides for negative tests
	nonceOverride    string
	audienceOverride string
	skipIDToken      bool
	signKey          *rsa.PrivateKey

	codes map[string]authRequest
}

type authRequest struct {
	challenge string
	nonce     string
}

func newMockIDP(t *testing.T) *mockIDP {
	key, err := rsa.GenerateKey(rand.Reader, 2048)
	require.NoError(t, err)
	m := &mockIDP{key: key, codes: map[string]authRequest{}, claims: map[string]interface{}{"preferred_username": "alice"}}
	mux := http.NewServeMux()
	mux.HandleFunc("/.well-known/openid-configuration", func(w http.ResponseWriter, r *http.Request) {
		_ = json.NewEncoder(w).Encode(map[string]interface{}{
			"issuer":                                m.srv.URL,
			"authorization_endpoint":                m.srv.URL + "/authorize",
			"token_endpoint":                        m.srv.URL + "/token",
			"jwks_uri":                              m.srv.URL + "/jwks",
			"id_token_signing_alg_values_supported": []string{"RS256"},
		})
	})
	mux.HandleFunc("/jwks", func(w http.ResponseWriter, r *http.Request) {
		_ = json.NewEncoder(w).Encode(jose.JSONWebKeySet{Keys: []jose.JSONWebKey{
			{Key: &key.PublicKey, KeyID: "k1", Algorithm: "RS256", Use: "sig"},
		}})
	})
	mux.HandleFunc("/token", m.token)
	m.srv = httptest.NewServer(mux)
	t.Cleanup(m.srv.Close)
	return m
}

func (m *mockIDP) signingKey() *rsa.PrivateKey {
	if m.signKey != nil {
		return m.signKey
	}
	return m.key
}

func (m *mockIDP) token(w http.ResponseWriter, r *http.Request) {
	_ = r.ParseForm()
	user, pass, basic := r.BasicAuth()
	if r.Form.Get("client_id") != "" {
		user, pass = r.Form.Get("client_id"), r.Form.Get("client_secret")
	} else if basic {
		user, _ = url.QueryUnescape(user)
		pass, _ = url.QueryUnescape(pass)
	}
	req, ok := m.codes[r.Form.Get("code")]
	delete(m.codes, r.Form.Get("code")) // codes are single use
	sum := sha256.Sum256([]byte(r.Form.Get("code_verifier")))
	if !ok || user != clientID || pass != clientSecret ||
		r.Form.Get("grant_type") != "authorization_code" ||
		base64.RawURLEncoding.EncodeToString(sum[:]) != req.challenge {
		w.WriteHeader(http.StatusBadRequest)
		_, _ = w.Write([]byte(`{"error":"invalid_grant"}`))
		return
	}
	resp := map[string]interface{}{"access_token": "at", "token_type": "Bearer", "expires_in": 60}
	if !m.skipIDToken {
		nonce, aud := req.nonce, clientID
		if m.nonceOverride != "" {
			nonce = m.nonceOverride
		}
		if m.audienceOverride != "" {
			aud = m.audienceOverride
		}
		signer, err := jose.NewSigner(jose.SigningKey{Algorithm: jose.RS256, Key: jose.JSONWebKey{Key: m.signingKey(), KeyID: "k1"}}, (&jose.SignerOptions{}).WithType("JWT"))
		if err != nil {
			panic(err)
		}
		claims := map[string]interface{}{"nonce": nonce}
		for k, v := range m.claims {
			claims[k] = v
		}
		raw, err := jwt.Signed(signer).Claims(jwt.Claims{
			Issuer:   m.srv.URL,
			Subject:  "sub-1",
			Audience: jwt.Audience{aud},
			Expiry:   jwt.NewNumericDate(time.Now().Add(time.Minute)),
			IssuedAt: jwt.NewNumericDate(time.Now()),
		}).Claims(claims).Serialize()
		if err != nil {
			panic(err)
		}
		resp["id_token"] = raw
	}
	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(resp)
}

// authorize plays the browser and identity provider: it validates the authorization request and returns a code.
func (m *mockIDP) authorize(t *testing.T, location string) string {
	u, err := url.Parse(location)
	require.NoError(t, err)
	assert.Equal(t, m.srv.URL+"/authorize", u.Scheme+"://"+u.Host+u.Path)
	q := u.Query()
	assert.Equal(t, "code", q.Get("response_type"))
	assert.Equal(t, clientID, q.Get("client_id"))
	assert.Equal(t, "S256", q.Get("code_challenge_method"))
	assert.NotEmpty(t, q.Get("code_challenge"))
	assert.NotEmpty(t, q.Get("state"))
	assert.NotEmpty(t, q.Get("nonce"))
	assert.Contains(t, q.Get("scope"), "openid")
	code := "code-" + q.Get("state")
	m.codes[code] = authRequest{challenge: q.Get("code_challenge"), nonce: q.Get("nonce")}
	return code
}

type env struct {
	t      *testing.T
	idp    *mockIDP
	db     *test.Database
	router *mux.Router
}

func newEnv(t *testing.T, mod func(*config.OIDC)) *env {
	idp := newMockIDP(t)
	db := test.InMemoryDB(t)
	t.Cleanup(func() { db.Close() })
	conf := config.OIDC{
		Issuer:        idp.srv.URL,
		ClientID:      clientID,
		ClientSecret:  clientSecret,
		RedirectURL:   "https://traggo.example.com/auth/oidc/callback",
		Scopes:        []string{"openid", "profile", "email"},
		UsernameClaim: "preferred_username",
	}
	if mod != nil {
		mod(&conf)
	}
	router := mux.NewRouter()
	router.Use(auth.Middleware(db.DB))
	require.NoError(t, Register(router, db.DB, conf, 4))
	return &env{t: t, idp: idp, db: db, router: router}
}

func (e *env) do(req *http.Request) *httptest.ResponseRecorder {
	rec := httptest.NewRecorder()
	e.router.ServeHTTP(rec, req)
	return rec
}

// start calls the login endpoint and returns the flow cookie and the redirect location.
func (e *env) start() (*http.Cookie, string) {
	rec := e.do(httptest.NewRequest("GET", LoginPath, nil))
	require.Equal(e.t, http.StatusFound, rec.Code, rec.Body.String())
	var flow *http.Cookie
	for _, c := range rec.Result().Cookies() {
		if c.Name == flowCookie {
			flow = c
		}
	}
	require.NotNil(e.t, flow)
	return flow, rec.Header().Get("Location")
}

func (e *env) callback(flow *http.Cookie, query url.Values) *httptest.ResponseRecorder {
	req := httptest.NewRequest("GET", "/auth/oidc/callback?"+query.Encode(), nil)
	if flow != nil {
		req.AddCookie(flow)
	}
	return e.do(req)
}

// login runs the complete flow.
func (e *env) login() *httptest.ResponseRecorder {
	flow, location := e.start()
	code := e.idp.authorize(e.t, location)
	u, _ := url.Parse(location)
	return e.callback(flow, url.Values{"code": {code}, "state": {u.Query().Get("state")}})
}

func sessionCookie(rec *httptest.ResponseRecorder) *http.Cookie {
	for _, c := range rec.Result().Cookies() {
		if c.Name == "traggo" && c.Value != "" {
			return c
		}
	}
	return nil
}

func TestInfo_disabled(t *testing.T) {
	router := mux.NewRouter()
	require.NoError(t, Register(router, nil, config.OIDC{}, 4))
	rec := httptest.NewRecorder()
	router.ServeHTTP(rec, httptest.NewRequest("GET", InfoPath, nil))
	assert.Equal(t, `{"enabled":false,"passwordLogin":true}`+"\n", rec.Body.String())

	rec = httptest.NewRecorder()
	router.ServeHTTP(rec, httptest.NewRequest("GET", LoginPath, nil))
	assert.Equal(t, http.StatusNotFound, rec.Code)
}

func TestInfo_enabled(t *testing.T) {
	e := newEnv(t, func(c *config.OIDC) { c.DisablePasswordLogin = true })
	rec := e.do(httptest.NewRequest("GET", InfoPath, nil))
	assert.Equal(t, `{"enabled":true,"passwordLogin":false,"loginUrl":"/auth/oidc/login"}`+"\n", rec.Body.String())
}

func TestNew_invalidRedirect(t *testing.T) {
	_, err := New(nil, config.OIDC{Issuer: "https://i", ClientID: "a", RedirectURL: "/callback"}, 4)
	assert.Error(t, err)
}

func TestLogin_redirectsToProvider(t *testing.T) {
	e := newEnv(t, nil)
	flow, location := e.start()
	assert.True(t, flow.HttpOnly)
	assert.True(t, flow.Secure)
	assert.Equal(t, http.SameSiteLaxMode, flow.SameSite)
	assert.Equal(t, int(flowMaxAge.Seconds()), flow.MaxAge)
	e.idp.authorize(t, location)
	u, _ := url.Parse(location)
	assert.Equal(t, "https://traggo.example.com/auth/oidc/callback", u.Query().Get("redirect_uri"))
}

func TestLogin_providerDown(t *testing.T) {
	e := newEnv(t, nil)
	e.idp.srv.Close()
	rec := e.do(httptest.NewRequest("GET", LoginPath, nil))
	assert.Equal(t, http.StatusBadGateway, rec.Code)
}

func TestCallback_existingUser(t *testing.T) {
	e := newEnv(t, nil)
	e.db.Create(&model.User{Name: "alice", Pass: []byte{1}})

	rec := e.login()
	assert.Equal(t, http.StatusFound, rec.Code, rec.Body.String())
	assert.Equal(t, "/", rec.Header().Get("Location"))
	session := sessionCookie(rec)
	require.NotNil(t, session)

	device := model.Device{}
	require.False(t, e.db.Where("token = ?", session.Value).Find(&device).RecordNotFound())
	assert.Equal(t, model.TypeLongExpiry, device.Type)
	assert.Equal(t, session.MaxAge, device.Type.Seconds())
	user := model.User{}
	e.db.First(&user, device.UserID)
	assert.Equal(t, "alice", user.Name)
	var count int
	e.db.Model(new(model.User)).Count(&count)
	assert.Equal(t, 1, count)

	// the flow cookie is cleared
	for _, c := range rec.Result().Cookies() {
		if c.Name == flowCookie {
			assert.Equal(t, -1, c.MaxAge)
		}
	}
}

func TestCallback_unknownUser_noAutoRegister(t *testing.T) {
	e := newEnv(t, nil)
	rec := e.login()
	assert.Equal(t, http.StatusForbidden, rec.Code)
	assert.Nil(t, sessionCookie(rec))
	var count int
	e.db.Model(new(model.User)).Count(&count)
	assert.Equal(t, 0, count)
}

func TestCallback_autoRegister(t *testing.T) {
	e := newEnv(t, func(c *config.OIDC) {
		c.AutoRegister = true
		c.AdminClaim = "groups"
		c.AdminValue = "traggo-admins"
	})
	e.idp.claims["groups"] = []string{"users", "traggo-admins"}

	rec := e.login()
	assert.Equal(t, http.StatusFound, rec.Code, rec.Body.String())
	assert.NotNil(t, sessionCookie(rec))
	user := model.User{}
	require.False(t, e.db.Where("name = ?", "alice").Find(&user).RecordNotFound())
	assert.True(t, user.Admin)
	assert.NotEmpty(t, user.Pass)
}

func TestCallback_autoRegister_notAdmin(t *testing.T) {
	e := newEnv(t, func(c *config.OIDC) {
		c.AutoRegister = true
		c.AdminClaim = "groups"
		c.AdminValue = "traggo-admins"
	})
	e.idp.claims["groups"] = []string{"users"}
	assert.Equal(t, http.StatusFound, e.login().Code)
	user := model.User{}
	e.db.Where("name = ?", "alice").Find(&user)
	assert.False(t, user.Admin)
}

func TestCallback_adminFlagSyncedOnlyWhenConfigured(t *testing.T) {
	e := newEnv(t, nil)
	e.db.Create(&model.User{Name: "alice", Pass: []byte{1}, Admin: true})
	assert.Equal(t, http.StatusFound, e.login().Code)
	user := model.User{}
	e.db.Where("name = ?", "alice").Find(&user)
	assert.True(t, user.Admin, "admin flag must stay untouched without admin claim config")

	e = newEnv(t, func(c *config.OIDC) { c.AdminClaim = "role"; c.AdminValue = "admin" })
	e.db.Create(&model.User{Name: "alice", Pass: []byte{1}, Admin: true})
	e.idp.claims["role"] = "user"
	assert.Equal(t, http.StatusFound, e.login().Code)
	e.db.Where("name = ?", "alice").Find(&user)
	assert.False(t, user.Admin)
}

func TestCallback_customUsernameClaim(t *testing.T) {
	e := newEnv(t, func(c *config.OIDC) { c.UsernameClaim = "email" })
	e.idp.claims["email"] = "alice@example.com"
	e.db.Create(&model.User{Name: "alice@example.com", Pass: []byte{1}})
	assert.Equal(t, http.StatusFound, e.login().Code)
}

func TestCallback_missingUsernameClaim(t *testing.T) {
	e := newEnv(t, func(c *config.OIDC) { c.UsernameClaim = "nope"; c.AutoRegister = true })
	rec := e.login()
	assert.Equal(t, http.StatusForbidden, rec.Code)
	assert.Nil(t, sessionCookie(rec))
}

func TestCallback_stateMismatch(t *testing.T) {
	e := newEnv(t, nil)
	e.db.Create(&model.User{Name: "alice", Pass: []byte{1}})
	flow, location := e.start()
	code := e.idp.authorize(t, location)
	rec := e.callback(flow, url.Values{"code": {code}, "state": {"other"}})
	assert.Equal(t, http.StatusBadRequest, rec.Code)
	assert.Nil(t, sessionCookie(rec))
}

func TestCallback_noOrTamperedFlowCookie(t *testing.T) {
	e := newEnv(t, nil)
	rec := e.callback(nil, url.Values{"code": {"c"}, "state": {"s"}})
	assert.Equal(t, http.StatusBadRequest, rec.Code)

	flow, location := e.start()
	code := e.idp.authorize(t, location)
	u, _ := url.Parse(location)
	flow.Value = "x" + flow.Value
	rec = e.callback(flow, url.Values{"code": {code}, "state": {u.Query().Get("state")}})
	assert.Equal(t, http.StatusBadRequest, rec.Code)
}

func TestCallback_cookieFromOtherInstanceRejected(t *testing.T) {
	e := newEnv(t, nil)
	other := newEnv(t, nil)
	flow, location := other.start()
	code := other.idp.authorize(t, location)
	u, _ := url.Parse(location)
	rec := e.callback(flow, url.Values{"code": {code}, "state": {u.Query().Get("state")}})
	assert.Equal(t, http.StatusBadRequest, rec.Code)
}

func TestCallback_providerError(t *testing.T) {
	e := newEnv(t, nil)
	flow, location := e.start()
	u, _ := url.Parse(location)
	rec := e.callback(flow, url.Values{"error": {"access_denied"}, "state": {u.Query().Get("state")}})
	assert.Equal(t, http.StatusUnauthorized, rec.Code)
	assert.Nil(t, sessionCookie(rec))
}

func TestCallback_codeReplayFails(t *testing.T) {
	e := newEnv(t, nil)
	e.db.Create(&model.User{Name: "alice", Pass: []byte{1}})
	flow, location := e.start()
	code := e.idp.authorize(t, location)
	u, _ := url.Parse(location)
	q := url.Values{"code": {code}, "state": {u.Query().Get("state")}}
	assert.Equal(t, http.StatusFound, e.callback(flow, q).Code)
	assert.Equal(t, http.StatusUnauthorized, e.callback(flow, q).Code)
}

func TestCallback_wrongPKCEVerifierFails(t *testing.T) {
	e := newEnv(t, nil)
	e.db.Create(&model.User{Name: "alice", Pass: []byte{1}})
	flow, location := e.start()
	code := e.idp.authorize(t, location)
	// the attacker redeems a code which was issued for another challenge
	e.idp.codes[code] = authRequest{challenge: "other", nonce: e.idp.codes[code].nonce}
	u, _ := url.Parse(location)
	rec := e.callback(flow, url.Values{"code": {code}, "state": {u.Query().Get("state")}})
	assert.Equal(t, http.StatusUnauthorized, rec.Code)
	assert.Nil(t, sessionCookie(rec))
}

func TestCallback_invalidIDToken(t *testing.T) {
	cases := map[string]func(*mockIDP){
		"nonce mismatch": func(m *mockIDP) { m.nonceOverride = "evil" },
		"wrong audience": func(m *mockIDP) { m.audienceOverride = "another-client" },
		"no id token":    func(m *mockIDP) { m.skipIDToken = true },
	}
	for name, mod := range cases {
		t.Run(name, func(t *testing.T) {
			e := newEnv(t, nil)
			e.db.Create(&model.User{Name: "alice", Pass: []byte{1}})
			mod(e.idp)
			rec := e.login()
			assert.Equal(t, http.StatusUnauthorized, rec.Code)
			assert.Nil(t, sessionCookie(rec))
		})
	}
}

func TestCallback_tokenSignedWithUnknownKey(t *testing.T) {
	e := newEnv(t, nil)
	e.db.Create(&model.User{Name: "alice", Pass: []byte{1}})
	key, err := rsa.GenerateKey(rand.Reader, 2048)
	require.NoError(t, err)
	e.idp.signKey = key // the jwks endpoint still serves the other public key
	rec := e.login()
	assert.Equal(t, http.StatusUnauthorized, rec.Code)
	assert.Nil(t, sessionCookie(rec))
}
