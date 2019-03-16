package timespan

import (
	"context"

	"github.com/traggo/server/auth"
	"github.com/traggo/server/generated/gqlmodel"
	"github.com/traggo/server/model"
)

// Timers returns all running timers for a user
func (r *ResolverForTimeSpan) Timers(ctx context.Context) ([]gqlmodel.TimeSpan, error) {
	user := auth.GetUser(ctx)

	var timeSpans []model.TimeSpan
	r.DB.Preload("Tags").Where("user_id = ?", user.ID).Where("end_user_time is null").Find(&timeSpans)

	var result []gqlmodel.TimeSpan
	for _, span := range timeSpans {
		result = append(result, timeSpanToExternal(span))
	}
	return result, nil
}
