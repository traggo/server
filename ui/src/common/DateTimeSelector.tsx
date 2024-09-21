import * as React from 'react';
import {KeyboardDateTimePicker} from '@material-ui/pickers';
import * as moment from 'moment';
import {uglyConvertToLocalTime} from '../timespan/timeutils';

interface DateTimeSelectorProps {
    selectedDate: moment.Moment;
    onSelectDate: (date: moment.Moment) => void;
    showDate: boolean;
    label: string;
    popoverOpen?: (open: boolean) => void;
    style?: object;
}

export const DateTimeSelector: React.FC<DateTimeSelectorProps> = React.memo(
    ({selectedDate, onSelectDate, showDate, label, popoverOpen = () => {}, style = null}) => {
        const [open, setOpen] = React.useState(false);
        const localeData = moment.localeData();
        const time = localeData.longDateFormat('LT').replace('A', 'a');
        const ampm = time.indexOf('a') !== -1;
        const format = showDate ? localeData.longDateFormat('L') + ' ' + time : time;
        const width = (showDate ? 185 : 105) + (ampm ? 20 : 0);

        return (
            <KeyboardDateTimePicker
                variant="inline"
                InputProps={{disableUnderline: true}}
                title={selectedDate.format()}
                style={{minWidth: width, maxWidth: width, ...style}}
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
