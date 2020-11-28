package timespan

import (
	"testing"

	"github.com/rs/zerolog/log"
	"github.com/stretchr/testify/require"
	"github.com/traggo/server/generated/gqlmodel"
	"github.com/traggo/server/model"
	"github.com/traggo/server/test"
	"github.com/traggo/server/test/fake"
)

func Test_Update_withoutEnd(t *testing.T) {
	db := test.InMemoryDB(t)
	defer db.Close()
	db.User(3)
	note := "A note"
	db.Create(&model.TimeSpan{
		ID:            1,
		UserID:        3,
		StartUserTime: test.Time("2019-06-10T18:30:00Z"),
		StartUTC:      test.Time("2019-06-10T16:30:00Z"),
		OffsetUTC:     7200,
		Tags:          []model.TimeSpanTag{},
		Note:          note,
	})

	resolver := ResolverForTimeSpan{DB: db.DB}
	timeSpan, err := resolver.UpdateTimeSpan(fake.User(3), 1, test.ModelTime("2019-06-10T19:30:00+02:00"), nil, nil, nil, "")

	require.Nil(t, err)
	expected := &gqlmodel.TimeSpan{
		ID:       1,
		Start:    test.ModelTime("2019-06-10T19:30:00+02:00"),
		OldStart: test.ModelTimeP("2019-06-10T18:30:00+02:00"),
		Tags:     []*gqlmodel.TimeSpanTag{},
		Note:     "",
	}
	require.Equal(t, expected, timeSpan)
	assertTimeSpanCount(t, db, 1)
	assertTimeSpanExist(t, db, model.TimeSpan{
		ID:            1,
		UserID:        3,
		StartUserTime: test.Time("2019-06-10T19:30:00Z"),
		StartUTC:      test.Time("2019-06-10T17:30:00Z"),
		OffsetUTC:     7200,
		Tags:          []model.TimeSpanTag{},
		Note:          "",
	})
}

func Test_Update(t *testing.T) {
	db := test.InMemoryDB(t)
	defer db.Close()
	db.User(5)
	note := "A note"
	db.Create(&model.TimeSpan{
		ID:            3,
		UserID:        5,
		StartUserTime: test.Time("2019-06-10T18:30:00Z"),
		StartUTC:      test.Time("2019-06-10T16:30:00Z"),
		OffsetUTC:     7200,
		Tags:          []model.TimeSpanTag{},
		Note:          note,
	})

	note = "An upated note"
	resolver := ResolverForTimeSpan{DB: db.DB}
	timeSpan, err := resolver.UpdateTimeSpan(fake.User(5), 3, test.ModelTime("2019-06-10T18:30:00+02:00"),
		test.ModelTimeP("2019-06-10T20:30:00+02:00"), nil, nil, note)
	log.Debug().Msg("oops")
	require.Nil(t, err)
	expected := &gqlmodel.TimeSpan{
		ID:       3,
		Start:    test.ModelTime("2019-06-10T18:30:00+02:00"),
		OldStart: test.ModelTimeP("2019-06-10T18:30:00+02:00"),
		End:      test.ModelTimeP("2019-06-10T20:30:00+02:00"),
		Tags:     []*gqlmodel.TimeSpanTag{},
		Note:     note,
	}
	require.Equal(t, expected, timeSpan)
	assertTimeSpanCount(t, db, 1)
	assertTimeSpanExist(t, db, model.TimeSpan{
		ID:            3,
		UserID:        5,
		StartUserTime: test.Time("2019-06-10T18:30:00Z"),
		StartUTC:      test.Time("2019-06-10T16:30:00Z"),
		EndUserTime:   test.TimeP("2019-06-10T20:30:00Z"),
		EndUTC:        test.TimeP("2019-06-10T18:30:00Z"),
		OffsetUTC:     7200,
		Tags:          []model.TimeSpanTag{},
		Note:          note,
	})
}

