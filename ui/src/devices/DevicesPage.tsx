import * as React from 'react';
import withStyles from '@material-ui/core/styles/withStyles';
import Paper from '@material-ui/core/Paper';
import {StyleRulesCallback, WithStyles} from '@material-ui/core/styles';
import {useMutation, useQuery} from 'react-apollo-hooks';
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
import IconButton from '@material-ui/core/IconButton';
import {RemoveDevice, RemoveDeviceVariables} from '../gql/__generated__/RemoveDevice';
import {UpdateDevice, UpdateDeviceVariables} from '../gql/__generated__/UpdateDevice';
import {useSnackbar} from 'notistack';
import {TextField} from '@material-ui/core';
import {InlineDateTimePicker} from 'material-ui-pickers';

const styles: StyleRulesCallback = (theme) => ({
    root: {
        ...theme.mixins.gutters(),
        paddingTop: theme.spacing.unit * 3,
        paddingBottom: theme.spacing.unit * 3,
        textAlign: 'center',
        maxWidth: 1200,
        minWidth: 800,
        margin: '0 auto',
    },
});

export const DevicesPage = withStyles(styles)(({classes}: WithStyles<typeof styles>) => {
    const {data, loading} = useQuery<Devices>(gqlDevice.Devices);
    const refetch = {refetchQueries: [{query: gqlDevice.Devices}, {query: gqlUser.CurrentUser}]};
    const {enqueueSnackbar} = useSnackbar();
    const removeDevice = useMutation<RemoveDevice, RemoveDeviceVariables>(gqlDevice.RemoveDevice, refetch);
    const [[editId, editName, editExpiresIn], setEditing] = React.useState<[number, string, string]>([-1, '', '']);
    const updateDevice = useMutation<UpdateDevice, UpdateDeviceVariables>(gqlDevice.UpdateDevice, refetch);
    if (loading || !data || !data.currentDevice || !data.devices) {
        return <CenteredSpinner />;
    }

    const devices = data.devices.map((device) => {
        const onClickDelete = () =>
            removeDevice({variables: {id: device.id}}).then(() => enqueueSnackbar('device deleted', {variant: 'success'}));
        const onClickSubmit = () => {
            setEditing([-1, '', '']);
            updateDevice({
                variables: {
                    id: editId,
                    name: editName,
                    expiresAt: editExpiresIn,
                },
            }).then(() => enqueueSnackbar('device edited', {variant: 'success'}));
        };
        const isCurrent = device.id === data.currentDevice!.id;
        const isEdited = editId === device.id;
        return (
            <TableRow selected={isCurrent}>
                <TableCell>{device.id}</TableCell>
                <TableCell>
                    {isEdited ? (
                        <TextField
                            value={editName}
                            onChange={(e) => setEditing([editId, e.target.value, editExpiresIn])}
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
                <TableCell title={device.expiresAt}>
                    {isEdited ? (
                        <InlineDateTimePicker
                            value={editExpiresIn}
                            onChange={(val) => setEditing([device.id, editName, val])}
                            keyboard={true}
                        />
                    ) : (
                        moment(device.expiresAt).fromNow()
                    )}
                </TableCell>
                <TableCell title={device.activeAt}>{moment(device.activeAt).fromNow()}</TableCell>
                <TableCell>
                    {isEdited ? (
                        <IconButton onClick={onClickSubmit}>
                            <DoneIcon />
                        </IconButton>
                    ) : (
                        <IconButton onClick={() => setEditing([device.id, device.name, device.expiresAt])}>
                            <EditIcon />
                        </IconButton>
                    )}
                    <IconButton onClick={onClickDelete}>
                        <DeleteIcon />
                    </IconButton>
                </TableCell>
            </TableRow>
        );
    });

    return (
        <Paper elevation={1} square={true} className={classes.root}>
            <Table padding={'dense'}>
                <TableHead>
                    <TableRow>
                        <TableCell>ID</TableCell>
                        <TableCell>Name</TableCell>
                        <TableCell>Created</TableCell>
                        <TableCell>Expires</TableCell>
                        <TableCell>Last Active</TableCell>
                        <TableCell style={{width: 150}} />
                    </TableRow>
                </TableHead>
                <TableBody>{devices}</TableBody>
            </Table>
        </Paper>
    );
});
