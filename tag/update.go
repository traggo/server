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
func (r *ResolverForTag) UpdateTag(ctx context.Context, key string, newKey *string, color string) (*gqlmodel.TagDefinition, error) {
	tag := model.TagDefinition{}
	userID := auth.GetUser(ctx).ID
	if r.DB.Where(&model.TagDefinition{UserID: userID, Key: key}).Find(&tag).RecordNotFound() {
		return nil, fmt.Errorf("tag with key '%s' does not exist", key)
	}

	tx := r.DB.Begin()

	newValue := model.TagDefinition{
		Key:    strings.ToLower(key),
		Color:  color,
		UserID: userID,
	}

	if newKey != nil && *newKey != key {
		if strings.Contains(*newKey, " ") {
			tx.Rollback()
			return nil, fmt.Errorf("tag must not contain spaces")
		}
		*newKey = strings.ToLower(*newKey)
		newValue.Key = *newKey
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
		usedInEntries := []model.DashboardEntry{}

		// Do not read the next statements, not proud of it.
		if err := tx.Where("keys LIKE ?", "%"+key).
			Or("keys like ?", "%"+key+"%").
			Or("keys like ?", key+"%").
			Find(&usedInEntries).Error; err != nil {
			tx.Rollback()
			return nil, err
		}

		for _, entry := range usedInEntries {
			tags := strings.Split(entry.Keys, ",")
			for index, tagInEntry := range tags {
				if tagInEntry == key {
					tags[index] = *newKey
				}
			}
			entry.Keys = strings.Join(tags, ",")
			if err := tx.Save(&entry).Error; err != nil {
				tx.Rollback()
				return nil, err
			}
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
