package config

import (
	"errors"
	"io/ioutil"
	"os"
	"path/filepath"
	"testing"

	"github.com/rs/zerolog"
	"github.com/stretchr/testify/assert"
	"github.com/traggo/server/config/mode"
)

func TestMain(m *testing.M) {
	absoluteFiles = []string{}
	os.Exit(m.Run())
}

func TestGet_withoutFiles_containsDefault(t *testing.T) {
	conf, logs := Get()

	assert.Empty(t, logs)
	assert.Equal(t, 3030, conf.Port)
}

func TestGet_withoutFiles_usesEnvVariables(t *testing.T) {
	assert.Nil(t, os.Setenv("TRAGGO_PORT", "9999"))
	defer os.Unsetenv("TRAGGO_PORT")

	conf, logs := Get()

	assert.Empty(t, logs)
	assert.Equal(t, 9999, conf.Port)
}

func TestGet_withFile_usesFile(t *testing.T) {
	defer os.Unsetenv("TRAGGO_PORT")
	assert.Nil(t, ioutil.WriteFile(".env", []byte("TRAGGO_PORT=5555"), 0777))
	defer os.Remove(".env")

	conf, logs := Get()

	expected := []FutureLog{{
		Msg:   `Loading file .env`,
		Level: zerolog.DebugLevel,
	}}
	assert.Equal(t, expected, logs)
	assert.Equal(t, 5555, conf.Port)
}

func TestGet_withInvalidEnvParam_errs(t *testing.T) {
	assert.Nil(t, os.Setenv("TRAGGO_PORT", "asdasd"))
	defer os.Unsetenv("TRAGGO_PORT")

	_, logs := Get()

	expected := []FutureLog{{
		Msg:   `cannot parse env params: envconfig.Process: assigning TRAGGO_PORT to Port: converting 'asdasd' to type int. details: strconv.ParseInt: parsing "asdasd": invalid syntax`,
		Level: zerolog.FatalLevel,
	}}
	assert.Equal(t, expected, logs)
}

func TestGet_prodMode_findFileRelativeToExecutable(t *testing.T) {
	defer os.Unsetenv("TRAGGO_PORT")
	oldMode := mode.Get()
	mode.Set(mode.Prod)
	defer func() {
		mode.Set(oldMode)
	}()
	old := osExecutable
	defer func() {
		osExecutable = old
	}()
	osExecutable = func() (string, error) {
		return "./testpath/some.exe", nil
	}
	assert.Nil(t, os.MkdirAll("./testpath", 0777))
	defer os.RemoveAll("./testpath")

	assert.Nil(t, ioutil.WriteFile("./testpath/.env", []byte("TRAGGO_PORT=6666"), 0777))

	conf, logs := Get()

	expected := []FutureLog{{
		Msg:   `Loading file ` + filepath.Join("testpath", ".env"),
		Level: zerolog.DebugLevel,
	}}
	assert.Equal(t, expected, logs)
	assert.Equal(t, 6666, conf.Port)
}

func TestGet_errsOnOsExecutable(t *testing.T) {
	defer os.Unsetenv("TRAGGO_PORT")
	oldMode := mode.Get()
	mode.Set(mode.Prod)
	defer func() {
		mode.Set(oldMode)
	}()
	old := osExecutable
	defer func() {
		osExecutable = old
	}()
	osExecutable = func() (string, error) {
		return "", errors.New("oops")
	}

	_, logs := Get()

	expected := []FutureLog{{
		Msg:   `Could not get path of executable using working directory instead. oops`,
		Level: zerolog.ErrorLevel,
	}}
	assert.Equal(t, expected, logs)
}

func TestGet_withoutFilePermission(t *testing.T) {
	defer os.Unsetenv("TRAGGO_PORT")
	old := osStat
	defer func() {
		osStat = old
	}()
	osStat = func(name string) (os.FileInfo, error) {
		if name == ".env" {
			return nil, os.ErrPermission
		}
		return os.Stat(name)
	}

	_, logs := Get()

	expected := []FutureLog{{
		Msg:   `cannot read file .env because permission denied`,
		Level: zerolog.WarnLevel,
	}}
	assert.Equal(t, expected, logs)
}

func TestGet_checksAbsoluteFiles(t *testing.T) {
	defer os.Unsetenv("TRAGGO_PORT")
	assert.Nil(t, os.MkdirAll("./absolutepath", 0777))
	defer os.RemoveAll("./absolutepath")
	path, err := filepath.Abs("./absolutepath/.env")
	assert.Nil(t, err)
	absoluteFiles = []string{path}
	defer func() {
		absoluteFiles = []string{}
	}()
	assert.Nil(t, ioutil.WriteFile("./absolutepath/.env", []byte("TRAGGO_PORT=7777"), 0777))

	conf, logs := Get()

	expected := []FutureLog{{
		Msg:   `Loading file ` + path,
		Level: zerolog.DebugLevel,
	}}
	assert.Equal(t, expected, logs)
	assert.Equal(t, 7777, conf.Port)
}
