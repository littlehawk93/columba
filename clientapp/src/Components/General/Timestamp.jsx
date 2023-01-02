import React from "react"
import Moment from "react-moment"

const timestampFormat = "MMM D YYYY - h:mm A";

class Timestamp extends React.Component
{
    render() {
        const { value } = this.props;

        return (<Moment format={timestampFormat}>{value}</Moment>);
    }
}

export default Timestamp;