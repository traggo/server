package entry

import (
	"fmt"

	"github.com/jinzhu/gorm"
	"github.com/traggo/server/generated/gqlmodel"
	"github.com/traggo/server/model"
)

func tagsDuplicates(src []*gqlmodel.InputTimeSpanTag, dst []*gqlmodel.InputTimeSpanTag) *gqlmodel.InputTimeSpanTag {
	existingTags := make(map[string]struct{})
	for _, tag := range src {
		existingTags[tag.Key+":"+tag.Value] = struct{}{}
	}

	for _, tag := range dst {
		if _, ok := existingTags[tag.Key+":"+tag.Value]; ok {
			return tag
		}

		existingTags[tag.Key] = struct{}{}
	}

	return nil
}

func tagsExist(db *gorm.DB, userID int, tags []*gqlmodel.InputTimeSpanTag) error {
	existingTags := make(map[string]struct{})

	for _, tag := range tags {
		if _, ok := existingTags[tag.Key]; ok {
			return fmt.Errorf("tag '%s' is present multiple times", tag.Key)
		}

		if db.Where("key = ?", tag.Key).Where("user_id = ?", userID).Find(new(model.TagDefinition)).RecordNotFound() {
			return fmt.Errorf("tag '%s' does not exist", tag.Key)
		}

		existingTags[tag.Key] = struct{}{}
	}
	return nil
}
