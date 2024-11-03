package convert

import (
	"strings"

	"github.com/traggo/server/generated/gqlmodel"
	"github.com/traggo/server/model"
)

func toExternalEntries(entries []model.DashboardEntry) ([]*gqlmodel.DashboardEntry, error) {
	result := []*gqlmodel.DashboardEntry{}
	for _, entry := range entries {
		converted, err := ToExternalEntry(entry)
		if err != nil {
			return nil, err
		}
		result = append(result, converted)
	}
	return result, nil
}

// ToExternalEntry converts entries.
func ToExternalEntry(entry model.DashboardEntry) (*gqlmodel.DashboardEntry, error) {
	mobilePosition, err := toExternalPosition(entry.MobilePosition)
	if err != nil {
		return &gqlmodel.DashboardEntry{}, err
	}
	desktopPos, err := toExternalPosition(entry.DesktopPosition)
	if err != nil {
		return &gqlmodel.DashboardEntry{}, err
	}
	pos := gqlmodel.ResponsiveDashboardEntryPos{
		Mobile:  enhancePos(mobilePosition, "mobile"),
		Desktop: enhancePos(desktopPos, "desktop"),
	}
	dateRange := &gqlmodel.RelativeOrStaticRange{
		From: entry.RangeFrom,
		To:   entry.RangeTo,
	}
	stats := &gqlmodel.StatsSelection{
		Interval: ExternalInterval(entry.Interval),
		Tags:     strings.Split(entry.Keys, ","),
		Range:    dateRange,
	}
	if entry.RangeID != model.NoRangeIDDefined {
		stats.RangeID = &entry.RangeID
		stats.Range = nil
	}
	return &gqlmodel.DashboardEntry{
		ID:             entry.ID,
		Title:          entry.Title,
		Total:          entry.Total,
		Pos:            &pos,
		StatsSelection: stats,
		EntryType:      ExternalEntryType(entry.Type),
	}, nil
}
