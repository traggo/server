import {Stats_stats_entries} from '../../gql/__generated__/Stats';
import {Cell, Legend, Pie, PieChart, ResponsiveContainer, Tooltip, TooltipProps} from 'recharts';
import * as React from 'react';
import {Colors} from './colors';
import {Typography} from '@material-ui/core';
import prettyMs from 'pretty-ms';
import Paper from '@material-ui/core/Paper';

interface DashboardPieChartProps {
    entries: Stats_stats_entries[];
}

export const DashboardPieChart: React.FC<DashboardPieChartProps> = ({entries}) => {
    const total = entries.map((e) => e.timeSpendInSeconds).reduce((x, y) => x + y, 0);
    return (
        <ResponsiveContainer>
            <PieChart>
                <Pie
                    isAnimationActive={false}
                    dataKey="timeSpendInSeconds"
                    nameKey={(entry) => {
                        // tslint:disable-next-line:no-any
                        return (entry.key + ':' + entry.value) as any;
                    }}
                    data={entries}
                    labelLine={false}
                    fill="#8884d8"
                    legendType={'square'}>
                    {entries.map((_, index) => (
                        <Cell key={index} fill={Colors[index % Colors.length]} />
                    ))}
                </Pie>
                <Tooltip content={<CustomTooltip total={total} />} />
                <Legend />
            </PieChart>
        </ResponsiveContainer>
    );
};

interface CustomTooltipProps extends TooltipProps {
    total: number;
}

const CustomTooltip = ({active, payload, total}: CustomTooltipProps) => {
    if (active && payload) {
        const first = payload[0];
        return (
            <Paper style={{padding: 10}} elevation={4}>
                <Typography>
                    {first.payload.key}:{first.payload.value}: {prettyMs(first.payload.timeSpendInSeconds * 1000)} (
                    {((first.payload.timeSpendInSeconds / total) * 100).toFixed(2)}%)
                </Typography>
            </Paper>
        );
    }

    return null;
};
