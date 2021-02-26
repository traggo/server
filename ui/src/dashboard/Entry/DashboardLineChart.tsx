import {Stats_stats} from '../../gql/__generated__/Stats';
import {CartesianGrid, Legend, Line, LineChart, ResponsiveContainer, Tooltip, XAxis, YAxis} from 'recharts';
import * as React from 'react';
import {Colors} from './colors';
import {ofSeconds} from './unit';
import {ofInterval} from './dateformat';
import {StatsInterval} from '../../gql/__generated__/globalTypes';
import moment from 'moment';
import {TagTooltip} from './TagTooltip';

interface DashboardPieChartProps {
    entries: Stats_stats[];
    interval: StatsInterval;
    total: boolean;
}

interface Indexed {
    start: string;
    end: string;
    data: Record<string, number>;
}

export const DashboardLineChart: React.FC<DashboardPieChartProps> = ({entries, interval, total}) => {
    const indexedEntries: Indexed[] = entries
        .map((entry) => {
            return {
                start: entry.start,
                end: entry.end,
                data: entry.entries!.reduce((all: Record<string, number>, current) => {
                    return {...all, [current.key + ':' + current.value]: current.timeSpendInSeconds};
                }, {}),
            };
        })
        .sort((left, right) => moment(left.start).diff(right.start));
    const dataMax = indexedEntries.reduce((max, entry) => {
        return Math.max(
            max,
            Object.values(entry.data).reduce((a: number, b: number) => Math.max(a, b), 0),
            0
        );
    }, 0);
    const unit = ofSeconds(dataMax);
    const dateFormat = ofInterval(interval);
    return (
        <ResponsiveContainer>
            <LineChart data={indexedEntries}>
                <CartesianGrid strokeDasharray="3 3" />
                <YAxis type="number" unit={unit.short} />
                <Tooltip content={<TagTooltip dateFormat={dateFormat} total={total} />} />
                <Legend />
                <XAxis dataKey={(entry) => dateFormat(moment(entry.start))} interval={'preserveStartEnd'} />

                {indexedEntries[0] &&
                    Object.keys(indexedEntries[0].data).map((key, index) => {
                        return (
                            <Line
                                key={key}
                                dataKey={(entry) => unit.toUnit(entry.data[key])}
                                strokeWidth={3}
                                stroke={Colors[index % Colors.length]}
                                name={key}
                            />
                        );
                    })}
            </LineChart>
        </ResponsiveContainer>
    );
};
