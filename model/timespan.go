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
	UserID        int
	Tags          []TimeSpanTag
}

// TimeSpanTag is a tag for a time range
type TimeSpanTag struct {
	TimeSpanID  int
	Key         string
	StringValue *string
}
