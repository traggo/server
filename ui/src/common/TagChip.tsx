import Chip from '@material-ui/core/Chip';
import * as React from 'react';
// @ts-ignore
import bestContrast from 'get-best-contrast-color';

export const TagChip = ({color, label, style={}}: {label: string; color: string, style?:React.CSSProperties}) => {
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
                height: 'fit-content',
                whiteSpace: 'normal',
                wordBreak: 'break-word',
                ... style
            }}
            label={label}
        />
    );
};
