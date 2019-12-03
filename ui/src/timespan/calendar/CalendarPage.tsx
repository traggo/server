import * as React from 'react';
import * as ReactDOM from 'react-dom';
import {Button, makeStyles, Paper, useTheme} from '@material-ui/core';
import moment from 'moment';
import {useApolloClient, useMutation, useQuery} from '@apollo/react-hooks';
import {TimeSpans_timeSpans_timeSpans} from '../../gql/__generated__/TimeSpans';
import * as gqlTimeSpan from '../../gql/timeSpan';
import {Trackers} from '../../gql/__generated__/Trackers';
import {Tags} from '../../gql/__generated__/Tags';
import * as gqlTag from '../../gql/tags';
import FullCalendar from '@fullcalendar/react';
import {calculateColor, ColorMode} from '../colorutils';
import '@fullcalendar/core/main.css';
import '@fullcalendar/daygrid/main.css';
import '@fullcalendar/timegrid/main.css';
import dayGridPlugin from '@fullcalendar/daygrid';
import timeGridPlugin from '@fullcalendar/timegrid';
import momentPlugin from '@fullcalendar/moment';
import interactionPlugin from '@fullcalendar/interaction';
import {OptionsInput} from '@fullcalendar/core';
import {UpdateTimeSpan, UpdateTimeSpanVariables} from '../../gql/__generated__/UpdateTimeSpan';
import Popper from '@material-ui/core/Popper';
import ClickAwayListener from '@material-ui/core/ClickAwayListener';
import {TimeSpan} from '../TimeSpan';
import {toTagSelectorEntry} from '../../tag/tagSelectorEntry';
import {AddTimeSpan, AddTimeSpanVariables} from '../../gql/__generated__/AddTimeSpan';
import {FullCalendarStyling} from './FullCalendarStyling';
import useInterval from '@rooks/use-interval';
import {EventApi} from '@fullcalendar/core/api/EventApi';
import {StopTimer, StopTimerVariables} from '../../gql/__generated__/StopTimer';
import {
    addTimeSpanInRangeToCache,
    addTimeSpanToCache,
    removeFromTimeSpanInRangeCache,
    removeFromTrackersCache,
} from '../../gql/utils';
import {StartTimer, StartTimerVariables} from '../../gql/__generated__/StartTimer';
import {timeRunningCalendar} from '../timeutils';
import {stripTypename} from '../../utils/strip';
import {TimeSpansInRange, TimeSpansInRangeVariables} from '../../gql/__generated__/TimeSpansInRange';

const toMoment = (date: Date): moment.Moment => {
    return moment(date).tz('utc');
};

declare global {
    interface Window {
        // tslint:disable-next-line:no-any
        __TRAGGO_CALENDAR: any;
    }
}

const useStyles = makeStyles((theme) => {
    return {
        start: {
            width: '50%',
            float: 'right',
            padding: '2px 0',
            right: 0,
            borderRadius: '0 0 10% 10px',
            textAlign: 'center',
            position: 'absolute',
            boxShadow: theme.shadows[1],
            borderTop: '1px solid red',
        },
    };
});

