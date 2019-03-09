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

// TimeP parses a time panics if not valid
func TimeP(value string) *time.Time {
	t := Time(value)
	return &t
}

// ModelTime parses a model.Time panics if not valid
func ModelTime(value string) model.Time {
	return model.Time(Time(value))
}

// ModelTimeP parses a model.Time panics if not valid
func ModelTimeP(value string) *model.Time {
	modelTime := ModelTime(value)
	return &modelTime
}
