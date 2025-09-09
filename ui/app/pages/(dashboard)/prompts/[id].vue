<script setup lang="ts">
import {
  FileText,
  Edit,
  Settings,
  Clock,
  Tag,
  User,
  Copy,
  Activity,
  Calendar,
  MoreVertical,
  Trash2,
  Archive
} from "lucide-vue-next"
import type { StatsCardProps } from "@/components/Cards/Stats.vue"

// Get the prompt ID from route params
const route = useRoute()
const promptId = route.params.id as string

// Environment context
const { selectedEnvironment } = useEnvironment()

// Sample prompt data - this would come from an API call
const prompt = ref({
  id: promptId,
  name: "Customer Support Assistant",
  description: "A helpful assistant for handling customer inquiries with empathy and professionalism",
  content: {
    userPrompt: `Help me resolve a customer inquiry with empathy and professionalism.

The customer is experiencing: {{customer_issue}}

Please provide a helpful response that:
- Shows empathy for their situation
- Asks clarifying questions if needed
- Offers clear solutions or next steps
- Maintains a professional yet friendly tone`,
    systemPrompt: `You are a customer support assistant for a technology company. Your role is to help customers with their inquiries in a helpful, empathetic, and professional manner.

Key Guidelines:
- Always greet customers warmly and thank them for contacting support
- Listen carefully to their concerns and ask clarifying questions when needed
- Provide clear, step-by-step solutions when possible
- If you cannot resolve an issue, escalate to the appropriate team
- End conversations by asking if there's anything else you can help with

Tone: Professional yet friendly, patient, and understanding
Language: Clear and accessible, avoiding technical jargon when speaking with non-technical customers`,
    config: {
      model: "gpt-4",
      temperature: 0.7,
      maxTokens: 1000,
      topP: 0.9
    }
  },
  category: "Customer Service",
  status: "approved" as "draft" | "review" | "approved" | "archived",
  deployments: [
    {
      environment: "Production",
      version: "1.2",
      deployedAt: "2025-01-15T10:30:00Z",
      isActive: true
    },
    {
      environment: "Staging",
      version: "1.2",
      deployedAt: "2025-01-16T09:20:00Z",
      isActive: true
    }
  ],
  createdBy: "Alice Johnson",
  createdAt: "2024-12-01T09:00:00Z",
  updatedAt: "2025-01-15T10:30:00Z",
  usageCount: 1247,
  tags: ["customer-service", "support", "professional"],
  version: "1.2",
  versions: [
    {
      id: "v1.2",
      version: "1.2",
      status: "approved" as const,
      deployments: ["Production", "Staging"],
      content: {
        userPrompt: `Help me resolve a customer inquiry with empathy and professionalism.

The customer is experiencing: {{customer_issue}}

Please provide a helpful response that:
- Shows empathy for their situation
- Asks clarifying questions if needed
- Offers clear solutions or next steps
- Maintains a professional yet friendly tone`,
        systemPrompt: `You are a customer support assistant for a technology company. Your role is to help customers with their inquiries in a helpful, empathetic, and professional manner.

Key Guidelines:
- Always greet customers warmly and thank them for contacting support
- Listen carefully to their concerns and ask clarifying questions when needed
- Provide clear, step-by-step solutions when possible
- If you cannot resolve an issue, escalate to the appropriate team
- End conversations by asking if there's anything else you can help with

Tone: Professional yet friendly, patient, and understanding
Language: Clear and accessible, avoiding technical jargon when speaking with non-technical customers`,
        config: {
          model: "gpt-4",
          temperature: 0.7,
          maxTokens: 1000,
          topP: 0.9
        }
      },
      createdAt: "2025-01-15T10:30:00Z",
      createdBy: "Alice Johnson",
      usageCount: 1247
    },
    {
      id: "v1.1",
      version: "1.1",
      status: "review" as const,
      deployments: [],
      content: {
        userPrompt: `Please help resolve this customer inquiry: {{customer_issue}}`,
        systemPrompt: `You are a customer support assistant. Be helpful and professional.`,
        config: {
          model: "gpt-3.5-turbo",
          temperature: 0.5,
          maxTokens: 500
        }
      },
      createdAt: "2025-01-08T09:15:00Z",
      createdBy: "Alice Johnson",
      usageCount: 892
    },
    {
      id: "v1.3",
      version: "1.3",
      status: "approved" as const,
      deployments: [],
      content: {
        userPrompt: `Enhanced customer support prompt with new features and improved clarity.`,
        systemPrompt: `You are an advanced customer support AI with enhanced capabilities for handling complex customer inquiries.`,
        config: {
          model: "gpt-4",
          temperature: 0.8,
          maxTokens: 1200
        }
      },
      createdAt: "2025-01-16T14:20:00Z",
      createdBy: "Alice Johnson",
      usageCount: 0
    },
    {
      id: "draft-1",
      version: "Draft",
      status: "draft" as const,
      deployments: [],
      content: {
        userPrompt: `Work in progress customer support prompt...`,
        systemPrompt: `Draft system prompt...`,
        config: {
          model: "gpt-4",
          temperature: 0.8,
          maxTokens: 1200
        }
      },
      createdAt: "2025-01-17T10:15:00Z",
      createdBy: "Alice Johnson",
      usageCount: 0
    }
  ],
  recentUsage: [
    {
      timestamp: "2025-01-15T14:30:00Z",
      application: "Customer Service Bot",
      usage: 45
    },
    {
      timestamp: "2025-01-14T16:45:00Z",
      application: "Help Desk Assistant",
      usage: 23
    },
    {
      timestamp: "2025-01-13T12:20:00Z",
      application: "Live Chat Support",
      usage: 67
    }
  ],
  recentActivity: [
    {
      timestamp: "2025-01-15T10:30:00Z",
      action: "Prompt Updated",
      details: "Updated guidelines section for clarity"
    },
    {
      timestamp: "2025-01-10T14:20:00Z",
      action: "High Usage",
      details: "Prompt used 150+ times in past 24 hours"
    },
    {
      timestamp: "2025-01-08T09:15:00Z",
      action: "Version Created",
      details: "Created version 1.2 with improved examples"
    }
  ]
})

