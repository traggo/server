import * as React from 'react';
import Paper from '@mui/material/Paper';
import makeStyles from '@mui/styles/makeStyles';
import {useMutation, useQuery} from '@apollo/client';
import * as gqlDevice from '../gql/device';
import * as gqlUser from '../gql/user';
import {
    DevicesQuery,
    RemoveDeviceMutation,
    RemoveDeviceMutationVariables,
    UpdateDeviceMutation,
    UpdateDeviceMutationVariables,
    DeviceType,
} from '../gql/__generated__';
import {CenteredSpinner} from '../common/CenteredSpinner';
import Table from '@mui/material/Table';
import TableBody from '@mui/material/TableBody';
import TableCell from '@mui/material/TableCell';
import TableHead from '@mui/material/TableHead';
import TableRow from '@mui/material/TableRow';
import moment from 'moment-timezone';
import DeleteIcon from '@mui/icons-material/Delete';
import EditIcon from '@mui/icons-material/Edit';
import DoneIcon from '@mui/icons-material/Done';
import CloseIcon from '@mui/icons-material/Close';
import IconButton from '@mui/material/IconButton';
import {useSnackbar} from 'notistack';
import {TextField} from '@mui/material';
import Button from '@mui/material/Button';
import {AddDeviceDialog} from './AddDeviceDialog';
import {deviceTypeToString} from './typeutils';
import Select from '@mui/material/NativeSelect';

const useStyles = makeStyles((theme) => ({
    root: {
        paddingLeft: theme.spacing(2),
        paddingRight: theme.spacing(2),
        [theme.breakpoints.up('sm')]: {
            paddingLeft: theme.spacing(3),
            paddingRight: theme.spacing(3),
        },
        paddingTop: theme.spacing(3),
        paddingBottom: theme.spacing(3),
        textAlign: 'center',
        maxWidth: 1200,
        margin: '0 auto',
    },
}));

export const DevicesPage = () => {
    const classes = useStyles();
    const {data, loading} = useQuery<DevicesQuery>(gqlDevice.Devices);
    const refetch = {refetchQueries: [{query: gqlDevice.Devices}, {query: gqlUser.CurrentUser}]};
    const {enqueueSnackbar} = useSnackbar();
    const [removeDevice] = useMutation<RemoveDeviceMutation, RemoveDeviceMutationVariables>(gqlDevice.RemoveDevice, refetch);
    const [[editId, editName, editDeviceType], setEditing] = React.useState<[number, string, DeviceType]>([
        -1,
        '',
        DeviceType.NoExpiry,
    ]);
    const [addActive, setAddActive] = React.useState(false);
    const [updateDevice] = useMutation<UpdateDeviceMutation, UpdateDeviceMutationVariables>(gqlDevice.UpdateDevice, refetch);
    if (loading || !data || !data.currentDevice || !data.devices) {
        return <CenteredSpinner />;
    }

    const devices = data.devices.map((device) => {
        const onClickDelete = () =>
            removeDevice({variables: {id: device.id}}).then(() => enqueueSnackbar('device deleted', {variant: 'success'}));
        const onClickSubmit = () => {
            setEditing([-1, '', DeviceType.NoExpiry]);
            updateDevice({
                variables: {
                    id: editId,
                    name: editName,
                    deviceType: editDeviceType,
                },
            }).then(() => enqueueSnackbar('device edited', {variant: 'success'}));
        };
        const isCurrent = device.id === data.currentDevice!.id;
        const isEdited = editId === device.id;
        const displayName = device.name + (isCurrent ? ' (current)' : '');
        return (
            <TableRow key={device.id} selected={isCurrent}>
                <TableCell>{device.id}</TableCell>
                <TableCell>
                    {isEdited ? (
                        <TextField
                            value={editName}
                            onChange={(e) => setEditing([editId, e.target.value, editDeviceType])}
                            onKeyDown={(e) => {
                                if (e.key === 'Enter') {
                                    onClickSubmit();
                                }
                            }}
                            onSubmit={onClickSubmit}
                            style={{minWidth: 128}}
                        />
                    ) : (
                        displayName
                    )}
                </TableCell>
                <TableCell title={device.createdAt}>{moment(device.createdAt).fromNow()}</TableCell>
                <TableCell title={device.type} style={{minWidth: 128}}>
                    {isEdited ? (
                        <Select
                            value={editDeviceType}
                            onChange={(e) => setEditing([editId, editName, e.target.value as DeviceType])}>
                            <option value={DeviceType.NoExpiry}>{deviceTypeToString(DeviceType.NoExpiry)}</option>
                            <option value={DeviceType.ShortExpiry}>{deviceTypeToString(DeviceType.ShortExpiry)}</option>
                            <option value={DeviceType.LongExpiry}>{deviceTypeToString(DeviceType.LongExpiry)}</option>
                        </Select>
                    ) : (
                        deviceTypeToString(device.type)
                    )}
                </TableCell>
                <TableCell title={device.activeAt}>{moment(device.activeAt).fromNow()}</TableCell>
                <TableCell align="right">
                    {isEdited ? (
                        <>
                            <IconButton onClick={onClickSubmit} title="Save" size="large">
                                <DoneIcon />
                            </IconButton>
                            <IconButton onClick={() => setEditing([-1, '', DeviceType.NoExpiry])} title="Cancel" size="large">
                                <CloseIcon />
                            </IconButton>
                        </>
                    ) : (
                        <>
                            <IconButton
                                onClick={() => setEditing([device.id, device.name, device.type])}
                                title="Edit"
                                size="large">
                                <EditIcon />
                            </IconButton>
                            <IconButton onClick={onClickDelete} title="Delete" size="large">
                                <DeleteIcon />
                            </IconButton>
                        </>
                    )}
                </TableCell>
            </TableRow>
        );
    });

    return (
        <Paper elevation={1} square={true} className={classes.root}>
            <Button
                color={'primary'}
                variant={'outlined'}
                size="small"
                onClick={() => setAddActive(true)}
                fullWidth
                style={{marginBottom: 10}}>
                Create Device
            </Button>
            <Table>
                <TableHead>
                    <TableRow>
                        <TableCell>ID</TableCell>
                        <TableCell>Name</TableCell>
                        <TableCell>Created</TableCell>
                        <TableCell>Expires after</TableCell>
                        <TableCell>Last Active</TableCell>
                        <TableCell style={{width: 150}} />
                    </TableRow>
                </TableHead>
                <TableBody>
                    {addActive ? <AddDeviceDialog initialName={''} open={true} close={() => setAddActive(false)} /> : null}
                    {devices}
                </TableBody>
            </Table>
        </Paper>
    );
};
