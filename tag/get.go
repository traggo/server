package tag

import (
	"context"

	"github.com/jinzhu/copier"
	"github.com/traggo/server/generated/gqlmodel"
	"github.com/traggo/server/schema"
)

// Tags returns all tags.
func (r *ResolverForTag) Tags(ctx context.Context) ([]gqlmodel.TagDefinition, error) {
	var tags []schema.TagDefinition
	find := r.DB.Find(&tags)
	var result []gqlmodel.TagDefinition
	copier.Copy(&result, &tags)
	return result, find.Error
}
