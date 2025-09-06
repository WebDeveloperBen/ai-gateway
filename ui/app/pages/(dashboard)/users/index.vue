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
import SearchFilter from "@/components/SearchFilter.vue"
import type { FilterConfig, SearchConfig, DisplayConfig } from "@/components/SearchFilter.vue"
import type { StatsCardProps } from "@/components/Cards/Stats.vue"

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

const modalState = reactive({
  showInvite: false,
  showEdit: false,
  editingUser: null as User | null,
  editLoading: false,
  showDeactivate: false,
  deactivatingUser: null as User | null,
  deactivateLoading: false,
  showDelete: false,
  deletingUser: null as User | null,
  deleteLoading: false
})

const searchFilterRef = useTemplateRef("searchFilterRef")
const route = useRoute()

// Check query parameter to auto-open modal
onMounted(() => {
  // auto open the invite user modal
  if (route.query.create === "user") {
    modalState.showInvite = true
    return
  }
  // allow deep linking to specific search filter records
  if (route.query.userId) {
    const user = users.find((u) => u.id === route.query.userId)
    if (user && searchFilterRef.value?.setFilters) {
      searchFilterRef.value.setFilters({ name: user.name })
    }
  }
})

const onUserInvited = (inviteData: any) => {
  console.log("User invited:", inviteData)
  // Clear the query parameter after invitation
  navigateTo("/users", { replace: true })
  // Could refresh the users list or add pending invite to list
}

const openEditModal = (user: User) => {
  modalState.editingUser = user
  modalState.showEdit = true
}

const handleUserSave = async (updatedUser: User) => {
  modalState.editLoading = true
  try {
    // TODO: Replace with actual API call
    // await $fetch(`/api/users/${updatedUser.id}`, { method: 'PUT', body: updatedUser })

    // Update local state
    const index = users.findIndex((u) => u.id === updatedUser.id)
    if (index !== -1) {
      users[index] = updatedUser
    }

    modalState.editLoading = false
    modalState.editingUser = null

    // TODO: Show success toast
    console.log("User updated:", updatedUser)
  } catch (error) {
    console.error("Failed to update user:", error)
    // TODO: Show error toast
  } finally {
    modalState.editLoading = false
  }
}

const handleEditCancel = () => {
  modalState.showEdit = false
  // Delay clearing the user to avoid flash during modal close animation
  setTimeout(() => {
    modalState.editingUser = null
  }, 150)
}

const openDeactivateModal = (user: User) => {
  modalState.deactivatingUser = user
  modalState.showDeactivate = true
}

const handleUserDeactivate = async () => {
  if (!modalState.deactivatingUser) return
  modalState.deactivateLoading = true
  try {
    const newStatus = modalState.deactivatingUser.status === "active" ? "inactive" : "active"

    // TODO: Replace with actual API call
    // await $fetch(`/api/users/${deactivatingUser.value.id}/status`, {
    //   method: 'PATCH',
    //   body: { status: newStatus }
    // })

    // Update local state
    const index = users.findIndex((u) => u.id === modalState.deactivatingUser!.id)
    if (index !== -1) {
      if (users[index]) users[index].status = newStatus
    }
    modalState.showDeactivate = false
    modalState.deactivatingUser = null

    // TODO: Show success toast
    console.log("User status updated:", newStatus)
  } catch (error) {
    console.error("Failed to update user status:", error)
    // TODO: Show error toast
  } finally {
    modalState.deactivateLoading = false
  }
}

const handleDeactivateCancel = () => {
  modalState.showDeactivate = false
  // Delay clearing the user to avoid flash during modal close animation
  setTimeout(() => {
    modalState.deactivatingUser = null
  }, 150)
}

const openDeleteModal = (user: User) => {
  modalState.deletingUser = user
  modalState.showDelete = true
}

const handleUserDelete = async () => {
  if (!modalState.deletingUser) return
  modalState.deleteLoading = true
  try {
    // TODO: Replace with actual API call
    // await $fetch(`/api/users/${deletingUser.value.id}`, { method: 'DELETE' })

    // Update local state - remove user from array
    const index = users.findIndex((u) => u.id === modalState.deletingUser!.id)
    if (index !== -1) {
      users.splice(index, 1)
    }
    modalState.showDelete = false
    modalState.deletingUser = null

    // TODO: Show success toast
    console.log("User deleted successfully")
  } catch (error) {
    console.error("Failed to delete user:", error)
    // TODO: Show error toast
  } finally {
    modalState.deleteLoading = false
  }
}

const handleDeleteCancel = () => {
  modalState.showDelete = false
  // Delay clearing the user to avoid flash during modal close animation
  setTimeout(() => {
    modalState.deletingUser = null
  }, 150)
}
</script>

<template>
  <div class="flex flex-col gap-6">
    <!-- Header -->
    <PageHeader title="Users" subtext="Manage users, roles and team access">
      <UiButton @click="modalState.showInvite = true" class="gap-2">
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
      ref="searchFilterRef"
      :items="users"
      :filters="filterConfigs"
      :search-config="searchConfig"
      :display-config="displayConfig"
      @filters-changed="handleFiltersChanged"
      @item-selected="handleItemSelected"
    />

    <!-- Users Table -->
    <CardsDataList title="All Users" :icon="Users">
      <template #actions> </template>

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
    </CardsDataList>

    <!-- Invite User Modal -->
    <LazyModalsUsersInvite v-model:open="modalState.showInvite" @invited="onUserInvited" />

    <!-- Edit User Modal -->
    <LazyModalsUsersEdit
      v-model:open="modalState.showEdit"
      :user="modalState.editingUser"
      :loading="modalState.editLoading"
      @save="handleUserSave"
      @cancel="handleEditCancel"
    />

    <!-- Deactivate User Modal -->
    <ConfirmationModal
      v-model:open="modalState.showDeactivate"
      :title="modalState.deactivatingUser?.status === 'active' ? 'Deactivate User' : 'Activate User'"
      :description="
        modalState.deactivatingUser?.status === 'active'
          ? `Are you sure you want to deactivate ${modalState.deactivatingUser?.name}? They will lose access to the system.`
          : `Are you sure you want to activate ${modalState.deactivatingUser?.name}? They will regain access to the system.`
      "
      :confirm-text="modalState.deactivatingUser?.status === 'active' ? 'Deactivate' : 'Activate'"
      :variant="modalState.deactivatingUser?.status === 'active' ? 'destructive' : 'default'"
      :loading="modalState.deactivateLoading"
      @confirm="handleUserDeactivate"
      @cancel="handleDeactivateCancel"
    />

    <!-- Delete User Modal -->
    <ConfirmationModal
      v-model:open="modalState.showDelete"
      title="Delete User"
      :description="`Are you sure you want to delete ${modalState.deletingUser?.name}? This action cannot be undone and will permanently remove their account and all associated data.`"
      confirm-text="Delete User"
      variant="destructive"
      :loading="modalState.deleteLoading"
      @confirm="handleUserDelete"
      @cancel="handleDeleteCancel"
    />
  </div>
</template>
