import * as React from 'react';
import Paper from '@material-ui/core/Paper';
import {makeStyles} from '@material-ui/core/styles';
import {useMutation, useQuery} from '@apollo/react-hooks';
import * as gqlDevice from '../gql/device';
import * as gqlUser from '../gql/user';
import {Devices} from '../gql/__generated__/Devices';
import {CenteredSpinner} from '../common/CenteredSpinner';
import Table from '@material-ui/core/Table';
import TableBody from '@material-ui/core/TableBody';
import TableCell from '@material-ui/core/TableCell';
import TableHead from '@material-ui/core/TableHead';
import TableRow from '@material-ui/core/TableRow';
import moment from 'moment-timezone';
import DeleteIcon from '@material-ui/icons/Delete';
import EditIcon from '@material-ui/icons/Edit';
import DoneIcon from '@material-ui/icons/Done';
import CloseIcon from '@material-ui/icons/Close';
import IconButton from '@material-ui/core/IconButton';
import {RemoveDevice, RemoveDeviceVariables} from '../gql/__generated__/RemoveDevice';
import {UpdateDevice, UpdateDeviceVariables} from '../gql/__generated__/UpdateDevice';
import {useSnackbar} from 'notistack';
import {TextField} from '@material-ui/core';
import Button from '@material-ui/core/Button';
import {AddDeviceDialog} from './AddDeviceDialog';
import {DeviceType} from '../gql/__generated__/globalTypes';
import {deviceTypeToString} from './typeutils';
import Select from '@material-ui/core/NativeSelect';

const useStyles = makeStyles((theme) => ({
    root: {
        ...theme.mixins.gutters(),
        paddingTop: theme.spacing(3),
        paddingBottom: theme.spacing(3),
        textAlign: 'center',
        maxWidth: 1200,
        margin: '0 auto',
    },
}));

export const DevicesPage = () => {
    const classes = useStyles();
    const {data, loading} = useQuery<Devices>(gqlDevice.Devices);
    const refetch = {refetchQueries: [{query: gqlDevice.Devices}, {query: gqlUser.CurrentUser}]};
    const {enqueueSnackbar} = useSnackbar();
    const [removeDevice] = useMutation<RemoveDevice, RemoveDeviceVariables>(gqlDevice.RemoveDevice, refetch);
    const [[editId, editName, editDeviceType], setEditing] = React.useState<[number, string, DeviceType]>([
        -1,
        '',
        DeviceType.NoExpiry,
    ]);
    const [addActive, setAddActive] = React.useState(false);
    const [updateDevice] = useMutation<UpdateDevice, UpdateDeviceVariables>(gqlDevice.UpdateDevice, refetch);
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
                        />
                    ) : (
                        device.name + (isCurrent ? ' (current)' : '')
                    )}
                </TableCell>
                <TableCell title={device.createdAt}>{moment(device.createdAt).fromNow()}</TableCell>
                <TableCell title={device.type}>
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
                            <IconButton onClick={onClickSubmit} title="Save">
                                <DoneIcon />
                            </IconButton>
                            <IconButton onClick={() => setEditing([-1, '', DeviceType.NoExpiry])} title="Cancel">
                                <CloseIcon />
                            </IconButton>
                        </>
                    ) : (
                        <>
                            <IconButton onClick={() => setEditing([device.id, device.name, device.type])} title="Edit">
                                <EditIcon />
                            </IconButton>
                            <IconButton onClick={onClickDelete} title="Delete">
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
