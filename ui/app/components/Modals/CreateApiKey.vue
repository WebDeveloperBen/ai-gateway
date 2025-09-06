<script setup lang="ts">
import { Plus, Settings, Key, Copy, Check } from "lucide-vue-next"
import { toast } from "vue-sonner"
import type { FormBuilder } from "@/components/Ui/FormBuilder/FormBuilder.vue"

const props = defineProps<{
  open: boolean
  appId?: string
}>()

const emit = defineEmits<{
  "update:open": [value: boolean]
  created: [apiKeyData: any]
}>()

// Get app config with fallback
const appConfig = useAppConfig()
const appName = computed(() => appConfig?.app?.name || "LLM Gateway")

// Sample applications data - in real app this would come from API
const applications = [
  { value: "app_1", label: "Customer Service Bot" },
  { value: "app_2", label: "Content Generator" },
  { value: "app_3", label: "Code Assistant" }
]

// Essential form fields for API key creation
const formFields: FormBuilder[] = [
  {
    variant: "Input",
    name: "name",
    label: "API Key Name",
    placeholder: "API Key",
    hint: "A descriptive name for this API key",
    required: true
  },
  {
    variant: "Select",
    name: "applicationId",
    label: "Application",
    hint: "Which application will use this API key?",
    options: applications,
    required: true
  },
  {
    variant: "Textarea",
    name: "description",
    label: "Description",
    placeholder: "This key is used for production customer service requests",
    hint: "Optional description of how this key will be used",
    rows: 2
  },
  {
    variant: "Input",
    name: "expiresIn",
    label: "Expires In (days)",
    placeholder: "365",
    hint: "Leave empty for no expiration",
    type: "number"
  }
]

// Form setup using vee-validate
const { handleSubmit, resetForm } = useForm({
  initialValues: {
    name: "",
    applicationId: "",
    description: "",
    expiresIn: ""
  }
})

// Dialog state
const dialogOpen = computed({
  get: () => props.open,
  set: (value) => emit("update:open", value)
})

// API Key creation state
const createdApiKey = ref<string | null>(null)
const createdApiKeyId = ref<string | null>(null)
const isCreating = ref(false)
const justCopied = ref(false)

const onSubmit = handleSubmit(async (values) => {
  try {
    isCreating.value = true

    // TODO: Implement API call to create API key
    console.log("Creating API key:", values)

    // Simulate API response - generate a full API key
    const generatedKey =
      "sk-proj-" +
      Math.random().toString(36).substring(2, 15) +
      Math.random().toString(36).substring(2, 15) +
      Math.random().toString(36).substring(2, 10)

    const apiKeyData = {
      id: "key_new",
      name: values.name,
      key: generatedKey,
      applicationId: values.applicationId,
      description: values.description,
      expiresIn: values.expiresIn,
      created: new Date().toISOString(),
      status: "active",
      permissions: ["read"], // Default permissions
      rateLimit: "1000/hour" // Default rate limit
    }

    // Store the created key to display it
    createdApiKey.value = generatedKey
    createdApiKeyId.value = apiKeyData.id

    // Show success toast
    toast({
      title: "API key created",
      description: `${values.name} has been successfully created. Make sure to copy it now as you won't be able to see it again.`,
      duration: 8000,
      icon: "lucide:check"
    })

    // Emit created event with API key data
    emit("created", apiKeyData)
  } catch (error) {
    console.error("Failed to create API key:", error)
    toast({
      title: "Error",
      description: "Failed to create API key. Please try again.",
      duration: 5000,
      icon: "lucide:x",
      variant: "destructive"
    })
  } finally {
    isCreating.value = false
  }
})

const closeDialog = (save: boolean) => {
  if (!save && !createdApiKey.value) {
    toast({
      title: "Creation cancelled",
      description: "API key creation has been cancelled.",
      duration: 3000,
      icon: "lucide:x"
    })
  }

  // Navigate to key detail page if a key was created
  if (createdApiKey.value && save && createdApiKeyId.value) {
    // Navigate to the specific key detail page
    const appIdToUse = props.appId || "app_1" // fallback to default app
    navigateTo(`/applications/${appIdToUse}/keys/${createdApiKeyId.value}`)
  }

  dialogOpen.value = false
  resetForm()
  createdApiKey.value = null
  createdApiKeyId.value = null
  justCopied.value = false
}

const handleFormSubmit = () => {
  onSubmit()
}

