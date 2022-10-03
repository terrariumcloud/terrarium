import React from 'react';
import { Link } from '@mui/material';
import AppBar from '@mui/material/AppBar';
import Toolbar from '@mui/material/Toolbar';
import Button from '@mui/material/Button';
import { Link as RouterLink} from 'react-router-dom';


function Header() {
  return (
    <AppBar position="static">
    <Toolbar>
        <Button color="inherit" component={RouterLink} to="/">
        Terrarium
      </Button>
    </Toolbar>
  </AppBar>
  );
}

export default Header;
