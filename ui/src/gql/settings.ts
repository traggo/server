import {gql} from 'apollo-boost';
import {useQuery} from '@apollo/react-hooks';
import {Settings as SettingsQueryResponse} from './__generated__/Settings';
import {DateLocale, Theme, WeekDay} from './__generated__/globalTypes';
import {stripTypename} from '../utils/strip';

export const Settings = gql`
    query Settings {
        userSettings {
            theme
            dateLocale
            firstDayOfTheWeek
        }
    }
`;

export const SetSettings = gql`
    mutation SetSettings($settings: InputUserSettings!) {
        setUserSettings(settings: $settings) {
            theme
        }
    }
`;

const defaultSettings = {
    theme: Theme.GruvboxDark,
    dateLocale: DateLocale.American,
    firstDayOfTheWeek: WeekDay.Monday,
} as const;

export const useSettings = (): {done: boolean} & Omit<SettingsQueryResponse['userSettings'], '__typename'> => {
    const data = useQuery<SettingsQueryResponse>(Settings);

    if (data.loading || data.error || !data.data) {
        return {...defaultSettings, done: false};
    }
    return {done: true, ...stripTypename(data.data.userSettings)};
};
