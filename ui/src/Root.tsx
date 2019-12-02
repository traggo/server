import * as React from 'react';
import './global.css';
import 'react-resizable/css/styles.css';
import 'typeface-roboto';
import {ThemeProvider} from './provider/ThemeProvider';
import {ApolloProvider} from './provider/ApolloProvider';
import {SnackbarProvider} from './provider/SnackbarProvider';
import {MuiPickersUtilsProvider} from '@material-ui/pickers';
import MomentUtils from '@date-io/moment';
import {Router} from './Router';
import {HashRouter} from 'react-router-dom';
import {BootUserSettings} from './provider/UserSettingsProvider';

export const Root = () => {
    return (
        <ApolloProvider>
            <BootUserSettings>
                <ThemeProvider>
                    <MuiPickersUtilsProvider utils={MomentUtils}>
                        <SnackbarProvider>
                            <HashRouter>
                                <Router />
                            </HashRouter>
                        </SnackbarProvider>
                    </MuiPickersUtilsProvider>
                </ThemeProvider>
            </BootUserSettings>
        </ApolloProvider>
    );
};
