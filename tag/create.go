package tag

import (
	"context"
	"fmt"
	"strings"

	"github.com/jinzhu/copier"
	"github.com/traggo/server/auth"
	"github.com/traggo/server/generated/gqlmodel"
	"github.com/traggo/server/model"
)

// CreateTag creates a tag.
func (r *ResolverForTag) CreateTag(ctx context.Context, key string, color string) (*gqlmodel.TagDefinition, error) {
	if strings.TrimSpace(key) == "" {
		return nil, fmt.Errorf("tag must not be empty")
	}
	if strings.Contains(key, " ") {
		return nil, fmt.Errorf("tag must not contain spaces")
	}

	userID := auth.GetUser(ctx).ID
	definition := &model.TagDefinition{
		Key:    strings.ToLower(key),
		Color:  color,
		UserID: userID,
	}

	if !r.DB.Where("user_id = ?", userID).Where("key = ?", strings.ToLower(key)).Find(new(model.TagDefinition)).RecordNotFound() {
		return nil, fmt.Errorf("tag with key '%s' does already exist", definition.Key)
	}

	create := r.DB.Create(&definition)
	gqlTag := &gqlmodel.TagDefinition{}
	copier.Copy(gqlTag, definition)
	return gqlTag, create.Error
}
