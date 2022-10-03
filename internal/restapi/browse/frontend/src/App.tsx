
import React from 'react';
import './App.css';
import Main from './layouts/main/Main';
import Browse from './layouts/modules/browse/Browse';
import ModuleInfo from './layouts/modules/moduleinfo/ModuleInfo';
import ModuleDetailDescription from './layouts/modules/moduleinfo/ModuleDetailDescription';
import ModuleDetailVersions from './layouts/modules/moduleinfo/ModuleDetailVersions';

import {
    createBrowserRouter,
    RouterProvider,
    Outlet
} from "react-router-dom";
import Header from './components/header/Header';
import Footer from './components/footer/Footer';
import { Container } from '@mui/material';


const Home = () => {
    return <>
        <Header />
        <Container style={{marginTop: ".3em", marginBottom: ".85em"}}>
            <Outlet />
        </Container>
        <Footer />
    </>
}

const routes = [
    {
        path: "",
        element: <Home />,
        children: [
            {
                index: true,
                element: <Main />,
            },
            {
                path: "terraform-modules",
                children: [
                    {
                        index: true,
                        element: <Browse/>,
                    },
                    {
                        path:":org/:name/:provider",
                        element: <ModuleInfo />,
                        children: [
                            {
                                index: true,
                                element: <ModuleDetailDescription/>,
                            },
                            {
                                path: "description",
                                element: <ModuleDetailDescription />,
                            },
                            {
                                path: "versions",
                                element: <ModuleDetailVersions />,
                            }
                        ]
                    }
                ],
            },
        ],
    },
]


// import your route components too
function App() {
  const router = createBrowserRouter(routes)


  return (
      <RouterProvider router={router} />
  );
}

export default App;
