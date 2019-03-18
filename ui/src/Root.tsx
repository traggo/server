import React from 'react';
import {LoginPage} from './login/LoginPage';
import './global.css';
import 'typeface-roboto';
import Query from 'react-apollo/Query';
import Mutation from 'react-apollo/Mutation';
import Button from '@material-ui/core/Button';
import * as gqlUser from './gql/user';
import Typography from '@material-ui/core/Typography';
import {Logout} from './gql/__generated__/Logout';
import {CurrentUser} from './gql/__generated__/CurrentUser';
import {ThemeProvider} from './provider/ThemeProvider';
import {ApolloProvider} from './provider/ApolloProvider';
import {SnackbarProvider} from './provider/SnackbarProvider';

export const Root = () => {
    return (
        <ApolloProvider>
            <ThemeProvider>
                <SnackbarProvider>
                    <Query query={gqlUser.CurrentUser}>
                        {({loading, error, data}) => {
                            if (loading) {
                                return 'loading';
                            }
                            if (error) {
                                return 'error';
                            }
                            if (data.user === null) {
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