// Settings modal state
const isSettingsModalOpen = ref(false)
const isEditModalOpen = ref(false)
const isDeleteModalOpen = ref(false)
const isSettingsLoading = ref(false)

// Version selection state
const selectedVersionId = ref<string>(prompt.value.versions[0]?.id || "")
const selectedVersion = computed(() => {
  return prompt.value.versions.find((v) => v.id === selectedVersionId.value) || prompt.value.versions[0]
})

// Current display content
const displayContent = computed(() => {
  return selectedVersion.value?.content || prompt.value.content
})

const displayStatus = computed(() => {
  return selectedVersion.value?.status || "draft"
})

// Environment-scoped data
const currentEnvironmentDeployment = computed(() => {
  if (!selectedEnvironment.value) return null
  return prompt.value.deployments.find(d => d.environment === selectedEnvironment.value?.name)
})

const isActiveInCurrentEnvironment = computed(() => {
  if (!selectedEnvironment.value) return false
  const deployment = currentEnvironmentDeployment.value
  return deployment?.isActive && deployment?.version === selectedVersion.value?.version
})

const environmentScopedVersions = computed(() => {
  if (!selectedEnvironment.value) return prompt.value.versions
  
  // Add deployment status for current environment to each version
  return prompt.value.versions.map(version => ({
    ...version,
    isActiveInEnvironment: currentEnvironmentDeployment.value?.version === version.version,
    deployments: version.deployments?.filter(env => env === selectedEnvironment.value?.name) || []
  }))
})

const openSettingsModal = () => {
  isSettingsModalOpen.value = true
}

const openEditModal = () => {
  isEditModalOpen.value = true
}

const openDeleteModal = () => {
  isDeleteModalOpen.value = true
}

const copyPromptContent = async () => {
  try {
    const content =
      displayContent.value.userPrompt ||
      displayContent.value.systemPrompt ||
      JSON.stringify(displayContent.value.config)
    await navigator.clipboard.writeText(content)
    // TODO: Show success toast
    console.log("Prompt content copied to clipboard")
  } catch (err) {
    console.error("Failed to copy prompt content:", err)
    // TODO: Show error toast
  }
}