export const CalendarPage: React.FC = () => {
    const apollo = useApolloClient();
    const classes = useStyles();
    const theme = useTheme();
    const timeSpansResult = useQuery<TimeSpansInRange, TimeSpansInRangeVariables>(gqlTimeSpan.TimeSpansInRange, {
        variables: {
            start: moment()
                .startOf('week')
                .format(),
            end: moment()
                .endOf('week')
                .format(),
        },
        fetchPolicy: 'cache-and-network',
    });
    const trackersResult = useQuery<Trackers>(gqlTimeSpan.Trackers, {fetchPolicy: 'cache-and-network'});
    const tagsResult = useQuery<Tags>(gqlTag.Tags);
    const [startTimer] = useMutation<StartTimer, StartTimerVariables>(gqlTimeSpan.StartTimer, {
        refetchQueries: [{query: gqlTimeSpan.Trackers}],
    });
    const [updateTimeSpanMutation] = useMutation<UpdateTimeSpan, UpdateTimeSpanVariables>(gqlTimeSpan.UpdateTimeSpan);
    const [currentDate, setCurrentDate] = React.useState(moment());
    const [stopTimer] = useMutation<StopTimer, StopTimerVariables>(gqlTimeSpan.StopTimer, {
        update: (cache, {data}) => {
            if (!data || !data.stopTimeSpan) {
                return;
            }
            removeFromTrackersCache(cache, data);
            addTimeSpanInRangeToCache(cache, data.stopTimeSpan, timeSpansResult.variables);
        },
    });
    useInterval(
        () => {
            setCurrentDate(moment());
        },
        60000,
        true
    );
    React.useEffect(() => {
        window.__TRAGGO_CALENDAR = {};
        return () => (window.__TRAGGO_CALENDAR = undefined);
    });
    const [ignore, setIgnore] = React.useState<boolean>(false);
    const [selected, setSelected] = React.useState<{selected: HTMLElement | null; data: TimeSpans_timeSpans_timeSpans | null}>({
        selected: null,
        data: null,
    });
    const [addTimeSpan] = useMutation<AddTimeSpan, AddTimeSpanVariables>(gqlTimeSpan.AddTimeSpan, {
        update: (cache, {data}) => {
            if (!data || !data.createTimeSpan) {
                return;
            }
            addTimeSpanInRangeToCache(cache, data.createTimeSpan, timeSpansResult.variables);
            addTimeSpanToCache(cache, data.createTimeSpan);
        },
    });

    const values = (() => {
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
            .map((ts) => {
                const colorKey = ts
                    .tags!.map((t) => t.key + ':' + t.value)
                    .sort((a, b) => a.localeCompare(b))
                    .join(' ');
                const color = calculateColor(colorKey, ColorMode.Bold, theme.palette.type);
                const borderColor = calculateColor(colorKey, ColorMode.None, theme.palette.type);
                return {
                    start: moment(ts.start).toDate(),
                    end: moment(ts.end || currentDate).toDate(),
                    hasEnd: !!ts.end,
                    editable: !!ts.end,
                    backgroundColor: color,
                    startEditable: true,
                    id: ts.id,
                    tags: ts.tags!.map(({value, key}) => ({key, value})),
                    title: ts.tags!.map((t) => t.key + ':' + t.value).join(' '),
                    extendedProps: {ts},
                    textColor: theme.palette.getContrastText(color),
                    borderColor,
                };
            });
    })();

    const onDrop: OptionsInput['eventDrop'] = (data) => {
        updateTimeSpanMutation({
            variables: {
                oldStart: moment(data.oldEvent.start!).format(),
                start: moment(data.event.start!).format(),
                end: moment(data.event.end!).format(),
                id: parseInt(data.event.id, 10),
                tags: stripTypename(data.event.extendedProps.ts.tags),
            },
        });
    };
    const onResize: OptionsInput['eventResize'] = (data) => {
        updateTimeSpanMutation({
            variables: {
                oldStart: moment(data.prevEvent.start!).format(),
                start: moment(data.event.start!).format(),
                end: moment(data.event.end!).format(),
                id: parseInt(data.event.id, 10),
                tags: stripTypename(data.event.extendedProps.ts.tags),
            },
        });
    };
    const onSelect: OptionsInput['select'] = (data) => {
        addTimeSpan({
            variables: {
                start: moment(data.start).format(),
                end: moment(data.end).format(),
                tags: [],
            },
        });
    };
    const onClick: OptionsInput['eventClick'] = (data) => {
        data.jsEvent.preventDefault();
        // tslint:disable-next-line:no-any
        setSelected({data: data.event.extendedProps.ts, selected: data.jsEvent.target as any});
    };

    return (
        <Paper style={{padding: 10, bottom: 10, top: 80, position: 'absolute'}} color="red">
            <FullCalendarStyling>
                <FullCalendar
                    defaultView="timeGridWeek"
                    rerenderDelay={10}
                    datesRender={(x) => {
                        const range = {start: moment(x.view.currentStart), end: moment(x.view.currentEnd)};
                        if (
                            !moment(timeSpansResult.variables.start).isSame(range.start) ||
                            !moment(timeSpansResult.variables.end).isSame(range.end)
                        ) {
                            timeSpansResult.refetch(range);
                        }
                    }}
                    editable={true}
                    events={values}
                    allDaySlot={false}
                    selectable={true}
                    selectMirror={true}
                    handleWindowResize={true}
                    height={'parent'}
                    selectMinDistance={20}
                    now={currentDate.toDate()}
                    defaultTimedEventDuration={{minute: 15}}
                    eventRender={(e) => {
                        const content = e.el.getElementsByClassName('fc-content').item(0);
                        if (content) {
                            content.innerHTML = getElementContent(e.event, () => {
                                stopTimer({
                                    variables: {id: e.event.extendedProps.ts.id, end: moment().format()},
                                });
                            });
                        }

                        e.el.setAttribute('data-has-end', '' + (!e.event.extendedProps.ts || !!e.event.extendedProps.ts.end));
                    }}
                    slotLabelInterval={{minute: 60}}
                    slotDuration={{minute: 15}}
                    scrollTime={{hour: 6, minute: 30}}
                    select={onSelect}
                    firstDay={moment.localeData().firstDayOfWeek()}
                    eventResize={onResize}
                    eventClick={onClick}
                    eventDrop={onDrop}
                    slotLabelFormat={(s) => toMoment(s.start.marker).format('LT')}
                    columnHeaderFormat={(s) => toMoment(s.start.marker).format('DD dddd')}
                    nowIndicator={true}
                    plugins={[dayGridPlugin, timeGridPlugin, interactionPlugin, momentPlugin]}
                    header={{
                        center: 'title',
                        left: 'prev,next today',
                        right: 'timeGridWeek,timeGridDay',
                    }}
                />
            </FullCalendarStyling>
            {!trackersResult.data || !(trackersResult.data.timers || []).length ? (
                <NowIndicator>
                    <Button
                        className={classes.start}
                        variant="contained"
                        color="primary"
                        onClick={() => {
                            startTimer({variables: {start: moment().format(), tags: []}}).then(() => {
                                setCurrentDate(moment());
                            });
                        }}>
                        Start
                    </Button>
                </NowIndicator>
            ) : null}
            {!!selected.selected && (
                <Popper open={true} anchorEl={selected.selected} style={{zIndex: 1200, maxWidth: 700}}>
                    <ClickAwayListener
                        onClickAway={() => {
                            if (ignore) {
                                return;
                            }
                            setSelected({selected: null, data: null});
                        }}>
                        <div>
                            <TimeSpan
                                elevation={10}
                                id={selected.data!.id}
                                key={selected.data!.id}
                                deleted={() => {
                                    removeFromTimeSpanInRangeCache(apollo.cache, selected.data!.id, timeSpansResult.variables);
                                    setSelected({selected: null, data: null});
                                }}
                                continued={() => setCurrentDate(moment())}
                                range={{
                                    from: moment(selected.data!.start),
                                    to: selected.data!.end ? moment(selected.data!.end) : undefined,
                                }}
                                initialTags={toTagSelectorEntry(
                                    tagsResult.data!.tags!,
                                    selected.data!.tags!.map((tag) => ({key: tag.key, value: tag.value}))
                                )}
                                dateSelectorOpen={setIgnore}
                                stopped={() => {
                                    setSelected({
                                        ...selected,
                                        data: {...selected.data!, end: currentDate.toDate()},
                                    });
                                }}
                            />
                        </div>
                    </ClickAwayListener>
                </Popper>
            )}
        </Paper>
    );
};

