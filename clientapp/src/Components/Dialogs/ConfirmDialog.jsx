import React from "react"
import Dialog from "@mui/material/Dialog"
import DialogTitle from "@mui/material/DialogTitle"
import DialogContent from "@mui/material/DialogContent"
import DialogContentText from "@mui/material/DialogContentText"
import DialogActions from "@mui/material/DialogActions"
import Button from "@mui/material/Button"

class ConfirmDialog extends React.Component
{
    constructor(props) {
        super(props);
    }

    render() {

        const { open, title, message, onConfirm, onCancel } = this.props;

        return (
            <Dialog
                open={open}
                onClose={onCancel}
            >
                <DialogTitle>{title}</DialogTitle>
                <DialogContent>
                    <DialogContentText>{message}</DialogContentText>
                </DialogContent>
                <DialogActions>
                    <Button onClick={onCancel} autoFocus>Cancel</Button>
                    <Button onClick={onConfirm}>Confirm</Button>
                </DialogActions>
            </Dialog>
        );
    }
}

export default ConfirmDialog;