import React from "react"
import { getAllServiceProviders } from "../../API/ServiceAPI"
import { ErrorContext } from "../../Context/Error";
import TextField from "@mui/material/TextField"
import MenuItem from "@mui/material/MenuItem"

class ServiceSelect extends React.Component
{
    constructor(props) {
        super(props);

        this.state = { options: [] };
    }

    componentDidMount () {
        this.updateValues();
    }

    updateValues = () => {

        getAllServiceProviders((results) => {
            this.setState({options: results});
        }, (error) => {
            this.props.onError(error);
        })
    }

    render() {

        const { options } = this.state;

        return (
            <TextField select {...this.props}>
                {options && (options.map((option) => {
                    return (<MenuItem key={option} value={option}>{option}</MenuItem>);
                }))}
            </TextField>
        )
    }
}

export default function FServiceSelect(props) {

    return (
        <ErrorContext.Consumer>
            {error => <ServiceSelect onError={error.onError} {...props} />}
        </ErrorContext.Consumer>
    );
}