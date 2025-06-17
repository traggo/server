package model

import (
	"database/sql/driver"
	"fmt"
)

// DashboardRange a named range of a dashboard.
type DashboardRange struct {
	ID          int `gorm:"primary_key;unique_index;AUTO_INCREMENT"`
	Name        string
	DashboardID int `gorm:"type:int REFERENCES dashboards(id) ON DELETE CASCADE"`
	Editable    bool
	From        string
	To          string
}

// Dashboard a dashboard
type Dashboard struct {
	ID      int `gorm:"primary_key;unique_index;AUTO_INCREMENT"`
	UserID  int `gorm:"type:int REFERENCES users(id) ON DELETE CASCADE"`
	Name    string
	Entries []DashboardEntry
	Ranges  []DashboardRange
}

// DashboardEntry an entry which represents a diagram in a dashboard.
type DashboardEntry struct {
	ID           int `gorm:"primary_key;unique_index;AUTO_INCREMENT"`
	DashboardID  int `gorm:"type:int REFERENCES dashboards(id) ON DELETE CASCADE"`
	Title        string
	Total        bool `gorm:"default:false"`
	Type         DashboardType
	Keys         string
	Interval     Interval
	RangeID      int
	RangeFrom    string
	RangeTo      string
	ExcludedTags []DashboardExcludedTag
	IncludedTags []DashboardIncludedTag

	MobilePosition  string
	DesktopPosition string
}

// DashboardExcludedTag a tag for filtering timespans
type DashboardExcludedTag struct {
	DashboardEntryID int `gorm:"type:int REFERENCES dashboard_entries(id) ON DELETE CASCADE"`
	Key              string
	StringValue      string
}

// DashboardIncludedTag a tag for filtering timespans
type DashboardIncludedTag struct {
	DashboardEntryID int `gorm:"type:int REFERENCES dashboard_entries(id) ON DELETE CASCADE"`
	Key              string
	StringValue      string
}

// DashboardType the dashboard type
type DashboardType string

// Value for db
func (t DashboardType) Value() (driver.Value, error) {
	return string(t), nil
}

// Scan for db
func (t *DashboardType) Scan(value interface{}) error {
	s, ok := value.([]byte)
	if !ok {
		return fmt.Errorf("expected string but was %#v", value)
	}
	*t = DashboardType(s)
	return nil
}

// Interval the interval in which the diagram data should be grouped.
type Interval string

// Value for db
func (t Interval) Value() (driver.Value, error) {
	return string(t), nil
}

// Scan for db
func (t *Interval) Scan(value interface{}) error {
	s, ok := value.([]byte)
	if !ok {
		return fmt.Errorf("expected string but was %#v", value)
	}
	*t = Interval(s)
	return nil
}

// No lint errors please.
const (
	TypePieChart        DashboardType = "piechart"
	TypeBarChart        DashboardType = "barchart"
	TypeLineChart       DashboardType = "linechart"
	HorizontalTable     DashboardType = "horizontaltable"
	VerticalTable       DashboardType = "verticaltable"
	TypeStackedBarChart DashboardType = "stackedbarchart"

	IntervalHourly  Interval = "hourly"
	IntervalDaily   Interval = "daily"
	IntervalWeekly  Interval = "weekly"
	IntervalMonthly Interval = "monthly"
	IntervalYearly  Interval = "yearly"
	IntervalSingle  Interval = "single"

	NoRangeIDDefined = -1
)
