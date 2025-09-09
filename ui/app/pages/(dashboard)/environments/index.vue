<script lang="ts">
const environments = ref<EnvironmentData[]>([
  {
    id: "env_1",
    name: "Production",
    description: "Live production environment for customer-facing applications",
    status: "active",
    memberCount: 15,
    applicationCount: 8,
    teams: ["Engineering", "DevOps"],
    owner: "Alice Johnson",
    monthlyRequests: 2500000,
    createdAt: "2024-10-01T08:00:00Z",
    lastActivity: "2025-01-15T14:30:00Z"
  },
  {
    id: "env_2",
    name: "Staging",
    description: "Pre-production testing and validation environment",
    status: "active",
    memberCount: 22,
    applicationCount: 12,
    teams: ["Engineering", "QA", "Product"],
    owner: "Bob Smith",
    monthlyRequests: 180000,
    createdAt: "2024-11-15T10:30:00Z",
    lastActivity: "2025-01-15T16:45:00Z"
  },
  {
    id: "env_3",
    name: "Development",
    description: "Development and experimentation environment",
    status: "active",
    memberCount: 12,
    applicationCount: 5,
    teams: ["Engineering"],
    owner: "Carol Williams",
    monthlyRequests: 95000,
    createdAt: "2024-09-20T09:15:00Z",
    lastActivity: "2025-01-15T11:20:00Z"
  },
  {
    id: "env_4",
    name: "Client Demo Hub",
    description: "Dedicated environment for client demonstrations and POCs",
    status: "inactive",
    memberCount: 8,
    applicationCount: 3,
    teams: ["Sales", "Customer Success"],
    owner: "David Brown",
    monthlyRequests: 12000,
    createdAt: "2024-12-10T11:00:00Z",
    lastActivity: "2025-01-14T10:30:00Z"
  },
  {
    id: "env_5",
    name: "AI Research Lab",
    description: "Environment for AI model research and experimentation",
    status: "active",
    memberCount: 6,
    applicationCount: 4,
    teams: ["Engineering", "Analytics"],
    owner: "Emma Davis",
    monthlyRequests: 45000,
    createdAt: "2024-08-15T14:20:00Z",
    lastActivity: "2024-12-20T15:45:00Z"
  }
])
</script>

<script setup lang="ts">
import {
  Globe,
  MoreVertical,
  Edit,
  Trash2,
  CheckCircle,
  XCircle,
  Layers,
  Users,
  Activity,
  Eye,
  UserPlus,
  Copy
} from "lucide-vue-next"
import type { FilterConfig, SearchConfig, DisplayConfig } from "@/components/SearchFilter.vue"
import type { StatsCardProps } from "@/components/Cards/Stats.vue"

useSeoMeta({ title: "Environments - LLM Gateway" })

// Filter state for SearchFilter component
const activeFilters = ref<Record<string, string>>({})

const showEditModal = ref(false)
const editingEnvironment = ref<EnvironmentData | null>(null)

// Global modals
const { openCreateEnvironmentModal } = useGlobalModals()

// Delete modal state
const showDeleteModal = ref(false)
const deletingEnvironment = ref<EnvironmentData | null>(null)
const deleteLoading = ref(false)

// Team assignment modal
const showTeamModal = ref(false)
const assigningEnvironment = ref<EnvironmentData | null>(null)

// Stats configuration
const stats: StatsCardProps[] = [
  {
    title: "Total Environments",
    value: environments.value.length.toString(),
    icon: Globe,
    description: "All registered environments",
    variant: "default"
  },
  {
    title: "Active Environments",
    value: environments.value.filter((env) => env.status === "active").length.toString(),
    icon: CheckCircle,
    description: "Currently running",
    variant: "chart-2"
  },
  {
    title: "Total Applications",
    value: environments.value.reduce((sum, env) => sum + env.applicationCount, 0).toString(),
    icon: Layers,
    description: "Across all environments",
    variant: "chart-3"
  },
  {
    title: "Monthly Requests",
    value: (environments.value.reduce((sum, env) => sum + env.monthlyRequests, 0) / 1000000).toFixed(1) + "M",
    icon: Activity,
    description: "Total API requests",
    variant: "chart-1"
  }
]

// Get unique teams for filtering
const allTeams = [...new Set(environments.value.flatMap((env) => env.teams))]

const filterConfigs: FilterConfig[] = [
  {
    key: "name",
    label: "Environment",
    options: environments.value.map((env) => ({ value: env.name, label: env.name, icon: Globe }))
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
    options: allTeams.map((team) => ({ value: team, label: team, icon: Users }))
  }
]

