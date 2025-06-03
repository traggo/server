import * as React from 'react';
import TextField from '@material-ui/core/TextField';
import Button from '@material-ui/core/Button';
import {useSnackbar} from 'notistack';
import {handleError} from '../utils/errors';
import * as gqlUser from '../gql/user';
import {Login, LoginVariables} from '../gql/__generated__/Login';
import {useMutation} from '@apollo/react-hooks';
import {Checkbox} from '@material-ui/core';
import {DeviceType} from '../gql/__generated__/globalTypes';
import FormControlLabel from '@material-ui/core/FormControlLabel';
import makeStyles from '@material-ui/core/styles/makeStyles';
import {Settings as SettingsGQL} from '../gql/settings';

const useStyles = makeStyles((theme) => ({
    button: {
        marginTop: theme.spacing(1),
    },
}));

export const LoginForm = () => {
    const classes = useStyles();
    const [login] = useMutation<Login, LoginVariables>(gqlUser.Login, {
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
