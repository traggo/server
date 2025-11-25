import React from 'react';
import FormControl from '@material-ui/core/FormControl';
import Box from '@material-ui/core/Box';
import InputLabel from '@material-ui/core/InputLabel';
import {TagSelector, TagSelectorProps} from './TagSelector';

interface FormTagSelectorProps extends TagSelectorProps {
    label: string;
    required?: boolean;
}

export const FormTagSelector = ({label, required = false, ...props}: FormTagSelectorProps) => {
    return (
        <Box mt={1}>
            <FormControl fullWidth required>
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
