import * as React from 'react';
import Button from '@material-ui/core/Button';
import TextField from '@material-ui/core/TextField';
import Dialog from '@material-ui/core/Dialog';
import DialogActions from '@material-ui/core/DialogActions';
import DialogContent from '@material-ui/core/DialogContent';
import DialogTitle from '@material-ui/core/DialogTitle';
import {useMutation} from '@apollo/react-hooks';
import {useSnackbar} from 'notistack';
import {handleError} from '../utils/errors';
import * as gqlDashboard from '../gql/dashboard';
import {CreateDashboard, CreateDashboardVariables} from '../gql/__generated__/CreateDashboard';

interface AddTagDialogProps {
    open: boolean;
    close: () => void;
}

export const AddDashboardDialog: React.FC<AddTagDialogProps> = ({close, open}) => {
    const [name, setName] = React.useState('');
    const {enqueueSnackbar} = useSnackbar();

    const [addUser] = useMutation<CreateDashboard, CreateDashboardVariables>(gqlDashboard.CreateDashboard, {
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
