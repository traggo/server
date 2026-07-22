import * as React from 'react';
import {TagSelectorEntry, toInputTags} from '../tag/tagSelectorEntry';
import {TagSelector} from '../tag/TagSelector';
import moment from 'moment';
import Paper from '@mui/material/Paper';
import {DateTimeSelector} from '../common/DateTimeSelector';
import {Button, TextField, Typography} from '@mui/material';
import makeStyles from '@mui/styles/makeStyles';
import {inUserTz} from './timeutils';
import {useMutation} from '@apollo/client';
import {
    StopTimerMutation,
    StopTimerMutationVariables,
    UpdateTimeSpanMutation,
    UpdateTimeSpanMutationVariables,
    RemoveTimeSpanMutation,
    RemoveTimeSpanMutationVariables,
    TimeSpansQuery,
    TrackersQuery,
    StartTimerMutation,
    StartTimerMutationVariables,
} from '../gql/__generated__';
import * as gqlTimeSpan from '../gql/timeSpan';
import IconButton from '@mui/material/IconButton';
import {MoreVert} from '@mui/icons-material';
import Menu from '@mui/material/Menu';
import MenuItem from '@mui/material/MenuItem';
import {useStateAndDelegateWithDelayOnChange} from '../utils/hooks';
import {isSameDate} from '../utils/time';
import {addTimeSpanToCache, removeFromTrackersCache} from '../gql/utils';
import {RelativeTime, RelativeToNow} from '../common/RelativeTime';

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

const useStyles = makeStyles(() => ({
    innerTimespan: {
        position: 'relative',
        width: '100%',
        display: 'flex',
        alignItems: 'center',
        textAlign: 'center',
        '@media (max-width: 750px)': {
            flexDirection: 'column',
        },
    },
    tagInput: {
        width: '100%',
        flex: '1',
        marginRight: 10,
        '@media (max-width: 750px)': {
            display: 'flex',
            width: 'calc(100% - 48px)',
            marginRight: '48px',
        },
    },
    timeSelection: {
        display: 'inline-flex',
        // flex-end (not center): the start/end pickers have a floating label above their value,
        // making them taller than the plain duration text - centering both as whole boxes leaves
        // the duration vertically centered in the row while the picker's actual value text (below
        // its label) sits lower, looking misaligned. Aligning bottoms lines up the value baselines.
        alignItems: 'flex-end',
        '@media (max-width: 750px)': {
            justifyContent: 'space-evenly',
            width: '100%',
            flexWrap: 'wrap',
        },
    },
    showMoreButton: {
        display: 'flex',
        alignItems: 'center',
        '@media (max-width: 750px)': {
            position: 'absolute',
            top: '0',
            right: '0',
        },
    },
}));

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
        const styles = useStyles();
        const [showNotes, toggleShowingNotes] = React.useState(initialNote !== '');
        const note = React.useRef<{value: string; handle?: number}>({value: initialNote});

        const [selectedEntries, setSelectedEntries] = React.useState<TagSelectorEntry[]>(initialTags);
        const [openMenu, setOpenMenu] = useStateAndDelegateWithDelayOnChange<null | HTMLElement>(null, (o) =>
            dateSelectorOpen(!!o)
        );
        const [stopTimer] = useMutation<StopTimerMutation, StopTimerMutationVariables>(gqlTimeSpan.StopTimer, {
            update: (cache, {data}) => {
                if (!data || !data.stopTimeSpan) {
                    return;
                }
                removeFromTrackersCache(cache, data);
                addTimeSpanToCache(cache, data.stopTimeSpan);
            },
        });
        const [startTimer] = useMutation<StartTimerMutation, StartTimerMutationVariables>(gqlTimeSpan.StartTimer, {
            refetchQueries: [{query: gqlTimeSpan.Trackers}],
        });
        const [updateTimeSpan] = useMutation<UpdateTimeSpanMutation, UpdateTimeSpanMutationVariables>(gqlTimeSpan.UpdateTimeSpan);
        const noteAwareUpdateTimeSpan = ({variables}: {variables: Omit<UpdateTimeSpanMutationVariables, 'note'>}) => {
            clearTimeout(note.current.handle);
            return updateTimeSpan({variables: {...variables, note: note.current.value}});
        };
        const [removeTimeSpan] = useMutation<RemoveTimeSpanMutation, RemoveTimeSpanMutationVariables>(
            gqlTimeSpan.RemoveTimeSpan,
            {
                update: (cache, {data}) => {
                    let oldData: TimeSpansQuery | null = null;
                    try {
                        oldData = cache.readQuery<TimeSpansQuery>({query: gqlTimeSpan.TimeSpans});
                    } catch {}

                    const oldTrackers = cache.readQuery<TrackersQuery>({query: gqlTimeSpan.Trackers});
                    if (!data || !data.removeTimeSpan) {
                        return;
                    }
                    const removedId = data.removeTimeSpan.id;
                    if (oldTrackers) {
                        cache.writeQuery<TrackersQuery>({
                            query: gqlTimeSpan.Trackers,
                            data: {
                                timers: (oldTrackers.timers || []).filter((tracker) => tracker.id !== removedId),
                            },
                        });
                    }
                    if (oldData) {
                        cache.writeQuery<TimeSpansQuery>({
                            query: gqlTimeSpan.TimeSpans,
                            data: {
                                timeSpans: {
                                    __typename: 'PagedTimeSpans',
                                    timeSpans: oldData.timeSpans.timeSpans.filter((ts) => ts.id !== removedId),
                                    cursor: oldData.timeSpans.cursor,
                                },
                            },
                        });
                    }
                },
            }
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
                    width: '100%',
                }}>
                <div className={styles.innerTimespan}>
                    <div className={styles.tagInput}>
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

                    <div className={styles.timeSelection}>
                        <div style={{alignItems: 'center', display: 'flex', justifyContent: 'space-evenly', flexWrap: 'wrap'}}>
                            <DateTimeSelector
                                popoverOpen={dateSelectorOpen}
                                selectedDate={from}
                                onSelectDate={(newFrom) => {
                                    if (!newFrom.isValid()) {
                                        return;
                                    }
                                    newFrom.set({second: 0});
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
                                        newTo.set({second: 0});
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
                        </div>

                        <div style={{alignItems: 'center', display: 'flex'}}>
                            <Typography
                                variant="subtitle1"
                                // Matches the start/end pickers' MUI standard-variant input, which has
                                // padding: 4px 0 5px below its own text - without this, aligning by
                                // flex-end lines up the two elements' boxes but not their actual text,
                                // since this Typography has no padding of its own to account for.
                                style={{minWidth: '70px', paddingBottom: '5px'}}
                                title="The amount of time between from and to">
                                {to ? <RelativeTime from={from} to={to} /> : <RelativeToNow from={from} />}
                            </Typography>
                        </div>
                    </div>

                    <IconButton
                        className={styles.showMoreButton}
                        onClick={(e: React.MouseEvent<HTMLElement>) => setOpenMenu(e.currentTarget)}
                        size="large">
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
                                toggleShowingNotes(!showNotes);
                            }}>
                            Show Notes
                        </MenuItem>
                        <MenuItem
                            onClick={() => {
                                setOpenMenu(null);
                                removeTimeSpan({variables: {id}}).then(() => deleted());
                            }}>
                            Delete
                        </MenuItem>
                    </Menu>
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
