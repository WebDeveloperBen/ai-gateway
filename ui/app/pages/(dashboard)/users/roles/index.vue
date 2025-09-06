<script lang="ts">
import {
  Users,
  Crown,
  Shield,
  UserCheck,
  Eye,
  MoreVertical,
  Edit,
  Trash2,
  UserX,
  CheckCircle,
  XCircle,
  User
} from "lucide-vue-next"
import type { FunctionalComponent } from "vue"
import SearchFilter from "~/components/SearchFilter.vue"
import type { FilterConfig, SearchConfig, DisplayConfig } from "~/components/SearchFilter.vue"
import type { StatsCardProps } from "~/components/Cards/Stats.vue"

interface RoleAssignment {
  userId: string
  userName: string
  userEmail: string
  role: string
  team: string
  status: "active" | "inactive"
  lastActive: string
  avatar?: string
  assignedDate: string
}

interface RoleStats {
  role: string
  totalUsers: number
  activeUsers: number
  teams: string[]
  description: string
  permissions: string[]
}

const roleAssignments: RoleAssignment[] = [
  {
    userId: "1",
    userName: "Alice Johnson",
    userEmail: "alice@company.com",
    role: "Admin",
    team: "Engineering",
    status: "active",
    lastActive: "2 hours ago",
    avatar: "https://images.unsplash.com/photo-1494790108755-2616b95b0c97?w=40&h=40&fit=crop&crop=face",
    assignedDate: "2024-12-01"
  },
  {
    userId: "2",
    userName: "Emma Davis",
    userEmail: "emma@company.com",
    role: "Admin",
    team: "Engineering",
    status: "active",
    lastActive: "30 minutes ago",
    avatar: "https://images.unsplash.com/photo-1472099645785-5658abf4ff4e?w=40&h=40&fit=crop&crop=face",
    assignedDate: "2024-11-15"
  },
  {
    userId: "3",
    userName: "Carol Williams",
    userEmail: "carol@company.com",
    role: "Owner",
    team: "Marketing",
    status: "active",
    lastActive: "1 hour ago",
    avatar: "https://images.unsplash.com/photo-1438761681033-6461ffad8d80?w=40&h=40&fit=crop&crop=face",
    assignedDate: "2024-10-01"
  },
  {
    userId: "4",
    userName: "Bob Smith",
    userEmail: "bob@company.com",
    role: "Developer",
    team: "Engineering",
    status: "active",
    lastActive: "1 day ago",
    assignedDate: "2024-12-10"
  },
  {
    userId: "5",
    userName: "David Brown",
    userEmail: "david@company.com",
    role: "Developer",
    team: "Product",
    status: "active",
    lastActive: "3 hours ago",
    assignedDate: "2024-11-20"
  },
  {
    userId: "6",
    userName: "Frank Wilson",
    userEmail: "frank@company.com",
    role: "Developer",
    team: "Customer Success",
    status: "active",
    lastActive: "2 days ago",
    assignedDate: "2024-12-05"
  },
  {
    userId: "7",
    userName: "Grace Lee",
    userEmail: "grace@company.com",
    role: "Viewer",
    team: "Analytics",
    status: "active",
    lastActive: "4 hours ago",
    assignedDate: "2024-12-15"
  },
  {
    userId: "8",
    userName: "Henry Chen",
    userEmail: "henry@company.com",
    role: "Viewer",
    team: "Marketing",
    status: "inactive",
    lastActive: "1 week ago",
    assignedDate: "2024-11-01"
  }
]

const roleStats: RoleStats[] = [
  {
    role: "Owner",
    totalUsers: 1,
    activeUsers: 1,
    teams: ["Marketing"],
    description: "Full system access with complete administrative privileges",
    permissions: ["All permissions", "Manage users", "Manage teams", "Manage policies", "Billing access"]
  },
  {
    role: "Admin",
    totalUsers: 2,
    activeUsers: 2,
    teams: ["Engineering"],
    description: "Administrative access with user and team management capabilities",
    permissions: ["Manage users", "Manage teams", "View analytics", "Manage API keys"]
  },
  {
    role: "Developer",
    totalUsers: 3,
    activeUsers: 3,
    teams: ["Engineering", "Product", "Customer Success"],
    description: "Development access with API usage and limited management features",
    permissions: ["API access", "View analytics", "Manage own API keys"]
  },
  {
    role: "Viewer",
    totalUsers: 2,
    activeUsers: 1,
    teams: ["Analytics", "Marketing"],
    description: "Read-only access for monitoring and reporting purposes",
    permissions: ["View dashboards", "View Applications", "View analytics", "Export reports"]
  }
]

