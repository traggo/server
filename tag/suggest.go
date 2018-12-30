package tag

import (
	"context"

	"github.com/jinzhu/copier"
	"github.com/traggo/server/generated/gqlmodel"
	"github.com/traggo/server/schema"
)

// SuggestTag suggests a tag.
func (r *ResolverForTag) SuggestTag(ctx context.Context, query string) ([]gqlmodel.TagDefinition, error) {
	var suggestions []schema.TagDefinition
	find := r.DB.Where("Key LIKE ?", query+"%").Find(&suggestions)
	var result []gqlmodel.TagDefinition
	copier.Copy(&result, &suggestions)
	return result, find.Error
}
