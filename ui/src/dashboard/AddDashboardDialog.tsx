import * as React from 'react';
import Button from '@mui/material/Button';
import TextField from '@mui/material/TextField';
import Dialog from '@mui/material/Dialog';
import DialogActions from '@mui/material/DialogActions';
import DialogContent from '@mui/material/DialogContent';
import DialogTitle from '@mui/material/DialogTitle';
import {useMutation} from '@apollo/client';
import {useSnackbar} from 'notistack';
import {handleError} from '../utils/errors';
import * as gqlDashboard from '../gql/dashboard';
import {CreateDashboardMutation, CreateDashboardMutationVariables} from '../gql/__generated__';

interface AddTagDialogProps {
    open: boolean;
    close: () => void;
}

export const AddDashboardDialog: React.FC<AddTagDialogProps> = ({close, open}) => {
    const [name, setName] = React.useState('');
    const {enqueueSnackbar} = useSnackbar();

    const [addUser] = useMutation<CreateDashboardMutation, CreateDashboardMutationVariables>(gqlDashboard.CreateDashboard, {
        refetchQueries: [{query: gqlDashboard.Dashboards}],
    });
    const submit = (e: React.FormEvent) => {
        e.preventDefault();
        addUser({variables: {name}})
            .then(() => {
                enqueueSnackbar('Dashboard created', {variant: 'success'});
                close();
            })
            .catch(handleError('Create Dashboard', enqueueSnackbar));
    };

    return (
        <Dialog open={open} onClose={close} aria-labelledby="form-dialog-title" fullWidth>
            <form onSubmit={submit} noValidate autoComplete="off">
                <DialogTitle id="form-dialog-title">Create Dashboard</DialogTitle>
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
                </DialogContent>
                <DialogActions>
                    <Button onClick={close} color="primary">
                        Cancel
                    </Button>
                    <Button type="submit" onClick={submit} color="primary">
                        Create Dashboard
                    </Button>
                </DialogActions>
            </form>
        </Dialog>
    );
};
