import * as React from 'react';
import {SnackbarProvider as Provider} from 'notistack';

export const SnackbarProvider: React.FC<React.PropsWithChildren> = ({children}) => {
    return (
        <Provider maxSnack={3} autoHideDuration={3500}>
            {children}
        </Provider>
    );
};
