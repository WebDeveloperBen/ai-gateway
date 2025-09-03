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
  X,
  Circle,
  CheckCircle,
  XCircle,
  User
} from "lucide-vue-next"
import type { FunctionalComponent } from "vue"

useSeoMeta({ title: "Users - LLM Gateway" })

interface User {
  id: string
  name: string
  email: string
  role: string
  status: "active" | "inactive"
  team: string
  lastActive: string
  avatar?: string
}

const searchQuery = ref("")
const selectedRole = ref("all")
const selectedTeam = ref("all")
const selectedStatus = ref("all")
const selectedUser = ref("all")
const showFilters = ref(false)

// Sample data
const users: User[] = [
  {
    id: "1",
    name: "Alice Johnson",
    email: "alice@company.com",
    role: "Admin",
    status: "active",
    team: "Engineering",
    lastActive: "2 hours ago",
    avatar: "https://images.unsplash.com/photo-1494790108755-2616b95b0c97?w=40&h=40&fit=crop&crop=face"
  },
  {
    id: "2",
    name: "Bob Smith",
    email: "bob@company.com",
    role: "Developer",
    status: "active",
    team: "Engineering",
    lastActive: "1 day ago"
  },
  {
    id: "3",
    name: "Carol Williams",
    email: "carol@company.com",
    role: "Viewer",
    status: "inactive",
    team: "Marketing",
    lastActive: "1 week ago",
    avatar: "https://images.unsplash.com/photo-1438761681033-6461ffad8d80?w=40&h=40&fit=crop&crop=face"
  },
  {
    id: "4",
    name: "David Brown",
    email: "david@company.com",
    role: "Developer",
    status: "active",
    team: "Product",
    lastActive: "3 hours ago"
  },
  {
    id: "5",
    name: "Emma Davis",
    email: "emma@company.com",
    role: "Admin",
    status: "active",
    team: "Engineering",
    lastActive: "30 minutes ago",
    avatar: "https://images.unsplash.com/photo-1472099645785-5658abf4ff4e?w=40&h=40&fit=crop&crop=face"
  }
]

const stats = [
  { label: "Total Users", value: "24", icon: Users, change: "+2 this month" },
  { label: "Active Users", value: "18", icon: UserCheck, change: "+5% from last month" },
  { label: "Admin Users", value: "3", icon: Crown, change: "No change" },
  { label: "Teams", value: "6", icon: Shield, change: "+1 this month" }
]

const roles = ["Admin", "Developer", "Viewer"]
const teams = ["Engineering", "Product", "Marketing", "Sales", "Design", "Operations"]

function getRoleIcon(role: string): FunctionalComponent {
  switch (role) {
    case "Admin":
      return Crown
    case "Developer":
      return Shield
    default:
      return Users
  }
}

function getRoleColor(role: string) {
  switch (role) {
    case "Admin":
      return "text-orange-600 bg-orange-50 border-orange-200"
    case "Developer":
      return "text-blue-600 bg-blue-50 border-blue-200"
    default:
      return "text-gray-600 bg-gray-50 border-gray-200"
  }
}

function getStatusColor(status: string) {
  return status === "active"
    ? "text-green-600 bg-green-50 border-green-200"
    : "text-gray-600 bg-gray-50 border-gray-200"
}

// Filter functions
function filterByRole(role: string) {
  selectedRole.value = role
  showFilters.value = false
}

function filterByTeam(team: string) {
  selectedTeam.value = team
  showFilters.value = false
}

function filterByStatus(status: string) {
  selectedStatus.value = status
  showFilters.value = false
}

function selectUser(user: User) {
  selectedUser.value = user.id
  searchQuery.value = ""
  showFilters.value = false
}

function clearSearch() {
  searchQuery.value = ""
  selectedRole.value = "all"
  selectedTeam.value = "all"
  selectedStatus.value = "all"
  selectedUser.value = "all"
  showFilters.value = false
}

const hasActiveFilters = computed(() => {
  return (
    selectedRole.value !== "all" ||
    selectedTeam.value !== "all" ||
    selectedStatus.value !== "all" ||
    selectedUser.value !== "all"
  )
})

const filteredUsers = computed(() => {
  let filtered = users

  // Filter by role
  if (selectedRole.value !== "all") {
    filtered = filtered.filter((user) => user.role === selectedRole.value)
  }

  // Filter by team
  if (selectedTeam.value !== "all") {
    filtered = filtered.filter((user) => user.team === selectedTeam.value)
  }

  // Filter by status
  if (selectedStatus.value !== "all") {
    filtered = filtered.filter((user) => user.status === selectedStatus.value)
  }

  // Filter by specific user
  if (selectedUser.value !== "all") {
    filtered = filtered.filter((user) => user.id === selectedUser.value)
  }

  // Apply search query last, but only if no specific user is selected
  if (searchQuery.value && selectedUser.value === "all") {
    filtered = filtered.filter(
      (user) =>
        user.name.toLowerCase().includes(searchQuery.value.toLowerCase()) ||
        user.email.toLowerCase().includes(searchQuery.value.toLowerCase())
    )
  }

  return filtered
})

