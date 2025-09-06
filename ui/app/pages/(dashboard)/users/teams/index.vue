<script lang="ts">
interface TeamData {
  id: string
  name: string
  description: string
  status: "active" | "inactive"
  memberCount: number
  owner: string
  adminCount: number
  developerCount: number
  viewerCount: number
  policies: string[]
  costCenter: string
  createdAt: string
  lastActivity: string
}

const teams = ref<TeamData[]>([
  {
    id: "team_1",
    name: "Engineering",
    description: "Core engineering team for product development",
    status: "active",
    memberCount: 12,
    owner: "Alice Johnson",
    adminCount: 2,
    developerCount: 8,
    viewerCount: 2,
    policies: ["Development Policy", "Security Policy"],
    costCenter: "ENG-001",
    createdAt: "2024-12-01T08:00:00Z",
    lastActivity: "2025-01-15T14:30:00Z"
  },
  {
    id: "team_2",
    name: "Marketing",
    description: "Content creation and marketing campaigns",
    status: "active",
    memberCount: 6,
    owner: "Carol Williams",
    adminCount: 1,
    developerCount: 2,
    viewerCount: 3,
    policies: ["Marketing Policy"],
    costCenter: "MKT-002",
    createdAt: "2024-11-15T10:30:00Z",
    lastActivity: "2025-01-14T16:45:00Z"
  },
  {
    id: "team_3",
    name: "Customer Success",
    description: "Customer support and success operations",
    status: "active",
    memberCount: 8,
    owner: "David Brown",
    adminCount: 1,
    developerCount: 3,
    viewerCount: 4,
    policies: ["Customer Data Policy", "Security Policy"],
    costCenter: "CS-003",
    createdAt: "2024-10-20T09:15:00Z",
    lastActivity: "2025-01-15T11:20:00Z"
  },
  {
    id: "team_4",
    name: "Analytics",
    description: "Data analysis and business intelligence",
    status: "inactive",
    memberCount: 3,
    owner: "Emma Davis",
    adminCount: 1,
    developerCount: 1,
    viewerCount: 1,
    policies: ["Data Policy"],
    costCenter: "ANL-004",
    createdAt: "2024-09-10T11:00:00Z",
    lastActivity: "2024-12-22T10:30:00Z"
  }
])
</script>
<script setup lang="ts">
import {
  Users,
  Plus,
  MoreVertical,
  Edit,
  Trash2,
  UserCheck,
  UserX,
  Crown,
  Shield,
  CheckCircle,
  XCircle,
  Eye,
  Building,
  Settings,
  Calendar,
  Activity,
  FileText,
  View,
  ViewIcon
} from "lucide-vue-next"
import type { FunctionalComponent } from "vue"
import SearchFilter from "~/components/SearchFilter.vue"
import type { FilterConfig, SearchConfig, DisplayConfig } from "~/components/SearchFilter.vue"
import type { StatsCardProps } from "~/components/Cards/Stats.vue"

useSeoMeta({ title: "Teams - LLM Gateway" })

// Filter state for SearchFilter component
const activeFilters = ref<Record<string, string>>({})

const showCreateModal = ref(false)
const showEditModal = ref(false)
const editingTeam = ref<TeamData | null>(null)
const editLoading = ref(false)

// Delete modal state
const showDeleteModal = ref(false)
const deletingTeam = ref<TeamData | null>(null)
const deleteLoading = ref(false)

// Stats configuration
const stats: StatsCardProps[] = [
  {
    title: "Total Teams",
    value: teams.value.length.toString(),
    icon: Building,
    description: "All registered teams",
    variant: "chart-1"
  },
  {
    title: "Active Teams",
    value: teams.value.filter((team) => team.status === "active").length.toString(),
    icon: CheckCircle,
    description: "Currently active teams",
    variant: "chart-2"
  },
  {
    title: "Total Members",
    value: teams.value.reduce((sum, team) => sum + team.memberCount, 0).toString(),
    icon: Users,
    description: "Across all teams",
    variant: "chart-3"
  },
  {
    title: "Active Policies",
    value: [...new Set(teams.value.flatMap((team) => team.policies))].length.toString(),
    icon: FileText,
    description: "Unique policies assigned",
    variant: "chart-4"
  }
]

