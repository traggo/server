package tag_test

import (
	"testing"

	"github.com/jinzhu/gorm"
	"github.com/traggo/server/schema"
	"github.com/traggo/server/tag"
	"github.com/traggo/server/test"
)

func Test_getTags_withoutTags(t *testing.T) {
	query := `
query {
   tags {
      key,
      color,
      type
   }
}`
	expected := `
{
  "tags": []
}`
	test.NewGQL(t).Queries(tag.Queries).Exec(query).Succeeds(expected)
}

func Test_getTags_withTags(t *testing.T) {
	query := `
query {
   tags {
      key,
      color,
      type
   }
}`
	createTags := func(db *gorm.DB) {
		db.Create(&schema.TagDefinition{Key: "my tag", Color: "#fff", Type: schema.TypeSingleValue})
		db.Create(&schema.TagDefinition{Key: "my tag 2", Color: "#fff", Type: schema.TypeSingleValue})
	}

	expected := `
{
  "tags": [
    {
      "color": "#fff",
      "key": "my tag",
      "type": "singlevalue"
    },
    {
      "color": "#fff",
      "key": "my tag 2",
      "type": "singlevalue"
    }
  ]
}`
	test.NewGQL(t).Queries(tag.Queries).BeforeExec(createTags).Exec(query).Succeeds(expected)

}

func Test_suggestTag_withMatchingTags(t *testing.T) {
	query := `
query {
   suggestTag(query: "pr") {
      key,
      color,
      type
   }
}`
	createTags := func(db *gorm.DB) {
		db.Create(&schema.TagDefinition{Key: "project", Color: "#fff", Type: schema.TypeSingleValue})
		db.Create(&schema.TagDefinition{Key: "priority", Color: "#fff", Type: schema.TypeSingleValue})
		db.Create(&schema.TagDefinition{Key: "wood", Color: "#fff", Type: schema.TypeSingleValue})
	}

	expected := `
{
  "suggestTag": [
    {
      "color": "#fff",
      "key": "project",
      "type": "singlevalue"
    },
    {
      "color": "#fff",
      "key": "priority",
      "type": "singlevalue"
    }
  ]
}`
	test.NewGQL(t).Queries(tag.Queries).BeforeExec(createTags).Exec(query).Succeeds(expected)

}

func Test_suggestTag_withoutMatchingTags(t *testing.T) {
	query := `
query {
   suggestTag(query: "fire") {
      key,
      color,
      type
   }
}`
	createTags := func(db *gorm.DB) {
		db.Create(&schema.TagDefinition{Key: "project", Color: "#fff", Type: schema.TypeSingleValue})
		db.Create(&schema.TagDefinition{Key: "wood", Color: "#fff", Type: schema.TypeSingleValue})
	}

	expected := `
{
  "suggestTag": []
}`
	test.NewGQL(t).Queries(tag.Queries).BeforeExec(createTags).Exec(query).Succeeds(expected)

}
