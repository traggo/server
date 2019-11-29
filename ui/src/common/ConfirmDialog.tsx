import Button from '@material-ui/core/Button';
import Dialog from '@material-ui/core/Dialog';
import DialogActions from '@material-ui/core/DialogActions';
import DialogContent from '@material-ui/core/DialogContent';
import DialogContentText from '@material-ui/core/DialogContentText';
import DialogTitle from '@material-ui/core/DialogTitle';
import React from 'react';

interface Props {
    title: string;
    fClose: () => void;
    fOnSubmit: () => void;
}

export const ConfirmDialog: React.FC<Props> = ({children, title, fClose, fOnSubmit}) => {
    const submitAndClose = () => {
        fOnSubmit();
        fClose();
    };
    return (
        <Dialog open={true} onClose={fClose} aria-labelledby="form-dialog-title" className="confirm-dialog">
            <DialogTitle id="form-dialog-title">{title}</DialogTitle>
            <DialogContent>
                <DialogContentText>{children}</DialogContentText>
            </DialogContent>
            <DialogActions>
                <Button onClick={fClose} className="cancel">
                    No
                </Button>
                <Button onClick={submitAndClose} autoFocus color="primary" variant="contained" className="confirm">
                    Yes
                </Button>
            </DialogActions>
        </Dialog>
    );
};
