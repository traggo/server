package setting

import (
	"testing"
	"time"

	"github.com/stretchr/testify/require"
	"github.com/traggo/server/generated/gqlmodel"
	"github.com/traggo/server/test"
	"github.com/traggo/server/test/fake"
)

const (
	invalidInternalWeekday = time.Weekday(100)
	invalidExternalWeekday = gqlmodel.WeekDay("abc")
)

func TestSettingsResolver(t *testing.T) {
	db := test.InMemoryDB(t)
	defer db.Close()

	db.User(1)

	resolver := &ResolverForSettings{DB: db.DB}

	settings, err := resolver.UserSettings(fake.User(1))
	require.NoError(t, err)
	require.Equal(t, gqlmodel.ThemeGruvboxDark, settings.Theme)

	_, err = resolver.SetUserSettings(fake.User(1), gqlmodel.InputUserSettings{
		Theme:              gqlmodel.ThemeGruvboxLight,
		DateLocale:         gqlmodel.DateLocaleGerman,
		FirstDayOfTheWeek:  gqlmodel.WeekDayWednesday,
		DateTimeInputStyle: gqlmodel.DateTimeInputStyleFancy,
	})
	require.NoError(t, err)

	settings, err = resolver.UserSettings(fake.User(1))
	require.NoError(t, err)
	require.Equal(t, &gqlmodel.UserSettings{
		Theme:              gqlmodel.ThemeGruvboxLight,
		DateLocale:         gqlmodel.DateLocaleGerman,
		FirstDayOfTheWeek:  gqlmodel.WeekDayWednesday,
		DateTimeInputStyle: gqlmodel.DateTimeInputStyleFancy,
	}, settings)
}

func TestWeekDayConvert(t *testing.T) {
	for _, day := range gqlmodel.AllWeekDay {
		require.Equal(t, day, toExternalWeekday(toInternalWeekday(day)))
	}
	require.Panics(t, func() {
		toInternalWeekday(invalidExternalWeekday)
	})
	require.Panics(t, func() {
		toExternalWeekday(invalidInternalWeekday)
	})
}

func TestShouldHandleInvalidInputs(t *testing.T) {
	toExternalTheme("aoeuaoeu")
	toInternalTheme("aoeuaoeu")
	toExternalDateLocale("aeu")
	toInternalDateLocale("aoeu")
	toExternalDateTimeInputStyle("aoeu")
	toInternalDateTimeInputStyle("aoeu")

}
