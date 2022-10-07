import React from 'react';
import SearchBar from '../../../components/search-bar/SearchBar';
import {useParams} from 'react-router-dom';
import {
    Chip,
    Link,
    List,
    Paper,
    Table,
    TableBody,
    TableCell,
    TableContainer,
    TableRow
} from '@mui/material';
import {ModuleMetadata, useModuleMetadata} from "../../../data/useModuleMetadata";
import {styled} from "@mui/material/styles";

function ModuleMetadataInfo({module}: { module: ModuleMetadata }) {
    const ListItem = styled('li')(({ theme }) => ({
        margin: theme.spacing(0.5),
    }));
    return (
        <>
            <TableContainer component={Paper} style={{marginTop: ".85em"}}>
                <Table aria-label="module metadata">
                    <TableBody>
                        <TableRow>
                            <TableCell component={"th"} scope={"row"}><b>Organization</b></TableCell>
                            <TableCell>{module.organization}</TableCell>
                        </TableRow>
                        <TableRow>
                            <TableCell component={"th"} scope={"row"}><b>Name</b></TableCell>
                            <TableCell>{module.name}</TableCell>
                        </TableRow>
                        <TableRow>
                            <TableCell component={"th"} scope={"row"}><b>Provider</b></TableCell>
                            <TableCell>{module.provider}</TableCell>
                        </TableRow>
                        <TableRow>
                            <TableCell component={"th"} scope={"row"}><b>Source Repository</b></TableCell>
                            <TableCell><Link href={module.source_url}>{module.source_url}</Link></TableCell>
                        </TableRow>
                        <TableRow>
                            <TableCell component={"th"} scope={"row"}><b>Description</b></TableCell>
                            <TableCell>{module.description}</TableCell>
                        </TableRow>
                        <TableRow>
                            <TableCell component={"th"} scope={"row"}>Available Versions</TableCell>
                            <TableCell>
                                <List component={"ul"}    sx={{
                                    display: 'flex',
                                    justifyContent: 'center',
                                    flexWrap: 'wrap',
                                    listStyle: 'none',
                                    p: 0.5,
                                    m: 0,
                                }}>
                                    {module.versions.map((version, index) => {
                                        return (
                                            <ListItem key={index}>
                                                <Chip label={version} />
                                            </ListItem>
                                        )
                                    })}
                                </List>
                            </TableCell>
                        </TableRow>
                    </TableBody>
                </Table>
            </TableContainer>
        </>
    );
}

function ModuleInfo() {
  const moduleId = useParams();
  const module = useModuleMetadata(moduleId.org, moduleId.name, moduleId.provider)
  return (
    <>
        <SearchBar />
        {module && <ModuleMetadataInfo module={module}/>}
    </>
  );
}

export default ModuleInfo;
