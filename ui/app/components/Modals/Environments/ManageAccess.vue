<script setup lang="ts">
import { UserPlus, Users, ShieldCheck, Search, Plus, X } from "lucide-vue-next"

interface User {
  id: string
  name: string
  email: string
  avatar?: string
}

interface Team {
  id: string
  name: string
  description: string
  memberCount: number
}

interface Role {
  id: string
  name: string
  description: string
  permissions: string[]
}

interface EnvironmentAccess {
  userId?: string
  teamId?: string
  roleId: string
  type: "user" | "team"
}

interface Props {
  open: boolean
  environment: any
}

const props = defineProps<Props>()
const emit = defineEmits<{
  "update:open": [value: boolean]
  updated: [data: { environmentId: string; access: EnvironmentAccess[] }]
}>()

const isOpen = computed({
  get: () => props.open,
  set: (value) => emit("update:open", value)
})

// Available roles (would come from API)
const availableRoles: Role[] = [
  {
    id: "owner",
    name: "Owner",
    description: "Full administrative and environment access",
    permissions: ["read", "write", "admin", "delete"]
  },
  {
    id: "admin",
    name: "Admin",
    description: "Administrative access with user management",
    permissions: ["read", "write", "admin"]
  },
  {
    id: "developer",
    name: "Developer",
    description: "Read/write access to applications and API keys",
    permissions: ["read", "write"]
  },
  {
    id: "viewer",
    name: "Viewer",
    description: "Read-only access to environment resources",
    permissions: ["read"]
  }
]

// Available users (would come from API)
const availableUsers: User[] = [
  { id: "1", name: "Alice Johnson", email: "alice@company.com" },
  { id: "2", name: "Bob Smith", email: "bob@company.com" },
  { id: "3", name: "Carol Williams", email: "carol@company.com" },
  { id: "4", name: "David Brown", email: "david@company.com" },
  { id: "5", name: "Emma Davis", email: "emma@company.com" },
  { id: "6", name: "Frank Wilson", email: "frank@company.com" },
  { id: "7", name: "Grace Miller", email: "grace@company.com" },
  { id: "8", name: "Henry Garcia", email: "henry@company.com" }
]

// Available teams (would come from API)
const availableTeams: Team[] = [
  { id: "1", name: "Engineering", description: "Software development team", memberCount: 12 },
  { id: "2", name: "DevOps", description: "Infrastructure and deployment team", memberCount: 5 },
  { id: "3", name: "Product", description: "Product management team", memberCount: 8 },
  { id: "4", name: "QA", description: "Quality assurance team", memberCount: 6 },
  { id: "5", name: "Customer Success", description: "Customer support team", memberCount: 10 },
  { id: "6", name: "Sales", description: "Sales team", memberCount: 15 },
  { id: "7", name: "Marketing", description: "Marketing team", memberCount: 7 },
  { id: "8", name: "Analytics", description: "Data analytics team", memberCount: 4 }
]

// Current access assignments (would come from API based on environment)
const currentAccess = ref<EnvironmentAccess[]>([
  { userId: "1", roleId: "owner", type: "user" },
  { userId: "2", roleId: "admin", type: "user" },
  { teamId: "1", roleId: "developer", type: "team" },
  { teamId: "2", roleId: "admin", type: "team" }
])

// Search and filter state
const searchQuery = ref("")
const selectedRole = ref("")
const isAdding = ref(false)

// Computed filtered lists
const filteredUsers = computed(() => {
  let users = availableUsers.filter((user) => !currentAccess.value.some((access) => access.userId === user.id))

  if (searchQuery.value) {
    const query = searchQuery.value.toLowerCase()
    users = users.filter((user) => user.name.toLowerCase().includes(query) || user.email.toLowerCase().includes(query))
  }

  return users
})

const filteredTeams = computed(() => {
  let teams = availableTeams.filter((team) => !currentAccess.value.some((access) => access.teamId === team.id))

  if (searchQuery.value) {
    const query = searchQuery.value.toLowerCase()
    teams = teams.filter(
      (team) => team.name.toLowerCase().includes(query) || team.description.toLowerCase().includes(query)
    )
  }

  return teams
})

// Helper functions
const getUserName = (userId: string) => {
  return availableUsers.find((u) => u.id === userId)?.name || "Unknown User"
}

const getTeamName = (teamId: string) => {
  return availableTeams.find((t) => t.id === teamId)?.name || "Unknown Team"
}

const getRoleName = (roleId: string) => {
  return availableRoles.find((r) => r.id === roleId)?.name || "Unknown Role"
}

const getInitials = (name: string) => {
  return name
    .split(" ")
    .map((n) => n[0])
    .join("")
    .toUpperCase()
}

