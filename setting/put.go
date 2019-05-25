package setting

import (
	"context"

	"github.com/traggo/server/auth"
	"github.com/traggo/server/model"
)

// SettingPut sets the value of a setting.
func (r *ResolverForSettings) SettingPut(ctx context.Context, namespace string, key string, value string) (string, error) {
	return value, r.DB.Save(&model.Setting{
		UserID:    auth.GetUser(ctx).ID,
		Namespace: namespace,
		Key:       key,
		Value:     value,
	}).Error
}
