import React from "react"
import { getActivePackages } from "./API/PackageAPI"
import AppBar from "@mui/material/AppBar"
import Toolbar from "@mui/material/Toolbar"
import Typography from "@mui/material/Typography"
import Grid from "@mui/material/Grid"
import NewPackageForm from "./Components/Package/NewPackageForm"
import PackageList from "./Components/Package/PackageList"
import { ErrorContext } from "./Context/Error"

class ColumbaApp extends React.Component
{
    constructor(props) {
        super(props);

        this.state = { error: null, packages: [] };
    }

    componentDidMount() {
        this.updatePackageList();
    }

    updatePackageList = () => {
        getActivePackages(this.onPackagesSuccess, this.onError);
    }

    onPackagesSuccess = (packages) => {
        this.setState({
            packages: packages
        });
    }

    onError = (error) => {
        this.setState({error: error});
    }

    onPackageCreated = (pkg) => {
        this.updatePackageList();
    }

    render() {

        const { packages } = this.state;

        return (
            <ErrorContext.Provider value={{error: this.state.error, level: "error", onError: this.onError}}>
                <AppBar position="fixed" color="primary">
                    <Toolbar>
                        <Typography variant="h5" component="div" sx={{flexGrow: 1, textAlign: "center" }}>Columba Package Tracking</Typography>
                    </Toolbar>
                </AppBar>
                <Grid container spacing={2} sx={{ marginTop: "60px" }}>
                    <Grid item md={2} lg={3}></Grid>
                    <Grid item xs={12} md={10} lg={8}>
                        <Grid container spacing={2}>
                            <Grid item xs={12}>
                                <NewPackageForm onPackageCreated={this.onPackageCreated} />
                            </Grid>
                            <Grid item xs={12}>
                                <PackageList packages={packages} />
                            </Grid>
                        </Grid>
                    </Grid>
                    <Grid item md={2} lg={3}></Grid>
                </Grid>
            </ErrorContext.Provider>
        );
    }
}

export default ColumbaApp;