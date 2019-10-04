package timespan

import (
	"context"
	"errors"

	"github.com/traggo/server/auth"
	"github.com/traggo/server/generated/gqlmodel"
	"github.com/traggo/server/model"
)

// TimeSpans returns all time spans for a user
func (r *ResolverForTimeSpan) TimeSpans(ctx context.Context, fromInclusive *model.Time, toInclusive *model.Time, cursor *gqlmodel.InputCursor) (*gqlmodel.PagedTimeSpans, error) {
	user := auth.GetUser(ctx)
	cursor = normalize(cursor)

	if cursor.StartID == nil {
		var s model.TimeSpan
		if err := r.DB.Model(new(model.TimeSpan)).Select("max(id) as id").Find(&s).Error; err != nil {
			return nil, err
		}
		cursor.StartID = &s.ID
	}

	call := r.DB.Preload("Tags").Where("user_id = ?", user.ID).Not("end_user_time is NULL").Order("start_user_time DESC").Limit(*cursor.PageSize)
	if cursor.Offset != nil && cursor.StartID != nil {
		call = call.Where("id <= ?", *cursor.StartID).Offset(*cursor.Offset)
	}
	if fromInclusive != nil {
		if toInclusive != nil {
			if fromInclusive.Time().After(toInclusive.Time()) {
				return nil, errors.New("fromInclusive must be before toInclusive")
			}

			call = call.Where("start_user_time <= ? AND end_user_time >= ?", toInclusive.OmitTimeZone(), fromInclusive.OmitTimeZone())
		} else {
			call = call.Where("start_user_time >= ? OR end_user_time >= ?", fromInclusive.OmitTimeZone(), fromInclusive.OmitTimeZone())
		}
	} else if toInclusive != nil {
		call = call.Where("end_user_time <= ? OR start_user_time <= ?", toInclusive.OmitTimeZone(), toInclusive.OmitTimeZone())
	}

	var timeSpans []model.TimeSpan
	call.Find(&timeSpans)

	var result []*gqlmodel.TimeSpan
	for _, span := range timeSpans {
		result = append(result, timeSpanToExternal(span))
	}
	return &gqlmodel.PagedTimeSpans{
		TimeSpans: result,
		Cursor: &gqlmodel.Cursor{
			HasMore:  len(timeSpans) != 0 && *cursor.Offset%*cursor.PageSize == 0,
			Offset:   *cursor.Offset + len(timeSpans),
			StartID:  *cursor.StartID,
			PageSize: *cursor.PageSize},
	}, nil
}

func normalize(cursor *gqlmodel.InputCursor) *gqlmodel.InputCursor {
	if cursor == nil {
		cursor = &gqlmodel.InputCursor{}
	}

	maxPageSize := 100
	if cursor.PageSize == nil || maxPageSize < *cursor.PageSize {
		cursor.PageSize = &maxPageSize
	}

	if cursor.Offset == nil {
		zero := 0
		cursor.Offset = &zero
	}

	return cursor
}