const stats: StatsCardProps[] = [
  { title: "Total Role Assignments", value: "8", icon: Users, description: "Across all roles", variant: "chart-1" },
  { title: "Active Assignments", value: "7", icon: CheckCircle, description: "Currently active", variant: "chart-2" },
  {
    title: "Unique Roles",
    value: "4",
    icon: Shield,
    description: "Owner, Admin, Developer, Viewer",
    variant: "chart-3"
  },
  {
    title: "Teams Covered",
    value: "4",
    icon: Crown,
    description: "All teams have role assignments",
    variant: "chart-4"
  }
]

const roles = ["Owner", "Admin", "Developer", "Viewer"]

const filterConfigs: FilterConfig[] = [
  {
    key: "role",
    label: "Role",
    options: roles.map((role) => ({
      value: role,
      label: role,
      icon: getRoleIcon(role)
    }))
  },
  {
    key: "status",
    label: "Status",
    options: [
      { value: "active", label: "Active", icon: CheckCircle },
      { value: "inactive", label: "Inactive", icon: XCircle }
    ]
  }
]

const searchConfig: SearchConfig<RoleAssignment> = {
  fields: ["role"],
  placeholder: "Search by role..."
}
</script>

<script setup lang="ts">
useSeoMeta({ title: "Roles - LLM Gateway" })

const activeFilters = ref<Record<string, string>>({})
const selectedRole = ref<string>("all")

const displayConfig: DisplayConfig<RoleAssignment> = {
  getItemText: (assignment) => `${assignment.role}`,
  getItemValue: (assignment) => assignment.role,
  getItemIcon: (assignment) => getRoleIcon(assignment.role)
}

function handleFiltersChanged(filters: Record<string, string>) {
  activeFilters.value = filters
}

function handleItemSelected(assignment: RoleAssignment) {
  console.log("Selected assignment:", assignment)
}

const filteredAssignments = computed(() => {
  let filtered = roleAssignments

  if (selectedRole.value && selectedRole.value !== "all") {
    filtered = filtered.filter((assignment) => assignment.role === selectedRole.value)
  }

  if (activeFilters.value.role && activeFilters.value.role !== "all") {
    filtered = filtered.filter((assignment) => assignment.role === activeFilters.value.role)
  }
  if (activeFilters.value.status && activeFilters.value.status !== "all") {
    filtered = filtered.filter((assignment) => assignment.status === activeFilters.value.status)
  }
  return filtered
})

const currentRoleStats = computed(() => {
  if (selectedRole.value === "all") return null
  return roleStats.find((stat) => stat.role === selectedRole.value)
})

const showAssignRoleModal = ref(false)
const showEditModal = ref(false)
const editingAssignment = ref<RoleAssignment | null>(null)

const openEditModal = (assignment: RoleAssignment) => {
  editingAssignment.value = assignment
  showEditModal.value = true
}

const onRoleAssigned = (assignmentData: any) => {
  console.log("Role assigned:", assignmentData)
  // TODO: Refresh the role assignments list or add to local state
}
</script>

