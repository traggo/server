import * as React from 'react';
import {TextField} from '@mui/material';
import {parseRelativeTime} from '../utils/time';
import Typography from '@mui/material/Typography';
import useTimeout from '@rooks/use-timeout';

interface RelativeDateTimeSelectorProps {
    value: string;
    onChange: (value: string, valid: boolean) => void;
    type: 'startOf' | 'endOf';
    label?: string;
    disabled?: boolean;
    disableUnderline?: boolean;
    style?: React.CSSProperties;
    small?: boolean;
}

export const RelativeDateTimeSelector: React.FC<RelativeDateTimeSelectorProps> = ({
    value,
    onChange: setValue,
    type,
    style,
    label,
    small = false,
    disableUnderline = false,
    disabled = false,
}) => {
    const [errVisible, setErrVisible] = React.useState(false);
    const [error, setError] = React.useState('');
    const {start, stop} = useTimeout(() => setErrVisible(true), 200);

    const parsed = parseRelativeTime(value, type);
    let helperText: React.ReactNode;
    if (small) {
        helperText = undefined;
    } else if (errVisible) {
        helperText = (
            <Typography color={'secondary'} variant={'caption'}>
                {error}
            </Typography>
        );
    } else {
        helperText = <Typography variant={'caption'}>{!parsed.success ? '...' : parsed.value.format('llll')}</Typography>;
    }
    return (
        <TextField
            fullWidth
            style={style}
            value={value}
            disabled={disabled}
            InputProps={{disableUnderline}}
            onChange={(e) => {
                const newValue = e.target.value;
                const result = parseRelativeTime(newValue, type);
                setErrVisible(false);
                stop();
                if (!result.success) {
                    setError(result.error);
                    start();
                } else {
                    setError('');
                }
                setValue(newValue, result.success);
            }}
            error={error !== ''}
            helperText={helperText}
            label={label}
        />
    );
};
