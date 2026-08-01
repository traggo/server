import * as React from 'react';
import {useQuery} from '@apollo/client';
import * as gqlTimeSpan from '../gql/timeSpan';
import * as gqlTag from '../gql/tags';
import {TimeSpan, TimeSpanProps} from './TimeSpan';
import {TagsQuery, TimeSpansQuery, TimeSpansQueryVariables} from '../gql/__generated__';
import {useInterval} from '../utils/hooks';
import moment from 'moment';
import {Typography} from '@mui/material';
import {GroupedTimeSpanProps, toGroupedTimeSpanProps} from './timespanutils';
import {TagSelectorEntry} from '../tag/tagSelectorEntry';
import {isSameDate} from '../utils/time';

interface DoneTrackersProps {
    addTagsToTracker?: (entries: TagSelectorEntry[]) => void;
}

export const DoneTrackers: React.FC<DoneTrackersProps> = ({addTagsToTracker}) => {
    const trackersResult = useQuery<TimeSpansQuery, TimeSpansQueryVariables>(gqlTimeSpan.TimeSpans, {
        variables: {cursor: {pageSize: 30}},
    });
    const loading = React.useRef(false);
    const tagsResult = useQuery<TagsQuery>(gqlTag.Tags);
    const [currentDate, setCurrentDate] = React.useState(moment());
    const loadMoreRef = React.useRef<HTMLDivElement | null>(null);
    useInterval(
        () => {
            if (!isSameDate(currentDate, moment())) {
                setCurrentDate(moment());
            }
        },
        1000,
        true
    );

    const fetchMore = React.useCallback(() => {
        if (!trackersResult || !trackersResult.data || trackersResult.loading || loading.current) {
            return;
        }
        loading.current = true;
        const {offset, pageSize, startId} = trackersResult.data.timeSpans.cursor;
        trackersResult
            .fetchMore({
                variables: {
                    cursor: {
                        startId,
                        offset,
                        pageSize,
                    },
                },
                updateQuery: (prev, {fetchMoreResult}): TimeSpansQuery => {
                    if (!fetchMoreResult) {
                        return prev;
                    }

                    return {
                        timeSpans: {
                            __typename: 'PagedTimeSpans',
                            timeSpans: [...prev.timeSpans.timeSpans, ...fetchMoreResult.timeSpans.timeSpans],
                            cursor: fetchMoreResult.timeSpans.cursor,
                        },
                    };
                },
            })
            .finally(() => {
                loading.current = false;
            });
    }, [trackersResult]);

    // Loads the next page once the sentinel div at the bottom of the list scrolls into view,
    // replacing react-infinite's windowed/virtualized approach - this list is a personal time
    // log, not large enough to need DOM virtualization.
    React.useEffect(() => {
        const target = loadMoreRef.current;
        if (!target) {
            return;
        }
        const observer = new IntersectionObserver(
            (entries) => {
                if (entries[0].isIntersecting) {
                    fetchMore();
                }
            },
            {rootMargin: '2000px 0px'}
        );
        observer.observe(target);
        return () => observer.disconnect();
    }, [fetchMore]);

    const values: GroupedTimeSpanProps = React.useMemo(() => {
        if (
            trackersResult.error ||
            trackersResult.loading ||
            !trackersResult.data ||
            trackersResult.data.timeSpans === null ||
            tagsResult.error ||
            tagsResult.loading ||
            !tagsResult.data ||
            tagsResult.data.tags == null
        ) {
            return [];
        }
        return toGroupedTimeSpanProps(trackersResult.data.timeSpans.timeSpans, tagsResult.data.tags, currentDate);
    }, [trackersResult, tagsResult, currentDate]);

    return (
        <div style={{marginTop: 10}}>
            {values.map(({key, timeSpans}) => (
                <DatedTimeSpans key={key} name={key} timeSpans={timeSpans} addTagsToTracker={addTagsToTracker} />
            ))}
            <div ref={loadMoreRef} />
        </div>
    );
};

const DatedTimeSpans: React.FC<
    {
        name: string;
        timeSpans: TimeSpanProps[];
    } & DoneTrackersProps
> = ({name, timeSpans, addTagsToTracker}) => {
    return (
        <div key={name}>
            <Typography key={name} align="center" variant={'h5'}>
                {name}
            </Typography>
            {timeSpans.map((timeSpanProps) => (
                <TimeSpan key={timeSpanProps.id} {...timeSpanProps} addTagsToTracker={addTagsToTracker} />
            ))}
        </div>
    );
};
