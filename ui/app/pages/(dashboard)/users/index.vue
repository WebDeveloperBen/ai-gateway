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
  User
} from "lucide-vue-next"
import type { FunctionalComponent } from "vue"
import SearchFilter from "~/components/SearchFilter.vue"
import type { FilterConfig, SearchConfig, DisplayConfig } from "~/components/SearchFilter.vue"

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

// Filter state
const activeFilters = ref<Record<string, string>>({})

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

// SearchFilter configuration
const filterConfigs: FilterConfig[] = [
  {
    key: "name",
    label: "User",
    options: users.map((u) => ({ value: u.name, label: u.name, icon: User }))
  },
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
    key: "team",
    label: "Team",
    options: teams.map((team) => ({
      value: team,
      label: team,
      icon: Users
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

const searchConfig: SearchConfig<User> = {
  fields: ["name", "email"],
  placeholder: "Search users, filter by role, team or status..."
}

const displayConfig: DisplayConfig<User> = {
  getItemText: (user) => `${user.name} - ${user.email}`,
  getItemValue: (user) => user.name,
  getItemIcon: () => User
}

// Event handlers
function handleFiltersChanged(filters: Record<string, string>) {
  activeFilters.value = filters
}

function handleItemSelected(user: User) {
  // Handle user selection if needed
  console.log("Selected user:", user)
}

const filteredUsers = computed(() => {
  let filtered = users

  if (activeFilters.value.name && activeFilters.value.name !== "all") {
    filtered = filtered.filter((user) => user.name === activeFilters.value.name)
  }
  if (activeFilters.value.role && activeFilters.value.role !== "all") {
    filtered = filtered.filter((user) => user.role === activeFilters.value.role)
  }
  if (activeFilters.value.team && activeFilters.value.team !== "all") {
    filtered = filtered.filter((user) => user.team === activeFilters.value.team)
  }
  if (activeFilters.value.status && activeFilters.value.status !== "all") {
    filtered = filtered.filter((user) => user.status === activeFilters.value.status)
  }
  return filtered
})

// Modal state
const showInviteModal = ref(false)

// Check query parameter to auto-open modal
const route = useRoute()
onMounted(() => {
  if (route.query.create === 'user') {
    showInviteModal.value = true
  }
})

const onUserInvited = (inviteData: any) => {
  console.log("User invited:", inviteData)
  // Clear the query parameter after invitation
  navigateTo("/users", { replace: true })
  // Could refresh the users list or add pending invite to list
}
</script>

<template>
  <div class="flex flex-col gap-6">
    <!-- Header -->
    <div class="flex items-center justify-between">
      <div>
        <h1 class="text-2xl font-semibold tracking-tight">Users</h1>
        <p class="text-muted-foreground">Manage users, roles and team access</p>
      </div>
      <UiButton @click="showInviteModal = true">
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

    <!-- Search & Filter Component -->
    <SearchFilter
      :items="users"
      :filters="filterConfigs"
      :search-config="searchConfig"
      :display-config="displayConfig"
      @filters-changed="handleFiltersChanged"
      @item-selected="handleItemSelected"
    />

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

    <!-- Invite User Modal -->
    <ModalsInviteUser v-model:open="showInviteModal" @invited="onUserInvited" />
  </div>
</template>
