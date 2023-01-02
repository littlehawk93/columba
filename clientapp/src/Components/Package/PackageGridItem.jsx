import React from "react"
import { formatLocationString } from "../../API/LocationAPI"
import { eventIsDelivered } from "../../API/EventAPI"
import Grid from "@mui/material/Grid"
import Card from '@mui/material/Card'
import CardHeader from "@mui/material/CardHeader"
import CardContent from "@mui/material/CardContent"
import CardActions from "@mui/material/CardActions"
import Button from '@mui/material/Button'
import Typography from '@mui/material/Typography'
import IconButton from "@mui/material/IconButton"
import RefreshIcon from "@mui/icons-material/Refresh"
import DeleteIcon from "@mui/icons-material/Delete"
import Timestamp from "../General/Timestamp"
import PackageComponent from "./PackageComponent"
import EventsDialog from "../Event/EventsDialog"

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

class PackageGridItemBase extends React.Component
{
    constructor(props) {
        super(props);

        this.state = {
            showEventDialog: false,
        };
    }

    onEventTableShow = (e) => {
        this.setState({
            showEventDialog: true
        });
    }

    onEventTableHide = (e) => {
        this.setState({
            showEventDialog: false
        });
    }

    render() {
        const { item, latestEvent, refreshing, popoverShowing, onRemoveClick, onCopyTrackingNumberClick, onRefreshClick } = this.props;

        const { showEventDialog } = this.state;

        return (
            <Grid item xs={12} md={6} lg={4}>
                <input id={item.id + "-tracking-number-clipboard"} type="hidden" value={item.tracking_number} />
                <Card variant="outlined" sx={{display: "flex", flexDirection: "column", justifyContent: "space-between", height: "100%"}}>
                    <CardHeader 
                        title={
                            <Typography variant="h5">{item.label ? item.label : item.tracking_number}</Typography>
                        }
                        action={
                            <div>
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
                            </div>
                        }
                    />
                    <CardContent sx={{flexGrow: 1}}>
                        <Button variant="text" onClick={onCopyTrackingNumberClick} disabled={popoverShowing}>{item.tracking_number}</Button>
                        {latestEvent && (
                            <Typography variant="h6">{latestEvent.event_text}</Typography>
                        )}
                        {latestEvent && (
                            <Typography variant="body1">{formatLocationString(latestEvent.location)}</Typography>
                        )}
                        {latestEvent && (
                            <Timestamp value={latestEvent.timestamp} />
                        )}
                    </CardContent>
                    <CardActions disableSpacing>
                        {item.tracking_url && (
                            <Button component="a" size="small" href={item.tracking_url} target="_blank">Track with {item.service}</Button>
                        )}
                        <Button sx={{marginLeft: "auto"}} disabled={refreshing} onClick={this.onEventTableShow}>Details</Button>
                    </CardActions>
                </Card>
                {showEventDialog && (<EventsDialog events={item.events} open={showEventDialog} onClose={this.onEventTableHide} />)} 
            </Grid>
        );
    }
}

export default function PackageGridItem(props) {

    return (
        <PackageComponent component={<PackageGridItemBase />} {...props} />
    )
}