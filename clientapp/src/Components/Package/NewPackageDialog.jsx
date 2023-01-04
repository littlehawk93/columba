import React from "react"
import { createPackage } from "../../API/PackageAPI"
import Dialog from "@mui/material/Dialog"
import DialogTitle from "@mui/material/DialogTitle"
import DialogContent from "@mui/material/DialogContent"
import DialogActions from "@mui/material/DialogActions"
import Button from "@mui/material/Button"
import Box from "@mui/material/Box"
import Grid from "@mui/material/Grid"
import TextField from "@mui/material/TextField"
import ServiceSelect from "../Service/ServiceSelect"
import { ErrorContext } from "../../Context/Error"

class NewPackageDialogBase extends React.Component
{
    constructor(props) {
        super(props);

        this.state = {
            label: {
                error: "",
                value: ""
            },
            trackingNum: {
                error: "",
                value: ""
            },
            service: {
                error: "",
                value: ""
            },
            processing: false
        }
    }

    componentDidMount() {
        this.setState({
                label: {
                    error: "",
                    value: ""
                },
                trackingNum: {
                    error: "",
                    value: ""
                },
                service: {
                    error: "",
                    value: ""
                },
                processing: false
            }
        );
    }

    onChange = (e) => {

        const { value } = e.target;
        const { name } = e.target;

        if (name === "np-label") {
            this.setState({
                label: {
                    error: "",
                    value: value
                }
            });
        } else if (name === "np-trackingnum") {
            this.setState({
                trackingNum: {
                    error: value === "" ? "Tracking number cannot be blank" : "",
                    value: value
                }
            });
        } else if (name === "np-service") {
            this.setState({
                service: {
                    error: value === "" ? "You must select a Service Provider" : "",
                    value: value
                }
            });
        }
    }

    onClose = (e, reason) => {

        const { processing } = this.state;

        if (!processing && this.props.onClose) {
            this.props.onClose();
        }
    }

    onSubmit = (e) => {

        var labelValue = this.state.label.value.trim();
        var trackingNumValue = this.state.trackingNum.value.trim();
        var serviceValue = this.state.service.value.trim();

        if(trackingNumValue === "") {
            this.setState({
                trackingNum: {
                    error: "Tracking number cannot be blank",
                    value: trackingNumValue
                }
            });
            return;
        } else if (serviceValue === "") {
            this.setState({
                service: {
                    error: "You must select a Service Provider",
                    value: serviceValue
                }
            });
            return;
        }

        const pkg = {
            label: labelValue,
            tracking_number: trackingNumValue,
            service: serviceValue
        };

        this.setState({processing: true}, () => {
            createPackage(pkg, this.onSuccess, this.onError);
        });
    }

    onSuccess = (pkg) => {
        this.setState({
            label: {
                error: "",
                value: ""
            },
            trackingNum: {
                error: "",
                value: ""
            },
            service: {
                error: "",
                value: ""
            },
            processing: false
        }, () => {
            if (this.props.onPackageCreated) {
                this.props.onPackageCreated(pkg);
            }
        });
    }

    onError = (error) => {
        this.setState({processing: false}, () => {
            if(this.props.onError) {
                this.props.onError(error);
            }
        });
    }

    render() {

        const { label, trackingNum, service, processing } = this.state;
        const { open } = this.props;

        return (
            <Dialog
                open={open}
                onClose={this.onClose}
                disableEscapeKeyDown
            >
                <DialogTitle>Add New Package</DialogTitle>
                <DialogContent>
                    <Box component="form" autoComplete="off">
                        <Grid container spacing={4}>
                            <Grid item xs={12}>
                                <TextField 
                                    fullWidth
                                    size="small"
                                    error={label.error !== ""}
                                    helperText={label.error}
                                    id="new-package-label"
                                    name="np-label"
                                    value={label.value}
                                    placeholder="e.g. Mom's Gift"
                                    label="Label (Optional)"
                                    variant="filled"
                                    onChange={this.onChange}
                                    disabled={processing}
                                />
                            </Grid>
                            <Grid item xs={12}>
                                <TextField 
                                    fullWidth
                                    size="small"
                                    error={trackingNum.error !== ""}
                                    helperText={trackingNum.error}
                                    id="new-package-trackingnum"
                                    name="np-trackingnum"
                                    value={trackingNum.value}
                                    placeholder="Enter Tracking Number"
                                    label="Tracking Number"
                                    variant="filled"
                                    onChange={this.onChange}
                                    disabled={processing}
                                />
                            </Grid>
                            <Grid item xs={12}>
                                <ServiceSelect 
                                    fullWidth
                                    size="small"
                                    error={service.error !== ""}
                                    helperText={service.error}
                                    id="new-package-service"
                                    name="np-service"
                                    value={service.value}
                                    label="Service"
                                    variant="filled"
                                    onChange={this.onChange}
                                    disabled={processing}
                                />
                            </Grid>
                        </Grid>
                    </Box>
                </DialogContent>
                <DialogActions>
                    <Button disabled={processing} onClick={this.onClose} autoFocus>Cancel</Button>
                    <Button disabled={processing} onClick={this.onSubmit}>Add Package</Button>
                </DialogActions>
            </Dialog> 
        );
    }
}

export default function NewPackageDialog(props) {

    return (
        <ErrorContext.Consumer>
            {error => <NewPackageDialogBase onError={error.onError} {...props} />}
        </ErrorContext.Consumer>
    )
}