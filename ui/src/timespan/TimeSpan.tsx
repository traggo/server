import * as React from 'react';
import {TagSelectorEntry, toInputTags} from '../tag/tagSelectorEntry';
import {TagSelector} from '../tag/TagSelector';
import moment from 'moment';
import Paper from '@material-ui/core/Paper';
import {DateTimeSelector} from '../common/DateTimeSelector';
import {Button} from '@material-ui/core';
import {inUserTz, timeRunning} from './timeutils';
import {useMutation} from 'react-apollo-hooks';
import {StopTimer, StopTimerVariables} from '../gql/__generated__/StopTimer';
import * as gqlTimeSpan from '../gql/timeSpan';
import {UpdateTimeSpan, UpdateTimeSpanVariables} from '../gql/__generated__/UpdateTimeSpan';
import IconButton from '@material-ui/core/IconButton';
import {MoreVert} from '@material-ui/icons';
import Menu from '@material-ui/core/Menu';
import MenuItem from '@material-ui/core/MenuItem';
import {RemoveTimeSpan, RemoveTimeSpanVariables} from '../gql/__generated__/RemoveTimeSpan';
import {useStateAndDelegateWithDelayOnChange} from '../utils/hooks';
import {TimeSpans} from '../gql/__generated__/TimeSpans';
import {isSameDate} from '../utils/time';
import {Trackers} from '../gql/__generated__/Trackers';

interface Range {
    from: moment.Moment;
    to?: moment.Moment;
}

export interface TimeSpanProps {
    id: number;
    range: Range & {oldFrom?: moment.Moment};
    initialTags: TagSelectorEntry[];
    now?: moment.Moment;
    dateSelectorOpen?: React.Dispatch<React.SetStateAction<boolean>>;
    deleted?: () => void;
    stopped?: () => void;
    addTagsToTracker?: (tags: TagSelectorEntry[]) => void;
}

