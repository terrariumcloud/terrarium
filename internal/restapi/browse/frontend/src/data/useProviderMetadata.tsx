import {useEffect, useState} from 'react';

export interface ProviderMetadata {
    organization: string
    name: string
    description?: string
    source_repo_url: string
    maturity?: string
    versions: string[]
}

interface ProviderMetadataResponse {
    data: ProviderMetadata
}

export const useProviderMetadata = (organization: string | undefined, name: string | undefined) => {
    const [providerMetadata, setProviderMetadata] = useState<ProviderMetadata|null>(null)
    const providerMetadataURI = `/api/providers/${organization}/${name}`
    useEffect(() => {
        fetch(providerMetadataURI)
            .then((response) => {
                return response.json();
            })
            .then((response: ProviderMetadataResponse) => {
                setProviderMetadata(response.data);
            })
    }, [providerMetadataURI])
    return providerMetadata
}