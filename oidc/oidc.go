// Package oidc implements the optional OpenID Connect login (authorization code flow with PKCE).
// After a successful login the regular traggo session cookie is issued.
package oidc

import (
	"context"
	"crypto/hmac"
	"crypto/rand"
	"crypto/sha256"
	"crypto/subtle"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"strings"
	"sync"
	"time"

	gooidc "github.com/coreos/go-oidc/v3/oidc"
	"github.com/gorilla/mux"
	"github.com/jinzhu/gorm"
	"github.com/rs/zerolog/log"
	"github.com/traggo/server/auth"
	"github.com/traggo/server/config"
	"github.com/traggo/server/device"
	"github.com/traggo/server/model"
	"github.com/traggo/server/user/password"
	"golang.org/x/oauth2"
)

const (
	// LoginPath starts the login by redirecting to the identity provider.
	LoginPath = "/auth/oidc/login"
	// InfoPath tells the ui whether OIDC is enabled.
	InfoPath = "/auth/oidc/info"

	flowCookie  = "traggo_oidc"
	flowMaxAge  = 10 * time.Minute
	deviceName  = "OIDC login"
	deviceType  = model.TypeLongExpiry
	httpTimeout = 10 * time.Second
)

// Info is returned by InfoPath.
type Info struct {
	Enabled       bool   `json:"enabled"`
	PasswordLogin bool   `json:"passwordLogin"`
	LoginURL      string `json:"loginUrl,omitempty"`
}

// Handler serves the OIDC endpoints.
type Handler struct {
	conf         config.OIDC
	db           *gorm.DB
	passStrength int
	redirect     *url.URL
	key          []byte

	mu       sync.Mutex
	provider *gooidc.Provider
}

// Register adds the OIDC routes to the router. The info route is always available, the login routes only when enabled.
// The auth middleware must be installed on the router.
func Register(r *mux.Router, db *gorm.DB, conf config.OIDC, passStrength int) error {
	h, err := New(db, conf, passStrength)
	if err != nil {
		return err
	}
	r.HandleFunc(InfoPath, h.Info).Methods(http.MethodGet)
	if conf.Enabled() {
		r.HandleFunc(LoginPath, h.Login).Methods(http.MethodGet)
		r.HandleFunc(h.redirect.Path, h.Callback).Methods(http.MethodGet)
	}
	return nil
}

// New creates a handler.
func New(db *gorm.DB, conf config.OIDC, passStrength int) (*Handler, error) {
	h := &Handler{conf: conf, db: db, passStrength: passStrength, key: make([]byte, 32)}
	if !conf.Enabled() {
		return h, nil
	}
	if err := conf.Validate(); err != nil {
		return nil, err
	}
	redirect, err := url.Parse(conf.RedirectURL)
	if err != nil || redirect.Host == "" || (redirect.Scheme != "https" && redirect.Scheme != "http") {
		return nil, errors.New("TRAGGO_OIDC_REDIRECT_URL must be an absolute http(s) url")
	}
	h.redirect = redirect
	if _, err := rand.Read(h.key); err != nil {
		return nil, err
	}
	return h, nil
}

// Info reports whether OIDC is enabled.
func (h *Handler) Info(w http.ResponseWriter, _ *http.Request) {
	info := Info{PasswordLogin: !h.conf.DisablePasswordLogin}
	if h.conf.Enabled() {
		info.Enabled = true
		info.LoginURL = LoginPath
	}
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Cache-Control", "no-store")
	_ = json.NewEncoder(w).Encode(info)
}

// Login redirects to the identity provider.
func (h *Handler) Login(w http.ResponseWriter, r *http.Request) {
	oauthConf, _, err := h.oauthConfig(r.Context())
	if err != nil {
		log.Error().Err(err).Msg("OIDC: provider not available")
		http.Error(w, "identity provider not available", http.StatusBadGateway)
		return
	}

	state, nonce, verifier := randomString(), randomString(), oauth2.GenerateVerifier()
	h.setFlowCookie(w, state+"."+nonce+"."+verifier)

	http.Redirect(w, r, oauthConf.AuthCodeURL(state,
		oauth2.S256ChallengeOption(verifier),
		gooidc.Nonce(nonce)), http.StatusFound)
}

