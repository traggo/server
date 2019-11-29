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
	tag := model.TagDefinition{}
	userId := auth.GetUser(ctx).ID
	if r.DB.Where(&model.TagDefinition{UserID: userId, Key: key}).Find(&tag).RecordNotFound() {
		return nil, fmt.Errorf("tag with key '%s' does not exist", tag.Key)
	}
	tx := r.DB.Begin()
	if err := tx.Where(model.TagDefinition{Key: key, UserID: userId}).
		Delete(new(model.TagDefinition)).Error; err != nil {
		tx.Rollback()
		return nil, err
	}

	timeSpansIdsOfUser := tx.Model(new(model.TimeSpan)).
		Select("id").
		Where(&model.TimeSpan{UserID: userId}).
		SubQuery()

	if err := tx.
		Where("time_span_id in ?", timeSpansIdsOfUser).
		Where(&model.TimeSpanTag{Key: key}).
		Delete(new(model.TimeSpanTag)).Error; err != nil {
		tx.Rollback()
		return nil, err
	}

	remove := tx.Commit()
	gqlTag := &gqlmodel.TagDefinition{}
	copier.Copy(gqlTag, &tag)
	return gqlTag, remove.Error
}
