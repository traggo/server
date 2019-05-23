package statistics

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/traggo/server/auth"
	"github.com/traggo/server/generated/gqlmodel"
	"github.com/traggo/server/model"
)

// Stats groups the time spans by tag key.
func (r *ResolverForStatistics) Stats(ctx context.Context, ranges []*gqlmodel.Range, tags []string, excludeTags []*gqlmodel.InputTimeSpanTag, requireTags []*gqlmodel.InputTimeSpanTag) ([]*gqlmodel.RangedStatisticsEntries, error) {
	var variables []interface{}
	var rangesStrs []string
	for _, r := range ranges {
		if r.Start.Time().After(r.End.Time()) {
			return nil, fmt.Errorf("range start must be before range end")
		}
		rangesStrs = append(rangesStrs, "(?, ?)")
		variables = append(variables, r.Start.OmitTimeZone(), r.End.OmitTimeZone())
	}
	queryRanges := strings.Join(rangesStrs, ", ")

	queryRequire, requireVars := build(requireTags, "1 = 1")
	variables = append(variables, requireVars...)

	queryExclude, excludeVars := build(excludeTags, "1 != 1")
	variables = append(variables, excludeVars...)

	query := fmt.Sprintf(`
WITH dates(query_start, query_end) AS (
    VALUES %s
)
SELECT query_start                       as query_start,
       query_end                         as query_end,
       key                               as key,
       string_value                      as string_value,
       sum(round((julianday(CASE WHEN end_user_time > query_end THEN query_end ELSE end_user_time END) -
                  julianday(CASE WHEN start_user_time < query_start THEN query_start ELSE start_user_time END))
                     * 24 * 60 * 60, 0)) as time_spend_in_seconds
FROM dates
         JOIN time_span_tags tst
         JOIN (
    SELECT *
    FROM time_spans tsx
    WHERE (EXISTS(SELECT tstx.key
                  FROM time_span_tags as tstx
                  WHERE (tstx.time_span_id = tsx.id)
                    AND (%s)))
      AND (NOT EXISTS(SELECT tstx.key
                      FROM time_span_tags as tstx
                      WHERE (tstx.time_span_id = id)
                        AND (%s)))
) as ts on ts.id = tst.time_span_id
WHERE (key in (?))
  AND (user_id = ?)
  AND (start_user_time <= query_end AND end_user_time >= query_start)
GROUP BY query_start,
         key,
         string_value;
`, queryRanges, queryRequire, queryExclude)
	variables = append(variables, tags, auth.GetUser(ctx).ID)

	var entries []statReturn
	r.DB.Raw(query, variables...).Scan(&entries)

	statisticsEntries, err := group(entries)
	if err != nil {
		return nil, err
	}

	return statisticsEntries, r.DB.Error
}

type statReturn struct {
	QueryStart         string
	QueryEnd           string
	Key                string
	StringValue        *string
	TimeSpendInSeconds int
}

func group(entries []statReturn) ([]*gqlmodel.RangedStatisticsEntries, error) {
	stats := map[string]*gqlmodel.RangedStatisticsEntries{}
	for _, entry := range entries {
		id := entry.QueryStart + "/" + entry.QueryEnd
		if _, ok := stats[id]; !ok {
			start, err := time.ParseInLocation("2006-01-02 15:04:05Z07:00", entry.QueryStart, time.UTC)
			if err != nil {
				return nil, err
			}
			end, err := time.ParseInLocation("2006-01-02 15:04:05Z07:00", entry.QueryEnd, time.UTC)
			if err != nil {
				return nil, err
			}
			stats[id] = &gqlmodel.RangedStatisticsEntries{Start: model.Time(start), End: model.Time(end)}
		}
		stats[id].Entries = append(stats[id].Entries, &gqlmodel.StatisticsEntry{
			Key:                entry.Key,
			StringValue:        entry.StringValue,
			TimeSpendInSeconds: entry.TimeSpendInSeconds,
		})
	}

	var result []*gqlmodel.RangedStatisticsEntries
	for _, value := range stats {
		result = append(result, value)
	}

	return result, nil
}

func build(tags []*gqlmodel.InputTimeSpanTag, noop string) (string, []interface{}) {
	if len(tags) == 0 {
		return noop, nil
	}

	var have []string
	var haveParams []interface{}
	for _, tag := range tags {
		have = append(have, "(tstx.key = ? AND tstx.string_value = ?)")
		haveParams = append(haveParams, tag.Key, tag.StringValue)
	}

	return strings.Join(have, " OR "), haveParams
}
