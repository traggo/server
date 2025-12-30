import * as React from 'react';
import {TagSelectorEntry} from '../tag/tagSelectorEntry';
import {TagSelector} from '../tag/TagSelector';
import moment from 'moment-timezone';
import {Button} from '@material-ui/core';
import {MoreVert} from '@material-ui/icons';
import IconButton from '@material-ui/core/IconButton';
import Paper from '@material-ui/core/Paper';
import Menu from '@material-ui/core/Menu';
import MenuItem from '@material-ui/core/MenuItem';
import ClickAwayListener from '@material-ui/core/ClickAwayListener';
import {DateTimeSelector} from '../common/DateTimeSelector';
import {useMutation} from '@apollo/react-hooks';
import * as gqlTimeSpan from '../gql/timeSpan';
import {StartTimer, StartTimerVariables} from '../gql/__generated__/StartTimer';
import {InputTimeSpanTag} from '../gql/__generated__/globalTypes';
import {AddTimeSpan, AddTimeSpanVariables} from '../gql/__generated__/AddTimeSpan';
import {useSnackbar} from 'notistack';
import {inUserTz} from './timeutils';
import {addTimeSpanToCache} from '../gql/utils';
import * as gqlStats from '../gql/statistics';

enum Type {
    Tracker,
    Manual,
}

export const calcShowDate = (from: moment.Moment, to: moment.Moment) => {
    const fromString = from.format('YYYYMMDD');
    return fromString !== to.format('YYYYMMDD') || moment().format('YYYYMMDD') !== fromString;
};

interface TrackerProps {
    onSelectedEntriesChanged: (entries: TagSelectorEntry[]) => void;
    selectedEntries: TagSelectorEntry[];
}

export const Tracker: React.FC<TrackerProps> = ({selectedEntries, onSelectedEntriesChanged: setSelectedEntries}) => {
    const [openMenu, setOpenMenu] = React.useState<null | HTMLElement>(null);
    const [type, setType] = React.useState<Type>(Type.Tracker);
    const [from, setFrom] = React.useState<moment.Moment>(moment().subtract(15, 'minute'));
    const [to, setTo] = React.useState<moment.Moment>(moment());
    const [showDate, setShowDate] = React.useState(false);
    const [hasInvalidRange, setHasInvalidRange] = React.useState(false);
    const [startTimer] = useMutation<StartTimer, StartTimerVariables>(gqlTimeSpan.StartTimer, {
        refetchQueries: [{query: gqlTimeSpan.Trackers}, {query: gqlStats.Stats2}],
    });
    const [addTimeSpan] = useMutation<AddTimeSpan, AddTimeSpanVariables>(gqlTimeSpan.AddTimeSpan, {
        refetchQueries: [{query: gqlStats.Stats2}],
        update: (cache, {data}) => {
            if (!data || !data.createTimeSpan) {
                return;
            }
            addTimeSpanToCache(cache, data.createTimeSpan);
        },
    });
    const {enqueueSnackbar} = useSnackbar();

    React.useEffect(() => {
        const handle = window.setInterval(() => {
            setShowDate(calcShowDate(from, to));
        }, 10000);
        return () => clearInterval(handle);
    }, [showDate, from, to]);

    const submit = () => {
        const tags = selectedEntries.map(
            (entry: TagSelectorEntry): InputTimeSpanTag => ({key: entry.tag.key, value: entry.value})
        );
        if (type === Type.Tracker) {
            startTimer({variables: {start: inUserTz(moment()).format(), tags, note: ''}}).then(() => {
                setSelectedEntries([]);
                enqueueSnackbar('tracker started', {variant: 'success'});
            });
        } else {
            addTimeSpan({variables: {start: inUserTz(from).format(), end: inUserTz(to).format(), tags, note: ''}}).then(() => {
                setSelectedEntries([]);
                enqueueSnackbar('time span added', {variant: 'success'});
            });
        }
    };

    return (
        <ClickAwayListener onClickAway={() => setOpenMenu(null)}>
            <Paper style={{display: 'flex', alignItems: 'center', padding: '10px'}}>
                <div style={{flex: '1', marginRight: 10}}>
                    <TagSelector
                        selectedEntries={selectedEntries}
                        onSelectedEntriesChanged={setSelectedEntries}
                        onCtrlEnter={submit}
                    />
                </div>
                {type === Type.Manual ? (
                    <div>
                        <DateTimeSelector
                            selectedDate={from}
                            onSelectDate={(newFrom) => {
                                if (!newFrom.isValid()) {
                                    return;
                                }
                                setFrom(newFrom);
                                setShowDate(calcShowDate(newFrom, to));
                                // Check if new range would be invalid
                                if (moment(newFrom).isAfter(to)) {
                                    setHasInvalidRange(true);
                                } else {
                                    setHasInvalidRange(false);
                                }
                            }}
                            showDate={showDate}
                            label="start"
                        />
                        <DateTimeSelector
                            selectedDate={to}
                            onSelectDate={(newTo) => {
                                if (!newTo.isValid()) {
                                    return;
                                }
                                setTo(newTo);
                                setShowDate(calcShowDate(from, newTo));
                                // Check if new range would be invalid
                                if (moment(newTo).isBefore(from)) {
                                    setHasInvalidRange(true);
                                } else {
                                    setHasInvalidRange(false);
                                }
                            }}
                            showDate={showDate}
                            label="end"
                        />
                    </div>
                ) : null}
                {hasInvalidRange && type === Type.Manual && (
                    <div style={{color: '#f44336', fontSize: '0.9rem', marginRight: '8px', display: 'flex', alignItems: 'center'}}>
                        <span title="Warning: End time is before start time">⚠️</span>
                    </div>
                )}
                <Button variant="text" style={{height: 50}} onClick={submit} disabled={hasInvalidRange}>
                    {type === Type.Manual ? 'add' : 'start'}
                </Button>
                <IconButton onClick={(e: React.MouseEvent<HTMLElement>) => setOpenMenu(e.currentTarget)}>
                    <MoreVert />
                </IconButton>
                <Menu aria-haspopup="true" anchorEl={openMenu} open={openMenu !== null}>
                    <MenuItem
                        selected={type === Type.Tracker}
                        onClick={() => {
                            setOpenMenu(null);
                            setType(Type.Tracker);
                        }}>
                        Tracker
                    </MenuItem>
                    <MenuItem
                        selected={type === Type.Manual}
                        onClick={() => {
                            setOpenMenu(null);
                            setType(Type.Manual);
                        }}>
                        Manual
                    </MenuItem>
                </Menu>
            </Paper>
        </ClickAwayListener>
    );
};
