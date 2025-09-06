<script setup lang="ts">
import { Shield, User } from "lucide-vue-next"
import { useForm } from "vee-validate"
import type { FormBuilder } from "~/components/Ui/FormBuilder/FormBuilder.vue"

interface User {
  id: string
  name: string
  email: string
  role: string
  team: string
  status: "active" | "inactive"
}

interface Team {
  id: string
  name: string
}

interface Props {
  open: boolean
}

const availableUsers: User[] = [
  { id: "1", name: "Alice Johnson", email: "alice@company.com", role: "Admin", team: "Engineering", status: "active" },
  { id: "2", name: "Bob Smith", email: "bob@company.com", role: "Developer", team: "Engineering", status: "active" },
  { id: "3", name: "Carol Williams", email: "carol@company.com", role: "Owner", team: "Marketing", status: "active" },
  { id: "4", name: "David Brown", email: "david@company.com", role: "Developer", team: "Product", status: "active" },
  { id: "5", name: "Emma Davis", email: "emma@company.com", role: "Admin", team: "Engineering", status: "active" },
  {
    id: "6",
    name: "Frank Wilson",
    email: "frank@company.com",
    role: "Developer",
    team: "Customer Success",
    status: "active"
  },
  { id: "7", name: "Grace Lee", email: "grace@company.com", role: "Viewer", team: "Analytics", status: "active" },
  { id: "8", name: "Henry Chen", email: "henry@company.com", role: "Viewer", team: "Marketing", status: "inactive" },
  {
    id: "9",
    name: "Isabel Martinez",
    email: "isabel@company.com",
    role: "Developer",
    team: "Product",
    status: "active"
  },
  { id: "10", name: "Jake Thompson", email: "jake@company.com", role: "Viewer", team: "Sales", status: "active" }
]

const availableTeams: Team[] = [
  { id: "1", name: "Engineering" },
  { id: "2", name: "Product" },
  { id: "3", name: "Marketing" },
  { id: "4", name: "Customer Success" },
  { id: "5", name: "Analytics" },
  { id: "6", name: "Sales" }
]

const roleOptions = [
  { value: "Owner", label: "Owner", description: "Full system access with complete administrative privileges" },
  { value: "Admin", label: "Admin", description: "Administrative access with user and team management capabilities" },
  {
    value: "Developer",
    label: "Developer",
    description: "Development access with API usage and limited management features"
  },
  { value: "Viewer", label: "Viewer", description: "Read-only access for monitoring and reporting purposes" }
]

const props = defineProps<Props>()
const emit = defineEmits<{
  "update:open": [value: boolean]
  assigned: [assignmentData: any]
}>()

const isOpen = computed({
  get: () => props.open,
  set: (value) => emit("update:open", value)
})

const userOptions = availableUsers.map((user) => ({
  value: user.id,
  label: `${user.name} (${user.email})`
}))

const teamOptions = availableTeams.map((team) => ({
  value: team.id,
  label: team.name
}))

const formFields: FormBuilder[] = [
  {
    variant: "Select",
    name: "userId",
    label: "Select User",
    placeholder: "Choose a user to assign a role to",
    hint: "Select the user who will receive the new role assignment",
    required: true,
    options: userOptions
  },
  {
    variant: "Select",
    name: "role",
    label: "Role Assignment",
    placeholder: "Choose the role to assign",
    hint: "This determines what permissions the user will have",
    required: true,
    options: roleOptions.map((option) => ({
      value: option.value,
      label: option.label,
      description: option.description
    }))
  },
  {
    variant: "Select",
    name: "teamId",
    label: "Team Assignment",
    placeholder: "Select which team this role applies to",
    hint: "The user will have this role within the selected team",
    required: true,
    options: teamOptions
  }
]

const { handleSubmit, resetForm, values, isSubmitting } = useForm<{
  userId: string
  role: string
  teamId: string
  name: string
}>({
  initialValues: {
    name: "",
    userId: "",
    role: "",
    teamId: ""
  }
})

const selectedUser = computed(() => {
  if (!values.userId) return null
  return availableUsers.find((user) => user.id === values.userId)
})

const selectedRole = computed(() => {
  if (!values.role) return null
  return roleOptions.find((role) => role.value === values.role)
})

