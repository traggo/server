import * as React from 'react';
import Button from '@mui/material/Button';
import TextField from '@mui/material/TextField';
import Dialog from '@mui/material/Dialog';
import DialogActions from '@mui/material/DialogActions';
import DialogContent from '@mui/material/DialogContent';
import DialogTitle from '@mui/material/DialogTitle';
import * as gqlUser from '../gql/user';
import {useMutation} from '@apollo/client';
import {useSnackbar} from 'notistack';
import {handleError} from '../utils/errors';
import {Checkbox} from '@mui/material';
import {CreateUserMutation, CreateUserMutationVariables} from '../gql/__generated__';
import FormControlLabel from '@mui/material/FormControlLabel';

interface AddTagDialogProps {
    open: boolean;
    close: () => void;
}

export const AddUserDialog: React.FC<AddTagDialogProps> = ({close, open}) => {
    const [name, setName] = React.useState('');
    const [pass, setPass] = React.useState('');
    const [admin, setAdmin] = React.useState(false);
    const {enqueueSnackbar} = useSnackbar();

    const [addUser] = useMutation<CreateUserMutation, CreateUserMutationVariables>(gqlUser.CreateUser, {
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
                        control={<Checkbox checked={admin} onChange={(e) => setAdmin(e.target.checked)} />}
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
