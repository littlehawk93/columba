import React from "react"
import Stack from "@mui/material/Stack"
import PackageCard from "./PackageCard";

class PackageCardList extends React.Component
{
    render() {

        const { packages } = this.props;

        return (
            <Stack spacing={2}>
                {packages && packages.map((pkg) => {
                    return (<PackageCard item={pkg} key={"package-card-" + pkg.id} onPackageRemoved={this.props.onPackageRemoved} />);
                })}
            </Stack>
        );
    }
}

export default PackageCardList;