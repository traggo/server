import * as React from 'react';
import withStyles from '@material-ui/core/styles/withStyles';
import Paper from '@material-ui/core/Paper';
import {StyleRulesCallback, WithStyles} from '@material-ui/core/styles';
import {useMutation, useQuery} from 'react-apollo-hooks';
import * as gqlUser from '../gql/user';
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
import {TextField} from '@material-ui/core';
import {Users} from '../gql/__generated__/Users';
import {RemoveUser, RemoveUserVariables} from '../gql/__generated__/RemoveUser';
import {UpdateUser, UpdateUserVariables} from '../gql/__generated__/UpdateUser';
import Checkbox from '@material-ui/core/Checkbox';
import Button from '@material-ui/core/Button';
import {AddUserDialog} from './AddUserDialog';

const styles: StyleRulesCallback = (theme) => ({
    root: {
        ...theme.mixins.gutters(),
        paddingTop: theme.spacing(3),
        paddingBottom: theme.spacing(3),
        textAlign: 'center',
        maxWidth: 800,
        minWidth: 500,
        margin: '0 auto',
    },
});

const NoEdit = [-1, '', '', false] as const;

export const UsersPage = withStyles(styles)(({classes}: WithStyles<typeof styles>) => {
    const {data, loading} = useQuery<Users>(gqlUser.Users);
    const [addUser, setAddUser] = React.useState(false);
    const refetch = {refetchQueries: [{query: gqlUser.Users}, {query: gqlUser.CurrentUser}]};
    const {enqueueSnackbar} = useSnackbar();
    const removeUser = useMutation<RemoveUser, RemoveUserVariables>(gqlUser.RemoveUser, refetch);
    const [[editId, editName, editPass, editAdmin], setEditing] = React.useState<Readonly<[number, string, string, boolean]>>(
        NoEdit
    );
    const updateUser = useMutation<UpdateUser, UpdateUserVariables>(gqlUser.UpdateUser, refetch);
    if (loading || !data || !data.currentUser || !data.users) {
        return <CenteredSpinner />;
    }

    const users = data.users.map((user) => {
        const onClickDelete = () =>
            removeUser({variables: {id: user.id}}).then(() => enqueueSnackbar('user deleted', {variant: 'success'}));
        const onClickSubmit = () => {
            setEditing(NoEdit);
            updateUser({
                variables: {
                    id: editId,
                    name: editName,
                    admin: editAdmin,
                    pass: editPass || undefined,
                },
            }).then(() => enqueueSnackbar('user edited', {variant: 'success'}));
        };
        const isCurrent = user.id === data.currentUser!.id;
        const isEdited = editId === user.id;
        return (
            <TableRow selected={isCurrent} key={user.id}>
                <TableCell>{user.id}</TableCell>
                <TableCell>
                    {isEdited ? (
                        <TextField
                            value={editName}
                            onChange={(e) => setEditing([editId, e.target.value, '', editAdmin])}
                            onKeyDown={(e) => {
                                if (e.key === 'Enter') {
                                    onClickSubmit();
                                }
                            }}
                            onSubmit={onClickSubmit}
                        />
                    ) : (
                        user.name + (isCurrent ? ' (current)' : '')
                    )}
                </TableCell>
                <TableCell>
                    {isEdited ? (
                        <TextField
                            value={editPass}
                            placeholder={'use old password'}
                            type="password"
                            onChange={(e) => setEditing([editId, editName, e.target.value, editAdmin])}
                        />
                    ) : (
                        '********'
                    )}
                </TableCell>
                <TableCell>
                    {isEdited ? (
                        <Checkbox checked={editAdmin} onChange={(e) => setEditing([editId, editName, '', e.target.checked])} />
                    ) : user.admin ? (
                        'Yes'
                    ) : (
                        'No'
                    )}
                </TableCell>
                <TableCell>
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
                            <IconButton onClick={() => setEditing([user.id, user.name, '', user.admin])} title="Edit">
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
                onClick={() => setAddUser(true)}
                fullWidth
                style={{marginBottom: 10}}>
                Create User
            </Button>
            {addUser && <AddUserDialog open={true} close={() => setAddUser(false)} />}
            <Table>
                <TableHead>
                    <TableRow>
                        <TableCell>ID</TableCell>
                        <TableCell>Name</TableCell>
                        <TableCell>Password</TableCell>
                        <TableCell>Admin</TableCell>
                        <TableCell style={{width: 150}} />
                    </TableRow>
                </TableHead>
                <TableBody>{users}</TableBody>
            </Table>
        </Paper>
    );
});
