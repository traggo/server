package entry

import (
	"context"
	"errors"
	"fmt"
	"strings"

	"github.com/traggo/server/time"

	"github.com/traggo/server/dashboard/convert"
	"github.com/traggo/server/dashboard/util"

	"github.com/traggo/server/auth"
	"github.com/traggo/server/generated/gqlmodel"
	"github.com/traggo/server/model"
)

// AddDashboardEntry adds a dashboard entry.
func (r *ResolverForEntry) AddDashboardEntry(ctx context.Context, dashboardID int, entryType gqlmodel.EntryType, title string, total bool, stats gqlmodel.InputStatsSelection, pos *gqlmodel.InputResponsiveDashboardEntryPos) (*gqlmodel.DashboardEntry, error) {
	userID := auth.GetUser(ctx).ID

	if _, err := util.FindDashboard(r.DB, userID, dashboardID); err != nil {
		return nil, err
	}

	entry := model.DashboardEntry{
		Keys:            strings.Join(stats.Tags, ","),
		Type:            convert.InternalEntryType(entryType),
		Title:           title,
		Total:           total,
		DashboardID:     dashboardID,
		Interval:        convert.InternalInterval(stats.Interval),
		MobilePosition:  convert.EmptyPos(),
		DesktopPosition: convert.EmptyPos(),
		RangeID:         -1,
	}

	if len(stats.Tags) == 0 {
		return nil, errors.New("at least one tag is required")
	}

	if err := convert.ApplyPos(&entry, pos); err != nil {
		return &gqlmodel.DashboardEntry{}, err
	}

	if stats.RangeID != nil {
		if _, err := util.FindDashboardRange(r.DB, *stats.RangeID); err != nil {
			return nil, err
		}
		entry.RangeID = *stats.RangeID
	} else if stats.Range != nil {
		entry.RangeID = model.NoRangeIDDefined
		if err := time.Validate(stats.Range.From); err != nil {
			return nil, fmt.Errorf("range from (%s) invalid: %s", stats.Range.From, err)
		}
		if err := time.Validate(stats.Range.To); err != nil {
			return nil, fmt.Errorf("range to (%s) invalid: %s", stats.Range.To, err)
		}
		entry.RangeFrom = stats.Range.From
		entry.RangeTo = stats.Range.To
	}

	if err := r.DB.Save(&entry).Error; err != nil {
		return nil, err
	}

	return convert.ToExternalEntry(entry)
}
