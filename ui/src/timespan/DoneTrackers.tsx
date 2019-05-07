import * as React from 'react';
import {useQuery} from 'react-apollo-hooks';
import * as gqlTimeSpan from '../gql/timeSpan';
import * as gqlTag from '../gql/tags';
import {TimeSpan} from './TimeSpan';
import {Tags} from '../gql/__generated__/Tags';
import useInterval from '@rooks/use-interval';
import moment from 'moment';
import {TimeSpans} from '../gql/__generated__/TimeSpans';
import {Typography} from '@material-ui/core';
import {GroupedTimeSpanProps, toGroupedTimeSpanProps} from './timespanutils';
import {TagSelectorEntry} from '../tag/tagSelectorEntry';

interface DoneTrackersProps {
    addTagsToTracker?: (entries: TagSelectorEntry[]) => void;
}

export const DoneTrackers: React.FC<DoneTrackersProps> = ({addTagsToTracker}) => {
    const trackersResult = useQuery<TimeSpans>(gqlTimeSpan.TimeSpans);
    const tagsResult = useQuery<Tags>(gqlTag.Tags);
    const [currentDate, setCurrentDate] = React.useState(moment());
    useInterval(
        () => {
            setCurrentDate(moment());
        },
        1000,
        true
    );

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
        return toGroupedTimeSpanProps(trackersResult.data.timeSpans, tagsResult.data.tags, currentDate);
    }, [trackersResult, tagsResult, currentDate]);

    return (
        <div style={{marginTop: 10}}>
            {values.map(({key, timeSpans}) => {
                return (
                    <div key={key}>
                        <Typography key={key} align="center" variant={'h5'}>
                            {key}
                        </Typography>
                        {timeSpans.map((timeSpanProps) => (
                            <TimeSpan key={timeSpanProps.id} {...timeSpanProps} addTagsToTracker={addTagsToTracker} />
                        ))}
                    </div>
                );
            })}
        </div>
    );
};
