import { useEffect, useState } from 'react';

export interface ReleaseTypeEntry {
    name: string
    createdAt: string
    description?: string
    Organization: string
    type: string
    version: string
    links?: { Title?: string; Url: string; }[]
}

export const useReleaseTypeList = (): string[] => {
    const releaseTypeListURI = "/api/types"
    const [releaseTypes, setReleaseTypes] = useState<string[]>([])
    useEffect(() => {
        fetch(releaseTypeListURI)
            .then((response) => {
                return response.json();
                // return dummyReleaseTypeResponseData
            })
            .then((response: string[]) => {
                setReleaseTypes(response);
            })
    }, [])
    return releaseTypes
}
