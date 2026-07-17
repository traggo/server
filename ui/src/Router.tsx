import * as React from 'react';
import {useQuery} from '@apollo/client';
import {CurrentUserQuery} from './gql/__generated__';
import * as gqlUser from './gql/user';
import {CenteredSpinner} from './common/CenteredSpinner';
import {LoginPage} from './login/LoginPage';
import {Typography} from '@mui/material';
import Button from '@mui/material/Button';
import {ApolloError} from '@apollo/client';
import Grid from '@mui/material/Grid';
import {DefaultPaper} from './common/DefaultPaper';
import {Page} from './common/Page';
import {Navigate, Route, Routes} from 'react-router-dom';
import {DailyPage} from './timespan/DailyPage';
import {DevicesPage} from './devices/DevicesPage';
import {UsersPage} from './user/UsersPage';
import {DashboardPage} from './dashboard/DashboardPage';
import {DashboardsPage} from './dashboard/DashboardsPage';
import {TagPage} from './tag/TagPage';
import {SettingsPage} from './setting/SettingsPage';
import {CalendarPage} from './timespan/calendar/CalendarPage';

export const Router = () => {
    const {loading, error, data, refetch} = useQuery<CurrentUserQuery>(gqlUser.CurrentUser);
    if (loading) {
        return <CenteredSpinner />;
    }
    if (error) {
        return <Error refetch={refetch} error={error} />;
    }
    const loggedIn = data && data.user;
    const admin = data && data.user && data.user.admin;

    return (
        <Routes>
            <Route path="/user/login" element={loggedIn ? <Navigate to="/" /> : <LoginPage />} />
            <Route
                path="/*"
                element={
                    loggedIn ? (
                        <Page>
                            <Routes>
                                <Route path="dashboards" element={<DashboardsPage />} />
                                <Route path="dashboard/:id/*" element={<DashboardPage />} />
                                <Route path="timesheet/list" element={<DailyPage />} />
                                <Route path="timesheet/calendar" element={<CalendarPage />} />
                                <Route path="user/settings" element={<SettingsPage />} />
                                <Route path="user/devices" element={<DevicesPage />} />
                                <Route path="user/tags" element={<TagPage />} />
                                {admin ? <Route path="admin/users" element={<UsersPage />} /> : null}
                                <Route path="/" element={<Navigate to="/timesheet/list" replace />} />
                            </Routes>
                        </Page>
                    ) : (
                        <Navigate to="/user/login" />
                    )
                }
            />
        </Routes>
    );
};

const Error: React.FC<{error: ApolloError; refetch: () => void}> = ({error, refetch}) => {
    return (
        <Grid container direction="row" alignItems="center" justifyContent="center" style={{height: '100%'}}>
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
