import React from "react"
import Moment from "react-moment"

const timestampFormat = "MMM\u00a0D\u00a0YYYY h:mm\u00a0A";

class Timestamp extends React.Component
{
    render() {
        const { value } = this.props;
        return value ? (<Moment format={timestampFormat}>{value}</Moment>) : (<div></div>);
    }
}

export default Timestamp;