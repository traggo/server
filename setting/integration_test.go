package setting

import (
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/traggo/server/generated/gqlmodel"
	"github.com/traggo/server/test"
	"github.com/traggo/server/test/fake"
)

func TestSettings(t *testing.T) {
	db := test.InMemoryDB(t)
	defer db.Close()

	resolver := &ResolverForSettings{DB: db.DB}

	settings, err := resolver.Settings(fake.User(1), "ui")
	require.NoError(t, err)
	require.Empty(t, settings)

	settings, err = resolver.Settings(fake.User(1), "mobile")
	require.NoError(t, err)
	require.Empty(t, settings)

	value, err := resolver.SettingGet(fake.User(1), "ui", "test")
	require.NoError(t, err)
	require.Equal(t, "", value)

	_, err = resolver.SettingPut(fake.User(1), "ui", "test", "myvalue")
	require.NoError(t, err)

	value, err = resolver.SettingGet(fake.User(1), "ui", "test")
	require.NoError(t, err)
	require.Equal(t, "myvalue", value)

	value, err = resolver.SettingGet(fake.User(2), "ui", "test")
	require.NoError(t, err)
	require.Equal(t, "", value)

	_, err = resolver.SettingPut(fake.User(2), "ui", "test", "myvalue2")
	require.NoError(t, err)

	value, err = resolver.SettingGet(fake.User(2), "ui", "test")
	require.NoError(t, err)
	require.Equal(t, "myvalue2", value)

	value, err = resolver.SettingGet(fake.User(1), "ui", "test")
	require.NoError(t, err)
	require.Equal(t, "myvalue", value)

	values, err := resolver.Settings(fake.User(1), "ui")
	require.NoError(t, err)
	require.Equal(t, []*gqlmodel.Setting{{Key: "test", Value: "myvalue"}}, values)

	_, err = resolver.SettingPut(fake.User(2), "ui", "test2", "myvalue2")
	require.NoError(t, err)

	_, err = resolver.SettingPut(fake.User(1), "ui", "test2", "myvalue")
	require.NoError(t, err)

	values, err = resolver.Settings(fake.User(1), "ui")
	require.NoError(t, err)
	require.Equal(t, []*gqlmodel.Setting{{Key: "test", Value: "myvalue"}, {Key: "test2", Value: "myvalue"}}, values)

	values, err = resolver.Settings(fake.User(2), "ui")
	require.NoError(t, err)
	require.Equal(t, []*gqlmodel.Setting{{Key: "test", Value: "myvalue2"}, {Key: "test2", Value: "myvalue2"}}, values)

	values, err = resolver.Settings(fake.User(1), "mobile")
	require.NoError(t, err)
	require.Empty(t, values)

	value, err = resolver.SettingGet(fake.User(1), "mobile", "test")
	require.NoError(t, err)
	require.Equal(t, "", value)

	_, err = resolver.SettingPut(fake.User(1), "mobile", "test", "mobilevalue")
	require.NoError(t, err)

	value, err = resolver.SettingGet(fake.User(1), "mobile", "test")
	require.NoError(t, err)
	require.Equal(t, "mobilevalue", value)

	value, err = resolver.SettingGet(fake.User(1), "ui", "test")
	require.NoError(t, err)
	require.Equal(t, "myvalue", value)
}
