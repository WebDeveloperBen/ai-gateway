<script setup lang="ts">
import { Edit, Save, User, Mail, Crown, Shield, Users, Building, UserCheck, UserX } from "lucide-vue-next"

interface UserData {
  id: string
  name: string
  email: string
  role: string
  status: "active" | "inactive"
  team: string
  lastActive: string
  avatar?: string
}

interface Props {
  open: boolean
  user: UserData | null
  loading?: boolean
}

const props = withDefaults(defineProps<Props>(), {
  loading: false
})

const emit = defineEmits<{
  "update:open": [value: boolean]
  save: [user: UserData]
  cancel: []
}>()

const roles = ["Admin", "Developer", "Viewer"]
const teams = ["Engineering", "Product", "Marketing", "Sales", "Design", "Operations"]

const form = ref({
  name: "",
  email: "",
  role: "",
  status: "active" as "active" | "inactive",
  team: ""
})

watch(() => props.user, (newUser) => {
  if (newUser) {
    form.value = {
      name: newUser.name,
      email: newUser.email,
      role: newUser.role,
      status: newUser.status,
      team: newUser.team
    }
  }
}, { immediate: true })

const handleSave = () => {
  if (!props.user) return
  
  emit("save", {
    ...props.user,
    ...form.value
  })
}

const handleCancel = () => {
  emit("cancel")
  emit("update:open", false)
}

const handleOpenChange = (open: boolean) => {
  emit("update:open", open)
  if (!open) {
    emit("cancel")
  }
}

const isFormValid = computed(() => {
  return form.value.name && form.value.email && form.value.role && form.value.team
})
</script>

<template>
  <UiDialog :open="props.open" @update:open="handleOpenChange">
    <UiDialogContent class="sm:max-w-2xl">
      <UiDialogHeader>
        <UiDialogTitle class="flex items-center gap-2">
          <Edit class="h-5 w-5" />
          Edit User
        </UiDialogTitle>
        <UiDialogDescription>
          Update user information, role assignments, and team membership.
        </UiDialogDescription>
      </UiDialogHeader>

      <div class="space-y-6 py-4">
        <div class="grid grid-cols-2 gap-4">
          <div class="space-y-2">
            <label class="text-sm font-medium flex items-center gap-2">
              <User class="h-4 w-4" />
              Full Name
            </label>
            <UiInput
              v-model="form.name"
              placeholder="Enter full name"
              :disabled="props.loading"
            />
          </div>

          <div class="space-y-2">
            <label class="text-sm font-medium flex items-center gap-2">
              <Mail class="h-4 w-4" />
              Email Address
            </label>
            <UiInput
              v-model="form.email"
              type="email"
              placeholder="Enter email address"
              :disabled="props.loading"
            />
          </div>
        </div>

        <div class="grid grid-cols-2 gap-4">
          <div class="space-y-2">
            <label class="text-sm font-medium flex items-center gap-2">
              <Crown class="h-4 w-4" />
              Role
            </label>
            <select
              v-model="form.role"
              :disabled="props.loading"
              class="flex h-10 w-full rounded-md border border-input bg-background px-3 py-2 text-sm ring-offset-background file:border-0 file:bg-transparent file:text-sm file:font-medium placeholder:text-muted-foreground focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-ring focus-visible:ring-offset-2 disabled:cursor-not-allowed disabled:opacity-50"
            >
              <option value="">Select role</option>
              <option v-for="role in roles" :key="role" :value="role">{{ role }}</option>
            </select>
          </div>

          <div class="space-y-2">
            <label class="text-sm font-medium flex items-center gap-2">
              <Building class="h-4 w-4" />
              Team
            </label>
            <select
              v-model="form.team"
              :disabled="props.loading"
              class="flex h-10 w-full rounded-md border border-input bg-background px-3 py-2 text-sm ring-offset-background file:border-0 file:bg-transparent file:text-sm file:font-medium placeholder:text-muted-foreground focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-ring focus-visible:ring-offset-2 disabled:cursor-not-allowed disabled:opacity-50"
            >
              <option value="">Select team</option>
              <option v-for="team in teams" :key="team" :value="team">{{ team }}</option>
            </select>
          </div>
        </div>

        <div class="space-y-2">
          <label class="text-sm font-medium flex items-center gap-2">
            <component :is="form.status === 'active' ? UserCheck : UserX" class="h-4 w-4" />
            Status
          </label>
          <select
            v-model="form.status"
            :disabled="props.loading"
            class="flex h-10 w-full rounded-md border border-input bg-background px-3 py-2 text-sm ring-offset-background file:border-0 file:bg-transparent file:text-sm file:font-medium placeholder:text-muted-foreground focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-ring focus-visible:ring-offset-2 disabled:cursor-not-allowed disabled:opacity-50"
          >
            <option value="active">Active</option>
            <option value="inactive">Inactive</option>
          </select>
        </div>
      </div>

      <UiDialogFooter>
        <UiButton variant="outline" @click="handleCancel" :disabled="props.loading">
          Cancel
        </UiButton>
        <UiButton @click="handleSave" :disabled="props.loading || !isFormValid">
          <span v-if="!props.loading" class="flex items-center gap-2">
            <Save class="h-4 w-4" />
            Save Changes
          </span>
          <span v-else class="flex items-center gap-2">
            <div class="animate-spin rounded-full h-4 w-4 border-b-2 border-current"></div>
            Saving...
          </span>
        </UiButton>
      </UiDialogFooter>
    </UiDialogContent>
  </UiDialog>
</template>