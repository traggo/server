import ApolloClient from 'apollo-boost';
import {bootPreferences, setPreferences} from '../user/preferences';
import * as React from 'react';
import {ApolloProvider as Provider} from 'react-apollo';
import {ApolloProvider as ApolloProviderHooks} from '@apollo/react-hooks';
import {Preferences, UserPreferences} from '../gql/preferences.local';

const client = new ApolloClient({
    uri: './graphql',
    resolvers: {
        Mutation: {
            toggleTheme: (_root, _, {cache}: Record<'cache', ApolloClient<{}>>) => {
                const data = cache.readQuery<UserPreferences>({query: Preferences});
                setPreferences(cache, {...data, theme: data && data.theme === 'light' ? 'dark' : 'light'});
                return null;
            },
        },
    },
});
bootPreferences(client);

export const ApolloProvider: React.FC = ({children}) => {
    return (
        <Provider client={client}>
            <ApolloProviderHooks client={client}>{children}</ApolloProviderHooks>
        </Provider>
    );
};
