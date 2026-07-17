import * as React from 'react';
import {DesktopDateTimePicker} from '@mui/x-date-pickers/DesktopDateTimePicker';
import moment from 'moment';
import {uglyConvertToLocalTime} from '../timespan/timeutils';
import {useSettings} from '../gql/settings';
import {DateTimeInputStyle} from '../gql/__generated__';

interface DateTimeSelectorProps {
    selectedDate: moment.Moment;
    onSelectDate: (date: moment.Moment) => void;
    showDate: boolean;
    label: string;
    popoverOpen?: (open: boolean) => void;
}

export const DateTimeSelector: React.FC<DateTimeSelectorProps> = React.memo(
    ({selectedDate, onSelectDate, showDate, label, popoverOpen = () => {}}) => {
        const {done, dateTimeInputStyle} = useSettings();
        const [open, setOpen] = React.useState(false);

        if (!done) {
            return <span>...</span>;
        }

        if (dateTimeInputStyle === DateTimeInputStyle.Native) {
            return (
                <input
                    type="datetime-local"
                    value={selectedDate.format(selectedDate.format('YYYY-MM-DDTHH:mm'))}
                    onChange={(e) => {
                        onSelectDate(moment(e.target.value));
                    }}
                />
            );
        }
        const localeData = moment.localeData();
        const time = localeData.longDateFormat('LT').replace('A', 'a');
        const ampm = time.indexOf('a') !== -1;
        const format = showDate ? localeData.longDateFormat('L') + ' ' + time : time;

        return (
            <DesktopDateTimePicker
                className="time-picker"
                sx={{width: (showDate ? 185 : 105) + (ampm ? 20 : 0)}}
                onOpen={() => {
                    popoverOpen(true);
                    setOpen(true);
                }}
                onClose={() => {
                    popoverOpen(false);
                    setOpen(false);
                }}
                slotProps={{
                    textField: {
                        title: selectedDate.format(),
                        variant: 'standard',
                        margin: 'none',
                        InputProps: {disableUnderline: true},
                    },
                }}
                value={uglyConvertToLocalTime(selectedDate)}
                onChange={(date: moment.Moment | null) => {
                    if (!date || !date.isValid()) {
                        return;
                    }

                    if (!showDate && !open) {
                        date = date.set({
                            date: selectedDate.date(),
                            month: selectedDate.month(),
                            year: selectedDate.year(),
                        });
                    }
                    if (uglyConvertToLocalTime(selectedDate).isSame(date)) {
                        return;
                    }

                    onSelectDate(date);
                }}
                ampm={ampm}
                format={format}
                label={label}
                openTo={showDate ? 'day' : 'hours'}
            />
        );
    }
);
