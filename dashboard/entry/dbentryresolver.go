package entry

import "github.com/jinzhu/gorm"

// ResolverForEntry resolves dashboard entry things.
type ResolverForEntry struct {
	DB *gorm.DB
}
