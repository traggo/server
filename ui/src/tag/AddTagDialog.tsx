import * as React from 'react';
import Button from '@material-ui/core/Button';
import TextField from '@material-ui/core/TextField';
import Dialog from '@material-ui/core/Dialog';
import DialogActions from '@material-ui/core/DialogActions';
import DialogContent from '@material-ui/core/DialogContent';
import DialogContentText from '@material-ui/core/DialogContentText';
import DialogTitle from '@material-ui/core/DialogTitle';
import {SliderPicker} from 'react-color';
import {InputLabel} from '@material-ui/core';
import FormControl from '@material-ui/core/FormControl';
import {MutationFetchResult} from 'react-apollo';
import {AddTag, AddTagVariables} from '../gql/__generated__/AddTag';
import * as gqlTags from '../gql/tags';
import {TagSelectorEntry} from './tagSelectorEntry';
import {useMutation} from '@apollo/react-hooks';
import {useSnackbar} from 'notistack';
import {handleError} from '../utils/errors';

interface AddTagDialogProps {
    initialName: string;
    open: boolean;
    close: () => void;
    onAdded?: (tag: TagSelectorEntry['tag']) => void;
}

export const AddTagDialog: React.FC<AddTagDialogProps> = ({close, open, initialName, onAdded = () => {}}) => {
    const [name, setName] = React.useState(initialName);
    const [color, setColor] = React.useState('#e6b3b3');
    const {enqueueSnackbar} = useSnackbar();

    const [addTag] = useMutation<AddTag, AddTagVariables>(gqlTags.AddTag, {refetchQueries: [{query: gqlTags.Tags}]});
    const submit = (e: React.FormEvent) => {
        e.preventDefault();
        addTag({variables: {name, color}})
            .then((result: MutationFetchResult<AddTag> | void) => {
                close();
                if (result && result.data && result.data.createTag) {
                    onAdded(result.data.createTag);
                }
            })
            .catch(handleError('Add Tag', enqueueSnackbar));
    };

    return (
        <Dialog open={open} onClose={close} aria-labelledby="form-dialog-title" fullWidth>
            <form onSubmit={submit} noValidate autoComplete="off">
                <DialogTitle id="form-dialog-title">Create Tag</DialogTitle>
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
                    <FormControl fullWidth margin="dense">
                        <InputLabel htmlFor="color-picker" shrink={true}>
                            Color
                        </InputLabel>
                        <div id="color-picker" style={{marginTop: 25}}>
                            <SliderPicker onChange={(c) => setColor(c.hex)} color={color} />
                        </div>
                    </FormControl>
                </DialogContent>
                <DialogActions>
                    <Button onClick={close} color="primary">
                        Cancel
                    </Button>
                    <Button type="submit" onClick={submit} color="primary">
                        Create Tag
                    </Button>
                </DialogActions>
            </form>
        </Dialog>
    );
};
