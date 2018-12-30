package tag

import (
	"testing"

	"github.com/jinzhu/gorm"
	"github.com/stretchr/testify/require"
	"github.com/traggo/server/model"
)

func assertTagExist(t *testing.T, db *gorm.DB, expected model.TagDefinition) {
	foundUser := new(model.TagDefinition)
	find := db.Where("key = ?", expected.Key).Find(foundUser)
	require.Nil(t, find.Error)
	require.NotNil(t, foundUser)
	require.Equal(t, expected, *foundUser)
}

func assertTagCount(t *testing.T, db *gorm.DB, expected int) {
	count := new(int)
	db.Model(new(model.TagDefinition)).Count(count)
	require.Equal(t, expected, *count)
}
