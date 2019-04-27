package timespan

import (
	"testing"

	"github.com/traggo/server/generated/gqlmodel"

	"github.com/stretchr/testify/require"
	"github.com/traggo/server/model"
	"github.com/traggo/server/test"
	"github.com/traggo/server/test/fake"
)

func TestRemoveTimeSpan_succeeds_removesTimeSpan(t *testing.T) {
	db := test.InMemoryDB(t)
	defer db.Close()
	db.User(5)
	db.Create(&model.TimeSpan{
		ID:            1,
		StartUserTime: test.Time("2019-06-11T18:00:00Z"),
		StartUTC:      test.Time("2019-06-11T18:00:00Z"),
		EndUserTime:   nil,
		EndUTC:        nil,
		OffsetUTC:     0,
		UserID:        3,
	})

	resolver := ResolverForTimeSpan{DB: db.DB}
	actual, err := resolver.RemoveTimeSpan(fake.User(3), 1)
	require.NoError(t, err)
	expected := &gqlmodel.TimeSpan{
		ID:    1,
		Start: test.ModelTime("2019-06-11T18:00:00Z"),
	}
	require.Equal(t, expected, actual)
	assertTimeSpanCount(t, db, 0)
}

func TestRemoveTimeSpan_fails_notExistingTimeSpan(t *testing.T) {
	db := test.InMemoryDB(t)
	defer db.Close()
	db.User(3)

	resolver := ResolverForTimeSpan{DB: db.DB}
	_, err := resolver.RemoveTimeSpan(fake.User(3), 5)
	require.EqualError(t, err, "timespan with id 5 does not exist")
}

func TestRemoveTimeSpan_fails_noPermission(t *testing.T) {
	db := test.InMemoryDB(t)
	defer db.Close()
	db.User(3)
	db.User(5)
	db.Create(&model.TimeSpan{
		StartUserTime: test.Time("2019-06-11T18:00:00Z"),
		StartUTC:      test.Time("2019-06-11T18:00:00Z"),
		EndUserTime:   nil,
		EndUTC:        nil,
		OffsetUTC:     0,
		ID:            1,
		UserID:        3,
	})

	resolver := ResolverForTimeSpan{DB: db.DB}
	_, err := resolver.RemoveTimeSpan(fake.User(5), 1)
	require.EqualError(t, err, "timespan with id 1 does not exist")
}
