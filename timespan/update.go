package timespan

import (
	"context"
	"fmt"

	"github.com/traggo/server/auth"
	"github.com/traggo/server/generated/gqlmodel"
	"github.com/traggo/server/model"
)

// UpdateTimeSpan update a time span
func (r *ResolverForTimeSpan) UpdateTimeSpan(ctx context.Context, id int, start model.Time, end *model.Time, tags []gqlmodel.InputTimeSpanTag) (*gqlmodel.TimeSpan, error) {
	timeSpan, err := timespanToInternal(auth.GetUser(ctx).ID, start, end, tags)
	if err != nil {
		return nil, err
	}

	if r.DB.Where("user_id = ?", auth.GetUser(ctx).ID).Find(&model.TimeSpan{ID: id}).RecordNotFound() {
		return nil, fmt.Errorf("time span with id %d does not exist", id)
	}

	if err := tagsExist(r.DB, auth.GetUser(ctx).ID, timeSpan.Tags); err != nil {
		return nil, err
	}

	timeSpan.ID = id

	r.DB.Save(&timeSpan)

	external := timeSpanToExternal(timeSpan)
	return &external, nil
}
