package timespan

import (
	"context"
	"errors"

	"github.com/traggo/server/auth"
	"github.com/traggo/server/generated/gqlmodel"
	"github.com/traggo/server/model"
)

// TimeSpans returns all time spans for a user
func (r *ResolverForTimeSpan) TimeSpans(ctx context.Context, fromInclusive *model.Time, toInclusive *model.Time) ([]gqlmodel.TimeSpan, error) {
	user := auth.GetUser(ctx)

	call := r.DB.Preload("Tags").Where("user_id = ?", user.ID)
	if fromInclusive != nil {
		if toInclusive != nil {
			if fromInclusive.Time().After(toInclusive.Time()) {
				return nil, errors.New("fromInclusive must be before toInclusive")
			}

			call = call.Where("start_user_time <= ? AND (end_user_time >= ? OR end_user_time is null)", toInclusive.OmitTimeZone(), fromInclusive.OmitTimeZone())
		} else {
			call = call.Where("start_user_time >= ? OR (end_user_time >= ? OR end_user_time is null)", fromInclusive.OmitTimeZone(), fromInclusive.OmitTimeZone())
		}
	} else if toInclusive != nil {
		call = call.Where("end_user_time <= ? OR start_user_time <= ?", toInclusive.OmitTimeZone(), toInclusive.OmitTimeZone())
	}

	var timeSpans []model.TimeSpan
	call.Find(&timeSpans)

	var result []gqlmodel.TimeSpan
	for _, span := range timeSpans {
		result = append(result, timeSpanToExternal(span))
	}
	return result, nil
}
