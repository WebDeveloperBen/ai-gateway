<script setup lang="ts">
import { Plus, X, Crown, Shield, UserCheck, Eye, Building } from "lucide-vue-next"
import { useForm } from "vee-validate"
import type { FormBuilder } from "@/components/Ui/FormBuilder/FormBuilder.vue"

interface Props {
  open: boolean
}

// Sample users for selection
const availableUsers: User[] = [
  { id: "1", name: "Alice Johnson", email: "alice@company.com", role: "Admin" },
  { id: "2", name: "Bob Smith", email: "bob@company.com", role: "Developer" },
  { id: "3", name: "Carol Williams", email: "carol@company.com", role: "Admin" },
  { id: "4", name: "David Brown", email: "david@company.com", role: "Developer" },
  { id: "5", name: "Emma Davis", email: "emma@company.com", role: "Admin" },
  { id: "6", name: "Frank Wilson", email: "frank@company.com", role: "Developer" },
  { id: "7", name: "Grace Lee", email: "grace@company.com", role: "Viewer" },
  { id: "8", name: "Henry Chen", email: "henry@company.com", role: "Developer" }
]

const props = defineProps<Props>()
const emit = defineEmits<{
  "update:open": [value: boolean]
  created: [teamData: any]
}>()

const isOpen = computed({
  get: () => props.open,
  set: (value) => emit("update:open", value)
})

// User options for owner selection
const ownerOptions = availableUsers.map((user) => ({
  value: user.id,
  label: `${user.name} (${user.email})`
}))

// FormBuilder field definitions
const formFields: FormBuilder[] = [
  {
    variant: "Input",
    name: "name",
    label: "Team Name",
    placeholder: "e.g., Engineering, Marketing, Customer Success",
    hint: "A descriptive name for your team",
    required: true
  },
  {
    variant: "Textarea",
    name: "description",
    label: "Description",
    placeholder: "Brief description of the team's purpose and responsibilities",
    hint: "What does this team do?",
    required: true,
    rows: 2
  },
  {
    variant: "Select",
    name: "owner",
    label: "Team Owner",
    hint: "The person responsible for managing this team",
    options: ownerOptions,
    required: true
  }
]

// Form setup using vee-validate (no validation schema - FormBuilder handles it)
const { handleSubmit, resetForm, setFieldValue } = useForm<{
  name: string
  description: string
  costCenter: string
  owner: string
  members: Array<{ user: User; role: "Admin" | "Developer" | "Viewer" }>
}>({
  initialValues: {
    name: "",
    description: "",
    costCenter: "",
    owner: "",
    members: []
  }
})

// Team member management
const selectedMembers = ref<Array<{ user: User; role: "Admin" | "Developer" | "Viewer" }>>([])
const showUserSelector = ref(false)

const isLoading = shallowRef(false)

function addMember(user: User) {
  if (selectedMembers.value.find((m) => m.user.id === user.id)) return

  const newMember = {
    user,
    role: "Developer" as AvailableRoles
  }
  selectedMembers.value.push(newMember)

  // Update form state
  setFieldValue("members", [...selectedMembers.value])
  showUserSelector.value = false
}

function removeMember(userId: string) {
  selectedMembers.value = selectedMembers.value.filter((m) => m.user.id !== userId)
  // Update form state
  setFieldValue("members", [...selectedMembers.value])
}

function updateMemberRole(userId: string, role: "Admin" | "Developer" | "Viewer") {
  const member = selectedMembers.value.find((m) => m.user.id === userId)
  if (member) {
    member.role = role
    // Update form state
    setFieldValue("members", [...selectedMembers.value])
  }
}

const availableUsersForSelection = computed(() => {
  const selectedUserIds = selectedMembers.value.map((m) => m.user.id)
  return availableUsers.filter((user) => !selectedUserIds.includes(user.id))
})

const memberCountByRole = computed(() => {
  const counts = { Admin: 0, Developer: 0, Viewer: 0 }
  selectedMembers.value.forEach((member) => {
    counts[member.role]++
  })
  return counts
})

