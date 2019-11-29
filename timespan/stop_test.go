package timespan

import (
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/traggo/server/generated/gqlmodel"
	"github.com/traggo/server/model"
	"github.com/traggo/server/test"
	"github.com/traggo/server/test/fake"
)

func Test_Stop_fail_notExisting(t *testing.T) {
	db := test.InMemoryDB(t)
	defer db.Close()
	db.User(5)
	db.Create(&model.TagDefinition{Key: "test", UserID: 5})

	resolver := ResolverForTimeSpan{DB: db.DB}
	timeSpan, err := resolver.StopTimeSpan(fake.User(5), 3, test.ModelTime("2019-06-10T18:35:00+02:00"))
	require.Nil(t, timeSpan)
	require.EqualError(t, err, "time span with id 3 does not exist")
	assertTimeSpanCount(t, db, 0)
}

func Test_Stop_fail_noPermission(t *testing.T) {
	db := test.InMemoryDB(t)
	db.User(3)
	db.User(2)
	db.User(5)
	db.Create(&model.TimeSpan{
		ID:            3,
		UserID:        3,
		StartUserTime: test.Time("2019-06-10T18:30:00Z"),
		StartUTC:      test.Time("2019-06-10T16:30:00Z"),
		OffsetUTC:     7200,
		Tags:          []model.TimeSpanTag{},
	})
	defer db.Close()
	db.Create(&model.TagDefinition{Key: "test", UserID: 5})

	resolver := ResolverForTimeSpan{DB: db.DB}
	timeSpan, err := resolver.StopTimeSpan(fake.User(2), 3, test.ModelTime("2019-06-10T18:35:00+02:00"))
	require.Nil(t, timeSpan)
	require.EqualError(t, err, "time span with id 3 does not exist")
	assertTimeSpanCount(t, db, 1)
}

func Test_Stop_fail_alreadyFinished(t *testing.T) {
	db := test.InMemoryDB(t)
	db.User(2)
	db.User(5)
	db.Create(&model.TimeSpan{
		ID:            3,
		UserID:        2,
		StartUserTime: test.Time("2019-06-10T18:30:00Z"),
		StartUTC:      test.Time("2019-06-10T16:30:00Z"),
		EndUserTime:   test.TimeP("2019-06-10T19:30:00Z"),
		EndUTC:        test.TimeP("2019-06-10T17:30:00Z"),
		OffsetUTC:     7200,
		Tags:          []model.TimeSpanTag{},
	})
	defer db.Close()
	db.Create(&model.TagDefinition{Key: "test", UserID: 5})

	resolver := ResolverForTimeSpan{DB: db.DB}
	timeSpan, err := resolver.StopTimeSpan(fake.User(2), 3, test.ModelTime("2019-06-10T18:35:00+02:00"))
	require.Nil(t, timeSpan)
	require.EqualError(t, err, "timespan with id 3 has already an end date")
	assertTimeSpanCount(t, db, 1)
}

func Test_Stop(t *testing.T) {
	db := test.InMemoryDB(t)
	db.User(2)
	db.User(5)
	db.Create(&model.TimeSpan{
		ID:            3,
		UserID:        2,
		StartUserTime: test.Time("2019-06-10T18:30:00Z"),
		StartUTC:      test.Time("2019-06-10T16:30:00Z"),
		OffsetUTC:     7200,
		Tags:          []model.TimeSpanTag{},
	})
	defer db.Close()
	db.Create(&model.TagDefinition{Key: "test", UserID: 5})

	resolver := ResolverForTimeSpan{DB: db.DB}
	timeSpan, err := resolver.StopTimeSpan(fake.User(2), 3, test.ModelTime("2019-06-10T18:35:00+02:00"))
	require.NoError(t, err)

	expected := &gqlmodel.TimeSpan{
		ID:    3,
		Start: test.ModelTime("2019-06-10T18:30:00+02:00"),
		End:   test.ModelTimeP("2019-06-10T18:35:00+02:00"),
	}
	require.Equal(t, expected, timeSpan)

	assertTimeSpanCount(t, db, 1)
	assertTimeSpanExist(t, db, model.TimeSpan{
		ID:            3,
		UserID:        2,
		StartUserTime: test.Time("2019-06-10T18:30:00Z"),
		StartUTC:      test.Time("2019-06-10T16:30:00Z"),
		EndUserTime:   test.TimeP("2019-06-10T18:35:00Z"),
		EndUTC:        test.TimeP("2019-06-10T16:35:00Z"),
		OffsetUTC:     7200,
		Tags:          []model.TimeSpanTag{},
	})
}
