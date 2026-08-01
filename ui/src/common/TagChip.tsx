import Chip from '@mui/material/Chip';
import * as React from 'react';
// @ts-expect-error no type declarations for this package
import bestContrast from 'get-best-contrast-color';
import {Theme} from '@mui/material';

import makeStyles from '@mui/styles/makeStyles';

const useStyles = makeStyles((theme: Theme) => ({
    chip: {
        margin: theme.spacing(0.5, 0.6),
        cursor: 'text',
        minHeight: '32px',
        height: 'fit-content',
        whiteSpace: 'normal',
        wordBreak: 'break-word',
    },
}));

interface TagChipProps {
    label: string;
    color: string;
    onClick?: (e: React.MouseEvent) => void;
}

export const TagChip: React.FC<TagChipProps> = ({color, label, onClick}) => {
    const classes = useStyles();
    const textColor = bestContrast(color, ['#fff', '#000']);
    return (
        <Chip
            tabIndex={-1}
            variant="outlined"
            className={classes.chip}
            style={{background: color, color: textColor}}
            label={label}
            onClick={onClick}
        />
    );
};
