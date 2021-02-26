import {Stats_stats} from '../../gql/__generated__/Stats';
import * as React from 'react';
import {StatsInterval} from '../../gql/__generated__/globalTypes';
import moment from 'moment';
import {Table, TableBody, TableCell, TableRow} from '@material-ui/core';
import TableHead from '@material-ui/core/TableHead';
import {ofInterval} from './dateformat';
import prettyMs from 'pretty-ms';

interface DashboardTableProps {
    entries: Stats_stats[];
    interval: StatsInterval;
    mode: 'vertical' | 'horizontal';
    total: boolean;
}

interface Indexed {
    start: string;
    end: string;
    data: Record<string, number>;
}

export const DashboardTable: React.FC<DashboardTableProps> = ({entries, interval, mode, total}) => {
    const indexedEntries: Indexed[] = entries
        .map((entry) => {
            const result = {
                start: entry.start,
                end: entry.end,
                data: entry.entries!.reduce((all: Record<string, number>, current) => {
                    return {...all, [current.key + ':' + current.value]: current.timeSpendInSeconds};
                }, {}),
            };
            if (total) {
                let sum = 0;
                Object.keys(result.data).forEach((item) => {
                    sum += result.data[item];
                });
                result.data['Total'] = sum;
            }
            return result;
        })
        .sort((left, right) => moment(left.start).diff(right.start));
    const dateFormat = ofInterval(interval);
    const keys = Object.keys((indexedEntries[0] && indexedEntries[0].data) || {});
    return (
        <div style={{overflow: 'auto', height: '100%', margin: 2}}>
            <Table size="small">
                {mode === 'vertical' ? (
                    <>
                        <TableHead>
                            <TableRow>
                                <TableCell>Date</TableCell>
                                {keys.map((key) => {
                                    return <TableCell key={key}>{key}</TableCell>;
                                })}
                            </TableRow>
                        </TableHead>
                        <TableBody>
                            {indexedEntries.map((entry) => {
                                return (
                                    <TableRow key={entry.start}>
                                        <TableCell>{dateFormat(moment(entry.start))}</TableCell>
                                        {keys.map((key) => (
                                            <TableCell key={key + entry.start}>{prettyMs(entry.data[key] * 1000)}</TableCell>
                                        ))}
                                    </TableRow>
                                );
                            })}
                        </TableBody>
                    </>
                ) : (
                    <>
                        <TableHead>
                            <TableRow>
                                <TableCell>Date</TableCell>
                                {indexedEntries.map((entry) => {
                                    return <TableCell key={entry.start}>{dateFormat(moment(entry.start))}</TableCell>;
                                })}
                            </TableRow>
                        </TableHead>
                        <TableBody>
                            {keys.map((key) => {
                                return (
                                    <TableRow key={key}>
                                        <TableCell>{key}</TableCell>
                                        {indexedEntries.map((entry) => (
                                            <TableCell key={key + entry.start}>{prettyMs(entry.data[key] * 1000)}</TableCell>
                                        ))}
                                    </TableRow>
                                );
                            })}
                        </TableBody>
                    </>
                )}
            </Table>
        </div>
    );
};