// Callback finishes the login.
func (h *Handler) Callback(w http.ResponseWriter, r *http.Request) {
	state, nonce, verifier, ok := h.readFlowCookie(r)
	h.clearFlowCookie(w)
	if !ok {
		http.Error(w, "login session expired or invalid, try again", http.StatusBadRequest)
		return
	}
	query := r.URL.Query()
	if subtle.ConstantTimeCompare([]byte(query.Get("state")), []byte(state)) != 1 {
		http.Error(w, "invalid state", http.StatusBadRequest)
		return
	}
	if e := query.Get("error"); e != "" {
		log.Info().Str("error", e).Str("description", query.Get("error_description")).Msg("OIDC: identity provider returned an error")
		http.Error(w, "login failed", http.StatusUnauthorized)
		return
	}
	code := query.Get("code")
	if code == "" {
		http.Error(w, "missing code", http.StatusBadRequest)
		return
	}

	ctx := r.Context()
	oauthConf, verifierIDT, err := h.oauthConfig(ctx)
	if err != nil {
		log.Error().Err(err).Msg("OIDC: provider not available")
		http.Error(w, "identity provider not available", http.StatusBadGateway)
		return
	}
	ctx = gooidc.ClientContext(ctx, &http.Client{Timeout: httpTimeout})
	token, err := oauthConf.Exchange(ctx, code, oauth2.VerifierOption(verifier))
	if err != nil {
		log.Info().Err(err).Msg("OIDC: code exchange failed")
		http.Error(w, "login failed", http.StatusUnauthorized)
		return
	}
	rawIDToken, _ := token.Extra("id_token").(string)
	if rawIDToken == "" {
		log.Info().Msg("OIDC: response contains no id_token")
		http.Error(w, "login failed", http.StatusUnauthorized)
		return
	}
	idToken, err := verifierIDT.Verify(ctx, rawIDToken)
	if err != nil {
		log.Info().Err(err).Msg("OIDC: id token verification failed")
		http.Error(w, "login failed", http.StatusUnauthorized)
		return
	}
	if subtle.ConstantTimeCompare([]byte(idToken.Nonce), []byte(nonce)) != 1 {
		log.Info().Msg("OIDC: nonce mismatch")
		http.Error(w, "login failed", http.StatusUnauthorized)
		return
	}
	claims := map[string]interface{}{}
	if err := idToken.Claims(&claims); err != nil {
		log.Info().Err(err).Msg("OIDC: unable to read claims")
		http.Error(w, "login failed", http.StatusUnauthorized)
		return
	}

	user, err := h.findOrCreateUser(claims)
	if err != nil {
		log.Info().Err(err).Msg("OIDC: user mapping failed")
		http.Error(w, err.Error(), http.StatusForbidden)
		return
	}

	dev, err := device.Create(h.db, user, deviceName, deviceType)
	if err != nil {
		log.Error().Err(err).Msg("OIDC: unable to create device")
		http.Error(w, "internal error", http.StatusInternalServerError)
		return
	}
	auth.GetCreateSession(r.Context())(dev.Token, dev.Type.Seconds())
	log.Info().Str("username", user.Name).Msg("OIDC: login successful")
	http.Redirect(w, r, "/", http.StatusFound)
}

var (
	errNoUsername = errors.New("identity provider did not return a username")
	errNoUser     = errors.New("no matching user, ask an administrator to create one")
)

