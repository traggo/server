package timespan

import (
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/traggo/server/model"
	"github.com/traggo/server/test"
)

func assertTimeSpanExist(t *testing.T, db *test.Database, expected model.TimeSpan) {
	found := new(model.TimeSpan)
	find := db.Preload("Tags").Where("user_id = ?", expected.UserID).Where("id = ?", expected.ID).Find(found)
	require.Nil(t, find.Error)
	require.NotNil(t, found)
	require.Equal(t, expected, *found)
}

func assertTimeSpanCount(t *testing.T, db *test.Database, expected int) {
	count := new(int)
	db.Model(new(model.TimeSpan)).Count(count)
	require.Equal(t, expected, *count)
}
