import React from 'react';
import Downshift from 'downshift';
import {makeStyles, Theme} from '@material-ui/core/styles';
import TextField from '@material-ui/core/TextField';
import Paper from '@material-ui/core/Paper';
import MenuItem from '@material-ui/core/MenuItem';
import Chip from '@material-ui/core/Chip';
import {useQuery} from '@apollo/react-hooks';
import {Tags} from '../gql/__generated__/Tags';
import {TagSelectorEntry, label} from './tagSelectorEntry';
import * as gqlTags from '../gql/tags';
import {useSuggest} from './suggest';

const useStyles = makeStyles((theme: Theme) => ({
    root: {
        flexGrow: 1,
        height: 250,
    },
    container: {
        flexGrow: 1,
        position: 'relative',
    },
    paper: {
        position: 'absolute',
        zIndex: 1,
        marginTop: theme.spacing(1),
        left: 0,
        right: 0,
    },
    chip: {
        margin: theme.spacing(0.5, 0.25),
    },
    inputRoot: {
        flexWrap: 'wrap',
    },
    inputInput: {
        width: 'auto',
        flexGrow: 1,
    },
}));

interface TagFilterSelectorProps {
    value: TagSelectorEntry[];
    type: string;
    onChange: (entries: TagSelectorEntry[]) => void;
    disabled?: boolean;
}

export const TagFilterSelector: React.FC<TagFilterSelectorProps> = ({value: selectedItem, type, onChange, disabled = false}) => {
    const classes = useStyles();
    const [inputValue, setInputValue] = React.useState('');

    const tagsResult = useQuery<Tags>(gqlTags.Tags);
    const suggestions = useSuggest(tagsResult, inputValue, [])
        .filter((t) => !t.tag.create && !t.tag.alreadyUsed)
        .reverse();

    if (tagsResult.error || tagsResult.loading || !tagsResult.data || !tagsResult.data.tags) {
        return null;
    }

    function handleKeyDown(event: React.KeyboardEvent) {
        if (selectedItem.length && !inputValue.length && event.key === 'Backspace') {
            onChange(selectedItem.slice(0, selectedItem.length - 1));
        }
    }

    function handleInputChange(event: React.ChangeEvent<{value: string}>) {
        setInputValue(event.target.value);
    }

    function handleChange(item: TagSelectorEntry) {
        if (!item.value) {
            setInputValue(item.tag.key + ':');
            return;
        }
        let newSelectedItem = [...selectedItem];
        if (newSelectedItem.indexOf(item) === -1) {
            newSelectedItem = [...newSelectedItem, item];
        }
        setInputValue('');
        onChange(newSelectedItem);
    }

    const handleDelete = (item: TagSelectorEntry) => () => {
        const newSelectedItem = [...selectedItem];
        newSelectedItem.splice(newSelectedItem.indexOf(item), 1);
        onChange(newSelectedItem);
    };

    return (
        <Downshift
            id="downshift-multiple"
            inputValue={inputValue}
            onChange={handleChange}
            itemToString={(item) => (item ? label(item) : '')}
            defaultIsOpen={false}>
            {({getInputProps, getItemProps, getLabelProps, isOpen, highlightedIndex}) => {
                const {onBlur, onChange: downshiftOnChange, onFocus, ...inputProps} = getInputProps({
                    onKeyDown: handleKeyDown,
                    placeholder: `Select ${type} Tags`,
                });
                return (
                    <div className={classes.container}>
                        <TextField
                            InputLabelProps={getLabelProps()}
                            required={true}
                            disabled={disabled}
                            InputProps={{
                                classes: {
                                    root: classes.inputRoot,
                                    input: classes.inputInput,
                                },
                                startAdornment: selectedItem.map((item) => (
                                    <Chip
                                        key={label(item)}
                                        tabIndex={-1}
                                        label={label(item)}
                                        className={classes.chip}
                                        onDelete={disabled ? undefined : handleDelete(item)}
                                    />
                                )),
                                onBlur,
                                onChange: (event) => {
                                    handleInputChange(event);
                                    downshiftOnChange!(event as React.ChangeEvent<HTMLInputElement>);
                                },
                                onFocus,
                            }}
                            label={`${type} Tags`}
                            fullWidth
                            inputProps={inputProps}
                        />
                        {isOpen ? (
                            <Paper className={classes.paper} square>
                                {suggestions.map((suggestion, index) => (
                                    <MenuItem
                                        {...getItemProps({item: suggestion})}
                                        key={label(suggestion)}
                                        selected={highlightedIndex === index}
                                        component="div">
                                        {label(suggestion)}
                                    </MenuItem>
                                ))}
                            </Paper>
                        ) : null}
                    </div>
                );
            }}
        </Downshift>
    );
};
