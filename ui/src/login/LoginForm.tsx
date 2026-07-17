import * as React from 'react';
import TextField from '@mui/material/TextField';
import Button from '@mui/material/Button';
import {useSnackbar} from 'notistack';
import {handleError} from '../utils/errors';
import * as gqlUser from '../gql/user';
import {LoginMutation, LoginMutationVariables} from '../gql/__generated__';
import {useMutation} from '@apollo/client';
import {Checkbox} from '@mui/material';
import {DeviceType} from '../gql/__generated__';
import FormControlLabel from '@mui/material/FormControlLabel';
import makeStyles from '@mui/styles/makeStyles';
import {Settings as SettingsGQL} from '../gql/settings';

const useStyles = makeStyles((theme) => ({
    button: {
        marginTop: theme.spacing(1),
    },
}));

export const LoginForm = () => {
    const classes = useStyles();
    const [login] = useMutation<LoginMutation, LoginMutationVariables>(gqlUser.Login, {
        update: (cache, {data}) => {
            cache.writeQuery({query: gqlUser.CurrentUser, data: {user: data && data.login && data.login.user}});
        },
        refetchQueries: [{query: SettingsGQL}],
    });

    const {enqueueSnackbar} = useSnackbar();
    const [username, setUsername] = React.useState('');
    const [password, setPassword] = React.useState('');
    const [remember, setRemember] = React.useState(false);
    const submit = (e: React.FormEvent) => {
        e.preventDefault();
        login({
            variables: {
                name: username,
                pass: password,
                deviceType: remember ? DeviceType.LongExpiry : DeviceType.ShortExpiry,
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
            <FormControlLabel
                style={{float: 'right'}}
                control={<Checkbox checked={remember} onChange={(e) => setRemember(e.target.checked)} />}
                label="Remember Me"
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
};
