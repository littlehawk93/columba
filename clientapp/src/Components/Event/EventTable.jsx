import React from "react"
import Table from "@mui/material/Table"
import TableHead from "@mui/material/TableHead"
import TableBody from "@mui/material/TableBody"
import TableRow from "@mui/material/TableRow"
import TableCell from "@mui/material/TableCell"
import Moment from "react-moment";

class EventRow extends React.Component {

    formatLocationString = (location) => {

        if(location) {

            var result = `${location.city}, ${location.state}, ${location.zip}`;
            return result.replace(/(^(\s|,))|((\s|,)$)/, "");
        }
        return "";
    }

    render() {

        const { event } = this.props;

        return (
            <TableRow>
                <TableCell>
                    {this.formatLocationString(event.location)}
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