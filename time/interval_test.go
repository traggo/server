package time

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/traggo/server/generated/gqlmodel"
	"github.com/traggo/server/model"
)

func TestInterval(t *testing.T) {
	for _, inter := range []model.Interval{
		model.IntervalSingle, model.IntervalHourly, model.IntervalDaily, model.IntervalWeekly, model.IntervalMonthly, model.IntervalYearly,
	} {
		assert.Equal(t, inter, InternalInterval(ExternalInterval(inter)))
	}
	assert.Panics(t, func() {
		InternalInterval(gqlmodel.StatsInterval("meh"))
	})
	assert.Panics(t, func() {
		ExternalInterval(model.Interval("meh"))
	})
}
