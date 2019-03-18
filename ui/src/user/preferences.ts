import ApolloClient from 'apollo-client/ApolloClient';
import {UserPreferences} from '../gql/preferences.local';

const prefKey = '__preferences';

export const bootPreferences = (client: ApolloClient<{}>) => {
    setPreferences(client, cachedPreferences());
};

export const defaultPreferences = (): UserPreferences => ({theme: 'light'});

// tslint:disable-next-line:no-any
const isTheme = (value: any): value is 'dark' | 'light' => value === 'dark' || value === 'light';

export const cachedPreferences = (): UserPreferences => {
    const def = defaultPreferences();
    const pref = JSON.parse(localStorage.getItem(prefKey) || '{}') || def;
    return {
        theme: isTheme(pref.theme) ? pref.theme : def.theme,
    };
};

export const setPreferences = (client: ApolloClient<{}>, pref: UserPreferences) => {
    client.writeData<UserPreferences>({data: pref});

    localStorage.setItem(prefKey, JSON.stringify(pref));
};
