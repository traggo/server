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
import ClickAwayListener from '@material-ui/core/ClickAwayListener';

export const calcShowDate = (from: moment.Moment, to?: moment.Moment): boolean => {
    const fromString = from.format('YYYYMMDD');
    return to !== undefined && fromString !== to.format('YYYYMMDD');
};

interface Range {
    from: moment.Moment;
    to?: moment.Moment;
}

export interface TimeSpanProps {
    id: number;
    range: Range;
    initialTags: TagSelectorEntry[];
    now?: moment.Moment;
    dateSelectorOpen?: (open: boolean) => void;
    deleted?: () => void;
}

export const TimeSpan: React.FC<TimeSpanProps> = ({
    range: {from, to},
    id,
    initialTags,
    now,
    dateSelectorOpen = () => {},
    deleted = () => {},
}) => {
    if (!to && !now) {
        throw new Error('now must be set when to is not set');
    }

    const [selectedEntries, setSelectedEntries] = React.useState<TagSelectorEntry[]>(initialTags);
    const [openMenu, setOpenMenuX] = React.useState<null | HTMLElement>(null);
    const setOpenMenu = (o: null | HTMLElement) => {
        setOpenMenuX(o);
        dateSelectorOpen(!!o);
    };
    const stopTimer = useMutation<StopTimer, StopTimerVariables>(gqlTimeSpan.StopTimer, {
        refetchQueries: [{query: gqlTimeSpan.Trackers}, {query: gqlTimeSpan.TimeSpans}],
    });
    const updateTimeSpan = useMutation<UpdateTimeSpan, UpdateTimeSpanVariables>(gqlTimeSpan.UpdateTimeSpan, {
        refetchQueries: [{query: gqlTimeSpan.Trackers}, {query: gqlTimeSpan.TimeSpans}],
    });
    const removeTimeSpan = useMutation<RemoveTimeSpan, RemoveTimeSpanVariables>(gqlTimeSpan.RemoveTimeSpan, {
        refetchQueries: [{query: gqlTimeSpan.Trackers}, {query: gqlTimeSpan.TimeSpans}],
    });

    const showDate = to !== undefined && calcShowDate(from, to);
    return (
        <Paper style={{display: 'flex', alignItems: 'center', padding: '10px', margin: '10px 0'}}>
            <div style={{flex: '1', marginRight: 10}}>
                <TagSelector
                    selectedEntries={selectedEntries}
                    onSelectedEntriesChanged={(entries) => {
                        setSelectedEntries(entries);
                        updateTimeSpan({
                            variables: {
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
                                    start: inUserTz(newFrom).format(),
                                    end: inUserTz(newTo).format(),
                                    tags: toInputTags(selectedEntries),
                                },
                            });
                        } else {
                            updateTimeSpan({
                                variables: {
                                    id,
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
                        stopTimer({variables: {id, end: inUserTz(require(now)).format()}});
                    }}>
                    Stop {timeRunning(from, require(now))}
                </Button>
            )}
            <ClickAwayListener onClickAway={() => setOpenMenu(null)}>
                <>
                    <IconButton onClick={(e: React.MouseEvent<HTMLElement>) => setOpenMenu(e.currentTarget)}>
                        <MoreVert />
                    </IconButton>
                    <Menu aria-haspopup="true" anchorEl={openMenu} open={openMenu !== null}>
                        <MenuItem
                            onClick={() => {
                                setOpenMenu(null);
                                removeTimeSpan({variables: {id}});
                                deleted();
                            }}>
                            Delete
                        </MenuItem>
                    </Menu>
                </>
            </ClickAwayListener>
        </Paper>
    );
};

const require: <T>(e?: T) => T = (e) => {
    if (!e) {
        throw new Error('unset');
    }
    return e;
};
