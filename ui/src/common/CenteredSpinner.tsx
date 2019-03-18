import {CircularProgress} from '@material-ui/core';
import Grid from '@material-ui/core/Grid';
import * as React from 'react';

export const CenteredSpinner = () => (
    <Grid container={true} direction="row" alignItems="center" justify="center" style={{height: '95%'}}>
        <Grid item>
            <CircularProgress size={100} />
        </Grid>
    </Grid>
);
