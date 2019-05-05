import * as React from 'react';
import {useQuery} from 'react-apollo-hooks';
import {CurrentUser} from './gql/__generated__/CurrentUser';
import * as gqlUser from './gql/user';
import {CenteredSpinner} from './common/CenteredSpinner';
import {LoginPage} from './login/LoginPage';
import {Typography} from '@material-ui/core';
import Button from '@material-ui/core/Button';
import {ApolloError} from 'apollo-boost';
import Grid from '@material-ui/core/Grid';
import {DefaultPaper} from './common/DefaultPaper';
import {Page} from './common/Page';
import {Redirect, Route, Switch} from 'react-router';
import {DailyPage} from './timespan/DailyPage';
import {Calendar} from './timespan/calendar/Calendar';

export const Router = () => {
    const {loading, error, data, refetch} = useQuery<CurrentUser>(gqlUser.CurrentUser);
    if (loading) {
        return <CenteredSpinner />;
    }
    if (error) {
        return <Error refetch={refetch} error={error} />;
    }
    const loggedIn = data && data.user;

    return (
        <Switch>
            <Route path="/user/login">{loggedIn ? <Redirect to="/" /> : <LoginPage />}</Route>
            {loggedIn ? null : <Redirect to="/user/login" />}

            <Page>
                <Route exact path="/dashboard" component={TODO} />
                <Route exact path="/timesheet/daily" component={DailyPage} />
                <Route exact path="/timesheet/weekly" component={Calendar} />
                <Route exact path="/user/settings" component={TODO} />
                <Route exact path="/user/devices" component={TODO} />
                <Route exact path="/admin/users" component={TODO} />
                <Route exact path="/" render={() => <Redirect to="/dashboard" />} />
            </Page>
        </Switch>
    );
};

const Error: React.FC<{error: ApolloError; refetch: () => void}> = ({error, refetch}) => {
    return (
        <Grid container direction="row" alignItems="center" justify="center" style={{height: '100%'}}>
            <Grid item>
                <DefaultPaper>
                    <Typography variant="h3" component="h1" gutterBottom={true}>
                        Error
                    </Typography>
                    <Typography component="p">
                        {error.networkError && error.networkError.name + ': ' + error.networkError.message}
                        {error.graphQLErrors.map((gqlError) => gqlError.message).join(', ')}
                    </Typography>
                    <Button style={{marginTop: 10}} size="large" variant="outlined" onClick={() => refetch()}>
                        Retry
                    </Button>
                </DefaultPaper>
            </Grid>
        </Grid>
    );
};

export const TODO = () => {
    return (
        <Typography align="center" variant={'h2'}>
            TODO
        </Typography>
    );
};
