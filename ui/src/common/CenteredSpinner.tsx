import CircularProgress from '@mui/material/CircularProgress';
import Grid from '@mui/material/Grid';
import * as React from 'react';

export const CenteredSpinner = () => (
    <Grid container={true} direction="row" alignItems="center" justifyContent="center" style={{height: '95%'}}>
        <Grid item>
            <CircularProgress size={100} />
        </Grid>
    </Grid>
);
