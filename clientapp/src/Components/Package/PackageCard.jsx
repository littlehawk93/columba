import React from "react"
import { formatLocationString } from "../../API/LocationAPI"
import { eventIsDelivered } from "../../API/EventAPI"
import Grid from "@mui/material/Grid"
import Card from '@mui/material/Card'
import CardHeader from "@mui/material/CardHeader"
import CardContent from "@mui/material/CardContent"
import CardActions from "@mui/material/CardActions"
import Collapse from "@mui/material/Collapse"
import Button from '@mui/material/Button'
import Typography from '@mui/material/Typography'
import IconButton from "@mui/material/IconButton"
import RefreshIcon from "@mui/icons-material/Refresh"
import DeleteIcon from "@mui/icons-material/Delete"
import ExpandMoreIcon from "@mui/icons-material/ExpandMore"
import EventTable from "../Event/EventTable"
import Timestamp from "../General/Timestamp"
import PackageComponent from "./PackageComponent"
import ContentCopyIcon from "@mui/icons-material/ContentCopy"

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

class PackageCardBase extends React.Component
{
    constructor(props) {
        super(props);

        this.state = {
            expanded: false,
        };
    }

    onToggleExpand = (e) => {

        var expanded = this.state.expanded;

        this.setState({
            expanded: !expanded
        });
    }

    render() {
        const { item, latestEvent, refreshing, popoverShowing, onRemoveClick, onCopyTrackingNumberClick, onRefreshClick } = this.props;
        const { expanded } = this.state;

        var title = item.label ? item.label : item.tracking_number;

        if (eventIsDelivered(latestEvent)) {
            title = title + " - Delivered";
        }

        const subtitle = item.label ? item.tracking_number + " (" + item.service + ")" : "";

        return (
            <Card variant="outlined">
                <input id={item.id + "-tracking-number-clipboard"} class="tracking-number-clipboard" type="text" value={item.tracking_number} />
                <CardHeader 
                    title={
                        <Typography variant="h5">{title}</Typography>
                    }
                    subheader={subtitle}
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
                {latestEvent && (
                    <CardContent>
                        <Grid container>
                            <Grid item xs={12}>
                                <Typography variant="h6">Latest Event</Typography>
                            </Grid>
                            <Grid item xs={12} sm={12} md={4}>
                                {formatLocationString(latestEvent.location)}
                            </Grid>
                            <Grid item xs={12} sm={12} md={4}>
                                {latestEvent.event_text}
                            </Grid>
                            <Grid item xs={12} sm={12} md={4}>
                                <Timestamp value={latestEvent.timestamp} />
                            </Grid>
                        </Grid>
                    </CardContent>
                )}
                <CardActions disableSpacing>
                    <Button startIcon={<ContentCopyIcon />} size="small" sx={{ marginRight: 2 }} onClick={onCopyTrackingNumberClick} disabled={popoverShowing}>Copy Tracking #</Button>
                    {item.tracking_url && (
                        <Button component="a" size="small" href={item.tracking_url} target="_blank">Track with {item.service}</Button>
                    )}
                    <IconButton onClick={this.onToggleExpand} sx={{ marginLeft: "auto" }}>
                        <ExpandMoreIcon sx={{transform: expanded ? "rotate(180deg)" : "rotate(0deg)"}} />
                    </IconButton>
                </CardActions>
                <Collapse in={expanded} unmountOnExit>
                    <CardContent>
                        <EventTable events={item.events} />
                    </CardContent>
                </Collapse>
            </Card>
        );
    }
}

export default function PackageCard(props) {

    return (
        <PackageComponent component={<PackageCardBase />} {...props} />
    )
}