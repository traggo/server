import {gql} from 'apollo-boost';

export interface UserPreferences {
    theme: 'dark' | 'light';
}

export const Preferences = gql`
    {
        theme @client
    }
`;
