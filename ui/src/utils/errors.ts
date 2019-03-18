import {ApolloError} from 'apollo-client/errors/ApolloError';
import {withSnackbarProps} from 'notistack';

export const handleError = (prefix: string, enqueue: withSnackbarProps['enqueueSnackbar']): ((error: ApolloError) => void) => {
    return (error) => {
        error.graphQLErrors.forEach((gqlError) => {
            enqueue(`${prefix}: ${gqlError.message}`, {variant: 'warning'});
        });
    };
};
