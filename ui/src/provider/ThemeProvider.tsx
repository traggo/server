import * as React from 'react';
import MuiThemeProvider from '@material-ui/core/styles/MuiThemeProvider';
import createMuiTheme from '@material-ui/core/styles/createMuiTheme';
import CssBaseline from '@material-ui/core/CssBaseline';
import {Query, QueryResult} from 'react-apollo';
import {Preferences, UserPreferences} from '../gql/preferences.local';

const dark = createMuiTheme({
    overrides: {
        MuiLink: {
            root: {
                color: '#3498db',
            },
        },
    },
    palette: {
        primary: {
            main: '#455a64',
        },
        secondary: {
            main: '#f44336',
        },
        type: 'dark',
    },
});

const light = createMuiTheme({
    overrides: {
        MuiLink: {
            root: {
                color: '#2980b9',
            },
        },
    },
    palette: {
        background: {default: '#eeeeee'},
        primary: {
            main: '#455a64',
        },
        secondary: {
            main: '#f44336',
        },
        type: 'light',
    },
});

export const ThemeProvider: React.FC = ({children}) => {
    return (
        <Query<UserPreferences> query={Preferences}>
            {(pref: QueryResult<UserPreferences>) => {
                return (
                    <MuiThemeProvider theme={(pref.data && pref.data.theme === 'dark' && dark) || light}>
                        <CssBaseline />
                        {children}
                    </MuiThemeProvider>
                );
            }}
        </Query>
    );
};
