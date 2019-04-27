import * as React from 'react';
import {TagSelectorEntry, toInputTags} from '../tag/tagSelectorEntry';
import {TagSelector} from '../tag/TagSelector';
import moment from 'moment';
import Paper from '@material-ui/core/Paper';
import {DateTimeSelector} from '../common/DateTimeSelector';
import {Button} from '@material-ui/core';
import {timeRunning} from './timeutils';
import {useMutation} from 'react-apollo-hooks';
import {StopTimer, StopTimerVariables} from '../gql/__generated__/StopTimer';
import * as gqlTimeSpan from '../gql/timeSpan';
import {UpdateTimeSpan, UpdateTimeSpanVariables} from '../gql/__generated__/UpdateTimeSpan';

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
}

export const TimeSpan: React.FC<TimeSpanProps> = ({range: {from, to}, id, initialTags, now}) => {
    if (!to && !now) {
        throw new Error('now must be set when to is not set');
    }

    const [selectedEntries, setSelectedEntries] = React.useState<TagSelectorEntry[]>(initialTags);
    const stopTimer = useMutation<StopTimer, StopTimerVariables>(gqlTimeSpan.StopTimer, {
        refetchQueries: [{query: gqlTimeSpan.Trackers}, {query: gqlTimeSpan.TimeSpans}],
    });
    const updateTimeSpan = useMutation<UpdateTimeSpan, UpdateTimeSpanVariables>(gqlTimeSpan.UpdateTimeSpan, {
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
                                start: from,
                                end: to,
                                tags: toInputTags(entries),
                            },
                        });
                    }}
                />
            </div>
            <DateTimeSelector
                selectedDate={from}
                onSelectDate={(newFrom) => {
                    if (to && moment(newFrom).isAfter(to)) {
                        const newTo = moment(newFrom).add(15, 'minute');
                        updateTimeSpan({
                            variables: {
                                id,
                                start: newFrom.format(),
                                end: newTo.format(),
                                tags: toInputTags(selectedEntries),
                            },
                        });
                    } else {
                        updateTimeSpan({
                            variables: {
                                id,
                                start: newFrom.format(),
                                end: to && to.format(),
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
                    selectedDate={to}
                    onSelectDate={(newTo) => {
                        if (moment(newTo).isBefore(from)) {
                            const newFrom = moment(newTo).subtract(15, 'minute');
                            updateTimeSpan({
                                variables: {
                                    id,
                                    start: newFrom.format(),
                                    end: newTo.format(),
                                    tags: toInputTags(selectedEntries),
                                },
                            });
                        } else {
                            updateTimeSpan({
                                variables: {
                                    id,
                                    start: from.format(),
                                    end: newTo.format(),
                                    tags: toInputTags(selectedEntries),
                                },
                            });
                        }
                    }}
                    showDate={showDate}
                    label="end"
                />
            ) : (
                <>
                    <Button
                        style={{minWidth: 120}}
                        onClick={() => {
                            stopTimer({variables: {id, end: now}});
                        }}>
                        Stop {timeRunning(from, require(now))}
                    </Button>
                </>
            )}
        </Paper>
    );
};

const require: <T>(e?: T) => T = (e) => {
    if (!e) {
        throw new Error('unset');
    }
    return e;
};
