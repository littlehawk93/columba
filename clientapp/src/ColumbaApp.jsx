import React from "react"
import { ThemeProvider, createTheme } from "@mui/material/styles"
import { CssBaseline, Hidden } from "@mui/material"
import { getActivePackages } from "./API/PackageAPI"
import { ErrorContext } from "./Context/Error"
import { ConfirmContext } from "./Context/Confirm"
import { PopoverContext } from "./Context/Popover"
import AppBar from "@mui/material/AppBar"
import Toolbar from "@mui/material/Toolbar"
import Typography from "@mui/material/Typography"
import Grid from "@mui/material/Grid"
import AppFooter from "./Components/General/AppFooter"
import NewPackageForm from "./Components/Package/NewPackageForm"
import PackageList from "./Components/Package/PackageCardList"
import PackageGrid from "./Components/Package/PackageGrid"
import PackageTable from "./Components/Package/PackageTable"
import ErrorSnackbar from "./Components/Snackbars/ErrorSnackbar"
import Select from "@mui/material/Select"
import MenuItem from "@mui/material/MenuItem"
import ToggleButtonGroup from "@mui/material/ToggleButtonGroup"
import ToggleButton from "@mui/material/ToggleButton"
import ViewStreamIcon from "@mui/icons-material/ViewStream"
import TableRowsIcon from "@mui/icons-material/TableRows"
import WindowIcon from "@mui/icons-material/Window"
import ConfirmDialog from "./Components/Dialogs/ConfirmDialog"
import Popover from "@mui/material/Popover"


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
            themeName: this.getTheme(),
            layoutName: this.getLayout(),
            popover: {
                show: false,
                anchorElem: null,
                text: "",
                vertical: "",
                horizontal: "",
                onClose: null,
                closeJob: null,
            },
            confirm: {
                show: false,
                title: "",
                message: "",
                onConfirm: null,
                onCancel: null,
            }
        };
    }

    componentDidMount() {
        this.updatePackageList();
    }

    componentWillUnmount() {

        const { popover } = this.state;

        if(popover.closeJob) {
            clearTimeout(popover.closeJob);
        }
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

    onPopoverOpen = (anchorElem, text, vertical, horizontal, onClose) => {
        this.setState({
            popover: {
                show: true,
                anchorElem: anchorElem,
                text: text,
                vertical: vertical,
                horizontal: horizontal,
                onClose: onClose,
                closeJob: setTimeout(() => {
                    this.onPopoverClose();
                }, 1200)
            }
        });
    }

    onPopoverClose = () => {
        
        const { onClose } = this.state.popover;
        
        this.setState({
            popover: {
                show: false,
                anchorElem: null,
                text: "",
                vertical: "",
                horizontal: "",
                onClose: null,
                closeJob: null,
            }
        }, () => {
            if (onClose) {
                onClose();
            }
        });
    }

    onConfirmDialogOpen = (title, message, onConfirm, onCancel) => {
        this.setState({
            confirm: {
                show: true,
                title: title,
                message: message,
                onConfirm: onConfirm,
                onCancel: onCancel,
            }
        });
    }

    onConfirmSuccess = () => {
        
        const { onConfirm } = this.state.confirm;

        this.setState({
            confirm: {
                show: false,
                title: "",
                message: "",
                onConfirm: null,
                onCancel: null,
            }
        }, () => {
            if (onConfirm) {
                onConfirm();
            }
        });
    }

    onConfirmCancel = () => {

        const { onCancel } = this.state.confirm;

        this.setState({
            confirm: {
                show: false,
                title: "",
                message: "",
                onConfirm: null,
                onCancel: null,
            }
        }, () => {
            if (onCancel) {
                onCancel();
            }
        });
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

    getLayout = () => {

        const layout = localStorage.getItem("layout-name");

        if(layout) {
            return layout;
        }
        return "cards";
    }

    setLayout = (layout) => {
        localStorage.setItem("layout-name", layout);
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

    onChangeLayout = (e, value) => {
        
        this.setLayout(value);

        this.setState({
            layoutName: value
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

        const { packages, theme, error, themeName, layoutName, popover, confirm } = this.state;

        return (
            <ThemeProvider theme={theme}>
                <CssBaseline />
                <ErrorContext.Provider value={{error: error, level: "error", onError: this.onError}}>
                    <PopoverContext.Provider value={{anchorElem: popover.anchorElem, text: popover.text, vertical: popover.vertical, horizontal: popover.horizontal, onShowPopover: this.onPopoverOpen}}>
                        <ConfirmContext.Provider value={{show: confirm.show, title: confirm.title, message: confirm.message, onConfirm: confirm.onConfirm, onCancel: confirm.onCancel, onShowConfirm: this.onConfirmDialogOpen}}>    
                            <AppBar position="fixed" color="primary">
                                <Toolbar>
                                    <Hidden smDown>
                                        <ToggleButtonGroup exclusive value={layoutName} onChange={this.onChangeLayout} aria-label="packages layout">
                                            <ToggleButton value="cards" aria-label="cards layout" title="Cards"><ViewStreamIcon sx={{color: theme.palette.common.white}} /></ToggleButton>
                                            <ToggleButton value="grid" aria-label="grid layout" title="Grid"><WindowIcon sx={{color: theme.palette.common.white}} /></ToggleButton>
                                            <ToggleButton value="table" aria-label="table layout" title="Table"><TableRowsIcon sx={{color: theme.palette.common.white}} /></ToggleButton>
                                        </ToggleButtonGroup>
                                    </Hidden>
                                    <Typography variant="h5" component="div" sx={{flexGrow: 1, textAlign: "center" }}>Columba Package Tracking</Typography>
                                    <Select
                                        color="secondary"
                                        sx={{color: theme.palette.common.white}}
                                        size="small"
                                        onChange={this.onChangeTheme}
                                        name="select-theme"
                                        id="select-theme-dropdown"
                                        variant="standard"
                                        value={themeName}
                                        >
                                        <MenuItem value="light">Light</MenuItem>
                                        <MenuItem value="dark">Dark</MenuItem>
                                    </Select>
                                </Toolbar>
                            </AppBar>
                            <Grid container spacing={2} sx={{ marginTop: "60px", marginBottom: "60px" }}>
                                <Grid item md={1} lg={2}></Grid>
                                <Grid item xs={12} md={10} lg={8}>
                                    <Grid container spacing={2}>
                                        <Grid item xs={12}>
                                            <NewPackageForm onPackageCreated={this.onPackageCreated} />
                                        </Grid>
                                        <Hidden smDown>
                                            <Grid item xs={12}>
                                                {layoutName === "cards" && (
                                                    <PackageList packages={packages} onPackageRemoved={this.onPackageRemoved} />
                                                )}
                                                {layoutName === "grid" && (
                                                    <PackageGrid packages={packages} onPackageRemoved={this.onPackageRemoved} />
                                                )}
                                                {layoutName === "table" && (
                                                    <PackageTable packages={packages} onPackageRemoved={this.onPackageRemoved} />
                                                )}
                                            </Grid>
                                        </Hidden>
                                        <Hidden smUp>
                                            <Grid item xs={12}>
                                                <PackageGrid packages={packages} onPackageRemoved={this.onPackageRemoved} />
                                            </Grid>
                                        </Hidden>
                                    </Grid>
                                </Grid>
                                <Grid item md={2} lg={3}></Grid>
                            </Grid>
                            <AppFooter />
                            {error && (<ErrorSnackbar open={error != null && error != ""} error={error} onClose={this.clearError} />)}
                            {confirm.show && (
                                <ConfirmDialog 
                                    open={confirm.show}
                                    title={confirm.title} 
                                    message={confirm.message} 
                                    onCancel={this.onConfirmCancel} 
                                    onConfirm={this.onConfirmSuccess} />
                            )}
                            {popover.show && (
                                <Popover
                                    open={popover.anchorElem != null}
                                    anchorEl={popover.anchorElem}
                                    anchorOrigin={{vertical: popover.vertical, horizontal: popover.horizontal}}
                                    onClose={this.onPopoverClose}
                                    >
                                        <Typography sx={{ p: 2 }}>{popover.text}</Typography>
                                </Popover>
                            )}
                        </ConfirmContext.Provider>
                    </PopoverContext.Provider>
                </ErrorContext.Provider>
            </ThemeProvider>
        );
    }
}

export default ColumbaApp;