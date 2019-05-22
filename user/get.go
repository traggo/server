package user

import (
	"context"

	"github.com/jinzhu/copier"
	"github.com/traggo/server/generated/gqlmodel"
	"github.com/traggo/server/model"
)

// Users returns all users.
func (r *ResolverForUser) Users(ctx context.Context) ([]*gqlmodel.User, error) {
	var users []model.User
	find := r.DB.Find(&users)
	var result []*gqlmodel.User
	copier.Copy(&result, &users)
	return result, find.Error
}
