import Chip from '@material-ui/core/Chip';
import * as React from 'react';
// @ts-ignore
import bestContrast from 'get-best-contrast-color';

export const TagChip = ({color, label}: {label: string; color: string}) => {
    const textColor = bestContrast(color, ['#fff', '#000']);
    return (
        <Chip
            tabIndex={-1}
            variant="outlined"
            style={{
                background: color,
                margin: '5px',
                color: textColor,
                cursor: 'text',
                minHeight: '32px',
                height: 'auto',
                whiteSpace: 'normal',
                wordBreak: 'break-word',
            }}
            label={label}
        />
    );
};