func Test_Update_fail_endBeforeStart(t *testing.T) {
	db := test.InMemoryDB(t)
	defer db.Close()
	db.User(5)
	db.Create(&model.TimeSpan{
		ID:            3,
		UserID:        5,
		StartUserTime: test.Time("2019-06-10T18:30:00Z"),
		StartUTC:      test.Time("2019-06-10T16:30:00Z"),
		EndUserTime:   test.TimeP("2019-06-10T18:30:00Z"),
		EndUTC:        test.TimeP("2019-06-10T16:30:00Z"),
		OffsetUTC:     7200,
		Tags:          []model.TimeSpanTag{},
	})

	resolver := ResolverForTimeSpan{DB: db.DB}
	timeSpan, err := resolver.UpdateTimeSpan(fake.User(5), 3, test.ModelTime("2019-06-10T18:30:00+02:00"),
		test.ModelTimeP("2019-06-10T17:30:00+02:00"), nil, nil, "")
	require.Nil(t, timeSpan)
	require.EqualError(t, err, "start must be before end")
	assertTimeSpanCount(t, db, 1)
}

func Test_Update_fail_notExistingTag(t *testing.T) {
	db := test.InMemoryDB(t)
	defer db.Close()
	db.User(5)
	db.Create(&model.TimeSpan{
		ID:            3,
		UserID:        5,
		StartUserTime: test.Time("2019-06-10T18:30:00Z"),
		StartUTC:      test.Time("2019-06-10T16:30:00Z"),
		EndUserTime:   test.TimeP("2019-06-10T18:30:00Z"),
		EndUTC:        test.TimeP("2019-06-10T16:30:00Z"),
		OffsetUTC:     7200,
		Tags:          []model.TimeSpanTag{},
	})

	resolver := ResolverForTimeSpan{DB: db.DB}
	timeSpan, err := resolver.UpdateTimeSpan(fake.User(5), 3, test.ModelTime("2019-06-10T18:30:00+02:00"),
		test.ModelTimeP("2019-06-10T18:35:00+02:00"), []*gqlmodel.InputTimeSpanTag{{Key: "test"}}, nil, "nil")
	require.Nil(t, timeSpan)
	require.EqualError(t, err, "tag 'test' does not exist")
	assertTimeSpanCount(t, db, 1)
	assertTimeSpanExist(t, db, model.TimeSpan{
		ID:            3,
		UserID:        5,
		StartUserTime: test.Time("2019-06-10T18:30:00Z"),
		StartUTC:      test.Time("2019-06-10T16:30:00Z"),
		EndUserTime:   test.TimeP("2019-06-10T18:30:00Z"),
		EndUTC:        test.TimeP("2019-06-10T16:30:00Z"),
		OffsetUTC:     7200,
		Tags:          []model.TimeSpanTag{},
	})
}

func Test_Update_withTag(t *testing.T) {
	db := test.InMemoryDB(t)
	defer db.Close()
	db.User(5)
	db.Create(&model.TimeSpan{
		ID:            3,
		UserID:        5,
		StartUserTime: test.Time("2019-06-10T18:30:00Z"),
		StartUTC:      test.Time("2019-06-10T16:30:00Z"),
		EndUserTime:   test.TimeP("2019-06-10T19:30:00Z"),
		EndUTC:        test.TimeP("2019-06-10T17:30:00Z"),
		OffsetUTC:     7200,
		Tags:          []model.TimeSpanTag{},
	})
	db.Create(&model.TagDefinition{Key: "test", UserID: 5})

	resolver := ResolverForTimeSpan{DB: db.DB}
	timeSpan, err := resolver.UpdateTimeSpan(fake.User(5), 3, test.ModelTime("2019-06-10T18:30:00+02:00"),
		test.ModelTimeP("2019-06-10T19:30:00+02:00"), []*gqlmodel.InputTimeSpanTag{{Key: "test"}}, nil, "")
	require.NotNil(t, timeSpan)
	require.NoError(t, err)
	assertTimeSpanCount(t, db, 1)
	assertTimeSpanExist(t, db, model.TimeSpan{
		ID:            3,
		UserID:        5,
		StartUserTime: test.Time("2019-06-10T18:30:00Z"),
		StartUTC:      test.Time("2019-06-10T16:30:00Z"),
		EndUserTime:   test.TimeP("2019-06-10T19:30:00Z"),
		EndUTC:        test.TimeP("2019-06-10T17:30:00Z"),
		OffsetUTC:     7200,
		Tags: []model.TimeSpanTag{
			{Key: "test", TimeSpanID: 3},
		},
	})
}

