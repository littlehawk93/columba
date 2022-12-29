import React from "react"
import AppBar from "@mui/material/AppBar"
import Toolbar from "@mui/material/Toolbar"
import Typography from "@mui/material/Typography"
import Grid from "@mui/material/Grid"
import NewPackageForm from "./Components/Package/NewPackageForm"
import { ErrorContext } from "./Context/Error"

class ColumbaApp extends React.Component
{
    constructor(props) {
        super(props);

        this.state = { error: null };
    }

    onError = (error) => {
        this.setState({error: error});
    }

    render() {

        return (
            <ErrorContext.Provider value={{error: this.state.error, level: "error", onError: this.onError}}>
                <AppBar position="fixed" color="primary">
                    <Toolbar>
                        <Typography variant="h5" component="div" sx={{flexGrow: 1, textAlign: "center" }}>Columba Package Tracking</Typography>
                    </Toolbar>
                </AppBar>
                <Grid container spacing={2} sx={{ marginTop: "60px" }}>
                    <Grid item md={2} lg={3}></Grid>
                    <Grid item xs={12} md={8} lg={6}>
                        <NewPackageForm />
                    </Grid>
                    <Grid item md={2} lg={3}></Grid>
                </Grid>
            </ErrorContext.Provider>
        );
    }
}

export default ColumbaApp;