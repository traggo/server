import Paper from '@mui/material/Paper';
import makeStyles from '@mui/styles/makeStyles';
import * as React from 'react';

const useStyles = makeStyles((theme) => ({
    root: {
        paddingLeft: theme.spacing(2),
        paddingRight: theme.spacing(2),
        [theme.breakpoints.up('sm')]: {
            paddingLeft: theme.spacing(3),
            paddingRight: theme.spacing(3),
        },
        paddingTop: theme.spacing(4),
        paddingBottom: theme.spacing(3),
        textAlign: 'center',
        maxWidth: 400,
        borderTop: `5px solid ${theme.palette.primary.main}`,
    },
}));
export const DefaultPaper: React.FC<React.PropsWithChildren> = ({children}) => {
    const classes = useStyles();
    return (
        <Paper elevation={10} square={true} className={classes.root}>
            {children}
        </Paper>
    );
};
