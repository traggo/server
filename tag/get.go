package tag

import (
	"context"

	"github.com/jinzhu/copier"
	"github.com/traggo/server/auth"
	"github.com/traggo/server/generated/gqlmodel"
	"github.com/traggo/server/model"
)

// Tags returns all tags.
func (r *ResolverForTag) Tags(ctx context.Context) ([]*gqlmodel.TagDefinition, error) {
	var tags []model.TagDefinition
	find := r.DB.Where("user_id = ?", auth.GetUser(ctx).ID).Find(&tags)
	var result []*gqlmodel.TagDefinition
	copier.Copy(&result, &tags)
	return result, find.Error
}
