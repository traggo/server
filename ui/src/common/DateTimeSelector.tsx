import * as React from 'react';
import {InlineDateTimePicker} from 'material-ui-pickers';
import TextField, {TextFieldProps} from '@material-ui/core/TextField';
import * as moment from 'moment';
import {uglyConvertToLocalTime} from '../timespan/timeutils';

interface DateTimeSelectorProps {
    selectedDate: moment.Moment;
    onSelectDate: (date: moment.Moment) => void;
    showDate: boolean;
    label: string;
    popoverOpen?: (open: boolean) => void;
}

const TextFieldWithoutUnderline: React.FC<TextFieldProps> = ({InputProps, ...other}) => {
    const o: TextFieldProps['InputProps'] = {...InputProps, disableUnderline: true};
    return <TextField {...other} InputProps={o} />;
};

export const DateTimeSelector: React.FC<DateTimeSelectorProps> = ({
    selectedDate,
    onSelectDate,
    showDate,
    label,
    popoverOpen = () => {},
}) => {
    return (
        <InlineDateTimePicker
            title={selectedDate.format()}
            style={{width: showDate ? 185 : 95}}
            keyboard
            PopoverProps={{onEntered: () => popoverOpen(true), onExited: () => popoverOpen(false)}}
            variant="standard"
            margin="none"
            value={uglyConvertToLocalTime(selectedDate).format()}
            onChange={(date) => {
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
            TextFieldComponent={TextFieldWithoutUnderline}
            mask={
                showDate
                    ? [/\d/, /\d/, /\d/, /\d/, '/', /\d/, /\d/, '/', /\d/, /\d/, ' ', /\d/, /\d/, ':', /\d/, /\d/]
                    : [/\d/, /\d/, ':', /\d/, /\d/]
            }
        />
    );
};
