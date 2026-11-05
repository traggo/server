import * as React from 'react';
import './global.css';
import 'react-resizable/css/styles.css';
import 'react-grid-layout/css/styles.css';
import 'typeface-roboto';
import {ThemeProvider} from './provider/ThemeProvider';
import {ApolloProvider} from './provider/ApolloProvider';
import {SnackbarProvider} from './provider/SnackbarProvider';
import {LocalizationProvider} from '@mui/x-date-pickers';
import {AdapterMoment} from '@mui/x-date-pickers/AdapterMoment';
import {Router} from './Router';
import {HashRouter} from 'react-router-dom';
import {BootUserSettings} from './provider/UserSettingsProvider';

export const Root = () => {
    return (
        <ApolloProvider>
            <BootUserSettings>
                <ThemeProvider>
                    <LocalizationProvider dateAdapter={AdapterMoment}>
                        <SnackbarProvider>
                            <HashRouter>
                                <Router />
                            </HashRouter>
                        </SnackbarProvider>
                    </LocalizationProvider>
                </ThemeProvider>
            </BootUserSettings>
        </ApolloProvider>
    );
};
