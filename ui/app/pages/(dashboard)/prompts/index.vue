<script lang="ts">
import { Plus, FileText, Tag, Clock, CheckCircle, XCircle, User } from "lucide-vue-next"
import type { StatsCardProps } from "@/components/Cards/Stats.vue"
import SearchFilter from "@/components/SearchFilter.vue"
import type { FilterConfig, SearchConfig, DisplayConfig } from "@/components/SearchFilter.vue"

interface PromptData {
  id: string
  name: string
  description: string
  content: string
  category: string
  status: "active" | "draft" | "archived"
  createdBy: string
  updatedAt: string
  usageCount: number
  tags: string[]
}

const prompts = ref<PromptData[]>([
  {
    id: "prompt_1",
    name: "Customer Support Assistant",
    description: "A helpful assistant for handling customer inquiries with empathy and professionalism",
    content: "You are a customer support assistant. Always be helpful, empathetic, and professional...",
    category: "Customer Service",
    status: "active" as const,
    createdBy: "Alice Johnson",
    updatedAt: "2025-01-15T10:30:00Z",
    usageCount: 1247,
    tags: ["customer-service", "support", "professional"]
  },
  {
    id: "prompt_2",
    name: "Code Review Helper",
    description: "Technical assistant for code review and suggestions",
    content: "Review the following code for best practices, security issues, and optimization opportunities...",
    category: "Development",
    status: "active" as const,
    createdBy: "Bob Smith",
    updatedAt: "2025-01-14T16:45:00Z",
    usageCount: 892,
    tags: ["development", "code-review", "technical"]
  },
  {
    id: "prompt_3",
    name: "Marketing Copy Generator",
    description: "Creative assistant for marketing and promotional content",
    content: "Create engaging marketing copy that converts. Focus on benefits, not features...",
    category: "Marketing",
    status: "draft" as const,
    createdBy: "Carol Williams",
    updatedAt: "2025-01-13T09:20:00Z",
    usageCount: 156,
    tags: ["marketing", "copywriting", "creative"]
  },
  {
    id: "prompt_4",
    name: "Data Analysis Helper",
    description: "Assistant for data interpretation and insights",
    content: "Analyze the provided data and extract meaningful insights...",
    category: "Analytics",
    status: "archived" as const,
    createdBy: "David Brown",
    updatedAt: "2024-12-20T14:15:00Z",
    usageCount: 45,
    tags: ["data", "analytics", "insights"]
  }
])
</script>
<script setup lang="ts">
useSeoMeta({ title: "Prompts - LLM Gateway" })

// Filter state for SearchFilter component
const activeFilters = ref<Record<string, string>>({})

// Get app config
const appConfig = useAppConfig()

const handlePromptSelect = (prompt: PromptData) => {
  navigateTo(`/prompts/${prompt.id}`)
}

// SearchFilter configuration
const categories = [...new Set(prompts.value.map((prompt) => prompt.category))]

const filterConfigs: FilterConfig[] = [
  {
    key: "name",
    label: "Prompt",
    options: prompts.value.map((prompt) => ({ value: prompt.name, label: prompt.name, icon: FileText }))
  },
  {
    key: "status",
    label: "Status",
    options: [
      { value: "active", label: "Active", icon: CheckCircle },
      { value: "draft", label: "Draft", icon: Clock },
      { value: "archived", label: "Archived", icon: XCircle }
    ]
  },
  {
    key: "category",
    label: "Category",
    options: categories.map((category) => ({ value: category, label: category, icon: Tag }))
  },
  {
    key: "createdBy",
    label: "Created By",
    options: [...new Set(prompts.value.map((prompt) => prompt.createdBy))].map((creator) => ({
      value: creator,
      label: creator,
      icon: User
    }))
  }
]

const searchConfig: SearchConfig<PromptData> = {
  fields: ["name", "description", "category"],
  placeholder: "Search prompts, filter by status, category..."
}

const displayConfig: DisplayConfig<PromptData> = {
  getItemText: (prompt) => `${prompt.name} - ${prompt.description}`,
  getItemValue: (prompt) => prompt.name,
  getItemIcon: () => FileText
}

// Event handlers for SearchFilter
function handleFiltersChanged(filters: Record<string, string>) {
  activeFilters.value = filters
}

function handleItemSelected(prompt: PromptData) {
  navigateTo(`/prompts/${prompt.id}`)
}

