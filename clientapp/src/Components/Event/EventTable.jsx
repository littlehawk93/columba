import React from "react"
import { eventIsDelivered } from "../../API/EventAPI"
import { formatLocationString } from "../../API/LocationAPI"
import { darken, lighten } from "@mui/material/styles"
import Table from "@mui/material/Table"
import TableHead from "@mui/material/TableHead"
import TableBody from "@mui/material/TableBody"
import TableRow from "@mui/material/TableRow"
import TableCell from "@mui/material/TableCell"
import Moment from "react-moment"

const getBackgroundColor = (color, mode) => mode == "dark" ? darken(color, 0.6) : lighten(color, 0.6);

const successStyle = {
    bgcolor: (theme) => getBackgroundColor(theme.palette.success.main, theme.palette.mode)
};

const warningStyle = {
    bgcolor: (theme) => getBackgroundColor(theme.palette.warning.main, theme.palette.mode)
};

const errorStyle = {
    bgcolor: (theme) => getBackgroundColor(theme.palette.error.main, theme.palette.mode)
};

class EventRow extends React.Component {

    render() {

        const { event } = this.props;

        return (
            <TableRow sx={eventIsDelivered(event) ? successStyle : {}}>
                <TableCell>
                    {formatLocationString(event.location)}
                </TableCell>
                <TableCell>
                    {event.event_text}
                </TableCell>
                <TableCell>
                    <Moment format="MMM D YYYY - h:mm A">{event.timestamp}</Moment>
                </TableCell>
            </TableRow>
        )
    }
}

class EventTable extends React.Component {

    render() {

        const { events } = this.props;

        return (
            <Table>
                <TableHead>
                    <TableRow>
                        <TableCell width="20%">Location</TableCell>
                        <TableCell width="60%">Tracking Event</TableCell>
                        <TableCell width="20%">Date / Time</TableCell>
                    </TableRow>
                </TableHead>
                <TableBody>
                    {events && events.map((event) => {
                        return (<EventRow event={event} key={event.id} />);
                    })}
                </TableBody>
            </Table>
        );
    }
}

export default EventTable;