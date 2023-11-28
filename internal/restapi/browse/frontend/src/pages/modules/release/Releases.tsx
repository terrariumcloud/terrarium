import React from 'react';

import { Button, Card, CardContent, CardActions, Stack, Typography, Paper, FormControl, InputLabel, Select, MenuItem, SelectChangeEvent, IconButton, ListItemText, Checkbox, CircularProgress } from '@mui/material';
import { Timeline, TimelineItem, TimelineSeparator, TimelineConnector, TimelineContent, TimelineDot, TimelineOppositeContent, timelineOppositeContentClasses } from '@mui/lab';
import RotateLeftIcon from '@mui/icons-material/RotateLeft';

import { useFilteredReleaseList, ReleaseEntry, ReleaseLinks } from '../../../data/useReleasesList'
import GenericSearchBar from '../../../components/search-bar/SimpleSearchBar';
import { useReleaseOrgList } from '../../../data/useReleaseOrgsList';
import { useReleaseTypeList } from '../../../data/useReleaseTypesList';

import '../../../styles/releases.css';

const ITEM_HEIGHT = 48;
const MenuProps = {
    PaperProps: {
        style: {
            maxHeight: ITEM_HEIGHT * 5.5,
            maxWidth: 250,
        },
    },
};

const timeFilters = [
    {
        label: "Hour",
        value: "1h"
    },
    {
        label: "Day",
        value: "24h"
    },
    {
        label: "Week",
        value: "168h"
    },
    {
        label: "Month",
        value: "730h"
    },
    {
        label: "3 Months",
        value: "2190h"
    },
    {
        label: "6 Months",
        value: "4380h"
    },
    {
        label: "Year",
        value: "8760h"
    },
]

