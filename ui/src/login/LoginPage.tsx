import * as React from 'react';
import {StyleRulesCallback, WithStyles} from '@material-ui/core/styles';
import Typography from '@material-ui/core/Typography';
import Grid from '@material-ui/core/Grid';
import withStyles from '@material-ui/core/styles/withStyles';
import {LoginForm} from './LoginForm';
import {ToggleTheme} from '../common/ToggleTheme';
import Link from '@material-ui/core/Link';
import {DefaultPaper} from '../common/DefaultPaper';

const styles: StyleRulesCallback = () => ({
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
                        <DefaultPaper>
                            <Typography variant="h1" component="h1" gutterBottom={true}>
                                traggo
                            </Typography>
                            <LoginForm />
                        </DefaultPaper>
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
