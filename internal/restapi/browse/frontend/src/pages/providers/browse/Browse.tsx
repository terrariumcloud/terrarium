import React from 'react';
import SearchBar from '../../../components/search-bar/SearchBar';
import { Button, Card, CardContent, CardActions, Stack, Typography, Paper } from '@mui/material';
import { useFilteredProviderList, ProviderEntry } from '../../../data/useFilteredProviderList'
import { Link as RouterLink } from 'react-router-dom';

function BrowseProviders() {
    const [filteredProviderList, filterText, setFilterText] = useFilteredProviderList();

    const ProviderCard = ({ provider }: { provider: ProviderEntry }) => {
        const providerPage: string = `${provider.organization}/${provider.name}`
        return (
            <Card >
                <CardContent>
                    <Typography color="text.primary">{provider.organization || "Synamedia"} / {provider.name}</Typography>
                    <Typography variant="body2">{provider.description || "A provider"}</Typography>
                </CardContent>
                <CardActions>
                    <Button size="small" href={provider.source_url || ""}>Source</Button>
                    <Button size="small" component={RouterLink} to={providerPage}>Provider Info</Button>
                </CardActions>
            </Card>
        )
    }

    return (
        <>
            <SearchBar filterValue={filterText} setFilter={setFilterText} />
            <Stack spacing={2} style={{ marginTop: ".8em", marginBottom: ".8em" }}>
                <Paper>
                    <Typography variant="h5">Matching providers: {filteredProviderList.length}</Typography>
                </Paper>
                {filteredProviderList.map((mod, index) => { return <ProviderCard provider={mod} key={index} /> })}
            </Stack>
        </>
    );
}

export default BrowseProviders;