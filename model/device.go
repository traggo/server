package model

import (
	"time"
)

// Device represents something which can connect to traggo
type Device struct {
	ID        int    `gorm:"primary_key;unique_index;AUTO_INCREMENT"`
	Token     string `gorm:"unique"`
	Name      string
	UserID    int
	CreatedAt time.Time
	ExpiresAt time.Time
	ActiveAt  time.Time
}
