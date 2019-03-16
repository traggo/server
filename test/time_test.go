package test

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestTime_panics(t *testing.T) {
	assert.Panics(t, func() {
		Time("asds")
	})
}

func TestTime_succeeds(t *testing.T) {
	Time("2009-06-30T18:30:00+02:00")
}

func TestModelTime_panics(t *testing.T) {
	assert.Panics(t, func() {
		ModelTime("asds")
	})
}

func TestModelTime_succeeds(t *testing.T) {
	ModelTime("2009-06-30T18:30:00+02:00")
}

func TestTimeP_panics(t *testing.T) {
	assert.Panics(t, func() {
		TimeP("asds")
	})
}

func TestTimeP_succeeds(t *testing.T) {
	TimeP("2009-06-30T18:30:00+02:00")
}

func TestModelTimeP_panics(t *testing.T) {
	assert.Panics(t, func() {
		ModelTimeP("asds")
	})
}

func TestModelTimeP_succeeds(t *testing.T) {
	ModelTimeP("2009-06-30T18:30:00+02:00")
}
