package setting

import (
	"context"

	"github.com/traggo/server/auth"
	"github.com/traggo/server/generated/gqlmodel"
	"github.com/traggo/server/model"
)

// Settings gets all settings for a namespace.
func (r *ResolverForSettings) Settings(ctx context.Context, namespace string) ([]*gqlmodel.Setting, error) {
	var result []*gqlmodel.Setting
	find := r.DB.
		Model(&model.Setting{}).
		Select("key, value").
		Where(&model.Setting{UserID: auth.GetUser(ctx).ID, Namespace: namespace}).
		Scan(&result)
	return result, find.Error
}
