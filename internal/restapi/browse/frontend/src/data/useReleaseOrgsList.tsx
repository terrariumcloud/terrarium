import { useEffect, useState } from 'react';
import dummyReleaseOrgs from '../assets/release-orgs-list.json'

const dummyReleaseOrgResponseData = {
    "releaseOrgs": dummyReleaseOrgs,
}

export interface ReleaseOrgResponse {
    releaseOrgs: string[]
}

export const useReleaseOrgList = (): string[] => {
    const releaseOrgListURI = "/api/release/organizations"
    const [releaseOrgs, setReleaseOrgs] = useState<string[]>([])
    // const [releaseOrgsCount, setReleaseOrgsCount] = useState<number>(0)
    useEffect(() => {
        fetch(releaseOrgListURI)
            .then((response) => {
                // return response.json();
                return dummyReleaseOrgResponseData
            })
            .then((response: ReleaseOrgResponse) => {
                setReleaseOrgs(response.releaseOrgs);
            })
    }, [])
    return releaseOrgs
}
