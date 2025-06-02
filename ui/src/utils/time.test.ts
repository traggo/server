import {isValidDate, normalizeDate, parseRelativeTime, userFriendlyDate} from './time';
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

it('should parse', () => {
    expectSuccess(parseRelativeTime('now-1d', 'startOf', moment('2019-10-20T15:55:00'))).toEqual('2019-10-19 15:55:00');
    expectSuccess(parseRelativeTime('now-120s', 'startOf', moment('2019-10-20T15:55:15'))).toEqual('2019-10-20 15:53:15');
    expectSuccess(parseRelativeTime('now-1d-1h', 'startOf', moment('2019-10-20T15:55:00'))).toEqual('2019-10-19 14:55:00');
    expectSuccess(parseRelativeTime('now/w', 'startOf', moment('2019-10-20T15:55:15'))).toEqual('2019-10-14 00:00:00');
    expectSuccess(parseRelativeTime('now/w', 'endOf', moment('2019-10-20T15:55:15'))).toEqual('2019-10-20 23:59:59');
    expectSuccess(parseRelativeTime('now-1w/w', 'startOf', moment('2019-10-20T15:55:15'))).toEqual('2019-10-07 00:00:00');
    expectSuccess(parseRelativeTime('now-1y+1w/w', 'startOf', moment('2019-10-20T15:55:15'))).toEqual('2018-10-22 00:00:00');
    expectSuccess(parseRelativeTime('now/d+5h', 'startOf', moment('2019-10-20T15:55:00'))).toEqual('2019-10-20 05:00:00');
    expectSuccess(parseRelativeTime('now/y', 'startOf', moment('2019-10-20T15:55:15'))).toEqual('2019-01-01 00:00:00');
});

it('should convert to RFC3339 and back', () => {
    // can't put exact dates since moment doesn't allow overriding `.local()`'s timezone for unit tests
    const userDate = '2025-01-01 10:10';
    const rfcDate = moment(userDate).utc().format();

    expect(normalizeDate(userDate)).toBe(rfcDate);
    expect(userFriendlyDate(rfcDate)).toBe(userDate);
});

it('should not modify relative ranges', () => {
    expect(normalizeDate('now-1d')).toBe('now-1d');
    expect(normalizeDate('now-120s')).toBe('now-120s');
    expect(normalizeDate('now-1d-1h')).toBe('now-1d-1h');
    expect(normalizeDate('now/w')).toBe('now/w');
    expect(normalizeDate('now/w')).toBe('now/w');
    expect(normalizeDate('now-1w/w')).toBe('now-1w/w');
    expect(normalizeDate('now-1y+1w/w')).toBe('now-1y+1w/w');
    expect(normalizeDate('now/d+5h')).toBe('now/d+5h');
    expect(normalizeDate('now/y')).toBe('now/y');
});

const expectSuccess = (value: ReturnType<typeof parseRelativeTime>) => {
    if (value.success) {
        return expect(value.value.format('YYYY-MM-DD HH:mm:ss'));
    }
    expect(value.error).toEqual('no error');
    return expect('');
};
