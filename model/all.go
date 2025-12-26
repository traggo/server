package model

// All returns all schema instances.
func All() []interface{} {
	return []interface{}{
		new(TagDefinition),
		new(User),
		new(Device),
		new(TimeSpan),
		new(TimeSpanTag),
		new(UserSetting),
		new(Dashboard),
		new(DashboardEntry),
		new(DashboardTagFilter),
		new(DashboardRange),
	}
}
