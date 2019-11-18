import * as React from 'react';
import {useQuery} from '@apollo/react-hooks';
import * as gqlTimeSpan from '../gql/timeSpan';
import * as gqlTag from '../gql/tags';
import {Trackers} from '../gql/__generated__/Trackers';
import {Tags} from '../gql/__generated__/Tags';
import {TimeSpan} from './TimeSpan';
import {toTimeSpanProps} from './timespanutils';
import {Typography} from '@material-ui/core';

export const ActiveTrackers = () => {
    const trackersResult = useQuery<Trackers>(gqlTimeSpan.Trackers);
    const tagsResult = useQuery<Tags>(gqlTag.Tags);
    const values = React.useMemo(() => {
        if (
            trackersResult.error ||
            trackersResult.loading ||
            !trackersResult.data ||
            trackersResult.data.timers === null ||
            tagsResult.error ||
            tagsResult.loading ||
            !tagsResult.data ||
            tagsResult.data.tags === null
        ) {
            return [];
        }
        return toTimeSpanProps(trackersResult.data.timers, tagsResult.data.tags);
    }, [tagsResult, trackersResult]);

    if (!values.length) {
        return null;
    }

    return (
        <>
            <Typography align="center" variant="h5" style={{marginTop: 10}}>
                Active Timers
            </Typography>
            {values.map((value) => {
                return <TimeSpan key={value.id} {...value} />;
            })}
        </>
    );
};
