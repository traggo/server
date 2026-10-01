import * as React from 'react';
import {TextField} from '@material-ui/core';
import {parseRelativeTime} from '../utils/time';
import Typography from '@material-ui/core/Typography';
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
    value: apiValue,
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
    const parsed = parseRelativeTime(apiValue, type);

    return (
        <TextField
            fullWidth
            style={style}
            value={parsed.success ? parsed.localized : apiValue}
            disabled={disabled}
            InputProps={{disableUnderline}}
            onChange={(e) => {
                const result = parseRelativeTime(e.target.value, type);
                setErrVisible(false);
                stop();
                if (!result.success) {
                    setError(result.error);
                    start();
                } else {
                    setError('');
                }
                setValue(result.success ? result.normalized : e.target.value, result.success);
            }}
            error={error !== ''}
            helperText={
                small ? (
                    undefined
                ) : errVisible ? (
                    <Typography color={'secondary'} variant={'caption'}>
                        {error}
                    </Typography>
                ) : (
                    <Typography variant={'caption'}>{!parsed.success ? '...' : parsed.preview.format('llll')}</Typography>
                )
            }
            label={label}
        />
    );
};
