import React from "react"
import { getLatestEvent } from "../../API/PackageAPI"
import { formatLocationString } from "../../API/LocationAPI"
import PackageComponent from "./PackageComponent"
import TableRow from "@mui/material/TableRow"
import TableCell from "@mui/material/TableCell"
import ButtonGroup from "@mui/material/ButtonGroup"
import IconButton from "@mui/material/IconButton"
import RefreshIcon from "@mui/icons-material/Refresh"
import DeleteIcon from "@mui/icons-material/Delete"
import Timestamp from "../General/Timestamp"

const refreshAnimation = {
    animation: "spin 2s linear infinite", 
    "@keyframes spin": {
        "0%": {
            transform: "rotate(0deg)",
        },
        "100%": {
            transform: "rotate(360deg);",
        },
    },
};

class PackageTableRowBase extends React.Component
{
    constructor(props) {
        super(props);

        this.state = { item: this.props.item };
    }

    render() {

        const { item } = this.state;
        const latestEvent = getLatestEvent(item);

        return (
            <TableRow>
                <TableCell>
                    {item.label}
                </TableCell>
                <TableCell>
                    {item.tracking_number}
                </TableCell>
                <TableCell>
                    {latestEvent.event_text}
                </TableCell>
                <TableCell>
                    {formatLocationString(latestEvent.location)}
                </TableCell>
                <TableCell>
                    <Timestamp value={latestEvent.timestamp} />
                </TableCell>
                <TableCell>
                    <ButtonGroup>
                        <IconButton><RefreshIcon /></IconButton>
                        <IconButton><DeleteIcon /></IconButton>
                    </ButtonGroup>
                </TableCell>
            </TableRow>
        );
    }
}

export default function PackageTableRow(props) {

    return (
        <PackageComponent component={<PackageTableRowBase />} {...props} />
    )
}