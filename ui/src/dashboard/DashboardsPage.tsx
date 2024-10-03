import * as React from 'react';
import Paper from '@material-ui/core/Paper';
import {useMutation, useQuery} from '@apollo/react-hooks';
import {CenteredSpinner} from '../common/CenteredSpinner';
import Table from '@material-ui/core/Table';
import TableBody from '@material-ui/core/TableBody';
import TableCell from '@material-ui/core/TableCell';
import TableHead from '@material-ui/core/TableHead';
import TableRow from '@material-ui/core/TableRow';
import DeleteIcon from '@material-ui/icons/Delete';
import EditIcon from '@material-ui/icons/Edit';
import DoneIcon from '@material-ui/icons/Done';
import CloseIcon from '@material-ui/icons/Close';
import IconButton from '@material-ui/core/IconButton';
import {useSnackbar} from 'notistack';
import * as gqlDashboard from '../gql/dashboard';
import {TextField} from '@material-ui/core';
import Button from '@material-ui/core/Button';
import {Dashboards} from '../gql/__generated__/Dashboards';
import {RemoveDashboard, RemoveDashboardVariables} from '../gql/__generated__/RemoveDashboard';
import {UpdateDashboard, UpdateDashboardVariables} from '../gql/__generated__/UpdateDashboard';
import {AddDashboardDialog} from './AddDashboardDialog';
import makeStyles from '@material-ui/core/styles/makeStyles';
import {ConfirmDialog} from '../common/ConfirmDialog';

const useStyles = makeStyles((theme) => ({
    root: {
        ...theme.mixins.gutters(),
        paddingTop: theme.spacing(3),
        paddingBottom: theme.spacing(3),
        textAlign: 'center',
        maxWidth: 800,
        margin: '0 auto',
    },
}));

const NoEdit = [-1, ''] as const;

export const DashboardsPage = () => {
    const classes = useStyles();
    const {loading, data} = useQuery<Dashboards>(gqlDashboard.Dashboards);
    const [addDashboard, setAddDashboard] = React.useState(false);
    const [removeDashboardConfirm, setRemoveDashboardConfirm] = React.useState<false | [number, string]>(false);
    const refetch = {refetchQueries: [{query: gqlDashboard.Dashboards}]};
    const {enqueueSnackbar} = useSnackbar();
    const [removeDashboard] = useMutation<RemoveDashboard, RemoveDashboardVariables>(gqlDashboard.RemoveDashboard, refetch);
    const [[editId, editName], setEditing] = React.useState<Readonly<[number, string]>>(NoEdit);
    const [updateDashboard] = useMutation<UpdateDashboard, UpdateDashboardVariables>(gqlDashboard.UpdateDashboard, refetch);
    if (loading || !data || !data.dashboards) {
        return <CenteredSpinner />;
    }
    const onClickDelete = () => {
        if (removeDashboardConfirm) {
            removeDashboard({variables: {id: removeDashboardConfirm[0]}}).then(() =>
                enqueueSnackbar('dashboard deleted', {variant: 'success'})
            );
        }
    };

    const dashboards = data.dashboards.map((dashboard) => {
        const onClickSubmit = () => {
            setEditing(NoEdit);
            updateDashboard({
                variables: {
                    id: editId,
                    name: editName,
                },
            }).then(() => enqueueSnackbar('dashboard edited', {variant: 'success'}));
        };
        const isEdited = editId === dashboard.id;
        return (
            <TableRow key={dashboard.id}>
                <TableCell>{dashboard.id}</TableCell>
                <TableCell>
                    {isEdited ? (
                        <TextField
                            value={editName}
                            onChange={(e) => setEditing([editId, e.target.value])}
                            onKeyDown={(e) => {
                                if (e.key === 'Enter') {
                                    onClickSubmit();
                                }
                            }}
                            onSubmit={onClickSubmit}
                        />
                    ) : (
                        dashboard.name
                    )}
                </TableCell>
                <TableCell align="right">
                    {isEdited ? (
                        <>
                            <IconButton onClick={onClickSubmit} title="Save">
                                <DoneIcon />
                            </IconButton>
                            <IconButton onClick={() => setEditing(NoEdit)} title="Cancel">
                                <CloseIcon />
                            </IconButton>
                        </>
                    ) : (
                        <>
                            <IconButton onClick={() => setEditing([dashboard.id, dashboard.name])} title="Edit">
                                <EditIcon />
                            </IconButton>
                            <IconButton onClick={() => setRemoveDashboardConfirm([dashboard.id, dashboard.name])} title="Delete">
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
                onClick={() => setAddDashboard(true)}
                fullWidth
                style={{marginBottom: 10}}>
                Create Dashboard
            </Button>
            {addDashboard && <AddDashboardDialog open={true} close={() => setAddDashboard(false)} />}
            <Table>
                <TableHead>
                    <TableRow>
                        <TableCell>ID</TableCell>
                        <TableCell>Name</TableCell>
                        <TableCell style={{width: 150}} />
                    </TableRow>
                </TableHead>
                <TableBody>{dashboards}</TableBody>
            </Table>
            {removeDashboardConfirm ? (
                <ConfirmDialog
                    title={`Delete Dashboard ${removeDashboardConfirm[1]}`}
                    fClose={() => setRemoveDashboardConfirm(false)}
                    fOnSubmit={onClickDelete}>
                    <b>This operation cannot be undone.</b> Deleting the dashboard will remove all its dashboard entries.
                </ConfirmDialog>
            ) : null}
        </Paper>
    );
};
