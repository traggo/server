package model

import "time"

// UserSetting a setting for a user.
type UserSetting struct {
	UserID            int `gorm:"primary_key;unique_index"`
	Theme             string
	DateLocale        string
	FirstDayOfTheWeek string
}

// Settings constants
const (
	ThemeGruvboxDark   = "GruvboxDark"
	ThemeGruvboxLight  = "GruvboxLight"
	ThemeMaterialDark  = "MaterialDark"
	ThemeMaterialLight = "MaterialLight"

	DateLocaleGerman      = "German"
	DateLocaleAmerican    = "American"
	DateLocaleAmerican24h = "American24h"
	DateLocaleBritish     = "British"
	DateLocaleAustralian  = "Australian"
)

var daysOfWeek = map[string]time.Weekday{
	"Sunday":    time.Sunday,
	"Monday":    time.Monday,
	"Tuesday":   time.Tuesday,
	"Wednesday": time.Wednesday,
	"Thursday":  time.Thursday,
	"Friday":    time.Friday,
	"Saturday":  time.Saturday,
}

// FirstDayOfTheWeekTimeWeekday returns the configured first day of the week.
func (s UserSetting) FirstDayOfTheWeekTimeWeekday() time.Weekday {
	return daysOfWeek[s.FirstDayOfTheWeek]
}

// LastDayOfTheWeekTimeWeekday returns the configured last day of the week.
func (s UserSetting) LastDayOfTheWeekTimeWeekday() time.Weekday {
	return daysOfWeek[(daysOfWeek[s.FirstDayOfTheWeek] - 1).String()]
}
