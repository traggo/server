import {Trackers_timers} from '../gql/__generated__/Trackers';
import {Tags_tags} from '../gql/__generated__/Tags';
import {toTagSelectorEntry} from '../tag/tagSelectorEntry';
import moment from 'moment';
import {TimeSpanProps} from './TimeSpan';
import {TimeSpans_timeSpans_timeSpans} from '../gql/__generated__/TimeSpans';

export const toTimeSpanProps = (timers: Trackers_timers[], tags: Tags_tags[]): TimeSpanProps[] => {
    return [...timers].map((timer) => {
        const tagEntries = toTagSelectorEntry(tags, timer.tags || []);
        const range: TimeSpanProps['range'] = {from: moment.parseZone(timer.start)};
        if (timer.end) {
            range.to = moment.parseZone(timer.end);
        }
        return {
            id: timer.id,
            range: {
                ...range,
                oldFrom: timer.oldStart ? moment(timer.oldStart) : undefined,
            },
            initialTags: tagEntries,
            note: timer.note,
        };
    });
};

type GroupedByIndex = Record<string, TimeSpans_timeSpans_timeSpans[]>;
const group = (startOfTomorrow: moment.Moment, startOfToday: moment.Moment, startOfYesterday: moment.Moment) => (
    a: GroupedByIndex,
    current: TimeSpans_timeSpans_timeSpans
): GroupedByIndex => {
    const startTime = moment(current.oldStart || current.start);
    let date = `${startTime.format('dddd')}, ${startTime.format('LL')}`;
    if (startTime.isBetween(startOfToday, startOfTomorrow)) {
        date = `${date} (today)`;
    } else if (startTime.isBetween(startOfYesterday, startOfToday)) {
        date = `${date} (yesterday)`;
    }
    a[date] = [...(a[date] || []), current];
    return a;
};

export type GroupedTimeSpanProps = Array<{key: string; timeSpans: TimeSpanProps[]}>;

export const toGroupedTimeSpanProps = (
    timeSpans: TimeSpans_timeSpans_timeSpans[],
    tags: Tags_tags[],
    now: moment.Moment
): GroupedTimeSpanProps => {
    const datesWithTimeSpans: GroupedByIndex = timeSpans.reduce(
        group(
            moment(now)
                .add(1, 'day')
                .startOf('day'),
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
