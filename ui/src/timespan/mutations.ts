import {RemoveTimeSpan, RemoveTimeSpanVariables} from '../gql/__generated__/RemoveTimeSpan';
import * as gqlTimeSpan from '../gql/timeSpan';
import {TimeSpans} from '../gql/__generated__/TimeSpans';
import {Trackers} from '../gql/__generated__/Trackers';
import {MutationHookOptions} from '@apollo/react-hooks/lib/types';

export const removeTimeSpanOptions: MutationHookOptions<RemoveTimeSpan, RemoveTimeSpanVariables> = {
    update: (cache, {data}) => {
        let oldData: TimeSpans | null = null;
        try {
            oldData = cache.readQuery<TimeSpans>({query: gqlTimeSpan.TimeSpans});
        } catch (e) {}

        const oldTrackers = cache.readQuery<Trackers>({query: gqlTimeSpan.Trackers});
        if (!data || !data.removeTimeSpan) {
            return;
        }
        const removedId = data.removeTimeSpan.id;
        if (oldTrackers) {
            cache.writeQuery<Trackers>({
                query: gqlTimeSpan.Trackers,
                data: {
                    timers: (oldTrackers.timers || []).filter((tracker) => tracker.id !== removedId),
                },
            });
        }
        if (oldData) {
            cache.writeQuery<TimeSpans>({
                query: gqlTimeSpan.TimeSpans,
                data: {
                    timeSpans: {
                        __typename: 'PagedTimeSpans',
                        timeSpans: oldData.timeSpans.timeSpans.filter((ts) => ts.id !== removedId),
                        cursor: oldData.timeSpans.cursor,
                    },
                },
            });
        }
    },
};
