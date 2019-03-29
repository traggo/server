import * as React from 'react';
import MuiThemeProvider from '@material-ui/core/styles/MuiThemeProvider';
import createMuiTheme from '@material-ui/core/styles/createMuiTheme';
import CssBaseline from '@material-ui/core/CssBaseline';
import Query, {QueryResult} from 'react-apollo/Query';
import {Preferences, UserPreferences} from '../gql/preferences.local';

const dark = createMuiTheme({
    typography: {
        useNextVariants: true,
    },
    palette: {
        primary: {
            main: '#78909c',
        },
        secondary: {
            main: '#f44336',
        },
        type: 'dark',
    },
});

const light = createMuiTheme({
    typography: {
        useNextVariants: true,
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