// SearchFilter configuration
const costCenters = [...new Set(teams.value.map((team) => team.costCenter))]
const policies = [...new Set(teams.value.flatMap((team) => team.policies))]

const filterConfigs: FilterConfig[] = [
  {
    key: "name",
    label: "Team",
    options: teams.value.map((team) => ({ value: team.name, label: team.name, icon: Building }))
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
    key: "costCenter",
    label: "Cost Center",
    options: costCenters.map((center) => ({ value: center, label: center, icon: Settings }))
  }
]

const searchConfig: SearchConfig<TeamData> = {
  fields: ["name", "description", "owner"],
  placeholder: "Search teams, filter by status, cost center..."
}

const displayConfig: DisplayConfig<TeamData> = {
  getItemText: (team) => `${team.name} - ${team.description}`,
  getItemValue: (team) => team.name,
  getItemIcon: () => Building
}

// Event handlers for SearchFilter
function handleFiltersChanged(filters: Record<string, string>) {
  activeFilters.value = filters
}

function handleItemSelected(team: TeamData) {
  navigateTo(`/users/teams/${team.id}`)
}

// Filtering logic
const filteredTeams = computed(() => {
  let filtered = teams.value

  if (activeFilters.value.name && activeFilters.value.name !== "all") {
    filtered = filtered.filter((team) => team.name === activeFilters.value.name)
  }
  if (activeFilters.value.status && activeFilters.value.status !== "all") {
    filtered = filtered.filter((team) => team.status === activeFilters.value.status)
  }
  if (activeFilters.value.costCenter && activeFilters.value.costCenter !== "all") {
    filtered = filtered.filter((team) => team.costCenter === activeFilters.value.costCenter)
  }

  return filtered
})

// Helper functions
function getRoleIcon(role: string): FunctionalComponent {
  switch (role) {
    case "Owner":
      return Crown
    case "Admin":
      return Shield
    case "Developer":
      return UserCheck
    case "Viewer":
      return Eye
    default:
      return Users
  }
}

function getStatusColor(status: string) {
  return status === "active"
    ? "text-green-600 bg-green-50 border-green-200"
    : "text-gray-600 bg-gray-50 border-gray-200"
}

// Modal handlers
const openEditModal = (team: TeamData) => {
  editingTeam.value = team
  showEditModal.value = true
}

const openDeleteModal = (team: TeamData) => {
  deletingTeam.value = team
  showDeleteModal.value = true
}

// Check query parameter to auto-open modal
const route = useRoute()
onMounted(() => {
  if (route.query.create === "team") {
    showCreateModal.value = true
  }
})

const onTeamCreated = (teamData: any) => {
  console.log("Team created:", teamData)
  navigateTo("/teams", { replace: true })
}
</script>

