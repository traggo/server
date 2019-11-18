import * as React from 'react';
import {useQuery} from '@apollo/react-hooks';
import * as gqlTimeSpan from '../gql/timeSpan';
import * as gqlTag from '../gql/tags';
import {TimeSpan, TimeSpanProps} from './TimeSpan';
import {Tags} from '../gql/__generated__/Tags';
import useInterval from '@rooks/use-interval';
import moment from 'moment';
import {TimeSpans, TimeSpansVariables} from '../gql/__generated__/TimeSpans';
import {Typography} from '@material-ui/core';
import {GroupedTimeSpanProps, toGroupedTimeSpanProps} from './timespanutils';
import {TagSelectorEntry} from '../tag/tagSelectorEntry';
import ReactInfinite from 'react-infinite';
import {isSameDate} from '../utils/time';

interface DoneTrackersProps {
    addTagsToTracker?: (entries: TagSelectorEntry[]) => void;
}

export const DoneTrackers: React.FC<DoneTrackersProps> = ({addTagsToTracker}) => {
    const trackersResult = useQuery<TimeSpans, TimeSpansVariables>(gqlTimeSpan.TimeSpans, {
        variables: {cursor: {pageSize: 30}},
    });
    const loading = React.useRef(false);
    const tagsResult = useQuery<Tags>(gqlTag.Tags);
    const [infiniteLoading, setInfiniteLoading] = React.useState(false);
    const [currentDate, setCurrentDate] = React.useState(moment());
    const [heights, setHeights] = React.useState<Record<string, number>>({});
    useInterval(
        () => {
            if (!isSameDate(currentDate, moment())) {
                setCurrentDate(moment());
            }
        },
        1000,
        true
    );

    const fetchMore = () => {
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
                updateQuery: (prev, {fetchMoreResult}): TimeSpans => {
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
            .then(() => {
                loading.current = false;
                return setInfiniteLoading(false);
            })
            .catch(() => {
                loading.current = false;
                return setInfiniteLoading(false);
            });
    };

    const values: GroupedTimeSpanProps = React.useMemo(() => {
        if (
            trackersResult.error ||
            trackersResult.loading ||
            !trackersResult.data ||
            trackersResult.data.timeSpans === null ||
            tagsResult.error ||
            tagsResult.loading ||
            !tagsResult.data ||
            tagsResult.data.tags === null
        ) {
            return [];
        }
        return toGroupedTimeSpanProps(trackersResult.data.timeSpans.timeSpans, tagsResult.data.tags, currentDate);
    }, [trackersResult, tagsResult, currentDate]);

    return (
        <div style={{marginTop: 10}}>
            <ReactInfinite
                key={1}
                useWindowAsScrollContainer
                preloadBatchSize={window.innerHeight}
                onInfiniteLoad={fetchMore}
                isInfiniteLoading={infiniteLoading}
                infiniteLoadBeginEdgeOffset={2000}
                loadingSpinnerDelegate={
                    <Typography align={'center'} variant={'h5'}>
                        .. loading time spans ..
                    </Typography>
                }
                elementHeight={values.map((m) => heights[m.key] || 500)}>
                {values.map(({key, timeSpans}) => {
                    return (
                        <DatedTimeSpans
                            key={key}
                            name={key}
                            timeSpans={timeSpans}
                            addTagsToTracker={addTagsToTracker}
                            setHeight={setHeights}
                        />
                    );
                })}
            </ReactInfinite>
        </div>
    );
};

const DatedTimeSpans: React.FC<{
    name: string;
    setHeight: (cb: (height: Record<string, number>) => Record<string, number>) => void;
    timeSpans: TimeSpanProps[];
} & DoneTrackersProps> = ({name, timeSpans, addTagsToTracker, setHeight}) => {
    const [ref, setRef] = React.useState<HTMLDivElement | null>();
    React.useEffect(() => {
        if (ref) {
            setHeight((old) => ({...old, [name]: ref.getBoundingClientRect().height}));
        }
    }, [ref, name, setHeight]);
    return (
        <div key={name} ref={(r) => setRef(r)}>
            <Typography key={name} align="center" variant={'h5'}>
                {name}
            </Typography>
            {timeSpans.map((timeSpanProps) => (
                <TimeSpan key={timeSpanProps.id} {...timeSpanProps} addTagsToTracker={addTagsToTracker} />
            ))}
        </div>
    );
};
