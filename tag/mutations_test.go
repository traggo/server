package tag_test

import (
	"testing"

	"github.com/jinzhu/gorm"
	"github.com/traggo/server/schema"
	"github.com/traggo/server/tag"
	"github.com/traggo/server/test"
)

func Test_createTag_succeeds(t *testing.T) {
	query := `
mutation {
   createTag(key: "new tag", color: "#123456", type:novalue) {
      key,
      color,
      type
   }
}`
	expected := `
{
  "createTag": {
    "color": "#123456",
    "key": "new tag",
    "type": "novalue"
  }
}`
	test.NewGQL(t).Mutations(tag.Mutations).Exec(query).Succeeds(expected)
}

func Test_createTag_withExistingTag_fails(t *testing.T) {
	createTag := func(db *gorm.DB) {
		db.Create(&schema.TagDefinition{Key: "existing tag", Color: "#fff", Type: schema.TypeSingleValue})
	}
	query := `
mutation {
   createTag(key: "existing tag", color: "#123456", type:novalue) {
      key
   }
}`
	expectedError := "tag with key 'existing tag' does already exist"

	test.NewGQL(t).Mutations(tag.Mutations).BeforeExec(createTag).Exec(query).Errs(expectedError)
}

func Test_removeTag_withExistingTag_succeeds(t *testing.T) {
	createTag := func(db *gorm.DB) {
		db.Create(&schema.TagDefinition{Key: "existing tag", Color: "#fff", Type: schema.TypeSingleValue})
	}
	query := `
mutation {
   removeTag(key: "existing tag") {
      key,
      color,
      type
   }
}`
	expected := `
{
  "removeTag": {
    "color": "#fff",
    "key": "existing tag",
    "type": "singlevalue"
  }
}`

	test.NewGQL(t).Mutations(tag.Mutations).BeforeExec(createTag).Exec(query).Succeeds(expected)
}

func Test_removeTag_withoutExistingTag_fails(t *testing.T) {
	query := `
mutation {
   removeTag(key: "not existing") {
      key,
      color,
      type
   }
}`

	expectedError := "tag with key 'not existing' does not exist"

	test.NewGQL(t).Mutations(tag.Mutations).Exec(query).Errs(expectedError)
}
