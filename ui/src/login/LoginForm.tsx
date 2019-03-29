import * as React from 'react';
import TextField from '@material-ui/core/TextField';
import Button from '@material-ui/core/Button';
import {withSnackbar, withSnackbarProps} from 'notistack';
import {handleError} from '../utils/errors';
import * as gqlUser from '../gql/user';
import {StyleRulesCallback, WithStyles} from '@material-ui/core/styles';
import withStyles from '@material-ui/core/styles/withStyles';
import {Login, LoginVariables} from '../gql/__generated__/Login';
import {useMutation} from 'react-apollo-hooks';

const styles: StyleRulesCallback = (theme) => ({
    button: {
        marginTop: theme.spacing.unit,
    },
});

interface LoginFormProps extends withSnackbarProps, WithStyles<typeof styles> {}

export const LoginForm = withStyles(styles)(
    withSnackbar<LoginFormProps>(({classes, enqueueSnackbar}) => {
        const login = useMutation<Login, LoginVariables>(gqlUser.Login, {
            update: (cache, {data}) => {
                cache.writeQuery({query: gqlUser.CurrentUser, data: {user: data && data.login && data.login.user}});
            },
        });

        const [username, setUsername] = React.useState('');
        const [password, setPassword] = React.useState('');
        const submit = (e: React.FormEvent) => {
            e.preventDefault();
            login({
                variables: {
                    name: username,
                    pass: password,
                    expiresAt: '2030-06-11T10:00:00Z',
                },
            })
                .then(() => enqueueSnackbar('Login successful', {variant: 'success'}))
                .catch(handleError('Login failed', enqueueSnackbar));
        };
        return (
            <form noValidate autoComplete="off" onSubmit={submit}>
                <TextField
                    required
                    margin="dense"
                    variant="outlined"
                    label="username"
                    autoFocus
                    fullWidth
                    onChange={(e) => setUsername(e.target.value)}
                />
                <TextField
                    required
                    type="password"
                    margin="dense"
                    variant="outlined"
                    label="password"
                    fullWidth
                    onChange={(e) => setPassword(e.target.value)}
                />
                <Button
                    type="submit"
                    className={classes.button}
                    fullWidth
                    size="large"
                    onClick={submit}
                    variant="contained"
                    color="primary">
                    Login
                </Button>
            </form>
        );
    })
);
