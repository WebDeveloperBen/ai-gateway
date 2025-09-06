<script lang="ts">
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
import type { StatsCardProps } from "~/components/Cards/Stats.vue"

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

const stats: StatsCardProps[] = [
  { title: "Total Users", value: "24", icon: Users, description: "+2 this month", variant: "chart-1" },
  { title: "Active Users", value: "18", icon: UserCheck, description: "+5% from last month", variant: "chart-2" },
  { title: "Admin Users", value: "3", icon: Crown, description: "No change", variant: "chart-3" },
  { title: "Teams", value: "6", icon: Shield, description: "+1 this month", variant: "chart-4" }
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
</script>
<script setup lang="ts">
useSeoMeta({ title: "Users - LLM Gateway" })

// Filter state
const activeFilters = ref<Record<string, string>>({})

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
const showEditModal = ref(false)
const editingUser = ref<User | null>(null)
const editLoading = ref(false)

// Deactivate modal state
const showDeactivateModal = ref(false)
const deactivatingUser = ref<User | null>(null)
const deactivateLoading = ref(false)

// Delete modal state
const showDeleteModal = ref(false)
const deletingUser = ref<User | null>(null)
const deleteLoading = ref(false)

// Check query parameter to auto-open modal
const route = useRoute()
onMounted(() => {
  if (route.query.create === "user") {
    showInviteModal.value = true
  }
})

const onUserInvited = (inviteData: any) => {
  console.log("User invited:", inviteData)
  // Clear the query parameter after invitation
  navigateTo("/users", { replace: true })
  // Could refresh the users list or add pending invite to list
}

const openEditModal = (user: User) => {
  editingUser.value = user
  showEditModal.value = true
}

const handleUserSave = async (updatedUser: User) => {
  editLoading.value = true
  try {
    // TODO: Replace with actual API call
    // await $fetch(`/api/users/${updatedUser.id}`, { method: 'PUT', body: updatedUser })

    // Update local state
    const index = users.findIndex((u) => u.id === updatedUser.id)
    if (index !== -1) {
      users[index] = updatedUser
    }

    showEditModal.value = false
    editingUser.value = null

    // TODO: Show success toast
    console.log("User updated:", updatedUser)
  } catch (error) {
    console.error("Failed to update user:", error)
    // TODO: Show error toast
  } finally {
    editLoading.value = false
  }
}

const handleEditCancel = () => {
  showEditModal.value = false
  // Delay clearing the user to avoid flash during modal close animation
  setTimeout(() => {
    editingUser.value = null
  }, 150)
}

const openDeactivateModal = (user: User) => {
  deactivatingUser.value = user
  showDeactivateModal.value = true
}

const handleUserDeactivate = async () => {
  if (!deactivatingUser.value) return

  deactivateLoading.value = true
  try {
    const newStatus = deactivatingUser.value.status === "active" ? "inactive" : "active"

    // TODO: Replace with actual API call
    // await $fetch(`/api/users/${deactivatingUser.value.id}/status`, {
    //   method: 'PATCH',
    //   body: { status: newStatus }
    // })

    // Update local state
    const index = users.findIndex((u) => u.id === deactivatingUser.value!.id)
    if (index !== -1) {
      if (users[index]) users[index].status = newStatus
    }

    showDeactivateModal.value = false
    deactivatingUser.value = null

    // TODO: Show success toast
    console.log("User status updated:", newStatus)
  } catch (error) {
    console.error("Failed to update user status:", error)
    // TODO: Show error toast
  } finally {
    deactivateLoading.value = false
  }
}

const handleDeactivateCancel = () => {
  showDeactivateModal.value = false
  // Delay clearing the user to avoid flash during modal close animation
  setTimeout(() => {
    deactivatingUser.value = null
  }, 150)
}

const openDeleteModal = (user: User) => {
  deletingUser.value = user
  showDeleteModal.value = true
}

const handleUserDelete = async () => {
  if (!deletingUser.value) return

  deleteLoading.value = true
  try {
    // TODO: Replace with actual API call
    // await $fetch(`/api/users/${deletingUser.value.id}`, { method: 'DELETE' })

    // Update local state - remove user from array
    const index = users.findIndex((u) => u.id === deletingUser.value!.id)
    if (index !== -1) {
      users.splice(index, 1)
    }

    showDeleteModal.value = false
    deletingUser.value = null

    // TODO: Show success toast
    console.log("User deleted successfully")
  } catch (error) {
    console.error("Failed to delete user:", error)
    // TODO: Show error toast
  } finally {
    deleteLoading.value = false
  }
}

const handleDeleteCancel = () => {
  showDeleteModal.value = false
  // Delay clearing the user to avoid flash during modal close animation
  setTimeout(() => {
    deletingUser.value = null
  }, 150)
}
</script>

<template>
  <div class="flex flex-col gap-6">
    <!-- Header -->
    <PageHeader
      title="Users"
      subtext="Manage users, roles and team access"
    >
      <UiButton @click="showInviteModal = true" class="gap-2">
        <Plus class="h-4 w-4" />
        Invite User
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
      :items="users"
      :filters="filterConfigs"
      :search-config="searchConfig"
      :display-config="displayConfig"
      @filters-changed="handleFiltersChanged"
      @item-selected="handleItemSelected"
    />

    <!-- Users Table -->
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
              <UiDropdownMenuItem @click="openEditModal(user)">
                <Edit class="mr-2 size-4" />
                Edit User
              </UiDropdownMenuItem>
              <UiDropdownMenuItem @click="openDeactivateModal(user)">
                <component :is="user.status === 'active' ? UserX : UserCheck" class="mr-2 size-4" />
                {{ user.status === "active" ? "Deactivate" : "Activate" }}
              </UiDropdownMenuItem>
              <UiDropdownMenuSeparator />
              <UiDropdownMenuItem class="text-destructive" @click="openDeleteModal(user)">
                <Trash2 class="mr-2 size-4" />
                Delete User
              </UiDropdownMenuItem>
            </UiDropdownMenuContent>
          </UiDropdownMenu>
        </div>
      </div>
    </div>

    <!-- Invite User Modal -->
    <ModalsInviteUser v-model:open="showInviteModal" @invited="onUserInvited" />

    <!-- Edit User Modal -->
    <UserEditModal
      v-model:open="showEditModal"
      :user="editingUser"
      :loading="editLoading"
      @save="handleUserSave"
      @cancel="handleEditCancel"
    />

    <!-- Deactivate User Modal -->
    <ConfirmationModal
      v-model:open="showDeactivateModal"
      :title="deactivatingUser?.status === 'active' ? 'Deactivate User' : 'Activate User'"
      :description="
        deactivatingUser?.status === 'active'
          ? `Are you sure you want to deactivate ${deactivatingUser?.name}? They will lose access to the system.`
          : `Are you sure you want to activate ${deactivatingUser?.name}? They will regain access to the system.`
      "
      :confirm-text="deactivatingUser?.status === 'active' ? 'Deactivate' : 'Activate'"
      :variant="deactivatingUser?.status === 'active' ? 'destructive' : 'default'"
      :loading="deactivateLoading"
      @confirm="handleUserDeactivate"
      @cancel="handleDeactivateCancel"
    />

    <!-- Delete User Modal -->
    <ConfirmationModal
      v-model:open="showDeleteModal"
      title="Delete User"
      :description="`Are you sure you want to delete ${deletingUser?.name}? This action cannot be undone and will permanently remove their account and all associated data.`"
      confirm-text="Delete User"
      variant="destructive"
      :loading="deleteLoading"
      @confirm="handleUserDelete"
      @cancel="handleDeleteCancel"
    />
  </div>
</template>
