import * as React from 'react';
import Button from '@material-ui/core/Button';
import TextField from '@material-ui/core/TextField';
import Dialog from '@material-ui/core/Dialog';
import DialogActions from '@material-ui/core/DialogActions';
import DialogContent from '@material-ui/core/DialogContent';
import DialogContentText from '@material-ui/core/DialogContentText';
import DialogTitle from '@material-ui/core/DialogTitle';
import * as gqlDevice from '../gql/device';
import {FetchResult} from 'react-apollo/Mutation';
import {useMutation} from 'react-apollo-hooks';
import {InlineDateTimePicker} from 'material-ui-pickers';
import moment from 'moment-timezone';
import {CreateDevice, CreateDeviceVariables} from '../gql/__generated__/CreateDevice';
import {useSnackbar} from 'notistack';
import {handleError} from '../utils/errors';

interface AddTagDialogProps {
    initialName: string;
    open: boolean;
    close: () => void;
}

export const AddDeviceDialog: React.FC<AddTagDialogProps> = ({close, open, initialName}) => {
    const [token, setToken] = React.useState('');
    const [name, setName] = React.useState(initialName);
    const [expiresAt, setExpiresAt] = React.useState(moment().add(10, 'year'));
    const {enqueueSnackbar} = useSnackbar();

    const addDevice = useMutation<CreateDevice, CreateDeviceVariables>(gqlDevice.CreateDevice, {
        refetchQueries: [{query: gqlDevice.Devices}],
    });
    const submit = (e: React.FormEvent) => {
        e.preventDefault();
        addDevice({variables: {expiresAt, name}})
            .then((result: FetchResult<CreateDevice>) => {
                console.log(result);
                if (result.data && result.data.device) {
                    enqueueSnackbar('Client created', {variant: 'success'});
                    setToken(result.data.device.token);
                }
            })
            .catch(handleError('Create Dialog', enqueueSnackbar));
    };

    return (
        <Dialog
            open={open}
            onClose={() => {
                if (token === '') {
                    close();
                }
            }}
            aria-labelledby="form-dialog-title"
            fullWidth>
            {token === '' ? (
                <form onSubmit={submit} noValidate autoComplete="off">
                    <DialogTitle id="form-dialog-title">Create Device</DialogTitle>
                    <DialogContent>
                        <DialogContentText />
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
                        <InlineDateTimePicker
                            value={expiresAt}
                            onChange={setExpiresAt}
                            keyboard={true}
                            fullWidth
                            margin="dense"
                            label="Expire Date"
                        />
                    </DialogContent>
                    <DialogActions>
                        <Button onClick={close} color="primary">
                            Cancel
                        </Button>
                        <Button type="submit" onClick={submit} color="primary">
                            Create Device
                        </Button>
                    </DialogActions>
                </form>
            ) : (
                <>
                    <DialogTitle id="form-dialog-title">Device Created</DialogTitle>
                    <DialogContent>
                        <DialogContentText>
                            The device has the following authentication token, copy it and save it somewhere because you cannot
                            obtain the token after closing this dialog
                            <br /> <b>{token}</b>
                        </DialogContentText>
                    </DialogContent>
                    <DialogActions>
                        <Button onClick={close} color="primary">
                            Close
                        </Button>
                    </DialogActions>
                </>
            )}
        </Dialog>
    );
};
