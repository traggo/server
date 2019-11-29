package tag

import (
	"context"
	"fmt"
	"strings"

	"github.com/jinzhu/copier"
	"github.com/traggo/server/auth"
	"github.com/traggo/server/generated/gqlmodel"
	"github.com/traggo/server/model"
)

// UpdateTag updates a tag.
func (r *ResolverForTag) UpdateTag(ctx context.Context, key string, newKey *string, color string, typeArg gqlmodel.TagDefinitionType) (*gqlmodel.TagDefinition, error) {
	tag := model.TagDefinition{}
	userID := auth.GetUser(ctx).ID
	if r.DB.Where(&model.TagDefinition{UserID: userID, Key: key}).Find(&tag).RecordNotFound() {
		return nil, fmt.Errorf("tag with key '%s' does not exist", key)
	}

	tx := r.DB.Begin()

	newValue := model.TagDefinition{
		Key:    strings.ToLower(key),
		Color:  color,
		Type:   model.TagDefinitionType(typeArg),
		UserID: userID,
	}

	if newKey != nil {
		newValue.Key = strings.ToLower(*newKey)
		timeSpansIdsOfUser := tx.Model(new(model.TimeSpan)).
			Select("id").
			Where(&model.TimeSpan{UserID: userID}).
			SubQuery()

		if err := tx.
			Model(new(model.TimeSpanTag)).
			Where("time_span_id in ?", timeSpansIdsOfUser).
			Where(&model.TimeSpanTag{Key: key}).
			Updates(&model.TimeSpanTag{Key: *newKey}).Error; err != nil {
			tx.Rollback()
			return nil, err
		}
	}

	if err := tx.Model(new(model.TagDefinition)).Where(&model.TagDefinition{UserID: userID, Key: key}).Updates(&newValue).Error; err != nil {
		tx.Rollback()
		return nil, err
	}

	update := tx.Commit()

	gqlTag := &gqlmodel.TagDefinition{}
	copier.Copy(gqlTag, &newValue)
	return gqlTag, update.Error
}
