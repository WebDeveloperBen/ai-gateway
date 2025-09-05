<script setup lang="ts">
import {
  Plus,
  MoreVertical,
  Key,
  Activity,
  AlertCircle,
  X,
  CheckCircle,
  XCircle,
  Circle,
  Layers
} from "lucide-vue-next"

// Sample data - replace with actual API calls
const applications = ref([
  {
    id: "app_1",
    name: "Customer Service Bot",
    description: "AI-powered customer support assistant for handling common inquiries",
    status: "active",
    apiKeyCount: 3,
    monthlyRequests: 45200,
    lastUsed: "2025-01-15T10:30:00Z",
    models: ["gpt-4", "gpt-3.5-turbo"],
    team: "Customer Success"
  },
  {
    id: "app_2",
    name: "Content Generator",
    description: "Automated content creation for marketing campaigns",
    status: "active",
    apiKeyCount: 2,
    monthlyRequests: 12800,
    lastUsed: "2025-01-14T16:45:00Z",
    models: ["gpt-4"],
    team: "Marketing"
  },
  {
    id: "app_3",
    name: "Code Assistant",
    description: "Development helper for code review and generation",
    status: "inactive",
    apiKeyCount: 1,
    monthlyRequests: 0,
    lastUsed: "2024-12-20T09:15:00Z",
    models: ["gpt-4"],
    team: "Engineering"
  }
])

const searchQuery = ref("")
const selectedStatus = ref("all")
const showFilters = ref(false)
const showCreateModal = ref(false)

// Get app config
const appConfig = useAppConfig()

// Command interface methods
const filterByStatus = (status: string) => {
  selectedStatus.value = status
  searchQuery.value = ""
  showFilters.value = false
}

const selectApplication = (app: any) => {
  // Navigate to application details
  navigateTo(`/applications/${app.id}`)
  searchQuery.value = ""
}

const clearAllFilters = () => {
  selectedStatus.value = "all"
  searchQuery.value = ""
}

const onApplicationCreated = (applicationId: string) => {
  console.log("Application created:", applicationId)
  // Here you could refresh the applications list or add the new app to the existing list
  // For now, we'll just close the modal (handled by the modal component)
}

const hasActiveFilters = computed(() => {
  return selectedStatus.value !== "all"
})

const filteredApplications = computed(() => {
  return applications.value.filter((app) => {
    const matchesSearch =
      app.name.toLowerCase().includes(searchQuery.value.toLowerCase()) ||
      app.description.toLowerCase().includes(searchQuery.value.toLowerCase())
    const matchesStatus = selectedStatus.value === "all" || app.status === selectedStatus.value
    return matchesSearch && matchesStatus
  })
})

const formatNumber = (num: number) => {
  return new Intl.NumberFormat().format(num)
}

const getStatusBadgeClass = (status: string) => {
  return status === "active"
    ? "bg-green-100 text-green-800 dark:bg-green-900 dark:text-green-200"
    : "bg-gray-100 text-gray-800 dark:bg-gray-900 dark:text-gray-200"
}
</script>

