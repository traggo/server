import * as React from 'react';
import './global.css';
import 'react-big-calendar/lib/css/react-big-calendar.css';
import 'react-big-calendar/lib/addons/dragAndDrop/styles.css';
import 'typeface-roboto';
import {ThemeProvider} from './provider/ThemeProvider';
import {ApolloProvider} from './provider/ApolloProvider';
import {SnackbarProvider} from './provider/SnackbarProvider';
import {MuiPickersUtilsProvider} from '@material-ui/pickers';
import MomentUtils from '@date-io/moment';
import {Router} from './Router';
import {HashRouter} from 'react-router-dom';

export const Root = () => {
    return (
        <ApolloProvider>
            <ThemeProvider>
                <MuiPickersUtilsProvider utils={MomentUtils}>
                    <SnackbarProvider>
                        <HashRouter>
                            <Router />
                        </HashRouter>
                    </SnackbarProvider>
                </MuiPickersUtilsProvider>
            </ThemeProvider>
        </ApolloProvider>
    );
};
