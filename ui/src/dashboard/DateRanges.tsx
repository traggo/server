import * as React from 'react';
import {useMutation} from '@apollo/client';
import * as gqlDashboard from '../gql/dashboard';
import {Dashboard, DashboardRange} from '../gql/types';
import {IconButton, Paper, Typography} from '@mui/material';
import {RelativeDateTimeSelector} from '../common/RelativeDateTimeSelector';
import ArrowForwardIcon from '@mui/icons-material/ArrowForward';
import MoreVert from '@mui/icons-material/MoreVert';
import PlusIcon from '@mui/icons-material/Add';
import Menu from '@mui/material/Menu';
import MenuItem from '@mui/material/MenuItem';
import {
    RemoveDashboardRangeMutation,
    RemoveDashboardRangeMutationVariables,
    UpdateDashboardRangeMutation,
    UpdateDashboardRangeMutationVariables,
    AddDashboardRangeMutation,
    AddDashboardRangeMutationVariables,
} from '../gql/__generated__';
import {stripTypename} from '../utils/strip';
import {Range} from '../utils/range';
import Input from '@mui/material/Input';

interface Props {
    changeMode: boolean;
    dashboard: Dashboard;
    setRanges: (cb: (ranges: Record<number, Range>) => Record<number, Range>) => void;
    ranges: Record<number, Range>;
}

