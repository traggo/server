import React from 'react';
import Downshift from 'downshift';
import {makeStyles, Theme} from '@material-ui/core/styles';
import TextField from '@material-ui/core/TextField';
import Paper from '@material-ui/core/Paper';
import MenuItem from '@material-ui/core/MenuItem';
import Chip from '@material-ui/core/Chip';
import {useQuery} from 'react-apollo-hooks';
import {Tags} from '../gql/__generated__/Tags';
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

interface TagKeySelectorProps {
    value: string[];
    onChange: (entries: string[]) => void;
    disabled?: boolean;
}

export const TagKeySelector: React.FC<TagKeySelectorProps> = ({value: selectedItem, onChange, disabled = false}) => {
    const classes = useStyles();
    const [inputValue, setInputValue] = React.useState('');

    const tagsResult = useQuery<Tags>(gqlTags.Tags);
    const suggestions = useSuggest(tagsResult, inputValue, selectedItem, true)
        .filter((t) => !t.tag.create && !t.tag.alreadyUsed)
        .map((t) => t.tag.key);

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

    function handleChange(item: string) {
        let newSelectedItem = [...selectedItem];
        if (newSelectedItem.indexOf(item) === -1) {
            newSelectedItem = [...newSelectedItem, item];
        }
        setInputValue('');
        onChange(newSelectedItem);
    }

    const handleDelete = (item: string) => () => {
        const newSelectedItem = [...selectedItem];
        newSelectedItem.splice(newSelectedItem.indexOf(item), 1);
        onChange(newSelectedItem);
    };

    return (
        <Downshift
            id="downshift-multiple"
            inputValue={inputValue}
            onChange={handleChange}
            selectedItem={selectedItem}
            defaultIsOpen={false}>
            {({getInputProps, getItemProps, getLabelProps, isOpen, highlightedIndex}) => {
                const {onBlur, onChange: downshiftOnChange, onFocus, ...inputProps} = getInputProps({
                    onKeyDown: handleKeyDown,
                    placeholder: 'Select Tags',
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
                                        key={item}
                                        tabIndex={-1}
                                        label={item}
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
                            label={'Tags'}
                            fullWidth
                            inputProps={inputProps}
                        />
                        {isOpen ? (
                            <Paper className={classes.paper} square>
                                {suggestions.map((suggestion, index) => (
                                    <MenuItem
                                        {...getItemProps({item: suggestion})}
                                        key={suggestion}
                                        selected={highlightedIndex === index}
                                        component="div">
                                        {suggestion}
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
