package dashboard

import (
	"github.com/jinzhu/gorm"
	"github.com/traggo/server/dashboard/dbrange"
	"github.com/traggo/server/dashboard/entry"
)

// ResolverForDashboard resolves dashboard specific things.
type ResolverForDashboard struct {
	DB *gorm.DB
	entry.ResolverForEntry
	dbrange.ResolverForRange
}

// NewResolverForDashboard creates a new resolver.
func NewResolverForDashboard(db *gorm.DB) ResolverForDashboard {
	return ResolverForDashboard{
		DB: db,
		ResolverForEntry: entry.ResolverForEntry{
			DB: db,
		},
		ResolverForRange: dbrange.ResolverForRange{
			DB: db,
		},
	}
}
