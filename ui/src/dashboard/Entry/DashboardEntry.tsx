import {DashboardItem} from '../../gql/types';
import * as React from 'react';
import {DashboardPieChart} from './DashboardPieChart';
import {StatEntry} from '../../gql/types';
import {useQuery} from '@apollo/client';
import * as gqlStats from '../../gql/statistics';
import {Paper} from '@material-ui/core';
import Typography from '@material-ui/core/Typography';
import moment from 'moment-timezone';
import {EntryType} from '../../gql/__generated__';
import {expectNever} from '../../utils/never';
import {Stats2Query, Stats2QueryVariables} from '../../gql/__generated__';
import {DashboardBarChart} from './DashboardBarChart';
import {DashboardLineChart} from './DashboardLineChart';
import {CenteredSpinner} from '../../common/CenteredSpinner';
import {Center} from '../../common/Center';
import {findRange, Range} from '../../utils/range';
import {DashboardTable} from './DashboardTable';

interface DashboardEntryProps {
    entry: DashboardItem;
    ranges: Record<number, Range>;
}

export const DashboardEntry = React.forwardRef<HTMLDivElement, DashboardEntryProps>(({entry, ranges}, ref) => {
    const range = findRange(
        {...entry.statsSelection, range: entry.statsSelection.range ?? null, rangeId: entry.statsSelection.rangeId ?? null},
        ranges
    );
    return (
        <Paper style={{width: '100%', height: '100%'}} ref={ref}>
            <Typography style={{lineHeight: '20px', paddingTop: 15}} component="h4" align="center" variant="h6">
                {entry.title}
            </Typography>
            <Typography style={{position: 'absolute', right: 10, top: 35, fontSize: 10, color: 'gray'}}>
                {range.from} to {range.to}
            </Typography>
            <div style={{height: 'calc(100% - 35px)'}}>
                <SpecificDashboardEntry range={range} entry={entry} />
            </div>
        </Paper>
    );
});

// tslint:disable-next-line:cyclomatic-complexity mccabe-complexity
const SpecificDashboardEntry: React.FC<{entry: DashboardItem; range: Range}> = ({entry, range}) => {
    const interval = entry.statsSelection.interval;
    const stats = useQuery<Stats2Query, Stats2QueryVariables>(gqlStats.Stats2, {
        variables: {
            now: moment().startOf('hour').format(),
            stats: {
                range,
                interval,
                tags: entry.statsSelection.tags,
                excludeTags: entry.statsSelection.excludeTags,
                includeTags: entry.statsSelection.includeTags,
            },
        },
    });

    if (stats.loading) {
        return <CenteredSpinner />;
    }

    if (stats.error) {
        return (
            <Center>
                <Typography>error: {stats.error.message}</Typography>
            </Center>
        );
    }
    const firstEntries: StatEntry[] =
        (stats.data && stats.data.stats && stats.data.stats[0] && stats.data.stats[0].entries) || [];
    if (firstEntries.length === 0) {
        return (
            <Center>
                <Typography>no data</Typography>
            </Center>
        );
    }

    const entries = (stats.data && stats.data.stats) || [];
    switch (entry.entryType) {
        case EntryType.PieChart:
            return <DashboardPieChart entries={firstEntries} />;
        case EntryType.BarChart:
            return <DashboardBarChart entries={entries} interval={interval} type="normal" total={entry.total} />;
        case EntryType.StackedBarChart:
            return <DashboardBarChart entries={entries} interval={interval} type="stacked" total={entry.total} />;
        case EntryType.LineChart:
            return <DashboardLineChart entries={entries} interval={interval} total={entry.total} />;
        case EntryType.VerticalTable:
            return <DashboardTable mode="vertical" entries={entries} interval={interval} total={entry.total} />;
        case EntryType.HorizontalTable:
            return <DashboardTable mode="horizontal" entries={entries} interval={interval} total={entry.total} />;
        default:
            return expectNever(entry.entryType);
    }
};
