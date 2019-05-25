package setting

import (
	"context"

	"github.com/traggo/server/auth"
	"github.com/traggo/server/model"
)

// SettingGet gets a single setting.
func (r *ResolverForSettings) SettingGet(ctx context.Context, namespace string, key string) (string, error) {
	setting := model.Setting{
		UserID:    auth.GetUser(ctx).ID,
		Namespace: namespace,
		Key:       key,
	}
	find := r.DB.Where(&setting).Find(&setting)

	if find.RecordNotFound() {
		return "", nil
	}

	return setting.Value, find.Error
}
