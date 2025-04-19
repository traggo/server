import * as React from 'react';
import {KeyboardDateTimePicker} from '@material-ui/pickers';
import * as moment from 'moment';
import {uglyConvertToLocalTime} from '../timespan/timeutils';
import {useSettings} from '../gql/settings';
import {DateTimeInputStyle} from '../gql/__generated__/globalTypes';

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

        if (!done) {
            return <span>...</span>;
        }

        if (dateTimeInputStyle === DateTimeInputStyle.Native) {
            return (
                <input
                    type="datetime-local"
                    value={selectedDate.format(selectedDate.format('YYYY-MM-DDTHH:mm'))}
                    onChange={(e) => {
                        onSelectDate(moment.default(e.target.value));
                    }}
                />
            );
        }

        const [open, setOpen] = React.useState(false);
        const localeData = moment.localeData();
        const time = localeData.longDateFormat('LT').replace('A', 'a');
        const ampm = time.indexOf('a') !== -1;
        const format = showDate ? localeData.longDateFormat('L') + ' ' + time : time;

        return (
            <KeyboardDateTimePicker
                className="time-picker"
                variant="inline"
                InputProps={{disableUnderline: true}}
                title={selectedDate.format()}
                style={{width: (showDate ? 185 : 105) + (ampm ? 20 : 0)}}
                PopoverProps={{
                    onEntered: () => {
                        popoverOpen(true);
                        setOpen(true);
                    },
                    onExited: () => {
                        popoverOpen(false);
                        setOpen(false);
                    },
                }}
                margin="none"
                value={uglyConvertToLocalTime(selectedDate).format()}
                onChange={(date: moment.Moment) => {
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
                openTo={showDate ? 'date' : 'hours'}
            />
        );
    }
);
