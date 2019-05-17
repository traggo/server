import * as React from 'react';
import Button from '@material-ui/core/Button';
import TextField from '@material-ui/core/TextField';
import Dialog from '@material-ui/core/Dialog';
import DialogActions from '@material-ui/core/DialogActions';
import DialogContent from '@material-ui/core/DialogContent';
import DialogTitle from '@material-ui/core/DialogTitle';
import * as gqlUser from '../gql/user';
import {useMutation} from 'react-apollo-hooks';
import {useSnackbar} from 'notistack';
import {handleError} from '../utils/errors';
import {Checkbox} from '@material-ui/core';
import {CreateUser, CreateUserVariables} from '../gql/__generated__/CreateUser';
import FormControlLabel from '@material-ui/core/FormControlLabel';

interface AddTagDialogProps {
    open: boolean;
    close: () => void;
}

export const AddUserDialog: React.FC<AddTagDialogProps> = ({close, open}) => {
    const [name, setName] = React.useState('');
    const [pass, setPass] = React.useState('');
    const [admin, setAdmin] = React.useState(false);
    const {enqueueSnackbar} = useSnackbar();

    const addUser = useMutation<CreateUser, CreateUserVariables>(gqlUser.CreateUser, {
        refetchQueries: [{query: gqlUser.Users}],
    });
    const submit = (e: React.FormEvent) => {
        e.preventDefault();
        addUser({variables: {name, pass, admin}})
            .then(() => {
                enqueueSnackbar('User created', {variant: 'success'});
                close();
            })
            .catch(handleError('Create User', enqueueSnackbar));
    };

    return (
        <Dialog open={open} onClose={close} aria-labelledby="form-dialog-title" fullWidth>
            <form onSubmit={submit} noValidate autoComplete="off">
                <DialogTitle id="form-dialog-title">Create User</DialogTitle>
                <DialogContent>
                    <TextField
                        autoFocus
                        margin="dense"
                        id="name"
                        label="Name"
                        type="text"
                        fullWidth
                        value={name}
                        onChange={(e) => setName(e.target.value)}
                    />
                    <TextField
                        margin="dense"
                        label="Password"
                        type="password"
                        fullWidth
                        value={pass}
                        onChange={(e) => setPass(e.target.value)}
                    />
                    <FormControlLabel
                        control={
                            <Checkbox checked={admin} onChange={(e) => setAdmin(e.target.checked)}>
                                Admin
                            </Checkbox>
                        }
                        label="Admin"
                    />
                </DialogContent>
                <DialogActions>
                    <Button onClick={close} color="primary">
                        Cancel
                    </Button>
                    <Button type="submit" onClick={submit} color="primary">
                        Create User
                    </Button>
                </DialogActions>
            </form>
        </Dialog>
    );
};
