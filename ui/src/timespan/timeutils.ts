// @ts-ignore
import Duration from 'duration';
import moment from 'moment-timezone';

export const timeRunning = (date: moment.Moment, now: moment.Moment) => {
    const d = new Duration(inUserTz(date).toDate(), inUserTz(now).toDate());

    if (d.minutes < 5) {
        return d.toString(1, 1);
    }
    if (d.hours < 24) {
        return d.toString(1, 2);
    }
    return d.toString(1, 3);
};

export const uglyConvertToLocalTime = (m: moment.Moment): moment.Moment => {
    const withoutTimeZone: string = m.format('YYYY-MM-DDTHH:mm:ss');
    return moment(withoutTimeZone);
};

export const inUserTz = (m: moment.Moment): moment.Moment => m.tz(moment.tz.guess());
