import React from "react"
import { getPackageEvents } from "../../API/EventAPI"
import { getLatestEvent, deletePackage } from "../../API/PackageAPI"
import { ErrorContext } from "../../Context/Error"
import { ConfirmContext } from "../../Context/Confirm"
import { PopoverContext } from "../../Context/Popover"

const autoRefreshIntervalMillis = 300000; // 5 minutes in milliseconds

class PackageComponentBase extends React.Component
{
    constructor(props) {
        super(props);

        this.state = {
            refreshing: false,
            showClipboardPopover: false,
            item: props.item,
            refreshJob: null,
        };
    }

    componentDidMount() {
        this.setState({
            refreshJob: setTimeout(this.onRefresh, autoRefreshIntervalMillis)
        });
    }

    componentWillUnmount() {

        if (this.state.refreshJob) {
            clearTimeout(this.state.refreshJob);
            this.setState({refreshJob: null});
        }
    }

    onRefresh = (e) => {

        if (this.state.refreshJob) {
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

    onRemove = (e) => {
        if (this.props.onShowConfirm) {
            this.props.onShowConfirm("Remove Package?", "Remove package from listing? This action cannot be un-done.", this.onRemoveConfirm, null);
        }
    }

    onRemoveConfirm = (e) => {
        deletePackage(this.state.item.id, this.onRemoveSuccess, this.onRemoveError);
    }

    onRemoveSuccess = () => {
        if(this.props.onPackageRemoved) {
            this.props.onPackageRemoved(this.state.item);
        }
    }

    onRemoveError = (error) => {
        if(this.props.onError) {
            this.props.onError(error);
        }
    }

    onCopyTrackingNumber = (e) => {

        const { item } = this.state;

        var input = document.getElementById(item.id + "-tracking-number-clipboard");

        if (input) {
            const { value } = input;
            navigator.clipboard.writeText(value).then(() => {
                this.setState({
                    showClipboardPopover: true 
                }, () => {
                    if(this.props.onShowPopover) {
                        this.props.onShowPopover(e.target, "Copied to Clipboard", "center", "right", this.onClipboardPopoverClose);
                    }
                });
            }).catch(() => {
                try
                {
                    input.focus();
                    input.select();

                    document.execCommand("copy");

                    this.setState({
                        showClipboardPopover: true 
                    }, () => {
                        if(this.props.onShowPopover) {
                            this.props.onShowPopover(e.target, "Copied to Clipboard", "center", "right", this.onClipboardPopoverClose);
                        }
                    });
                }
                catch (err) {}
            });
        }
    }

    onClipboardPopoverClose = () => {
        this.setState({showClipboardPopover: false });
    }

    render() {

        const { component } = this.props;
        const { item, refreshing, showClipboardPopover } = this.state;

        return React.cloneElement(component, {
            onError: this.props.onError,
            item: item,
            latestEvent: getLatestEvent(item),
            onCopyTrackingNumberClick: this.onCopyTrackingNumber,
            onRefreshClick: this.onRefresh,
            onRemoveClick: this.onRemove,
            refreshing: refreshing,
            popoverShowing: showClipboardPopover,
        });
    }
}

export default function PackageComponent(props) {

    return (
        <ErrorContext.Consumer>
            {error => 
                <PopoverContext.Consumer>
                    {popover => 
                        <ConfirmContext.Consumer>
                            {confirm =>
                                <PackageComponentBase onError={error.onError} onShowConfirm={confirm.onShowConfirm} onShowPopover={popover.onShowPopover} {...props} />
                            }
                        </ConfirmContext.Consumer>
                    }
                </PopoverContext.Consumer>
            }
        </ErrorContext.Consumer>
    )
}