const selectedTeam = computed(() => {
  if (!values.teamId) return null
  return availableTeams.find((team) => team.id === values.teamId)
})

const userCurrentRoleInTeam = computed(() => {
  if (!selectedUser.value || !selectedTeam.value) return null

  // Check if user already has a role in the selected team
  // For now, we'll assume user.team represents their primary team
  // In a real app, you'd have a user-team-role relationship table
  if (selectedUser.value.team === selectedTeam.value.name) {
    return {
      role: selectedUser.value.role,
      team: selectedUser.value.team
    }
  }
  return null
})

const isLoading = shallowRef(false)
const onSubmit = handleSubmit(async (formData) => {
  isLoading.value = true
  try {
    const user = availableUsers.find((u) => u.id === formData.userId)
    const team = availableTeams.find((t) => t.id === formData.teamId)

    if (!user || !team) {
      throw new Error("Invalid user or team selection")
    }

    const assignmentData = {
      id: `assignment_${Date.now()}`,
      userId: user.id,
      userName: user.name,
      userEmail: user.email,
      role: formData.role,
      team: team.name,
      status: "active" as const,
      assignedDate: new Date().toISOString().split("T")[0],
      lastActive: "Just assigned"
    }

    // TODO: Replace with actual API call
    // await $fetch('/api/role-assignments', { method: 'POST', body: assignmentData })

    console.log("Assigning role:", assignmentData)

    emit("assigned", assignmentData)
    handleClose()
  } catch (error) {
    console.error("Failed to assign role:", error)
    // TODO: Show error toast
  } finally {
    isLoading.value = false
  }
})

const handleFormSubmit = () => {
  onSubmit()
}

function handleClose() {
  isOpen.value = false
  // Reset form after modal close animation
  setTimeout(() => {
    resetForm()
  }, 150)
}
</script>

