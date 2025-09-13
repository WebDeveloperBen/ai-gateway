<script setup lang="ts">
import { Database, Plus, Calendar, Tag, FileText, Bookmark, Settings, GitBranch, Globe, Folder, HardDrive, User } from "lucide-vue-next"
import { toast } from "vue-sonner"
import { useForm } from "vee-validate"
import type { FormBuilder } from "@/components/Ui/FormBuilder/FormBuilder.vue"
import type { DataSourceData } from "@/models/datasource"

const props = defineProps<{
  open: boolean
}>()

const emit = defineEmits<{
  "update:open": [value: boolean]
  created: [dataSource: DataSourceData]
}>()

// Data source type options
const typeOptions = [
  { value: "documentation", label: "Documentation" },
  { value: "knowledge-base", label: "Knowledge Base" },
  { value: "policies", label: "Policies" },
  { value: "api-docs", label: "API Documentation" }
]

// Source platform options
const sourceOptions = [
  { value: "confluence", label: "Confluence" },
  { value: "notion", label: "Notion" },
  { value: "sharepoint", label: "SharePoint" },
  { value: "github", label: "GitHub" },
  { value: "google-drive", label: "Google Drive" }
]

// Common schedule options
const scheduleOptions = [
  { value: "0 1 * * *", label: "Daily at 1:00 AM" },
  { value: "0 2 * * *", label: "Daily at 2:00 AM" },
  { value: "0 3 * * *", label: "Daily at 3:00 AM" },
  { value: "0 3 * * 1", label: "Weekly on Monday at 3:00 AM" },
  { value: "0 4 * * 0", label: "Weekly on Sunday at 4:00 AM" },
  { value: "0 2 1 * *", label: "Monthly on 1st at 2:00 AM" }
]

// Sample users for owner selection
const availableUsers = [
  { value: "Alice Johnson", label: "Alice Johnson" },
  { value: "Bob Smith", label: "Bob Smith" },
  { value: "Carol Williams", label: "Carol Williams" },
  { value: "David Brown", label: "David Brown" },
  { value: "Emma Davis", label: "Emma Davis" }
]

// FormBuilder field definitions
const formFields: FormBuilder[] = [
  {
    variant: "Input",
    name: "name",
    label: "Data Source Name",
    placeholder: "Customer Documentation",
    hint: "A descriptive name for your data source",
    required: true
  },
  {
    variant: "Textarea",
    name: "description",
    label: "Description",
    placeholder: "Product documentation and user guides for customer support",
    hint: "What kind of content does this data source contain?",
    required: true,
    rows: 3
  },
  {
    variant: "Select",
    name: "type",
    label: "Content Type",
    hint: "What type of content will be indexed from this source?",
    options: typeOptions,
    required: true
  },
  {
    variant: "Select",
    name: "source",
    label: "Source Platform",
    hint: "Which platform hosts this content?",
    options: sourceOptions,
    required: true
  },
  {
    variant: "Input",
    name: "url",
    label: "Source URL",
    placeholder: "https://company.atlassian.net/wiki",
    hint: "The URL or endpoint for accessing the content",
    type: "url",
    required: true
  },
  {
    variant: "Select",
    name: "schedule",
    label: "Sync Schedule",
    hint: "How often should this data source be synced?",
    options: scheduleOptions,
    required: true
  },
  {
    variant: "Select",
    name: "owner",
    label: "Data Source Owner",
    hint: "Who is responsible for maintaining this data source?",
    options: availableUsers,
    required: true
  },
  {
    variant: "Input",
    name: "tags",
    label: "Tags",
    placeholder: "documentation, support, customer",
    hint: "Comma-separated tags for organizing this data source"
  }
]

// Form setup using vee-validate
const { handleSubmit, resetForm, values } = useForm<{
  name: string
  description: string
  type: string
  source: string
  url: string
  schedule: string
  owner: string
  tags: string
}>({
  initialValues: {
    name: "",
    description: "",
    type: "",
    source: "",
    url: "",
    schedule: "",
    owner: "",
    tags: ""
  }
})

// Dialog state
const dialogOpen = computed({
  get: () => props.open,
  set: (value) => emit("update:open", value)
})

const isLoading = ref(false)

const onSubmit = handleSubmit(async (formData) => {
  isLoading.value = true
  
  try {
    // Convert tags string to array
    const tagsArray = formData.tags 
      ? formData.tags.split(',').map(tag => tag.trim()).filter(tag => tag.length > 0)
      : []

    // Create data source object
    const dataSource: DataSourceData = {
      id: `ds_${Date.now()}`,
      name: formData.name,
      description: formData.description,
      type: formData.type,
      source: formData.source,
      url: formData.url,
      status: "active",
      schedule: formData.schedule,
      lastSync: new Date().toISOString(),
      nextSync: calculateNextSync(formData.schedule),
      documentsCount: 0,
      owner: formData.owner,
      tags: tagsArray
    }

    console.log("Creating data source:", dataSource)

    // TODO: Implement API call to create data source
    
    // Show success toast
    toast({
      title: "Data source created",
      description: `${formData.name} has been successfully created and will begin syncing soon.`,
      duration: 5000,
      icon: "lucide:check"
    })

    // Emit created event
    emit("created", dataSource)
    
    // Close modal and reset form
    closeDialog(true)
  } catch (error) {
    console.error("Failed to create data source:", error)
    toast({
      title: "Error",
      description: "Failed to create data source. Please try again.",
      duration: 5000,
      icon: "lucide:x",
      variant: "destructive"
    })
  } finally {
    isLoading.value = false
  }
})

