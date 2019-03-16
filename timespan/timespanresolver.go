package timespan

import "github.com/jinzhu/gorm"

// ResolverForTimeSpan resolves time span specific things.
type ResolverForTimeSpan struct {
	DB *gorm.DB
}
