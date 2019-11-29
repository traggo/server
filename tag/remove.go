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
	userID := auth.GetUser(ctx).ID
	if r.DB.Where(&model.TagDefinition{UserID: userID, Key: key}).Find(&tag).RecordNotFound() {
		return nil, fmt.Errorf("tag with key '%s' does not exist", key)
	}

	usedInEntries := []model.DashboardEntry{}
	// Do not read the next statements, not proud of it.
	if err := r.DB.Where("keys LIKE ?", "%"+key).
		Or("keys like ?", "%"+key+"%").
		Or("keys like ?", key+"%").
		Find(&usedInEntries).Error; err != nil {
		return nil, err
	}

	if len(usedInEntries) > 0 {
		dashboard := &model.Dashboard{ID: usedInEntries[0].DashboardID}
		r.DB.Find(dashboard)
		return nil, fmt.Errorf("tag '%s' is used in dashboard '%s' entry '%s', remove this reference before deleting the tag", key, dashboard.Name, usedInEntries[0].Title)
	}

	tx := r.DB.Begin()
	if err := tx.Where(model.TagDefinition{Key: key, UserID: userID}).
		Delete(new(model.TagDefinition)).Error; err != nil {
		tx.Rollback()
		return nil, err
	}

	timeSpansIdsOfUser := tx.Model(new(model.TimeSpan)).
		Select("id").
		Where(&model.TimeSpan{UserID: userID}).
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
