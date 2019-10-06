import * as React from 'react';
import TextField from '@material-ui/core/TextField';
import FormControl from '@material-ui/core/FormControl';
import InputLabel from '@material-ui/core/InputLabel';
import Select from '@material-ui/core/NativeSelect/NativeSelect';
import {EntryType, StatsInterval} from '../../gql/__generated__/globalTypes';
import {TagKeySelector} from '../../tag/TagKeySelector';
import {Dashboards_dashboards_items} from '../../gql/__generated__/Dashboards';
import {RelativeDateTimeSelector} from '../../common/RelativeDateTimeSelector';
import {parseRelativeTime} from '../../utils/time';
import {Grid, Typography, Switch} from '@material-ui/core';

interface EditPopupProps {
    entry: Dashboards_dashboards_items;
    onChange: (entry: Dashboards_dashboards_items | null) => void;
    disabled?: boolean;
    ranges: Record<string, string>;
}

export const isValidDashboardEntry = (item: Dashboards_dashboards_items): boolean => {
    return (
        item.statsSelection.tags !== null &&
        item.statsSelection.tags.length > 0 &&
        ((item.statsSelection.range &&
            parseRelativeTime(item.statsSelection.range.from, 'startOf').success &&
            parseRelativeTime(item.statsSelection.range.to, 'startOf').success) ||
            !!item.statsSelection.rangeId) &&
        item.title.length > 0
    );
};
export const DashboardEntryForm: React.FC<EditPopupProps> = ({entry, onChange: setEntry, disabled = false, ranges}) => {
    const [staticRange, setStaticRange] = React.useState(!entry.statsSelection.rangeId);
    return (
        <>
            <TextField
                label={'Title'}
                value={entry.title}
                disabled={disabled}
                onChange={(e) => {
                    entry.title = e.target.value;
                    setEntry(entry);
                }}
                fullWidth
            />
            <FormControl margin={'normal'} fullWidth>
                <InputLabel>Type</InputLabel>
                <Select
                    fullWidth
                    value={entry.entryType}
                    disabled={disabled}
                    onChange={(e) => {
                        entry.entryType = e.target.value as EntryType;
                        setEntry(entry);
                    }}>
                    {Object.values(EntryType).map((type) => (
                        <option key={type} value={type}>
                            {type}
                        </option>
                    ))}
                </Select>
            </FormControl>
            <FormControl margin={'normal'} fullWidth>
                <InputLabel>Interval</InputLabel>
                <Select
                    fullWidth
                    disabled={disabled}
                    value={entry.statsSelection.interval}
                    onChange={(e) => {
                        entry.statsSelection.interval = e.target.value as StatsInterval;
                        setEntry(entry);
                    }}>
                    {Object.values(StatsInterval).map((type) => (
                        <option key={type} value={type}>
                            {type}
                        </option>
                    ))}
                </Select>
            </FormControl>
            <Typography component="div">
                <Grid component="label" container alignItems="center" spacing={1}>
                    <Grid item>Range: </Grid>
                    <Grid item>Global</Grid>
                    <Grid item>
                        <Switch
                            checked={staticRange}
                            onChange={(e) => {
                                if (e.target.checked) {
                                    setStaticRange(true);
                                    entry.statsSelection.rangeId = null;
                                    setEntry(entry);
                                } else {
                                    setStaticRange(false);

                                    if (entry.statsSelection.range === null) {
                                        entry.statsSelection.range = {
                                            __typename: 'RelativeOrStaticRange',
                                            from: 'now/w',
                                            to: 'now/w',
                                        };
                                        setEntry(entry);
                                    }
                                }
                            }}
                        />
                    </Grid>
                    <Grid item>Static</Grid>
                </Grid>
            </Typography>
            {staticRange ? (
                <>
                    <RelativeDateTimeSelector
                        label={'From'}
                        disabled={disabled}
                        value={entry.statsSelection.range!.from}
                        onChange={(startDate) => {
                            entry.statsSelection.range!.from = startDate;
                            setEntry(entry);
                        }}
                        type="startOf"
                    />
                    <RelativeDateTimeSelector
                        label={'To'}
                        disabled={disabled}
                        value={entry.statsSelection.range!.to}
                        onChange={(startDate) => {
                            entry.statsSelection.rangeId = null;
                            entry.statsSelection.range!.to = startDate;
                            setEntry(entry);
                        }}
                        type="endOf"
                    />
                </>
            ) : (
                <Select
                    fullWidth
                    disabled={disabled}
                    value={entry.statsSelection.rangeId || ''}
                    onChange={(e) => {
                        if (e.target.value === '') {
                            return;
                        }
                        entry.statsSelection.rangeId = parseInt(e.target.value, 10);
                        entry.statsSelection.range = null;
                        setEntry(entry);
                    }}>
                    {entry.statsSelection.rangeId ? (
                        undefined
                    ) : (
                        <option key={''} value={''}>
                            Select a date range
                        </option>
                    )}
                    {Object.keys(ranges).map((key) => (
                        <option key={key} value={key}>
                            {ranges[key]}
                        </option>
                    ))}
                </Select>
            )}
            <TagKeySelector
                value={entry.statsSelection.tags || []}
                disabled={disabled}
                onChange={(tags) => {
                    entry.statsSelection.tags = tags;
                    setEntry(entry);
                }}
            />
        </>
    );
};
