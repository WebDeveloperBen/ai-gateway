<script setup lang="ts">
import { Plus, Search, Copy, Eye, EyeOff, MoreVertical, Key, Calendar, Activity, X } from "lucide-vue-next"

// Sample API Keys data
const apiKeys = ref([
  {
    id: "key_1",
    name: "Production Key",
    key: "sk-proj-abc123...def456",
    applicationId: "app_1",
    applicationName: "Customer Service Bot",
    created: "2024-12-15T10:00:00Z",
    lastUsed: "2025-01-15T14:30:00Z",
    status: "active",
    permissions: ["read", "write"],
    requestCount: 25400,
    rateLimit: "1000/hour"
  },
  {
    id: "key_2",
    name: "Development Key",
    key: "sk-proj-xyz789...uvw012",
    applicationId: "app_1",
    applicationName: "Customer Service Bot",
    created: "2024-11-20T09:15:00Z",
    lastUsed: "2025-01-14T11:45:00Z",
    status: "active",
    permissions: ["read"],
    requestCount: 8900,
    rateLimit: "100/hour"
  },
  {
    id: "key_3",
    name: "Testing Key",
    key: "sk-proj-mno345...rst678",
    applicationId: "app_2",
    applicationName: "Content Generator",
    created: "2024-10-10T16:20:00Z",
    lastUsed: "2024-12-30T08:00:00Z",
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
const visibleKeys = ref<Set<string>>(new Set())

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
  // Navigate to key details or perform action
  console.log("Selected key:", key)
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
    const matchesSearch =
      key.name.toLowerCase().includes(searchQuery.value.toLowerCase()) ||
      key.applicationName.toLowerCase().includes(searchQuery.value.toLowerCase())
    const matchesApp = selectedApp.value === "all" || key.applicationId === selectedApp.value
    const matchesStatus = selectedStatus.value === "all" || key.status === selectedStatus.value
    return matchesSearch && matchesApp && matchesStatus
  })
})

const toggleKeyVisibility = (keyId: string) => {
  if (visibleKeys.value.has(keyId)) {
    visibleKeys.value.delete(keyId)
  } else {
    visibleKeys.value.add(keyId)
  }
}

const copyToClipboard = async (text: string) => {
  try {
    await navigator.clipboard.writeText(text)
    // You could add a toast notification here
  } catch (err) {
    console.error("Failed to copy text: ", err)
  }
}

const formatNumber = (num: number) => {
  return new Intl.NumberFormat().format(num)
}

const getStatusBadgeClass = (status: string) => {
  return status === "active"
    ? "bg-green-100 text-green-800 dark:bg-green-900 dark:text-green-200"
    : "bg-gray-100 text-gray-800 dark:bg-gray-900 dark:text-gray-200"
}

const maskKey = (key: string) => {
  return key.substring(0, 12) + "..." + key.substring(key.length - 6)
}
</script>

