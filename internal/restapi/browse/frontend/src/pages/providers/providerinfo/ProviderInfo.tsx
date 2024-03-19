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
import {ProviderMetadata, useProviderMetadata} from "../../../data/useProviderMetadata";
import {styled} from "@mui/material/styles";

function ProviderMetadataInfo({provider}: { provider: ProviderMetadata }) {
    const ListItem = styled('li')(({ theme }) => ({
        margin: theme.spacing(0.5),
    }));
    return (
        <>
            <TableContainer component={Paper} style={{marginTop: ".85em"}}>
                <Table aria-label="provider metadata">
                    <TableBody>
                        <TableRow>
                            <TableCell component={"th"} scope={"row"}><b>Organization</b></TableCell>
                            <TableCell>{provider.organization}</TableCell>
                        </TableRow>
                        <TableRow>
                            <TableCell component={"th"} scope={"row"}><b>Name</b></TableCell>
                            <TableCell>{provider.name}</TableCell>
                        </TableRow>
                        <TableRow>
                            <TableCell component={"th"} scope={"row"}><b>Source Repository</b></TableCell>
                            <TableCell><Link href={provider.source_url}>{provider.source_url}</Link></TableCell>
                        </TableRow>
                        <TableRow>
                            <TableCell component={"th"} scope={"row"}><b>Description</b></TableCell>
                            <TableCell>{provider.description}</TableCell>
                        </TableRow>
                        <TableRow>
                            <TableCell component={"th"} scope={"row"}><b>Maturity</b></TableCell>
                            <TableCell>{provider.maturity}</TableCell>
                        </TableRow>
                        <TableRow>
                            <TableCell component={"th"} scope={"row"}><b>Available Versions</b></TableCell>
                            <TableCell>
                                <List component={"ul"}    sx={{
                                    display: 'flex',
                                    justifyContent: 'center',
                                    flexWrap: 'wrap',
                                    listStyle: 'none',
                                    p: 0.5,
                                    m: 0,
                                }}>
                                    {provider.versions.map((version, index) => {
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

function ProviderInfo() {
  const providerId = useParams();
  const provider = useProviderMetadata(providerId.org, providerId.name)
  return (
    <>
        <SearchBar />
        {provider && <ProviderMetadataInfo provider={provider}/>}
    </>
  );
}

export default ProviderInfo;