const searchConfig: SearchConfig<EnvironmentData> = {
  fields: ["name", "description", "owner"],
  placeholder: "Search environments, filter by status, teams..."
}

const displayConfig: DisplayConfig<EnvironmentData> = {
  getItemText: (env) => `${env.name} - ${env.description}`,
  getItemValue: (env) => env.name,
  getItemIcon: () => Globe
}

// Event handlers for SearchFilter
function handleFiltersChanged(filters: Record<string, string>) {
  activeFilters.value = filters
}

function handleItemSelected(environment: EnvironmentData) {
  navigateTo(`/environments/${environment.id}`)
}

// Filtering logic
const filteredEnvironments = computed(() => {
  let filtered = environments.value

  if (activeFilters.value.name && activeFilters.value.name !== "all") {
    filtered = filtered.filter((env) => env.name === activeFilters.value.name)
  }
  if (activeFilters.value.status && activeFilters.value.status !== "all") {
    filtered = filtered.filter((env) => env.status === activeFilters.value.status)
  }
  if (activeFilters.value.team && activeFilters.value.team !== "all") {
    filtered = filtered.filter((env) => env.teams.includes(activeFilters.value.team!))
  }

  return filtered
})

function getStatusColor(status: string) {
  switch (status) {
    case "active":
      return "text-emerald-600 dark:text-emerald-400 bg-emerald-50 dark:bg-emerald-950/30 border-emerald-200 dark:border-emerald-800/50"
    case "inactive":
      return "text-gray-600 dark:text-gray-400 bg-gray-50 dark:bg-gray-900/30 border-gray-200 dark:border-gray-700/50"
    default:
      return "text-gray-600 dark:text-gray-400 bg-gray-50 dark:bg-gray-900/30 border-gray-200 dark:border-gray-700/50"
  }
}

// Modal handlers
const openEditModal = (environment: EnvironmentData) => {
  editingEnvironment.value = environment
  showEditModal.value = true
}

const openDeleteModal = (environment: EnvironmentData) => {
  deletingEnvironment.value = environment
  showDeleteModal.value = true
}

const openTeamModal = (environment: EnvironmentData) => {
  assigningEnvironment.value = environment
  showTeamModal.value = true
}

// Check query parameter to auto-open modal
const route = useRoute()
onMounted(() => {
  if (route.query.create === "environment") {
    openCreateEnvironmentModal()
  }
})

const onTeamsAssigned = (data: { environmentId: string; teamIds: string[] }) => {
  console.log("Teams assigned:", data)
  // Update the environment in the list
  const env = environments.value.find((e) => e.id === data.environmentId)
  if (env) {
    const teamNames = availableTeams
      .map((t) => t.name)
      .filter((name) => availableTeams.find((t) => t.name === name && data.teamIds.includes(t.id)))
    env.teams = teamNames
  }
  showTeamModal.value = false
  assigningEnvironment.value = null
}

// Sample teams for team assignment (would normally come from API)
const availableTeams = [
  { id: "1", name: "Engineering" },
  { id: "2", name: "Product" },
  { id: "3", name: "Marketing" },
  { id: "4", name: "Customer Success" },
  { id: "5", name: "DevOps" },
  { id: "6", name: "QA" },
  { id: "7", name: "Sales" },
  { id: "8", name: "Analytics" }
]
</script>

