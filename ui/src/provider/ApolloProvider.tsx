import * as React from 'react';
import {ApolloClient, ApolloProvider as Provider, HttpLink, InMemoryCache} from '@apollo/client';

const client = new ApolloClient({
    link: new HttpLink({uri: './graphql'}),
    cache: new InMemoryCache(),
});

export const ApolloProvider: React.FC<React.PropsWithChildren> = ({children}) => {
    return <Provider client={client}>{children}</Provider>;
};