<template>
  <div class="flex flex-col gap-6">
    <!-- Header -->
    <PageHeader title="Teams" subtext="Organize users into teams with role-based access and policy management">
      <UiButton @click="showCreateModal = true" class="gap-2">
        <Plus class="h-4 w-4" />
        Create Team
      </UiButton>
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
      :items="teams"
      :filters="filterConfigs"
      :search-config="searchConfig"
      :display-config="displayConfig"
      @filters-changed="handleFiltersChanged"
      @item-selected="handleItemSelected"
    />

    <!-- Teams List -->
    <div class="space-y-4">
      <div
        v-for="team in filteredTeams"
        :key="team.id"
        class="flex items-center justify-between p-6 border rounded-lg hover:bg-muted/50 cursor-pointer transition-colors"
        @click="navigateTo(`/users/teams/${team.id}`)"
      >
        <div class="flex items-center gap-6">
          <div class="p-3 rounded-lg bg-primary/10">
            <Building class="size-6 text-primary" />
          </div>

          <div class="space-y-2">
            <div class="flex items-center gap-3">
              <h3 class="font-semibold text-lg">{{ team.name }}</h3>
              <div
                class="inline-flex items-center gap-1 px-2 py-1 rounded-md text-xs font-medium border"
                :class="getStatusColor(team.status)"
              >
                <div class="size-1.5 rounded-full" :class="team.status === 'active' ? 'bg-green-500' : 'bg-gray-400'" />
                {{ team.status === "active" ? "Active" : "Inactive" }}
              </div>
            </div>

            <p class="text-muted-foreground">{{ team.description }}</p>

            <div class="flex items-center gap-6 text-sm text-muted-foreground">
              <div class="flex items-center gap-1">
                <Users class="size-4" />
                <span>{{ team.memberCount }} members</span>
              </div>
              <div class="flex items-center gap-1">
                <Crown class="size-4" />
                <span>{{ team.owner }}</span>
              </div>
              <div class="flex items-center gap-1">
                <FileText class="size-4" />
                <span>{{ team.policies.length }} policies</span>
              </div>
            </div>

            <!-- Role Distribution -->
            <div class="flex items-center gap-4 text-xs">
              <div class="flex items-center gap-1">
                <Crown class="size-3 text-orange-600" />
                <span class="text-muted-foreground">1 Owner</span>
              </div>
              <div class="flex items-center gap-1">
                <Shield class="size-3 text-blue-600" />
                <span class="text-muted-foreground">{{ team.adminCount }} Admins</span>
              </div>
              <div class="flex items-center gap-1">
                <UserCheck class="size-3 text-green-600" />
                <span class="text-muted-foreground">{{ team.developerCount }} Developers</span>
              </div>
              <div class="flex items-center gap-1">
                <Eye class="size-3 text-purple-600" />
                <span class="text-muted-foreground">{{ team.viewerCount }} Viewers</span>
              </div>
            </div>

            <!-- Policies -->
            <div v-if="team.policies.length > 0" class="flex items-center gap-2">
              <span class="text-xs text-muted-foreground">Policies:</span>
              <div class="flex flex-wrap gap-1">
                <UiBadge v-for="policy in team.policies" :key="policy" variant="outline" class="text-xs">
                  {{ policy }}
                </UiBadge>
              </div>
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
            <UiDropdownMenuContent align="end" class="w-48">
              <UiDropdownMenuItem @click="navigateTo(`/users/teams/${team.id}`)">
                <Eye class="mr-2 size-4" />
                View Details
              </UiDropdownMenuItem>
              <UiDropdownMenuItem @click="openEditModal(team)">
                <Edit class="mr-2 size-4" />
                Edit Team
              </UiDropdownMenuItem>
              <UiDropdownMenuSeparator />
              <UiDropdownMenuItem class="text-destructive" @click="openDeleteModal(team)">
                <Trash2 class="mr-2 size-4" />
                Delete Team
              </UiDropdownMenuItem>
            </UiDropdownMenuContent>
          </UiDropdownMenu>
        </div>
      </div>
    </div>

    <!-- Create Team Modal -->
    <ModalsCreateTeam v-model:open="showCreateModal" @created="onTeamCreated" />

    <!-- Delete Team Modal -->
    <ConfirmationModal
      v-model:open="showDeleteModal"
      title="Delete Team"
      :description="`Are you sure you want to delete ${deletingTeam?.name}? This action cannot be undone and will remove all team associations.`"
      confirm-text="Delete Team"
      variant="destructive"
      :loading="deleteLoading"
      @confirm="() => {}"
      @cancel="
        () => {
          showDeleteModal = false
          deletingTeam = null
        }
      "
    />
  </div>
</template>
