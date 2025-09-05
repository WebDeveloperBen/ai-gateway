<script setup lang="ts">
import { Plus, Key, X, CheckCircle, XCircle, Circle, Layers, Activity } from "lucide-vue-next"

// Get the app ID from route params
const route = useRoute()
const appId = route.params.appId as string

// Sample API Keys data
const apiKeys = ref([
  {
    id: "key_1",
    name: "Production Key",
    keyPrefix: "sk-",
    applicationId: "app_1",
    applicationName: "Customer Service Bot",
    created: "2024-12-15T10:00:00Z",
    lastUsed: "2025-01-15T14:30:00Z",
    expires: "2025-12-15T10:00:00Z",
    status: "active",
    permissions: ["read", "write"],
    requestCount: 25400,
    rateLimit: "1000/hour"
  },
  {
    id: "key_2",
    name: "Development Key",
    keyPrefix: "sk-",
    applicationId: "app_1",
    applicationName: "Customer Service Bot",
    created: "2024-11-20T09:15:00Z",
    lastUsed: "2025-01-14T11:45:00Z",
    expires: "2025-11-20T09:15:00Z",
    status: "active",
    permissions: ["read"],
    requestCount: 8900,
    rateLimit: "100/hour"
  },
  {
    id: "key_3",
    name: "Testing Key",
    keyPrefix: "sk-",
    applicationId: "app_2",
    applicationName: "Content Generator",
    created: "2024-10-10T16:20:00Z",
    lastUsed: "2024-12-30T08:00:00Z",
    expires: "2025-10-10T16:20:00Z",
    status: "inactive",
    permissions: ["read", "write"],
    requestCount: 450,
    rateLimit: "50/hour"
  }
])

const searchQuery = ref("")
const selectedApp = ref("all")
const selectedStatus = ref("all")
const showCreateModal = ref(false)
const showFilters = ref(false)

// Command interface methods
const filterByApp = (appId: string) => {
  selectedApp.value = appId
  searchQuery.value = ""
  showFilters.value = false
}

const filterByStatus = (status: string) => {
  selectedStatus.value = status
  searchQuery.value = ""
  showFilters.value = false
}

const selectKey = (key: any) => {
  // Navigate to key details
  navigateTo(`/applications/${appId}/keys/${key.id}`)
  searchQuery.value = ""
}

const clearAllFilters = () => {
  selectedApp.value = "all"
  selectedStatus.value = "all"
  searchQuery.value = ""
}

const hasActiveFilters = computed(() => {
  return selectedApp.value !== "all" || selectedStatus.value !== "all"
})

const applications = computed(() => {
  const uniqueApps = [...new Set(apiKeys.value.map((key) => ({ id: key.applicationId, name: key.applicationName })))]
  return uniqueApps
})

const filteredKeys = computed(() => {
  return apiKeys.value.filter((key) => {
    // Always filter by the current appId for app-specific pages
    const matchesApp = key.applicationId === appId
    const matchesSearch =
      key.name.toLowerCase().includes(searchQuery.value.toLowerCase()) ||
      key.applicationName.toLowerCase().includes(searchQuery.value.toLowerCase())
    const matchesStatus = selectedStatus.value === "all" || key.status === selectedStatus.value
    return matchesApp && matchesSearch && matchesStatus
  })
})

// Stats card configuration - easily replaceable with API call
const statsCards = computed(() => {
  // This computed property can easily be replaced with:
  // const { data: statsCards } = await useFetch(`/api/applications/${appId}/stats`)

  return [
    {
      title: "Total Keys",
      value: filteredKeys.value.length,
      icon: Key,
      description: "All API keys for this application"
    },
    {
      title: "Active Keys",
      value: filteredKeys.value.filter((key) => key.status === "active").length,
      icon: CheckCircle,
      description: "Currently active and usable"
    },
    {
      title: "Total Requests",
      value: filteredKeys.value.reduce((sum, key) => sum + key.requestCount, 0),
      icon: Activity,
      description: "API calls made this month"
    },
    {
      title: "Applications",
      value: 1,
      icon: Layers,
      description: "Connected applications"
    }
  ]
})

const navigateToKey = async (key: any) => {
  const keyId = typeof key === "string" ? key : key.id
  console.log("Navigating to key:", keyId)
  try {
    await navigateTo(`/applications/${appId}/keys/${keyId}`, { replace: false })
  } catch (error) {
    console.error("Navigation error:", error)
  }
}

const handleRegenerateKey = (key: any) => {
  console.log("Regenerating key:", key.id)
  // TODO: Implement key regeneration
}

const handleDeleteKey = (key: any) => {
  console.log("Deleting key:", key.id)
  // TODO: Implement key deletion
}

