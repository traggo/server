import * as React from 'react';
import Button from '@mui/material/Button';
import TextField from '@mui/material/TextField';
import Dialog from '@mui/material/Dialog';
import DialogActions from '@mui/material/DialogActions';
import DialogContent from '@mui/material/DialogContent';
import DialogContentText from '@mui/material/DialogContentText';
import DialogTitle from '@mui/material/DialogTitle';
import * as gqlDevice from '../gql/device';
import {FetchResult, useMutation} from '@apollo/client';
import {CreateDeviceMutation, CreateDeviceMutationVariables, DeviceType} from '../gql/__generated__';
import {useSnackbar} from 'notistack';
import {handleError} from '../utils/errors';
import {deviceTypeToString} from './typeutils';
import Select from '@mui/material/NativeSelect';
import FormControl from '@mui/material/FormControl';
import InputLabel from '@mui/material/InputLabel';

interface AddTagDialogProps {
    initialName: string;
    open: boolean;
    close: () => void;
}

export const AddDeviceDialog: React.FC<AddTagDialogProps> = ({close, open, initialName}) => {
    const [token, setToken] = React.useState('');
    const [name, setName] = React.useState(initialName);
    const [deviceType, setDeviceType] = React.useState(DeviceType.NoExpiry);
    const {enqueueSnackbar} = useSnackbar();

    const [addDevice] = useMutation<CreateDeviceMutation, CreateDeviceMutationVariables>(gqlDevice.CreateDevice, {
        refetchQueries: [{query: gqlDevice.Devices}],
    });
    const submit = (e: React.FormEvent) => {
        e.preventDefault();
        addDevice({variables: {deviceType, name}})
            .then((result: FetchResult<CreateDeviceMutation>) => {
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
                        <FormControl margin={'normal'} fullWidth>
                            <InputLabel>Expires after</InputLabel>
                            <Select value={deviceType} onChange={(e) => setDeviceType(e.target.value as DeviceType)} fullWidth>
                                <option value={DeviceType.NoExpiry}>{deviceTypeToString(DeviceType.NoExpiry)}</option>
                                <option value={DeviceType.ShortExpiry}>{deviceTypeToString(DeviceType.ShortExpiry)}</option>
                                <option value={DeviceType.LongExpiry}>{deviceTypeToString(DeviceType.LongExpiry)}</option>
                            </Select>
                        </FormControl>
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
