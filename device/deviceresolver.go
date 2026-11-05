package device

import "github.com/jinzhu/gorm"

// ResolverForDevice resolves device specific things.
type ResolverForDevice struct {
	DB *gorm.DB
	// PasswordLoginDisabled rejects the username/password login.
	PasswordLoginDisabled bool
}