// Helper function to calculate next sync time
function calculateNextSync(cronExpression: string): string {
  // Simple calculation - in reality, you'd use a proper cron parser
  const now = new Date()
  const nextSync = new Date(now.getTime() + 24 * 60 * 60 * 1000) // Next day for simplicity
  return nextSync.toISOString()
}

const handleFormSubmit = () => {
  onSubmit()
}

const closeDialog = (save: boolean) => {
  if (!save) {
    toast({
      title: "Creation cancelled",
      description: "Data source creation has been cancelled.",
      duration: 3000,
      icon: "lucide:x"
    })
  }
  dialogOpen.value = false
  resetForm()
}

// Watch for modal close to reset form
watch(dialogOpen, (isOpen) => {
  if (!isOpen) {
    resetForm()
  }
})

// Helper to get type icon
function getTypeIcon(type: string) {
  switch (type) {
    case "documentation": return FileText
    case "knowledge-base": return Bookmark
    case "policies": return Settings
    case "api-docs": return GitBranch
    default: return FileText
  }
}

// Helper to get source icon
function getSourceIcon(source: string) {
  switch (source) {
    case "confluence": return Globe
    case "notion": return Bookmark
    case "sharepoint": return Folder
    case "github": return GitBranch
    case "google-drive": return HardDrive
    default: return Globe
  }
}
</script>

<template>
  <UiDialog v-model:open="dialogOpen">
    <UiDialogContent class="sm:max-w-2xl max-h-[90vh] overflow-y-auto">
      <template #header>
        <div class="flex items-center gap-3 pb-4">
          <div
            class="w-10 h-10 rounded-full bg-gradient-to-br from-primary to-primary/80 flex items-center justify-center shadow-sm"
          >
            <Database class="w-5 h-5 text-primary-foreground" />
          </div>
          <div>
            <h2 class="text-lg font-semibold">Add New Data Source</h2>
            <p class="text-sm text-muted-foreground">
              Configure a new data source for RAG ingestion with automated scheduling
            </p>
          </div>
        </div>
      </template>

      <template #content>
        <form @submit.prevent="handleFormSubmit" class="space-y-8">
          <!-- Form Fields -->
          <div>
            <UiFormBuilder :fields="formFields" />

            <!-- Preview Section -->
            <div v-if="values.name" class="mt-6 p-4 bg-muted/20 rounded-lg border">
              <h3 class="text-sm font-medium mb-3 flex items-center gap-2">
                <Database class="size-4" />
                Data Source Preview
              </h3>

              <div class="space-y-3">
                <div class="flex items-center gap-3">
                  <div class="p-2 rounded-lg bg-primary/10">
                    <component 
                      :is="getTypeIcon(values.type)" 
                      class="size-4 text-primary" 
                    />
                  </div>
                  <div>
                    <div class="font-medium text-sm">{{ values.name }}</div>
                    <div class="text-xs text-muted-foreground">{{ values.description || "Data source for RAG ingestion" }}</div>
                  </div>
                </div>

                <div class="grid grid-cols-2 gap-4 text-xs text-muted-foreground">
                  <div v-if="values.type" class="flex items-center gap-2">
                    <component :is="getTypeIcon(values.type)" class="size-3" />
                    <span>{{ typeOptions.find(t => t.value === values.type)?.label }}</span>
                  </div>
                  <div v-if="values.source" class="flex items-center gap-2">
                    <component :is="getSourceIcon(values.source)" class="size-3" />
                    <span>{{ sourceOptions.find(s => s.value === values.source)?.label }}</span>
                  </div>
                  <div v-if="values.schedule" class="flex items-center gap-2">
                    <Calendar class="size-3" />
                    <span>{{ scheduleOptions.find(s => s.value === values.schedule)?.label }}</span>
                  </div>
                  <div v-if="values.owner" class="flex items-center gap-2">
                    <User class="size-3" />
                    <span>{{ values.owner }}</span>
                  </div>
                </div>

                <div v-if="values.tags" class="flex items-center gap-2">
                  <Tag class="size-3 text-muted-foreground" />
                  <div class="flex flex-wrap gap-1">
                    <UiBadge 
                      v-for="tag in values.tags.split(',').map(t => t.trim()).filter(t => t)" 
                      :key="tag" 
                      variant="secondary" 
                      class="text-xs"
                    >
                      {{ tag }}
                    </UiBadge>
                  </div>
                </div>
              </div>
            </div>

            <!-- Info callout -->
            <div class="mt-6 bg-blue-50 dark:bg-blue-950 border border-blue-200 dark:border-blue-800 rounded-lg p-4">
              <div class="flex gap-3">
                <Database class="w-5 h-5 text-blue-600 dark:text-blue-400 flex-shrink-0 mt-0.5" />
                <div class="text-sm">
                  <p class="font-medium text-blue-900 dark:text-blue-100 mb-1">What happens next?</p>
                  <p class="text-blue-700 dark:text-blue-300">
                    Your data source will be registered and the first sync will begin automatically. You can monitor the progress and manage the sync schedule from the data sources dashboard.
                  </p>
                </div>
              </div>
            </div>
          </div>
        </form>
      </template>

      <template #footer>
        <UiDialogFooter>
          <UiButton variant="outline" type="button" class="mt-2 sm:mt-0" @click="closeDialog(false)">
            Cancel
          </UiButton>
          <UiButton type="submit" class="gap-2" @click="handleFormSubmit" :loading="isLoading">
            <Plus class="w-4 h-4" />
            Create Data Source
          </UiButton>
        </UiDialogFooter>
      </template>
    </UiDialogContent>
  </UiDialog>
</template>