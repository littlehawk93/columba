import React from "react"
import { getPackageEvents } from "../../API/EventAPI"
import { deletePackage } from "../../API/PackageAPI"
import { ErrorContext } from "../../Context/Error"
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
import ConfirmDialog from "../Dialogs/ConfirmDialog"

const autoRefreshIntervalMillis = 300000; // 5 minutes in milliseconds

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

class PackageCard extends React.Component
{
    constructor(props) {
        super(props);

        this.state = {
            showRemoveConfirmDialog: false,
            refreshing: false,
            expanded: false,
            item: props.item,
            refreshJob: setTimeout(this.onRefresh, autoRefreshIntervalMillis)
        };
    }

    componentWillUnmount() {

        if (this.state.refreshJob) {
            clearTimeout(this.state.refreshJob);
            this.setState({refreshJob: null});
        }
    }

    onRefresh = (e) => {

        if (e) {
            clearTimeout(this.state.refreshJob);
        }

        this.setState({
            refreshing: true,
            refreshJob: setTimeout(this.onRefresh, autoRefreshIntervalMillis)
        }, () => {
            getPackageEvents(this.state.item.id, this.onEventsReceived, this.onError);
        });
    }

    onEventsReceived = (events) => {

        var item = this.state.item;

        item.events = events;

        this.setState({
            item: item,
            refreshing: false
        });
    }

    onError = (error) => {

        this.setState({refreshing: false}, () => {
            if(this.props.onError) {
                this.props.onError(error);
            }
        });
    }

    onToggleExpand = (e) => {

        var expanded = this.state.expanded;

        this.setState({
            expanded: !expanded
        });
    }

    onRemove = (e) => {

        this.setState({
            showRemoveConfirmDialog: true
        });
    }

    onRemoveCancel = (e) => {

        this.setState({
            showRemoveConfirmDialog: false
        });
    }

    onRemoveConfirm = (e) => {
        this.setState({
            showRemoveConfirmDialog: false
        }, () => {
            deletePackage(this.state.item.id, this.onRemoveSuccess, this.onRemoveError);
        });
    }

    onRemoveSuccess = (pkg) => {
        if(this.props.onRemove) {
            this.props.onRemove(this.state.item);
        }
    }

    onRemoveError = (error) => {
        if(this.props.onError) {
            this.props.onError(error);
        }
    }

    render() {

        const { item } = this.state;

        return (
            <Card variant="outlined">
                <CardHeader 
                    title={
                        <Typography variant="h5">{item.label ? item.label : item.tracking_number}</Typography>
                    }
                    subheader={
                        (item.label ? item.tracking_number : "") + " (" + item.service + ")"
                    }
                    action={
                        <div>
                            <IconButton 
                                title="Refresh Tracking Data" 
                                disabled={this.state.refreshing}
                                onClick={this.onRefresh}
                            >
                                <RefreshIcon sx={this.state.refreshing ? refreshAnimation : null }/>
                            </IconButton>
                            <IconButton 
                                title="Remove Package" 
                                disabled={this.state.refreshing}
                                onClick={this.onRemove}
                            >
                                <DeleteIcon />
                            </IconButton>
                        </div>
                    }
                />
                <CardActions>
                    {item.tracking_url && (<Button component="a" size="small" href={item.tracking_url} target="_blank">Track with {item.service}</Button>)}
                    <IconButton
                        sx={{marginLeft: "auto"}}
                        onClick={this.onToggleExpand}
                    >
                        <ExpandMoreIcon sx={{transform: this.state.expanded ? "rotate(180deg)" : "rotate(0deg)"}} />
                    </IconButton>
                </CardActions>
                <Collapse in={this.state.expanded} unmountOnExit>
                    <CardContent>
                        <EventTable events={item.events} />
                    </CardContent>
                </Collapse>
                <ConfirmDialog 
                    open={this.state.showRemoveConfirmDialog}
                    title="Remove Package" 
                    message="Remove package from listing? This action cannot be un-done." 
                    onCancel={this.onRemoveCancel} 
                    onConfirm={this.onRemoveConfirm} />
            </Card>
        );
    }
}

export default function FPackageCard(props) {

    return (
        <ErrorContext.Consumer>
            {error => <PackageCard onError={error.onError} {...props} />}
        </ErrorContext.Consumer>
    )
}