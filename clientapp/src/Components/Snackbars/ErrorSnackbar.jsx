import React from "react"
import Snackbar from "@mui/material/Snackbar"
import Alert from "@mui/material/Alert"

class ErrorSnackbar extends React.Component
{
    render() {

        const { error, open, onClose } = this.props;

        return (
            <Snackbar
                anchorOrigin={{vertical: "bottom", horizontal: "center"}}
                open={open} 
                autoHideDuration={6000}
                onClose={onClose}
            >
                <Alert onClose={onClose} severity="error">{error}</Alert>
            </Snackbar>
        );
    }
}

export default ErrorSnackbar;