import Chip from '@material-ui/core/Chip';
import * as React from 'react';
// @ts-ignore
import bestContrast from 'get-best-contrast-color';
import {makeStyles, Theme} from '@material-ui/core';

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
    onClick?: () => void;
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
