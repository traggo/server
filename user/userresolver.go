package user

import (
	"github.com/jinzhu/gorm"
	"github.com/traggo/server/user/password"
)

var createPassword = password.CreatePassword

// ResolverForUser resolves user specific things.
type ResolverForUser struct {
	DB           *gorm.DB
	PassStrength int
}
