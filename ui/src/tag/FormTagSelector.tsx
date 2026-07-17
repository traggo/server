import React from 'react';
import FormControl from '@mui/material/FormControl';
import Box from '@mui/material/Box';
import InputLabel from '@mui/material/InputLabel';
import {TagSelector, TagSelectorProps} from './TagSelector';

interface FormTagSelectorProps extends TagSelectorProps {
    label: string;
    required?: boolean;
}

export const FormTagSelector = ({label, required = false, ...props}: FormTagSelectorProps) => {
    return (
        <Box mt={1}>
            <FormControl fullWidth required={required}>
                <Box pt={2}>
                    <InputLabel shrink> {label} </InputLabel>
                    <Box className="MuiInput-underline">
                        <TagSelector {...props} />
                    </Box>
                </Box>
            </FormControl>
        </Box>
    );
};
