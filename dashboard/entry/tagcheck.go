package entry

import (
	"fmt"

	"github.com/traggo/server/model"
)

func tagsDuplicates(tags []model.DashboardTagFilter) error {
	existingTags := make(map[model.DashboardTagFilter]struct{})

	for _, tag := range tags {
		if _, ok := existingTags[tag]; ok {
			tagType := "exclude"
			if tag.Include {
				tagType = "include"
			}

			return fmt.Errorf("%s tags: tag '%s' is present multiple times", tagType, tag.Key+":"+tag.StringValue)
		} else {
			copyTag := tag
			copyTag.Include = !copyTag.Include

			if _, ok := existingTags[copyTag]; ok {
				return fmt.Errorf("tag '%s' is present in both exclude tags and include tags", tag.Key+":"+tag.StringValue)
			}
		}

		existingTags[tag] = struct{}{}
	}

	return nil
}
