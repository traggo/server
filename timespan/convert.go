package timespan

import (
	"fmt"
	"time"

	"github.com/traggo/server/generated/gqlmodel"
	"github.com/traggo/server/model"
)

func timespanToInternal(userID int, start model.Time, end *model.Time, tags []gqlmodel.InputTimeSpanTag) (model.TimeSpan, error) {
	_, offset := start.Time().Zone()
	span := model.TimeSpan{
		StartUserTime: start.OmitTimeZone(),
		StartUTC:      start.UTC(),
		UserID:        userID,
		Tags:          tagsToInternal(tags),
		OffsetUTC:     offset,
	}

	if end != nil {
		if start.Time().After(end.Time()) {
			return span, fmt.Errorf("start must be before end")
		}
		endUser := end.OmitTimeZone()
		span.EndUserTime = &endUser
		endUTC := end.UTC()
		span.EndUTC = &endUTC
	}

	return span, nil
}

func timeSpanToExternal(span model.TimeSpan) gqlmodel.TimeSpan {
	location := time.FixedZone("unknown", span.OffsetUTC)

	result := gqlmodel.TimeSpan{
		Start: model.Time(span.StartUTC.In(location)),
		End:   nil,
		ID:    span.ID,
		Tags:  tagsToExternal(span.Tags),
	}
	if span.EndUTC != nil && !span.EndUTC.IsZero() {
		end := *span.EndUTC
		endModel := model.Time(end.In(location))
		result.End = &endModel
	}

	return result
}

func tagsToExternal(tags []model.TimeSpanTag) []gqlmodel.TimeSpanTag {
	var result []gqlmodel.TimeSpanTag
	for _, tag := range tags {
		result = append(result, gqlmodel.TimeSpanTag{
			Key:         tag.Key,
			StringValue: tag.StringValue,
		})
	}
	return result
}

func tagsToInternal(gqls []gqlmodel.InputTimeSpanTag) []model.TimeSpanTag {
	result := make([]model.TimeSpanTag, 0)
	for _, tag := range gqls {
		result = append(result, model.TimeSpanTag{
			Key:         tag.Key,
			StringValue: tag.StringValue,
		})
	}
	return result
}

func tagsToInputTag(tags []model.TimeSpanTag) []gqlmodel.InputTimeSpanTag {
	result := make([]gqlmodel.InputTimeSpanTag, 0)
	for _, tag := range tags {
		result = append(result, gqlmodel.InputTimeSpanTag{
			Key:         tag.Key,
			StringValue: tag.StringValue,
		})
	}
	return result
}
