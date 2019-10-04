package time

import (
	"fmt"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/traggo/server/model"
	"github.com/traggo/server/test"
)

func TestParseRange(t *testing.T) {

	entries := []testEntry{
		now("2019-05-13T15:55:23Z").
			Range("now-1d", "now").
			Interval(model.IntervalSingle).
			Expect(r("2019-05-12T15:55:23Z", "2019-05-13T15:55:23Z")),
		now("2019-05-13T15:55:23Z").
			Range("now-1d", "2019-05-13T18:52:23Z").
			Interval(model.IntervalSingle).
			Expect(r("2019-05-12T15:55:23Z", "2019-05-13T18:52:23Z")),
		now("2019-05-13T15:55:23Z").
			Range("now-1w", "now").
			Interval(model.IntervalSingle).
			Expect(r("2019-05-06T15:55:23Z", "2019-05-13T15:55:23Z")),
		now("2019-05-13T15:55:23Z").
			Range("now-10d", "now").
			Interval(model.IntervalSingle).
			Expect(r("2019-05-03T15:55:23Z", "2019-05-13T15:55:23Z")),
		now("2019-05-13T15:55:23Z").
			Range("now-1M", "now").
			Interval(model.IntervalSingle).
			Expect(r("2019-04-13T15:55:23Z", "2019-05-13T15:55:23Z")),
		now("2019-05-13T15:55:23Z").
			Range("now-1y", "now").
			Interval(model.IntervalSingle).
			Expect(r("2018-05-13T15:55:23Z", "2019-05-13T15:55:23Z")),
		now("2019-05-13T15:55:23Z").
			Range("now-5s-3m-5h-2d-2w-3M-1y", "now").
			Interval(model.IntervalSingle).
			Expect(r("2018-01-27T10:52:18Z", "2019-05-13T15:55:23Z")),
		now("2019-05-13T15:55:23Z").
			Range("now/s", "now").
			Interval(model.IntervalSingle).
			Expect(r("2019-05-13T15:55:23Z", "2019-05-13T15:55:23Z")),
		now("2019-05-13T15:55:23Z").
			Range("now/m", "now").
			Interval(model.IntervalSingle).
			Expect(r("2019-05-13T15:55:00Z", "2019-05-13T15:55:23Z")),
		now("2019-05-13T15:55:23Z").
			Range("now/h", "now").
			Interval(model.IntervalSingle).
			Expect(r("2019-05-13T15:00:00Z", "2019-05-13T15:55:23Z")),
		now("2019-05-13T15:55:23Z").
			Range("now/d", "now").
			Interval(model.IntervalSingle).
			Expect(r("2019-05-13T00:00:00Z", "2019-05-13T15:55:23Z")),
		now("2019-05-15T15:55:23Z").
			Range("now/w", "now").
			Interval(model.IntervalSingle).
			Expect(r("2019-05-13T00:00:00Z", "2019-05-15T15:55:23Z")),
		now("2019-05-15T15:55:23Z").
			Range("now/M", "now").
			Interval(model.IntervalSingle).
			Expect(r("2019-05-01T00:00:00Z", "2019-05-15T15:55:23Z")),
		now("2019-05-15T15:55:23Z").
			Range("now/y", "now").
			Interval(model.IntervalSingle).
			Expect(r("2019-01-01T00:00:00Z", "2019-05-15T15:55:23Z")),
		now("2019-05-13T15:55:23Z").
			Range("now", "now/s").
			Interval(model.IntervalSingle).
			Expect(r("2019-05-13T15:55:23Z", "2019-05-13T15:55:23.999999999Z")),
		now("2019-05-13T15:55:23Z").
			Range("now", "now/m").
			Interval(model.IntervalSingle).
			Expect(r("2019-05-13T15:55:23Z", "2019-05-13T15:55:59.999999999Z")),
		now("2019-05-13T15:55:23Z").
			Range("now", "now/h").
			Interval(model.IntervalSingle).
			Expect(r("2019-05-13T15:55:23Z", "2019-05-13T15:59:59.999999999Z")),
		now("2019-05-13T15:55:23Z").
			Range("now", "now/d").
			Interval(model.IntervalSingle).
			Expect(r("2019-05-13T15:55:23Z", "2019-05-13T23:59:59.999999999Z")),
		now("2019-05-13T15:55:23Z").
			Range("now", "now/d").
			Interval(model.IntervalSingle).
			Expect(r("2019-05-13T15:55:23Z", "2019-05-13T23:59:59.999999999Z")),
		now("2019-05-15T15:55:23Z").
			Range("now", "now/w").
			Interval(model.IntervalSingle).
			Expect(r("2019-05-15T15:55:23Z", "2019-05-19T23:59:59.999999999Z")),
		now("2019-05-15T15:55:23Z").
			Range("now", "now/M").
			Interval(model.IntervalSingle).
			Expect(r("2019-05-15T15:55:23Z", "2019-05-31T23:59:59.999999999Z")),
		now("2019-05-15T15:55:23Z").
			Range("now", "now/y").
			Interval(model.IntervalSingle).
			Expect(r("2019-05-15T15:55:23Z", "2019-12-31T23:59:59.999999999Z")),
		now("2019-05-15T15:55:23Z").
			Range("now-1y/y", "now-1y/y"). // last year
			Interval(model.IntervalSingle).
			Expect(r("2018-01-01T00:00:00Z", "2018-12-31T23:59:59.999999999Z")),
		now("2019-05-15T15:55:23Z").
			Range("now-1w/w", "now-1w/w"). // last week
			Interval(model.IntervalSingle).
			Expect(r("2019-05-06T00:00:00Z", "2019-05-12T23:59:59.999999999Z")),
		now("2019-05-15T15:55:23Z").
			Range("now-1w/w", "now-1w/w-2d"). // last work monday to friday
			Interval(model.IntervalSingle).
			Expect(r("2019-05-06T00:00:00Z", "2019-05-10T23:59:59.999999999Z")),
		now("2019-05-15T15:55:23Z").
			Range("now-1w/w+5d", "now-1w/w"). // last weekend
			Interval(model.IntervalSingle).
			Expect(r("2019-05-11T00:00:00Z", "2019-05-12T23:59:59.999999999Z")),
		now("2019-05-13T15:55:23Z").
			Range("now-1d", "now").
			Interval(model.IntervalYearly).
			Expect(r("2019-05-12T15:55:23Z", "2019-05-13T15:55:23Z")),
		now("2019-05-13T15:55:23Z").
			Range("now-1d", "now").
			Interval(model.IntervalMonthly).
			Expect(r("2019-05-12T15:55:23Z", "2019-05-13T15:55:23Z")),
		now("2019-05-13T15:55:23Z").
			Range("now-1d", "now").
			Interval(model.IntervalWeekly).
			Expect(r("2019-05-12T15:55:23Z", "2019-05-13T15:55:23Z")),
		now("2019-05-13T15:55:23Z").
			Range("now-5h/h", "now/h").
			Interval(model.IntervalHourly).
			Expect(
				r("2019-05-13T10:00:00Z", "2019-05-13T10:59:59Z"),
				r("2019-05-13T11:00:00Z", "2019-05-13T11:59:59Z"),
				r("2019-05-13T12:00:00Z", "2019-05-13T12:59:59Z"),
				r("2019-05-13T13:00:00Z", "2019-05-13T13:59:59Z"),
				r("2019-05-13T14:00:00Z", "2019-05-13T14:59:59Z"),
				r("2019-05-13T15:00:00Z", "2019-05-13T15:59:59Z")),
		now("2019-05-13T15:55:23Z").
			Range("now-5d/d", "now/d").
			Interval(model.IntervalDaily).
			Expect(
				r("2019-05-08T00:00:00Z", "2019-05-08T23:59:59Z"),
				r("2019-05-09T00:00:00Z", "2019-05-09T23:59:59Z"),
				r("2019-05-10T00:00:00Z", "2019-05-10T23:59:59Z"),
				r("2019-05-11T00:00:00Z", "2019-05-11T23:59:59Z"),
				r("2019-05-12T00:00:00Z", "2019-05-12T23:59:59Z"),
				r("2019-05-13T00:00:00Z", "2019-05-13T23:59:59Z")),
		now("2019-02-13T15:55:23Z").
			Range("now-5M/M", "now/M").
			Interval(model.IntervalMonthly).
			Expect(
				r("2018-09-01T00:00:00Z", "2018-09-30T23:59:59Z"),
				r("2018-10-01T00:00:00Z", "2018-10-31T23:59:59Z"),
				r("2018-11-01T00:00:00Z", "2018-11-30T23:59:59Z"),
				r("2018-12-01T00:00:00Z", "2018-12-31T23:59:59Z"),
				r("2019-01-01T00:00:00Z", "2019-01-31T23:59:59Z"),
				r("2019-02-01T00:00:00Z", "2019-02-28T23:59:59Z")),
	}

	for _, entry := range entries {
		t.Run(fmt.Sprintf("now=%s;%s_to_%s;%s", entry.now, entry.from, entry.to, entry.interval), func(t *testing.T) {
			result, err := ParseRange(entry.now, RelativeRange{From: entry.from, To: entry.to}, entry.interval)
			assert.NoError(t, err)
			assert.Equal(t, entry.expectRanges, result)
		})
	}
}

func TestRanges(t *testing.T) {
	assert.Panics(t, func() {
		ranges(time.Time{}, time.Time{}, "meh")
	})
}

func now(now string) testEntry {
	return testEntry{now: test.Time(now)}
}

type testEntry struct {
	now          time.Time
	from         string
	to           string
	interval     model.Interval
	expectRanges []StaticRange
}

func r(from string, to string) StaticRange {
	return StaticRange{
		From: test.Time(from),
		To:   test.Time(to),
	}
}

func (t testEntry) Expect(ranges ...StaticRange) testEntry {
	t.expectRanges = ranges
	return t
}

func (t testEntry) Range(from string, to string) testEntry {
	t.from = from
	t.to = to
	return t
}

func (t testEntry) Interval(interval model.Interval) testEntry {
	t.interval = interval
	return t
}
