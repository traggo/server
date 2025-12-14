import * as React from 'react';
import {Tags} from '../gql/__generated__/Tags';
import * as gqlTags from '../gql/tags';
import {useQuery} from '@apollo/react-hooks';
import {useError} from '../utils/errors';
import Typography from '@material-ui/core/Typography';
import Tooltip from '@material-ui/core/Tooltip';
import MenuItem from '@material-ui/core/MenuItem';
import {AddTagDialog} from './AddTagDialog';
import {TagSelectorEntry, itemLabel, label, addValues} from './tagSelectorEntry';
import {useSuggest} from './suggest';
import Paper from '@material-ui/core/Paper';
import ClickAwayListener from '@material-ui/core/ClickAwayListener';
import Input from '@material-ui/core/Input';
import {useStateAndDelegateWithDelayOnChange} from '../utils/hooks';
import {TagChip} from '../common/TagChip';
import {makeStyles, Theme} from '@material-ui/core/styles';

const useStyles = makeStyles((theme: Theme) => ({
    root: {
        width: '100%',
    },
    inputRoot: {display: 'flex', flexWrap: 'wrap', cursor: 'text', width: '100%'},
    inputInput: {height: 40, minWidth: 150, flexGrow: 1},
    paper: {
        position: 'absolute',
        zIndex: 1,
        marginTop: theme.spacing(1),
    },
}));

export interface TagSelectorProps {
    onSelectedEntriesChanged: (entries: TagSelectorEntry[]) => void;
    selectedEntries: TagSelectorEntry[];
    dialogOpen?: React.Dispatch<React.SetStateAction<boolean>>;
    onCtrlEnter?: () => void;
    createTags?: boolean;
    allowDuplicateKeys?: boolean;
    onlySelectKeys?: boolean;
    removeWhenClicked?: boolean;
}

export const TagSelector: React.FC<TagSelectorProps> = ({
    selectedEntries,
    onSelectedEntriesChanged: setSelectedEntries,
    dialogOpen = () => {},
    onCtrlEnter,
    createTags = true,
    allowDuplicateKeys = false,
    onlySelectKeys = false,
    removeWhenClicked = false,
}) => {
    const classes = useStyles();
    const [tooltipErrorActive, tooltipError, showTooltipError] = useError(4000);
    const [open, setOpen] = React.useState(false);
    const [currentValue, setCurrentValueInternal] = React.useState('');
    const [highlightedIndex, setHighlightedIndex] = React.useState<number>(0);
    const [addDialogOpen, setAddDialogOpen] = useStateAndDelegateWithDelayOnChange<boolean>(false, dialogOpen);
    const input = React.useRef<null | HTMLDivElement>(null);
    const container = React.useRef<null | HTMLDivElement>(null);

    const tagsResult = useQuery<Tags>(gqlTags.Tags);

    const suggestions = useSuggest(tagsResult, currentValue, selectedEntries, onlySelectKeys, allowDuplicateKeys, createTags);

    if (tagsResult.error || tagsResult.loading || !tagsResult.data || !tagsResult.data.tags) {
        return null;
    }

    const setCurrentValue = (newValue: string) => {
        if (currentValue.indexOf(' ') !== -1) {
            throw new Error('old value should never contain a space');
        }
        setHighlightedIndex(0);

        if (newValue.indexOf(' ') !== -1) {
            const {errors, entries} = addValues(newValue, tagsResult, selectedEntries, onlySelectKeys, allowDuplicateKeys);

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

    const focusInput = () => {
        if (input.current) {
            input.current.focus();
        }
    };

    const trySubmit = (entry: TagSelectorEntry) => {
        if (entry.tag.alreadyUsed) {
            return;
        }
        if (entry.tag.create) {
            setAddDialogOpen(true);
            return;
        }

        focusInput();

        if (!onlySelectKeys && !entry.value) {
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

    const onTagClicked = (entry: TagSelectorEntry) => {
        if (!removeWhenClicked) {
            return;
        }
        const tagIndex = selectedEntries.indexOf(entry);
        selectedEntries.splice(tagIndex, 1);

        setSelectedEntries(selectedEntries);
    };

    const onKeyDown = (event: React.KeyboardEvent) => {
        if (!currentValue && selectedEntries.length && event.key === 'Backspace') {
            event.preventDefault();
            const last = selectedEntries[selectedEntries.length - 1];
            setSelectedEntries(selectedEntries.slice(0, selectedEntries.length - 1));
            setCurrentValue(event.ctrlKey ? '' : itemLabel(last, onlySelectKeys));
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
            if (!currentValue && event.ctrlKey && onCtrlEnter) {
                onCtrlEnter();
            } else {
                trySubmit(suggestions[highlightedIndex]);
            }
        }
        if (event.key === 'Escape' && input.current) {
            event.preventDefault();
            setOpen(false);
        }
        if (event.key === 'Tab') {
            setOpen(false);
        }
    };

    return (
        <ClickAwayListener onClickAway={() => setOpen(false)}>
            <div className={classes.root}>
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
                    <div ref={(ref) => (container.current = ref)} className={classes.inputRoot} onClick={focusInput}>
                        {toChips(selectedEntries, onlySelectKeys, onTagClicked)}
                        <Input
                            margin="none"
                            value={currentValue}
                            inputRef={(ref) => (input.current = ref)}
                            onFocus={() => setOpen(true)}
                            onKeyDown={onKeyDown}
                            disableUnderline={true}
                            onChange={(e) => setCurrentValue(e.target.value)}
                            placeholder="Enter Tags"
                            className={classes.inputInput}
                        />
                    </div>
                </Tooltip>

                {open ? (
                    <Paper
                        className={classes.paper}
                        style={{width: (container.current && container.current.clientWidth) || 300}}
                        square>
                        {suggestions.map((entry, index) => (
                            <Item
                                key={label(entry)}
                                entry={entry}
                                onClick={trySubmit}
                                selected={index === highlightedIndex}
                                onlySelectKeys={onlySelectKeys}
                            />
                        ))}
                    </Paper>
                ) : null}
                {addDialogOpen && (
                    <AddTagDialog
                        onAdded={(tag) => trySubmit({tag, value: ''})}
                        open={true}
                        initialName={currentValue.split(':')[0]}
                        close={() => {
                            setTimeout(focusInput, 50);
                            setAddDialogOpen(false);
                        }}
                    />
                )}
            </div>
        </ClickAwayListener>
    );
};

interface ItemProps {
    entry: TagSelectorEntry;
    selected: boolean;
    onlySelectKeys: boolean;
    onClick: (entry: TagSelectorEntry) => void;
}

const Item: React.FC<ItemProps> = ({entry, selected, onlySelectKeys, onClick}) => {
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
            {itemLabel(entry, onlySelectKeys)}
        </MenuItem>
    );
};

const toChips = (entries: TagSelectorEntry[], onlySelectKeys: boolean, onClick: (entry: TagSelectorEntry) => void) => {
    return entries.map((entry) => (
        <TagChip
            key={itemLabel(entry, onlySelectKeys)}
            label={itemLabel(entry, onlySelectKeys)}
            color={entry.tag.color}
            onClick={() => onClick(entry)}
        />
    ));
};
