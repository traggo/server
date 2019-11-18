import * as React from 'react';
import {useMutation} from '@apollo/react-hooks';
import * as gqlDashboard from '../gql/dashboard';
import {Dashboards_dashboards, Dashboards_dashboards_ranges} from '../gql/__generated__/Dashboards';
import {IconButton, Paper, Typography} from '@material-ui/core';
import {RelativeDateTimeSelector} from '../common/RelativeDateTimeSelector';
import ArrowForwardIcon from '@material-ui/icons/ArrowForward';
import MoreVert from '@material-ui/icons/MoreVert';
import PlusIcon from '@material-ui/icons/Add';
import Menu from '@material-ui/core/Menu';
import MenuItem from '@material-ui/core/MenuItem';
import {RemoveDashboardRange, RemoveDashboardRangeVariables} from '../gql/__generated__/RemoveDashboardRange';
import {UpdateDashboardRange, UpdateDashboardRangeVariables} from '../gql/__generated__/UpdateDashboardRange';
import {stripTypename} from '../utils/strip';
import {Range} from '../utils/range';
import {AddDashboardRange, AddDashboardRangeVariables} from '../gql/__generated__/AddDashboardRange';
import Input from '@material-ui/core/Input';

interface Props {
    changeMode: boolean;
    dashboard: Dashboards_dashboards;
    setRanges: (cb: (ranges: Record<number, Range>) => Record<number, Range>) => void;
    ranges: Record<number, Range>;
}

export const DateRanges: React.FC<Props> = ({changeMode, dashboard, ranges, setRanges}) => {
    const saveRef = React.useRef<number | null>(null);
    const [editedRanges, setEditedRanges] = React.useState<Record<number, Range>>({});
    const [editedNames, setEditedNames] = React.useState<Record<number, string>>({});
    const [openMenu, setOpenMenu] = React.useState<null | [HTMLButtonElement, number]>(null);

    const [removeRange] = useMutation<RemoveDashboardRange, RemoveDashboardRangeVariables>(gqlDashboard.RemoveDashboardRange, {
        refetchQueries: [{query: gqlDashboard.Dashboards}],
    });
    const [updateRange] = useMutation<UpdateDashboardRange, UpdateDashboardRangeVariables>(gqlDashboard.UpdateDashboardRange, {
        refetchQueries: [{query: gqlDashboard.Dashboards}],
    });
    const [addRange] = useMutation<AddDashboardRange, AddDashboardRangeVariables>(gqlDashboard.AddDashboardRange, {
        refetchQueries: [{query: gqlDashboard.Dashboards}],
    });
    const saveRanges = (range: Dashboards_dashboards_ranges, newRange: Range, newName: string) => {
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
                                ) : (
                                    undefined
                                )}
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
                                ) : (
                                    undefined
                                )}
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
                        }}>
                        <PlusIcon />
                    </IconButton>
                </Paper>
            ) : (
                undefined
            )}
        </>
    );
};
