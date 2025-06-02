import {normalizeDate} from './range';
import moment from 'moment';

moment.updateLocale('en', {
    week: {
        dow: 1, // monday
        doy: moment.localeData('en').firstDayOfYear(),
    },
});
moment.tz.setDefault('UTC');

it('should convert to RFC3339', () => {
    expect(normalizeDate('2025-01-01 10:10')).toBe('2025-01-01T10:10:00Z');
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
