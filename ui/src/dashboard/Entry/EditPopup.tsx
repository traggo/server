import * as React from 'react';
import Popper from '@material-ui/core/Popper';
import ClickAwayListener from '@material-ui/core/ClickAwayListener';
import {Paper} from '@material-ui/core';
import Typography from '@material-ui/core/Typography';
import Button from '@material-ui/core/Button';
import {Dashboards_dashboards_items} from '../../gql/__generated__/Dashboards';
import {useMutation} from '@apollo/react-hooks';
import * as gqlDashboard from '../../gql/dashboard';
import {UpdateDashboardEntry, UpdateDashboardEntryVariables} from '../../gql/__generated__/UpdateDashboardEntry';
import {Fade} from '../../common/Fade';
import {DashboardEntryForm, isValidDashboardEntry} from './DashboardEntryForm';

interface EditPopupProps {
    entry: Dashboards_dashboards_items;
    anchorEl: HTMLElement;
    onChange: (entry: Dashboards_dashboards_items | null) => void;
    doPreview: (preview: boolean) => void;
    preview: boolean;
    ranges: Record<number, string>;
}
export const EditPopup: React.FC<EditPopupProps> = ({entry, anchorEl, onChange: setEdit, doPreview, preview, ranges}) => {
    const [updateEntry] = useMutation<UpdateDashboardEntry, UpdateDashboardEntryVariables>(gqlDashboard.UpdateDashboardEntry, {
        refetchQueries: [{query: gqlDashboard.Dashboards}],
    });
    const valid = isValidDashboardEntry(entry);
    return (
        <Popper
            key="popup"
            open={true}
            anchorEl={anchorEl}
            placement={'right-start'}
            disablePortal={false}
            style={{zIndex: 99999}}
            keepMounted={true}>
            <Fade fullyVisible={!preview} opacity={0.2}>
                <ClickAwayListener onClickAway={() => setEdit(null)}>
                    <Paper style={{padding: 10, maxWidth: 500}}>
                        <Typography variant="h5">Edit</Typography>
                        <DashboardEntryForm ranges={ranges} entry={entry} onChange={setEdit} disabled={preview} />
                        <div style={{textAlign: 'right', display: 'flex', justifyContent: 'flex-end', paddingTop: 10}}>
                            <Button
                                color={'primary'}
                                variant={'outlined'}
                                style={{marginRight: 10, flex: 1}}
                                disabled={!valid}
                                onClick={() => doPreview(!preview)}>
                                {preview ? 'Exit Preview' : 'Preview'}
                            </Button>
                            <Button
                                color={'secondary'}
                                variant={'outlined'}
                                style={{marginRight: 10}}
                                onClick={() => setEdit(null)}>
                                Cancel
                            </Button>
                            <Button
                                color={'primary'}
                                variant={'contained'}
                                disabled={!valid}
                                onClick={() => {
                                    updateEntry({
                                        variables: {
                                            entryId: entry.id,
                                            entryType: entry.entryType,
                                            title: entry.title,
                                            total: entry.total,
                                            stats: {
                                                tags: entry.statsSelection.tags,
                                                interval: entry.statsSelection.interval,
                                                range: entry.statsSelection.range
                                                    ? {
                                                          from: entry.statsSelection.range.from,
                                                          to: entry.statsSelection.range.to,
                                                      }
                                                    : null,
                                                rangeId: entry.statsSelection.rangeId,
                                            },
                                        },
                                    }).then(() => setEdit(null));
                                }}>
                                Save
                            </Button>
                        </div>
                    </Paper>
                </ClickAwayListener>
            </Fade>
        </Popper>
    );
};
