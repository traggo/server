import {StatsInterval} from '../../gql/__generated__/globalTypes';
import * as moment from 'moment';
import {expectNever} from '../../utils/never';

export const ofInterval = (interval: StatsInterval): ((date: moment.Moment) => string) => {
    switch (interval) {
        case StatsInterval.Weekly:
        case StatsInterval.Monthly:
        case StatsInterval.Yearly:
            return (d) => d.format('l');
        case StatsInterval.Hourly:
            return (d) => d.format('lll');
        case StatsInterval.Single:
        case StatsInterval.Daily:
            return (d) => d.format('dddd') + ', ' + d.format('l');
        default:
            return expectNever(interval);
    }
};
