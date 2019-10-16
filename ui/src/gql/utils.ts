import {DataProxy} from 'apollo-cache';
import {TimeSpans} from './__generated__/TimeSpans';
import * as gqlTimeSpan from './timeSpan';
import moment from 'moment';
import {AddTimeSpan_createTimeSpan} from './__generated__/AddTimeSpan';

export const addTimeSpanToCache = (cache: DataProxy, ts: AddTimeSpan_createTimeSpan) => {
    const oldTimeSpans = cache.readQuery<TimeSpans>({query: gqlTimeSpan.TimeSpans});
    if (!oldTimeSpans) {
        return;
    }
    cache.writeQuery<TimeSpans>({
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
