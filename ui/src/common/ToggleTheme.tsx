import {Preferences} from '../gql/preferences.local';
import * as React from 'react';
import Button from '@material-ui/core/Button';
import {useMutation} from 'react-apollo-hooks';

export const ToggleTheme = ({className}: {className: string}) => {
    const toggle = useMutation(Preferences);
    return (
        <Button className={className} onClick={() => toggle()}>
            Toggle Theme
        </Button>
    );
};