export const DateRanges: React.FC<Props> = ({changeMode, dashboard, ranges, setRanges}) => {
    const saveRef = React.useRef<number | null>(null);
    const [editedRanges, setEditedRanges] = React.useState<Record<number, Range>>({});
    const [editedNames, setEditedNames] = React.useState<Record<number, string>>({});
    const [openMenu, setOpenMenu] = React.useState<null | [HTMLButtonElement, number]>(null);

    const [removeRange] = useMutation<RemoveDashboardRangeMutation, RemoveDashboardRangeMutationVariables>(
        gqlDashboard.RemoveDashboardRange,
        {
            refetchQueries: [{query: gqlDashboard.Dashboards}],
        }
    );
    const [updateRange] = useMutation<UpdateDashboardRangeMutation, UpdateDashboardRangeMutationVariables>(
        gqlDashboard.UpdateDashboardRange,
        {
            refetchQueries: [{query: gqlDashboard.Dashboards}],
        }
    );
    const [addRange] = useMutation<AddDashboardRangeMutation, AddDashboardRangeMutationVariables>(
        gqlDashboard.AddDashboardRange,
        {
            refetchQueries: [{query: gqlDashboard.Dashboards}],
        }
    );
    const saveRanges = (range: DashboardRange, newRange: Range, newName: string) => {
        if (saveRef.current) {
            clearTimeout(saveRef.current);
        }
        saveRef.current = window.setTimeout(() => {
            updateRange({
                variables: {
                    rangeId: range.id,
                    range: {
                        editable: !range.editable,
                        range: newRange,
                        name: newName,
                    },
                },
            });
        }, 500);
    };
    React.useEffect(() => {
        if (changeMode) {
            setEditedRanges({});
            setEditedNames({});
        }
    }, [changeMode]);

    return (
        <>
            {dashboard.ranges
                .filter((range) => range.editable || changeMode)
                .map((range) => {
                    const name = editedNames[range.id] !== undefined ? editedNames[range.id] : range.name;
                    const dateRange = (changeMode ? editedRanges[range.id] : ranges[range.id]) || stripTypename(range.range);
                    return (
                        <React.Fragment key={range.id}>
                            <Paper
                                style={{display: 'inline-block', padding: '3px 10px', marginLeft: 10, marginBottom: 10}}
                                elevation={1}>
                                <div style={{display: 'inline-block'}}>
                                    <Typography variant={'subtitle2'} component="div">
                                        {changeMode ? (
                                            <Input
                                                disableUnderline={true}
                                                fullWidth
                                                onChange={(e) => {
                                                    const value = e.target.value;
                                                    setEditedNames((old) => ({
                                                        ...old,
                                                        [range.id]: value,
                                                    }));
                                                    saveRanges(range, dateRange, value);
                                                }}
                                                value={name}
                                            />
                                        ) : (
                                            name
                                        )}
                                    </Typography>
                                    <RelativeDateTimeSelector
                                        small={!changeMode}
                                        disableUnderline={!changeMode}
                                        style={{width: changeMode ? 170 : 120}}
                                        value={dateRange.from}
                                        onChange={(value, valid) => {
                                            (changeMode ? setEditedRanges : setRanges)((old) => ({
                                                ...old,
                                                [range.id]: {...dateRange, from: value},
                                            }));
                                            if (valid && changeMode) {
                                                saveRanges(range, {...dateRange, from: value}, name);
                                            }
                                        }}
                                        type="startOf"
                                    />
                                    <ArrowForwardIcon style={{margin: '0 10px'}} />
                                    <RelativeDateTimeSelector
                                        small={!changeMode}
                                        disableUnderline={!changeMode}
                                        style={{width: changeMode ? 170 : 120}}
                                        value={dateRange.to}
                                        onChange={(value, valid) => {
                                            (changeMode ? setEditedRanges : setRanges)((old) => ({
                                                ...old,
                                                [range.id]: {...dateRange, to: value},
                                            }));
                                            if (valid && changeMode) {
                                                saveRanges(
                                                    range,
                                                    {
                                                        ...dateRange,
                                                        to: value,
                                                    },
                                                    name
                                                );
                                            }
                                        }}
                                        type="endOf"
                                    />
                                </div>
                                {changeMode ? (
                                    <>
                                        <Typography style={{height: 50, display: 'inline-block'}} component="div">
                                            <IconButton size="medium" onClick={(e) => setOpenMenu([e.currentTarget, range.id])}>
                                                <MoreVert />
                                            </IconButton>
                                        </Typography>
                                    </>
                                ) : undefined}
                                {openMenu && openMenu[1] === range.id ? (
                                    <Menu
                                        key="uff"
                                        aria-haspopup="true"
                                        anchorEl={openMenu[0]}
                                        onClose={() => setOpenMenu(null)}
                                        open={openMenu !== null}>
                                        <MenuItem
                                            onClick={() => {
                                                const used = dashboard.items.filter(
                                                    (entry) => entry.statsSelection.rangeId === range.id
                                                );
                                                if (used.length === 0) {
                                                    removeRange({variables: {rangeId: range.id}});
                                                } else {
                                                    alert('Range is used by ' + used.map((entry) => entry.title).join(', '));
                                                }
                                                setOpenMenu(null);
                                            }}>
                                            Delete
                                        </MenuItem>
                                        <MenuItem
                                            onClick={() => {
                                                updateRange({
                                                    variables: {
                                                        rangeId: range.id,
                                                        range: {
                                                            editable: !range.editable,
                                                            range: stripTypename(range.range),
                                                            name: range.name,
                                                        },
                                                    },
                                                });
                                                setOpenMenu(null);
                                            }}>
                                            {range.editable ? 'make static' : 'make editable'}
                                        </MenuItem>
                                    </Menu>
                                ) : undefined}
                            </Paper>
                        </React.Fragment>
                    );
                })}
            {changeMode ? (
                <Paper style={{display: 'inline-block', padding: '3px 10px', marginLeft: 10, marginBottom: 10}} elevation={1}>
                    <IconButton
                        onClick={() => {
                            addRange({
                                variables: {
                                    dashboardId: dashboard.id,
                                    range: {name: 'new range', editable: true, range: {from: 'now/w', to: 'now/w'}},
                                },
                            });
                        }}
                        size="large">
                        <PlusIcon />
                    </IconButton>
                </Paper>
            ) : undefined}
        </>
    );
};
