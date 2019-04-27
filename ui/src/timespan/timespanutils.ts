import {Trackers_timers} from '../gql/__generated__/Trackers';
import {Tags_tags} from '../gql/__generated__/Tags';
import {toTagSelectorEntry} from '../tag/tagSelectorEntry';
import moment from 'moment';
import {TimeSpanProps} from './TimeSpan';
import {Omit} from '../common/tsutil';
import {TimeSpans_timeSpans} from '../gql/__generated__/TimeSpans';

export const toTimeSpanProps = (timers: Trackers_timers[], tags: Tags_tags[]): Array<Omit<TimeSpanProps, 'now'>> => {
    return [...timers].map((timer) => {
        const tagEntries = toTagSelectorEntry(tags, timer.tags || []);
        const range: TimeSpanProps['range'] = {from: moment.parseZone(timer.start)};
        if (timer.end) {
            range.to = moment.parseZone(timer.end);
        }
        return {
            id: timer.id,
            range,
            initialTags: tagEntries,
        };
    });
};

type GroupedByIndex = Record<string, TimeSpans_timeSpans[]>;
const group = (startOfToday: moment.Moment, startOfYesterday: moment.Moment) => (
    a: GroupedByIndex,
    current: TimeSpans_timeSpans
): GroupedByIndex => {
    const startTime = moment(current.start);
    let date = startTime.format('DD. MMMM YY');
    if (startTime.isAfter(startOfToday)) {
        date = `Today, ${date}`;
    } else if (startTime.isAfter(startOfYesterday)) {
        date = `Yesterday, ${date}`;
    }
    a[date] = [...(a[date] || []), current];
    return a;
};

export type GroupedTimeSpanProps = Array<{key: string; timeSpans: Array<Omit<TimeSpanProps, 'now'>>}>;

export const toGroupedTimeSpanProps = (
    timeSpans: TimeSpans_timeSpans[],
    tags: Tags_tags[],
    now: moment.Moment
): GroupedTimeSpanProps => {
    const datesWithTimeSpans: GroupedByIndex = timeSpans.reduce(
        group(
            moment(now).startOf('day'),
            moment(now)
                .subtract(1, 'day')
                .startOf('day')
        ),
        {}
    );
    return Object.keys(datesWithTimeSpans).map((key) => {
        const groupedTimeSpans = datesWithTimeSpans[key];
        return {
            key,
            timeSpans: toTimeSpanProps(groupedTimeSpans, tags),
        };
    });
};
