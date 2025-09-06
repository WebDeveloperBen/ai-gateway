<script lang="ts">
interface ApplicationData {
  id: string
  name: string
  description: string
  status: "active" | "inactive"
  apiKeyCount: number
  monthlyRequests: number
  lastUsed: string
  models: string[]
  team: string
}

const applications = ref<ApplicationData[]>([
  {
    id: "app_1",
    name: "Customer Service Bot",
    description: "AI-powered customer support assistant for handling common inquiries",
    status: "active" as const,
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
    status: "active" as const,
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
    status: "inactive" as const,
    apiKeyCount: 1,
    monthlyRequests: 0,
    lastUsed: "2024-12-20T09:15:00Z",
    models: ["gpt-4"],
    team: "Engineering"
  }
])
</script>
<script setup lang="ts">
import { Plus, Key, Activity, CheckCircle, XCircle, Layers, Users } from "lucide-vue-next"
import type { StatsCardProps } from "@/components/Cards/Stats.vue"
import SearchFilter from "@/components/SearchFilter.vue"
import type { FilterConfig, SearchConfig, DisplayConfig } from "@/components/SearchFilter.vue"

const showCreateModal = ref(false)

// Filter state for SearchFilter component
const activeFilters = ref<Record<string, string>>({})

// Get app config
const appConfig = useAppConfig()

const handleApplicationSelect = (app: ApplicationData) => {
  navigateTo(`/applications/${app.id}`)
}

const onApplicationCreated = (applicationId: string) => {
  console.log("Application created:", applicationId)
  // Clear the query parameter after creation
  navigateTo("/applications", { replace: true })
  // Here you could refresh the applications list or add the new app to the existing list
}

// Check query parameter to auto-open modal
const route = useRoute()
onMounted(() => {
  if (route.query.create === "application") {
    showCreateModal.value = true
  }
})

// SearchFilter configuration
const teams = [...new Set(applications.value.map((app) => app.team))]

const filterConfigs: FilterConfig[] = [
  {
    key: "name",
    label: "Application",
    options: applications.value.map((app) => ({ value: app.name, label: app.name, icon: Layers }))
  },
  {
    key: "status",
    label: "Status",
    options: [
      { value: "active", label: "Active", icon: CheckCircle },
      { value: "inactive", label: "Inactive", icon: XCircle }
    ]
  },
  {
    key: "team",
    label: "Team",
    options: teams.map((team) => ({ value: team, label: team, icon: Users }))
  }
]

const searchConfig: SearchConfig<ApplicationData> = {
  fields: ["name", "description"],
  placeholder: "Search applications, filter by status, team..."
}

const displayConfig: DisplayConfig<ApplicationData> = {
  getItemText: (app) => `${app.name} - ${app.description}`,
  getItemValue: (app) => app.name,
  getItemIcon: () => Layers
}

// Event handlers for SearchFilter
function handleFiltersChanged(filters: Record<string, string>) {
  activeFilters.value = filters
}

function handleItemSelected(app: ApplicationData) {
  navigateTo(`/applications/${app.id}`)
}

const filteredApplications = computed(() => {
  let filtered = applications.value

  // Apply SearchFilter filters
  if (activeFilters.value.name && activeFilters.value.name !== "all") {
    filtered = filtered.filter((app) => app.name === activeFilters.value.name)
  }
  if (activeFilters.value.status && activeFilters.value.status !== "all") {
    filtered = filtered.filter((app) => app.status === activeFilters.value.status)
  }
  if (activeFilters.value.team && activeFilters.value.team !== "all") {
    filtered = filtered.filter((app) => app.team === activeFilters.value.team)
  }

  return filtered
})

// Icon and variant arrays for index matching with API data
const statsIcons = [Layers, CheckCircle, Key, Activity]
const statsVariants: StatsCardProps["variant"][] = ["default", "chart-2", "chart-3", "chart-1"]

// Stats configuration - ready for API replacement
const statsCards = computed(() => {
  // This can be replaced with: const { data: statsData } = await useFetch('/api/applications/stats')
  // Then: statsData.map((stat, index) => ({ ...stat, icon: statsIcons[index], variant: statsVariants[index] }))

  return [
    {
      title: "Total Applications",
      value: applications.value.length,
      description: "All registered applications"
    },
    {
      title: "Active Applications",
      value: applications.value.filter((app) => app.status === "active").length,
      description: "Currently running applications"
    },
    {
      title: "Total API Keys",
      value: applications.value.reduce((sum, app) => sum + app.apiKeyCount, 0),
      description: "Across all applications"
    },
    {
      title: "Monthly Requests",
      value: applications.value.reduce((sum, app) => sum + app.monthlyRequests, 0),
      description: "Total API calls this month"
    }
  ].map((stat, index) => ({
    ...stat,
    icon: statsIcons[index],
    variant: statsVariants[index]
  }))
})
</script>

<template>
  <div class="flex flex-col gap-6">
    <!-- Header -->
    <PageHeader
      title="Applications"
      :subtext="`Manage your AI-powered applications and their ${appConfig.app.name} access`"
    >
      <UiButton @click="showCreateModal = true" class="gap-2">
        <Plus class="h-4 w-4" />
        New Application
      </UiButton>
    </PageHeader>

    <!-- Stats Cards -->
    <div class="grid grid-cols-1 sm:grid-cols-2 lg:grid-cols-4 gap-6">
      <CardsStats
        v-for="card in statsCards"
        :key="card.title"
        :title="card.title"
        :value="card.value"
        :icon="card.icon"
        :description="card.description"
        :variant="card.variant"
      />
    </div>

    <!-- Search & Filter Component -->
    <SearchFilter
      :items="applications"
      :filters="filterConfigs"
      :search-config="searchConfig"
      :display-config="displayConfig"
      @filters-changed="handleFiltersChanged"
      @item-selected="handleItemSelected"
    />

    <ApplicationsList :applications="filteredApplications" @select-application="handleApplicationSelect" />

    <!-- Create Application Modal -->
    <LazyModalsApplicationsCreate v-model:open="showCreateModal" @created="onApplicationCreated" />
  </div>
</template>
