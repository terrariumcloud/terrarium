import { useEffect, useState } from 'react';
import dummyReleaseTypes from '../assets/release-types-list.json'

const dummyReleaseTypeResponseData = {
    "releaseTypes": dummyReleaseTypes,
}

export interface ReleaseTypeResponse {
    releaseTypes: string[]
}

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
    const releaseTypeListURI = "/api/release/types"
    const [releaseTypes, setReleaseTypes] = useState<string[]>([])
    useEffect(() => {
        fetch(releaseTypeListURI)
            .then((response) => {
                return dummyReleaseTypeResponseData
            })
            .then((response: ReleaseTypeResponse) => {
                setReleaseTypes(response.releaseTypes);
            })
    }, [])
    return releaseTypes
}
