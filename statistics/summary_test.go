package statistics

import (
	"fmt"
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
	"github.com/traggo/server/generated/gqlmodel"
	"github.com/traggo/server/test"
	"github.com/traggo/server/test/fake"
)

func testData(db *test.Database) {
	user := db.User(1)

	user.TimeSpan("2019-06-11T10:00:00Z", "2019-06-11T11:00:00Z"). // 60 min
									Tag("proj", "gotify").Tag("issue", "1").Tag("type", "work")

	user.TimeSpan("2019-06-11T11:01:00Z", "2019-06-11T11:10:00Z"). // 9 min
									Tag("proj", "gotify").Tag("issue", "2").Tag("type", "work")

	user.TimeSpan("2019-06-11T11:15:00Z", "2019-06-11T11:30:00Z"). // 15 min
									Tag("type", "break")

	user.TimeSpan("2019-06-11T11:31:00Z", "2019-06-11T11:44:00Z"). // 13 min
									Tag("proj", "gotify").Tag("issue", "3").Tag("type", "review")

	user.TimeSpan("2019-06-11T11:45:00Z", "2019-06-11T11:51:00Z"). // 6 min
									Tag("proj", "gotify").Tag("issue", "3").Tag("type", "support")

	user.TimeSpan("2019-06-11T12:00:00Z", "2019-06-11T12:21:00Z"). // 21 min
									Tag("proj", "gotify").Tag("issue", "2").Tag("type", "work")

	user.TimeSpan("2019-06-11T12:30:00Z", "2019-06-11T12:47:00Z"). // 17 min
									Tag("proj", "gotify").Tag("issue", "1").Tag("type", "support")

	user.TimeSpan("2019-06-11T12:48:00Z", "2019-06-11T13:00:00Z"). // 12 min
									Tag("proj", "gotify").Tag("issue", "3").Tag("type", "review")

	user.TimeSpan("2019-06-11T13:01:00Z", "2019-06-11T13:30:00Z"). // 29 min
									Tag("type", "break")

	user.TimeSpan("2019-06-11T13:34:00Z", "2019-06-11T13:39:00Z"). // 5 min
									Tag("proj", "gotify").Tag("issue", "12").Tag("type", "work")

	user.TimeSpan("2019-06-11T13:40:00Z", "2019-06-11T14:00:00Z"). // 20 min
									Tag("proj", "gotify").Tag("issue", "55").Tag("type", "review")

	user.TimeSpan("2019-06-11T14:01:00Z", "2019-06-11T14:11:00Z"). // 10 min
									Tag("proj", "gotify").Tag("issue", "12").Tag("type", "review")

	user.TimeSpan("2019-06-11T14:22:00Z", "2019-06-11T14:29:00Z"). // 7 min
									Tag("proj", "gotify").Tag("issue", "55").Tag("type", "support")

	user.TimeSpan("2019-06-11T14:30:00Z", "2019-06-11T14:44:00Z"). // 14 min
									Tag("proj", "traggo").Tag("issue", "55").Tag("type", "review")

	user.TimeSpan("2019-06-11T14:47:00Z", "2019-06-11T15:07:00Z"). // 20 min
									Tag("proj", "traggo").Tag("issue", "55").Tag("type", "review")

	user.TimeSpan("2019-06-11T15:11:00Z", "2019-06-11T15:41:00Z"). // 30 min
									Tag("proj", "traggo").Tag("issue", "33").Tag("type", "work")

	user.TimeSpan("2019-06-11T16:01:00Z", "2019-06-11T17:00:00Z"). // 59 min
									Tag("type", "break")

	user.TimeSpan("2019-06-15T12:00:00Z", "2019-06-16T12:00:00Z"). // 1 day
									Tag("type", "break")
}

