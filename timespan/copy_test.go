package timespan

import (
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/traggo/server/generated/gqlmodel"
	"github.com/traggo/server/model"
	"github.com/traggo/server/test"
	"github.com/traggo/server/test/fake"
)

func TestCopy_noPermission(t *testing.T) {
	db := test.InMemoryDB(t)
	db.User(3)
	db.User(5)
	db.User(2)
	db.Create(&model.TimeSpan{
		ID:            3,
		UserID:        3,
		StartUserTime: test.Time("2019-06-10T18:30:00Z"),
		StartUTC:      test.Time("2019-06-10T16:30:00Z"),
		EndUserTime:   test.TimeP("2019-06-10T18:30:00Z"),
		EndUTC:        test.TimeP("2019-06-10T16:30:00Z"),
		Tags:          []model.TimeSpanTag{},
	})
	defer db.Close()
	db.Create(&model.TagDefinition{Key: "test", UserID: 5})

	resolver := ResolverForTimeSpan{DB: db.DB}
	timeSpan, err := resolver.CopyTimeSpan(fake.User(2), 3, test.ModelTime("2019-06-10T18:30:00+02:00"),
		test.ModelTimeP("2019-06-10T18:35:00+02:00"))
	require.Nil(t, timeSpan)
	require.EqualError(t, err, "time span with id 3 does not exist")
	assertTimeSpanCount(t, db, 1)
}

func TestCopy(t *testing.T) {
	db := test.InMemoryDB(t)
	db.User(3)
	db.Create(&model.TagDefinition{Key: "test", UserID: 3})
	db.Create(&model.TimeSpan{
		ID:            3,
		UserID:        3,
		StartUserTime: test.Time("2019-06-10T18:30:00Z"),
		StartUTC:      test.Time("2019-06-10T16:30:00Z"),
		EndUserTime:   test.TimeP("2019-06-10T18:30:00Z"),
		EndUTC:        test.TimeP("2019-06-10T16:30:00Z"),
		OffsetUTC:     7200,
		Tags: []model.TimeSpanTag{
			{Key: "test", TimeSpanID: 3},
		},
		Note: "A special note",
	})
	defer db.Close()

	resolver := ResolverForTimeSpan{DB: db.DB}
	timeSpan, err := resolver.CopyTimeSpan(fake.User(3), 3, test.ModelTime("2019-06-15T10:30:00+02:00"), nil)
	require.NoError(t, err)

	expected := &gqlmodel.TimeSpan{
		ID:    4,
		Start: test.ModelTime("2019-06-15T10:30:00+02:00"),
		End:   nil,
		Tags: []*gqlmodel.TimeSpanTag{
			{Key: "test"},
		},
		Note: "A special note",
	}

	require.Equal(t, expected, timeSpan)

	assertTimeSpanCount(t, db, 2)
	assertTimeSpanExist(t, db, model.TimeSpan{
		ID:            4,
		UserID:        3,
		StartUserTime: test.Time("2019-06-15T10:30:00Z"),
		StartUTC:      test.Time("2019-06-15T08:30:00Z"),
		EndUserTime:   nil,
		EndUTC:        nil,
		OffsetUTC:     7200,
		Tags: []model.TimeSpanTag{
			{Key: "test", TimeSpanID: 4},
		},
		Note: "A special note",
	})
}

func TestCopy_notFound(t *testing.T) {
	db := test.InMemoryDB(t)
	db.User(3)
	db.User(2)
	db.User(5)
	defer db.Close()
	db.Create(&model.TagDefinition{Key: "test", UserID: 5})

	resolver := ResolverForTimeSpan{DB: db.DB}
	timeSpan, err := resolver.CopyTimeSpan(fake.User(2), 3, test.ModelTime("2019-06-10T18:30:00+02:00"),
		test.ModelTimeP("2019-06-10T18:35:00+02:00"))
	require.Nil(t, timeSpan)
	require.EqualError(t, err, "time span with id 3 does not exist")
	assertTimeSpanCount(t, db, 0)
}
