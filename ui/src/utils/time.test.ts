import {isValidDate, parseRelativeTime} from './time';
import moment from 'moment';

moment.updateLocale('en', {
    week: {
        dow: 1, // monday
        doy: moment.localeData('en').firstDayOfYear(),
    },
});

it('should test for valid date', () => {
    expect(isValidDate('2017-05-05', 'YYYY-MM-DD HH:mm')).toBe(false);
    expect(isValidDate('2017-05-05T15:23', 'YYYY-MM-DD HH:mm')).toBe(false);
    expect(isValidDate('2017-05-05 15:23', 'YYYY-MM-DD HH:mm')).toBe(true);

    expect(isValidDate('2019-10-20T15:55:00Z')).toBe(true);
    expect(isValidDate('2017-05-05 15:23')).toBe(false);
});

// 2018-10-15 Monday
// 2018-10-22 Monday

// 2019-10-07 Monday
// 2019-10-14 Monday
// 2019-10-21 Monday

it.each([
    {
        value: 'now-1d',
        divide: 'startOf',
        now: moment('2019-10-20T15:55:00'),
        expected: '2019-10-19 15:55:00',
        localized: 'now-1d',
        normalized: 'now-1d',
        shouldParse: true,
    },
    {
        value: 'now-120s',
        divide: 'startOf',
        now: moment('2019-10-20T15:55:15'),
        expected: '2019-10-20 15:53:15',
        localized: 'now-120s',
        normalized: 'now-120s',
        shouldParse: true,
    },
    {
        value: 'now-1d-1h',
        divide: 'startOf',
        now: moment('2019-10-20T15:55:00'),
        expected: '2019-10-19 14:55:00',
        localized: 'now-1d-1h',
        normalized: 'now-1d-1h',
        shouldParse: true,
    },
    {
        value: 'now/w',
        divide: 'startOf',
        now: moment('2019-10-20T15:55:15'),
        expected: '2019-10-14 00:00:00',
        localized: 'now/w',
        normalized: 'now/w',
        shouldParse: true,
    },
    {
        value: 'now/w',
        divide: 'endOf',
        now: moment('2019-10-20T15:55:15'),
        expected: '2019-10-20 23:59:59',
        localized: 'now/w',
        normalized: 'now/w',
        shouldParse: true,
    },
    {
        value: 'now-1w/w',
        divide: 'startOf',
        now: moment('2019-10-20T15:55:15'),
        expected: '2019-10-07 00:00:00',
        localized: 'now-1w/w',
        normalized: 'now-1w/w',
        shouldParse: true,
    },
    {
        value: 'now-1y+1w/w',
        divide: 'startOf',
        now: moment('2019-10-20T15:55:15'),
        expected: '2018-10-22 00:00:00',
        localized: 'now-1y+1w/w',
        normalized: 'now-1y+1w/w',
        shouldParse: true,
    },
    {
        value: 'now/d+5h',
        divide: 'startOf',
        now: moment('2019-10-20T15:55:00'),
        expected: '2019-10-20 05:00:00',
        localized: 'now/d+5h',
        normalized: 'now/d+5h',
        shouldParse: true,
    },
    {
        value: 'now/y',
        divide: 'startOf',
        now: moment('2019-10-20T15:55:15'),
        expected: '2019-01-01 00:00:00',
        localized: 'now/y',
        normalized: 'now/y',
        shouldParse: true,
    },
    {
        value: '2025-01-01 10:10',
        divide: 'startOf',
        expected: '2025-01-01 10:10:00',
        localized: '2025-01-01 10:10',
        normalized: moment('2025-01-01 10:10')
            .utc()
            .format(),
        shouldParse: true,
    },
    {
        value: '2025-01-01 10:10',
        divide: 'endOf',
        expected: '2025-01-01 10:10:00',
        localized: '2025-01-01 10:10',
        normalized: moment('2025-01-01 10:10')
            .utc()
            .format(),
        shouldParse: true,
    },
    {
        value: '2025-01-02',
        divide: 'startOf',
        expected: '2025-01-02 00:00:00',
        localized: '2025-01-02',
        normalized: moment('2025-01-02 00:00')
            .utc()
            .format(),
        shouldParse: true,
    },
    {
        value: '2025-01-02',
        divide: 'endOf',
        expected: '2025-01-02 23:59:59',
        localized: '2025-01-02',
        normalized: moment('2025-01-02 23:59:59')
            .utc()
            .format(),
        shouldParse: true,
    },
    {
        value: moment('2025-01-02 10:00:00').format(),
        divide: 'startOf',
        expected: '2025-01-02 10:00:00',
        localized: '2025-01-02 10:00',
        normalized: moment('2025-01-02 10:00:00')
            .utc()
            .format(),
        shouldParse: true,
    },
    {
        value: moment('2025-01-02 00:00:00').format(),
        divide: 'startOf',
        expected: '2025-01-02 00:00:00',
        localized: '2025-01-02',
        normalized: moment('2025-01-02 00:00:00')
            .utc()
            .format(),
        shouldParse: true,
    },
    {
        value: moment('2025-01-02 23:59:59')
            .utc()
            .format(),
        divide: 'endOf',
        expected: '2025-01-02 23:59:59',
        localized: '2025-01-02',
        normalized: moment('2025-01-02 23:59:59')
            .utc()
            .format(),
        shouldParse: true,
    },
    {
        value: 'invalid',
        divide: 'endOf',
        expected: "Expected valid date (e.g. 2020-01-01 16:30) or 'now' at index 0",
        localized: 'invalid',
        normalized: undefined,
        shouldParse: false,
    },
])('should parse', ({value, divide, now, expected, normalized, localized, shouldParse}) => {
    const result = parseRelativeTime(value, divide as 'startOf' | 'endOf', now);
    expect(result.success).toBe(shouldParse);
    if (result.success) {
        expect(result.preview.format('YYYY-MM-DD HH:mm:ss')).toEqual(expected);
        expect(result.normalized).toEqual(normalized);
        expect(result.localized).toEqual(localized);
    } else {
        expect(result.error).toEqual(expected);
    }
});