func Test_Update_fail_tagAddedMultipleTimes(t *testing.T) {
	db := test.InMemoryDB(t)
	defer db.Close()
	db.User(5)
	db.Create(&model.TimeSpan{
		ID:            3,
		UserID:        5,
		StartUserTime: test.Time("2019-06-10T18:30:00Z"),
		StartUTC:      test.Time("2019-06-10T16:30:00Z"),
		EndUserTime:   test.TimeP("2019-06-10T18:30:00Z"),
		EndUTC:        test.TimeP("2019-06-10T16:30:00Z"),
		OffsetUTC:     7200,
		Tags:          []model.TimeSpanTag{},
	})
	db.Create(&model.TagDefinition{Key: "test", UserID: 5})

	resolver := ResolverForTimeSpan{DB: db.DB}
	timeSpan, err := resolver.UpdateTimeSpan(fake.User(5), 3, test.ModelTime("2019-06-10T18:30:00+02:00"),
		test.ModelTimeP("2019-06-10T18:35:00+02:00"), []*gqlmodel.InputTimeSpanTag{{Key: "test"}, {Key: "test"}}, nil, "")
	require.Nil(t, timeSpan)
	require.EqualError(t, err, "tag 'test' is present multiple times")
	assertTimeSpanCount(t, db, 1)
}

func Test_Update_fail_notExisting(t *testing.T) {
	db := test.InMemoryDB(t)
	defer db.Close()
	db.User(5)
	db.Create(&model.TagDefinition{Key: "test", UserID: 5})

	resolver := ResolverForTimeSpan{DB: db.DB}
	timeSpan, err := resolver.UpdateTimeSpan(fake.User(5), 3, test.ModelTime("2019-06-10T18:30:00+02:00"),
		test.ModelTimeP("2019-06-10T18:35:00+02:00"), []*gqlmodel.InputTimeSpanTag{{Key: "test"}, {Key: "test"}}, nil, "")
	require.Nil(t, timeSpan)
	require.EqualError(t, err, "time span with id 3 does not exist")
	assertTimeSpanCount(t, db, 0)
}

func Test_Update_fail_noPermission(t *testing.T) {
	db := test.InMemoryDB(t)
	db.User(3)
	db.User(2)
	db.User(5)
	db.Create(&model.TimeSpan{
		ID:            3,
		UserID:        3,
		StartUserTime: test.Time("2019-06-10T18:30:00Z"),
		StartUTC:      test.Time("2019-06-10T16:30:00Z"),
		EndUserTime:   test.TimeP("2019-06-10T18:30:00Z"),
		EndUTC:        test.TimeP("2019-06-10T16:30:00Z"),
		OffsetUTC:     7200,
		Tags:          []model.TimeSpanTag{},
	})
	defer db.Close()
	db.Create(&model.TagDefinition{Key: "test", UserID: 5})

	resolver := ResolverForTimeSpan{DB: db.DB}
	timeSpan, err := resolver.UpdateTimeSpan(fake.User(2), 3, test.ModelTime("2019-06-10T18:30:00+02:00"),
		test.ModelTimeP("2019-06-10T18:35:00+02:00"), []*gqlmodel.InputTimeSpanTag{{Key: "test"}, {Key: "test"}}, nil, "")
	require.Nil(t, timeSpan)
	require.EqualError(t, err, "time span with id 3 does not exist")
	assertTimeSpanCount(t, db, 1)
}
