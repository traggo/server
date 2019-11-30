import * as React from 'react';
import {useSettings} from '../gql/settings';
import {CenteredSpinner} from '../common/CenteredSpinner';
import moment from 'moment';
import {DateLocale, WeekDay} from '../gql/__generated__/globalTypes';
import {expectNever} from '../utils/never';

const setLocale = (locale: DateLocale) => {
    switch (locale) {
        case DateLocale.English:
            moment.locale('en');
            return true;
        case DateLocale.German:
            moment.locale('de');
            return true;
        default:
            return expectNever(locale);
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
    const {done, firstDayOfTheWeek, dateLocale} = useSettings();

    React.useEffect(() => {
        if (!done) {
            return;
        }
        setLocale(dateLocale);
        moment.updateLocale(moment.locale(), {
            week: {
                dow: weekDayToMoment(firstDayOfTheWeek),
                doy: moment.localeData(moment.locale()).firstDayOfYear(),
            },
        });
    }, [dateLocale, firstDayOfTheWeek, done]);

    if (!done) {
        return <CenteredSpinner />;
    }

    return <>{children}</>;
};
