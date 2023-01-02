import React from "react"
import Table from "@mui/material/Table"
import TableHead from "@mui/material/TableHead"
import TableRow from "@mui/material/TableRow"
import TableCell from "@mui/material/TableCell"
import TableBody from "@mui/material/TableBody"
import PackageTableRow from "./PackageTableRow"

class PackageTable extends React.Component
{
    render() {

        const { packages } = this.props;

        return (
            <Table hover striped>
                <TableHead>
                    <TableRow>
                        <TableCell>Label</TableCell>
                        <TableCell>Tracking #</TableCell>
                        <TableCell>Last Event</TableCell>
                        <TableCell>Location</TableCell>
                        <TableCell>Time</TableCell>
                        <TableCell></TableCell>
                    </TableRow>
                </TableHead>
                <TableBody>
                    {packages && packages.map((pkg) => {
                        return (<PackageTableRow key={"package-table-row-" + pkg.id} item={pkg} />);
                    })}
                </TableBody>
            </Table>
        );
    }
}

export default PackageTable;