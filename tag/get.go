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
	userID := auth.GetUser(ctx).ID

	timeSpansIdsOfUser := r.DB.Model(new(model.TimeSpan)).
		Select("id").
		Where(&model.TimeSpan{UserID: userID}).
		SubQuery()
	usages := r.DB.Select("COUNT (*)").Where("time_span_tags.time_span_id in ?", timeSpansIdsOfUser).
		Where("tag_definitions.key = time_span_tags.key").
		Model(new(model.TimeSpanTag)).
		Group("time_span_tags.key").
		SubQuery()
	find := r.DB.Select("tag_definitions.*, ? as usages", usages).Where("user_id = ?", userID).Order("usages desc").Find(&tags)
	result := []*gqlmodel.TagDefinition{}
	copier.Copy(&result, &tags)
	return result, find.Error
}