// Actions
const addUserAccess = (userId: string) => {
  if (!selectedRole.value) return

  currentAccess.value.push({
    userId,
    roleId: selectedRole.value,
    type: "user"
  })
}

const addTeamAccess = (teamId: string) => {
  if (!selectedRole.value) return

  currentAccess.value.push({
    teamId,
    roleId: selectedRole.value,
    type: "team"
  })
}

const removeAccess = (index: number) => {
  currentAccess.value.splice(index, 1)
}

const updateAccess = (index: number, newRoleId: string) => {
  currentAccess.value[index].roleId = newRoleId
}

const handleSave = () => {
  emit("updated", {
    environmentId: props.environment?.id,
    access: currentAccess.value
  })
  isOpen.value = false
}

const resetForm = () => {
  searchQuery.value = ""
  selectedRole.value = ""
  isAdding.value = false
}

// Watch for modal close to reset form
watch(isOpen, (newValue) => {
  if (!newValue) {
    resetForm()
  }
})
</script>

<template>
  <UiDialog v-model:open="isOpen">
    <UiDialogContent class="min-w-6xl h-[95dvh] p-0 flex flex-col overflow-hidden">
      <!-- Header with gradient -->
      <div class="relative overflow-hidden">
        <div class="absolute inset-0 bg-gradient-to-r from-primary/5 via-primary/10 to-primary/5"></div>
        <div class="relative p-6 border-b">
          <div class="flex items-center justify-between">
            <div>
              <UiDialogTitle class="flex items-center gap-3 text-2xl font-bold">
                <div class="p-2 rounded-lg bg-primary/10">
                  <ShieldCheck class="h-6 w-6 text-primary" />
                </div>
                Environment Access Control
              </UiDialogTitle>
              <UiDialogDescription class="mt-2 text-base">
                Manage who can access
                <span class="font-semibold text-foreground">{{ environment?.name }}</span> environment
              </UiDialogDescription>
            </div>
            <div class="flex items-center gap-2">
              <UiBadge variant="outline" class="px-3 py-1"> {{ currentAccess.length }} Active Assignments </UiBadge>
            </div>
          </div>
        </div>
      </div>

      <!-- Main content area -->
      <div class="flex-1 flex min-h-0">
        <!-- Current Access Panel -->
        <div class="w-2/5 flex flex-col border-r bg-muted/20">
          <div class="px-4 py-2 border-b bg-background">
            <div class="flex items-center justify-between">
              <div>
                <h3 class="font-semibold text-lg">Current Access</h3>
                <p class="text-sm text-muted-foreground mt-1">Active user and team permissions</p>
              </div>
              <div class="flex items-center gap-2">
                <div class="flex items-center gap-1">
                  <div class="w-2 h-2 rounded-full bg-green-500"></div>
                  <span class="text-xs text-muted-foreground"
                    >{{ currentAccess.filter((a) => a.type === "user").length }} Users</span
                  >
                </div>
                <div class="flex items-center gap-1">
                  <div class="w-2 h-2 rounded-full bg-blue-500"></div>
                  <span class="text-xs text-muted-foreground"
                    >{{ currentAccess.filter((a) => a.type === "team").length }} Teams</span
                  >
                </div>
              </div>
            </div>
          </div>

          <div class="flex-1 overflow-hidden">
            <UiScrollArea class="h-full">
              <div class="p-4 space-y-3">
                <div
                  v-for="(access, index) in currentAccess"
                  :key="index"
                  class="group relative overflow-hidden p-4 border rounded-xl bg-background hover:shadow-md transition-all duration-200 hover:scale-[1.02]"
                >
                  <div class="flex items-center justify-between">
                    <div class="flex items-center gap-4">
                      <div class="relative">
                        <UiAvatar class="size-10 ring-2 ring-background shadow-sm">
                          <UiAvatarImage
                            v-if="access.type === 'user'"
                            :src="`https://avatar.vercel.sh/${access.userId}`"
                          />
                          <UiAvatarFallback
                            :class="
                              access.type === 'user' ? 'bg-green-100 text-green-700' : 'bg-blue-100 text-blue-700'
                            "
                          >
                            <component :is="access.type === 'team' ? Users : UserPlus" class="size-5" />
                          </UiAvatarFallback>
                        </UiAvatar>
                        <div
                          :class="[
                            'absolute -bottom-1 -right-1 w-4 h-4 rounded-full border-2 border-background',
                            access.type === 'user' ? 'bg-green-500' : 'bg-blue-500'
                          ]"
                        ></div>
                      </div>
                      <div class="flex-1">
                        <div class="font-semibold text-foreground">
                          {{ access.type === "user" ? getUserName(access.userId!) : getTeamName(access.teamId!) }}
                        </div>
                        <div class="text-sm text-muted-foreground flex items-center gap-2">
                          <span>{{
                            access.type === "user"
                              ? "Individual User"
                              : `Team • ${availableTeams.find((t) => t.id === access.teamId)?.memberCount || 0} members`
                          }}</span>
                        </div>
                      </div>
                    </div>

                    <div class="flex items-center gap-3">
                      <!-- Enhanced Role selector -->
                      <div class="relative">
                        <select
                          :value="access.roleId"
                          @change="updateAccess(index, ($event.target as HTMLSelectElement).value)"
                          class="appearance-none bg-muted/50 border border-border rounded-lg px-3 py-1 pr-6 text-sm font-medium focus:outline-none focus:ring-2 focus:ring-primary focus:border-transparent min-w-[100px]"
                        >
                          <option v-for="role in availableRoles" :key="role.id" :value="role.id">
                            {{ role.name }}
                          </option>
                        </select>
                        <div class="absolute inset-y-0 right-0 flex items-center pr-2 pointer-events-none">
                          <svg
                            class="w-3 h-3 text-muted-foreground"
                            fill="none"
                            stroke="currentColor"
                            viewBox="0 0 24 24"
                          >
                            <path
                              stroke-linecap="round"
                              stroke-linejoin="round"
                              stroke-width="2"
                              d="M19 9l-7 7-7-7"
                            ></path>
                          </svg>
                        </div>
                      </div>

                      <!-- Enhanced Remove button -->
                      <UiButton
                        variant="ghost"
                        size="sm"
                        @click="removeAccess(index)"
                        class="opacity-0 group-hover:opacity-100 transition-opacity text-destructive hover:text-destructive hover:bg-destructive/10"
                      >
                        <X class="size-4" />
                      </UiButton>
                    </div>
                  </div>
                </div>

                <div v-if="currentAccess.length === 0" class="text-center py-12">
                  <div class="mx-auto w-20 h-20 bg-muted rounded-full flex items-center justify-center mb-4">
                    <ShieldCheck class="size-8 text-muted-foreground" />
                  </div>
                  <h4 class="font-semibold text-foreground mb-2">No access assignments</h4>
                  <p class="text-sm text-muted-foreground">
                    Start by adding users or teams from the panel on the right
                  </p>
                </div>
              </div>
            </UiScrollArea>
          </div>
        </div>

        <!-- Add Access Panel -->
        <div class="w-3/5 flex flex-col">
          <div class="py-2 px-4 border-b bg-background">
            <div class="flex items-center justify-between mb-4">
              <div>
                <h3 class="font-semibold text-lg">Grant Access</h3>
                <p class="text-sm text-muted-foreground mt-1">Add new users and teams to this environment</p>
              </div>
            </div>

            <!-- Enhanced Role Selection -->
            <div class="space-y-3">
              <div>
                <label class="text-sm font-semibold text-foreground mb-2 block">Choose Permission Level</label>
                <div class="grid grid-cols-2 gap-2">
                  <button
                    v-for="role in availableRoles"
                    :key="role.id"
                    @click="selectedRole = role.id"
                    :class="[
                      'p-2 text-left border-2 rounded-lg transition-all',
                      selectedRole === role.id
                        ? 'border-primary bg-primary/5 ring-2 ring-primary/20'
                        : 'border-border hover:border-primary/50 hover:bg-muted/50'
                    ]"
                  >
                    <div class="font-medium text-sm">{{ role.name }}</div>
                    <div class="text-xs text-muted-foreground mt-1">{{ role.description }}</div>
                  </button>
                </div>
              </div>

              <!-- Enhanced Search -->
              <div class="relative">
                <Search class="absolute left-3 top-1/2 transform -translate-y-1/2 size-4 text-muted-foreground" />
                <UiInput
                  v-model="searchQuery"
                  placeholder="Search users or teams..."
                  class="pl-10 h-9 bg-muted/50 border-border focus:bg-background"
                />
              </div>
            </div>
          </div>

          <!-- Combined Users and Teams List -->
          <div class="flex-1 overflow-hidden">
            <UiScrollArea class="h-full">
              <div class="p-4 space-y-2">
                <!-- Section: Users -->
                <div v-if="filteredUsers.length > 0" class="space-y-2">
                  <div class="flex items-center gap-2 px-2 py-1">
                    <UserPlus class="size-4 text-green-600" />
                    <h4 class="font-semibold text-sm text-foreground">Users</h4>
                    <UiBadge variant="secondary" size="sm">{{ filteredUsers.length }}</UiBadge>
                  </div>
                  <div
                    v-for="user in filteredUsers"
                    :key="`user-${user.id}`"
                    class="group flex items-center justify-between p-3 border rounded-lg hover:shadow-sm cursor-pointer transition-all duration-200 hover:scale-[1.02]"
                    @click="selectedRole && addUserAccess(user.id)"
                    :class="{
                      'opacity-50 cursor-not-allowed': !selectedRole,
                      'hover:border-primary/50 hover:bg-primary/5': selectedRole
                    }"
                  >
                    <div class="flex items-center gap-3">
                      <UiAvatar class="size-8">
                        <UiAvatarImage :src="`https://avatar.vercel.sh/${user.id}`" />
                        <UiAvatarFallback class="bg-green-100 text-green-700 font-semibold">{{
                          getInitials(user.name)
                        }}</UiAvatarFallback>
                      </UiAvatar>
                      <div>
                        <div class="font-medium text-foreground">{{ user.name }}</div>
                        <div class="text-sm text-muted-foreground">{{ user.email }}</div>
                      </div>
                    </div>
                    <div class="flex items-center gap-2">
                      <div v-if="selectedRole" class="text-xs text-muted-foreground px-2 py-1 bg-muted rounded">
                        as {{ getRoleName(selectedRole) }}
                      </div>
                      <Plus class="size-4 text-muted-foreground group-hover:text-primary transition-colors" />
                    </div>
                  </div>
                </div>

                <!-- Section: Teams -->
                <div v-if="filteredTeams.length > 0" class="space-y-2" :class="{ 'mt-6': filteredUsers.length > 0 }">
                  <div class="flex items-center gap-2 px-2 py-1">
                    <Users class="size-4 text-blue-600" />
                    <h4 class="font-semibold text-sm text-foreground">Teams</h4>
                    <UiBadge variant="secondary" size="sm">{{ filteredTeams.length }}</UiBadge>
                  </div>
                  <div
                    v-for="team in filteredTeams"
                    :key="`team-${team.id}`"
                    class="group flex items-center justify-between p-3 border rounded-lg hover:shadow-sm cursor-pointer transition-all duration-200 hover:scale-[1.02]"
                    @click="selectedRole && addTeamAccess(team.id)"
                    :class="{
                      'opacity-50 cursor-not-allowed': !selectedRole,
                      'hover:border-primary/50 hover:bg-primary/5': selectedRole
                    }"
                  >
                    <div class="flex items-center gap-3">
                      <div class="p-2 rounded-lg bg-blue-100 text-blue-700">
                        <Users class="size-4" />
                      </div>
                      <div>
                        <div class="font-medium text-foreground">{{ team.name }}</div>
                        <div class="text-sm text-muted-foreground">
                          {{ team.description }} • {{ team.memberCount }} members
                        </div>
                      </div>
                    </div>
                    <div class="flex items-center gap-2">
                      <div v-if="selectedRole" class="text-xs text-muted-foreground px-2 py-1 bg-muted rounded">
                        as {{ getRoleName(selectedRole) }}
                      </div>
                      <Plus class="size-4 text-muted-foreground group-hover:text-primary transition-colors" />
                    </div>
                  </div>
                </div>

                <!-- Empty State -->
                <div v-if="filteredUsers.length === 0 && filteredTeams.length === 0" class="text-center py-12">
                  <div class="mx-auto w-16 h-16 bg-muted rounded-full flex items-center justify-center mb-4">
                    <UserPlus class="size-6 text-muted-foreground" />
                  </div>
                  <h4 class="font-semibold text-foreground mb-2">No users or teams found</h4>
                  <p class="text-sm text-muted-foreground">
                    {{
                      searchQuery
                        ? "Try adjusting your search terms"
                        : "All available users and teams are already assigned"
                    }}
                  </p>
                </div>
              </div>
            </UiScrollArea>
          </div>
        </div>
      </div>

      <!-- Footer Action Bar -->
      <div class="border-t bg-background/95 backdrop-blur supports-[backdrop-filter]:bg-background/60">
        <div class="flex items-center justify-between p-6">
          <div class="flex items-center gap-4">
            <div class="flex items-center gap-2 text-sm text-muted-foreground">
              <div class="flex items-center gap-1">
                <div class="w-2 h-2 rounded-full bg-green-500"></div>
                <span>{{ currentAccess.filter((a) => a.type === "user").length }} Users</span>
              </div>
              <div class="flex items-center gap-1">
                <div class="w-2 h-2 rounded-full bg-blue-500"></div>
                <span>{{ currentAccess.filter((a) => a.type === "team").length }} Teams</span>
              </div>
              <span>•</span>
              <span>{{ currentAccess.length }} Total Assignments</span>
            </div>
          </div>

          <div class="flex items-center gap-3">
            <UiButton variant="outline" @click="isOpen = false" class="px-6"> Cancel </UiButton>
            <UiButton @click="handleSave" class="px-6 bg-primary hover:bg-primary/90">
              <ShieldCheck class="size-4 mr-2" />
              Save Changes
            </UiButton>
          </div>
        </div>
      </div>
    </UiDialogContent>
  </UiDialog>
</template>
