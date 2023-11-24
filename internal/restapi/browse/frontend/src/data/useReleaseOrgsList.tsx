import { useEffect, useState } from 'react';

// const dummyReleaseOrgResponseData = {
//     "releaseOrgs": [
//         "saas",
//         "spvss-ivp",
//         "cie"
//     ],
// }

export const useReleaseOrgList = (): string[] => {
    const releaseOrgListURI = "/api/organizations"
    const [releaseOrgs, setReleaseOrgs] = useState<string[]>([])
    useEffect(() => {
        fetch(releaseOrgListURI)
            .then((response) => {
                return response.json();
                // return dummyReleaseOrgResponseData
            })
            .then((response: string[]) => {
                setReleaseOrgs(response);
            })
    }, [])
    return releaseOrgs
}
