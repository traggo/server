import {gql, useQuery} from '@apollo/client';
import {SettingsQuery as SettingsQueryResponse, DateLocale, Theme, WeekDay, DateTimeInputStyle} from './__generated__';
import {stripTypename} from '../utils/strip';

export const Settings = gql`
    query Settings {
        userSettings {
            theme
            dateLocale
            firstDayOfTheWeek
            dateTimeInputStyle
        }
    }
`;

export const SetSettings = gql`
    mutation SetSettings($settings: InputUserSettings!) {
        setUserSettings(settings: $settings) {
            theme
            dateTimeInputStyle
        }
    }
`;

const defaultSettings = {
    theme: Theme.GruvboxDark,
    dateLocale: DateLocale.American,
    firstDayOfTheWeek: WeekDay.Monday,
    dateTimeInputStyle: DateTimeInputStyle.Fancy,
} as const;

export const useSettings = (): {done: boolean} & Omit<SettingsQueryResponse['userSettings'], '__typename'> => {
    const data = useQuery<SettingsQueryResponse>(Settings);

    if (data.loading || data.error || !data.data) {
        return {...defaultSettings, done: false};
    }
    return {done: true, ...stripTypename(data.data.userSettings)};
};
