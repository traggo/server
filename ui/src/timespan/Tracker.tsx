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
import {useMutation} from 'react-apollo-hooks';
import * as gqlTimeSpan from '../gql/timeSpan';
import {StartTimer, StartTimerVariables} from '../gql/__generated__/StartTimer';
import {InputTimeSpanTag} from '../gql/__generated__/globalTypes';
import {AddTimeSpan, AddTimeSpanVariables} from '../gql/__generated__/AddTimeSpan';
import {useSnackbar} from 'notistack';
import {inUserTz} from './timeutils';

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
    const [manualSelected, setManualSelected] = React.useState(false);
    const [from, setFrom] = React.useState<moment.Moment>(moment());
    const [to, setTo] = React.useState<moment.Moment>(moment().add(15, 'minute'));
    const [showDate, setShowDate] = React.useState(false);
    const startTimer = useMutation<StartTimer, StartTimerVariables>(gqlTimeSpan.StartTimer, {
        refetchQueries: [{query: gqlTimeSpan.Trackers}],
    });
    const addTimeSpan = useMutation<AddTimeSpan, AddTimeSpanVariables>(gqlTimeSpan.AddTimeSpan, {
        refetchQueries: [{query: gqlTimeSpan.TimeSpans}],
    });
    const {enqueueSnackbar} = useSnackbar();

    React.useEffect(() => {
        if (manualSelected) {
            return;
        }
        const handle = window.setInterval(() => {
            const newFrom = moment();
            setFrom(newFrom);
            const newTo = moment().add(15, 'minute');
            setTo(newTo);
            setShowDate(calcShowDate(newFrom, newTo));
        }, 10000);
        return () => clearInterval(handle);
    }, [manualSelected]);

    React.useEffect(() => {
        const handle = window.setInterval(() => {
            setShowDate(calcShowDate(from, to));
        }, 10000);
        return () => clearInterval(handle);
    }, [showDate, from, to]);

    const submit = () => {
        const tags = selectedEntries.map(
            (entry: TagSelectorEntry): InputTimeSpanTag => ({key: entry.tag.key, stringValue: entry.value})
        );
        if (type === Type.Tracker) {
            startTimer({variables: {start: inUserTz(moment()).format(), tags}}).then(() => {
                setSelectedEntries([]);
                enqueueSnackbar('tracker started', {variant: 'success'});
            });
        } else {
            addTimeSpan({variables: {start: inUserTz(from).format(), end: inUserTz(to).format(), tags}}).then(() => {
                setSelectedEntries([]);
                enqueueSnackbar('time span added', {variant: 'success'});
            });
        }
    };

    return (
        <ClickAwayListener onClickAway={() => setOpenMenu(null)}>
            <Paper style={{display: 'flex', alignItems: 'center', padding: '10px'}}>
                <div style={{flex: '1', marginRight: 10}}>
                    <TagSelector selectedEntries={selectedEntries} onSelectedEntriesChanged={setSelectedEntries} />
                </div>
                {type === Type.Manual ? (
                    <div>
                        <DateTimeSelector
                            selectedDate={from}
                            onSelectDate={(newFrom) => {
                                setManualSelected(true);
                                setFrom(newFrom);
                                if (moment(newFrom).isAfter(to)) {
                                    const newTo = moment(newFrom).add(15, 'minute');
                                    setTo(newTo);
                                    setShowDate(calcShowDate(newFrom, newTo));
                                } else {
                                    setShowDate(calcShowDate(newFrom, to));
                                }
                            }}
                            showDate={showDate}
                            label="start"
                        />
                        <DateTimeSelector
                            selectedDate={to}
                            onSelectDate={(newTo) => {
                                setManualSelected(true);
                                setTo(newTo);
                                if (moment(newTo).isBefore(from)) {
                                    const newFrom = moment(newTo).subtract(15, 'minute');
                                    setFrom(newFrom);
                                    setShowDate(calcShowDate(newFrom, newTo));
                                } else {
                                    setShowDate(calcShowDate(from, newTo));
                                }
                            }}
                            showDate={showDate}
                            label="end"
                        />
                    </div>
                ) : null}
                <Button variant="text" style={{height: 50}} onClick={submit}>
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
