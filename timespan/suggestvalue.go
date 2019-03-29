package timespan

import (
	"context"

	"github.com/traggo/server/auth"

	"github.com/traggo/server/model"
)

// SuggestTagValue suggests a tag value.
func (r *ResolverForTimeSpan) SuggestTagValue(ctx context.Context, key string, query string) ([]string, error) {
	var suggestions []model.TimeSpanTag
	find := r.DB.
		Select("DISTINCT time_span_tags.string_value").
		Joins("JOIN time_spans on time_spans.id = time_span_tags.time_span_id").
		Where("user_id = ?", auth.GetUser(ctx).ID).
		Where("key = ?", key).Where("string_value LIKE ?", query+"%").
		Find(&suggestions)
	var result []string
	for _, value := range suggestions {
		result = append(result, *value.StringValue)
	}

	return result, find.Error
}
