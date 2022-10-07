import {useEffect, useState} from 'react';

export interface ModuleMetadata {
    organization: string
    name: string
    provider: string
    description?: string
    source_url: string
    maturity?: string
    versions: string[]
}

interface ModuleMetadataResponse {
    data: ModuleMetadata
}

// const fakeData: ModuleMetadata = {
//     description: "This is the description for the module it is supposedly a long text",
//     name: "test-module",
//     organization: "cie",
//     provider: "aws",
//     source_url: "https://github.com/...",
//     versions: [
//         "1.0.0",
//         "0.9.17",
//         "0.9.16",
//         "0.9.15",
//         "0.9.14",
//         "0.9.13",
//         "0.9.12",
//     ]
// }

export const useModuleMetadata = (organization: string | undefined, name: string | undefined, provider: string | undefined) => {
    const [moduleMetadata, setModuleMetadata] = useState<ModuleMetadata|null>(null)
    const moduleMetadataURI = `/api/modules/${organization}/${name}/${provider}`
    useEffect(() => {
        fetch(moduleMetadataURI)
            .then((response) => {
                return response.json();
            })
            .then((response: ModuleMetadataResponse) => {
                setModuleMetadata(response.data);
            })
    }, [moduleMetadataURI])
    return moduleMetadata
}