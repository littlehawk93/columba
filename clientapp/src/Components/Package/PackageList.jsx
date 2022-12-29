import React from "react"
import Stack from "@mui/material/Stack"
import PackageCard from "./PackageCard";

class PackageList extends React.Component
{
    render() {

        const { packages } = this.props;

        return (
            <Stack spacing={2}>
                {packages && packages.map((pkg) => {
                    return (<PackageCard item={pkg} key={pkg.id} />);
                })}
            </Stack>
        );
    }
}

export default PackageList;