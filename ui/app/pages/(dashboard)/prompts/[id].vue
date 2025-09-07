<script setup lang="ts">
import { FileText, Edit, Settings, Clock, Tag, User, Copy, Activity, Calendar, MoreVertical, Trash2, Archive } from "lucide-vue-next"
import type { StatsCardProps } from "@/components/Cards/Stats.vue"

// Get the prompt ID from route params
const route = useRoute()
const promptId = route.params.id as string

// Sample prompt data - this would come from an API call
const prompt = ref({
  id: promptId,
  name: "Customer Support Assistant",
  description: "A helpful assistant for handling customer inquiries with empathy and professionalism",
  content: `You are a customer support assistant for a technology company. Your role is to help customers with their inquiries in a helpful, empathetic, and professional manner.

Key Guidelines:
- Always greet customers warmly and thank them for contacting support
- Listen carefully to their concerns and ask clarifying questions when needed
- Provide clear, step-by-step solutions when possible
- If you cannot resolve an issue, escalate to the appropriate team
- End conversations by asking if there's anything else you can help with

Tone: Professional yet friendly, patient, and understanding
Language: Clear and accessible, avoiding technical jargon when speaking with non-technical customers

Example Interaction:
Customer: "My software keeps crashing when I try to export files."
Assistant: "I understand how frustrating that must be. Let me help you troubleshoot this issue. Can you tell me which file format you're trying to export and approximately how large the files are?"`,
  category: "Customer Service",
  status: "active" as "active" | "draft" | "archived",
  createdBy: "Alice Johnson",
  createdAt: "2024-12-01T09:00:00Z",
  updatedAt: "2025-01-15T10:30:00Z",
  usageCount: 1247,
  tags: ["customer-service", "support", "professional"],
  version: "1.2",
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
    await navigator.clipboard.writeText(prompt.value.content)
    // TODO: Show success toast
    console.log("Prompt content copied to clipboard")
  } catch (err) {
    console.error("Failed to copy prompt content:", err)
    // TODO: Show error toast
  }
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
      description: prompt.value.status === "active" ? "Production ready" : prompt.value.status === "draft" ? "In development" : "Archived"
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
</script>

<template>
  <div class="flex flex-col gap-6">
    <!-- Header -->
    <PageBadgeHeader :badge-status="prompt.status" :title="prompt.name" :subtext="prompt.description">
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
    </PageBadgeHeader>

    <!-- Prompt Overview Cards -->
    <div class="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-4 gap-6">
      <CardsStats
        v-for="card in statsCards"
        :key="card.title"
        :title="card.title"
        :value="card.value"
        :icon="card.icon"
        :description="card.description"
        :variant="card.variant"
      />
    </div>

    <!-- Prompt Content -->
    <UiCard>
      <UiCardHeader>
        <UiCardTitle class="flex items-center gap-2">
          <FileText class="h-5 w-5" />
          Prompt Content
        </UiCardTitle>
        <UiCardDescription>The full prompt template content</UiCardDescription>
      </UiCardHeader>
      <UiCardContent>
        <div class="relative">
          <pre class="whitespace-pre-wrap text-sm bg-muted p-4 rounded-lg border font-mono leading-relaxed">{{ prompt.content }}</pre>
          <UiButton 
            variant="ghost" 
            size="sm" 
            class="absolute top-2 right-2 gap-2"
            @click="copyPromptContent"
          >
            <Copy class="h-3 w-3" />
            Copy
          </UiButton>
        </div>
      </UiCardContent>
    </UiCard>

    <!-- Prompt Details -->
    <UiCard>
      <UiCardHeader>
        <UiCardTitle class="flex items-center gap-2">
          <Settings class="h-5 w-5" />
          Prompt Details
        </UiCardTitle>
      </UiCardHeader>
      <UiCardContent class="space-y-6">
        <!-- Prompt Info Grid -->
        <div class="grid grid-cols-1 md:grid-cols-3 gap-6">
          <div>
            <label class="text-sm font-medium text-muted-foreground flex items-center gap-2">
              <Calendar class="h-4 w-4" />
              Created
            </label>
            <p class="text-sm font-medium mt-1">{{ formatDate(prompt.createdAt) }}</p>
          </div>

          <div>
            <label class="text-sm font-medium text-muted-foreground flex items-center gap-2">
              <Clock class="h-4 w-4" />
              Last Updated
            </label>
            <p class="text-sm font-medium mt-1">{{ formatDate(prompt.updatedAt) }}</p>
          </div>

          <div>
            <label class="text-sm font-medium text-muted-foreground flex items-center gap-2">
              <User class="h-4 w-4" />
              Created By
            </label>
            <div class="mt-1 flex items-center gap-2">
              <div class="w-6 h-6 bg-primary/10 rounded-full flex items-center justify-center">
                <span class="text-xs font-medium text-primary">{{ 
                  prompt.createdBy.split(" ").map(n => n[0]).join("")
                }}</span>
              </div>
              <span class="text-sm font-medium">{{ prompt.createdBy }}</span>
            </div>
          </div>
        </div>

        <!-- Category and Tags -->
        <div class="grid grid-cols-1 md:grid-cols-2 gap-6">
          <div>
            <label class="text-sm font-medium text-muted-foreground flex items-center gap-2">
              <Tag class="h-4 w-4" />
              Category
            </label>
            <p class="text-sm font-medium mt-1">{{ prompt.category }}</p>
          </div>

          <div>
            <label class="text-sm font-medium text-muted-foreground">Tags</label>
            <div class="mt-1 flex flex-wrap gap-1">
              <UiBadge 
                v-for="tag in prompt.tags" 
                :key="tag" 
                variant="secondary"
                class="text-xs"
              >
                {{ tag }}
              </UiBadge>
            </div>
          </div>
        </div>
      </UiCardContent>
    </UiCard>

    <!-- Recent Usage -->
    <UiCard>
      <UiCardHeader>
        <UiCardTitle class="flex items-center gap-2">
          <Activity class="h-5 w-5" />
          Recent Usage
        </UiCardTitle>
        <UiCardDescription>Applications using this prompt</UiCardDescription>
      </UiCardHeader>
      <UiCardContent>
        <div class="space-y-4">
          <div 
            v-for="usage in prompt.recentUsage"
            :key="usage.timestamp"
            class="flex items-center justify-between p-3 border rounded-lg"
          >
            <div class="flex items-center gap-3">
              <div class="w-8 h-8 bg-primary/10 rounded-lg flex items-center justify-center">
                <FileText class="h-4 w-4 text-primary" />
              </div>
              <div>
                <p class="text-sm font-medium">{{ usage.application }}</p>
                <p class="text-xs text-muted-foreground">{{ formatDate(usage.timestamp) }}</p>
              </div>
            </div>
            <div class="text-right">
              <p class="text-sm font-medium">{{ usage.usage }}</p>
              <p class="text-xs text-muted-foreground">uses</p>
            </div>
          </div>
        </div>
      </UiCardContent>
    </UiCard>

    <!-- Recent Activity -->
    <UiCard>
      <UiCardHeader>
        <UiCardTitle class="flex items-center gap-2">
          <Activity class="h-5 w-5" />
          Recent Activity
        </UiCardTitle>
        <UiCardDescription>Latest changes and events for this prompt</UiCardDescription>
      </UiCardHeader>
      <UiCardContent>
        <RecentActivity :activities="prompt.recentActivity" />
      </UiCardContent>
    </UiCard>

    <!-- Delete Confirmation Modal -->
    <ConfirmationModal
      v-model:open="isDeleteModalOpen"
      title="Delete Prompt"
      :description="`Are you sure you want to delete '${prompt.name}'? This action cannot be undone and will permanently remove the prompt and all its versions.`"
      confirm-text="Delete Prompt"
      variant="destructive"
      @confirm="handleDelete"
    />

    <!-- Edit Modal (commented out until created) -->
    <!-- <LazyModalsPromptsEdit 
      v-model:open="isEditModalOpen" 
      :prompt="prompt"
      @saved="handlePromptSave"
    /> -->

    <!-- Settings Modal (commented out until created) -->
    <!-- <LazyModalsPromptsSettings
      v-model:open="isSettingsModalOpen"
      :prompt="prompt"
      :loading="isSettingsLoading"
      @save="handleSettingsSave"
      @cancel="handleSettingsCancel"
    /> -->
  </div>
</template>