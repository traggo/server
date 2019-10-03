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
}

export const DateTimeSelector: React.FC<DateTimeSelectorProps> = ({
    selectedDate,
    onSelectDate,
    showDate,
    label,
    popoverOpen = () => {},
}) => {
    return (
        <KeyboardDateTimePicker
            variant="inline"
            InputProps={{disableUnderline: true}}
            title={selectedDate.format()}
            style={{width: showDate ? 185 : 95}}
            PopoverProps={{onEntered: () => popoverOpen(true), onExited: () => popoverOpen(false)}}
            margin="none"
            value={uglyConvertToLocalTime(selectedDate).format()}
            onChange={(date: moment.Moment) => {
                if (uglyConvertToLocalTime(selectedDate).isSame(date)) {
                    return;
                }

                onSelectDate(date);
            }}
            ampm={false}
            format={showDate ? 'YYYY/MM/DD HH:mm' : 'HH:mm'}
            label={label}
            openTo={showDate ? 'date' : 'hours'}
            onError={console.log}
        />
    );
};
