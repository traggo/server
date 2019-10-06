package dbrange

import "github.com/jinzhu/gorm"

// ResolverForRange resolves range specific things.
type ResolverForRange struct {
	DB *gorm.DB
}
