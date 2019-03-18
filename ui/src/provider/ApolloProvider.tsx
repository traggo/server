import ApolloClient from 'apollo-boost';
import {bootPreferences} from '../user/preferences';
import * as React from 'react';
import {default as Provider} from 'react-apollo/ApolloProvider';

const client = new ApolloClient({
    uri: './graphql',
});
bootPreferences(client);

export const ApolloProvider: React.FC = ({children}) => {
    return <Provider client={client}>{children}</Provider>;
};
