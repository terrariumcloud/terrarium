import React from 'react';

import { Timeline, TimelineItem, TimelineSeparator, TimelineConnector, TimelineContent, TimelineDot, TimelineOppositeContent, timelineOppositeContentClasses } from '@mui/lab';
import { Button, Card, CardContent, CardActions, Stack, Typography, Paper, FormControl, InputLabel, Select, MenuItem, SelectChangeEvent, IconButton } from '@mui/material';
import { useFilteredReleaseList, ReleaseEntry } from '../../../data/useReleasesList'
import SearchRoundedIcon from '@mui/icons-material/SearchRounded';
import { Link as RouterLink } from 'react-router-dom';
import dummyRelease from '../../../assets/releases-list.json'

function Releases() {
    const [filteredModuleList, filterText, setFilterText] = useFilteredReleaseList();
    const [type, setType] = React.useState('');
    const [org, setOrg] = React.useState('');

    const handleTypeChange = (event: SelectChangeEvent) => {
        setType(event.target.value);
    };
    const handleOrgChange = (event: SelectChangeEvent) => {
        setOrg(event.target.value);
    };

    const ReleaseCard = ({ module }: { module: ReleaseEntry }) => {
        return (
            <TimelineItem>
                <TimelineOppositeContent color="textSecondary">
                    {new Date(module.createdAt).toLocaleString('en-US', {
                        weekday: 'short', // long, short, narrow
                        day: 'numeric', // numeric, 2-digit
                        month: 'long', // numeric, 2-digit, long, short, narrow
                        hour: 'numeric', // numeric, 2-digit
                        minute: 'numeric', // numeric, 2-digit
                    })}
                </TimelineOppositeContent>
                <TimelineSeparator>
                    <TimelineDot />
                    <TimelineConnector />
                </TimelineSeparator>
                <TimelineContent>
                    <Card >
                        <CardContent>
                            <Typography color="text.primary" variant='h6' display={'inline'}>{module.name}</Typography> 
                            <Typography color="text.secondary" variant='subtitle1' display={'inline'}>{` | ${module.Organization} | ${module.type}`}</Typography>
                            <Typography variant='body2'>{`${module.version}`}</Typography>
                            <Typography variant='body2'>{`${module.description}`}</Typography>
                        </CardContent>
                        <CardActions>
                            <Button size="small" href={""}>Source</Button>
                            <Button size="small" component={RouterLink} to={'/'}>Module Info</Button>
                        </CardActions>
                    </Card>
                </TimelineContent>
            </TimelineItem>
        )
    }

    return (
        <>
            <div style={{ paddingTop: "10px" }}>
                <Typography color="text.primary" variant='h5'>Latest Releases to Terrarium</Typography>
            </div>
            <Stack spacing={2} style={{ marginTop: ".8em", marginBottom: ".8em" }}>
                <Paper sx={{ display: 'flex', flexDirection: 'column', flexWrap: 'wrap' }}>
                    <div className="headingcolor">
                        Filter Releases
                    </div>
                    <div className="flex wrap" style={{ justifyContent: 'center', alignItems: "baseline", paddingBottom: "8px" }}>
                        <FormControl variant="standard" sx={{ m: 1, minWidth: 120 }}>
                            <InputLabel id="type-filter-label">Type</InputLabel>
                            <Select
                                labelId="type-filter-label"
                                value={type}
                                label="Type"
                                autoWidth
                                onChange={handleTypeChange}
                            >
                                <MenuItem value="all">
                                    <em>All</em>
                                </MenuItem>
                                <MenuItem value="bundle">Bundles</MenuItem>
                                <MenuItem value="module">Modules</MenuItem>
                            </Select>
                        </FormControl>
                        <FormControl variant="standard" sx={{ m: 1, minWidth: 180 }}>
                            <InputLabel id="type-filter-label">Organization</InputLabel>
                            <Select
                                labelId="type-filter-label"
                                value={org}
                                label="Org.."
                                autoWidth
                                onChange={handleOrgChange}
                            >
                                <MenuItem value="all">
                                    <em>All</em>
                                </MenuItem>
                                <MenuItem value="cie">cie</MenuItem>
                                <MenuItem value="brooklyn">brooklyn</MenuItem>
                                <MenuItem value="mediacloud">mediacloud</MenuItem>
                            </Select>
                        </FormControl>
                        <IconButton sx={{
                            color: 'white',
                            backgroundColor: '#1976d27d',
                            '&:hover': {
                                backgroundColor: "#1976d2",
                            }
                        }} aria-label="delete">
                            <SearchRoundedIcon />
                        </IconButton>
                    </div>
                </Paper>
                <Timeline
                    sx={{
                        [`& .${timelineOppositeContentClasses.root}`]: {
                            flex: 0.15,
                        },
                    }}
                >
                    {dummyRelease.map((mod, index) => { return <ReleaseCard module={mod} key={index} /> })}
                </Timeline>
            </Stack>
        </>
    );
}

export default Releases;
