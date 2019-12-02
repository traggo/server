import moment from 'moment-timezone';
import prettyMs from 'pretty-ms';

export const timeRunning = (date: moment.Moment, now: moment.Moment) => {
    const d = inUserTz(now).unix() - inUserTz(date).unix();
    return prettyMs(d * 1000, {unitCount: 2}).substring(1);
};

export const timeRunningCalendar = (date: moment.Moment, now: moment.Moment) => {
    const d = inUserTz(now).unix() - inUserTz(date).unix();

    if (d < 60) {
        return '<1m';
    }

    return prettyMs(d * 1000, {compact: true}).substring(1);
};

export const uglyConvertToLocalTime = (m: moment.Moment): moment.Moment => {
    const withoutTimeZone: string = m.format('YYYY-MM-DDTHH:mm:ss');
    return moment(withoutTimeZone);
};

export const inUserTz = (m: moment.Moment): moment.Moment => m.tz(moment.tz.guess());
