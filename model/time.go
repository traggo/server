package model

import (
	"fmt"
	"io"
	"time"
)

// Time scalar type for graphql
type Time time.Time

// Time returns the wrapped time
func (t Time) Time() time.Time {
	return time.Time(t)
}

// OmitTimeZone omits the time zone and removes a utc date.
func (t Time) OmitTimeZone() time.Time {
	x := t.Time()
	return time.Date(x.Year(), x.Month(), x.Day(), x.Hour(), x.Minute(), x.Second(), x.Nanosecond(), time.UTC)
}

// UTC changes the timezone to utc.
func (t Time) UTC() time.Time {
	return t.Time().UTC()
}

// MarshalGQL implements the graphql.Marshaler interface
func (t Time) MarshalGQL(w io.Writer) {
	w.Write([]byte(t.Time().Format(time.RFC3339)))
}

// UnmarshalGQL implements the graphql.Unmarshaler interface
func (t *Time) UnmarshalGQL(v interface{}) error {
	raw, ok := v.(string)
	if !ok {
		return fmt.Errorf("time must be a string")
	}

	parse, err := time.Parse(time.RFC3339, raw)
	if err != nil {
		return err
	}
	*t = Time(parse)
	return nil
}
