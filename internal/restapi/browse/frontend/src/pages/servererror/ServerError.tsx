import {Typography, Card, CardActionArea, CardMedia, CardContent, Button, Link, Stack} from "@mui/material";
import React from "react";
import serverErrorPicture1 from './server-error-1.jpg'
import {Link as RouterLink} from "react-router-dom";





const ServerError = () => {
    return(
        <>
            <Card>
                <CardActionArea>
                    <CardMedia
                        component="img"
                        image={serverErrorPicture1}
                        alt="Picture taken from https://pxhere.com/en/photo/1356678"
                    />
                    <CardContent>
                        <Typography gutterBottom variant="h5" component="div">
                            Terrarium is facing an issue please contact us...
                        </Typography>
                        <Typography variant="body2" color="text.secondary">
                            The free high-resolution photo of sunlight, morning, smoke, fire, glow, burn, brand, natural disaster, disaster, event, forest fire, waldsterben, wildfire, environmental protection, embers, conflagration, atmospheric phenomenon, geological phenomenon
                            , taken with an unknown camera 04/04 2017 The picture taken with
                            The image is released free of copyrights under Creative Commons CC0.
                            You may download, modify, distribute, and use them royalty free for anything you like, even in commercial applications. Attribution is not required.
                            It is was taken from <Link href="https://pxhere.com/en/photo/1356678">PxHere</Link>.
                        </Typography>
                        <Stack direction="row-reverse" spacing="5">
                            <Button component={RouterLink} to="/">Home</Button>
                        </Stack>
                    </CardContent>
                </CardActionArea>
            </Card>
        </>
    )
}
export default ServerError