import moment from 'moment';
import * as React from 'react';
import {timeRunning} from '../timespan/timeutils';
import useInterval from '@rooks/use-interval';

export const RelativeToNow: React.FC<{from: moment.Moment}> = ({from}) => {
    const [now, setNow] = React.useState(moment());

    useInterval(
        () => {
            setNow(moment());
        },
        1000,
        true
    );
    return <>{timeRunning(from, now)}</>;
};
