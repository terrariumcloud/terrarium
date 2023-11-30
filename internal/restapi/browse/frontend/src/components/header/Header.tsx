import React from 'react';
import AppBar from '@mui/material/AppBar';
import Toolbar from '@mui/material/Toolbar';
import Button from '@mui/material/Button';
import { Link as RouterLink } from 'react-router-dom';
import { Typography } from '@mui/material';


function Header() {
  return (
    <AppBar position="static">
      <Toolbar>
        <Typography className='logo-title' color="inherit" component={RouterLink} to="/" style={{ flexGrow: 1, textTransform: "uppercase" }}>
          Terrarium
        </Typography>
        <Button color="inherit" component={RouterLink} to="/terraform-modules" style={{ textTransform: "capitalize" }}>
          Modules
        </Button>
        <Button color="inherit" component={RouterLink} to="/releases" style={{ textTransform: "capitalize" }}>
          Releases
        </Button>
      </Toolbar>
    </AppBar>
  );
}

export default Header;
