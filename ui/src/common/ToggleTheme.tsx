import {ToggleTheme as ToggleThemeGQL} from '../gql/preferences.local';
import * as React from 'react';
import Button from '@material-ui/core/Button';
import {useMutation} from '@apollo/react-hooks';

export const ToggleTheme = ({className}: {className: string}) => {
    const [toggle] = useMutation(ToggleThemeGQL);
    return (
        <Button className={className} onClick={() => toggle()}>
            Toggle Theme
        </Button>
    );
};
