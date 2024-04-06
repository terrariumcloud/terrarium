import {useEffect, useState} from 'react';

export interface ProviderResponseData {
    providers: ProviderEntry[]
    count: number
}

export interface ProviderResponse {
    providers: ProviderEntry[]
}

export interface ProviderEntry {
    organization: string
    name: string
    description?: string
    source_repo_url: string
    maturity?: string
}

export const useProviderList = ():ProviderEntry[] => {
    const providerListURI = "/api/providers"
    const [providers, setProviders] = useState<ProviderEntry[]>([])
    useEffect(() => {
        fetch(providerListURI)
            .then((response) => {
                return response.json();
            })
            .then((response: ProviderResponse) => {
                setProviders(response.providers);
            })
    }, [])
    return providers
}

export const useFilteredProviderList = ():[ProviderEntry[], string, ((value: (((prevState: string) => string) | string)) => void)] =>  {
    const providers = useProviderList()
    const [filterText, setFilterText] = useState<string>("")
    const filteredProviders = providers
        .filter((providerInfo) => {
            const filterValue = filterText.toLowerCase()

            if (filterText === "") {
                return true
            }
            const providerSearchText = providerInfo.organization + " "
                + providerInfo.name

            return providerSearchText.toLowerCase().includes(filterValue)
        })
    return [filteredProviders, filterText, setFilterText]
}
