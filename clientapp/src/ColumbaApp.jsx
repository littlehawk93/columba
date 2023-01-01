import React from "react"
import { ThemeProvider, createTheme } from "@mui/material/styles"
import { CssBaseline } from "@mui/material"
import { getActivePackages } from "./API/PackageAPI"
import AppBar from "@mui/material/AppBar"
import Toolbar from "@mui/material/Toolbar"
import Typography from "@mui/material/Typography"
import Grid from "@mui/material/Grid"
import NewPackageForm from "./Components/Package/NewPackageForm"
import PackageList from "./Components/Package/PackageList"
import Select from "@mui/material/Select"
import MenuItem from "@mui/material/MenuItem"
import Snackbar from "@mui/material/Snackbar"
import Alert from "@mui/material/Alert"
import { ErrorContext } from "./Context/Error"

class ColumbaApp extends React.Component
{
    constructor(props) {
        super(props);

        this.state = { 
            error: null, 
            packages: [],
            theme: createTheme({
                palette: {
                    mode: this.getTheme()
                }
            }),
            themeName: this.getTheme()
        };
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

    clearError = () => {
        this.setState({error: null});
    }

    onPackageCreated = (pkg) => {
        var packages = this.state.packages;

        packages.push(pkg);

        this.setState({packages: packages});
    }

    getTheme = () => {

        const theme = localStorage.getItem("theme-name");

        if(theme) {
            return theme;
        }
        return "light";
    }

    setTheme = (theme) => {

        localStorage.setItem("theme-name", theme);
    }

    onChangeTheme = (e) => {

        const theme = e.target.value;

        this.setTheme(theme);

        this.setState({
            theme: createTheme({
                palette: {
                    mode: theme
                }
            }),
            themeName: theme
        });
    }

    onPackageRemoved = (pkg) => {

        var packages = this.state.packages;

        for(var i=0;i<packages.length;i++) {

            if (packages[i].id == pkg.id) {
                packages.splice(i, 1);
                break;
            }
        }

        this.setState({
            packages: packages
        });
    }

    render() {

        const { packages, theme, error, themeName } = this.state;

        return (
            <ThemeProvider theme={theme}>
                <CssBaseline />
                <ErrorContext.Provider value={{error: error, level: "error", onError: this.onError}}>
                    <AppBar position="fixed" color="primary">
                        <Toolbar>
                            <Typography variant="h5" component="div" sx={{flexGrow: 1, textAlign: "center" }}>Columba Package Tracking</Typography>
                            <Select
                                size="small"
                                onChange={this.onChangeTheme}
                                name="select-theme"
                                id="select-theme-dropdown"
                                variant="standard"
                                value={themeName}
                                >
                                <MenuItem value="light">Light Theme</MenuItem>
                                <MenuItem value="dark">Dark Theme</MenuItem>
                            </Select>
                        </Toolbar>
                    </AppBar>
                    <Grid container spacing={2} sx={{ marginTop: "60px" }}>
                        <Grid item md={1} lg={2}></Grid>
                        <Grid item xs={12} md={10} lg={8}>
                            <Grid container spacing={2}>
                                <Grid item xs={12}>
                                    <NewPackageForm onPackageCreated={this.onPackageCreated} />
                                </Grid>
                                <Grid item xs={12}>
                                    <PackageList packages={packages} onPackageRemoved={this.onPackageRemoved} />
                                </Grid>
                            </Grid>
                        </Grid>
                        <Grid item md={2} lg={3}></Grid>
                    </Grid>
                    {error && (
                        <Snackbar
                            anchorOrigin={{vertical: "bottom", horizontal: "center"}}
                            open={error != null && error != ""} 
                            autoHideDuration={6000}
                            onClose={this.clearError}
                        >
                            <Alert onClose={this.clearError} severity="error">{error}</Alert>
                        </Snackbar>
                    )}
                </ErrorContext.Provider>
            </ThemeProvider>
        );
    }
}

export default ColumbaApp;