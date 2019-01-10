package auth

import (
	"time"

	"github.com/jinzhu/gorm"
	"github.com/rs/zerolog/log"
	"github.com/traggo/server/model"
)

var timeNow = time.Now

// CleanUp clean up expired devices
func CleanUp(db *gorm.DB, interval time.Duration, close chan bool) {
	for {
		select {
		case <-time.After(interval):
			affected := db.Where("expires_at < ?", timeNow()).Delete(new(model.Device)).RowsAffected
			if affected > 0 {
				log.Debug().Int64("amount", affected).Msg("removed devices")
			}
		case <-close:
			return
		}
	}
}
