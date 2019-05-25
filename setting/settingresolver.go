package setting

import "github.com/jinzhu/gorm"

// ResolverForSettings resolves setting specific things.
type ResolverForSettings struct {
	DB *gorm.DB
}