const copyToClipboard = async (text: string) => {
  try {
    await navigator.clipboard.writeText(text)
    justCopied.value = true

    toast.success("copied", {})

    // Reset the copied state after 2 seconds
    setTimeout(() => {
      justCopied.value = false
    }, 2000)
  } catch (err) {
    console.error("Failed to copy text: ", err)
    toast({
      title: "Copy failed",
      description: "Failed to copy to clipboard. Please copy manually.",
      duration: 3000,
      icon: "lucide:x",
      variant: "destructive"
    })
  }
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
        <UiDialogTitle class="sr-only">{{ createdApiKey ? "API Key Created" : "Create New API Key" }}</UiDialogTitle>
        <div class="flex items-center gap-3 pb-4">
          <div
            class="w-10 h-10 rounded-full bg-gradient-to-br from-emerald-500 to-teal-600 flex items-center justify-center shadow-sm"
          >
            <Key class="w-5 h-5 text-white" />
          </div>
          <div>
            <h2 class="text-lg font-semibold">{{ createdApiKey ? "API Key Created" : "Create New API Key" }}</h2>
            <p class="text-sm text-muted-foreground">
              {{
                createdApiKey
                  ? "Your API key has been generated successfully"
                  : `Generate a new API key for your ${appName} applications`
              }}
            </p>
          </div>
        </div>
      </template>

      <template #content>
        <!-- Show created API key -->
        <div v-if="createdApiKey" class="space-y-6">
          <div class="text-center">
            <div
              class="w-16 h-16 bg-green-100 dark:bg-green-900 rounded-full flex items-center justify-center mx-auto mb-4"
            >
              <Key class="w-8 h-8 text-green-600 dark:text-green-400" />
            </div>
            <h3 class="text-lg font-semibold text-green-900 dark:text-green-100">API Key Created Successfully!</h3>
            <p class="text-sm text-green-700 dark:text-green-300 mt-1">
              Please copy your API key now. You won't be able to see it again.
            </p>
          </div>

          <div class="bg-muted/50 rounded-lg p-4 border-2 border-dashed">
            <div class="flex items-center gap-3">
              <code class="text-sm font-mono flex-1 break-all">{{ createdApiKey }}</code>
              <UiButton
                variant="outline"
                size="sm"
                @click="copyToClipboard(createdApiKey)"
                class="gap-2"
                :disabled="justCopied"
              >
                <Check v-if="justCopied" class="h-4 w-4 text-green-600" />
                <Copy v-else class="h-4 w-4" />
                {{ justCopied ? "Copied!" : "Copy" }}
              </UiButton>
            </div>
          </div>

          <div class="bg-red-50 dark:bg-red-950 border border-red-200 dark:border-red-800 rounded-lg p-4">
            <div class="flex gap-3">
              <Settings class="w-5 h-5 text-red-600 dark:text-red-400 flex-shrink-0 mt-0.5" />
              <div class="text-sm">
                <p class="font-medium text-red-900 dark:text-red-100 mb-1">⚠️ Security Warning</p>
                <p class="text-red-700 dark:text-red-300">
                  This API key will never be shown again after you close this dialog. Make sure to copy and store it in
                  a secure location.
                </p>
              </div>
            </div>
          </div>
        </div>

        <!-- Show form -->
        <form v-else @submit.prevent="handleFormSubmit" class="space-y-8">
          <div>
            <UiFormBuilder :fields="formFields" />

            <!-- Security Warning -->
            <div
              class="mt-6 bg-amber-50 dark:bg-amber-950 border border-amber-200 dark:border-amber-800 rounded-lg p-4"
            >
              <div class="flex gap-3">
                <Settings class="w-5 h-5 text-amber-600 dark:text-amber-400 flex-shrink-0 mt-0.5" />
                <div class="text-sm">
                  <p class="font-medium text-amber-900 dark:text-amber-100 mb-1">Important Security Notice</p>
                  <p class="text-amber-700 dark:text-amber-300">
                    Your API key will be shown only once after creation. Make sure to copy and store it securely before
                    closing this dialog.
                  </p>
                </div>
              </div>
            </div>
          </div>
        </form>
      </template>

      <template #footer>
        <UiDialogFooter>
          <!-- Show different buttons based on state -->
          <template v-if="createdApiKey">
            <UiButton variant="outline" type="button" @click="closeDialog(false)"> Close </UiButton>
            <UiButton type="button" class="gap-2" @click="closeDialog(true)"> Continue to Key Management </UiButton>
          </template>
          <template v-else>
            <UiButton variant="outline" type="button" class="mt-2 sm:mt-0" @click="closeDialog(false)">
              Cancel
            </UiButton>
            <UiButton type="button" class="gap-2" @click="handleFormSubmit" :disabled="isCreating">
              <Plus class="w-4 h-4" />
              {{ isCreating ? "Creating..." : "Create API Key" }}
            </UiButton>
          </template>
        </UiDialogFooter>
      </template>
    </UiDialogContent>
  </UiDialog>
</template>
