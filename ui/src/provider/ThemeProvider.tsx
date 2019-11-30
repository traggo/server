import * as React from 'react';
import {MuiThemeProvider, Theme} from '@material-ui/core/styles';
import createMuiTheme from '@material-ui/core/styles/createMuiTheme';
import CssBaseline from '@material-ui/core/CssBaseline';
import {useSettings} from '../gql/settings';
import {Theme as SettingTheme} from '../gql/__generated__/globalTypes';

const themes: Record<SettingTheme, Theme> = {
    [SettingTheme.GruvboxDark]: createMuiTheme({
        overrides: {
            MuiLink: {
                root: {
                    color: '#458588',
                },
            },
            MuiIconButton: {
                root: {
                    color: 'inherit',
                },
            },
            MuiListItemIcon: {
                root: {
                    color: 'inherit',
                },
            },
            MuiToolbar: {
                root: {
                    background: '#a89984',
                },
            },
        },
        palette: {
            background: {
                default: '#282828',
                paper: '#32302f',
            },
            text: {
                primary: '#fbf1d4',
            },
            primary: {
                main: '#a89984',
            },
            secondary: {
                main: '#f44336',
            },
            type: 'dark',
        },
    }),
    [SettingTheme.GruvboxLight]: createMuiTheme({
        overrides: {
            MuiLink: {
                root: {
                    color: '#458588',
                },
            },
            MuiIconButton: {
                root: {
                    color: 'inherit',
                },
            },
            MuiListItemIcon: {
                root: {
                    color: 'inherit',
                },
            },
            MuiToolbar: {
                root: {
                    background: '#7c6f64',
                },
            },
        },
        palette: {
            background: {
                default: '#fbf1c7',
                paper: '#f9f5d7',
            },
            text: {
                primary: '#282828',
            },
            primary: {
                main: '#7c6f64',
            },
            secondary: {
                main: '#f44336',
            },
            type: 'light',
        },
    }),
    [SettingTheme.MaterialLight]: createMuiTheme({
        overrides: {
            MuiLink: {
                root: {
                    color: '#2980b9',
                },
            },
            MuiIconButton: {
                root: {
                    color: 'inherit',
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
    }),
    [SettingTheme.MaterialDark]: createMuiTheme({
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
    }),
};

export const ThemeProvider: React.FC = ({children}) => {
    const {theme} = useSettings();
    return (
        <MuiThemeProvider theme={themes[theme] || themes[SettingTheme.GruvboxDark]}>
            <CssBaseline />
            {children}
        </MuiThemeProvider>
    );
};
