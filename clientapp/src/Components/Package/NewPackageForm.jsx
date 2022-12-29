import React from "react"
import { createPackage } from "../../API/PackageAPI"
import Box from "@mui/material/Box"
import Grid from "@mui/material/Grid"
import TextField from "@mui/material/TextField"
import ServiceSelect from "../Service/ServiceSelect"
import Fab from "@mui/material/Fab"
import Button from "@mui/material/Button"
import AddIcon from "@mui/icons-material/Add"
import Hidden from "@mui/material/Hidden"
import { ErrorContext } from "../../Context/Error"

class NewPackageForm extends React.Component
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

        var pkg = {
            label: labelValue,
            tracking_number: trackingNumValue,
            service: serviceValue
        };

        this.setState({processing: true}, () => {
            createPackage(pkg, this.onSuccess, this.onError);
        })
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

        return (
            <Box component="form" autoComplete="off">
                <Grid container spacing={2}>
                    <Grid item xs={12} sm={5} md={3} lg={3}>
                        <TextField 
                            fullWidth
                            size="small"
                            error={this.state.label.error !== ""}
                            helperText={this.state.label.error}
                            id="new-package-label"
                            name="np-label"
                            value={this.state.label.value}
                            placeholder="e.g. Mom's Gift"
                            label="Label (Optional)"
                            variant="filled"
                            onChange={this.onChange}
                            disabled={this.state.processing}
                        />
                    </Grid>
                    <Grid item xs={12} sm={7} md={4} lg={4}>
                        <TextField 
                            fullWidth
                            size="small"
                            error={this.state.trackingNum.error !== ""}
                            helperText={this.state.trackingNum.error}
                            id="new-package-trackingnum"
                            name="np-trackingnum"
                            value={this.state.trackingNum.value}
                            placeholder="Enter Tracking Number"
                            label="Tracking Number"
                            variant="filled"
                            onChange={this.onChange}
                            disabled={this.state.processing}
                        />
                    </Grid>
                    <Grid item xs={12} sm={6} md={2} lg={2}>
                        <ServiceSelect 
                            fullWidth
                            size="small"
                            error={this.state.service.error !== ""}
                            helperText={this.state.service.error}
                            id="new-package-service"
                            name="np-service"
                            value={this.state.service.value}
                            label="Service"
                            variant="filled"
                            onChange={this.onChange}
                            disabled={this.state.processing}
                        />
                    </Grid>
                    <Grid item xs={12} sm={3} md={1} lg={1} alignItems="flex-end">
                        <Hidden mdDown>
                            <Fab 
                                onClick={this.onSubmit}
                                size="medium"
                                color="primary" 
                                aria-label="Add Package"
                                disabled={this.state.processing}
                            >
                                <AddIcon />
                            </Fab>
                        </Hidden>
                        <Hidden mdUp>
                            <Button
                                onClick={this.onSubmit}
                                size="large"
                                variant="contained"
                                color="primary"
                                aria-label="Add Package"
                                startIcon={<AddIcon />}
                                disabled={this.state.processing}
                            >
                                Add Package
                            </Button>
                        </Hidden>
                    </Grid>
                </Grid>
            </Box>
        );
    }
}

export default function FNewPackageForm(props) {

    return (
        <ErrorContext.Consumer>
            {error => <NewPackageForm onError={error.onError} {...props} />}
        </ErrorContext.Consumer>
    )
}