const filteredPrompts = computed(() => {
  let filtered = prompts.value

  // Apply SearchFilter filters
  if (activeFilters.value.name && activeFilters.value.name !== "all") {
    filtered = filtered.filter((prompt) => prompt.name === activeFilters.value.name)
  }
  if (activeFilters.value.status && activeFilters.value.status !== "all") {
    filtered = filtered.filter((prompt) => prompt.status === activeFilters.value.status)
  }
  if (activeFilters.value.category && activeFilters.value.category !== "all") {
    filtered = filtered.filter((prompt) => prompt.category === activeFilters.value.category)
  }
  if (activeFilters.value.createdBy && activeFilters.value.createdBy !== "all") {
    filtered = filtered.filter((prompt) => prompt.createdBy === activeFilters.value.createdBy)
  }

  return filtered
})

// Icon and variant arrays for index matching with API data
const statsIcons = [FileText, CheckCircle, Clock, Tag]
const statsVariants: StatsCardProps["variant"][] = ["default", "chart-2", "chart-3", "chart-1"]

// Stats configuration - ready for API replacement
const statsCards = computed(() => {
  return [
    {
      title: "Total Prompts",
      value: prompts.value.length,
      description: "All prompt templates"
    },
    {
      title: "Active Prompts",
      value: prompts.value.filter((prompt) => prompt.status === "active").length,
      description: "Ready for production use"
    },
    {
      title: "Draft Prompts",
      value: prompts.value.filter((prompt) => prompt.status === "draft").length,
      description: "In development"
    },
    {
      title: "Total Usage",
      value: prompts.value.reduce((sum, prompt) => sum + prompt.usageCount, 0),
      description: "Across all prompts"
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
    <PageHeader title="Prompts" :subtext="`Manage your prompt library and templates for ${appConfig.app.name}`">
      <UiButton @click="navigateTo('/prompts/create')" class="gap-2">
        <Plus class="h-4 w-4" />
        New Prompt
      </UiButton>
    </PageHeader>

    <!-- Stats Cards -->
    <div class="grid grid-cols-1 sm:grid-cols-2 lg:grid-cols-4 gap-6">
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

    <!-- Search & Filter Component -->
    <SearchFilter
      :items="prompts"
      :filters="filterConfigs"
      :search-config="searchConfig"
      :display-config="displayConfig"
      @filters-changed="handleFiltersChanged"
      @item-selected="handleItemSelected"
    />

    <!-- Prompts List -->
    <CardsDataList title="All Prompts" :icon="FileText">
      <div class="space-y-4">
        <div
          v-for="prompt in filteredPrompts"
          :key="prompt.id"
          class="flex items-center justify-between p-4 border rounded-lg hover:bg-muted/50 cursor-pointer"
          @click="handlePromptSelect(prompt)"
        >
          <div class="flex items-center gap-4 flex-1">
            <div class="flex-shrink-0">
              <div class="w-10 h-10 rounded-lg bg-primary/10 flex items-center justify-center">
                <FileText class="h-5 w-5 text-primary" />
              </div>
            </div>

            <div class="space-y-1 flex-1 min-w-0">
              <div class="flex items-center gap-2">
                <p class="font-medium truncate">{{ prompt.name }}</p>
                <div
                  class="inline-flex items-center gap-1 px-2 py-1 rounded-md text-xs font-medium border flex-shrink-0"
                  :class="getStatusColor(prompt.status)"
                >
                  <div
                    class="size-1.5 rounded-full"
                    :class="
                      prompt.status === 'active'
                        ? 'bg-green-500'
                        : prompt.status === 'draft'
                          ? 'bg-yellow-500'
                          : 'bg-gray-400'
                    "
                  />
                  {{ prompt.status === "active" ? "Active" : prompt.status === "draft" ? "Draft" : "Archived" }}
                </div>
              </div>
              <p class="text-sm text-muted-foreground line-clamp-2">{{ prompt.description }}</p>
              <div class="flex items-center gap-4 text-xs text-muted-foreground">
                <span class="flex items-center gap-1">
                  <Tag class="h-3 w-3" />
                  {{ prompt.category }}
                </span>
                <span>•</span>
                <span class="flex items-center gap-1">
                  <User class="h-3 w-3" />
                  {{ prompt.createdBy }}
                </span>
                <span>•</span>
                <span>Updated {{ formatDate(prompt.updatedAt) }}</span>
              </div>
            </div>

            <div class="text-right flex-shrink-0">
              <div class="text-sm font-medium">{{ prompt.usageCount.toLocaleString() }}</div>
              <div class="text-xs text-muted-foreground">uses</div>
            </div>
          </div>
        </div>
      </div>
    </CardsDataList>
  </div>
</template>