<template>
  <UiDialog v-model:open="isOpen">
    <UiDialogContent class="sm:max-w-lg min-w-2xl">
      <template #header>
        <UiDialogTitle class="sr-only">Assign Role</UiDialogTitle>
        <div class="flex items-center gap-3 pb-4">
          <div
            class="w-10 h-10 rounded-full bg-gradient-to-br from-primary to-primary/80 flex items-center justify-center shadow-sm"
          >
            <Shield class="w-5 h-5 text-primary-foreground" />
          </div>
          <div>
            <h2 class="text-lg font-semibold">Assign Role</h2>
            <p class="text-sm text-muted-foreground">Assign a role to a user for a specific team</p>
          </div>
        </div>
      </template>

      <template #content>
        <form @submit.prevent="handleFormSubmit" class="space-y-6">
          <!-- Form Fields -->
          <UiFormBuilder :fields="formFields" />

          <!-- Assignment Preview -->
          <div
            v-if="selectedUser && selectedRole && selectedTeam"
            class="relative overflow-hidden rounded-xl border bg-gradient-to-br from-background to-muted/50"
          >
            <div class="absolute inset-0 bg-grid-pattern opacity-[0.03]"></div>
            <div class="relative p-6">
              <div class="flex items-center gap-2 mb-5">
                <div class="w-2 h-2 rounded-full bg-primary"></div>
                <h3 class="text-sm font-semibold text-foreground">Assignment Preview</h3>
              </div>

              <div class="space-y-5">
                <!-- User Info Card -->
                <div class="flex items-center gap-4 p-3 rounded-lg bg-card/60 border border-border/50">
                  <UiAvatar class="size-11 ring-2 ring-primary/20">
                    <UiAvatarFallback
                      class="text-sm font-semibold bg-gradient-to-br from-primary to-primary/80 text-primary-foreground"
                    >
                      {{
                        selectedUser.name
                          .split(" ")
                          .map((n) => n[0])
                          .join("")
                      }}
                    </UiAvatarFallback>
                  </UiAvatar>
                  <div class="flex-1">
                    <div class="font-semibold text-sm text-foreground">{{ selectedUser.name }}</div>
                    <div class="text-xs text-muted-foreground">{{ selectedUser.email }}</div>
                  </div>
                  <User class="w-4 h-4 text-muted-foreground" />
                </div>

                <!-- Role Transition -->
                <div class="grid grid-cols-5 gap-3 items-center">
                  <!-- Current State -->
                  <div
                    class="col-span-2 p-3 rounded-lg"
                    :class="
                      userCurrentRoleInTeam
                        ? 'bg-muted/60 border border-border/50'
                        : 'bg-destructive/10 border border-destructive/30'
                    "
                  >
                    <div
                      class="text-xs font-medium mb-2"
                      :class="userCurrentRoleInTeam ? 'text-muted-foreground' : 'text-destructive'"
                    >
                      {{ userCurrentRoleInTeam ? "Current Role" : "Current Status" }}
                    </div>

                    <!-- User has existing role in this team -->
                    <template v-if="userCurrentRoleInTeam">
                      <div class="flex items-center gap-2 mb-1">
                        <component
                          :is="getRoleIcon(userCurrentRoleInTeam.role)"
                          :class="getRoleColor(userCurrentRoleInTeam.role)"
                          class="size-4"
                        />
                        <span class="font-medium text-sm text-foreground">
                          {{ userCurrentRoleInTeam.role }}
                        </span>
                      </div>
                      <div class="text-xs text-muted-foreground">{{ userCurrentRoleInTeam.team }}</div>
                    </template>

                    <!-- User has no role in this team -->
                    <template v-else>
                      <div class="flex items-center gap-2 mb-1">
                        <svg class="w-4 h-4 text-destructive" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                          <path
                            stroke-linecap="round"
                            stroke-linejoin="round"
                            stroke-width="2"
                            d="M18.364 5.636l-3.536 3.536m0 5.656l3.536 3.536M9.172 9.172L5.636 5.636m3.536 9.192L5.636 18.364M12 12h0"
                          />
                        </svg>
                        <span class="font-medium text-sm text-destructive"> No access </span>
                      </div>
                      <div class="text-xs text-destructive/80">to {{ selectedTeam.name }}</div>
                    </template>
                  </div>

                  <!-- Arrow -->
                  <div class="flex justify-center">
                    <div
                      class="w-8 h-8 rounded-full bg-emerald-100 dark:bg-emerald-950/30 flex items-center justify-center"
                    >
                      <svg
                        class="w-4 h-4 text-emerald-600 dark:text-emerald-400"
                        fill="none"
                        stroke="currentColor"
                        viewBox="0 0 24 24"
                      >
                        <path
                          stroke-linecap="round"
                          stroke-linejoin="round"
                          stroke-width="2"
                          d="M13 7l5 5m0 0l-5 5m5-5H6"
                        />
                      </svg>
                    </div>
                  </div>

                  <!-- New Assignment -->
                  <div
                    class="col-span-2 p-3 rounded-lg bg-gradient-to-br from-emerald-50 to-emerald-100/50 dark:from-emerald-950/20 dark:to-emerald-950/10 border border-emerald-200 dark:border-emerald-800/50 ring-1 ring-emerald-500/10"
                  >
                    <div class="text-xs font-medium text-emerald-600 dark:text-emerald-400 mb-2">New Assignment</div>
                    <div class="flex items-center gap-2 mb-1">
                      <component
                        :is="getRoleIcon(selectedRole.value)"
                        :class="getRoleColor(selectedRole.value)"
                        class="size-4"
                      />
                      <span class="font-medium text-sm text-foreground">{{ selectedRole.label }}</span>
                    </div>
                    <div class="text-xs text-muted-foreground">{{ selectedTeam.name }}</div>
                  </div>
                </div>

                <!-- Role Description -->
                <div class="p-4 rounded-lg bg-primary/5 border border-primary/20">
                  <div class="flex items-center gap-2 mb-2">
                    <Shield class="w-3.5 h-3.5 text-primary" />
                    <div class="text-xs font-medium text-foreground">Role Permissions</div>
                  </div>
                  <p class="text-xs text-muted-foreground leading-relaxed">{{ selectedRole.description }}</p>
                </div>
              </div>
            </div>
          </div>

          <div class="flex justify-end space-x-2">
            <UiButton type="button" variant="outline" @click="handleClose" :disabled="isSubmitting"> Cancel </UiButton>
            <UiButton type="button" @click="handleFormSubmit" :loading="isSubmitting" class="gap-2">
              <Shield class="h-4 w-4" />
              Assign Role
            </UiButton>
          </div>
        </form>
      </template>
    </UiDialogContent>
  </UiDialog>
</template>
