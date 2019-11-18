import * as React from 'react';
import BigCalendar, {DateRangeFormatFunction, View} from 'react-big-calendar';
import moment from 'moment';
import {useMutation, useQuery} from '@apollo/react-hooks';
import * as gqlTimeSpan from '../../gql/timeSpan';
import * as gqlTag from '../../gql/tags';
import {TimeSpans} from '../../gql/__generated__/TimeSpans';
import {Tags} from '../../gql/__generated__/Tags';
import useInterval from '@rooks/use-interval';
import {Paper} from '@material-ui/core';
import Popper from '@material-ui/core/Popper';
import {TimeSpan} from '../TimeSpan';
import {InputTag, toTagSelectorEntry} from '../../tag/tagSelectorEntry';
import ClickAwayListener from '@material-ui/core/ClickAwayListener';
import Typography from '@material-ui/core/Typography';
import {AddTimeSpan, AddTimeSpanVariables} from '../../gql/__generated__/AddTimeSpan';
import withDragAndDrop from 'react-big-calendar/lib/addons/dragAndDrop';
import {UpdateTimeSpan, UpdateTimeSpanVariables} from '../../gql/__generated__/UpdateTimeSpan';
import {Trackers} from '../../gql/__generated__/Trackers';

const localizer = BigCalendar.momentLocalizer(moment);
const DADCalendar = withDragAndDrop(BigCalendar);

const selectFormat: DateRangeFormatFunction = (range) => `${moment(range.start)} - ${moment(range.end)}`;

interface CalendarEntry {
    id: number;
    start: Date;
    end: Date;
    hasEnd: boolean;
    title: string;
    tags: InputTag[];
}

export const Calendar = () => {
    const timeSpansResult = useQuery<TimeSpans>(gqlTimeSpan.TimeSpans);
    const trackersResult = useQuery<Trackers>(gqlTimeSpan.Trackers);
    const tagsResult = useQuery<Tags>(gqlTag.Tags);
    const [currentDate, setCurrentDate] = React.useState(moment());
    const [addTimeSpan] = useMutation<AddTimeSpan, AddTimeSpanVariables>(gqlTimeSpan.AddTimeSpan, {
        refetchQueries: [{query: gqlTimeSpan.TimeSpans}],
    });
    const [updateTimeSpanMutation] = useMutation<UpdateTimeSpan, UpdateTimeSpanVariables>(gqlTimeSpan.UpdateTimeSpan);
    useInterval(
        () => {
            setCurrentDate(moment());
        },
        5000,
        true
    );

    const values: CalendarEntry[] = React.useMemo(() => {
        if (
            timeSpansResult.error ||
            timeSpansResult.loading ||
            !timeSpansResult.data ||
            timeSpansResult.data.timeSpans === null ||
            trackersResult.error ||
            trackersResult.loading ||
            !trackersResult.data ||
            trackersResult.data.timers === null ||
            tagsResult.error ||
            tagsResult.loading ||
            !tagsResult.data ||
            tagsResult.data.tags === null
        ) {
            return [];
        }
        return timeSpansResult.data.timeSpans.timeSpans
            .concat(trackersResult.data.timers)
            .sort((a, b) => a.start.toString().localeCompare(b.start.toString()))
            .map((ts) => ({
                start: moment(ts.start).toDate(),
                end: (ts.end && moment(ts.end).toDate()) || currentDate.toDate(),
                hasEnd: !!ts.end,
                id: ts.id,
                tags: ts.tags!,
                title: ts.tags!.map((t) => t.key + ':' + t.stringValue).join(' '),
            }));
    }, [trackersResult, tagsResult, timeSpansResult, currentDate]);
    const [selected, setSelected] = React.useState<{selected: HTMLElement | null; data: CalendarEntry | null}>({
        selected: null,
        data: null,
    });
    const [ignore, setIgnore] = React.useState<boolean>(false);
    const [view, setView] = React.useState<View>('week');

    const updateTimeSpan = (data: {event: CalendarEntry; start: Date | string; end: Date | string}) => {
        updateTimeSpanMutation({
            variables: {
                id: data.event.id,
                start: moment(data.start).format(),
                end: moment(data.end).format(),
                tags: data.event.tags.map((tag) => ({key: tag.key, stringValue: tag.stringValue})),
            },
        });
    };

    return (
        <div style={{height: 'calc(100% - 64px', minHeight: 2000, margin: '0 auto', maxWidth: view === 'week' ? '100%' : 600}}>
            <Paper style={{height: '100%', padding: 10}}>
                <Typography component="div" style={{height: '100%'}}>
                    <DADCalendar<CalendarEntry>
                        elementProps={{style: {height: '100%'}}}
                        localizer={localizer}
                        events={values}
                        onView={setView}
                        view={view}
                        startAccessor="start"
                        endAccessor="end"
                        views={['week', 'day']}
                        step={10}
                        onEventResize={updateTimeSpan}
                        onEventDrop={updateTimeSpan}
                        onSelectSlot={(data) => {
                            if (data.action === 'click') {
                                return false;
                            }
                            addTimeSpan({
                                variables: {start: moment(data.start).format(), end: moment(data.end).format(), tags: []},
                            });
                            return true;
                        }}
                        timeslots={6}
                        selectable={true}
                        onSelectEvent={(data, event) => setSelected({selected: event.currentTarget, data})}
                        defaultView={'week'}
                        formats={{selectRangeFormat: selectFormat, timeGutterFormat: 'HH:mm'}}
                    />
                </Typography>
                {!!selected.selected && (
                    <Popper open={true} anchorEl={selected.selected} style={{zIndex: 1200, maxWidth: 700}}>
                        <ClickAwayListener
                            onClickAway={() => {
                                if (ignore) {
                                    return;
                                }
                                setSelected({selected: null, data: null});
                            }}>
                            <TimeSpan
                                id={selected.data!.id}
                                deleted={() => setSelected({selected: null, data: null})}
                                range={{
                                    from: moment(selected.data!.start),
                                    to: selected.data!.hasEnd ? moment(selected.data!.end) : undefined,
                                }}
                                initialTags={toTagSelectorEntry(tagsResult.data!.tags!, selected.data!.tags)}
                                dateSelectorOpen={setIgnore}
                                stopped={() => {
                                    setSelected({
                                        ...selected,
                                        data: {...selected.data!, hasEnd: true, end: currentDate.toDate()},
                                    });
                                }}
                            />
                        </ClickAwayListener>
                    </Popper>
                )}
            </Paper>
        </div>
    );
};
