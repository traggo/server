package time

import (
	"fmt"
	"time"

	"github.com/jmattheis/go-timemath"
	"github.com/traggo/server/model"
)

const nowKey = "now"

// RelativeRange represents a relative range.
type RelativeRange struct {
	From string
	To   string
}

// StaticRange represents a concrete range.
type StaticRange struct {
	From time.Time
	To   time.Time
}

// ParseRange parses a range and converts it to static ranges.
func ParseRange(now time.Time, r RelativeRange, interval model.Interval) ([]StaticRange, error) {
	from, err := ParseTime(now, r.From, true, time.Monday)
	if err != nil {
		return nil, fmt.Errorf("range from: %s", err)
	}
	to, err := ParseTime(now, r.To, false, time.Sunday)
	if err != nil {
		return nil, fmt.Errorf("range to: %s", err)
	}

	return ranges(from, to, interval), nil
}

// ParseTime parses time.
func ParseTime(now time.Time, value string, startOf bool, weekday time.Weekday) (time.Time, error) {
	parse, err := time.Parse(time.RFC3339, value)
	if err == nil {
		return parse, nil
	}

	return timemath.Parse(now, value, startOf, weekday)
}
