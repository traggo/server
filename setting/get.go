package setting

import (
	"context"
	"time"

	"github.com/jinzhu/gorm"
	"github.com/traggo/server/auth"
	"github.com/traggo/server/model"
)

// Get returns the settings
func Get(ctx context.Context, db *gorm.DB) (model.UserSetting, error) {
	internal := model.UserSetting{}
	user := auth.GetUser(ctx)
	defaultSettings := model.UserSetting{
		Theme:              model.ThemeGruvboxDark,
		DateLocale:         model.DateLocaleAmerican,
		FirstDayOfTheWeek:  time.Monday.String(),
		DateTimeInputStyle: model.DateTimeInputFancy,
	}

	if user == nil {
		return defaultSettings, nil
	}
	find := db.Where(&model.UserSetting{UserID: user.ID}).Find(&internal)

	if find.RecordNotFound() {
		return defaultSettings, nil
	}

	if find.Error != nil {
		return defaultSettings, find.Error
	}

	return internal, nil
}
