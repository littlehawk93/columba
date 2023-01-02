import React from "react"
import Dialog from "@mui/material/Dialog"
import DialogTitle from "@mui/material/DialogTitle"
import DialogContent from "@mui/material/DialogContent"
import DialogActions from "@mui/material/DialogActions"
import Button from "@mui/material/Button"
import EventTable from "./EventTable"

class EventsDialog extends React.Component
{
    render() {

        const { open, events, onClose } = this.props;

        return (
            <Dialog
                open={open}
                onClose={onClose}
                disableEscapeKeyDown
            >
                <DialogTitle>Tracking Events</DialogTitle>
                <DialogContent>
                    <EventTable events={events} />
                </DialogContent>
                <DialogActions>
                    <Button onClick={onClose} autoFocus>Close</Button>
                </DialogActions>
            </Dialog>
        );
    }
}

export default EventsDialog;