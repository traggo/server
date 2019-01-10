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
