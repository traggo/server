import * as React from 'react';
import {TagSelectorEntry, toInputTags} from '../tag/tagSelectorEntry';
import {TagSelector} from '../tag/TagSelector';
import moment from 'moment';
import Paper from '@material-ui/core/Paper';
import {DateTimeSelector} from '../common/DateTimeSelector';
import {Button, TextField, Tooltip, Typography} from '@material-ui/core';
import {inUserTz} from './timeutils';
import {useMutation} from '@apollo/react-hooks';
import {StopTimer, StopTimerVariables} from '../gql/__generated__/StopTimer';
import * as gqlTimeSpan from '../gql/timeSpan';
import {UpdateTimeSpan, UpdateTimeSpanVariables} from '../gql/__generated__/UpdateTimeSpan';
import IconButton from '@material-ui/core/IconButton';
import {MoreVert} from '@material-ui/icons';
import Menu from '@material-ui/core/Menu';
import MenuItem from '@material-ui/core/MenuItem';
import {RemoveTimeSpan, RemoveTimeSpanVariables} from '../gql/__generated__/RemoveTimeSpan';
import {useStateAndDelegateWithDelayOnChange} from '../utils/hooks';
import {isSameDate} from '../utils/time';
import {addTimeSpanToCache, removeFromTrackersCache} from '../gql/utils';
import {StartTimer, StartTimerVariables} from '../gql/__generated__/StartTimer';
import {RelativeTime, RelativeToNow} from '../common/RelativeTime';
import ShowNotesIcon from '@material-ui/icons/KeyboardArrowDown';
import HideNotesIcon from '@material-ui/icons/KeyboardArrowUp';
import {removeTimeSpanOptions} from './mutations';

interface Range {
    from: moment.Moment;
    to?: moment.Moment;
}

export interface TimeSpanProps {
    id: number;
    range: Range & {oldFrom?: moment.Moment};
    initialTags: TagSelectorEntry[];
    note: string;
    dateSelectorOpen?: React.Dispatch<React.SetStateAction<boolean>>;
    rangeChange?: (r: Range) => void;
    deleted?: () => void;
    stopped?: () => void;
    continued?: () => void;
    addTagsToTracker?: (tags: TagSelectorEntry[]) => void;
    elevation?: number;
}

