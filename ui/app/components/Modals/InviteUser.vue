<script setup lang="ts">
import { Mail } from "lucide-vue-next"
import type { FormBuilder } from "../Ui/FormBuilder/FormBuilder.vue"

const props = defineProps<{
  open: boolean
}>()

const emit = defineEmits<{
  "update:open": [value: boolean]
  invited: [userData: any]
}>()

const isOpen = computed({
  get: () => props.open,
  set: (value) => emit("update:open", value)
})

const roleOptions = [
  { value: "admin", label: "Admin", description: "Full system access and user management" },
  { value: "developer", label: "Developer", description: "Manage applications and API keys" },
  { value: "viewer", label: "Viewer", description: "Read-only access to dashboards and analytics" }
]

const formFields: FormBuilder[] = [
  {
    variant: "Input",
    name: "email",
    label: "Email Address",
    type: "email",
    placeholder: "user@example.com",
    required: true,
    rules: "required|email"
  },
  {
    variant: "Input",
    name: "name",
    label: "Full Name",
    type: "text",
    placeholder: "John Smith",
    required: true,
    rules: "required|min:2"
  },
  {
    variant: "Select",
    name: "role",
    label: "Role",
    placeholder: "Select a role",
    required: true,
    rules: "required",
    options: roleOptions.map((option) => ({
      value: option.value,
      label: option.label
    }))
  }
]

// Form handling using VeeValidate
const { handleSubmit, isSubmitting, resetForm } = useForm({
  initialValues: {
    email: "",
    name: "",
    role: "developer"
  }
})

const onSubmit = handleSubmit(async (values) => {
  try {
    // TODO: Replace with actual API call
    await new Promise((resolve) => setTimeout(resolve, 1000))

    const inviteData = {
      id: `invite_${Date.now()}`,
      email: values.email,
      name: values.name,
      role: values.role,
      status: "pending",
      invitedAt: new Date().toISOString()
    }

    emit("invited", inviteData)
    closeDialog()
  } catch (error) {
    console.error("Failed to send invitation:", error)
  }
})

const closeDialog = (resetFormData = true) => {
  if (resetFormData) {
    resetForm()
  }
  isOpen.value = false
}
</script>

<template>
  <UiDialog v-model:open="isOpen">
    <UiDialogContent class="sm:max-w-md">
      <template #header>
        <div class="flex items-center gap-3 pb-4">
          <div
            class="w-10 h-10 rounded-full bg-gradient-to-br from-blue-500 to-purple-600 flex items-center justify-center shadow-sm"
          >
            <Mail class="w-5 h-5 text-white" />
          </div>
          <div>
            <h2 class="text-lg font-semibold">Invite User</h2>
            <p class="text-sm text-muted-foreground">Send an invitation to join your LLM Gateway organization</p>
          </div>
        </div>
      </template>
      <template #content>
        <form @submit="onSubmit" class="space-y-6">
          <UiFormBuilder :fields="formFields" />
          <div class="flex justify-end space-x-2">
            <UiButton type="button" variant="outline" @click="closeDialog(false)" :disabled="isSubmitting">
              Cancel
            </UiButton>
            <UiButton type="submit" :loading="isSubmitting" class="gap-2">
              <Mail class="h-4 w-4" />
              Send Invitation
            </UiButton>
          </div>
        </form>
      </template>
    </UiDialogContent>
  </UiDialog>
</template>