const onApiKeyCreated = (apiKeyData: any) => {
  console.log("API key created:", apiKeyData)

  // Add missing fields that the keys page expects
  const enrichedApiKey = {
    ...apiKeyData,
    applicationName:
      applications.value.find((app) => app.id === apiKeyData.applicationId)?.name || "Unknown Application",
    lastUsed: apiKeyData.created, // Use creation date as last used for new keys
    permissions: apiKeyData.permissions || ["read"], // Default permissions
    requestCount: 0, // New keys start with 0 requests
    rateLimit: apiKeyData.rateLimit || "1000/hour" // Default rate limit
  }

  // Add the new key to the existing list
  apiKeys.value.unshift(enrichedApiKey)
}
</script>

<template>
  <div class="flex flex-col gap-6">
    <!-- Header -->
    <div class="flex items-center justify-between">
      <div>
        <h1 class="text-3xl font-bold tracking-tight text-primary">
          {{ apiKeys.find((k) => k.applicationId === appId)?.applicationName || "Application" }}
        </h1>
        <p class="text-muted-foreground">Manage API keys for this specific application</p>
      </div>
      <UiButton @click="showCreateModal = true" class="gap-2">
        <Plus class="h-4 w-4" />
        Create API Key
      </UiButton>
    </div>

    <!-- Stats Cards -->
    <div class="grid grid-cols-1 sm:grid-cols-2 lg:grid-cols-4 gap-6">
      <CardsStats
        v-for="card in statsCards"
        :key="card.title"
        :title="card.title"
        :value="card.value"
        :icon="card.icon"
        :description="card.description"
      />
    </div>

    <!-- Search & Filter Command -->
    <UiCommand class="rounded-lg border">
      <UiCommandInput v-model="searchQuery" placeholder="Search API keys, filter by application or status..." />
      <UiCommandList v-if="searchQuery || showFilters">
        <UiCommandEmpty>No API keys found.</UiCommandEmpty>

        <template v-if="!searchQuery">
          <UiCommandGroup heading="Filter by Application">
            <UiCommandItem
              v-for="app in applications"
              :key="app.id"
              :value="`app:${app.name}`"
              :text="app.name"
              :icon="Layers"
              @select="filterByApp(app.id)"
            />
          </UiCommandGroup>
          <UiCommandSeparator />

          <UiCommandGroup heading="Filter by Status">
            <UiCommandItem value="status:active" text="Active" :icon="CheckCircle" @select="filterByStatus('active')" />
            <UiCommandItem
              value="status:inactive"
              text="Inactive"
              :icon="XCircle"
              @select="filterByStatus('inactive')"
            />
            <UiCommandItem value="status:all" text="All Status" :icon="Circle" @select="filterByStatus('all')" />
          </UiCommandGroup>
        </template>

        <template v-else>
          <UiCommandGroup heading="API Keys">
            <UiCommandItem
              v-for="key in filteredKeys.slice(0, 8)"
              :key="key.id"
              :value="key.name"
              :text="key.name"
              :icon="Key"
              @select="selectKey(key)"
            />
          </UiCommandGroup>
        </template>
      </UiCommandList>
    </UiCommand>

    <!-- Active Filters -->
    <div v-if="hasActiveFilters" class="flex items-center gap-2">
      <span class="text-sm text-muted-foreground">Active filters:</span>
      <UiBadge v-if="selectedApp !== 'all'" variant="secondary" class="gap-1">
        {{ applications.find((app) => app.id === selectedApp)?.name || "Unknown App" }}
        <button @click="selectedApp = 'all'" class="ml-1 hover:bg-muted-foreground/20 rounded-sm">
          <X class="h-3 w-3" />
        </button>
      </UiBadge>
      <UiBadge v-if="selectedStatus !== 'all'" variant="secondary" class="gap-1">
        {{ selectedStatus }}
        <button @click="selectedStatus = 'all'" class="ml-1 hover:bg-muted-foreground/20 rounded-sm">
          <X class="h-3 w-3" />
        </button>
      </UiBadge>
      <UiButton v-if="hasActiveFilters" variant="ghost" size="sm" @click="clearAllFilters" class="h-6 px-2 text-xs">
        Clear all
      </UiButton>
    </div>

    <!-- API Keys List -->
    <ApiKeysList
      :keys="filteredKeys"
      @select-api-key="navigateToKey"
      @regenerate-key="handleRegenerateKey"
      @delete-key="handleDeleteKey"
    />

    <!-- Create API Key Modal -->
    <ModalsCreateApiKey v-model:open="showCreateModal" @created="onApiKeyCreated" />
  </div>
</template>
