import Header from "../header/Header";
import { Container } from "@mui/material";
import { Outlet } from "react-router-dom";
import Footer from "../footer/Footer";
import React from "react";

const Home = () => {
    return <>
        <Header />
        <Container style={{ minHeight: "100vh", paddingBottom: "80px", marginTop: ".3em", marginBottom: ".85em" }}>
            <Outlet />
        </Container>
        <Footer />
    </>
}
export default Home

