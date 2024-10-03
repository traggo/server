import * as React from 'react';
import makeStyles from '@material-ui/core/styles/makeStyles';
import {Paper} from '@material-ui/core';
import {SetSettings as SetSettingsGQL, Settings as SettingsGQL, useSettings} from '../gql/settings';
import {useMutation} from '@apollo/react-hooks';
import {SetSettings, SetSettingsVariables} from '../gql/__generated__/SetSettings';
import FormControl from '@material-ui/core/FormControl';
import InputLabel from '@material-ui/core/InputLabel';
import Select from '@material-ui/core/NativeSelect/NativeSelect';
import {DateLocale, Theme, WeekDay} from '../gql/__generated__/globalTypes';
import {useSnackbar} from 'notistack';
import {handleError} from '../utils/errors';

const useStyles = makeStyles((theme) => ({
    root: {
        ...theme.mixins.gutters(),
        paddingTop: theme.spacing(1),
        paddingBottom: theme.spacing(3),
        maxWidth: 500,
        margin: '0 auto',
    },
}));

export const SettingsPage: React.FC = () => {
    const classes = useStyles();
    const {done, ...settings} = useSettings();
    const {enqueueSnackbar} = useSnackbar();
    const [setSettings] = useMutation<SetSettings, SetSettingsVariables>(SetSettingsGQL, {
        refetchQueries: [{query: SettingsGQL}],
    });
    return (
        <Paper elevation={1} className={classes.root}>
            <FormControl margin={'normal'} fullWidth>
                <InputLabel>Date Locale</InputLabel>
                <Select
                    fullWidth
                    value={settings.dateLocale}
                    onChange={(e) => {
                        setSettings({
                            variables: {
                                settings: {
                                    ...settings,
                                    dateLocale: e.target.value as DateLocale,
                                },
                            },
                        })
                            .then(() => {
                                enqueueSnackbar('date locale changed', {
                                    variant: 'success',
                                });
                                enqueueSnackbar('a reload of the page is required for the new date locale to fully function', {
                                    variant: 'info',
                                    preventDuplicate: true,
                                    persist: true,
                                });
                            })
                            .catch(handleError('set date locale', enqueueSnackbar));
                    }}>
                    {Object.values(DateLocale).map((type) => (
                        <option key={type} value={type}>
                            {type}
                        </option>
                    ))}
                </Select>
            </FormControl>
            <FormControl margin={'normal'} fullWidth>
                <InputLabel>Theme</InputLabel>
                <Select
                    fullWidth
                    value={settings.theme}
                    onChange={(e) => {
                        setSettings({variables: {settings: {...settings, theme: e.target.value as Theme}}})
                            .then(() =>
                                enqueueSnackbar('theme changed', {
                                    variant: 'success',
                                })
                            )
                            .catch(handleError('set theme', enqueueSnackbar));
                    }}>
                    {Object.values(Theme).map((type) => (
                        <option key={type} value={type}>
                            {type}
                        </option>
                    ))}
                </Select>
            </FormControl>
            <FormControl margin={'normal'} fullWidth>
                <InputLabel>First day of the week</InputLabel>
                <Select
                    fullWidth
                    value={settings.firstDayOfTheWeek}
                    onChange={(e) => {
                        setSettings({
                            variables: {
                                settings: {
                                    ...settings,
                                    firstDayOfTheWeek: e.target.value as WeekDay,
                                },
                            },
                        })
                            .then(() =>
                                enqueueSnackbar('first day of the week changed', {
                                    variant: 'success',
                                })
                            )
                            .catch(handleError('set first day of the week', enqueueSnackbar));
                    }}>
                    {[
                        WeekDay.Sunday,
                        WeekDay.Monday,
                        WeekDay.Tuesday,
                        WeekDay.Wednesday,
                        WeekDay.Thursday,
                        WeekDay.Friday,
                        WeekDay.Saturday,
                    ].map((type) => (
                        <option key={type} value={type}>
                            {type}
                        </option>
                    ))}
                </Select>
            </FormControl>
        </Paper>
    );
};
