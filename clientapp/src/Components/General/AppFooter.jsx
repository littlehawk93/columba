import React from "react"
import { useTheme } from "@mui/material/styles"
import { darken, lighten } from "@mui/material/styles"
import Grid from "@mui/material/Grid"
import Hidden from "@mui/material/Hidden"
import Typography from "@mui/material/Typography"
import IconButton from "@mui/material/IconButton"
import GitHubIcon from "@mui/icons-material/GitHub"


const getModifiedColor = (color, mode) => mode == "dark" ? darken(color, 0.6) : lighten(color, 0.6);

const toolbarStyles = {
    backgroundColor: (theme) => getModifiedColor(theme.palette.primary.light    , theme.palette.mode),
}

class AppFooterBase extends React.Component
{
    render() {
        const { theme } = this.props;
        return (
            <footer style={{
                display: "flex",
                position: "fixed",
                bottom: 0,
                width: "100%",
                backgroundColor: toolbarStyles.backgroundColor(theme),
                padding: "0px 24px",
            }}>
                <Grid container>
                    <Hidden mdDown>
                        <Grid item xs={4}>
                            <Typography variant="subtitle2" textAlign="left" sx={{ lineHeight: '60px' }}>Feedback? Visit Columba's <a style={{ color: theme.palette.text.primary }} href="https://github.com/littlehawk93/columba/issues" target="_blank">project page</a></Typography>
                        </Grid>
                        <Grid item xs={4}>
                            <Typography variant="subtitle2" textAlign="center" sx={{ lineHeight: '60px' }}>Created By <a style={{ color: theme.palette.text.primary }} href="https://github.com/littlehawk93" target="_blank">littlehawk93</a></Typography>
                        </Grid>
                        <Grid item xs={4}>
                            <Typography variant="subtitle2" textAlign="right" sx={{ lineHeight: '60px' }}>Like Columba? <a style={{ color: theme.palette.text.primary }} href="https://www.buymeacoffee.com/littlehawk93" target="_blank">Buy me a Coffee</a></Typography>
                        </Grid>
                    </Hidden>
                    <Hidden smDown mdUp>
                        <Grid item xs={4}>
                            <Typography variant="subtitle2" textAlign="left" sx={{ lineHeight: '60px' }}>Feedback? Visit Columba's <a style={{ color: theme.palette.text.primary }} href="https://github.com/littlehawk93/columba/issues" target="_blank">project page</a></Typography>
                        </Grid>
                        <Grid item xs={4}>
                            <Typography variant="subtitle2" textAlign="center" sx={{ lineHeight: '60px' }}>Created By <a style={{ color: theme.palette.text.primary }} href="https://github.com/littlehawk93" target="_blank">littlehawk93</a></Typography>
                        </Grid>
                        <Grid item xs={4}>
                            <Typography variant="subtitle2" textAlign="right" sx={{ lineHeight: '60px' }}>Like Columba? <a style={{ color: theme.palette.text.primary }} href="https://www.buymeacoffee.com/littlehawk93" target="_blank">Buy me a Coffee</a></Typography>
                        </Grid>
                    </Hidden>
                    <Hidden mdUp>
                        <Grid item xs={4}>
                            <IconButton href="https://github.com/littlehawk93/columba/issues" target="_blank">
                                <GitHubIcon />
                            </IconButton>
                        </Grid>
                        <Grid item xs={4}>
                            <Typography variant="subtitle2" textAlign="center" sx={{ lineHeight: '60px' }}>Created By <a style={{ color: theme.palette.text.primary }} href="https://github.com/littlehawk93" target="_blank">littlehawk93</a></Typography>
                        </Grid>
                        <Grid item xs={4}>
                            <Typography variant="subtitle2" textAlign="center" sx={{ lineHeight: '60px' }}>Like Columba? <a style={{ color: theme.palette.text.primary }} href="https://www.buymeacoffee.com/littlehawk93" target="_blank">Buy me a Coffee</a></Typography>
                        </Grid>
                    </Hidden>
                </Grid>
            </footer>
        );
    }
}

export default function AppFooter(props) {

    const theme = useTheme();
    return (<AppFooterBase theme={theme} />);
};