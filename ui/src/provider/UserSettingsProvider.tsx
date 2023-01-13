import * as React from 'react';
import {useSettings} from '../gql/settings';
import {CenteredSpinner} from '../common/CenteredSpinner';
import moment, {LocaleSpecification} from 'moment';
import 'moment/locale/en-gb'; // DD/MM/YYYY but 24hr
import 'moment/locale/en-au'; // DD/MM/YYYY but am/pm
import {DateFormat, DateLocale, WeekDay} from '../gql/__generated__/globalTypes';
import {expectNever} from '../utils/never';

const setLocale = (locale: DateLocale, dateFormat: DateFormat, spec: LocaleSpecification) => {
    switch (locale) {
        case DateLocale.English:
            if (dateFormat === DateFormat.MMDDYYYY) {
                moment.locale('en', spec);
            } else if (dateFormat === DateFormat.DDMMYYYY) {
                moment.locale('en-au', spec);
            } else {
                throw new Error('Unexpected date format');
            }
            return;
        case DateLocale.English24h:
            if (dateFormat === DateFormat.MMDDYYYY) {
                moment.locale('en', {
                    ...spec,
                    longDateFormat: {
                        LTS: 'HH:mm:ss',
                        LT: 'HH:mm',
                        L: 'MM/DD/YYYY',
                        LL: 'MMMM D, YYYY',
                        LLL: 'MMMM D, YYYY HH:mm',
                        LLLL: 'dddd, MMMM D, YYYY HH:mm',
                    },
                });
            } else if (dateFormat === DateFormat.DDMMYYYY) {
                moment.locale('en-gb', spec);
            } else {
                throw new Error('Unexpected date format');
            }
            return;
        case DateLocale.German:
            moment.locale('de', spec);
            return;
        default:
            expectNever(locale);
            return;
    }
};

const weekDayToMoment = (s: WeekDay): number => {
    switch (s) {
        case WeekDay.Sunday:
            return 0;
        case WeekDay.Monday:
            return 1;
        case WeekDay.Tuesday:
            return 2;
        case WeekDay.Wednesday:
            return 3;
        case WeekDay.Thursday:
            return 4;
        case WeekDay.Friday:
            return 5;
        case WeekDay.Saturday:
            return 6;
        default:
            return expectNever(s);
    }
};

export const BootUserSettings: React.FC = ({children}): React.ReactElement => {
    const {done, firstDayOfTheWeek, dateLocale, dateFormat} = useSettings();

    React.useEffect(() => {
        if (!done) {
            return;
        }
        setLocale(dateLocale, dateFormat, {
            week: {
                dow: weekDayToMoment(firstDayOfTheWeek),
                doy: moment.localeData(moment.locale()).firstDayOfYear(),
            },
        });
    }, [dateLocale, firstDayOfTheWeek, done, dateFormat]);

    if (!done) {
        return <CenteredSpinner />;
    }

    return <>{children}</>;
};
