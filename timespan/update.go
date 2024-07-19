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
	// Setting the seconds to 0 when a change comes in, to not have 'wrong' Calculations
	var timeSpan model.TimeSpan
	var err error
	incStart := start.Time().Local()
	fmt.Println(time.Time.Location(incStart))
	fixedStart := model.Time(time.Date(incStart.Year(), incStart.Month(), incStart.Day(), incStart.Hour(), incStart.Minute(), 0, 0, time.Local))
	if end != nil {
		incEnd := end.Time().Local()
		fixedEnd := model.Time(time.Date(incEnd.Year(), incEnd.Month(), incEnd.Day(), incEnd.Hour(), incEnd.Minute(), 0, 0, time.Local))
		timeSpan, err = timespanToInternal(auth.GetUser(ctx).ID, fixedStart, &fixedEnd, tags, note)
		if err != nil {
			return nil, err
		}
	} else {
		timeSpan, err = timespanToInternal(auth.GetUser(ctx).ID, fixedStart, end, tags, note) // end == nil
		if err != nil {
			return nil, err
		}
	}
	oldTimeSpan := model.TimeSpan{ID: id}
	if r.DB.Where("user_id = ?", auth.GetUser(ctx).ID).Find(&oldTimeSpan).RecordNotFound() {
		return nil, fmt.Errorf("time span with id %d does not exist", id)
	}

	if err = tagsExist(r.DB, auth.GetUser(ctx).ID, timeSpan.Tags); err != nil {
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
