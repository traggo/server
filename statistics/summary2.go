package statistics

import (
	"context"
	"github.com/traggo/server/time"

	"github.com/traggo/server/generated/gqlmodel"
	"github.com/traggo/server/model"
)

// Stats2 another version of the stats endpoint
func (r *ResolverForStatistics) Stats2(ctx context.Context, now model.Time, stats gqlmodel.InputStatsSelection) ([]*gqlmodel.RangedStatisticsEntries, error) {

	var ranges []*gqlmodel.Range

	staticRanges, err := time.ParseRange(now.OmitTimeZone(), time.RelativeRange{
		From: stats.Range.From,
		To:   stats.Range.To,
	}, time.InternalInterval(stats.Interval))
	if err != nil {
		return nil, err
	}
	for _, r := range staticRanges {
		ranges = append(ranges, &gqlmodel.Range{Start: model.Time(r.From), End: model.Time(r.To)})
	}

	return r.Stats(ctx, ranges, stats.Tags, stats.ExcludeTags, stats.IncludeTags)
}

