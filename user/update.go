package user

import (
	"context"
	"fmt"

	"github.com/jinzhu/copier"
	"github.com/traggo/server/generated/gqlmodel"
	"github.com/traggo/server/schema"
)

// UpdateUser updates a user.
func (r *ResolverForUser) UpdateUser(ctx context.Context, id int, name string, pass *string, admin bool) (*gqlmodel.User, error) {
	user := new(schema.User)
	if r.DB.Find(user, id).RecordNotFound() {
		return nil, fmt.Errorf("user with id %d does not exist", id)
	}

	user.Name = name
	user.Admin = admin

	if pass != nil {
		user.Pass = createPassword(*pass, r.PassStrength)
	}

	update := r.DB.Save(user)
	gqlUser := &gqlmodel.User{}
	copier.Copy(gqlUser, user)
	return gqlUser, update.Error
}
