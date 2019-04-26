package timespan

import (
	"fmt"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
	"github.com/traggo/server/generated/gqlmodel"
	"github.com/traggo/server/model"
	"github.com/traggo/server/test"
	"github.com/traggo/server/test/fake"
)

var (
	timeSpan1 = &model.TimeSpan{
		ID:            1,
		UserID:        5,
		StartUserTime: test.Time("2019-06-10T18:30:00Z"),
		EndUserTime:   test.TimeP("2019-06-10T18:40:00Z"),
		StartUTC:      test.Time("2019-06-10T18:30:00Z"),
		EndUTC:        test.TimeP("2019-06-10T18:40:00Z"),
		OffsetUTC:     0,
		Tags: []model.TimeSpanTag{
			{Key: "test"},
			{Key: "test2"},
		},
	}
	modelTimeSpan1 = gqlmodel.TimeSpan{
		ID:    1,
		Start: test.ModelTime("2019-06-10T18:30:00Z"),
		End:   test.ModelTimeP("2019-06-10T18:40:00Z"),
		Tags: []gqlmodel.TimeSpanTag{
			{Key: "test"},
			{Key: "test2"},
		},
	}
	timeSpan2 = &model.TimeSpan{
		ID:            2,
		UserID:        5,
		StartUserTime: test.Time("2019-06-10T18:40:00Z"),
		EndUserTime:   test.TimeP("2019-06-10T18:50:00Z"),
		StartUTC:      test.Time("2019-06-10T18:40:00Z"),
		EndUTC:        test.TimeP("2019-06-10T18:50:00Z"),
		OffsetUTC:     0,
		Tags: []model.TimeSpanTag{
			{Key: "test"},
			{Key: "test2"},
		},
	}
	modelTimeSpan2 = gqlmodel.TimeSpan{
		ID:    2,
		Start: test.ModelTime("2019-06-10T18:40:00Z"),
		End:   test.ModelTimeP("2019-06-10T18:50:00Z"),
		Tags: []gqlmodel.TimeSpanTag{
			{Key: "test"},
			{Key: "test2"},
		},
	}
	timeSpan3 = &model.TimeSpan{
		ID:            3,
		UserID:        5,
		StartUserTime: test.Time("2019-06-10T18:50:00Z"),
		EndUserTime:   test.TimeP("2019-06-10T19:00:00Z"),
		StartUTC:      test.Time("2019-06-10T18:50:00Z"),
		EndUTC:        test.TimeP("2019-06-10T19:00:00Z"),
		OffsetUTC:     0,
		Tags: []model.TimeSpanTag{
			{Key: "test"},
			{Key: "test2"},
		},
	}
	modelTimeSpan3 = gqlmodel.TimeSpan{
		ID:    3,
		Start: test.ModelTime("2019-06-10T18:50:00Z"),
		End:   test.ModelTimeP("2019-06-10T19:00:00Z"),
		Tags: []gqlmodel.TimeSpanTag{
			{Key: "test"},
			{Key: "test2"},
		},
	}
	timeSpan4 = &model.TimeSpan{
		ID:            4,
		UserID:        5,
		StartUserTime: test.Time("2019-06-10T18:00:00Z"),
		EndUserTime:   test.TimeP("2019-06-10T19:00:00Z"),
		StartUTC:      test.Time("2019-06-10T18:00:00Z"),
		EndUTC:        test.TimeP("2019-06-10T19:00:00Z"),
		OffsetUTC:     0,
		Tags: []model.TimeSpanTag{
			{Key: "test"},
			{Key: "test2"},
		},
	}
	modelTimeSpan4 = gqlmodel.TimeSpan{
		ID:    4,
		Start: test.ModelTime("2019-06-10T18:00:00Z"),
		End:   test.ModelTimeP("2019-06-10T19:00:00Z"),
		Tags: []gqlmodel.TimeSpanTag{
			{Key: "test"},
			{Key: "test2"},
		},
	}
	runningTimeSpan = &model.TimeSpan{
		ID:            5,
		UserID:        5,
		StartUserTime: test.Time("2019-06-11T18:00:00Z"),
		StartUTC:      test.Time("2019-06-11T18:00:00Z"),
		EndUserTime:   nil,
		EndUTC:        nil,
		OffsetUTC:     0,
		Tags: []model.TimeSpanTag{
			{Key: "test"},
			{Key: "test2"},
		},
	}
	modelRunningTimeSpan = gqlmodel.TimeSpan{
		ID:    5,
		Start: test.ModelTime("2019-06-11T18:00:00Z"),
		End:   nil,
		Tags: []gqlmodel.TimeSpanTag{
			{Key: "test"},
			{Key: "test2"},
		},
	}
	timeSpanOtherUser = &model.TimeSpan{
		ID:            6,
		UserID:        2,
		StartUserTime: test.Time("2019-06-10T18:30:00Z"),
		EndUserTime:   test.TimeP("2019-06-10T18:40:00Z"),
		StartUTC:      test.Time("2019-06-10T18:30:00Z"),
		EndUTC:        test.TimeP("2019-06-10T18:40:00Z"),
		OffsetUTC:     0,
		Tags: []model.TimeSpanTag{
			{Key: "test"},
			{Key: "test2"},
		},
	}
)

