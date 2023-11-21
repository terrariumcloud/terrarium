import React from 'react';

import { Button, Card, CardContent, CardActions, Stack, Typography, Paper, FormControl, InputLabel, Select, MenuItem, SelectChangeEvent, IconButton, ListItemText, Checkbox } from '@mui/material';
import { Timeline, TimelineItem, TimelineSeparator, TimelineConnector, TimelineContent, TimelineDot, TimelineOppositeContent, timelineOppositeContentClasses } from '@mui/lab';
import RefreshRoundedIcon from '@mui/icons-material/RefreshRounded';

import { useFilteredReleaseList, ReleaseEntry, ReleaseLinks } from '../../../data/useReleasesList'
import GenericSearchBar from '../../../components/search-bar/GenericSearchBar';

const typesList = [
    'Oliver Hansen',
    'Van Henry',
    'April Tucker',
    'Ralph Hubbard',
    'Virginia Andrews',
    'Kelly Snyder',
];

const orgsList = [
    'cie', 'spvss', 'Miriam Wagner',
    'Bradley Wilkerson',
    'Virginia Andrews',
    'Kelly Snyder',
];

const ITEM_HEIGHT = 48;

const MenuProps = {
    PaperProps: {
        style: {
            maxHeight: ITEM_HEIGHT * 5.5,
            maxWidth: 250,
        },
    },
};

function Releases() {
    const [filteredModuleList, filterText, setFilterText] = useFilteredReleaseList();

    const [type, setType] = React.useState<string[]>([]);
    const [org, setOrg] = React.useState<string[]>([]);

    const handleTypeChange = (event: SelectChangeEvent<typeof type>) => {
        const {
            target: { value },
        } = event;
        setType(
            // On autofill we get a stringified value.
            typeof value === 'string' ? value.split(',') : value,
        );
    };

    const handleOrgChange = (event: SelectChangeEvent<typeof org>) => {
        const {
            target: { value },
        } = event;
        setOrg(
            // On autofill we get a stringified value.
            typeof value === 'string' ? value.split(',') : value,
        );
    };

    const LinkButton = ({ linkObj }: { linkObj: ReleaseLinks }) => {
        let domain, disabled = false
        try {
            if (!/^https?:\/\//i.test(linkObj.Url)) linkObj.Url = 'http://' + linkObj.Url;

            domain = (new URL(linkObj.Url)).hostname.replace('www.', '')
        } catch {
            domain = "Invalid-Link"
            disabled = true
        }
        return (<Button size="small" href={linkObj.Url} disabled={disabled}>{linkObj.Title ? linkObj.Title : domain}</Button>)
    }

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
                            <Typography variant='body2'>{`Version: ${module.version}`}</Typography>
                            {module.description && <><br /><Typography variant='body2'>{`${module.description}`}</Typography></>}
                        </CardContent>
                        <CardActions>
                            {module.links?.length &&
                                module.links.map((linkObject, index) => { return <LinkButton linkObj={linkObject} key={index} /> })}
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

                {/* Filters Section */}
                <Paper sx={{ display: 'flex', flexDirection: 'column', flexWrap: 'wrap' }}>
                    <div className="headingcolor">
                        Filter Releases
                    </div>
                    <div className="flex wrap" style={{ justifyContent: 'center', alignItems: "baseline", paddingBottom: "8px", backgroundColor: "rgba(66, 165, 245, 0.15)" }}>
                        <GenericSearchBar filterValue={filterText} setFilter={setFilterText} />
                        <FormControl variant="standard" sx={{ m: 1, minWidth: 120, maxWidth: 300 }}>
                            <InputLabel id="type-filter-label">Type</InputLabel>
                            <Select
                                labelId="type-filter-label"
                                value={type}
                                label="Type"
                                onChange={handleTypeChange}
                                multiple
                                renderValue={(selected) => selected.join(', ')}
                                MenuProps={MenuProps}
                            >
                                {typesList.map((typeEntry) => (
                                    <MenuItem key={typeEntry} value={typeEntry}>
                                        <Checkbox checked={type.indexOf(typeEntry) > -1} />
                                        <ListItemText primary={typeEntry} />
                                    </MenuItem>
                                ))}
                            </Select>
                        </FormControl>
                        <FormControl variant="standard" sx={{ m: 1, minWidth: 180, maxWidth: 300 }}>
                            <InputLabel id="org-filter-label">Organization</InputLabel>
                            <Select
                                labelId="org-filter-label"
                                value={org}
                                label="Org.."
                                onChange={handleOrgChange}
                                multiple
                                renderValue={(selected) => selected.join(', ')}
                                MenuProps={MenuProps}
                            >
                                {orgsList.map((orgEntry) => (
                                    <MenuItem key={orgEntry} value={orgEntry}>
                                        <Checkbox checked={org.indexOf(orgEntry) > -1} />
                                        <ListItemText primary={orgEntry} />
                                    </MenuItem>
                                ))}
                            </Select>
                        </FormControl>
                        <IconButton sx={{
                            color: 'white',
                            backgroundColor: '#1976d27d',
                            '&:hover': {
                                backgroundColor: "#1976d2",
                            }
                        }} aria-label="delete">
                            <RefreshRoundedIcon />
                        </IconButton>
                    </div>
                </Paper>

                {/* Release Timeline Section */}
                <Timeline
                    sx={{
                        [`& .${timelineOppositeContentClasses.root}`]: {
                            flex: 0.15,
                        },
                    }}
                >
                    {filteredModuleList.map((mod, index) => { return <ReleaseCard module={mod} key={index} /> })}
                </Timeline>
            </Stack>
        </>
    );
}

export default Releases;
