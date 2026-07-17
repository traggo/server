import * as React from 'react';
import Popper from '@mui/material/Popper';
import ClickAwayListener from '@mui/material/ClickAwayListener';
import {Paper} from '@mui/material';
import Typography from '@mui/material/Typography';
import Button from '@mui/material/Button';
import {DashboardItem} from '../../gql/types';
import {useMutation} from '@apollo/client';
import * as gqlDashboard from '../../gql/dashboard';
import {UpdateDashboardEntryMutation, UpdateDashboardEntryMutationVariables} from '../../gql/__generated__';
import {Fade} from '../../common/Fade';
import {DashboardEntryForm, isValidDashboardEntry} from './DashboardEntryForm';
import {handleError} from '../../utils/errors';
import {useSnackbar} from 'notistack';

interface EditPopupProps {
    entry: DashboardItem;
    anchorEl: HTMLElement;
    onChange: (entry: DashboardItem | null) => void;
    doPreview: (preview: boolean) => void;
    preview: boolean;
    ranges: Record<number, string>;
}
export const EditPopup: React.FC<EditPopupProps> = ({entry, anchorEl, onChange: setEdit, doPreview, preview, ranges}) => {
    const [updateEntry] = useMutation<UpdateDashboardEntryMutation, UpdateDashboardEntryMutationVariables>(
        gqlDashboard.UpdateDashboardEntry,
        {
            refetchQueries: [{query: gqlDashboard.Dashboards}],
        }
    );
    const valid = isValidDashboardEntry(entry);

    const {enqueueSnackbar} = useSnackbar();

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
                                                excludeTags: entry.statsSelection.excludeTags,
                                                includeTags: entry.statsSelection.includeTags,
                                            },
                                        },
                                    })
                                        .then(() => setEdit(null))
                                        .catch(handleError('Edit Dashboard Entry', enqueueSnackbar));
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