<template>
  <div class="flex flex-col gap-6">
    <!-- Header -->
    <PageHeader title="Roles" subtext="Manage user roles and permissions across teams">
      <div class="flex items-center gap-2">
        <ButtonsCreate title="Assign Role" :action="() => (showAssignRoleModal = true)" />
      </div>
    </PageHeader>

    <!-- Role Overview Cards -->
    <div class="grid gap-4 md:grid-cols-2 lg:grid-cols-4">
      <UiCard
        v-for="role in roleStats"
        :key="role.role"
        class="cursor-pointer transition-all hover:shadow-md"
        :class="selectedRole === role.role ? 'ring-2 ring-primary' : ''"
        @click="selectedRole = selectedRole === role.role ? 'all' : role.role"
      >
        <UiCardHeader class="pb-2">
          <div class="flex items-center justify-between">
            <div class="flex items-center gap-2">
              <component :is="getRoleIcon(role.role)" class="size-5" :class="getRoleColor(role.role)" />
              <UiCardTitle class="text-base">{{ role.role }}</UiCardTitle>
            </div>
            <UiBadge variant="outline" class="text-xs"> {{ role.activeUsers }}/{{ role.totalUsers }} </UiBadge>
          </div>
        </UiCardHeader>
        <UiCardContent class="space-y-3">
          <p class="text-sm text-muted-foreground">{{ role.description }}</p>
          <div class="space-y-2">
            <div class="flex items-center gap-2 text-sm">
              <Users class="size-4 text-muted-foreground" />
              <span class="text-muted-foreground">{{ role.totalUsers }} users</span>
            </div>
            <div class="flex items-center gap-2 text-sm">
              <Shield class="size-4 text-muted-foreground" />
              <span class="text-muted-foreground">{{ role.teams.length }} teams</span>
            </div>
          </div>
          <div class="pt-2">
            <div class="text-xs text-muted-foreground mb-1">Key permissions:</div>
            <div class="flex flex-wrap gap-1">
              <UiBadge
                v-for="permission in role.permissions.slice(0, 2)"
                :key="permission"
                variant="secondary"
                class="text-xs"
              >
                {{ permission }}
              </UiBadge>
              <UiBadge v-if="role.permissions.length > 2" variant="secondary" class="text-xs">
                +{{ role.permissions.length - 2 }}
              </UiBadge>
            </div>
          </div>
        </UiCardContent>
      </UiCard>
    </div>

    <!-- Role Details Panel -->
    <div v-if="currentRoleStats" class="rounded-lg border bg-muted/20 p-4">
      <div class="flex items-center justify-between mb-4">
        <div class="flex items-center gap-3">
          <component
            :is="getRoleIcon(currentRoleStats.role)"
            class="size-6"
            :class="getRoleColor(currentRoleStats.role)"
          />
          <div>
            <h3 class="text-lg font-semibold">{{ currentRoleStats.role }} Role</h3>
            <p class="text-sm text-muted-foreground">{{ currentRoleStats.description }}</p>
          </div>
        </div>
        <UiButton variant="outline" size="sm" @click="selectedRole = 'all'"> View All </UiButton>
      </div>

      <div class="grid gap-4 md:grid-cols-3">
        <div class="space-y-2">
          <h4 class="text-sm font-medium text-muted-foreground">User Statistics</h4>
          <div class="text-2xl font-bold">{{ currentRoleStats.activeUsers }}</div>
          <div class="text-xs text-muted-foreground">
            {{ currentRoleStats.activeUsers }}/{{ currentRoleStats.totalUsers }} active
          </div>
        </div>
        <div class="space-y-2">
          <h4 class="text-sm font-medium text-muted-foreground">Teams</h4>
          <div class="flex flex-wrap gap-1">
            <UiBadge v-for="team in currentRoleStats.teams" :key="team" variant="outline" class="text-xs">
              {{ team }}
            </UiBadge>
          </div>
        </div>
        <div class="space-y-2">
          <h4 class="text-sm font-medium text-muted-foreground">Permissions</h4>
          <div class="space-y-1 max-h-20 overflow-y-auto">
            <div
              v-for="permission in currentRoleStats.permissions"
              :key="permission"
              class="text-xs text-muted-foreground flex items-center gap-1"
            >
              <CheckCircle class="size-3 text-green-600" />
              {{ permission }}
            </div>
          </div>
        </div>
      </div>
    </div>

    <!-- Search & Filter Component -->
    <SearchFilter
      :items="roleAssignments"
      :filters="filterConfigs"
      :search-config="searchConfig"
      :display-config="displayConfig"
      @filters-changed="handleFiltersChanged"
      @item-selected="handleItemSelected"
    />

    <!-- Role Assignments Table -->
    <CardsDataList
      :title="selectedRole === 'all' ? 'All Role Assignments' : `${selectedRole} Assignments`"
      :icon="Shield"
    >
      <template #actions>
        <ButtonsCreate title="Assign Role" :action="() => (showAssignRoleModal = true)" />
      </template>

      <div class="space-y-4">
        <div
          v-for="assignment in filteredAssignments"
          :key="assignment.userId"
          class="flex items-center justify-between p-4 border rounded-lg hover:bg-muted/50"
        >
          <div class="flex items-center gap-4">
            <UiAvatar class="size-10">
              <UiAvatarImage v-if="assignment.avatar" :src="assignment.avatar" :alt="assignment.userName" />
              <UiAvatarFallback>{{
                assignment.userName
                  .split(" ")
                  .map((n) => n[0])
                  .join("")
              }}</UiAvatarFallback>
            </UiAvatar>

            <div class="space-y-1">
              <div class="flex items-center gap-2">
                <p class="font-medium">{{ assignment.userName }}</p>
                <div
                  class="inline-flex items-center gap-1 px-2 py-1 rounded-md text-xs font-medium border"
                  :class="getRoleColor(assignment.role)"
                >
                  <component :is="getRoleIcon(assignment.role)" class="size-3" />
                  {{ assignment.role }}
                </div>
              </div>
              <div class="flex items-center gap-4 text-sm text-muted-foreground">
                <span>{{ assignment.userEmail }}</span>
                <span>•</span>
                <span>{{ assignment.team }} team</span>
                <span>•</span>
                <span>Assigned {{ new Date(assignment.assignedDate).toLocaleDateString() }}</span>
              </div>
            </div>
          </div>

          <div class="flex items-center gap-2">
            <div
              class="inline-flex items-center gap-1 px-2 py-1 rounded-md text-xs font-medium border"
              :class="getStatusColor(assignment.status)"
            >
              <div
                class="size-1.5 rounded-full"
                :class="assignment.status === 'active' ? 'bg-green-500' : 'bg-gray-400'"
              />
              {{ assignment.status === "active" ? "Active" : "Inactive" }}
            </div>

            <UiDropdownMenu>
              <UiDropdownMenuTrigger as-child>
                <UiButton variant="ghost" size="sm">
                  <MoreVertical class="size-4" />
                </UiButton>
              </UiDropdownMenuTrigger>
              <UiDropdownMenuContent align="end" class="w-48">
                <UiDropdownMenuItem @click="openEditModal(assignment)">
                  <Edit class="mr-2 size-4" />
                  Change Role
                </UiDropdownMenuItem>
                <UiDropdownMenuItem @click="() => {}">
                  <component :is="assignment.status === 'active' ? UserX : UserCheck" class="mr-2 size-4" />
                  {{ assignment.status === "active" ? "Deactivate" : "Activate" }}
                </UiDropdownMenuItem>
                <UiDropdownMenuSeparator />
                <UiDropdownMenuItem class="text-destructive" @click="() => {}">
                  <Trash2 class="mr-2 size-4" />
                  Remove Role
                </UiDropdownMenuItem>
              </UiDropdownMenuContent>
            </UiDropdownMenu>
          </div>
        </div>

        <div v-if="filteredAssignments.length === 0" class="text-center py-8 text-muted-foreground">
          <Shield class="size-12 mx-auto mb-4 text-muted-foreground/50" />
          <h3 class="text-lg font-medium mb-2">No role assignments found</h3>
          <p class="text-sm mb-4">No users match the current filters or selected role.</p>
          <UiButton
            variant="outline"
            @click="
              () => {
                selectedRole = 'all'
                activeFilters = {}
              }
            "
          >
            Clear Filters
          </UiButton>
        </div>
      </div>
    </CardsDataList>

    <!-- Assign Role Modal -->
    <ModalsAssignRole v-model:open="showAssignRoleModal" @assigned="onRoleAssigned" />
  </div>
</template>
