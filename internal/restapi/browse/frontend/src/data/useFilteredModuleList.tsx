import {useEffect, useState} from 'react';

export interface ModuleResponseData {
    modules: ModuleEntry[]
    count: number
}

export interface ModuleResponse {
    modules: ModuleEntry[]
}

export interface ModuleEntry {
    organization: string
    name: string
    provider: string
    description?: string
    source_url: string
    maturity?: string
}

export const useModuleList = ():ModuleEntry[] => {
    const moduleListURI = "/api/modules"
    const [modules, setModules] = useState<ModuleEntry[]>([])
// const [modulesCount, setModulesCount] = useState<number>(0)
    useEffect(() => {
        fetch(moduleListURI)
            .then((response) => {
                return response.json();
            })
            .then((response: ModuleResponse) => {
                setModules(response.modules);
            })
    }, [])
    return modules
}

export const useFilteredModuleList = ():[ModuleEntry[], string, ((value: (((prevState: string) => string) | string)) => void)] =>  {
    const modules = useModuleList()
    const [filterText, setFilterText] = useState<string>("")
    const filteredModules = modules
        .filter((moduleInfo) => {
            const filterValue = filterText.toLowerCase()

            if (filterText === "") {
                return true
            }
            const moduleSearchText = moduleInfo.organization + " "
                + moduleInfo.name + " "
                + moduleInfo.provider

            return moduleSearchText.toLowerCase().includes(filterValue)
        })
    return [filteredModules, filterText, setFilterText]
}
