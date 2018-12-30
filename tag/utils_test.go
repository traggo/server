package tag

import (
	"testing"

	"github.com/jinzhu/gorm"
	"github.com/stretchr/testify/require"
	"github.com/traggo/server/schema"
)

func assertTagExist(t *testing.T, db *gorm.DB, expected schema.TagDefinition) {
	foundUser := new(schema.TagDefinition)
	find := db.Where("key = ?", expected.Key).Find(foundUser)
	require.Nil(t, find.Error)
	require.NotNil(t, foundUser)
	require.Equal(t, expected, *foundUser)
}

func assertTagCount(t *testing.T, db *gorm.DB, expected int) {
	count := new(int)
	db.Model(new(schema.TagDefinition)).Count(count)
	require.Equal(t, expected, *count)
}
