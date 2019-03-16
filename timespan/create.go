package timespan

import (
	"context"
	"fmt"

	"github.com/traggo/server/auth"
	"github.com/traggo/server/generated/gqlmodel"
	"github.com/traggo/server/model"
)

// CreateTimeSpan creates a time span
func (r *ResolverForTimeSpan) CreateTimeSpan(ctx context.Context, start model.Time, end *model.Time, tags []gqlmodel.InputTimeSpanTag) (*gqlmodel.TimeSpan, error) {
	timeSpan, err := timespanToInternal(auth.GetUser(ctx).ID, start, end, tags)
	if err != nil {
		return nil, err
	}

	existingTags := make(map[string]struct{})

	for _, tag := range timeSpan.Tags {
		if _, ok := existingTags[tag.Key]; ok {
			return nil, fmt.Errorf("tag '%s' is present multiple times", tag.Key)
		}

		if r.DB.Where("key = ?", tag.Key).Where("user_id = ?", auth.GetUser(ctx).ID).Find(new(model.TagDefinition)).RecordNotFound() {
			return nil, fmt.Errorf("tag '%s' does not exist", tag.Key)
		}

		existingTags[tag.Key] = struct{}{}
	}

	r.DB.Create(&timeSpan)

	external := timeSpanToExternal(timeSpan)
	return &external, nil
}
