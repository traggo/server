package user

import (
	"context"
	"fmt"

	"github.com/jinzhu/copier"
	"github.com/traggo/server/generated/gqlmodel"
	"github.com/traggo/server/schema"
)

// CreateUser creates a user.
func (r *ResolverForUser) CreateUser(ctx context.Context, name string, pass string, admin bool) (*gqlmodel.User, error) {
	newUser := &schema.User{
		Name:  name,
		Pass:  createPassword(pass, r.PassStrength),
		Admin: admin,
	}

	if !r.DB.Where("name = ?", newUser.Name).Find(&schema.User{}).RecordNotFound() {
		return nil, fmt.Errorf("user with name '%s' does already exist", newUser.Name)
	}

	create := r.DB.Create(&newUser)
	gqlUser := &gqlmodel.User{}
	copier.Copy(gqlUser, newUser)
	return gqlUser, create.Error
}
