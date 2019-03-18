import withStyles from '@material-ui/core/styles/withStyles';
import {Paper} from '@material-ui/core';
import {StyleRulesCallback, WithStyles} from '@material-ui/core/es';
import * as React from 'react';

const styles: StyleRulesCallback = (theme) => ({
    root: {
        ...theme.mixins.gutters(),
        paddingTop: theme.spacing.unit * 4,
        paddingBottom: theme.spacing.unit * 3,
        textAlign: 'center',
        maxWidth: 400,
        borderTop: `5px solid ${theme.palette.primary.main}`,
    },
});
export const DefaultPaper = withStyles(styles)(({children, classes}: WithStyles<typeof styles> & {children: React.ReactNode}) => {
    return (
        <Paper elevation={10} square={true} className={classes.root}>
            {children}
        </Paper>
    );
});
