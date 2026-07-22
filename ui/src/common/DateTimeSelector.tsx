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
        // The field fires onChange on every section edit (hour, minute, meridiem individually),
        // not just once a full value has been typed - if we forwarded those straight to
        // onSelectDate, a half-typed intermediate value could get persisted/trigger the parent's
        // before/after-the-other-field correction. So onChange only updates this local draft (for
        // responsive typing), and onSelectDate only fires once editing is actually done: on blur,
        // or on accepting a value from the popup calendar/clock.
        const [draft, setDraft] = React.useState<moment.Moment>(() => uglyConvertToLocalTime(selectedDate));
        const isEditing = React.useRef(false);

        React.useEffect(() => {
            if (!isEditing.current) {
                setDraft(uglyConvertToLocalTime(selectedDate));
            }
        }, [selectedDate]);

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

        const commit = (date: moment.Moment) => {
            if (!showDate && !open) {
                date = date.clone().set({
                    date: selectedDate.date(),
                    month: selectedDate.month(),
                    year: selectedDate.year(),
                });
            }
            if (uglyConvertToLocalTime(selectedDate).isSame(date)) {
                return;
            }
            onSelectDate(date);
        };

        return (
            <DesktopDateTimePicker
                className="time-picker"
                sx={{width: (showDate ? 185 : 105) + (ampm ? 20 : 0)}}
                onOpen={() => {
                    popoverOpen(true);
                    setOpen(true);
                    isEditing.current = true;
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
                        onFocus: () => {
                            isEditing.current = true;
                        },
                        onBlur: () => {
                            isEditing.current = false;
                            if (draft.isValid()) {
                                commit(draft);
                            }
                        },
                    },
                }}
                value={draft}
                onChange={(date: moment.Moment | null) => {
                    if (!date || !date.isValid()) {
                        return;
                    }
                    setDraft(date);
                }}
                onAccept={(date: moment.Moment | null) => {
                    isEditing.current = false;
                    if (date && date.isValid()) {
                        commit(date);
                    }
                }}
                ampm={ampm}
                format={format}
                label={label}
                openTo={showDate ? 'day' : 'hours'}
            />
        );
    }
);
