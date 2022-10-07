import React from 'react';

import {Button, Grid, List, ListItem, Typography} from '@mui/material';
import modulesImage from './modules.png';
import providersImage from './providers.jpg';

import { Link as RouterLink } from 'react-router-dom';


function Main() {


  return (
      <>
          <Typography variant="h5" align='center' gutterBottom style={{marginTop: "1em"}}>About Terrarium</Typography>
          <Typography variant="body2" gutterBottom>
              Terrarium is an open source initiative enabling enterprise to build process and tooling enabling
              the adoption of Terraform in an enterprise environment.
          </Typography>

          <Typography variant="body2" gutterBottom>
              The philosophy for the project is to enable collaboration across team at an enterprise level, to promote
              best practices and integration with company governance covering but not limited to, secure development
              lifecycle, reporting and dependency management at scale across independent team.
          </Typography>

          <Grid container direction="row" alignItems="strech" spacing={5} style={{marginTop: "1em"}}>
              <Grid item xs={7} >
                  <Typography variant="h5" align='left' gutterBottom>Modules</Typography>
                  <Typography variant="body2" gutterBottom>
                      With Terrarium you can share module between multiple terraform code base. All the module that are
                      currently available are available for consumption to include a module you just need a few things:
                  </Typography>
                  <List>
                      <ListItem>You need to instantiate a module in your terraform code.</ListItem>
                      <ListItem>You need to reference the terrarium address for the module.</ListItem>
                      <ListItem>You need to specify the version of the module.</ListItem>
                  </List>
                  <Typography variant="body2" gutterBottom>
                      To discover the list of available module and their versions, Terrarium comes with module discovery
                      and search capabilities which we hopefully made friendly enough to be of use.
                  </Typography>
                  <Button component={RouterLink} to="terraform-modules">Browse for modules</Button>
              </Grid>
              <Grid item xs={5}>
                  <img src={modulesImage} width="100%" alt="Terraform modules illustration"/>
              </Grid>
          </Grid>
          
          <Grid container direction="row" alignItems="strech" spacing={5} style={{marginTop: "1em"}}>
              <Grid item xs={5}>
                  <img src={providersImage} width="100%" alt="Work in progress support for providers"/>
              </Grid>
              <Grid item xs={7} >
                  <Typography variant="h5" align='right' gutterBottom>Providers</Typography>
                  <Typography variant="body2" gutterBottom>
                      With Terrarium you will be able to publish and consume Terraform Provider. But it is not yet available.
                  </Typography>
              </Grid>
          </Grid>
      </>
  );
}

export default Main;
