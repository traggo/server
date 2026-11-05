import {DataProxy} from '@apollo/client';
import * as gqlTimeSpan from './timeSpan';
import moment from 'moment';
import {CreatedTimeSpan} from './types';
import {
    TimeSpansQuery,
    TimeSpansInRangeQuery,
    TimeSpansInRangeQueryVariables,
    TrackersQuery,
    StopTimerMutation,
} from './__generated__';

export const addTimeSpanToCache = (cache: DataProxy, ts: CreatedTimeSpan) => {
    let oldTimeSpans: TimeSpansQuery | null = null;
    try {
        oldTimeSpans = cache.readQuery<TimeSpansQuery>({query: gqlTimeSpan.TimeSpans});
    } catch {}
    if (!oldTimeSpans) {
        return;
    }
    cache.writeQuery<TimeSpansQuery>({
        query: gqlTimeSpan.TimeSpans,
        data: {
            timeSpans: {
                __typename: 'PagedTimeSpans',
                timeSpans: oldTimeSpans.timeSpans.timeSpans
                    .concat([ts])
                    .sort((a, b) => moment(b.start).unix() - moment(a.start).unix()),
                cursor: oldTimeSpans.timeSpans.cursor,
            },
        },
    });
};
export const addTimeSpanInRangeToCache = (cache: DataProxy, ts: CreatedTimeSpan, vars: TimeSpansInRangeQueryVariables) => {
    const oldTimeSpans = cache.readQuery<TimeSpansInRangeQuery>({query: gqlTimeSpan.TimeSpansInRange, variables: vars});
    if (!oldTimeSpans) {
        return;
    }
    cache.writeQuery<TimeSpansInRangeQuery>({
        query: gqlTimeSpan.TimeSpansInRange,
        variables: vars,
        data: {
            timeSpans: {
                __typename: 'PagedTimeSpans',
                timeSpans: oldTimeSpans.timeSpans.timeSpans
                    .concat([ts])
                    .sort((a, b) => moment(b.start).unix() - moment(a.start).unix()),
                cursor: oldTimeSpans.timeSpans.cursor,
            },
        },
    });
};
export const removeFromTrackersCache = (cache: DataProxy, data: StopTimerMutation) => {
    const oldTrackers = cache.readQuery<TrackersQuery>({query: gqlTimeSpan.Trackers});
    if (!oldTrackers || !data || !data.stopTimeSpan) {
        return;
    }
    cache.writeQuery<TrackersQuery>({
        query: gqlTimeSpan.Trackers,
        data: {
            timers: (oldTrackers.timers || []).filter((tracker) => tracker.id !== data.stopTimeSpan!.id),
        },
    });
};

export const removeFromTimeSpanInRangeCache = (cache: DataProxy, id: number, vars: TimeSpansInRangeQueryVariables) => {
    const old = cache.readQuery<TimeSpansInRangeQuery>({query: gqlTimeSpan.TimeSpansInRange, variables: vars});
    if (!old) {
        return;
    }
    cache.writeQuery<TimeSpansInRangeQuery>({
        query: gqlTimeSpan.TimeSpansInRange,
        variables: vars,
        data: {
            timeSpans: {
                __typename: 'PagedTimeSpans',
                timeSpans: old.timeSpans.timeSpans.filter((ts) => ts.id !== id),
                cursor: old.timeSpans.cursor,
            },
        },
    });
};
