package tag

import (
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/traggo/server/model"
	"github.com/traggo/server/test"
)

func assertTagExist(t *testing.T, db *test.Database, expected model.TagDefinition) {
	foundUser := new(model.TagDefinition)
	find := db.Where("key = ?", expected.Key).Find(foundUser)
	require.Nil(t, find.Error)
	require.NotNil(t, foundUser)
	require.Equal(t, expected, *foundUser)
}

func assertTagCount(t *testing.T, db *test.Database, expected int) {
	count := new(int)
	db.Model(new(model.TagDefinition)).Count(count)
	require.Equal(t, expected, *count)
}
