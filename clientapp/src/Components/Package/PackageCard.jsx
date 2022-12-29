import React from "react"
import Card from '@mui/material/Card';
import CardActions from '@mui/material/CardActions';
import CardContent from '@mui/material/CardContent';
import Button from '@mui/material/Button';
import Typography from '@mui/material/Typography';

class PackageCard extends React.Component
{
    render() {

        const { item } = this.props;

        return (
            <Card variant="outlined">
                <CardContent>
                    {item.label && (<Typography gutterBottom variant="h5" component="div">{item.label}</Typography>)}
                    <Typography variant="body2" color="text.secondary">{item.tracking_number} ({item.service})</Typography>
                </CardContent>
                <CardActions>
                    {item.tracking_url && (<Button component="a" size="small" href={item.tracking_url} target="_blank">Track with {item.service}</Button>)}
                </CardActions>
            </Card>
        );
    }
}

export default PackageCard;