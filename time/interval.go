package time

import (
	"github.com/traggo/server/generated/gqlmodel"
	"github.com/traggo/server/model"
)

// InternalInterval converts gqlmodel to internal
func InternalInterval(interval gqlmodel.StatsInterval) model.Interval {
	switch interval {
	case gqlmodel.StatsIntervalHourly:
		return model.IntervalHourly
	case gqlmodel.StatsIntervalDaily:
		return model.IntervalDaily
	case gqlmodel.StatsIntervalWeekly:
		return model.IntervalWeekly
	case gqlmodel.StatsIntervalMonthly:
		return model.IntervalMonthly
	case gqlmodel.StatsIntervalYearly:
		return model.IntervalYearly
	case gqlmodel.StatsIntervalSingle:
		return model.IntervalSingle
	default:
		panic("unknown interval type " + interval)
	}
}

// ExternalInterval converts internal to gqlmodel
func ExternalInterval(interval model.Interval) gqlmodel.StatsInterval {
	switch interval {
	case model.IntervalHourly:
		return gqlmodel.StatsIntervalHourly
	case model.IntervalDaily:
		return gqlmodel.StatsIntervalDaily
	case model.IntervalWeekly:
		return gqlmodel.StatsIntervalWeekly
	case model.IntervalMonthly:
		return gqlmodel.StatsIntervalMonthly
	case model.IntervalYearly:
		return gqlmodel.StatsIntervalYearly
	case model.IntervalSingle:
		return gqlmodel.StatsIntervalSingle
	default:
		panic("unknown interval type " + interval)
	}
}
