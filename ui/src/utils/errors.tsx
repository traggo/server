import * as React from 'react';
import {ApolloError} from 'apollo-client/errors/ApolloError';
import {withSnackbarProps} from 'notistack';

export const handleError = (prefix: string, enqueue: withSnackbarProps['enqueueSnackbar']): ((error: ApolloError) => void) => {
    return (error) => {
        error.graphQLErrors.forEach((gqlError) => {
            enqueue(`${prefix}: ${gqlError.message}`, {variant: 'warning'});
        });
    };
};

export const useError = (timeout = 1000): [boolean, string, (s: string) => void] => {
    const [error, setError] = React.useState('');
    const [active, setActive] = React.useState(false);

    React.useLayoutEffect(() => {
        if (!active) {
            return;
        }
        const handle = setTimeout(() => {
            setActive(false);
        }, timeout);
        return () => clearTimeout(handle);
    }, [active]);

    return [
        active,
        error,
        (message: string) => {
            setError(message);
            setActive(true);
        },
    ];
};