export const TimeSpan: React.FC<TimeSpanProps> = ({
    range: {from, to, oldFrom},
    id,
    initialTags,
    now,
    dateSelectorOpen = () => {},
    deleted = () => {},
    stopped = () => {},
    addTagsToTracker,
}) => {
    if (!to && !now) {
        throw new Error('now must be set when to is not set');
    }

    const [selectedEntries, setSelectedEntries] = React.useState<TagSelectorEntry[]>(initialTags);
    const [openMenu, setOpenMenu] = useStateAndDelegateWithDelayOnChange<null | HTMLElement>(null, (o) => dateSelectorOpen(!!o));
    const stopTimer = useMutation<StopTimer, StopTimerVariables>(gqlTimeSpan.StopTimer, {
        update: (cache, {data}) => {
            const oldTimeSpans = cache.readQuery<TimeSpans>({query: gqlTimeSpan.TimeSpans});
            const oldTrackers = cache.readQuery<Trackers>({query: gqlTimeSpan.Trackers});
            if (!oldTimeSpans || !oldTrackers || !data || !data.stopTimeSpan) {
                return;
            }
            cache.writeQuery<Trackers>({
                query: gqlTimeSpan.Trackers,
                data: {
                    timers: (oldTrackers.timers || []).filter((tracker) => tracker.id !== data.stopTimeSpan!.id),
                },
            });
            cache.writeQuery<TimeSpans>({
                query: gqlTimeSpan.TimeSpans,
                data: {
                    timeSpans: {
                        __typename: 'PagedTimeSpans',
                        timeSpans: oldTimeSpans.timeSpans.timeSpans
                            .concat([data.stopTimeSpan])
                            .sort((a, b) => moment(b.start).unix() - moment(a.start).unix()),
                        cursor: oldTimeSpans.timeSpans.cursor,
                    },
                },
            });
        },
    });
    const updateTimeSpan = useMutation<UpdateTimeSpan, UpdateTimeSpanVariables>(gqlTimeSpan.UpdateTimeSpan);
    const removeTimeSpan = useMutation<RemoveTimeSpan, RemoveTimeSpanVariables>(gqlTimeSpan.RemoveTimeSpan, {
        update: (cache, {data}) => {
            const oldData = cache.readQuery<TimeSpans>({query: gqlTimeSpan.TimeSpans});
            if (!oldData || !data || !data.removeTimeSpan) {
                return;
            }
            const removedId = data.removeTimeSpan.id;
            cache.writeQuery<TimeSpans>({
                query: gqlTimeSpan.TimeSpans,
                data: {
                    timeSpans: {
                        __typename: 'PagedTimeSpans',
                        timeSpans: oldData.timeSpans.timeSpans.filter((ts) => ts.id !== removedId),
                        cursor: oldData.timeSpans.cursor,
                    },
                },
            });
        },
    });

    const wasMoved = !isSameDate(from, oldFrom);
    const showDate = to !== undefined && (!isSameDate(from, to) || wasMoved);
    return (
        <Paper style={{display: 'flex', alignItems: 'center', padding: '10px', margin: '10px 0', opacity: wasMoved ? 0.5 : 1}}>
            <div style={{flex: '1', marginRight: 10}}>
                <TagSelector
                    dialogOpen={dateSelectorOpen}
                    selectedEntries={selectedEntries}
                    onSelectedEntriesChanged={(entries) => {
                        setSelectedEntries(entries);
                        updateTimeSpan({
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
                    if (to && moment(newFrom).isAfter(to)) {
                        const newTo = moment(newFrom).add(15, 'minute');
                        updateTimeSpan({
                            variables: {
                                oldStart: oldFrom,
                                id,
                                start: inUserTz(newFrom).format(),
                                end: inUserTz(newTo).format(),
                                tags: toInputTags(selectedEntries),
                            },
                        });
                    } else {
                        updateTimeSpan({
                            variables: {
                                id,
                                oldStart: oldFrom,
                                start: inUserTz(newFrom).format(),
                                end: to && inUserTz(to).format(),
                                tags: toInputTags(selectedEntries),
                            },
                        });
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
                        if (moment(newTo).isBefore(from)) {
                            const newFrom = moment(newTo).subtract(15, 'minute');
                            updateTimeSpan({
                                variables: {
                                    id,
                                    oldStart: oldFrom,
                                    start: inUserTz(newFrom).format(),
                                    end: inUserTz(newTo).format(),
                                    tags: toInputTags(selectedEntries),
                                },
                            });
                        } else {
                            updateTimeSpan({
                                variables: {
                                    id,
                                    oldStart: oldFrom,
                                    start: inUserTz(from).format(),
                                    end: inUserTz(newTo).format(),
                                    tags: toInputTags(selectedEntries),
                                },
                            });
                        }
                    }}
                    showDate={showDate}
                    label="end"
                />
            ) : (
                <Button
                    style={{minWidth: 120}}
                    onClick={() => {
                        stopTimer({variables: {id, end: inUserTz(require(now)).format()}}).then(stopped);
                    }}>
                    Stop {timeRunning(from, require(now))}
                </Button>
            )}
            <>
                <IconButton onClick={(e: React.MouseEvent<HTMLElement>) => setOpenMenu(e.currentTarget)}>
                    <MoreVert />
                </IconButton>
                <Menu aria-haspopup="true" anchorEl={openMenu} open={openMenu !== null} onClose={() => setOpenMenu(null)}>
                    <MenuItem
                        onClick={() => {
                            setOpenMenu(null);
                            removeTimeSpan({variables: {id}});
                            deleted();
                        }}>
                        Delete
                    </MenuItem>
                    {addTagsToTracker ? (
                        <MenuItem
                            onClick={() => {
                                setOpenMenu(null);
                                addTagsToTracker(selectedEntries);
                            }}>
                            Copy tags
                        </MenuItem>
                    ) : null}
                </Menu>
            </>
        </Paper>
    );
};

const require: <T>(e?: T) => T = (e) => {
    if (!e) {
        throw new Error('unset');
    }
    return e;
};
