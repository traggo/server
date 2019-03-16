package timespan

import (
	"context"
	"fmt"

	"github.com/traggo/server/auth"
	"github.com/traggo/server/generated/gqlmodel"
	"github.com/traggo/server/model"
)

// StopTimeSpan sets an end date to an existing time span.
func (r *ResolverForTimeSpan) StopTimeSpan(ctx context.Context, id int, end model.Time) (*gqlmodel.TimeSpan, error) {
	old := &model.TimeSpan{ID: id}

	if r.DB.Preload("Tags").Where("user_id = ?", auth.GetUser(ctx).ID).Find(old).RecordNotFound() {
		return nil, fmt.Errorf("time span with id %d does not exist", id)
	}

	if old.EndUTC != nil {
		return nil, fmt.Errorf("timespan with id %d has already an end date", id)
	}

	utc := end.UTC()
	old.EndUTC = &utc
	userTime := end.OmitTimeZone()
	old.EndUserTime = &userTime
	r.DB.Save(old)

	external := timeSpanToExternal(*old)
	return &external, nil
}
