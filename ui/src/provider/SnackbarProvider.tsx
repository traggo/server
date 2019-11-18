import * as React from 'react';
import {SnackbarProvider as Provider} from 'notistack';
import {makeStyles} from '@material-ui/core/styles';

const useStyles = makeStyles(() => ({
    error: {
        background: '#E53935',
        color: '#fff',
    },
    warning: {
        background: '#d35400',
        color: '#fff',
    },
    info: {
        background: '#2980b9',
        color: '#fff',
    },
}));

export const SnackbarProvider: React.FC = ({children}) => {
    const classes = useStyles();
    return (
        <Provider
            maxSnack={3}
            classes={{variantError: classes.error, variantWarning: classes.warning, variantInfo: classes.info}}
            autoHideDuration={3500}>
            {children}
        </Provider>
    );
};