export const TimeSpan: React.FC<TimeSpanProps> = React.memo(
    ({
        range: {from, to, oldFrom},
        id,
        initialTags,
        note: initialNote,
        dateSelectorOpen = () => {},
        rangeChange = () => {},
        deleted = () => {},
        stopped = () => {},
        continued = () => {},
        elevation = 1,
        addTagsToTracker,
    }) => {
        const [showNotes, toggleShowingNotes] = React.useState(initialNote !== '');
        const note = React.useRef<{value: string; handle?: number}>({value: initialNote});

        const [selectedEntries, setSelectedEntries] = React.useState<TagSelectorEntry[]>(initialTags);
        const [openMenu, setOpenMenu] = useStateAndDelegateWithDelayOnChange<null | HTMLElement>(null, (o) =>
            dateSelectorOpen(!!o)
        );
        const [stopTimer] = useMutation<StopTimer, StopTimerVariables>(gqlTimeSpan.StopTimer, {
            update: (cache, {data}) => {
                if (!data || !data.stopTimeSpan) {
                    return;
                }
                removeFromTrackersCache(cache, data);
                addTimeSpanToCache(cache, data.stopTimeSpan);
            },
        });
        const [startTimer] = useMutation<StartTimer, StartTimerVariables>(gqlTimeSpan.StartTimer, {
            refetchQueries: [{query: gqlTimeSpan.Trackers}],
        });
        const [updateTimeSpan] = useMutation<UpdateTimeSpan, UpdateTimeSpanVariables>(gqlTimeSpan.UpdateTimeSpan);
        const noteAwareUpdateTimeSpan = ({variables}: {variables: Omit<UpdateTimeSpanVariables, 'note'>}) => {
            clearTimeout(note.current.handle);
            return updateTimeSpan({variables: {...variables, note: note.current.value}});
        };
        const [removeTimeSpan] = useMutation<RemoveTimeSpan, RemoveTimeSpanVariables>(
            gqlTimeSpan.RemoveTimeSpan,
            removeTimeSpanOptions
        );

        const updateNote = (newValue: string) => {
            window.clearTimeout(note.current.handle);
            const handle = window.setTimeout(
                () =>
                    updateTimeSpan({
                        variables: {
                            oldStart: oldFrom,
                            id,
                            start: inUserTz(from).format(),
                            end: to && inUserTz(to).format(),
                            tags: toInputTags(selectedEntries),
                            note: newValue,
                        },
                    }),
                200
            );
            note.current = {handle, value: newValue};
        };

        const wasMoved = !isSameDate(from, oldFrom);
        const showDate = to !== undefined && (!isSameDate(from, to) || wasMoved);
        return (
            <Paper
                elevation={elevation}
                style={{
                    display: 'flex',
                    flexDirection: 'column',
                    padding: '10px',
                    margin: '10px 0',
                    opacity: wasMoved ? 0.5 : 1,
                }}>
                <div
                    style={{
                        display: 'flex',
                        alignItems: 'center',
                    }}>
                    <Tooltip title="Toggle notes">
                        <IconButton onClick={() => toggleShowingNotes(!showNotes)}>
                            {showNotes ? <HideNotesIcon /> : <ShowNotesIcon />}
                        </IconButton>
                    </Tooltip>
                    <div style={{flex: '1', marginRight: 10}}>
                        <TagSelector
                            dialogOpen={dateSelectorOpen}
                            selectedEntries={selectedEntries}
                            onSelectedEntriesChanged={(entries) => {
                                setSelectedEntries(entries);
                                noteAwareUpdateTimeSpan({
                                    variables: {
                                        oldStart: oldFrom,
                                        id,
                                        start: inUserTz(from).format(),
                                        end: to && inUserTz(to).format(),
                                        tags: toInputTags(entries),
                                    },
                                });
                            }}
                        />
                    </div>
                    <DateTimeSelector
                        popoverOpen={dateSelectorOpen}
                        selectedDate={from}
                        onSelectDate={(newFrom) => {
                            if (!newFrom.isValid()) {
                                return;
                            }
                            if (to && moment(newFrom).isAfter(to)) {
                                const newTo = moment(newFrom).add(15, 'minute');
                                noteAwareUpdateTimeSpan({
                                    variables: {
                                        oldStart: oldFrom,
                                        id,
                                        start: inUserTz(newFrom).format(),
                                        end: inUserTz(newTo).format(),
                                        tags: toInputTags(selectedEntries),
                                    },
                                }).then(() => rangeChange({from: newFrom, to: newTo}));
                            } else {
                                noteAwareUpdateTimeSpan({
                                    variables: {
                                        id,
                                        oldStart: oldFrom,
                                        start: inUserTz(newFrom).format(),
                                        end: to && inUserTz(to).format(),
                                        tags: toInputTags(selectedEntries),
                                    },
                                }).then(() => rangeChange({from: newFrom, to}));
                            }
                        }}
                        showDate={showDate}
                        label="start"
                    />
                    {to !== undefined ? (
                        <DateTimeSelector
                            popoverOpen={dateSelectorOpen}
                            selectedDate={to}
                            onSelectDate={(newTo) => {
                                if (!newTo.isValid()) {
                                    return;
                                }
                                if (moment(newTo).isBefore(from)) {
                                    const newFrom = moment(newTo).subtract(15, 'minute');
                                    noteAwareUpdateTimeSpan({
                                        variables: {
                                            id,
                                            oldStart: oldFrom,
                                            start: inUserTz(newFrom).format(),
                                            end: inUserTz(newTo).format(),
                                            tags: toInputTags(selectedEntries),
                                        },
                                    }).then(() => rangeChange({from: newFrom, to: newTo}));
                                } else {
                                    noteAwareUpdateTimeSpan({
                                        variables: {
                                            id,
                                            oldStart: oldFrom,
                                            start: inUserTz(from).format(),
                                            end: inUserTz(newTo).format(),
                                            tags: toInputTags(selectedEntries),
                                        },
                                    }).then(() => rangeChange({from, to: newTo}));
                                }
                            }}
                            showDate={showDate}
                            label="end"
                        />
                    ) : (
                        <Button
                            onClick={() => {
                                stopTimer({variables: {id, end: inUserTz(moment()).format()}}).then(stopped);
                            }}>
                            Stop
                        </Button>
                    )}
                    <>
                        {
                            <Typography
                                variant="subtitle1"
                                style={{width: 70, textAlign: 'right'}}
                                title="The amount of time between from and to">
                                {to ? <RelativeTime from={from} to={to} /> : <RelativeToNow from={from} />}
                            </Typography>
                        }
                        <IconButton onClick={(e: React.MouseEvent<HTMLElement>) => setOpenMenu(e.currentTarget)}>
                            <MoreVert />
                        </IconButton>
                        <Menu aria-haspopup="true" anchorEl={openMenu} open={openMenu !== null} onClose={() => setOpenMenu(null)}>
                            {to ? (
                                <MenuItem
                                    onClick={() => {
                                        setOpenMenu(null);
                                        startTimer({
                                            variables: {
                                                start: inUserTz(moment()).format(),
                                                tags: toInputTags(selectedEntries),
                                                note: note.current.value,
                                            },
                                        }).then(() => continued());
                                    }}>
                                    Continue
                                </MenuItem>
                            ) : null}
                            {addTagsToTracker ? (
                                <MenuItem
                                    onClick={() => {
                                        setOpenMenu(null);
                                        addTagsToTracker(selectedEntries);
                                    }}>
                                    Copy tags
                                </MenuItem>
                            ) : null}
                            <MenuItem
                                onClick={() => {
                                    setOpenMenu(null);
                                    removeTimeSpan({variables: {id}}).then(() => deleted());
                                }}>
                                Delete
                            </MenuItem>
                        </Menu>
                    </>
                </div>
                {showNotes ? (
                    <div>
                        <TextField
                            label="Note"
                            fullWidth
                            multiline
                            defaultValue={initialNote}
                            onChange={(e) => updateNote(e.target.value)}
                        />
                    </div>
                ) : null}
            </Paper>
        );
    }
);