function Releases() {
    const typesList = useReleaseTypeList()
    const orgsList = useReleaseOrgList()

    const [selectedTypes, setType] = React.useState<string[]>([]);
    const [selectedOrg, setOrg] = React.useState<string[]>([]);
    const [selectedTime, setTime] = React.useState('168h');

    const [filteredModuleList, filterText, setFilterText, isLoading] = useFilteredReleaseList(selectedTypes, selectedOrg, selectedTime);

    const handleTypeChange = (event: SelectChangeEvent<typeof selectedTypes>) => {
        const {
            target: { value },
        } = event;
        setType(
            // On autofill we get a stringified value.
            typeof value === 'string' ? value.split(',') : value,
        );
    };

    const handleOrgChange = (event: SelectChangeEvent<typeof selectedOrg>) => {
        const {
            target: { value },
        } = event;
        setOrg(
            // On autofill we get a stringified value.
            typeof value === 'string' ? value.split(',') : value,
        );
    };

    const handleTimeChange = (event: SelectChangeEvent<string>) => {
        setTime(event.target.value);
    };

    const resetFilters = () => {
        setType([])
        setOrg([])
        setFilterText("")
    }

    const LinkButton = ({ linkObj }: { linkObj: ReleaseLinks }) => {
        let domain, disabled = false
        try {
            if (linkObj.url) {
                domain = (new URL(linkObj.url)).hostname.replace('www.', '')
            } else {
                domain = "no-url"
                disabled = true
            }
        } catch {
            domain = "invalid-url"
            disabled = true
        }
        return (<Button variant="outlined" size="small" href={linkObj.url} disabled={disabled} style={{ margin: '8px' }}>
            <div className='word-wrapper'>
                {linkObj.title ? disabled ? `${linkObj.title} (${domain})` : linkObj.title : domain}
            </div>
        </Button>)
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
                            <Typography color="text.secondary" variant='subtitle1' display={'inline'}>{module.organization && ` | ${module.organization}`}{module.type && ` | ${module.type}`}</Typography>
                            <Typography variant='body2'>{`Version: ${module.version}`}</Typography>
                            {module.description && <><br /><Typography variant='body2'>{`${module.description}`}</Typography></>}
                        </CardContent>
                        <CardActions className='wrap'>
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
            <div className='flex' style={{ paddingTop: "10px" }}>
                <Typography color="text.primary" variant='h5' sx={{ flexGrow: '1' }}>
                    Latest Releases to Terrarium
                </Typography>
                <div className='flex' style={{ alignItems: 'baseline' }}>
                    <Typography color="grey" variant='body1' sx={{ p: '4px 0px 4px 0px' }}>
                        From the Past
                    </Typography>
                    <FormControl variant="standard" sx={{ ml: 1, minWidth: '80px' }}>
                        <Select
                            sx={{
                                color: "grey",
                                "& .MuiSvgIcon-root": {
                                    color: "grey",
                                },
                            }}
                            labelId="time-select-label"
                            id="time-select"
                            value={selectedTime}
                            label="Time"
                            onChange={handleTimeChange}
                        >
                            {timeFilters.map((timeEntry) => (<MenuItem value={timeEntry.value}>{timeEntry.label}</MenuItem>))}
                        </Select>
                    </FormControl>
                </div>
            </div>
            <Stack spacing={2} style={{ marginTop: ".8em", marginBottom: ".8em" }}>

                {/* Filters Section */}
                <Paper sx={{ display: 'flex', flexDirection: 'column', flexWrap: 'wrap' }}>
                    <div className="headingcolor">
                        Filter Releases
                    </div>
                    <div className="flex wrap" style={{ justifyContent: 'center', alignItems: "baseline", paddingBottom: "8px", backgroundColor: "rgba(66, 165, 245, 0.15)" }}>
                        <GenericSearchBar filterValue={filterText} setFilter={setFilterText} />
                        <FormControl variant="standard" sx={{ m: 1, minWidth: 140, maxWidth: 300 }}>
                            <InputLabel id="type-filter-label">Type</InputLabel>
                            <Select
                                labelId="type-filter-label"
                                value={selectedTypes}
                                label="Type"
                                onChange={handleTypeChange}
                                multiple
                                renderValue={(selected) => selected.join(', ')}
                                MenuProps={MenuProps}
                            >
                                {typesList.map((typeEntry) => (
                                    <MenuItem key={typeEntry} value={typeEntry}>
                                        <Checkbox checked={selectedTypes.indexOf(typeEntry) > -1} />
                                        <ListItemText primary={typeEntry} />
                                    </MenuItem>
                                ))}
                            </Select>
                        </FormControl>
                        <FormControl variant="standard" sx={{ m: 1, minWidth: 180, maxWidth: 300 }}>
                            <InputLabel id="org-filter-label">Organization</InputLabel>
                            <Select
                                labelId="org-filter-label"
                                value={selectedOrg}
                                label="Org.."
                                onChange={handleOrgChange}
                                multiple
                                renderValue={(selected) => selected.join(', ')}
                                MenuProps={MenuProps}
                            >
                                {orgsList.map((orgEntry) => (
                                    <MenuItem key={orgEntry} value={orgEntry}>
                                        <Checkbox checked={selectedOrg.indexOf(orgEntry) > -1} />
                                        <ListItemText primary={orgEntry} />
                                    </MenuItem>
                                ))}
                            </Select>
                        </FormControl>
                        <IconButton title="Reset Filters" onClick={resetFilters} sx={{
                            color: 'white',
                            backgroundColor: '#1976d27d',
                            '&:hover': {
                                backgroundColor: "#1976d2",
                            }
                        }} aria-label="delete">
                            <RotateLeftIcon />
                        </IconButton>
                    </div>
                </Paper>

                {/* Release Timeline Section */}
                {isLoading && <CircularProgress />}
                {!isLoading &&
                    <Timeline
                        sx={{
                            [`& .${timelineOppositeContentClasses.root}`]: {
                                flex: 0.15,
                            },
                        }}
                    >
                        {filteredModuleList.map((mod, index) => { return <ReleaseCard module={mod} key={index} /> })}
                    </Timeline>
                }
            </Stack>
        </>
    );
}

export default Releases;
