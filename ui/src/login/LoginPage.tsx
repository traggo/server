import * as React from 'react';
import Typography from '@material-ui/core/Typography';
import Grid from '@material-ui/core/Grid';
import {LoginForm} from './LoginForm';
import Link from '@material-ui/core/Link';
import {DefaultPaper} from '../common/DefaultPaper';
import * as gqlVersion from '../gql/version';
import {useQuery} from '@apollo/react-hooks';
import {Version} from '../gql/__generated__/Version';
import makeStyles from '@material-ui/core/styles/makeStyles';

const useStyles = makeStyles(() => ({
    footerLink: {
        margin: '0 2px',
    },
}));

export const LoginPage = () => {
    const classes = useStyles();
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
        </Grid>
    );
};
