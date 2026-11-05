import {Tracker, Tag, TimeSpanItem} from '../gql/types';
import {toTagSelectorEntry} from '../tag/tagSelectorEntry';
import moment from 'moment';
import {TimeSpanProps} from './TimeSpan';

export const toTimeSpanProps = (timers: Tracker[], tags: Tag[]): TimeSpanProps[] => {
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

type GroupedByIndex = Record<string, TimeSpanItem[]>;
const group =
    (startOfTomorrow: moment.Moment, startOfToday: moment.Moment, startOfYesterday: moment.Moment) =>
    (a: GroupedByIndex, current: TimeSpanItem): GroupedByIndex => {
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

export const toGroupedTimeSpanProps = (timeSpans: TimeSpanItem[], tags: Tag[], now: moment.Moment): GroupedTimeSpanProps => {
    const datesWithTimeSpans: GroupedByIndex = timeSpans.reduce(
        group(
            moment(now).add(1, 'day').startOf('day'),
            moment(now).startOf('day'),
            moment(now).subtract(1, 'day').startOf('day')
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
