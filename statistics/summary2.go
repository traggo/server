package statistics

import (
	"context"

	"github.com/traggo/server/generated/gqlmodel"
	"github.com/traggo/server/model"
	"github.com/traggo/server/setting"
	"github.com/traggo/server/time"
)

// Stats2 another version of the stats endpoint
func (r *ResolverForStatistics) Stats2(ctx context.Context, now model.Time, stats gqlmodel.InputStatsSelection) ([]*gqlmodel.RangedStatisticsEntries, error) {

	settings, err := setting.Get(ctx, r.DB)
	if err != nil {
		return nil, err
	}

	var ranges []*gqlmodel.Range

	staticRanges, err := time.ParseRange(now.OmitTimeZone(),
		time.RelativeRange{From: stats.Range.From, To: stats.Range.To},
		time.InternalInterval(stats.Interval),
		settings.FirstDayOfTheWeekTimeWeekday(),
		settings.LastDayOfTheWeekTimeWeekday())
	if err != nil {
		return nil, err
	}
	for _, r := range staticRanges {
		ranges = append(ranges, &gqlmodel.Range{Start: model.Time(r.From), End: model.Time(r.To)})
	}

	return r.Stats(ctx, ranges, stats.Tags, stats.ExcludeTags, stats.IncludeTags)
}
