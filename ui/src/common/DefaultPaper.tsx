import Paper from '@material-ui/core/Paper';
import {makeStyles} from '@material-ui/core/styles';
import * as React from 'react';

const useStyles = makeStyles((theme) => ({
    root: {
        ...theme.mixins.gutters(),
        paddingTop: theme.spacing(4),
        paddingBottom: theme.spacing(3),
        textAlign: 'center',
        maxWidth: 400,
        borderTop: `5px solid ${theme.palette.primary.main}`,
    },
}));
export const DefaultPaper: React.FC = ({children}) => {
    const classes = useStyles();
    return (
        <Paper elevation={10} square={true} className={classes.root}>
            {children}
        </Paper>
    );
};