<template>
  <div class="flex flex-col gap-6">
    <!-- Header -->
    <div class="flex items-center justify-between">
      <div>
        <h1 class="text-3xl font-bold tracking-tight">Applications</h1>
        <p class="text-muted-foreground">
          Manage your AI-powered applications and their {{ appConfig.app.name }} access
        </p>
      </div>
      <UiButton @click="showCreateModal = true" class="gap-2">
        <Plus class="h-4 w-4" />
        New Application
      </UiButton>
    </div>

    <!-- Stats Cards -->
    <div class="grid grid-cols-1 sm:grid-cols-2 lg:grid-cols-4 gap-4">
      <UiCard>
        <UiCardHeader class="flex flex-row items-center justify-between space-y-0 pb-2">
          <UiCardTitle class="text-sm font-medium">Total Applications</UiCardTitle>
          <Activity class="h-4 w-4 text-muted-foreground" />
        </UiCardHeader>
        <UiCardContent>
          <div class="text-2xl font-bold">{{ applications.length }}</div>
        </UiCardContent>
      </UiCard>

      <UiCard>
        <UiCardHeader class="flex flex-row items-center justify-between space-y-0 pb-2">
          <UiCardTitle class="text-sm font-medium">Active Applications</UiCardTitle>
          <Activity class="h-4 w-4 text-muted-foreground" />
        </UiCardHeader>
        <UiCardContent>
          <div class="text-2xl font-bold">
            {{ applications.filter((app) => app.status === "active").length }}
          </div>
        </UiCardContent>
      </UiCard>

      <UiCard>
        <UiCardHeader class="flex flex-row items-center justify-between space-y-0 pb-2">
          <UiCardTitle class="text-sm font-medium">Total API Keys</UiCardTitle>
          <Key class="h-4 w-4 text-muted-foreground" />
        </UiCardHeader>
        <UiCardContent>
          <div class="text-2xl font-bold">
            {{ applications.reduce((sum, app) => sum + app.apiKeyCount, 0) }}
          </div>
        </UiCardContent>
      </UiCard>

      <UiCard>
        <UiCardHeader class="flex flex-row items-center justify-between space-y-0 pb-2">
          <UiCardTitle class="text-sm font-medium">Monthly Requests</UiCardTitle>
          <Activity class="h-4 w-4 text-muted-foreground" />
        </UiCardHeader>
        <UiCardContent>
          <div class="text-2xl font-bold">
            {{ formatNumber(applications.reduce((sum, app) => sum + app.monthlyRequests, 0)) }}
          </div>
        </UiCardContent>
      </UiCard>
    </div>

    <!-- Search & Filter Command -->
    <UiCommand class="rounded-lg border shadow-sm">
      <UiCommandInput v-model="searchQuery" placeholder="Search applications or filter by status..." />
      <UiCommandList v-if="searchQuery || showFilters">
        <UiCommandEmpty>No applications found.</UiCommandEmpty>

        <template v-if="!searchQuery">
          <UiCommandGroup heading="Filter by Status">
            <UiCommandItem
              value="status:active"
              text="Active Applications"
              :icon="CheckCircle"
              @select="filterByStatus('active')"
            />
            <UiCommandItem
              value="status:inactive"
              text="Inactive Applications"
              :icon="XCircle"
              @select="filterByStatus('inactive')"
            />
            <UiCommandItem value="status:all" text="All Applications" :icon="Circle" @select="filterByStatus('all')" />
          </UiCommandGroup>
          <UiCommandSeparator />

          <UiCommandGroup heading="Quick Actions">
            <UiCommandItem
              value="action:create"
              text="Create New Application"
              :icon="Plus"
              @select="showCreateModal = true"
            />
          </UiCommandGroup>
        </template>

        <template v-else>
          <UiCommandGroup heading="Applications">
            <UiCommandItem
              v-for="app in filteredApplications.slice(0, 8)"
              :key="app.id"
              :value="app.name"
              :text="app.name"
              :icon="Layers"
              @select="selectApplication(app)"
            />
          </UiCommandGroup>
        </template>
      </UiCommandList>
    </UiCommand>

    <!-- Active Filters -->
    <div v-if="hasActiveFilters" class="flex items-center gap-2">
      <span class="text-sm text-muted-foreground">Active filters:</span>
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

    <!-- Applications List -->
    <div class="flex flex-col gap-4">
      <div v-if="filteredApplications.length === 0" class="text-center py-12">
        <AlertCircle class="mx-auto h-12 w-12 text-muted-foreground" />
        <h3 class="mt-4 text-lg font-medium">No applications found</h3>
        <p class="text-muted-foreground">Try adjusting your search or create a new application.</p>
      </div>

      <UiCard 
        v-for="app in filteredApplications" 
        :key="app.id" 
        interactive
        @click="navigateTo(`/applications/${app.id}`)"
      >
        <UiCardHeader>
          <div class="flex items-start justify-between">
            <div class="space-y-1 flex-1">
              <div class="flex items-center gap-3">
                <UiCardTitle class="text-lg">{{ app.name }}</UiCardTitle>
                <UiBadge :class="getStatusBadgeClass(app.status)">
                  {{ app.status }}
                </UiBadge>
              </div>
              <UiCardDescription class="text-sm">
                {{ app.description }}
              </UiCardDescription>
              <div class="flex items-center gap-4 text-xs text-muted-foreground">
                <span>Team: {{ app.team }}</span>
                <span>•</span>
                <span>Models: {{ app.models.join(", ") }}</span>
              </div>
            </div>
            <UiDropdownMenu>
              <UiDropdownMenuTrigger as-child>
                <UiButton variant="ghost" size="sm" @click.stop>
                  <MoreVertical class="h-4 w-4" />
                </UiButton>
              </UiDropdownMenuTrigger>
              <UiDropdownMenuContent align="end">
                <UiDropdownMenuItem as-child>
                  <NuxtLink :to="`/applications/${app.id}`"> View Details </NuxtLink>
                </UiDropdownMenuItem>
                <UiDropdownMenuItem as-child>
                  <NuxtLink :to="`/applications/${app.id}/keys`"> Manage API Keys </NuxtLink>
                </UiDropdownMenuItem>
                <UiDropdownMenuItem> View Analytics </UiDropdownMenuItem>
                <UiDropdownMenuSeparator />
                <UiDropdownMenuItem class="text-red-600"> Delete Application </UiDropdownMenuItem>
              </UiDropdownMenuContent>
            </UiDropdownMenu>
          </div>
        </UiCardHeader>
        <UiCardContent>
          <div class="grid grid-cols-1 md:grid-cols-3 gap-4 text-sm">
            <div>
              <p class="text-muted-foreground">API Keys</p>
              <p class="font-medium">{{ app.apiKeyCount }}</p>
            </div>
            <div>
              <p class="text-muted-foreground">Monthly Requests</p>
              <p class="font-medium">{{ formatNumber(app.monthlyRequests) }}</p>
            </div>
            <div>
              <p class="text-muted-foreground">Last Used</p>
              <p class="font-medium">
                {{ new Date(app.lastUsed).toLocaleDateString() }}
              </p>
            </div>
          </div>
        </UiCardContent>
      </UiCard>
    </div>

    <!-- Create Application Modal -->
    <ModalsCreateApplication v-model:open="showCreateModal" @created="onApplicationCreated" />
  </div>
</template>
