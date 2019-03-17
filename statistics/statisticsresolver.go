package statistics

import "github.com/jinzhu/gorm"

// ResolverForStatistics resolves statistic things.
type ResolverForStatistics struct {
	DB *gorm.DB
}
