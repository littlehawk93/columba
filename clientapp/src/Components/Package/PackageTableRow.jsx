import React from "react"
import { eventIsDelivered } from "../../API/EventAPI"
import { formatLocationString } from "../../API/LocationAPI"
import { darken, lighten } from "@mui/material/styles"
import PackageComponent from "./PackageComponent"
import TableRow from "@mui/material/TableRow"
import TableCell from "@mui/material/TableCell"
import ButtonGroup from "@mui/material/ButtonGroup"
import IconButton from "@mui/material/IconButton"
import Button from "@mui/material/Button"
import RefreshIcon from "@mui/icons-material/Refresh"
import DeleteIcon from "@mui/icons-material/Delete"
import Timestamp from "../General/Timestamp"

const getBackgroundColor = (color, mode) => mode == "dark" ? darken(color, 0.6) : lighten(color, 0.6);

const successStyle = {
    bgcolor: (theme) => getBackgroundColor(theme.palette.success.main, theme.palette.mode)
};

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

        const { item, latestEvent, refreshing, popoverShowing, onRemoveClick, onCopyTrackingNumberClick, onRefreshClick } = this.props;

        return (
            <TableRow sx={eventIsDelivered(latestEvent) ? successStyle : {}}>
                <TableCell>
                    {item.label}
                </TableCell>
                <TableCell>
                    <input id={item.id + "-tracking-number-clipboard"} type="hidden" value={item.tracking_number} />
                    <Button onClick={onCopyTrackingNumberClick} disabled={popoverShowing} variant="text">{item.tracking_number}</Button>
                </TableCell>
                <TableCell>
                    {latestEvent ? latestEvent.event_text : ""}
                </TableCell>
                <TableCell>
                    {latestEvent ? formatLocationString(latestEvent.location) : ""}
                </TableCell>
                <TableCell>
                    {latestEvent && (<Timestamp value={latestEvent.timestamp} />)}
                </TableCell>
                <TableCell>
                    <ButtonGroup>
                        <IconButton 
                                title="Refresh Tracking Data" 
                                disabled={refreshing}
                                onClick={onRefreshClick}
                            >
                                <RefreshIcon sx={refreshing ? refreshAnimation : null }/>
                            </IconButton>
                            <IconButton 
                                title="Remove Package" 
                                disabled={refreshing}
                                onClick={onRemoveClick}
                            >
                                <DeleteIcon />
                        </IconButton>
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