// Form submission
const onSubmit = handleSubmit(async (formData) => {
  isLoading.value = true

  try {
    // Validate owner is selected
    const ownerUser = availableUsers.find((u) => u.id === formData.owner)
    if (!ownerUser) {
      throw new Error("Owner must be selected")
    }

    const teamData = {
      ...formData,
      members: [
        { userId: ownerUser.id, role: "Owner" },
        ...selectedMembers.value.map((m) => ({ userId: m.user.id, role: m.role }))
      ],
      totalMembers: selectedMembers.value.length + 1, // +1 for owner
      roleDistribution: {
        owner: 1,
        admin: memberCountByRole.value.Admin,
        developer: memberCountByRole.value.Developer,
        viewer: memberCountByRole.value.Viewer
      }
    }

    // TODO: Replace with actual API call
    // await $fetch('/api/teams', { method: 'POST', body: teamData })

    console.log("Creating team:", teamData)

    emit("created", teamData)
    handleClose()
  } catch (error) {
    console.error("Failed to create team:", error)
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
    selectedMembers.value = []
    showUserSelector.value = false
  }, 150)
}
</script>

<template>
  <UiDialog v-model:open="isOpen">
    <UiDialogContent class="min-w-3xl max-h-[90vh] overflow-y-auto">
      <template #header>
        <UiDialogTitle class="sr-only">Create New Team</UiDialogTitle>
        <div class="flex items-center gap-3 pb-4">
          <div
            class="w-10 h-10 rounded-full bg-gradient-to-br from-blue-500 to-indigo-600 flex items-center justify-center shadow-sm"
          >
            <Building class="w-5 h-5 text-white" />
          </div>
          <div>
            <h2 class="text-lg font-semibold">Create New Team</h2>
            <p class="text-sm text-muted-foreground">Organize users into teams with specific roles and policies</p>
          </div>
        </div>
      </template>

      <template #content>
        <form @submit.prevent="handleFormSubmit" class="space-y-6">
          <!-- Basic Team Information -->
          <UiFormBuilder :fields="formFields" />

          <!-- Team Members -->
          <div class="space-y-4">
            <div class="flex items-center justify-between">
              <div>
                <h3 class="text-sm font-medium">Team Members</h3>
                <p class="text-xs text-muted-foreground">Add users to the team and assign their roles</p>
              </div>
              <UiButton
                type="button"
                variant="outline"
                size="sm"
                @click="showUserSelector = true"
                :disabled="availableUsersForSelection.length === 0"
              >
                <Plus class="mr-2 size-4" />
                Add Member
              </UiButton>
            </div>

            <!-- Member Count Summary -->
            <div v-if="selectedMembers.length > 0" class="flex items-center gap-4 p-3 bg-muted/50 rounded-lg text-sm">
              <div class="flex items-center gap-1">
                <Crown class="size-4 text-orange-600" />
                <span>1 Owner</span>
              </div>
              <div class="flex items-center gap-1">
                <Shield class="size-4 text-blue-600" />
                <span>{{ memberCountByRole.Admin }} Admins</span>
              </div>
              <div class="flex items-center gap-1">
                <UserCheck class="size-4 text-green-600" />
                <span>{{ memberCountByRole.Developer }} Developers</span>
              </div>
              <div class="flex items-center gap-1">
                <Eye class="size-4 text-purple-600" />
                <span>{{ memberCountByRole.Viewer }} Viewers</span>
              </div>
            </div>

            <!-- Selected Members List -->
            <div v-if="selectedMembers.length > 0" class="space-y-3 max-h-48 overflow-y-auto">
              <div
                v-for="member in selectedMembers"
                :key="member.user.id"
                class="group relative bg-gradient-to-r from-background via-background to-muted/20 border border-border/60 rounded-lg p-4 hover:border-border transition-all duration-200 hover:shadow-sm"
              >
                <div class="flex items-center justify-between">
                  <div class="flex items-center gap-3 flex-1 min-w-0">
                    <div class="relative">
                      <UiAvatar class="size-10 ring-2 ring-background shadow-sm">
                        <UiAvatarFallback class="text-xs font-medium bg-gradient-to-br from-primary/10 to-primary/5">
                          {{
                            member.user.name
                              .split(" ")
                              .map((n) => n[0])
                              .join("")
                          }}
                        </UiAvatarFallback>
                      </UiAvatar>
                      <div
                        class="absolute -bottom-1 -right-1 p-1 bg-background rounded-full shadow-sm border border-border/40"
                      >
                        <component :is="getRoleIcon(member.role)" class="size-3" :class="getRoleColor(member.role)" />
                      </div>
                    </div>
                    <div class="flex-1 min-w-0">
                      <div class="font-medium text-sm text-foreground truncate">{{ member.user.name }}</div>
                      <div class="text-xs text-muted-foreground truncate">{{ member.user.email }}</div>
                    </div>
                  </div>

                  <div class="flex items-center gap-3 ml-4">
                    <div class="flex items-center gap-2 bg-muted/50 rounded-md px-2 py-1 min-w-0">
                      <component
                        :is="getRoleIcon(member.role)"
                        class="size-3 flex-shrink-0"
                        :class="getRoleColor(member.role)"
                      />
                      <UiNativeSelect
                        :model-value="member.role"
                        @update:model-value="
                          (role: 'Admin' | 'Developer' | 'Viewer') => updateMemberRole(member.user.id, role)
                        "
                        class="bg-transparent border-0 text-xs font-medium h-8 pr-1 focus-visible:ring-0 focus-visible:ring-offset-0 focus-visible:outline-none"
                      >
                        <option v-for="role in availableRoles" :key="role.value" :value="role.value">
                          {{ role.label }}
                        </option>
                      </UiNativeSelect>
                    </div>

                    <UiButton
                      type="button"
                      variant="ghost"
                      size="sm"
                      class="h-8 w-8 p-0 opacity-60 hover:opacity-100 hover:bg-destructive/10 hover:text-destructive transition-all"
                      @click="removeMember(member.user.id)"
                    >
                      <X class="size-4" />
                    </UiButton>
                  </div>
                </div>
              </div>
            </div>

            <!-- User Selector -->
            <UiDialog v-model:open="showUserSelector">
              <UiDialogContent>
                <UiDialogHeader>
                  <UiDialogTitle>Add Team Members</UiDialogTitle>
                  <UiDialogDescription>
                    Select users to add to the team. You can assign roles after adding them.
                  </UiDialogDescription>
                </UiDialogHeader>

                <div class="space-y-2 max-h-60 overflow-y-auto">
                  <div
                    v-for="user in availableUsersForSelection"
                    :key="user.id"
                    class="flex items-center gap-3 p-3 border rounded-lg hover:bg-muted/50 cursor-pointer"
                    @click="addMember(user)"
                  >
                    <UiAvatar class="size-8">
                      <UiAvatarFallback class="text-xs">
                        {{
                          user.name
                            .split(" ")
                            .map((n) => n[0])
                            .join("")
                        }}
                      </UiAvatarFallback>
                    </UiAvatar>
                    <div class="flex-1">
                      <div class="font-medium text-sm">{{ user.name }}</div>
                      <div class="text-xs text-muted-foreground">{{ user.email }}</div>
                    </div>
                    <div class="flex items-center gap-1 text-xs text-muted-foreground">
                      <component :is="getRoleIcon(user.role)" class="size-3" />
                      Current: {{ user.role }}
                    </div>
                  </div>
                </div>

                <div class="flex justify-end">
                  <UiButton variant="outline" @click="showUserSelector = false"> Close </UiButton>
                </div>
              </UiDialogContent>
            </UiDialog>
          </div>
        </form>
      </template>

      <template #footer>
        <UiDialogFooter>
          <UiButton variant="outline" type="button" @click="handleClose" :disabled="isLoading"> Cancel </UiButton>
          <UiButton type="button" @click="handleFormSubmit" :disabled="isLoading" class="gap-2">
            <Plus class="w-4 h-4" />
            {{ isLoading ? "Creating..." : "Create Team" }}
          </UiButton>
        </UiDialogFooter>
      </template>
    </UiDialogContent>
  </UiDialog>
</template>
