import * as React from 'react';
import {StyleRulesCallback, WithStyles} from '@material-ui/core/styles';
import Typography from '@material-ui/core/Typography';
import Grid from '@material-ui/core/Grid';
import withStyles from '@material-ui/core/styles/withStyles';
import {LoginForm} from './LoginForm';
import {ToggleTheme} from '../common/ToggleTheme';
import Link from '@material-ui/core/Link';
import {DefaultPaper} from '../common/DefaultPaper';
import * as gqlVersion from '../gql/version';
import {useQuery} from 'react-apollo-hooks';
import {Version} from '../gql/__generated__/Version';

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

export const LoginPage = withStyles(styles, {withTheme: true})(({classes}: WithStyles<typeof styles>) => {
    const {data: {version = gqlVersion.VersionDefault.version} = gqlVersion.VersionDefault} = useQuery<Version>(
        gqlVersion.Version
    );
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
                    {version.name}@{version.commit.slice(0, 8)}
                </Typography>
            </div>
            <ToggleTheme className={classes.themeButton} />
        </Grid>
    );
});
