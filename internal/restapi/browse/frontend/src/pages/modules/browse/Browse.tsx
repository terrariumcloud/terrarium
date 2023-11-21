import React from 'react';
import SearchBar from '../../../components/search-bar/SearchBar';
import {Button, Card, CardContent, CardActions, Stack, Typography, Paper} from '@mui/material';
import { useFilteredModuleList, ModuleEntry } from '../../../data/useFilteredModuleList'
import { Link as RouterLink } from 'react-router-dom';

function Browse() {
    const [filteredModuleList, filterText, setFilterText] = useFilteredModuleList();

    const ModuleCard = ({module}:{module: ModuleEntry}) => {
        const modulePage: string = `${module.organization}/${module.name}/${module.provider}`
        return (
            <Card >
                <CardContent>
                <Typography color="text.primary">{module.organization || "Synamedia"} / {module.name} / {module.provider}</Typography>
                <Typography variant="body2">{module.description || "A module"}</Typography>
                </CardContent>
                <CardActions>
                    <Button size="small" href={module.source_url || ""}>Source</Button>
                    <Button size="small" component={RouterLink} to={modulePage}>Module Info</Button>
                </CardActions>
            </Card>
    )
    }

    return (
        <>
            <SearchBar filterValue={filterText} setFilter={setFilterText} />
            <Stack spacing={2} style={{marginTop: ".8em", marginBottom: ".8em"}}>
                <Paper>
                    <Typography variant="h5">Matching modules: {filteredModuleList.length}</Typography>
                </Paper>
                {filteredModuleList.map((mod, index) => { return <ModuleCard module={mod} key={index} /> })}
            </Stack>
        </>
    );
}

export default Browse;