func TestSummary(t *testing.T) {
	d := []sData{
		{
			From: "2019-06-11T10:00:00Z",
			To:   "2019-06-11T20:00:00Z",
			Key:  "proj",
			Expected: []*gqlmodel.StatisticsEntry{
				entry("proj", "gotify", 180*time.Minute),
				entry("proj", "traggo", 64*time.Minute),
			},
		},
		{
			From: "2019-06-11T10:00:00Z",
			To:   "2019-06-11T20:00:00Z",
			Key:  "proj",
			Has: []*gqlmodel.InputTimeSpanTag{
				tag("type", "review"),
			},
			Expected: []*gqlmodel.StatisticsEntry{
				entry("proj", "gotify", 55*time.Minute),
				entry("proj", "traggo", 34*time.Minute),
			},
		},
		{
			From: "2019-06-11T10:00:00Z",
			To:   "2019-06-11T20:00:00Z",
			Key:  "proj",
			Has: []*gqlmodel.InputTimeSpanTag{
				tag("type", "work"),
			},
			Expected: []*gqlmodel.StatisticsEntry{
				entry("proj", "gotify", 95*time.Minute),
				entry("proj", "traggo", 30*time.Minute),
			},
		},
		{
			From: "2019-06-11T10:00:00Z",
			To:   "2019-06-11T20:00:00Z",
			Key:  "proj",
			NotHas: []*gqlmodel.InputTimeSpanTag{
				tag("type", "work"),
			},
			Expected: []*gqlmodel.StatisticsEntry{
				entry("proj", "gotify", 85*time.Minute),
				entry("proj", "traggo", 34*time.Minute),
			},
		},
		{
			From: "2019-06-11T10:00:00Z",
			To:   "2019-06-11T20:00:00Z",
			Key:  "proj",
			Has: []*gqlmodel.InputTimeSpanTag{
				tag("type", "review"),
			},
			NotHas: []*gqlmodel.InputTimeSpanTag{
				tag("issue", "12"),
			},
			Expected: []*gqlmodel.StatisticsEntry{
				entry("proj", "gotify", 45*time.Minute),
				entry("proj", "traggo", 34*time.Minute),
			},
		},
		{
			From: "2019-06-11T10:00:00Z",
			To:   "2019-06-11T20:00:00Z",
			Key:  "type",
			Expected: []*gqlmodel.StatisticsEntry{
				entry("type", "break", 103*time.Minute),
				entry("type", "review", 89*time.Minute),
				entry("type", "support", 30*time.Minute),
				entry("type", "work", 125*time.Minute),
			},
		},
		{
			From: "2019-06-11T10:00:00Z",
			To:   "2019-06-11T20:00:00Z",
			Key:  "type",
			Has: []*gqlmodel.InputTimeSpanTag{
				tag("proj", "gotify"),
			},
			Expected: []*gqlmodel.StatisticsEntry{
				entry("type", "review", 55*time.Minute),
				entry("type", "support", 30*time.Minute),
				entry("type", "work", 95*time.Minute),
			},
		},
		{
			From: "2019-06-11T10:00:00Z",
			To:   "2019-06-11T20:00:00Z",
			Key:  "type",
			Has: []*gqlmodel.InputTimeSpanTag{
				tag("proj", "traggo"),
			},
			Expected: []*gqlmodel.StatisticsEntry{
				entry("type", "review", 34*time.Minute),
				entry("type", "work", 30*time.Minute),
			},
		},
		{
			From: "2019-06-11T10:00:00Z",
			To:   "2019-06-11T20:00:00Z",
			Key:  "issue",
			Expected: []*gqlmodel.StatisticsEntry{
				entry("issue", "1", 77*time.Minute),
				entry("issue", "12", 15*time.Minute),
				entry("issue", "2", 30*time.Minute),
				entry("issue", "3", 31*time.Minute),
				entry("issue", "33", 30*time.Minute),
				entry("issue", "55", 61*time.Minute),
			},
		},
		{
			From: "2019-06-11T10:00:00Z",
			To:   "2019-06-11T20:00:00Z",
			Key:  "issue",
			NotHas: []*gqlmodel.InputTimeSpanTag{
				tag("proj", "traggo"),
				tag("type", "work"),
			},
			Expected: []*gqlmodel.StatisticsEntry{
				entry("issue", "1", 17*time.Minute),
				entry("issue", "12", 10*time.Minute),
				entry("issue", "3", 31*time.Minute),
				entry("issue", "55", 27*time.Minute),
			},
		},
		{
			From: "2019-06-15T12:00:00Z",
			To:   "2019-06-16T00:00:00Z",
			Key:  "type",
			Expected: []*gqlmodel.StatisticsEntry{
				entry("type", "break", 12*time.Hour),
			},
		},
	}

	for _, data := range d {
		t.Run(data.String(), func(t *testing.T) {
			db := test.InMemoryDB(t)
			defer db.Close()
			testData(db)

			resolver := ResolverForStatistics{DB: db.DB}
			stat := gqlmodel.StatInput{Key: data.Key, MustHave: data.Has, MustNotHave: data.NotHas}

			stats, err := resolver.TimeSpanSummary(fake.User(1), test.ModelTime(data.From),
				test.ModelTime(data.To), stat)
			require.NoError(t, err)

			require.Equal(t, data.Expected, stats)
		})
	}
}

func TestSummary_fails_toBeforeFrom(t *testing.T) {
	db := test.InMemoryDB(t)
	defer db.Close()

	resolver := ResolverForStatistics{DB: db.DB}
	stat := gqlmodel.StatInput{Key: "proj"}

	_, err := resolver.TimeSpanSummary(fake.User(1), test.ModelTime("2019-06-15T13:00:00Z"),
		test.ModelTime("2019-06-15T12:00:00Z"), stat)
	require.Error(t, err, "toInclusive must be after fromInclusive")
}

type sData struct {
	From     string
	To       string
	Key      string
	Has      []*gqlmodel.InputTimeSpanTag
	NotHas   []*gqlmodel.InputTimeSpanTag
	Expected []*gqlmodel.StatisticsEntry
}

func (d sData) String() string {
	var expect []string
	for _, entry := range d.Expected {
		duration := time.Duration(entry.TimeSpendInSeconds) * time.Second
		if entry.StringValue == nil {
			expect = append(expect, fmt.Sprintf("%s=%s", entry.Key, duration))
		} else {
			expect = append(expect, fmt.Sprintf("%s:%s=%s", entry.Key, *entry.StringValue, duration))
		}
	}

	fTag := func(tags []*gqlmodel.InputTimeSpanTag) string {
		var result []string
		for _, tag := range tags {
			if tag.StringValue == nil {
				result = append(result, fmt.Sprintf("%s", tag.Key))
			} else {
				result = append(result, fmt.Sprintf("%s:%s", tag.Key, *tag.StringValue))
			}
		}
		return strings.Join(result, ",")
	}

	return fmt.Sprintf("From=%s|To=%s|Key=%s|Has=[%s]|NotHas=[%s]|Expected=[%s]",
		d.From, d.To, d.Key, fTag(d.Has), fTag(d.NotHas), strings.Join(expect, ","))
}

func entry(key string, value string, duration time.Duration) *gqlmodel.StatisticsEntry {
	return &gqlmodel.StatisticsEntry{
		Key:                key,
		StringValue:        &value,
		TimeSpendInSeconds: int(duration.Seconds()),
	}
}

func tag(key string, value string) *gqlmodel.InputTimeSpanTag {
	return &gqlmodel.InputTimeSpanTag{
		Key:         key,
		StringValue: &value,
	}
}
