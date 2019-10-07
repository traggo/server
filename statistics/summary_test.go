package statistics

import (
	"testing"
	"time"

	"github.com/stretchr/testify/require"
	"github.com/traggo/server/generated/gqlmodel"
	"github.com/traggo/server/model"
	"github.com/traggo/server/test"
	"github.com/traggo/server/test/fake"
)

func TestSummary_fails(t *testing.T) {
	db := test.InMemoryDB(t)
	defer db.Close()

	resolver := ResolverForStatistics{DB: db.DB}

	_, err := resolver.Stats(fake.User(1), []*gqlmodel.Range{rangex("2019-06-11T10:00:00Z", "2019-06-11T09:00:00Z")}, []string{"type"}, nil, nil)
	require.Error(t, err, "range start must be before range end ")
}

func TestStats(t *testing.T) {
	data := []*testEntry{
		name("simple").
			Load(ts("2019-06-11T10:00:00Z", 10*time.Second, tag("type", "review"))).
			Keys("type").
			Ranges(rangex("2019-06-11T10:00:00Z", "2019-06-11T10:00:10Z")).
			Expected(
				result("2019-06-11T10:00:00Z", "2019-06-11T10:00:10Z",
					entry("type", "review", 10*time.Second))),
		name("boundaries 1").
			Load(ts("2019-06-11T10:00:00Z", 10*time.Second, tag("type", "review"))).
			Keys("type").
			Ranges(rangex("2019-06-11T10:00:00Z", "2019-06-11T10:00:05Z")).
			Expected(
				result("2019-06-11T10:00:00Z", "2019-06-11T10:00:05Z",
					entry("type", "review", 5*time.Second))),
		name("boundaries 2").
			Load(ts("2019-06-11T10:00:00Z", 10*time.Second, tag("type", "review"))).
			Keys("type").
			Ranges(rangex("2019-06-11T10:00:05Z", "2019-06-11T10:00:10Z")).
			Expected(
				result("2019-06-11T10:00:05Z", "2019-06-11T10:00:10Z",
					entry("type", "review", 5*time.Second))),
		name("boundaries 3").
			Load(ts("2019-06-11T10:00:00Z", 10*time.Second, tag("type", "review"))).
			Keys("type").
			Ranges(rangex("2019-06-11T10:00:04Z", "2019-06-11T10:00:06Z")).
			Expected(
				result("2019-06-11T10:00:04Z", "2019-06-11T10:00:06Z",
					entry("type", "review", 2*time.Second))),
		name("exclude 1").
			Load(ts("2019-06-11T10:00:00Z", 10*time.Second, tag("type", "review"), tag("issue", "13"))).
			Keys("type").
			Exclude(tag("issue", "13")).
			Ranges(rangex("2019-06-11T10:00:04Z", "2019-06-11T10:00:06Z")).
			Expected(
				result("2019-06-11T10:00:04Z", "2019-06-11T10:00:06Z", []*gqlmodel.StatisticsEntry{}...)),
		name("empty range").
			Load(ts("2019-06-11T10:00:00Z", 10*time.Second, tag("type", "review"), tag("issue", "13"))).
			Keys("type").
			Ranges(rangex("2019-06-11T10:00:04Z", "2019-06-11T10:00:06Z"),
				rangex("2019-06-12T10:00:04Z", "2019-06-12T10:00:06Z")).
			Expected(
				result("2019-06-11T10:00:04Z", "2019-06-11T10:00:06Z", entry("type", "review", 2*time.Second)),
				result("2019-06-12T10:00:04Z", "2019-06-12T10:00:06Z", entry("type", "review", 0*time.Second))),
		name("include 1").
			Load(ts("2019-06-11T10:00:00Z", 10*time.Second, tag("type", "review"), tag("issue", "13"))).
			Keys("type").
			Include(tag("issue", "13")).
			Ranges(rangex("2019-06-11T10:00:00Z", "2019-06-11T10:00:10Z")).
			Expected(
				result("2019-06-11T10:00:00Z", "2019-06-11T10:00:10Z",
					entry("type", "review", 10*time.Second))),
		name("multiple 1").
			Load(
				ts("2019-06-11T10:00:00Z", 10*time.Second, tag("type", "review"), tag("issue", "13")),
				ts("2019-06-11T10:01:00Z", 10*time.Second, tag("type", "review"), tag("issue", "14")),
			).Keys("type").
			Ranges(rangex("2019-06-11T10:00:00Z", "2019-06-11T10:02:00Z")).
			Expected(
				result("2019-06-11T10:00:00Z", "2019-06-11T10:02:00Z",
					entry("type", "review", 20*time.Second))),
		name("multiple 2").
			Load(
				ts("2019-06-11T10:00:00Z", 10*time.Second, tag("type", "review"), tag("issue", "13")),
				ts("2019-06-11T10:01:00Z", 10*time.Second, tag("type", "review"), tag("issue", "14")),
			).Keys("type").
			Exclude(tag("issue", "14")).
			Ranges(rangex("2019-06-11T10:00:00Z", "2019-06-11T10:02:00Z")).
			Expected(
				result("2019-06-11T10:00:00Z", "2019-06-11T10:02:00Z",
					entry("type", "review", 10*time.Second))),
		name("multiple 3").
			Load(
				ts("2019-06-11T10:00:00Z", 10*time.Second, tag("proj", "gotify"), tag("type", "review"), tag("issue", "13")),
				ts("2019-06-11T10:01:00Z", 60*time.Second, tag("proj", "gotify"), tag("type", "review"), tag("issue", "12")),
				ts("2019-06-11T10:02:00Z", 7*time.Second, tag("proj", "gotify"), tag("type", "review"), tag("issue", "13")),
				ts("2019-06-11T10:03:00Z", 14*time.Second, tag("proj", "gotify"), tag("type", "work"), tag("issue", "15")),
				ts("2019-06-11T10:04:00Z", 16*time.Second, tag("proj", "gotify"), tag("type", "work"), tag("issue", "14")),
				ts("2019-06-11T10:05:00Z", 33*time.Second, tag("proj", "gotify"), tag("type", "support"), tag("issue", "10")),
				ts("2019-06-11T10:05:00Z", 12*time.Second, tag("proj", "gotify"), tag("type", "support"), tag("issue", "13")),
				ts("2019-06-11T10:06:00Z", 90*time.Second, tag("type", "break")),
				ts("2019-06-11T10:08:00Z", 33*time.Second, tag("proj", "traggo"), tag("type", "work"), tag("issue", "10")),
				ts("2019-06-11T10:09:00Z", 21*time.Second, tag("proj", "traggo"), tag("type", "work"), tag("issue", "15")),
				ts("2019-06-11T10:10:00Z", 18*time.Second, tag("proj", "traggo"), tag("type", "work"), tag("issue", "12")),
				ts("2019-06-11T10:11:00Z", 56*time.Second, tag("proj", "traggo"), tag("type", "review"), tag("issue", "11")),
				ts("2019-06-11T10:12:00Z", 44*time.Second, tag("type", "break")),
			).Keys("type").
			Ranges(rangex("2019-06-11T10:00:00Z", "2019-06-11T10:15:00Z")).
			Expected(
				result("2019-06-11T10:00:00Z", "2019-06-11T10:15:00Z",
					entry("type", "break", 134*time.Second),
					entry("type", "review", 133*time.Second),
					entry("type", "support", 45*time.Second),
					entry("type", "work", 102*time.Second))),
		name("multiple 4").
			Load(
				ts("2019-06-11T10:00:00Z", 10*time.Second, tag("proj", "gotify"), tag("type", "review"), tag("issue", "13")),
				ts("2019-06-11T10:01:00Z", 60*time.Second, tag("proj", "gotify"), tag("type", "review"), tag("issue", "12")),
				ts("2019-06-11T10:02:00Z", 7*time.Second, tag("proj", "gotify"), tag("type", "review"), tag("issue", "13")),
				ts("2019-06-11T10:03:00Z", 14*time.Second, tag("proj", "gotify"), tag("type", "work"), tag("issue", "15")),
				ts("2019-06-11T10:04:00Z", 16*time.Second, tag("proj", "gotify"), tag("type", "work"), tag("issue", "14")),
				ts("2019-06-11T10:05:00Z", 33*time.Second, tag("proj", "gotify"), tag("type", "support"), tag("issue", "10")),
				ts("2019-06-11T10:05:00Z", 12*time.Second, tag("proj", "gotify"), tag("type", "support"), tag("issue", "13")),
				ts("2019-06-11T10:06:00Z", 90*time.Second, tag("type", "break")),
				ts("2019-06-11T10:08:00Z", 32*time.Second, tag("proj", "traggo"), tag("type", "work"), tag("issue", "10")),
				ts("2019-06-11T10:09:00Z", 21*time.Second, tag("proj", "traggo"), tag("type", "work"), tag("issue", "15")),
				ts("2019-06-11T10:10:00Z", 18*time.Second, tag("proj", "traggo"), tag("type", "work"), tag("issue", "12")),
				ts("2019-06-11T10:11:00Z", 56*time.Second, tag("proj", "traggo"), tag("type", "review"), tag("issue", "11")),
				ts("2019-06-11T10:12:00Z", 44*time.Second, tag("type", "break")),
			).Keys("type", "issue").
			Ranges(rangex("2019-06-11T10:00:00Z", "2019-06-11T10:15:00Z")).
			Expected(
				result("2019-06-11T10:00:00Z", "2019-06-11T10:15:00Z",
					entry("issue", "10", 65*time.Second),
					entry("issue", "11", 56*time.Second),
					entry("issue", "12", 78*time.Second),
					entry("issue", "13", 29*time.Second),
					entry("issue", "14", 16*time.Second),
					entry("issue", "15", 35*time.Second),
					entry("type", "break", 134*time.Second),
					entry("type", "review", 133*time.Second),
					entry("type", "support", 45*time.Second),
					entry("type", "work", 101*time.Second))),
		name("multiple 5").
			Load(
				ts("2019-06-11T10:00:00Z", 10*time.Second, tag("proj", "gotify"), tag("type", "review"), tag("issue", "13")),
				ts("2019-06-11T10:01:00Z", 60*time.Second, tag("proj", "gotify"), tag("type", "review"), tag("issue", "12")),
				ts("2019-06-11T10:02:00Z", 7*time.Second, tag("proj", "gotify"), tag("type", "review"), tag("issue", "13")),
				ts("2019-06-11T10:03:00Z", 14*time.Second, tag("proj", "gotify"), tag("type", "work"), tag("issue", "15")),
				ts("2019-06-11T10:04:00Z", 16*time.Second, tag("proj", "gotify"), tag("type", "work"), tag("issue", "14")),
				ts("2019-06-11T10:05:00Z", 33*time.Second, tag("proj", "gotify"), tag("type", "support"), tag("issue", "10")),
				ts("2019-06-11T10:05:00Z", 12*time.Second, tag("proj", "gotify"), tag("type", "support"), tag("issue", "13")),
				ts("2019-06-11T10:06:00Z", 90*time.Second, tag("type", "break")),
				ts("2019-06-11T10:08:00Z", 32*time.Second, tag("proj", "traggo"), tag("type", "work"), tag("issue", "10")),
				ts("2019-06-11T10:09:00Z", 21*time.Second, tag("proj", "traggo"), tag("type", "work"), tag("issue", "15")),
				ts("2019-06-11T10:10:00Z", 18*time.Second, tag("proj", "traggo"), tag("type", "work"), tag("issue", "12")),
				ts("2019-06-11T10:11:00Z", 56*time.Second, tag("proj", "traggo"), tag("type", "review"), tag("issue", "11")),
				ts("2019-06-11T10:12:00Z", 44*time.Second, tag("type", "break")),
			).Keys("type").
			Ranges(rangex("2019-06-11T10:00:00Z", "2019-06-11T10:15:00Z")).
			Include(tag("issue", "13"), tag("issue", "12")).
			Expected(
				result("2019-06-11T10:00:00Z", "2019-06-11T10:15:00Z",
					entry("type", "review", 77*time.Second),
					entry("type", "support", 12*time.Second),
					entry("type", "work", 18*time.Second))),
		name("multiple 6").
			Load(
				ts("2019-06-11T10:00:00Z", 10*time.Second, tag("proj", "gotify"), tag("type", "review"), tag("issue", "13")),
				ts("2019-06-11T10:01:00Z", 60*time.Second, tag("proj", "gotify"), tag("type", "review"), tag("issue", "12")),
				ts("2019-06-11T10:02:00Z", 7*time.Second, tag("proj", "gotify"), tag("type", "review"), tag("issue", "13")),
				ts("2019-06-11T10:03:00Z", 14*time.Second, tag("proj", "gotify"), tag("type", "work"), tag("issue", "15")),
				ts("2019-06-11T10:04:00Z", 16*time.Second, tag("proj", "gotify"), tag("type", "work"), tag("issue", "14")),
				ts("2019-06-11T10:05:00Z", 33*time.Second, tag("proj", "gotify"), tag("type", "support"), tag("issue", "10")),
				ts("2019-06-11T10:05:00Z", 12*time.Second, tag("proj", "gotify"), tag("type", "support"), tag("issue", "13")),
				ts("2019-06-11T10:06:00Z", 90*time.Second, tag("type", "break")),
				ts("2019-06-11T10:08:00Z", 32*time.Second, tag("proj", "traggo"), tag("type", "work"), tag("issue", "10")),
				ts("2019-06-11T10:09:00Z", 21*time.Second, tag("proj", "traggo"), tag("type", "work"), tag("issue", "15")),
				ts("2019-06-11T10:10:00Z", 18*time.Second, tag("proj", "traggo"), tag("type", "work"), tag("issue", "12")),
				ts("2019-06-11T10:11:00Z", 56*time.Second, tag("proj", "traggo"), tag("type", "review"), tag("issue", "11")),
				ts("2019-06-11T10:12:00Z", 44*time.Second, tag("type", "break")),
			).Keys("type").
			Ranges(rangex("2019-06-11T10:00:00Z", "2019-06-11T10:15:00Z")).
			Exclude(tag("issue", "10"), tag("issue", "11"), tag("issue", "14"), tag("issue", "15")).
			Expected(
				result("2019-06-11T10:00:00Z", "2019-06-11T10:15:00Z",
					entry("type", "break", 134*time.Second),
					entry("type", "review", 77*time.Second),
					entry("type", "support", 12*time.Second),
					entry("type", "work", 18*time.Second))),
	}

	for _, entry := range data {
		t.Run(entry.name, func(t *testing.T) {
			db := test.InMemoryDB(t)
			defer db.Close()

			for _, timespan := range entry.timespans {
				db.Create(&timespan)
			}

			resolver := ResolverForStatistics{DB: db.DB}

			stats, err := resolver.Stats(fake.User(1), entry.ranges, entry.keys, entry.exclude, entry.include)
			require.NoError(t, err)
			require.Len(t, stats, len(entry.expected))
			require.Equal(t, entry.expected, stats)
		})
	}
}

