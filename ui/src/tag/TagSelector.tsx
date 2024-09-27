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

export interface TagSelectorProps {
    onSelectedEntriesChanged: (entries: TagSelectorEntry[]) => void;
    selectedEntries: TagSelectorEntry[];
    dialogOpen?: React.Dispatch<React.SetStateAction<boolean>>;
    onCtrlEnter?: () => void;
}

export const TagSelector: React.FC<TagSelectorProps> = ({
    selectedEntries,
    onSelectedEntriesChanged: setSelectedEntries,
    dialogOpen = () => {},
    onCtrlEnter,
}) => {
    const [tooltipErrorActive, tooltipError, showTooltipError] = useError(4000);
    const [open, setOpen] = React.useState(false);
    const [currentValue, setCurrentValueInternal] = React.useState('');
    const [highlightedIndex, setHighlightedIndex] = React.useState<number>(0);
    const [addDialogOpen, setAddDialogOpen] = useStateAndDelegateWithDelayOnChange<boolean>(false, dialogOpen);
    const input = React.useRef<null | HTMLDivElement>(null);
    const container = React.useRef<null | HTMLDivElement>(null);

    const tagsResult = useQuery<Tags>(gqlTags.Tags);
    const suggestions = useSuggest(
        tagsResult,
        currentValue,
        selectedEntries.map((t) => t.tag.key)
    );

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

        if (!entry.value) {
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
    };

    return (
        <ClickAwayListener onClickAway={() => setOpen(false)}>
            <div style={{width: '100%'}}>
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
                    <div
                        ref={(ref) => (container.current = ref)}
                        style={{display: 'flex', flexWrap: 'wrap', cursor: 'text', width: '100%'}}
                        onClick={focusInput}>
                        {toChips(selectedEntries)}
                        <Input
                            margin="none"
                            value={currentValue}
                            inputRef={(ref) => (input.current = ref)}
                            onFocus={() => setOpen(true)}
                            onKeyDown={onKeyDown}
                            disableUnderline={true}
                            onChange={(e) => setCurrentValue(e.target.value)}
                            placeholder="Enter Tags"
                            style={{height: 40, minWidth: 150, flexGrow: 1}}
                        />
                    </div>
                </Tooltip>

                {open ? (
                    <Paper
                        style={{
                            position: 'absolute',
                            width: (container.current && container.current.clientWidth) || 300,
                            zIndex: 1000,
                        }}>
                        {suggestions.map((entry, index) => (
                            <Item key={label(entry)} entry={entry} onClick={trySubmit} selected={index === highlightedIndex} />
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
    return entries.map((entry) => <TagChip key={label(entry)} label={label(entry)} color={entry.tag.color} />);
};
