import {AddTimeSpanMutation, DashboardsQuery, StatsQuery, TagsQuery, TimeSpansQuery, TrackersQuery} from './__generated__';

// graphql-code-generator's typescript-operations plugin inlines nested selections as anonymous
// object types instead of apollo-codegen's old per-nesting-level named interfaces, so these are
// hand-derived aliases for the nested shapes the app references by name in several places.
type ArrayItem<T> = T extends Array<infer U> ? U : never;

export type Dashboard = ArrayItem<NonNullable<DashboardsQuery['dashboards']>>;
export type DashboardItem = ArrayItem<Dashboard['items']>;
export type DashboardItemRange = NonNullable<DashboardItem['statsSelection']['range']>;
export type DashboardRange = ArrayItem<Dashboard['ranges']>;

export type StatsRangeEntries = ArrayItem<NonNullable<StatsQuery['stats']>>;
export type StatEntry = ArrayItem<NonNullable<StatsRangeEntries['entries']>>;

export type Tag = ArrayItem<NonNullable<TagsQuery['tags']>>;

export type TimeSpanItem = ArrayItem<TimeSpansQuery['timeSpans']['timeSpans']>;

export type Tracker = ArrayItem<NonNullable<TrackersQuery['timers']>>;

export type CreatedTimeSpan = NonNullable<AddTimeSpanMutation['createTimeSpan']>;
