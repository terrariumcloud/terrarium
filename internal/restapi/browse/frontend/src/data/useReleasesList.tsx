import {useEffect, useState} from 'react';

export interface ReleaseResponseData {
    releases: ReleaseEntry[]
    count: number
}

export interface ReleaseResponse {
    releases: ReleaseEntry[]
}

export interface ReleaseEntry {
    name: string
    createdAt: string
    description?: string
    Organization: string
    type: string
    version: string
}

export const useReleaseList = ():ReleaseEntry[] => {
    const releaseListURI = "/api/releases"
    const [releases, setReleases] = useState<ReleaseEntry[]>([])
// const [releasesCount, setReleasesCount] = useState<number>(0)
    useEffect(() => {
        fetch(releaseListURI)
            .then((response) => {
                return response.json();
            })
            .then((response: ReleaseResponse) => {
                setReleases(response.releases);
            })
    }, [])
    return releases
}

export const useFilteredReleaseList = ():[ReleaseEntry[], string, ((value: (((prevState: string) => string) | string)) => void)] =>  {
    const releases = useReleaseList()
    const [filterText, setFilterText] = useState<string>("")
    const filteredReleases = releases
        .filter((releaseInfo) => {
            const filterValue = filterText.toLowerCase()

            if (filterText === "") {
                return true
            }
            const releaseSearchText = releaseInfo.Organization + " "
                + releaseInfo.name
            console.log("Release Information", releaseInfo)
            return releaseSearchText.toLowerCase().includes(filterValue)
        })
    return [filteredReleases, filterText, setFilterText]
}
