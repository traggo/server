package tag

import "github.com/jinzhu/gorm"

// ResolverForTag resolves tag specific things.
type ResolverForTag struct {
	DB *gorm.DB
}
