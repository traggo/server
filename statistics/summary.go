package statistics

import (
	"context"
	"errors"
	"fmt"
	"strings"

	"github.com/traggo/server/auth"
	"github.com/traggo/server/generated/gqlmodel"
	"github.com/traggo/server/model"
)

// TimeSpanSummary groups the time spans by tag key.
func (r *ResolverForStatistics) TimeSpanSummary(ctx context.Context, from model.Time, to model.Time, stat gqlmodel.StatInput) ([]gqlmodel.StatisticsEntry, error) {
	if from.Time().After(to.Time()) {
		return nil, fmt.Errorf("fromInclusive must be before toInclusive")
	}

	base := r.DB.Select("tstx.key").
		Table("time_span_tags as tstx").
		Where("tstx.time_span_id = tsx.id")

	filteredTimeSpans := r.DB.Select(`
		(CASE WHEN tsx.start_user_time < ? THEN ? ELSE tsx.start_user_time END) as start_user_time,
		(CASE WHEN tsx.end_user_time > ? THEN ? ELSE tsx.end_user_time END) as end_user_time,
		tsx.*`, from.OmitTimeZone(), from.OmitTimeZone(), to.OmitTimeZone(), to.OmitTimeZone()).
		Table("time_spans as tsx")

	if query, params, err := build(stat.MustHave); err == nil {
		filteredTimeSpans = filteredTimeSpans.Where("EXISTS ?", base.Where(query, params...).SubQuery())
	}
	if query, params, err := build(stat.MustNotHave); err == nil {
		filteredTimeSpans = filteredTimeSpans.Where("NOT EXISTS ?", base.Where(query, params...).SubQuery())
	}

	var entries []gqlmodel.StatisticsEntry
	r.DB.Select(`max(tst.key) as key,
			max(tst.string_value) as string_value, 
			sum(round((julianday(ts.end_user_time) - julianday(ts.start_user_time)) * 24 * 60 * 60, 0)) as time_spend_in_seconds`).
		Table("time_span_tags as tst").
		Joins("JOIN ? as ts on ts.id = tst.time_span_id", filteredTimeSpans.SubQuery()).
		Group("tst.key, tst.string_value").
		Where("tst.key = ?", stat.Key).
		Where("ts.user_id = ?", auth.GetUser(ctx).ID).
		Where("ts.start_user_time <= ? AND ts.end_user_time >= ?", to.OmitTimeZone(), from.OmitTimeZone()).
		Scan(&entries)

	return entries, nil
}

func build(tags []gqlmodel.InputTimeSpanTag) (string, []interface{}, error) {
	if len(tags) == 0 {
		return "", nil, errors.New("empty")
	}

	var have []string
	var haveParams []interface{}
	for _, tag := range tags {
		have = append(have, "(tstx.key = ? AND tstx.string_value = ?)")
		haveParams = append(haveParams, tag.Key, tag.StringValue)
	}

	return strings.Join(have, " OR "), haveParams, nil
}
