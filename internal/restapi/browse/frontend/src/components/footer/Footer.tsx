import React from 'react';
import { styled, alpha } from '@mui/material/styles';
import Typography from '@mui/material/Typography';

const StyledFooter = styled('footer')(({ theme }) => ({
  position: 'relative',
  backgroundColor: alpha(theme.palette.primary.light, 0.15),
  marginLeft: 0,
  width: '100%',
  textAlign: 'center',
  marginTop: "-80px",
  height: "80px",
  justifyContent: "center",
  display: "flex",
  flexDirection: "column",
  boxShadow: "0px -2px 4px -1px rgba(0,0,0,0.2), 0px -4px 5px 0px rgba(0,0,0,0.14), 0px -1px 10px 0px rgba(0,0,0,0.12)"
}));

function Footer() {
  return (
    <StyledFooter className="App-footer">
      <Typography variant="body1" color="inherit" noWrap>Terrarium is the place where you can grow your Terraform eco-system in house.</Typography>
      <Typography variant="h6" color="gray" noWrap style={{ fontSize: "0.75rem" }}>
        Provided by <a className="nodecorlink" href='https://www.synamedia.com/'>Synamedia</a>
      </Typography>
    </StyledFooter>
  );
}

export default Footer;
