<script setup lang="ts">
import { Plus, Settings, Server } from "lucide-vue-next"
import { toast } from "vue-sonner"
import type { FormBuilder } from "@/components/Ui/FormBuilder/FormBuilder.vue"

const props = defineProps<{
  open: boolean
}>()

const emit = defineEmits<{
  "update:open": [value: boolean]
  created: [applicationId: string]
}>()

// Get app config with fallback
const appConfig = useAppConfig()
const appName = computed(() => appConfig?.app?.name || "LLM Gateway")

const environments = [
  { value: "production", label: "Production" },
  { value: "staging", label: "Staging" },
  { value: "development", label: "Development" },
  { value: "testing", label: "Testing" }
]
const teams = ["Engineering", "Marketing", "Customer Success", "Sales", "Operations", "Data Science"]

// Essential form fields for application creation
const formFields: FormBuilder[] = [
  {
    variant: "Input",
    name: "name",
    label: "Application Name",
    placeholder: "Customer Service Bot",
    required: true
  },
  {
    variant: "Textarea",
    name: "description",
    label: "Description",
    placeholder: "AI-powered customer support assistant for handling common inquiries",
    hint: "Brief description of what this application does",
    rows: 3
  },
  {
    variant: "Select",
    name: "team",
    label: "Team",
    hint: "Which team owns this application?",
    options: teams.map((team) => ({ value: team, label: team })),
    required: true
  },
  {
    variant: "Select",
    name: "environment",
    label: "Environment",
    hint: "What environment will this application run in?",
    options: environments,
    required: true
  },
  {
    variant: "Input",
    name: "owner",
    label: "Owner Email",
    placeholder: "john.doe@company.com",
    hint: "Primary contact for this application",
    type: "email",
    required: true,
    wrapperClass: "mb-6"
  }
]

// Form setup using vee-validate
const { handleSubmit, resetForm } = useForm({
  initialValues: {
    name: "",
    description: "",
    team: "",
    environment: "",
    owner: "",
    project: ""
  }
})

// Dialog state
const dialogOpen = computed({
  get: () => props.open,
  set: (value) => emit("update:open", value)
})

const onSubmit = handleSubmit(async (values) => {
  try {
    // TODO: Implement API call to create application
    console.log("Creating application:", values)

    // Simulate API response
    const applicationId = "new-app-id" // This would come from the API response

    // Show success toast
    toast({
      title: "Application created",
      description: `${values.name} has been successfully created.`,
      duration: 5000,
      icon: "lucide:check"
    })

    // Close modal and reset form
    closeDialog(true)

    // Emit created event with application ID
    emit("created", applicationId)
  } catch (error) {
    console.error("Failed to create application:", error)
    toast({
      title: "Error",
      description: "Failed to create application. Please try again.",
      duration: 5000,
      icon: "lucide:x",
      variant: "destructive"
    })
  }
})

const closeDialog = (save: boolean) => {
  if (!save) {
    toast({
      title: "Creation cancelled",
      description: "Application creation has been cancelled.",
      duration: 3000,
      icon: "lucide:x"
    })
  }
  dialogOpen.value = false
  resetForm()
}

const handleFormSubmit = () => {
  onSubmit()
}

// Watch for modal close to reset form
watch(dialogOpen, (isOpen) => {
  if (!isOpen) {
    resetForm()
  }
})
</script>

<template>
  <UiDialog v-model:open="dialogOpen">
    <UiDialogContent class="sm:max-w-2xl max-h-[90vh] overflow-y-auto">
      <template #header>
        <div class="flex items-center gap-3 pb-4">
          <div
            class="w-10 h-10 rounded-full bg-gradient-to-br from-blue-500 to-purple-600 flex items-center justify-center shadow-sm"
          >
            <Server class="w-5 h-5 text-white" />
          </div>
          <div>
            <h2 class="text-lg font-semibold">Create New Application</h2>
            <p class="text-sm text-muted-foreground">
              Provide the essential information to register your application with the {{ appName }}
            </p>
          </div>
        </div>
      </template>
      <template #content>
        <form @submit="onSubmit" class="space-y-8">
          <div>
            <UiFormBuilder :fields="formFields" />

            <!-- Info callout -->
            <div class="bg-blue-50 dark:bg-blue-950 border border-blue-200 dark:border-blue-800 rounded-lg p-4">
              <div class="flex gap-3">
                <Settings class="w-5 h-5 text-blue-600 dark:text-blue-400 flex-shrink-0 mt-0.5" />
                <div class="text-sm">
                  <p class="font-medium text-blue-900 dark:text-blue-100 mb-1">What's next?</p>
                  <p class="text-blue-700 dark:text-blue-300">
                    After creating your application, configure models and policies, then generate your first API key to start making requests.
                  </p>
                </div>
              </div>
            </div>
          </div>
        </form>
      </template>

      <template #footer>
        <UiDialogFooter>
          <UiButton variant="outline" type="button" class="mt-2 sm:mt-0" @click="closeDialog(false)"> Cancel </UiButton>
          <UiButton type="submit" class="gap-2" @click="handleFormSubmit">
            <Plus class="w-4 h-4" />
            Create Application
          </UiButton>
        </UiDialogFooter>
      </template>
    </UiDialogContent>
  </UiDialog>
</template>
