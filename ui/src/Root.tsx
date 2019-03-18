import React from 'react';
import {LoginPage} from './login/LoginPage';
import './global.css';
import 'typeface-roboto';
import Query, {QueryResult} from 'react-apollo/Query';
import Mutation from 'react-apollo/Mutation';
import Button from '@material-ui/core/Button';
import * as gqlUser from './gql/user';
import Typography from '@material-ui/core/Typography';
import {Logout} from './gql/__generated__/Logout';
import {CurrentUser} from './gql/__generated__/CurrentUser';
import {ThemeProvider} from './provider/ThemeProvider';
import {ApolloProvider} from './provider/ApolloProvider';
import {SnackbarProvider} from './provider/SnackbarProvider';
import {CenteredSpinner} from './common/CenteredSpinner';
import Grid from '@material-ui/core/Grid';
import {DefaultPaper} from './common/DefaultPaper';
import {ApolloError} from 'apollo-client/errors/ApolloError';

export const Root = () => {
    return (
        <ApolloProvider>
            <ThemeProvider>
                <SnackbarProvider>
                    <Query<CurrentUser> query={gqlUser.CurrentUser}>
                        {({loading, error, data, refetch}: QueryResult<CurrentUser>) => {
                            if (loading) {
                                return <CenteredSpinner />;
                            }
                            if (error) {
                                return <Error refetch={refetch} error={error} />;
                            }
                            if (!data || data.user === null) {
                                return <LoginPage />;
                            } else {
                                return (
                                    <React.Fragment>
                                        <Typography>Hello {data.user.name}</Typography>
                                        <Mutation<Logout>
                                            mutation={gqlUser.Logout}
                                            update={(cache) =>
                                                cache.writeQuery<CurrentUser>({
                                                    query: gqlUser.CurrentUser,
                                                    data: {user: null},
                                                })
                                            }>
                                            {(logout) => <Button onClick={() => logout()}>Logout</Button>}
                                        </Mutation>
                                    </React.Fragment>
                                );
                            }
                        }}
                    </Query>
                </SnackbarProvider>
            </ThemeProvider>
        </ApolloProvider>
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
