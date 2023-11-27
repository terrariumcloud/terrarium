import { useEffect, useState } from 'react';

export const useReleaseOrgList = (): string[] => {
    const releaseOrgListURI = "/api/organizations"
    const [releaseOrgs, setReleaseOrgs] = useState<string[]>([])
    useEffect(() => {
        fetch(releaseOrgListURI)
            .then((response) => {
                return response.json();
            })
            .then((response: string[]) => {
                setReleaseOrgs(response);
            })
    }, [])
    return releaseOrgs
}