<template>
  <div class="flex flex-col gap-6">
    <!-- Header -->
    <PageHeader title="Environments" subtext="Create and manage custom environments with team access controls">
      <ButtonsCreate title="Create Environment" :action="openCreateEnvironmentModal" />
    </PageHeader>

    <!-- Stats Cards -->
    <div class="grid gap-4 md:grid-cols-2 lg:grid-cols-4">
      <CardsStats
        v-for="stat in stats"
        :key="stat.title"
        :title="stat.title"
        :value="stat.value"
        :icon="stat.icon"
        :description="stat.description"
        :variant="stat.variant"
      />
    </div>

    <!-- Search & Filter Component -->
    <SearchFilter
      :items="environments"
      :filters="filterConfigs"
      :search-config="searchConfig"
      :display-config="displayConfig"
      @filters-changed="handleFiltersChanged"
      @item-selected="handleItemSelected"
    />

    <!-- Environments List -->
    <CardsDataList title="All Environments" :icon="Globe">
      <template #actions>
        <ButtonsCreate title="Create Environment" :action="openCreateEnvironmentModal" />
      </template>

      <div class="space-y-4">
        <div
          v-for="environment in filteredEnvironments"
          :key="environment.id"
          class="flex items-center justify-between p-6 border rounded-lg hover:bg-muted/50 cursor-pointer transition-colors"
          @click="navigateTo(`/environments/${environment.id}`)"
        >
          <div class="flex items-center gap-6">
            <div class="p-3 rounded-lg bg-primary/10">
              <Globe class="size-6 text-primary" />
            </div>

            <div class="space-y-2">
              <div class="flex items-center gap-3">
                <h3 class="font-semibold text-lg">{{ environment.name }}</h3>
                <div
                  class="inline-flex items-center gap-1 px-2 py-1 rounded-md text-xs font-medium border"
                  :class="getStatusColor(environment.status)"
                >
                  <div
                    class="size-1.5 rounded-full"
                    :class="{
                      'bg-emerald-500': environment.status === 'active',
                      'bg-gray-400': environment.status === 'inactive'
                    }"
                  />
                  {{ environment.status.charAt(0).toUpperCase() + environment.status.slice(1) }}
                </div>
              </div>

              <p class="text-muted-foreground">{{ environment.description }}</p>

              <div class="flex items-center gap-6 text-sm text-muted-foreground">
                <div class="flex items-center gap-1">
                  <Users class="size-4" />
                  <span>{{ environment.memberCount }} members</span>
                </div>
                <div class="flex items-center gap-1">
                  <Layers class="size-4" />
                  <span>{{ environment.applicationCount }} apps</span>
                </div>
                <div class="flex items-center gap-1">
                  <Activity class="size-4" />
                  <span>{{ formatNumberPretty(environment.monthlyRequests) }} req/month</span>
                </div>
              </div>

              <!-- Teams -->
              <div class="flex items-center gap-2">
                <span class="text-xs text-muted-foreground">Teams:</span>
                <div class="flex flex-wrap gap-1">
                  <UiBadge v-for="team in environment.teams" :key="team" variant="secondary" class="text-xs">
                    {{ team }}
                  </UiBadge>
                </div>
              </div>

              <!-- Owner Info -->
              <div class="flex items-center gap-4 text-xs text-muted-foreground">
                <span>Owner: {{ environment.owner }}</span>
              </div>
            </div>
          </div>

          <div class="flex items-center gap-2" @click.stop>
            <UiDropdownMenu>
              <UiDropdownMenuTrigger as-child>
                <UiButton variant="ghost" size="sm">
                  <MoreVertical class="size-4" />
                </UiButton>
              </UiDropdownMenuTrigger>
              <UiDropdownMenuContent align="end" class="w-56">
                <UiDropdownMenuItem @click="navigateTo(`/environments/${environment.id}`)">
                  <Eye class="mr-2 size-4" />
                  View Details
                </UiDropdownMenuItem>
                <UiDropdownMenuItem @click="openTeamModal(environment)">
                  <UserPlus class="mr-2 size-4" />
                  Manage Teams
                </UiDropdownMenuItem>
                <UiDropdownMenuItem @click="navigateTo(`/environments/${environment.id}/applications`)">
                  <Layers class="mr-2 size-4" />
                  Manage Applications
                </UiDropdownMenuItem>
                <UiDropdownMenuSeparator />
                <UiDropdownMenuItem @click="openEditModal(environment)">
                  <Edit class="mr-2 size-4" />
                  Edit Environment
                </UiDropdownMenuItem>
                <UiDropdownMenuItem>
                  <Copy class="mr-2 size-4" />
                  Clone Environment
                </UiDropdownMenuItem>
                <UiDropdownMenuSeparator />
                <UiDropdownMenuItem class="text-destructive" @click="openDeleteModal(environment)">
                  <Trash2 class="mr-2 size-4" />
                  Delete Environment
                </UiDropdownMenuItem>
              </UiDropdownMenuContent>
            </UiDropdownMenu>
          </div>
        </div>
      </div>
    </CardsDataList>

    <!-- Team Assignment Modal -->
    <LazyModalsTeamsAssign
      v-model:open="showTeamModal"
      :environment="assigningEnvironment"
      @assigned="onTeamsAssigned"
    />

    <!-- Delete Environment Modal -->
    <ConfirmationModal
      v-model:open="showDeleteModal"
      title="Delete Environment"
      :description="`Are you sure you want to delete ${deletingEnvironment?.name}? This action cannot be undone and will remove all applications and access controls.`"
      confirm-text="Delete Environment"
      variant="destructive"
      :loading="deleteLoading"
      @confirm="() => {}"
      @cancel="
        () => {
          showDeleteModal = false
          deletingEnvironment = null
        }
      "
    />
  </div>
</template>
