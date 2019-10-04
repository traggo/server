package time

import (
	"time"

	"github.com/jmattheis/go-timemath"
	"github.com/traggo/server/model"
)

func ranges(from time.Time, to time.Time, interval model.Interval) []StaticRange {
	switch interval {
	case model.IntervalSingle:
		return []StaticRange{{From: from, To: to}}
	case model.IntervalHourly:
		return rangeForUnit(from, to, timemath.Hour)
	case model.IntervalDaily:
		return rangeForUnit(from, to, timemath.Day)
	case model.IntervalWeekly:
		return rangeForUnit(from, to, timemath.Week)
	case model.IntervalMonthly:
		return rangeForUnit(from, to, timemath.Month)
	case model.IntervalYearly:
		return rangeForUnit(from, to, timemath.Year)
	default:
		panic("unknown interval type")
	}
}

func rangeForUnit(from time.Time, to time.Time, u timemath.Unit) []StaticRange {
	var result []StaticRange
	newFrom := from
	for newFrom.Before(to) {
		newTo := timemath.Second.Subtract(u.Add(newFrom, 1), 1)
		if newTo.After(to) {
			newTo = to
		}
		result = append(result, StaticRange{From: newFrom, To: newTo})
		newFrom = timemath.Second.Add(newTo, 1)
	}
	return result
}
