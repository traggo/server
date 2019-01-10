package test

import (
	"time"

	"github.com/traggo/server/model"
)

// Time parses a time panics if not valid
func Time(value string) time.Time {
	parse, err := time.ParseInLocation(time.RFC3339, value, time.UTC)
	if err != nil {
		panic(err)
	}
	return parse
}

// ModelTime parses a model.Time panics if not valid
func ModelTime(value string) model.Time {
	return model.Time(Time(value))
}
