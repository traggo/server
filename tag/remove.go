package tag

import (
	"context"
	"fmt"

	"github.com/jinzhu/copier"
	"github.com/traggo/server/auth"
	"github.com/traggo/server/generated/gqlmodel"
	"github.com/traggo/server/model"
)

// RemoveTag removes a tag.
func (r *ResolverForTag) RemoveTag(ctx context.Context, key string) (*gqlmodel.TagDefinition, error) {
	tag := model.TagDefinition{Key: key}
	if r.DB.Where("user_id = ?", auth.GetUser(ctx).ID).Find(&tag).RecordNotFound() {
		return nil, fmt.Errorf("tag with key '%s' does not exist", tag.Key)
	}

	remove := r.DB.Delete(&tag)
	gqlTag := &gqlmodel.TagDefinition{}
	copier.Copy(gqlTag, &tag)
	return gqlTag, remove.Error
}
