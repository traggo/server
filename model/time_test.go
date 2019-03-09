package model

import (
	"bytes"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestTime_MarshalGQL(t *testing.T) {
	var buffer bytes.Buffer
	expected := "2009-06-30T18:30:00+02:00"
	parse, err := time.Parse(time.RFC3339, expected)
	assert.Nil(t, err)
	toTest := Time(parse)
	toTest.MarshalGQL(&buffer)
	actual := buffer.String()
	assert.Equal(t, expected, actual)
}

func TestTime_UnmarshalGQL_success(t *testing.T) {
	date := "2009-06-30T18:30:00+02:00"
	parse, err := time.Parse(time.RFC3339, date)
	assert.Nil(t, err)
	expected := Time(parse)

	actual := &Time{}
	err = actual.UnmarshalGQL(date)
	assert.Nil(t, err)

	assert.Equal(t, expected, *actual)
}

func TestTime_OmitTimeZone(t *testing.T) {
	date := "2009-06-30T18:30:00+02:00"
	tzDate, err := time.Parse(time.RFC3339, date)
	assert.Nil(t, err)
	utcDate := "2009-06-30T18:30:00Z"
	withoutTz, err := time.Parse(time.RFC3339, utcDate)
	assert.Nil(t, err)

	assert.Equal(t, withoutTz, Time(tzDate).OmitTimeZone())
}

func TestTime_UTC(t *testing.T) {
	date := "2009-06-30T18:30:00+02:00"
	parse, err := time.Parse(time.RFC3339, date)
	assert.Nil(t, err)
	without := "2009-06-30T16:30:00Z"
	withoutTz, err := time.Parse(time.RFC3339, without)
	assert.Nil(t, err)

	assert.Equal(t, withoutTz, Time(parse).UTC())
}

func TestTime_UnmarshalGQL_failInvalidType(t *testing.T) {
	actual := &Time{}
	err := actual.UnmarshalGQL(1)
	assert.EqualError(t, err, "time must be a string")
}

func TestTime_UnmarshalGQL_invalidFormat(t *testing.T) {
	actual := &Time{}
	err := actual.UnmarshalGQL("lol")
	assert.EqualError(t, err, `parsing time "lol" as "2006-01-02T15:04:05Z07:00": cannot parse "lol" as "2006"`)
}
