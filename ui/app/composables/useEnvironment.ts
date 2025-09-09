interface Environment {
  name: string
  logo: any
  description: string
  status: 'active' | 'maintenance'
}

// Global state for the selected environment
const selectedEnvironment = ref<Environment | null>(null)

export const useEnvironment = () => {
  const setEnvironment = (environment: Environment) => {
    selectedEnvironment.value = environment
  }

  const getEnvironment = () => {
    return selectedEnvironment.value
  }

  const isEnvironmentSelected = () => {
    return selectedEnvironment.value !== null
  }

  // Helper to get environment-scoped data
  const getEnvironmentData = <T>(data: T[], filterFn?: (item: T, env: Environment) => boolean) => {
    if (!selectedEnvironment.value) return data
    
    if (filterFn) {
      return data.filter(item => filterFn(item, selectedEnvironment.value!))
    }
    
    return data
  }

  return {
    selectedEnvironment: readonly(selectedEnvironment),
    setEnvironment,
    getEnvironment,
    isEnvironmentSelected,
    getEnvironmentData
  }
}