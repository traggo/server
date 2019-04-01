import * as React from 'react';
import {TextField} from '@material-ui/core';
import {Tags} from '../gql/__generated__/Tags';
import Chip from '@material-ui/core/Chip';
import {TagDefinitionType} from '../gql/__generated__/globalTypes';
import * as gqlTags from '../gql/tags';
import {useQuery} from 'react-apollo-hooks';
import {useError} from '../utils/errors';
import Typography from '@material-ui/core/Typography';
import Tooltip from '@material-ui/core/Tooltip';
import MenuItem from '@material-ui/core/MenuItem';
import {AddTagDialog} from './AddTagDialog';
import {TagSelectorEntry, itemLabel, label, addValues} from './tagSelectorEntry';
import {useSuggest} from './suggest';
import Paper from '@material-ui/core/Paper';
import ClickAwayListener from '@material-ui/core/ClickAwayListener';

export const TagSelector = () => {
    const [tooltipErrorActive, tooltipError, showTooltipError] = useError(4000);
    const [open, setOpen] = React.useState(false);
    const [currentValue, setCurrentValueInternal] = React.useState('');
    const [highlightedIndex, setHighlightedIndex] = React.useState<number>(0);
    const [addDialogOpen, setAddDialogOpen] = React.useState(false);
    const [selectedEntries, setSelectedEntries] = React.useState<TagSelectorEntry[]>([]);
    const input = React.useRef<null | HTMLDivElement>(null);

    const tagsResult = useQuery<Tags>(gqlTags.Tags);
    const suggestions = useSuggest(tagsResult, currentValue, selectedEntries);

    if (tagsResult.error || tagsResult.loading || !tagsResult.data || !tagsResult.data.tags) {
        return null;
    }

    const setCurrentValue = (newValue: string) => {
        if (currentValue.indexOf(' ') !== -1) {
            throw new Error('old value should never contain a space');
        }
        setHighlightedIndex(0);

        if (newValue.indexOf(' ') !== -1) {
            const {errors, entries} = addValues(newValue, tagsResult, selectedEntries);

            setSelectedEntries([...selectedEntries, ...entries]);

            if (errors.length > 0) {
                if (errors.length === 1) {
                    showTooltipError(errors[0].error);
                } else {
                    showTooltipError('skipped tags because:\n' + errors.map((e) => e.error).join('\n'));
                }
                setCurrentValueInternal(errors[errors.length - 1].value);
            } else {
                setCurrentValueInternal('');
            }
            return;
        }

        setCurrentValueInternal(newValue);
    };

    const trySubmit = (entry: TagSelectorEntry) => {
        if (entry.tag.alreadyUsed) {
            return;
        }
        if (entry.tag.create) {
            setAddDialogOpen(true);
            return;
        }

        if (input.current !== null) {
            input.current.focus();
        }
        if (entry.tag.type === TagDefinitionType.singlevalue && !entry.value) {
            const newValue = entry.tag.key + ':';
            if (currentValue !== newValue) {
                setHighlightedIndex(0);
                setCurrentValue(newValue);
            } else {
                showTooltipError('enter a value after the colon');
            }
            return;
        }

        setSelectedEntries([...selectedEntries, entry]);
        setCurrentValue('');
        return;
    };

    const onKeyDown = (event: React.KeyboardEvent) => {
        if (!currentValue && selectedEntries.length && event.key === 'Backspace') {
            event.preventDefault();
            const last = selectedEntries[selectedEntries.length - 1];
            setSelectedEntries(selectedEntries.slice(0, selectedEntries.length - 1));
            setCurrentValue(event.ctrlKey ? '' : label(last));
        }
        if (event.key === 'ArrowUp') {
            event.preventDefault();
            setHighlightedIndex(highlightedIndex === 0 ? suggestions.length - 1 : highlightedIndex - 1);
        }
        if (event.key === 'ArrowDown') {
            event.preventDefault();
            setHighlightedIndex(highlightedIndex === suggestions.length - 1 ? 0 : highlightedIndex + 1);
        }
        if (event.key === 'Enter' && highlightedIndex < suggestions.length) {
            event.preventDefault();
            trySubmit(suggestions[highlightedIndex]);
        }
    };

    return (
        <ClickAwayListener onClickAway={() => setOpen(false)}>
            <div>
                <Tooltip
                    disableFocusListener
                    disableHoverListener
                    disableTouchListener
                    open={tooltipErrorActive}
                    placement={'top'}
                    title={
                        <Typography color="inherit" style={{whiteSpace: 'pre-line'}}>
                            {tooltipError}
                        </Typography>
                    }>
                    <TextField
                        fullWidth
                        margin="normal"
                        variant="outlined"
                        value={currentValue}
                        onChange={(e) => setCurrentValue(e.target.value)}
                        InputProps={{
                            inputRef: (ref) => (input.current = ref),
                            onFocus: () => setOpen(true),
                            onKeyDown,
                            startAdornment: toChips(selectedEntries),
                        }}
                        label="Tags"
                    />
                </Tooltip>

                {open ? (
                    <Paper>
                        {suggestions.map((entry, index) => (
                            <Item key={label(entry)} entry={entry} onClick={trySubmit} selected={index === highlightedIndex} />
                        ))}
                    </Paper>
                ) : null}
                {addDialogOpen && (
                    <AddTagDialog
                        onAdded={(tag) => trySubmit({tag})}
                        open={true}
                        initialName={currentValue.split(':')[0]}
                        close={() => setAddDialogOpen(false)}
                    />
                )}
            </div>
        </ClickAwayListener>
    );
};

interface ItemProps {
    entry: TagSelectorEntry;
    selected: boolean;
    onClick: (entry: TagSelectorEntry) => void;
}

const Item: React.FC<ItemProps> = ({entry, selected, onClick}) => {
    return (
        <MenuItem
            key={entry.tag.key}
            selected={selected}
            component="a"
            onClick={() => onClick(entry)}
            disabled={entry.tag.alreadyUsed}
            style={{
                fontWeight: selected ? 500 : 400,
            }}>
            {itemLabel(entry)}
        </MenuItem>
    );
};

const toChips = (entries: TagSelectorEntry[]) => {
    return entries.map((entry) => <TagChip key={label(entry)} entry={entry} />);
};

const TagChip = ({entry}: {entry: TagSelectorEntry}) => {
    return <Chip tabIndex={-1} style={{background: entry.tag.color, marginRight: 10}} label={label(entry)} />;
};
