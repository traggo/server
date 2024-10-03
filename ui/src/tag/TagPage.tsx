import * as React from 'react';
import Paper from '@material-ui/core/Paper';
import {makeStyles} from '@material-ui/core/styles';
import {useMutation, useQuery} from '@apollo/react-hooks';
import * as gqlTag from '../gql/tags';
import * as gqlDashboard from '../gql/dashboard';
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
import Button from '@material-ui/core/Button';
import {Tags} from '../gql/__generated__/Tags';
import {RemoveTag, RemoveTagVariables} from '../gql/__generated__/RemoveTag';
import {UpdateTag, UpdateTagVariables} from '../gql/__generated__/UpdateTag';
import {AddTagDialog} from './AddTagDialog';
import {SliderPicker} from 'react-color';
import {TagChip} from '../common/TagChip';
import {handleError} from '../utils/errors';
import {ConfirmDialog} from '../common/ConfirmDialog';

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

export const TagPage = () => {
    const classes = useStyles();
    const {data, loading} = useQuery<Tags>(gqlTag.Tags);
    const [removeTagConfirm, setRemoveTagConfirm] = React.useState('');
    const refetch = {refetchQueries: [{query: gqlTag.Tags}, {query: gqlDashboard.Dashboards}]};
    const {enqueueSnackbar} = useSnackbar();
    const [removeTag] = useMutation<RemoveTag, RemoveTagVariables>(gqlTag.RemoveTag, refetch);
    const [[editKey, editKeyNew, editColor], setEditing] = React.useState<[string, string, string]>(['', '', '']);
    const [addActive, setAddActive] = React.useState(false);
    const [updateTag] = useMutation<UpdateTag, UpdateTagVariables>(gqlTag.UpdateTag, refetch);
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
            }).then(() => enqueueSnackbar('tag edited', {variant: 'success'}));
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
                        />
                    ) : (
                        tag.key
                    )}
                </TableCell>
                <TableCell>
                    {isEdited ? (
                        <SliderPicker onChange={(c) => setEditing([editKey, editKeyNew, c.hex])} color={editColor} />
                    ) : (
                        <TagChip label={tag.color} color={tag.color} />
                    )}
                </TableCell>
                <TableCell align="right">{tag.usages}</TableCell>
                <TableCell align="right">
                    {isEdited ? (
                        <>
                            <IconButton onClick={onClickSubmit} title="Save">
                                <DoneIcon />
                            </IconButton>
                            <IconButton onClick={() => setEditing(['', '', ''])} title="Cancel">
                                <CloseIcon />
                            </IconButton>
                        </>
                    ) : (
                        <>
                            <IconButton onClick={() => setEditing([tag.key, tag.key, tag.color])} title="Edit">
                                <EditIcon />
                            </IconButton>
                            <IconButton onClick={() => setRemoveTagConfirm(tag.key)} title="Delete">
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