// Version sidebar handlers
const handleVersionSelect = (versionId: string) => {
  selectedVersionId.value = versionId
}

const handleVersionDeploy = async (versionId: string) => {
  try {
    // TODO: Replace with actual API call
    console.log("Deploying version:", versionId)
    // TODO: Show success toast
  } catch (error) {
    console.error("Failed to deploy version:", error)
    // TODO: Show error toast
  }
}

const handleVersionClone = async (versionId: string) => {
  try {
    // TODO: Replace with actual API call
    console.log("Cloning version:", versionId)
    // TODO: Show success toast
  } catch (error) {
    console.error("Failed to clone version:", error)
    // TODO: Show error toast
  }
}

const handleVersionArchive = async (versionId: string) => {
  try {
    // TODO: Replace with actual API call
    console.log("Archiving version:", versionId)
    // TODO: Show success toast
  } catch (error) {
    console.error("Failed to archive version:", error)
    // TODO: Show error toast
  }
}

const handleVersionDelete = async (versionId: string) => {
  try {
    // TODO: Replace with actual API call
    console.log("Deleting version:", versionId)
    // TODO: Show success toast
  } catch (error) {
    console.error("Failed to delete version:", error)
    // TODO: Show error toast
  }
}

const handlePromoteVersion = async (versionId: string, environment: string) => {
  try {
    // TODO: Replace with actual API call
    console.log(`Promoting version ${versionId} to ${environment}`)
    
    // Simulate updating the deployments
    const version = prompt.value.versions.find(v => v.id === versionId)
    if (version && !version.deployments?.includes(environment)) {
      version.deployments = [...(version.deployments || []), environment]
    }
    
    // Add to main deployments array
    const existingDeploymentIndex = prompt.value.deployments.findIndex(d => d.environment === environment)
    const newDeployment = {
      environment,
      version: version?.version || '',
      deployedAt: new Date().toISOString(),
      isActive: true
    }
    
    if (existingDeploymentIndex >= 0) {
      prompt.value.deployments[existingDeploymentIndex] = newDeployment
    } else {
      prompt.value.deployments.push(newDeployment)
    }
    
    // TODO: Show success toast
    console.log(`Successfully promoted to ${environment}`)
  } catch (error) {
    console.error("Failed to promote version:", error)
    // TODO: Show error toast
  }
}

const handleOpenPlayground = (versionId?: string) => {
  const version = versionId ? prompt.value.versions.find((v) => v.id === versionId) : selectedVersion.value

  if (version) {
    // TODO: Navigate to playground with pre-filled prompt
    console.log("Opening playground with version:", version.version)
    navigateTo(`/playground?promptId=${prompt.value.id}&versionId=${version.id}`)
  }
}

// Tabbed content handlers
const handleContentUpdate = (content: any) => {
  // TODO: Handle content updates for editing
  console.log("Content updated:", content)
}

const handleContentCopy = (type: string, content: string) => {
  console.log(`${type} content copied:`, content)
  // TODO: Show success toast
}

const handleSettingsSave = async (updatedPrompt: typeof prompt.value) => {
  isSettingsLoading.value = true
  try {
    // TODO: Replace with actual API call
    // await $fetch(`/api/prompts/${promptId}`, { method: 'PUT', body: updatedPrompt })

    // Update local state
    prompt.value = { ...prompt.value, ...updatedPrompt }
    isSettingsModalOpen.value = false

    // TODO: Show success toast
  } catch (error) {
    console.error("Failed to update prompt:", error)
    // TODO: Show error toast
  } finally {
    isSettingsLoading.value = false
  }
}

const handleSettingsCancel = () => {
  isSettingsModalOpen.value = false
}

const handleDelete = async () => {
  try {
    // TODO: Replace with actual API call
    // await $fetch(`/api/prompts/${promptId}`, { method: 'DELETE' })

    // Navigate back to prompts list
    navigateTo("/prompts")

    // TODO: Show success toast
  } catch (error) {
    console.error("Failed to delete prompt:", error)
    // TODO: Show error toast
  }
}

// Icon and variant arrays for index matching with API data
const statsIcons = [FileText, Activity, Clock, Tag]
const statsVariants: StatsCardProps["variant"][] = ["default", "chart-1", "chart-2", "chart-3"]

