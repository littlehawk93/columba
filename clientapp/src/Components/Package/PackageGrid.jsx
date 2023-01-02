import React from "react"
import Grid from "@mui/material/Grid"
import PackageGridItem from "./PackageGridItem"

class PackageGrid extends React.Component
{
    render() {
        const { packages } = this.props;

        return (
            <Grid container spacing={2}>
                {packages && packages.map((pkg) => {
                    return (<PackageGridItem key={"package-grid-item-" + pkg.id} item={pkg} onPackageRemoved={this.props.onPackageRemoved} />);
                })}
            </Grid>
        );
    }
}

export default PackageGrid;