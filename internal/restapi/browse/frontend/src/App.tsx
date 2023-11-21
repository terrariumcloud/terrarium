
import React from 'react';
import './App.css';
import Main from './pages/main/Main';
import Browse from './pages/modules/browse/Browse';
import ModuleInfo from './pages/modules/moduleinfo/ModuleInfo';
import ModuleDetailDescription from './pages/modules/moduleinfo/ModuleDetailDescription';
import ModuleDetailVersions from './pages/modules/moduleinfo/ModuleDetailVersions';

import {
    createBrowserRouter,
    RouterProvider
} from "react-router-dom";
import Home from './components/home/Home';
import NotFound from './pages/notfound/NotFound';
import ServerError from './pages/servererror/ServerError'

const routes = [
    {
        path: "",
        element: <Home />,

        children: [
            {
                index: true,
                element: <Main />,
                errorElement: <ServerError />,
            },
            {
                path: "terraform-modules",
                errorElement: <ServerError />,
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
                    },
                ],
            },
            // Last route...
            {
                path: "*",
                element: <NotFound />,
            }
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
