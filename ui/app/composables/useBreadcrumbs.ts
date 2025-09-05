export const useBreadcrumbs = () => {
  const route = useRoute()
  
  const breadcrumbItems = computed(() => {
    const pathSegments = route.path.split('/').filter(Boolean)
    const items: Array<{ label: string; link: string }> = [
      { label: 'Dashboard', link: '/' }
    ]

    // Build breadcrumbs based on route segments
    let currentPath = ''
    
    for (let i = 0; i < pathSegments.length; i++) {
      const segment = pathSegments[i]
      currentPath += `/${segment}`
      
      let label = segment
      let link = currentPath
      
      // Handle dynamic segments and custom labels
      if (segment === 'applications') {
        label = 'Applications'
      } else if (segment === 'keys') {
        label = 'API Keys'
      } else if (segment === 'users') {
        label = 'Users'
      } else if (segment === 'models') {
        label = 'Models'
      } else if (segment === 'environments') {
        label = 'Environments'
      } else if (segment === 'settings') {
        label = 'Settings'
      } else if (segment.startsWith('app_')) {
        // For application IDs, we'll show the application name
        // In a real app, this would fetch from an API or store
        label = 'Application Details'
      } else if (segment.startsWith('key_')) {
        // For key IDs, we'll show the key name
        label = 'API Key Details'
      } else if (segment.startsWith('user_')) {
        // For user IDs
        label = 'User Details'
      } else {
        // Capitalize first letter for other segments
        label = segment.charAt(0).toUpperCase() + segment.slice(1)
      }
      
      // Don't add link to current page
      if (i === pathSegments.length - 1) {
        link = ''
      }
      
      items.push({ label, link })
    }
    
    return items
  })
  
  return {
    breadcrumbItems
  }
}