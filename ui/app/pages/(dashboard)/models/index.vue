<script lang="ts">
interface ModelData {
  id: string
  name: string
  provider: string
  deployment: string
  endpoint: string
  status: "active" | "inactive" | "error"
  modelType: string
  version: string
  maxTokens: number
  requestCount: number
  lastUsed: string
  applications: string[]
  costPerToken: number
}

const models = ref<ModelData[]>([
  {
    id: "model_1",
    name: "GPT-4 Production",
    provider: "OpenAI",
    deployment: "gpt-4-deployment-prod",
    endpoint: "https://api.openai.com/v1",
    status: "active" as const,
    modelType: "gpt-4",
    version: "gpt-4-0613",
    maxTokens: 8192,
    requestCount: 12450,
    lastUsed: "2025-01-15T14:30:00Z",
    applications: ["Customer Service Bot", "Content Generator"],
    costPerToken: 0.00003
  },
  {
    id: "model_2", 
    name: "GPT-3.5 Turbo Dev",
    provider: "Azure OpenAI",
    deployment: "gpt-35-turbo-dev",
    endpoint: "https://myorg.openai.azure.com/",
    status: "active" as const,
    modelType: "gpt-3.5-turbo",
    version: "gpt-3.5-turbo-0613",
    maxTokens: 4096,
    requestCount: 8920,
    lastUsed: "2025-01-15T10:15:00Z",
    applications: ["Code Assistant"],
    costPerToken: 0.000002
  },
  {
    id: "model_3",
    name: "Claude-3 Sonnet",
    provider: "Anthropic",
    deployment: "claude-3-sonnet",
    endpoint: "https://api.anthropic.com/v1",
    status: "inactive" as const,
    modelType: "claude-3-sonnet",
    version: "20240229",
    maxTokens: 200000,
    requestCount: 0,
    lastUsed: "2024-12-28T16:45:00Z",
    applications: [],
    costPerToken: 0.000015
  },
  {
    id: "model_4",
    name: "GPT-4 Error Test",
    provider: "Azure OpenAI",
    deployment: "gpt-4-error-test",
    endpoint: "https://test.openai.azure.com/",
    status: "error" as const,
    modelType: "gpt-4",
    version: "gpt-4-0613",
    maxTokens: 8192,
    requestCount: 156,
    lastUsed: "2025-01-14T09:22:00Z",
    applications: [],
    costPerToken: 0.00003
  }
])
</script>
<script setup lang="ts">
import { Plus, Cpu, CheckCircle, XCircle, AlertTriangle, Activity, DollarSign, Zap, Layers } from "lucide-vue-next"
import type { StatsCardProps } from "@/components/Cards/Stats.vue"
import SearchFilter from "@/components/SearchFilter.vue"
import type { FilterConfig, SearchConfig, DisplayConfig } from "@/components/SearchFilter.vue"

const showCreateModal = ref(false)
const activeFilters = ref<Record<string, string>>({})
const appConfig = useAppConfig()

const handleModelSelect = (model: ModelData) => {
  navigateTo(`/models/${model.id}`)
}

const onModelCreated = (modelId: string) => {
  console.log("Model created:", modelId)
  navigateTo("/models", { replace: true })
}

const route = useRoute()
onMounted(() => {
  if (route.query.create === "model") {
    showCreateModal.value = true
  }
})

const providers = [...new Set(models.value.map((model) => model.provider))]
const modelTypes = [...new Set(models.value.map((model) => model.modelType))]

const filterConfigs: FilterConfig[] = [
  {
    key: "name",
    label: "Model",
    options: models.value.map((model) => ({ value: model.name, label: model.name, icon: Cpu }))
  },
  {
    key: "provider", 
    label: "Provider",
    options: providers.map((provider) => ({ value: provider, label: provider, icon: Layers }))
  },
  {
    key: "modelType",
    label: "Model Type", 
    options: modelTypes.map((type) => ({ value: type, label: type, icon: Cpu }))
  },
  {
    key: "status",
    label: "Status",
    options: [
      { value: "active", label: "Active", icon: CheckCircle },
      { value: "inactive", label: "Inactive", icon: XCircle },
      { value: "error", label: "Error", icon: AlertTriangle }
    ]
  }
]

const searchConfig: SearchConfig<ModelData> = {
  fields: ["name", "provider", "modelType", "deployment"],
  placeholder: "Search models, filter by provider, type, status..."
}

const displayConfig: DisplayConfig<ModelData> = {
  getItemText: (model) => `${model.name} - ${model.provider} ${model.modelType}`,
  getItemValue: (model) => model.name,
  getItemIcon: () => Cpu
}

function handleFiltersChanged(filters: Record<string, string>) {
  activeFilters.value = filters
}

function handleItemSelected(model: ModelData) {
  navigateTo(`/models/${model.id}`)
}

const filteredModels = computed(() => {
  let filtered = models.value

  if (activeFilters.value.name && activeFilters.value.name !== "all") {
    filtered = filtered.filter((model) => model.name === activeFilters.value.name)
  }
  if (activeFilters.value.provider && activeFilters.value.provider !== "all") {
    filtered = filtered.filter((model) => model.provider === activeFilters.value.provider)
  }
  if (activeFilters.value.modelType && activeFilters.value.modelType !== "all") {
    filtered = filtered.filter((model) => model.modelType === activeFilters.value.modelType)
  }
  if (activeFilters.value.status && activeFilters.value.status !== "all") {
    filtered = filtered.filter((model) => model.status === activeFilters.value.status)
  }

  return filtered
})

const statsIcons = [Cpu, CheckCircle, Activity, DollarSign]
const statsVariants: StatsCardProps["variant"][] = ["default", "chart-2", "chart-1", "chart-3"]

const statsCards = computed(() => {
  const totalCost = models.value.reduce((sum, model) => sum + (model.requestCount * model.costPerToken), 0)
  
  return [
    {
      title: "Total Models",
      value: models.value.length,
      description: "Registered model deployments"
    },
    {
      title: "Active Models", 
      value: models.value.filter((model) => model.status === "active").length,
      description: "Currently available for routing"
    },
    {
      title: "Total Requests",
      value: models.value.reduce((sum, model) => sum + model.requestCount, 0),
      description: "Across all model deployments"
    },
    {
      title: "Estimated Cost",
      value: `$${totalCost.toFixed(2)}`,
      description: "Based on token usage"
    }
  ].map((stat, index) => ({
    ...stat,
    icon: statsIcons[index],
    variant: statsVariants[index]
  }))
})
</script>

<template>
  <div class="flex flex-col gap-6">
    <PageHeader
      title="Models"
      :subtext="`Manage your AI model deployments and ${appConfig.app.name} routing configuration`"
    >
      <UiButton @click="showCreateModal = true" class="gap-2">
        <Plus class="h-4 w-4" />
        Register Model
      </UiButton>
    </PageHeader>

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

    <SearchFilter
      :items="models"
      :filters="filterConfigs"
      :search-config="searchConfig"
      :display-config="displayConfig"
      @filters-changed="handleFiltersChanged"
      @item-selected="handleItemSelected"
    />

    <ModelsList :models="filteredModels" @select-model="handleModelSelect" />

    <LazyModalsModelsCreate v-model:open="showCreateModal" @created="onModelCreated" />
  </div>
</template>