// Stats configuration - ready for API replacement
const statsCards = computed(() => {
  return [
    {
      title: "Total Usage",
      value: prompt.value.usageCount,
      description: "All-time uses"
    },
    {
      title: "This Month",
      value: prompt.value.recentUsage.reduce((sum, usage) => sum + usage.usage, 0),
      description: "Recent usage"
    },
    {
      title: "Version",
      value: prompt.value.version,
      description: "Current version"
    },
    {
      title: "Status",
      value: prompt.value.status,
      description:
        prompt.value.status === "active"
          ? "Production ready"
          : prompt.value.status === "draft"
            ? "In development"
            : "Archived"
    }
  ].map((stat, index) => ({
    ...stat,
    icon: statsIcons[index],
    variant: statsVariants[index]
  }))
})

const formatDate = (dateString: string) => {
  return new Date(dateString).toLocaleDateString("en-US", {
    month: "short",
    day: "numeric",
    year: "numeric"
  })
}

const getStatusColor = (status: string) => {
  switch (status) {
    case "active":
      return "border-green-200 bg-green-50 text-green-700"
    case "draft":
      return "border-yellow-200 bg-yellow-50 text-yellow-700"
    case "archived":
      return "border-gray-200 bg-gray-50 text-gray-700"
    default:
      return "border-gray-200 bg-gray-50 text-gray-700"
  }
}

const getStatusBadgeClass = (status: string) => {
  switch (status) {
    case "approved":
      return "bg-green-100 text-green-800 dark:bg-green-900/20 dark:text-green-400 border-green-200 dark:border-green-700"
    case "review":
      return "bg-blue-100 text-blue-800 dark:bg-blue-900/20 dark:text-blue-400 border-blue-200 dark:border-blue-700"
    case "draft":
      return "bg-orange-100 text-orange-800 dark:bg-orange-900/20 dark:text-orange-400 border-orange-200 dark:border-orange-700"
    case "archived":
      return "bg-gray-100 text-gray-800 dark:bg-gray-900/20 dark:text-gray-400 border-gray-200 dark:border-gray-700"
    default:
      return "bg-gray-100 text-gray-800 dark:bg-gray-900/20 dark:text-gray-400 border-gray-200 dark:border-gray-700"
  }
}

const getDeploymentEnvironments = computed(() => {
  return selectedVersion.value?.deployments || []
})
</script>

