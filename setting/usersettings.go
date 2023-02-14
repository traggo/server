package setting

import (
	"context"
	"time"

	"github.com/traggo/server/auth"
	"github.com/traggo/server/generated/gqlmodel"
	"github.com/traggo/server/model"
)

// SetUserSettings sets the user settings.
func (r *ResolverForSettings) SetUserSettings(ctx context.Context, settings gqlmodel.InputUserSettings) (*gqlmodel.UserSettings, error) {
	internal := model.UserSetting{
		Theme:             toInternalTheme(settings.Theme),
		FirstDayOfTheWeek: toInternalWeekday(settings.FirstDayOfTheWeek).String(),
		UserID:            auth.GetUser(ctx).ID,
		DateLocale:        toInternalDateLocale(settings.DateLocale),
	}

	save := r.DB.Save(internal)

	return toExternal(internal), save.Error
}

// UserSettings returns the user settings.
func (r *ResolverForSettings) UserSettings(ctx context.Context) (*gqlmodel.UserSettings, error) {
	settings, err := Get(ctx, r.DB)
	return toExternal(settings), err
}

func toExternal(internal model.UserSetting) *gqlmodel.UserSettings {
	return &gqlmodel.UserSettings{
		Theme:             toExternalTheme(internal.Theme),
		DateLocale:        toExternalDateLocale(internal.DateLocale),
		FirstDayOfTheWeek: toExternalWeekday(internal.FirstDayOfTheWeekTimeWeekday()),
	}
}

func toInternalDateLocale(locale gqlmodel.DateLocale) string {
	switch locale.String() {
	case model.DateLocaleAmerican, model.DateLocaleGerman, model.DateLocaleAmerican24h, model.DateLocaleAustralian, model.DateLocaleBritish:
		return locale.String()
	default:
		return model.DateLocaleAmerican
	}
}

func toExternalDateLocale(dateLocale string) gqlmodel.DateLocale {
	if gqlmodel.DateLocale(dateLocale).IsValid() {
		return gqlmodel.DateLocale(dateLocale)
	}
	if dateLocale == "English24h" {
		return gqlmodel.DateLocaleAmerican24h
	}
	return gqlmodel.DateLocaleAmerican
}

func toInternalTheme(theme gqlmodel.Theme) string {
	switch theme.String() {
	case model.ThemeGruvboxDark, model.ThemeGruvboxLight, model.ThemeMaterialLight, model.ThemeMaterialDark:
		return theme.String()
	default:
		return model.ThemeGruvboxDark
	}
}

func toExternalTheme(theme string) gqlmodel.Theme {
	if gqlmodel.Theme(theme).IsValid() {
		return gqlmodel.Theme(theme)
	}
	return gqlmodel.ThemeGruvboxDark
}

func toExternalWeekday(weekday time.Weekday) gqlmodel.WeekDay {
	switch weekday {
	case time.Monday:
		return gqlmodel.WeekDayMonday
	case time.Tuesday:
		return gqlmodel.WeekDayTuesday
	case time.Wednesday:
		return gqlmodel.WeekDayWednesday
	case time.Thursday:
		return gqlmodel.WeekDayThursday
	case time.Friday:
		return gqlmodel.WeekDayFriday
	case time.Saturday:
		return gqlmodel.WeekDaySaturday
	case time.Sunday:
		return gqlmodel.WeekDaySunday
	default:
		panic("unknown weekday")
	}
}

func toInternalWeekday(weekday gqlmodel.WeekDay) time.Weekday {
	switch weekday {
	case gqlmodel.WeekDayMonday:
		return time.Monday
	case gqlmodel.WeekDayTuesday:
		return time.Tuesday
	case gqlmodel.WeekDayWednesday:
		return time.Wednesday
	case gqlmodel.WeekDayThursday:
		return time.Thursday
	case gqlmodel.WeekDayFriday:
		return time.Friday
	case gqlmodel.WeekDaySaturday:
		return time.Saturday
	case gqlmodel.WeekDaySunday:
		return time.Sunday
	default:
		panic("unknown weekday")
	}
}
