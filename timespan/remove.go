package timespan

import (
	"context"
	"fmt"

	"github.com/traggo/server/auth"
	"github.com/traggo/server/generated/gqlmodel"
	"github.com/traggo/server/model"
)

// RemoveTimeSpan removes a timespan.
func (r *ResolverForTimeSpan) RemoveTimeSpan(ctx context.Context, id int) (*gqlmodel.TimeSpan, error) {
	timeSpan := model.TimeSpan{ID: id}
	if r.DB.Preload("Tags").Where("user_id = ?", auth.GetUser(ctx).ID).Find(&timeSpan).RecordNotFound() {
		return nil, fmt.Errorf("timespan with id %d does not exist", timeSpan.ID)
	}

	remove := r.DB.Delete(&timeSpan)
	external := timeSpanToExternal(timeSpan)
	return &external, remove.Error
}