type data struct {
	DB       []*model.TimeSpan
	From     *model.Time
	To       *model.Time
	Expected []gqlmodel.TimeSpan
}

func (d data) String() string {
	str := "DB=["
	for _, v := range d.DB {
		str += fmt.Sprint(v.ID) + ","
	}
	str = str + "]"
	if d.From != nil {
		str = str + "|From=" + d.From.Time().Format(time.RFC3339)
	}
	if d.To != nil {
		str = str + "|To=" + d.To.Time().Format(time.RFC3339)
	}

	str += "|Expected=["
	for _, v := range d.Expected {
		str += fmt.Sprint(v.ID) + ","
	}
	str = str + "]"
	return str
}

func TestGet(t *testing.T) {

	d := []data{
		{
			DB:       []*model.TimeSpan{timeSpan1, timeSpanOtherUser},
			From:     nil,
			To:       nil,
			Expected: []gqlmodel.TimeSpan{modelTimeSpan1},
		},
		{
			DB:       []*model.TimeSpan{timeSpan1, timeSpan2, timeSpan3, timeSpan4, runningTimeSpan},
			From:     nil,
			To:       nil,
			Expected: []gqlmodel.TimeSpan{modelTimeSpan3, modelTimeSpan2, modelTimeSpan1, modelTimeSpan4},
		},
		{
			DB:       []*model.TimeSpan{timeSpan1, timeSpan2, timeSpan3, timeSpan4},
			From:     test.ModelTimeP("2019-06-10T18:40:01+02:00"),
			To:       nil,
			Expected: []gqlmodel.TimeSpan{modelTimeSpan3, modelTimeSpan2, modelTimeSpan4},
		},
		{
			DB:       []*model.TimeSpan{timeSpan1, timeSpan2, timeSpan3, timeSpan4},
			From:     modelTimeSpan1.End,
			To:       nil,
			Expected: []gqlmodel.TimeSpan{modelTimeSpan3, modelTimeSpan2, modelTimeSpan1, modelTimeSpan4},
		},
		{
			DB:       []*model.TimeSpan{timeSpan1, timeSpan2, timeSpan3, timeSpan4},
			From:     nil,
			To:       test.ModelTimeP("2019-06-10T18:45:00+02:00"),
			Expected: []gqlmodel.TimeSpan{modelTimeSpan2, modelTimeSpan1, modelTimeSpan4},
		},
		{
			DB:       []*model.TimeSpan{timeSpan1, timeSpan2, timeSpan3, timeSpan4},
			From:     test.ModelTimeP("2019-06-10T18:35:00+02:00"),
			To:       test.ModelTimeP("2019-06-10T18:45:00+02:00"),
			Expected: []gqlmodel.TimeSpan{modelTimeSpan2, modelTimeSpan1, modelTimeSpan4},
		},
		{
			DB:       []*model.TimeSpan{timeSpan1, timeSpan2, timeSpan3, timeSpan4},
			From:     test.ModelTimeP("2019-06-10T18:55:00+02:00"),
			To:       test.ModelTimeP("2019-06-10T18:56:00+02:00"),
			Expected: []gqlmodel.TimeSpan{modelTimeSpan3, modelTimeSpan4},
		},
		{
			DB:       []*model.TimeSpan{runningTimeSpan},
			From:     test.ModelTimeP("2019-06-11T18:30:00Z"),
			To:       test.ModelTimeP("2019-06-11T19:00:00Z"),
			Expected: nil,
		},
	}

	for i, testData := range d {
		t.Run(fmt.Sprintf("%d>%s", i, testData), func(t *testing.T) {
			db := test.InMemoryDB(t)
			db.User(5)
			defer db.Close()
			for _, entry := range testData.DB {
				db.Create(entry)
			}

			resolver := ResolverForTimeSpan{DB: db.DB}
			timeSpans, err := resolver.TimeSpans(fake.User(5), testData.From, testData.To)

			require.NoError(t, err)
			require.Equal(t, testData.Expected, timeSpans)
		})
	}
}

func TestGet_fail_toBeforeStart(t *testing.T) {
	db := test.InMemoryDB(t)
	defer db.Close()

	resolver := ResolverForTimeSpan{DB: db.DB}
	timeSpans, err := resolver.TimeSpans(fake.User(5), test.ModelTimeP("2019-06-10T18:30:00+02:00"),
		test.ModelTimeP("2019-06-10T17:30:00+02:00"))
	require.Nil(t, timeSpans)
	require.EqualError(t, err, "fromInclusive must be before toInclusive")
}