// Separate computed for dropdown search results - searches users, roles, teams, status
const searchResults = computed(() => {
  if (!searchQuery.value) return { users: [], roles: [], teams: [], statuses: [] }

  const query = searchQuery.value.toLowerCase()

  const matchingUsers = users
    .filter((user) => user.name.toLowerCase().includes(query) || user.email.toLowerCase().includes(query))
    .slice(0, 6)

  const matchingRoles = roles.filter((role) => role.toLowerCase().includes(query))

  const matchingTeams = teams.filter((team) => team.toLowerCase().includes(query))

  const matchingStatuses = [
    { value: "active", label: "Active" },
    { value: "inactive", label: "Inactive" }
  ].filter((status) => status.label.toLowerCase().includes(query))

  return {
    users: matchingUsers,
    roles: matchingRoles,
    teams: matchingTeams,
    statuses: matchingStatuses
  }
})
</script>

<template>
  <div class="flex flex-col gap-6">
    <!-- Header -->
    <div class="flex items-center justify-between">
      <div>
        <h1 class="text-2xl font-semibold tracking-tight">Users</h1>
        <p class="text-muted-foreground">Manage users, roles and team access</p>
      </div>
      <UiButton>
        <Plus class="mr-2 size-4" />
        Invite User
      </UiButton>
    </div>

    <!-- Stats Cards -->
    <div class="grid gap-4 md:grid-cols-2 lg:grid-cols-4">
      <UiCard v-for="stat in stats" :key="stat.label" class="p-6">
        <div class="flex items-center gap-4">
          <div class="rounded-lg bg-primary/10 p-3">
            <component :is="stat.icon" class="size-5 text-primary" />
          </div>
          <div class="flex-1">
            <p class="text-2xl font-semibold">{{ stat.value }}</p>
            <p class="text-sm text-muted-foreground">{{ stat.label }}</p>
          </div>
        </div>
        <div class="mt-4">
          <p class="text-xs text-muted-foreground">{{ stat.change }}</p>
        </div>
      </UiCard>
    </div>

    <!-- Search & Filter Command -->
    <UiCommand class="rounded-lg border shadow-sm">
      <div class="flex w-full items-center px-3" cmdk-input-wrapper="">
        <UiCommandInput
          v-model="searchQuery"
          placeholder="Search users, filter by role, team or status..."
          class="flex h-10 w-full rounded-md bg-transparent py-3 text-sm outline-none placeholder:text-muted-foreground disabled:cursor-not-allowed disabled:opacity-50"
        />
        <UiButton
          v-if="searchQuery"
          variant="ghost"
          size="sm"
          class="h-auto p-1 hover:bg-transparent"
          @click="searchQuery = ''"
        >
          <X class="size-4 text-muted-foreground hover:text-foreground" />
        </UiButton>
      </div>
      <UiCommandList v-if="searchQuery || showFilters">
        <UiCommandEmpty>No users found.</UiCommandEmpty>

        <template v-if="!searchQuery">
          <UiCommandGroup heading="Filter by Role">
            <UiCommandItem
              v-for="role in roles"
              :key="role"
              :value="`role:${role}`"
              :text="role"
              :icon="getRoleIcon(role)"
              @select="filterByRole(role)"
            />
            <UiCommandItem value="role:all" text="All Roles" :icon="Users" @select="filterByRole('all')" />
          </UiCommandGroup>
          <UiCommandSeparator />

          <UiCommandGroup heading="Filter by Team">
            <UiCommandItem
              v-for="team in teams"
              :key="team"
              :value="`team:${team}`"
              :text="team"
              :icon="Users"
              @select="filterByTeam(team)"
            />
            <UiCommandItem value="team:all" text="All Teams" :icon="Users" @select="filterByTeam('all')" />
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
          <UiCommandGroup v-if="searchResults.users.length" heading="Users">
            <UiCommandItem
              v-for="user in searchResults.users"
              :key="user.id"
              :value="user.name"
              :text="`${user.name} - ${user.email}`"
              :icon="User"
              @select="selectUser(user)"
            />
          </UiCommandGroup>

          <UiCommandGroup v-if="searchResults.roles.length" heading="Roles">
            <UiCommandItem
              v-for="role in searchResults.roles"
              :key="role"
              :value="`role:${role}`"
              :text="role"
              :icon="getRoleIcon(role)"
              @select="filterByRole(role)"
            />
          </UiCommandGroup>

          <UiCommandGroup v-if="searchResults.teams.length" heading="Teams">
            <UiCommandItem
              v-for="team in searchResults.teams"
              :key="team"
              :value="`team:${team}`"
              :text="team"
              :icon="Users"
              @select="filterByTeam(team)"
            />
          </UiCommandGroup>

          <UiCommandGroup v-if="searchResults.statuses.length" heading="Status">
            <UiCommandItem
              v-for="status in searchResults.statuses"
              :key="status.value"
              :value="`status:${status.value}`"
              :text="status.label"
              :icon="status.value === 'active' ? CheckCircle : XCircle"
              @select="filterByStatus(status.value)"
            />
          </UiCommandGroup>
        </template>
      </UiCommandList>
    </UiCommand>

    <!-- Active Filters -->
    <div v-if="hasActiveFilters" class="flex items-center justify-between gap-4">
      <div class="flex items-center gap-2">
        <span class="text-sm text-muted-foreground">Active filters:</span>
        <UiBadge v-if="selectedUser !== 'all'" variant="secondary" class="gap-1">
          User: {{ users.find((u) => u.id === selectedUser)?.name }}
          <button @click="selectedUser = 'all'" class="ml-1 hover:bg-muted-foreground/20 rounded-sm">
            <X class="h-3 w-3" />
          </button>
        </UiBadge>
        <UiBadge v-if="selectedRole !== 'all'" variant="secondary" class="gap-1">
          Role: {{ selectedRole }}
          <button @click="selectedRole = 'all'" class="ml-1 hover:bg-muted-foreground/20 rounded-sm">
            <X class="h-3 w-3" />
          </button>
        </UiBadge>
        <UiBadge v-if="selectedTeam !== 'all'" variant="secondary" class="gap-1">
          Team: {{ selectedTeam }}
          <button @click="selectedTeam = 'all'" class="ml-1 hover:bg-muted-foreground/20 rounded-sm">
            <X class="h-3 w-3" />
          </button>
        </UiBadge>
        <UiBadge v-if="selectedStatus !== 'all'" variant="secondary" class="gap-1">
          Status: {{ selectedStatus === "active" ? "Active" : "Inactive" }}
          <button @click="selectedStatus = 'all'" class="ml-1 hover:bg-muted-foreground/20 rounded-sm">
            <X class="h-3 w-3" />
          </button>
        </UiBadge>
      </div>
      <UiButton variant="outline" size="sm" @click="clearSearch">Clear All</UiButton>
    </div>

    <!-- Users Table -->
    <UiCard>
      <div class="p-6">
        <div class="space-y-4">
          <div
            v-for="user in filteredUsers"
            :key="user.id"
            class="flex items-center justify-between p-4 border rounded-lg hover:bg-muted/50"
          >
            <div class="flex items-center gap-4">
              <UiAvatar class="size-10">
                <UiAvatarImage v-if="user.avatar" :src="user.avatar" :alt="user.name" />
                <UiAvatarFallback>{{
                  user.name
                    .split(" ")
                    .map((n) => n[0])
                    .join("")
                }}</UiAvatarFallback>
              </UiAvatar>

              <div class="space-y-1">
                <div class="flex items-center gap-2">
                  <p class="font-medium">{{ user.name }}</p>
                  <div
                    class="inline-flex items-center gap-1 px-2 py-1 rounded-md text-xs font-medium border"
                    :class="getRoleColor(user.role)"
                  >
                    <component :is="getRoleIcon(user.role)" class="size-3" />
                    {{ user.role }}
                  </div>
                </div>
                <div class="flex items-center gap-4 text-sm text-muted-foreground">
                  <span>{{ user.email }}</span>
                  <span>•</span>
                  <span>{{ user.team }}</span>
                  <span>•</span>
                  <span>Last active {{ user.lastActive }}</span>
                </div>
              </div>
            </div>

            <div class="flex items-center gap-2">
              <div
                class="inline-flex items-center gap-1 px-2 py-1 rounded-md text-xs font-medium border"
                :class="getStatusColor(user.status)"
              >
                <div class="size-1.5 rounded-full" :class="user.status === 'active' ? 'bg-green-500' : 'bg-gray-400'" />
                {{ user.status === "active" ? "Active" : "Inactive" }}
              </div>

              <UiDropdownMenu>
                <UiDropdownMenuTrigger as-child>
                  <UiButton variant="ghost" size="sm">
                    <MoreVertical class="size-4" />
                  </UiButton>
                </UiDropdownMenuTrigger>
                <UiDropdownMenuContent align="end" class="w-48">
                  <UiDropdownMenuItem>
                    <Edit class="mr-2 size-4" />
                    Edit User
                  </UiDropdownMenuItem>
                  <UiDropdownMenuItem>
                    <component :is="user.status === 'active' ? UserX : UserCheck" class="mr-2 size-4" />
                    {{ user.status === "active" ? "Deactivate" : "Activate" }}
                  </UiDropdownMenuItem>
                  <UiDropdownMenuSeparator />
                  <UiDropdownMenuItem class="text-destructive">
                    <Trash2 class="mr-2 size-4" />
                    Delete User
                  </UiDropdownMenuItem>
                </UiDropdownMenuContent>
              </UiDropdownMenu>
            </div>
          </div>
        </div>
      </div>
    </UiCard>
  </div>
</template>
