package auth

import (
	"time"

	"github.com/jinzhu/gorm"
	"github.com/jmattheis/go-timemath"
	"github.com/rs/zerolog/log"
	"github.com/traggo/server/model"
)

var timeNow = time.Now

// CleanUp clean up expired devices
func CleanUp(db *gorm.DB, interval time.Duration, close chan bool) {
	for {
		select {
		case <-time.After(interval):
			affected := db.Where("type = ? AND active_at < ?", model.TypeLongExpiry, timemath.Second.Subtract(timeNow(), model.TypeLongExpiry.Seconds())).
				Or("type = ? AND active_at < ?", model.TypeShortExpiry, timemath.Second.Subtract(timeNow(), model.TypeShortExpiry.Seconds())).
				Delete(new(model.Device)).RowsAffected
			if affected > 0 {
				log.Debug().Int64("amount", affected).Msg("removed devices")
			}
		case <-close:
			return
		}
	}
}
