package timespan

import (
	"fmt"
	"time"

	"github.com/traggo/server/generated/gqlmodel"
	"github.com/traggo/server/model"
)

func timespanToInternal(userID int, start model.Time, end *model.Time, tags []*gqlmodel.InputTimeSpanTag, note string) (model.TimeSpan, error) {
	_, offset := start.Time().Zone()
	span := model.TimeSpan{
		StartUserTime: start.OmitTimeZone(),
		StartUTC:      start.UTC(),
		UserID:        userID,
		Tags:          tagsToInternal(tags),
		OffsetUTC:     offset,
		Note:          note,
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

func timeSpanToExternal(span model.TimeSpan) *gqlmodel.TimeSpan {
	location := time.FixedZone("unknown", span.OffsetUTC)

	result := gqlmodel.TimeSpan{
		Start: model.Time(span.StartUTC.In(location)),
		End:   nil,
		ID:    span.ID,
		Tags:  tagsToExternal(span.Tags),
		Note:  span.Note,
	}
	if span.EndUTC != nil && !span.EndUTC.IsZero() {
		end := *span.EndUTC
		endModel := model.Time(end.In(location))
		result.End = &endModel
	}

	return &result
}

func tagsToExternal(tags []model.TimeSpanTag) []*gqlmodel.TimeSpanTag {
	result := []*gqlmodel.TimeSpanTag{}
	for _, tag := range tags {
		result = append(result, &gqlmodel.TimeSpanTag{
			Key:   tag.Key,
			Value: tag.StringValue,
		})
	}
	return result
}

func tagsToInternal(gqls []*gqlmodel.InputTimeSpanTag) []model.TimeSpanTag {
	result := make([]model.TimeSpanTag, 0)
	for _, tag := range gqls {
		result = append(result, tagToInternal(*tag))
	}
	return result
}

func tagToInternal(gqls gqlmodel.InputTimeSpanTag) model.TimeSpanTag {
	return model.TimeSpanTag{
		Key:         gqls.Key,
		StringValue: gqls.Value,
	}
}

func tagsToInputTag(tags []model.TimeSpanTag) []*gqlmodel.InputTimeSpanTag {
	result := make([]*gqlmodel.InputTimeSpanTag, 0)
	for _, tag := range tags {
		result = append(result, &gqlmodel.InputTimeSpanTag{
			Key:   tag.Key,
			Value: tag.StringValue,
		})
	}
	return result
}
