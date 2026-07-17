import * as React from 'react';
import Paper from '@mui/material/Paper';
import makeStyles from '@mui/styles/makeStyles';
import {useMutation, useQuery} from '@apollo/client';
import * as gqlTag from '../gql/tags';
import * as gqlDashboard from '../gql/dashboard';
import {CenteredSpinner} from '../common/CenteredSpinner';
import Table from '@mui/material/Table';
import TableBody from '@mui/material/TableBody';
import TableCell from '@mui/material/TableCell';
import TableHead from '@mui/material/TableHead';
import TableRow from '@mui/material/TableRow';
import DeleteIcon from '@mui/icons-material/Delete';
import EditIcon from '@mui/icons-material/Edit';
import DoneIcon from '@mui/icons-material/Done';
import CloseIcon from '@mui/icons-material/Close';
import IconButton from '@mui/material/IconButton';
import {useSnackbar} from 'notistack';
import {TextField} from '@mui/material';
import Button from '@mui/material/Button';
import {
    TagsQuery,
    RemoveTagMutation,
    RemoveTagMutationVariables,
    UpdateTagMutation,
    UpdateTagMutationVariables,
} from '../gql/__generated__';
import {AddTagDialog} from './AddTagDialog';
import {HexColorPicker} from 'react-colorful';
import {TagChip} from '../common/TagChip';
import {handleError} from '../utils/errors';
import {ConfirmDialog} from '../common/ConfirmDialog';

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

export const TagPage = () => {
    const classes = useStyles();
    const {data, loading} = useQuery<TagsQuery>(gqlTag.Tags);
    const [removeTagConfirm, setRemoveTagConfirm] = React.useState('');
    const refetch = {refetchQueries: [{query: gqlTag.Tags}, {query: gqlDashboard.Dashboards}]};
    const {enqueueSnackbar} = useSnackbar();
    const [removeTag] = useMutation<RemoveTagMutation, RemoveTagMutationVariables>(gqlTag.RemoveTag, refetch);
    const [[editKey, editKeyNew, editColor], setEditing] = React.useState<[string, string, string]>(['', '', '']);
    const [addActive, setAddActive] = React.useState(false);
    const [updateTag] = useMutation<UpdateTagMutation, UpdateTagMutationVariables>(gqlTag.UpdateTag, refetch);
    if (loading || !data || !data.tags) {
        return <CenteredSpinner />;
    }
    const onClickDelete = () => {
        return removeTag({variables: {key: removeTagConfirm}})
            .then(() => enqueueSnackbar('tag deleted', {variant: 'success'}))
            .catch(handleError('Delete Tag', enqueueSnackbar));
    };

    const tags = data.tags.map((tag) => {
        const onClickSubmit = () => {
            setEditing(['', '', '']);
            updateTag({
                variables: {
                    key: editKey,
                    newKey: editKeyNew,
                    color: editColor,
                },
            })
                .then(() => enqueueSnackbar('tag edited', {variant: 'success'}))
                .catch(handleError('Edit Tag', enqueueSnackbar));
        };
        const isEdited = editKey === tag.key;
        return (
            <TableRow key={tag.key}>
                <TableCell>
                    {isEdited ? (
                        <TextField
                            value={editKeyNew}
                            onChange={(e) => setEditing([editKey, e.target.value, editColor])}
                            onKeyDown={(e) => {
                                if (e.key === 'Enter') {
                                    onClickSubmit();
                                }
                            }}
                            onSubmit={onClickSubmit}
                            style={{minWidth: 128}}
                        />
                    ) : (
                        tag.key
                    )}
                </TableCell>
                <TableCell style={{minWidth: 128}}>
                    {isEdited ? (
                        <HexColorPicker onChange={(c) => setEditing([editKey, editKeyNew, c])} color={editColor} />
                    ) : (
                        <TagChip label={tag.color} color={tag.color} />
                    )}
                </TableCell>
                <TableCell align="right">{tag.usages}</TableCell>
                <TableCell align="right">
                    {isEdited ? (
                        <>
                            <IconButton onClick={onClickSubmit} title="Save" size="large">
                                <DoneIcon />
                            </IconButton>
                            <IconButton onClick={() => setEditing(['', '', ''])} title="Cancel" size="large">
                                <CloseIcon />
                            </IconButton>
                        </>
                    ) : (
                        <>
                            <IconButton onClick={() => setEditing([tag.key, tag.key, tag.color])} title="Edit" size="large">
                                <EditIcon />
                            </IconButton>
                            <IconButton onClick={() => setRemoveTagConfirm(tag.key)} title="Delete" size="large">
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
                Create Tag
            </Button>
            <Table>
                <TableHead>
                    <TableRow>
                        <TableCell>Key</TableCell>
                        <TableCell>Color</TableCell>
                        <TableCell style={{width: 100}}>Usages</TableCell>
                        <TableCell style={{width: 150}} />
                    </TableRow>
                </TableHead>
                <TableBody>
                    {addActive ? <AddTagDialog initialName={''} open={true} close={() => setAddActive(false)} /> : null}
                    {tags}
                    {removeTagConfirm ? (
                        <ConfirmDialog
                            title={`Delete Tag ${removeTagConfirm}`}
                            fClose={() => setRemoveTagConfirm('')}
                            fOnSubmit={onClickDelete}>
                            <b>This operation cannot be undone.</b> Deleting the tag will remove all references in time spans.
                        </ConfirmDialog>
                    ) : null}
                </TableBody>
            </Table>
        </Paper>
    );
};