func ts(start string, d time.Duration, tags ...model.TimeSpanTag) model.TimeSpan {
	end := test.Time(start).Add(d)
	return model.TimeSpan{
		StartUserTime: test.Time(start),
		EndUserTime:   &end,
		UserID:        1,
		Tags:          tags,
	}
}

type testEntry struct {
	timespans []model.TimeSpan
	keys      []string
	exclude   []*gqlmodel.InputTimeSpanTag
	include   []*gqlmodel.InputTimeSpanTag
	name      string
	ranges    []*gqlmodel.Range
	expected  []*gqlmodel.RangedStatisticsEntries
}

func name(name string) *testEntry {
	return &testEntry{
		name: name,
	}
}

func (e *testEntry) Expected(expected ...*gqlmodel.RangedStatisticsEntries) *testEntry {
	e.expected = expected
	return e
}

func (e *testEntry) Keys(keys ...string) *testEntry {
	e.keys = keys
	return e
}

func (e *testEntry) Load(timespans ...model.TimeSpan) *testEntry {
	e.timespans = timespans
	return e
}

func (e *testEntry) Ranges(ranges ...*gqlmodel.Range) *testEntry {
	e.ranges = ranges
	return e
}

func (e *testEntry) Include(includes ...model.TimeSpanTag) *testEntry {
	for _, entry := range includes {
		e.include = append(e.include, &gqlmodel.InputTimeSpanTag{
			Key:         entry.Key,
			StringValue: entry.StringValue,
		})
	}
	return e
}
func (e *testEntry) Exclude(excludes ...model.TimeSpanTag) *testEntry {
	for _, entry := range excludes {
		e.exclude = append(e.exclude, &gqlmodel.InputTimeSpanTag{
			Key:         entry.Key,
			StringValue: entry.StringValue,
		})
	}
	return e
}

func tag(key string, value string) model.TimeSpanTag {
	return model.TimeSpanTag{
		Key:         key,
		StringValue: &value,
	}
}

func rangex(start, stop string) *gqlmodel.Range {
	return &gqlmodel.Range{
		Start: test.ModelTime(start),
		End:   test.ModelTime(stop),
	}
}

func result(start, stop string, entries ...*gqlmodel.StatisticsEntry) *gqlmodel.RangedStatisticsEntries {
	return &gqlmodel.RangedStatisticsEntries{
		Start:   test.ModelTimeUTC(start),
		End:     test.ModelTimeUTC(stop),
		Entries: entries,
	}
}

func entry(key string, value string, duration time.Duration) *gqlmodel.StatisticsEntry {
	return &gqlmodel.StatisticsEntry{
		Key:                key,
		StringValue:        &value,
		TimeSpendInSeconds: duration.Seconds(),
	}
}
