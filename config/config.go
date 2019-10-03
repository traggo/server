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
	absoluteFiles = []string{"/etc/traggo/.env"}
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
	result = append(result, absoluteFiles...)
	return result
}