const NowIndicator: React.FC = ({children}) => {
    const [indicator, setIndicator] = React.useState<Element | null>(null);

    useInterval(
        () => {
            if (!indicator) {
                setIndicator(document.getElementsByClassName('fc-now-indicator').item(0));
            }
        },
        500,
        true
    );
    if (!indicator) {
        return <></>;
    }
    return ReactDOM.createPortal(children, indicator);
};

const getElementContent = (event: EventApi, stop: () => void): string => {
    if (!event.start || !event.end) {
        return '';
    }
    const start = moment(event.start);
    const end = moment(event.end);
    const diff = end.diff(start, 'minute');

    const lines = Math.floor(diff / 15);
    const hasEnd = !event.extendedProps.ts || event.extendedProps.ts.end;

    let stopButton = '';
    if (!hasEnd) {
        const id = event.extendedProps.ts.id;
        if (!window.__TRAGGO_CALENDAR) {
            window.__TRAGGO_CALENDAR = {};
        }
        window.__TRAGGO_CALENDAR[id] = (e: Event) => {
            e.preventDefault();
            e.stopPropagation();
            stop();
            return false;
        };
        stopButton = `<div class="stop"><a onClick="return window.__TRAGGO_CALENDAR[${id}](event)">STOP ${timeRunningCalendar(
            start,
            end
        )}</a></div>`;
    }
    const clamp = (amount: number) =>
        `<span class="ellipsis" title="${event.title}" style="-webkit-line-clamp: ${amount}">${event.title}</span>`;

    const running = hasEnd ? `<span style="float: right">${timeRunningCalendar(start, end)}</span>` : '';
    const date = `${start.format('LT')} - ${hasEnd ? end.format('LT') : 'now'} ${running}`;
    if (lines < 2) {
        return event.title ? `<span class="ellipsis-single" title="${event.title}">${event.title}</span>${stopButton}` : date;
    }
    if (lines === 2) {
        if (hasEnd) {
            return `${date}<span class="ellipsis-single" title="${event.title}">${event.title}</span>${stopButton}`;
        } else {
            return `${clamp(2)}${stopButton}`;
        }
    }

    return `${date}<br/>${clamp(lines - 1)}${stopButton}`;
};
