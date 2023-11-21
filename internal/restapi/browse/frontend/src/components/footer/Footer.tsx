import React from 'react';
import { styled, alpha } from '@mui/material/styles';
import Typography from '@mui/material/Typography';

const StyledFooter = styled('div')(({ theme }) => ({
  position: 'relative',
  backgroundColor: alpha(theme.palette.primary.light, 0.15),
  marginLeft: 0,
  width: '100%',
  textAlign: 'center',
}));

function Footer() {
  return (
    <StyledFooter className="App-footer">
      <Typography variant="body1"   color="inherit" noWrap>Terrarium is the place where you can grow your Terraform eco-system in house.</Typography>
      <Typography variant="h6" color="inherit" noWrap>
        Provided by Synamedia
      </Typography>
    </StyledFooter>
  );
}

export default Footer;
