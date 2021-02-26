package entry

import (
	"context"
	"fmt"
	"strings"

	"github.com/traggo/server/time"

	"github.com/traggo/server/model"

	"github.com/traggo/server/auth"
	"github.com/traggo/server/dashboard/convert"
	"github.com/traggo/server/dashboard/util"
	"github.com/traggo/server/generated/gqlmodel"
)

// UpdateDashboardEntry updates a dashboard entry.
func (r *ResolverForEntry) UpdateDashboardEntry(ctx context.Context, id int, entryType *gqlmodel.EntryType, title *string, total *bool, stats *gqlmodel.InputStatsSelection, pos *gqlmodel.InputResponsiveDashboardEntryPos) (*gqlmodel.DashboardEntry, error) {
	userID := auth.GetUser(ctx).ID

	entry, err := util.FindDashboardEntry(r.DB, id)
	if err != nil {
		return nil, err
	}

	if _, err := util.FindDashboard(r.DB, userID, entry.DashboardID); err != nil {
		return nil, err
	}

	if title != nil {
		entry.Title = *title
	}

	if total != nil {
		entry.Total = *total
	}

	if stats != nil {
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
		entry.Keys = strings.Join(stats.Tags, ",")
		entry.Interval = convert.InternalInterval(stats.Interval)
	}

	if entryType != nil {
		entry.Type = convert.InternalEntryType(*entryType)
	}

	if err := convert.ApplyPos(&entry, pos); err != nil {
		return &gqlmodel.DashboardEntry{}, err
	}

	r.DB.Save(entry)

	return convert.ToExternalEntry(entry)
}
