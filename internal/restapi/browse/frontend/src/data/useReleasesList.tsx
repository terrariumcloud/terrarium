import { useEffect, useState } from 'react';

export interface ReleaseResponseData {
    releases: ReleaseEntry[]
    count: number
}

export interface ReleaseResponse {
    releases: ReleaseEntry[]
}

export interface ReleaseLinks { title?: string; url?: string; }

export interface ReleaseEntry {
    name: string
    createdAt: string
    description?: string
    organization: string
    type: string
    version: string
    links?: { title?: string; url?: string; }[]
}

export const useReleaseList = (selectedTime: string, setIsLoading: ((value: (((prevState: boolean) => boolean) | boolean)) => void)): ReleaseEntry[] => {
    const ageQuery = "age=" + selectedTime
    const releaseListURI = "/api/releases?" + ageQuery
    const [releases, setReleases] = useState<ReleaseEntry[]>([])
    useEffect(() => {
        setIsLoading(true)
        fetch(releaseListURI)
            .then((response) => {
                return response.json();
            })
            .then((response: ReleaseResponse) => {
                if (response.releases == null) response.releases = []
                setReleases(response.releases.reverse());
                setIsLoading(false)
            }).catch(() => {
                setIsLoading(false)
            })
    }, [releaseListURI, setIsLoading])
    return releases
}

export const useFilteredReleaseList = (selectedTypes: string[], selectedOrgs: string[], selectedTime: string): [ReleaseEntry[], string, ((value: (((prevState: string) => string) | string)) => void), boolean] => {

    const [filterText, setFilterText] = useState<string>("");
    const [isLoading, setIsLoading] = useState(false);

    const releases = useReleaseList(selectedTime, setIsLoading)

    const releasesFilteredonTypes = releases
        .filter((releaseInfo) => {
            if (!selectedTypes.length) return true
            return selectedTypes.includes(releaseInfo.type)
        })
    const releasesFilteredonTypesOrgs = releasesFilteredonTypes
        .filter((releaseInfo) => {
            if (!selectedOrgs.length) return true
            return selectedOrgs.includes(releaseInfo.organization)
        })
    const filteredReleases = releasesFilteredonTypesOrgs
        .filter((releaseInfo) => {
            const filterValue = filterText.toLowerCase()

            if (filterText === "") {
                return true
            }

            const releaseSearchText = releaseInfo.organization + " "
                + releaseInfo.name + " " + releaseInfo.type + " " + releaseInfo.version

            return releaseSearchText.toLowerCase().includes(filterValue)
        })
    return [filteredReleases, filterText, setFilterText, isLoading]
}
