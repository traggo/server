import * as React from 'react';
import {Grid, StyleRulesCallback, WithStyles} from '@material-ui/core';
import Paper from '@material-ui/core/Paper';
import Typography from '@material-ui/core/Typography';
import withStyles from '@material-ui/core/styles/withStyles';
import {LoginForm} from './LoginForm';
import {ToggleTheme} from '../common/ToggleTheme';
import Link from '@material-ui/core/Link';

const styles: StyleRulesCallback = (theme) => ({
    root: {
        ...theme.mixins.gutters(),
        paddingTop: theme.spacing.unit * 4,
        paddingBottom: theme.spacing.unit * 3,
        textAlign: 'center',
        maxWidth: 400,
        borderTop: `5px solid ${theme.palette.primary.main}`,
    },
    footerLink: {
        margin: '0 2px',
    },
    themeButton: {
        position: 'absolute',
        top: 5,
        right: 5,
    },
});

export const LoginPage = withStyles(styles, {withTheme: true})(
    class extends React.Component<WithStyles<typeof styles>> {
        public render(): React.ReactNode {
            const {classes} = this.props;
            return (
                <Grid container={true} direction="row" alignItems="center" justify="center" style={{height: '95%'}}>
                    <Grid item>
                        <Paper elevation={10} square={true} className={classes.root}>
                            <Typography variant="h1" component="h1" gutterBottom={true}>
                                traggo
                            </Typography>
                            <LoginForm />
                        </Paper>
                    </Grid>
                    <div style={{position: 'absolute', bottom: 10, display: 'flex'}}>
                        <Typography variant="subtitle1" component="span" className={classes.footerLink}>
                            <Link href="https://github.com/traggo/server">Source Code</Link>
                        </Typography>
                        <Typography variant="subtitle1" component="span" className={classes.footerLink}>
                            |
                        </Typography>
                        <Typography variant="subtitle1" component="span" className={classes.footerLink}>
                            <Link href="https://github.com/traggo/server/issues">Bug Tracker</Link>
                        </Typography>
                        <Typography variant="subtitle1" component="span" className={classes.footerLink}>
                            |
                        </Typography>
                        <Typography variant="subtitle1" component="span" className={classes.footerLink}>
                            v1.0.0@12efefef
                        </Typography>
                    </div>
                    <ToggleTheme className={classes.themeButton} />
                </Grid>
            );
        }
    }
);
