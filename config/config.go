package config

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/joho/godotenv"
	"github.com/kelseyhightower/envconfig"
	"github.com/rs/zerolog"
	"github.com/traggo/server/config/mode"
)

var (
	prefix        = "traggo"
	files         = []string{".env.development.local", ".env.development", ".env.local", ".env"}
	absoluteFiles = []string{"/etc/traggo/server.ini"}
	osExecutable  = os.Executable
	osStat        = os.Stat
)

// Config represents the application configuration.
type Config struct {
	Port               int      `default:"3030"`
	LogLevel           LogLevel `default:"info" split_words:"true"`
	DefaultUserName    string   `default:"admin" split_words:"true"`
	DefaultUserPass    string   `default:"admin" split_words:"true"`
	PassStrength       int      `default:"10" split_words:"true"`
	DatabaseDialect    string   `default:"sqlite3" split_words:"true"`
	DatabaseConnection string   `default:"data/traggo.db" split_words:"true"`

	// OIDC is the optional OpenID Connect login, it is disabled unless TRAGGO_OIDC_ISSUER is set.
	OIDC OIDC `envconfig:"oidc"`
}

// Secret is a string which is masked when marshalled or printed, so it does not end up in logs.
type Secret string

// MarshalJSON masks the value.
func (Secret) MarshalJSON() ([]byte, error) { return []byte(`"***"`), nil }

// String masks the value.
func (Secret) String() string { return "***" }

// OIDC holds the OpenID Connect settings.
type OIDC struct {
	// Issuer is the issuer URL, the discovery document is fetched from <issuer>/.well-known/openid-configuration.
	Issuer       string   `split_words:"true"`
	ClientID     string   `envconfig:"client_id"`
	ClientSecret Secret   `envconfig:"client_secret"`
	RedirectURL  string   `envconfig:"redirect_url"`
	Scopes       []string `default:"openid,profile,email"`
	// UsernameClaim is the ID token claim which is matched with the traggo user name.
	UsernameClaim string `envconfig:"username_claim" default:"preferred_username"`
	// AdminClaim and AdminValue make users admins if the claim equals / contains the value, both empty: admin flag is not managed.
	AdminClaim string `envconfig:"admin_claim"`
	AdminValue string `envconfig:"admin_value"`
	// AutoRegister creates users which do not exist yet.
	AutoRegister bool `envconfig:"auto_register"`
	// DisablePasswordLogin rejects the username/password login.
	DisablePasswordLogin bool `envconfig:"disable_password_login"`
}

// Enabled returns true when OIDC is configured.
func (o OIDC) Enabled() bool {
	return o.Issuer != ""
}

// Validate checks that all required values are set when OIDC is enabled.
func (o OIDC) Validate() error {
	if !o.Enabled() {
		if o.DisablePasswordLogin {
			return fmt.Errorf("TRAGGO_OIDC_DISABLE_PASSWORD_LOGIN requires TRAGGO_OIDC_ISSUER")
		}
		return nil
	}
	if o.ClientID == "" || o.RedirectURL == "" {
		return fmt.Errorf("TRAGGO_OIDC_CLIENT_ID and TRAGGO_OIDC_REDIRECT_URL are required when TRAGGO_OIDC_ISSUER is set")
	}
	if (o.AdminClaim == "") != (o.AdminValue == "") {
		return fmt.Errorf("TRAGGO_OIDC_ADMIN_CLAIM and TRAGGO_OIDC_ADMIN_VALUE must be set together")
	}
	return nil
}

// Get loads the application config.
func Get() (Config, []FutureLog) {
	var logs []FutureLog
	dir, log := getExecutableOrWorkDir()
	if log != nil {
		logs = append(logs, *log)
	}

	for _, file := range getFiles(dir) {
		_, fileErr := osStat(file)
		if fileErr == nil {
			if err := godotenv.Load(file); err != nil {
				logs = append(logs, FutureLog{
					Level: zerolog.FatalLevel,
					Msg:   fmt.Sprintf("cannot load file %s: %s", file, err)})
			} else {
				logs = append(logs, FutureLog{
					Level: zerolog.DebugLevel,
					Msg:   fmt.Sprintf("Loading file %s", file)})
			}
		} else if os.IsNotExist(fileErr) {
			continue
		} else {
			logs = append(logs, FutureLog{
				Level: zerolog.WarnLevel,
				Msg:   fmt.Sprintf("cannot read file %s because %s", file, fileErr)})
		}
	}

	config := Config{}
	err := envconfig.Process(prefix, &config)
	if err != nil {
		logs = append(logs, FutureLog{
			Level: zerolog.FatalLevel,
			Msg:   fmt.Sprintf("cannot parse env params: %s", err)})
	}

	if err := config.OIDC.Validate(); err != nil {
		logs = append(logs, FutureLog{
			Level: zerolog.FatalLevel,
			Msg:   fmt.Sprintf("invalid oidc config: %s", err)})
	}

	return config, logs
}

func getExecutableOrWorkDir() (string, *FutureLog) {
	dir, err := getExecutableDir()
	// when using `go run main.go` the executable lives in th temp directory therefore the env.development
	// will not be read, this enforces that the current work directory is used in dev mode.
	if err != nil || mode.Get() == mode.Dev {
		return filepath.Dir("."), err
	}
	return dir, nil
}

func getExecutableDir() (string, *FutureLog) {
	ex, err := osExecutable()
	if err != nil {
		return "", &FutureLog{
			Level: zerolog.ErrorLevel,
			Msg:   "Could not get path of executable using working directory instead. " + err.Error()}
	}
	return filepath.Dir(ex), nil
}

func getFiles(relativeTo string) []string {
	var result []string
	for _, file := range files {
		result = append(result, filepath.Join(relativeTo, file))
	}
	homeDir, err := os.UserHomeDir()
	if err == nil {
		result = append(result, filepath.Join(homeDir, ".config/traggo/server.ini"))
	}
	result = append(result, absoluteFiles...)
	return result
}
