package model

import (
	"database/sql/driver"
	"fmt"
)

// DashboardRange a named range of a dashboard.
type DashboardRange struct {
	ID          int `gorm:"primary_key;unique_index;AUTO_INCREMENT"`
	Name        string
	DashboardID int
	Editable    bool
	From        string
	To          string
}

// Dashboard a dashboard
type Dashboard struct {
	ID      int `gorm:"primary_key;unique_index;AUTO_INCREMENT"`
	UserID  int
	Name    string
	Entries []DashboardEntry
	Ranges  []DashboardRange
}

// DashboardEntry an entry which represents a diagram in a dashboard.
type DashboardEntry struct {
	ID          int `gorm:"primary_key;unique_index;AUTO_INCREMENT"`
	DashboardID int
	Title       string
	Type        DashboardType
	Keys        string
	Interval    Interval
	RangeID     int
	RangeFrom   string
	RangeTo     string

	MobilePosition  string
	DesktopPosition string
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
	TypeStackedBarChart DashboardType = "stackedbarchart"

	IntervalHourly  Interval = "hourly"
	IntervalDaily   Interval = "daily"
	IntervalWeekly  Interval = "weekly"
	IntervalMonthly Interval = "monthly"
	IntervalYearly  Interval = "yearly"
	IntervalSingle  Interval = "single"

	NoRangeIDDefined = -1
)
