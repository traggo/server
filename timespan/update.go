package timespan

import (
	"context"
	"fmt"
	"time"

	"github.com/traggo/server/auth"
	"github.com/traggo/server/generated/gqlmodel"
	"github.com/traggo/server/model"
)

// UpdateTimeSpan update a time span
func (r *ResolverForTimeSpan) UpdateTimeSpan(ctx context.Context, id int, start model.Time, end *model.Time, tags []*gqlmodel.InputTimeSpanTag, oldStart *model.Time, note string) (*gqlmodel.TimeSpan, error) {
	timeSpan, err := timespanToInternal(auth.GetUser(ctx).ID, start, end, tags, note)
	if err != nil {
		return nil, err
	}

	oldTimeSpan := model.TimeSpan{ID: id}
	if r.DB.Where("user_id = ?", auth.GetUser(ctx).ID).Find(&oldTimeSpan).RecordNotFound() {
		return nil, fmt.Errorf("time span with id %d does not exist", id)
	}

	if err := tagsExist(r.DB, auth.GetUser(ctx).ID, timeSpan.Tags); err != nil {
		return nil, err
	}

	timeSpan.ID = id

	r.DB.Where("time_span_id = ?", timeSpan.ID).Delete(new(model.TimeSpanTag))
	r.DB.Save(&timeSpan)

	external := timeSpanToExternal(timeSpan)
	if oldStart == nil {
		location := time.FixedZone("unknown", oldTimeSpan.OffsetUTC)
		old := model.Time(oldTimeSpan.StartUTC.In(location))
		oldStart = &old
	}
	external.OldStart = oldStart
	return external, nil
}
