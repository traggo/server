import withApollo from 'react-apollo/withApollo';
import {Preferences, UserPreferences} from '../gql/preferences.local';
import {setPreferences} from '../user/preferences';
import * as React from 'react';
import {Button} from '@material-ui/core';

export const ToggleTheme = withApollo<{className?: string}>(({client, className}) => {
    return (
        <Button
            className={className}
            onClick={() => {
                const data = client.readQuery<UserPreferences>({query: Preferences});
                if (data === null) {
                    throw new Error('illegal state');
                }
                data.theme = data.theme === 'light' ? 'dark' : 'light';
                setPreferences(client, data);
            }}>
            Toggle Theme
        </Button>
    );
});
