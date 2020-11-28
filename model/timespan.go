package model

import "time"

// TimeSpan is basically a tagged time range.
type TimeSpan struct {
	ID            int `gorm:"primary_key;unique_index;AUTO_INCREMENT"`
	StartUTC      time.Time
	EndUTC        *time.Time
	StartUserTime time.Time
	EndUserTime   *time.Time
	OffsetUTC     int
	UserID        int `gorm:"type:int REFERENCES users(id) ON DELETE CASCADE"`
	Tags          []TimeSpanTag
	Note          string
}

// TimeSpanTag is a tag for a time range
type TimeSpanTag struct {
	TimeSpanID  int `gorm:"type:int REFERENCES time_spans(id) ON DELETE CASCADE"`
	Key         string
	StringValue string
}
