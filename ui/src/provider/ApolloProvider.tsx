import ApolloClient from 'apollo-boost';
import * as React from 'react';
import {ApolloProvider as Provider} from 'react-apollo';
import {ApolloProvider as ApolloProviderHooks} from '@apollo/react-hooks';

const client = new ApolloClient({
    uri: './graphql',
});

export const ApolloProvider: React.FC = ({children}) => {
    return (
        <Provider client={client}>
            <ApolloProviderHooks client={client}>{children}</ApolloProviderHooks>
        </Provider>
    );
};