<template>
  <div class="flex flex-col gap-6">
    <!-- Enhanced Header -->
    <UiCard>
      <UiCardContent class="p-6">
        <div class="flex items-start justify-between">
          <!-- Main Info -->
          <div class="flex-1">
            <div class="flex items-center gap-3 mb-2">
              <h1 class="text-3xl font-bold text-primary tracking-tight">{{ prompt.name }}</h1>
              <UiBadge :class="getStatusBadgeClass(displayStatus)" size="md">
                {{ displayStatus }}
              </UiBadge>

              <!-- Version Info -->
              <div class="flex items-center gap-2">
                <UiBadge variant="outline" size="md"> v{{ selectedVersion?.version || prompt.version }} </UiBadge>

                <!-- Environment Context -->
                <div v-if="selectedEnvironment" class="flex items-center gap-2">
                  <UiBadge variant="outline" size="sm" class="text-xs">
                    {{ selectedEnvironment.name }}
                  </UiBadge>
                  
                  <!-- Show if this version is active in current environment -->
                  <div v-if="isActiveInCurrentEnvironment" class="flex items-center gap-1">
                    <div class="w-2 h-2 bg-green-500 rounded-full animate-pulse" title="Active in this environment"></div>
                    <span class="text-xs font-medium text-green-700 dark:text-green-400">ACTIVE</span>
                  </div>

                  <!-- Show if viewing different version than deployed -->
                  <div v-else-if="currentEnvironmentDeployment" class="flex items-center gap-1">
                    <div class="w-2 h-2 bg-orange-500 rounded-full" title="Not active in this environment"></div>
                    <span class="text-xs font-medium text-orange-700 dark:text-orange-400">INACTIVE</span>
                  </div>

                  <!-- Show if no deployment in this environment -->
                  <div v-else class="flex items-center gap-1">
                    <div class="w-2 h-2 bg-gray-500 rounded-full" title="Not deployed to this environment"></div>
                    <span class="text-xs font-medium text-gray-700 dark:text-gray-400">NOT DEPLOYED</span>
                  </div>
                </div>
              </div>
            </div>
            <p class="text-muted-foreground text-lg mb-4">{{ prompt.description }}</p>

            <!-- Key Metrics Row -->
            <div class="flex items-center gap-6 text-sm">
              <div class="flex items-center gap-2">
                <Activity class="h-4 w-4 text-blue-600" />
                <span class="font-medium">{{ prompt.usageCount.toLocaleString() }}</span>
                <span class="text-muted-foreground">total uses</span>
              </div>
              <div class="flex items-center gap-2">
                <Clock class="h-4 w-4 text-green-600" />
                <span class="font-medium">{{ formatDate(prompt.updatedAt) }}</span>
                <span class="text-muted-foreground">last updated</span>
              </div>
              <div class="flex items-center gap-2">
                <User class="h-4 w-4 text-purple-600" />
                <span class="font-medium">{{ prompt.createdBy }}</span>
                <span class="text-muted-foreground">author</span>
              </div>
            </div>
          </div>

          <!-- Action Buttons -->
          <div class="flex items-center gap-2">
            <UiButton variant="outline" class="gap-2" @click="copyPromptContent">
              <Copy class="h-4 w-4" />
              Copy Content
            </UiButton>
            <UiButton variant="default" class="gap-2" @click="openEditModal">
              <Edit class="h-4 w-4" />
              Edit Prompt
            </UiButton>

            <UiDropdownMenu>
              <UiDropdownMenuTrigger as-child>
                <UiButton variant="outline" size="sm">
                  <MoreVertical class="h-4 w-4" />
                </UiButton>
              </UiDropdownMenuTrigger>
              <UiDropdownMenuContent align="end" class="w-48">
                <UiDropdownMenuItem @click="openSettingsModal">
                  <Settings class="mr-2 h-4 w-4" />
                  Settings
                </UiDropdownMenuItem>
                <UiDropdownMenuItem v-if="prompt.status !== 'archived'" @click="prompt.status = 'archived'">
                  <Archive class="mr-2 h-4 w-4" />
                  Archive
                </UiDropdownMenuItem>
                <UiDropdownMenuSeparator />
                <UiDropdownMenuItem class="text-destructive" @click="openDeleteModal">
                  <Trash2 class="mr-2 h-4 w-4" />
                  Delete Prompt
                </UiDropdownMenuItem>
              </UiDropdownMenuContent>
            </UiDropdownMenu>
          </div>
        </div>
      </UiCardContent>
    </UiCard>

    <!-- Main Content Area with Sidebar -->
    <div class="flex gap-6">
      <!-- Version Selection Sidebar -->
      <PromptsVersionSidebar
        :versions="environmentScopedVersions"
        :selected-version-id="selectedVersionId"
        @select-version="handleVersionSelect"
        @deploy-version="handleVersionDeploy"
        @promote-version="handlePromoteVersion"
        @clone-version="handleVersionClone"
        @archive-version="handleVersionArchive"
        @delete-version="handleVersionDelete"
        @open-playground="handleOpenPlayground"
      />

      <!-- Prompt Content with Tabs -->
      <div class="flex-1 min-h-96">
        <PromptsTabbedContent
          :content="displayContent"
          :status="displayStatus"
          :is-editable="false"
          :show-playground-button="true"
          @update="handleContentUpdate"
          @copy="handleContentCopy"
          @open-playground="() => handleOpenPlayground(selectedVersionId)"
        />
      </div>
    </div>

    <!-- Delete Confirmation Modal -->
    <ConfirmationModal
      v-model:open="isDeleteModalOpen"
      title="Delete Prompt"
      :description="`Are you sure you want to delete '${prompt.name}'? This action cannot be undone and will permanently remove the prompt and all its versions.`"
      confirm-text="Delete Prompt"
      variant="destructive"
      @confirm="handleDelete"
    />
  </div>
</template>
