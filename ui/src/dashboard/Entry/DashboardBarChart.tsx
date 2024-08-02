import {Stats_stats} from '../../gql/__generated__/Stats';
import {Bar, BarChart, CartesianGrid, Legend, ResponsiveContainer, Tooltip, XAxis, YAxis} from 'recharts';
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
    type: 'stacked' | 'normal';
    total: boolean;
}

interface Indexed {
    start: string;
    end: string;
    data: Record<string, number>;
}

export const DashboardBarChart: React.FC<DashboardPieChartProps> = ({entries, interval, type, total}) => {
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
            <BarChart data={indexedEntries}>
                <CartesianGrid strokeDasharray="3 3" />
                <YAxis type="number" unit={unit.short} />
                <Tooltip content={<TagTooltip dateFormat={dateFormat} total={total} />} />
                <Legend />
                <XAxis dataKey={(entry) => dateFormat(moment(entry.start))} interval={'preserveStartEnd'} />

                {indexedEntries[0] &&
                    Object.keys(indexedEntries[0].data).map((key, index) => {
                        return (
                            <Bar
                                key={key}
                                dataKey={(entry) => unit.toUnit(entry.data[key])}
                                fill={Colors[index % Colors.length]}
                                stackId={type === 'stacked' ? 'a' : undefined}
                                name={key}
                            />
                        );
                    })}
            </BarChart>
        </ResponsiveContainer>
    );
};
