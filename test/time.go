package test

import (
	"time"

	"github.com/traggo/server/model"
)

// Time parses a time panics if not valid
func Time(value string) time.Time {
	parse, firstErr := time.ParseInLocation(time.RFC3339, value, time.UTC)
	if firstErr != nil {
		var err error
		parse, err = time.ParseInLocation(time.RFC3339Nano, value, time.UTC)
		if err != nil {
			panic(firstErr)
		}
	}
	return parse
}

func timeWithCustomTZ(value string) time.Time {
	parse, err := time.Parse(time.RFC3339, value)
	if err != nil {
		panic(err)
	}
	_, offset := parse.Zone()
	return parse.In(time.FixedZone("unknown", offset))
}

// TimeP parses a time panics if not valid
func TimeP(value string) *time.Time {
	t := Time(value)
	return &t
}

// ModelTimeUTC parses a model.Time in utc time. panics if not valid
func ModelTimeUTC(value string) model.Time {
	return model.Time(Time(value))
}

// ModelTime parses a model.Time panics if not valid
func ModelTime(value string) model.Time {
	return model.Time(timeWithCustomTZ(value))
}

// ModelTimeP parses a model.Time panics if not valid
func ModelTimeP(value string) *model.Time {
	modelTime := ModelTime(value)
	return &modelTime
}