<template>
  <div class="flex flex-col gap-6">
    <!-- Header -->
    <div class="flex items-center justify-between">
      <div>
        <h1 class="text-3xl font-bold tracking-tight">API Keys</h1>
        <p class="text-muted-foreground">Manage API keys for your applications</p>
      </div>
      <UiButton @click="showCreateModal = true" class="gap-2">
        <Plus class="h-4 w-4" />
        Create API Key
      </UiButton>
    </div>

    <!-- Stats Cards -->
    <div class="grid grid-cols-1 sm:grid-cols-2 lg:grid-cols-4 gap-4">
      <UiCard>
        <UiCardHeader class="flex flex-row items-center justify-between space-y-0 pb-2">
          <UiCardTitle class="text-sm font-medium">Total Keys</UiCardTitle>
          <Key class="h-4 w-4 text-muted-foreground" />
        </UiCardHeader>
        <UiCardContent>
          <div class="text-2xl font-bold">{{ apiKeys.length }}</div>
        </UiCardContent>
      </UiCard>

      <UiCard>
        <UiCardHeader class="flex flex-row items-center justify-between space-y-0 pb-2">
          <UiCardTitle class="text-sm font-medium">Active Keys</UiCardTitle>
          <Activity class="h-4 w-4 text-muted-foreground" />
        </UiCardHeader>
        <UiCardContent>
          <div class="text-2xl font-bold">{{ apiKeys.filter((key) => key.status === "active").length }}</div>
        </UiCardContent>
      </UiCard>

      <UiCard>
        <UiCardHeader class="flex flex-row items-center justify-between space-y-0 pb-2">
          <UiCardTitle class="text-sm font-medium">Total Requests</UiCardTitle>
          <Activity class="h-4 w-4 text-muted-foreground" />
        </UiCardHeader>
        <UiCardContent>
          <div class="text-2xl font-bold">
            {{ formatNumber(apiKeys.reduce((sum, key) => sum + key.requestCount, 0)) }}
          </div>
        </UiCardContent>
      </UiCard>

      <UiCard>
        <UiCardHeader class="flex flex-row items-center justify-between space-y-0 pb-2">
          <UiCardTitle class="text-sm font-medium">Applications</UiCardTitle>
          <Activity class="h-4 w-4 text-muted-foreground" />
        </UiCardHeader>
        <UiCardContent>
          <div class="text-2xl font-bold">{{ applications.length }}</div>
        </UiCardContent>
      </UiCard>
    </div>

    <!-- Search & Filter Command -->
    <UiCommand class="rounded-lg border shadow-sm">
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
              icon="lucide:layers"
              @select="filterByApp(app.id)"
            />
          </UiCommandGroup>
          <UiCommandSeparator />

          <UiCommandGroup heading="Filter by Status">
            <UiCommandItem
              value="status:active"
              text="Active"
              icon="lucide:check-circle"
              @select="filterByStatus('active')"
            />
            <UiCommandItem
              value="status:inactive"
              text="Inactive"
              icon="lucide:x-circle"
              @select="filterByStatus('inactive')"
            />
            <UiCommandItem value="status:all" text="All Status" icon="lucide:circle" @select="filterByStatus('all')" />
          </UiCommandGroup>
        </template>

        <template v-else>
          <UiCommandGroup heading="API Keys">
            <UiCommandItem
              v-for="key in filteredKeys.slice(0, 8)"
              :key="key.id"
              :value="key.name"
              :text="key.name"
              icon="lucide:key"
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

    <!-- API Keys Table -->
    <UiCard>
      <UiCardHeader>
        <UiCardTitle>API Keys</UiCardTitle>
      </UiCardHeader>
      <UiCardContent>
        <div class="flex flex-col gap-4">
          <div v-if="filteredKeys.length === 0" class="text-center py-8">
            <Key class="mx-auto h-12 w-12 text-muted-foreground" />
            <h3 class="mt-4 text-lg font-medium">No API keys found</h3>
            <p class="text-muted-foreground">Try adjusting your search or create a new API key.</p>
          </div>

          <div v-for="key in filteredKeys" :key="key.id" class="border rounded-lg p-4 space-y-3">
            <div class="flex items-start justify-between">
              <div class="space-y-1">
                <div class="flex items-center gap-3">
                  <h4 class="font-medium">{{ key.name }}</h4>
                  <UiBadge :class="getStatusBadgeClass(key.status)">
                    {{ key.status }}
                  </UiBadge>
                </div>
                <p class="text-sm text-muted-foreground">{{ key.applicationName }}</p>
              </div>
              <UiDropdownMenu>
                <UiDropdownMenuTrigger as-child>
                  <UiButton variant="ghost" size="sm">
                    <MoreVertical class="h-4 w-4" />
                  </UiButton>
                </UiDropdownMenuTrigger>
                <UiDropdownMenuContent align="end">
                  <UiDropdownMenuItem> Edit Key </UiDropdownMenuItem>
                  <UiDropdownMenuItem> View Usage </UiDropdownMenuItem>
                  <UiDropdownMenuItem> Regenerate </UiDropdownMenuItem>
                  <UiDropdownMenuSeparator />
                  <UiDropdownMenuItem class="text-red-600"> Delete Key </UiDropdownMenuItem>
                </UiDropdownMenuContent>
              </UiDropdownMenu>
            </div>

            <div class="flex items-center gap-2 bg-muted/50 rounded p-2">
              <code class="text-sm font-mono flex-1">
                {{ visibleKeys.has(key.id) ? key.key : maskKey(key.key) }}
              </code>
              <UiButton variant="ghost" size="sm" @click="toggleKeyVisibility(key.id)">
                <Eye v-if="!visibleKeys.has(key.id)" class="h-4 w-4" />
                <EyeOff v-else class="h-4 w-4" />
              </UiButton>
              <UiButton variant="ghost" size="sm" @click="copyToClipboard(key.key)">
                <Copy class="h-4 w-4" />
              </UiButton>
            </div>

            <div class="grid grid-cols-2 md:grid-cols-5 gap-4 text-sm">
              <div>
                <p class="text-muted-foreground">Created</p>
                <p class="font-medium">{{ new Date(key.created).toLocaleDateString() }}</p>
              </div>
              <div>
                <p class="text-muted-foreground">Last Used</p>
                <p class="font-medium">{{ new Date(key.lastUsed).toLocaleDateString() }}</p>
              </div>
              <div>
                <p class="text-muted-foreground">Requests</p>
                <p class="font-medium">{{ formatNumber(key.requestCount) }}</p>
              </div>
              <div>
                <p class="text-muted-foreground">Rate Limit</p>
                <p class="font-medium">{{ key.rateLimit }}</p>
              </div>
              <div>
                <p class="text-muted-foreground">Permissions</p>
                <p class="font-medium">{{ key.permissions.join(", ") }}</p>
              </div>
            </div>
          </div>
        </div>
      </UiCardContent>
    </UiCard>
  </div>
</template>
