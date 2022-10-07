import {Typography, Card, CardActionArea, CardMedia, CardContent, Button, Link, Stack} from "@mui/material";
import React from "react";
import notFoundPicture1 from './not-found-1.jpg'
import {Link as RouterLink} from "react-router-dom";





const NotFound = () => {
    return(
        <>
            <Card>
                <CardActionArea>
                    <CardMedia
                        component="img"
                        image={notFoundPicture1}
                        alt="Picture taken from https://pxhere.com/en/photo/887567"
                    />
                    <CardContent>
                        <Typography gutterBottom variant="h5" component="div">
                            Looks like you took a wrong turn...
                        </Typography>
                        <Typography variant="body2" color="text.secondary">
                            The free high-resolution photo of forest, wood, trail, bridge, overgrown, walkway, waterway, nowhere, footbridge, habitat, vermont, dead end, intervale
                            , taken with an SCH-I545 02/28 2017 The picture taken with 5.0mm, f/2.2s, 1/692s, ISO 50
                            The image is released free of copyrights under Creative Commons CC0.
                            You may download, modify, distribute, and use them royalty free for anything you like, even in commercial applications. Attribution is not required.
                            It is was taken from <Link href="https://pxhere.com/en/photo/887567">PxHere</Link>.
                        </Typography>
                        <Stack direction="row-reverse" spacing="5">
                            <Button component={RouterLink} to="/terraform-modules">Browse for modules</Button>
                            <Button component={RouterLink} to="/">Home</Button>
                        </Stack>
                    </CardContent>
                </CardActionArea>
            </Card>
        </>
    )
}
export default NotFound