func (h *Handler) findOrCreateUser(claims map[string]interface{}) (*model.User, error) {
	name, _ := claims[h.conf.UsernameClaim].(string)
	name = strings.TrimSpace(name)
	if name == "" {
		return nil, errNoUsername
	}

	admin := h.isAdmin(claims)
	user := &model.User{}
	if h.db.Where("name = ?", name).Find(user).RecordNotFound() {
		if !h.conf.AutoRegister {
			return nil, errNoUser
		}
		// the random password can not be guessed, the user can only log in via OIDC unless an admin sets a password.
		user = &model.User{
			Name:  name,
			Pass:  password.CreatePassword(randomString(), h.passStrength),
			Admin: admin,
		}
		if err := h.db.Create(user).Error; err != nil {
			return nil, fmt.Errorf("unable to create user: %w", err)
		}
		log.Info().Str("username", name).Bool("admin", admin).Msg("OIDC: created user")
		return user, nil
	}

	if h.conf.AdminClaim != "" && user.Admin != admin {
		user.Admin = admin
		if err := h.db.Model(user).Update("admin", admin).Error; err != nil {
			return nil, fmt.Errorf("unable to update user: %w", err)
		}
		log.Info().Str("username", name).Bool("admin", admin).Msg("OIDC: updated admin flag")
	}
	return user, nil
}

// isAdmin is true when the admin claim equals the admin value or is a list containing it.
func (h *Handler) isAdmin(claims map[string]interface{}) bool {
	if h.conf.AdminClaim == "" {
		return false
	}
	switch v := claims[h.conf.AdminClaim].(type) {
	case string:
		return v == h.conf.AdminValue
	case []interface{}:
		for _, item := range v {
			if s, ok := item.(string); ok && s == h.conf.AdminValue {
				return true
			}
		}
	}
	return false
}

// oauthConfig lazily discovers the provider, so traggo starts even if the identity provider is down.
func (h *Handler) oauthConfig(ctx context.Context) (*oauth2.Config, *gooidc.IDTokenVerifier, error) {
	h.mu.Lock()
	defer h.mu.Unlock()
	if h.provider == nil {
		ctx = gooidc.ClientContext(ctx, &http.Client{Timeout: httpTimeout})
		provider, err := gooidc.NewProvider(ctx, h.conf.Issuer)
		if err != nil {
			return nil, nil, err
		}
		h.provider = provider
	}
	return &oauth2.Config{
		ClientID:     h.conf.ClientID,
		ClientSecret: string(h.conf.ClientSecret),
		Endpoint:     h.provider.Endpoint(),
		RedirectURL:  h.conf.RedirectURL,
		Scopes:       h.conf.Scopes,
	}, h.provider.Verifier(&gooidc.Config{ClientID: h.conf.ClientID}), nil
}

func (h *Handler) sign(value string) string {
	mac := hmac.New(sha256.New, h.key)
	mac.Write([]byte(value))
	return base64.RawURLEncoding.EncodeToString(mac.Sum(nil))
}

func (h *Handler) flowCookie(value string, maxAge int) *http.Cookie {
	return &http.Cookie{
		Name:     flowCookie,
		Value:    value,
		Path:     "/",
		MaxAge:   maxAge,
		HttpOnly: true,
		Secure:   h.redirect.Scheme == "https",
		SameSite: http.SameSiteLaxMode, // lax: the cookie must be sent on the redirect back from the identity provider
	}
}

func (h *Handler) setFlowCookie(w http.ResponseWriter, value string) {
	payload := base64.RawURLEncoding.EncodeToString([]byte(value))
	http.SetCookie(w, h.flowCookie(payload+"."+h.sign(payload), int(flowMaxAge.Seconds())))
}

func (h *Handler) clearFlowCookie(w http.ResponseWriter) {
	http.SetCookie(w, h.flowCookie("", -1))
}

func (h *Handler) readFlowCookie(r *http.Request) (state, nonce, verifier string, ok bool) {
	c, err := r.Cookie(flowCookie)
	if err != nil {
		return "", "", "", false
	}
	parts := strings.Split(c.Value, ".")
	if len(parts) != 2 || !hmac.Equal([]byte(h.sign(parts[0])), []byte(parts[1])) {
		return "", "", "", false
	}
	raw, err := base64.RawURLEncoding.DecodeString(parts[0])
	if err != nil {
		return "", "", "", false
	}
	values := strings.Split(string(raw), ".")
	if len(values) != 3 {
		return "", "", "", false
	}
	return values[0], values[1], values[2], true
}

func randomString() string {
	b := make([]byte, 32)
	if _, err := rand.Read(b); err != nil {
		panic(err)
	}
	return base64.RawURLEncoding.EncodeToString(b)
}
