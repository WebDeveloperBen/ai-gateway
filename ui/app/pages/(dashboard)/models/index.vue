<script lang="ts">
interface ProviderData {
  id: string
  name: string
  type: string
  status: "active" | "inactive" | "error"
  endpoint: string
  modelsCount: number
  requestCount: number
  lastUsed: string
  apiKeyMasked: string
  scope: "organization" | "workspace"
  createdAt: string
}

const providers = ref<ProviderData[]>([
  {
    id: "provider_1",
    name: "OpenAI",
    type: "OpenAI",
    status: "active" as const,
    endpoint: "https://api.openai.com/v1",
    modelsCount: 4,
    requestCount: 15420,
    lastUsed: "2025-01-15T14:30:00Z",
    apiKeyMasked: "sk-...J2k9",
    scope: "organization",
    createdAt: "2024-12-01T10:00:00Z"
  },
  {
    id: "provider_2",
    name: "Azure OpenAI",
    type: "Azure OpenAI",
    status: "active" as const,
    endpoint: "https://myorg.openai.azure.com/",
    modelsCount: 2,
    requestCount: 8920,
    lastUsed: "2025-01-15T10:15:00Z",
    apiKeyMasked: "abc...xyz",
    scope: "workspace",
    createdAt: "2024-11-15T14:20:00Z"
  },
  {
    id: "provider_4",
    name: "Custom AI Endpoint",
    type: "Custom",
    status: "error" as const,
    endpoint: "https://api.custom.ai/v1",
    modelsCount: 0,
    requestCount: 156,
    lastUsed: "2025-01-14T09:22:00Z",
    apiKeyMasked: "key...789",
    scope: "workspace",
    createdAt: "2025-01-10T11:45:00Z"
  }
])
</script>
<script setup lang="ts">
import {
  Plus,
  Server,
  CheckCircle,
  XCircle,
  AlertTriangle,
  Activity,
  Layers,
  Eye,
  Settings,
  Trash2,
  MoreVertical
} from "lucide-vue-next"
import type { StatsCardProps } from "@/components/Cards/Stats.vue"
import SearchFilter from "@/components/SearchFilter.vue"
import type { FilterConfig, SearchConfig, DisplayConfig } from "@/components/SearchFilter.vue"

const showCreateModal = ref(false)
const activeFilters = ref<Record<string, string>>({})
const appConfig = useAppConfig()

const handleProviderSelect = (provider: ProviderData) => {
  navigateTo(`/models/providers/${provider.id}`)
}

const onProviderCreated = (providerId: string) => {
  console.log("Provider created:", providerId)
  navigateTo("/models/providers", { replace: true })
}

const route = useRoute()
onMounted(() => {
  if (route.query.create === "provider") {
    showCreateModal.value = true
  }
})

const providerTypes = [...new Set(providers.value.map((provider) => provider.type))]

const filterConfigs: FilterConfig[] = [
  {
    key: "name",
    label: "Provider",
    options: providers.value.map((provider) => ({ value: provider.name, label: provider.name, icon: Server }))
  },
  {
    key: "type",
    label: "Type",
    options: providerTypes.map((type) => ({ value: type, label: type, icon: Layers }))
  },
  {
    key: "scope",
    label: "Scope",
    options: [
      { value: "organization", label: "Organization", icon: Layers },
      { value: "workspace", label: "Workspace", icon: Server }
    ]
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

const searchConfig: SearchConfig<ProviderData> = {
  fields: ["name", "type", "endpoint"],
  placeholder: "Search providers, filter by type, status, scope..."
}

const displayConfig: DisplayConfig<ProviderData> = {
  getItemText: (provider) => `${provider.name} - ${provider.type}`,
  getItemValue: (provider) => provider.name,
  getItemIcon: () => Server
}

function handleFiltersChanged(filters: Record<string, string>) {
  activeFilters.value = filters
}

function handleItemSelected(provider: ProviderData) {
  navigateTo(`/models/providers/${provider.id}`)
}

const filteredProviders = computed(() => {
  let filtered = providers.value

  if (activeFilters.value.name && activeFilters.value.name !== "all") {
    filtered = filtered.filter((provider) => provider.name === activeFilters.value.name)
  }
  if (activeFilters.value.type && activeFilters.value.type !== "all") {
    filtered = filtered.filter((provider) => provider.type === activeFilters.value.type)
  }
  if (activeFilters.value.scope && activeFilters.value.scope !== "all") {
    filtered = filtered.filter((provider) => provider.scope === activeFilters.value.scope)
  }
  if (activeFilters.value.status && activeFilters.value.status !== "all") {
    filtered = filtered.filter((provider) => provider.status === activeFilters.value.status)
  }

  return filtered
})

const statsIcons = [Server, CheckCircle, Activity, Layers]
const statsVariants: StatsCardProps["variant"][] = ["default", "chart-2", "chart-1", "chart-3"]

const statsCards = computed(() => {
  return [
    {
      title: "Total Providers",
      value: providers.value.length,
      description: "Registered AI service providers"
    },
    {
      title: "Active Providers",
      value: providers.value.filter((provider) => provider.status === "active").length,
      description: "Currently operational"
    },
    {
      title: "Total Requests",
      value: providers.value.reduce((sum, provider) => sum + provider.requestCount, 0),
      description: "Across all providers"
    },
    {
      title: "Available Models",
      value: providers.value.reduce((sum, provider) => sum + provider.modelsCount, 0),
      description: "Total model deployments"
    }
  ].map((stat, index) => ({
    ...stat,
    icon: statsIcons[index],
    variant: statsVariants[index]
  }))
})

const getProviderIcon = (type: string) => {
  const icons: Record<string, string> = {
    OpenAI: "🤖",
    "Azure OpenAI": "☁️",
    Anthropic: "🧠",
    "Google AI": "🔍",
    Cohere: "⚡",
    Custom: "🔧"
  }
  return icons[type] || "🔧"
}

const getStatusBadgeClass = (status: string) => {
  switch (status) {
    case "active":
      return "!bg-chart-2/10 !text-chart-2 !border-chart-2/20"
    case "error":
      return "!bg-destructive/10 !text-destructive !border-destructive/20"
    case "inactive":
    default:
      return "!bg-muted !text-muted-foreground !border-border"
  }
}

const getTypeBadgeClass = (type: string) => {
  const classes: Record<string, string> = {
    OpenAI: "bg-primary/10 text-primary border border-primary/20",
    "Azure OpenAI": "bg-chart-1/10 text-chart-1 border border-chart-1/20",
    Anthropic: "bg-chart-3/10 text-chart-3 border border-chart-3/20",
    "Google AI": "bg-chart-4/10 text-chart-4 border border-chart-4/20",
    Cohere: "bg-chart-2/10 text-chart-2 border border-chart-2/20"
  }
  return classes[type] || "bg-muted/10 text-muted-foreground border border-border"
}
</script>

<template>
  <div class="flex flex-col gap-6">
    <PageHeader
      title="Providers"
      :subtext="`Manage AI service providers and their integration with ${appConfig.app.name}`"
    >
      <UiButton @click="showCreateModal = true" class="gap-2">
        <Plus class="h-4 w-4" />
        Add Provider
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
      :items="providers"
      :filters="filterConfigs"
      :search-config="searchConfig"
      :display-config="displayConfig"
      @filters-changed="handleFiltersChanged"
      @item-selected="handleItemSelected"
    />

    <div class="flex flex-col gap-4">
      <div v-if="filteredProviders.length === 0" class="text-center py-12">
        <AlertTriangle class="mx-auto h-12 w-12 text-muted-foreground" />
        <h3 class="mt-4 text-lg font-medium">No providers found</h3>
        <p class="text-muted-foreground">Try adjusting your search or add a new provider.</p>
      </div>

      <UiCard
        v-for="provider in filteredProviders"
        :key="provider.id"
        interactive
        @click="handleProviderSelect(provider)"
      >
        <UiCardHeader>
          <div class="flex items-start justify-between">
            <div class="space-y-1 flex-1">
              <div class="flex items-center gap-2">
                <span class="text-xl">{{ getProviderIcon(provider.type) }}</span>
                <UiCardTitle class="text-lg">{{ provider.name }}</UiCardTitle>
                <UiBadge variant="outline" :class="getStatusBadgeClass(provider.status)">
                  {{ provider.status }}
                </UiBadge>
              </div>
              <UiCardDescription class="text-sm">
                {{ provider.endpoint }}
              </UiCardDescription>
              <div class="flex items-center gap-4 text-xs">
                <div class="flex items-center gap-1">
                  <span class="text-muted-foreground">Type:</span>
                  <UiBadge :class="getTypeBadgeClass(provider.type)" class="text-xs">
                    {{ provider.type }}
                  </UiBadge>
                </div>
                <span class="text-muted-foreground">•</span>
                <div class="flex items-center gap-1">
                  <span class="text-muted-foreground">Scope:</span>
                  <span class="font-medium text-foreground capitalize">{{ provider.scope }}</span>
                </div>
                <span class="text-muted-foreground">•</span>
                <div class="flex items-center gap-1">
                  <span class="text-muted-foreground">API Key:</span>
                  <span class="font-mono text-foreground">{{ provider.apiKeyMasked }}</span>
                </div>
              </div>
            </div>
            <UiDropdownMenu>
              <UiDropdownMenuTrigger as-child>
                <UiButton variant="ghost" size="sm" @click.stop>
                  <MoreVertical class="h-4 w-4" />
                </UiButton>
              </UiDropdownMenuTrigger>
              <UiDropdownMenuContent align="end" class="w-48">
                <UiDropdownMenuItem @click="navigateTo(`/models/providers/${provider.id}`)">
                  <Eye class="mr-2 size-4" />
                  View Details
                </UiDropdownMenuItem>
                <UiDropdownMenuItem @click="navigateTo(`/models/providers/${provider.id}/config`)">
                  <Settings class="mr-2 size-4" />
                  Configure
                </UiDropdownMenuItem>
                <UiDropdownMenuItem @click="navigateTo(`/models/providers/${provider.id}/models`)">
                  <Layers class="mr-2 size-4" />
                  Manage Models
                </UiDropdownMenuItem>
                <UiDropdownMenuSeparator />
                <UiDropdownMenuItem class="text-destructive">
                  <Trash2 class="mr-2 size-4" />
                  Remove Provider
                </UiDropdownMenuItem>
              </UiDropdownMenuContent>
            </UiDropdownMenu>
          </div>
        </UiCardHeader>
        <UiCardContent>
          <div class="grid grid-cols-1 md:grid-cols-3 gap-4 text-sm">
            <div class="flex items-center gap-2">
              <Layers class="h-4 w-4 text-chart-3" />
              <div>
                <p class="text-muted-foreground">Models</p>
                <p class="font-medium">{{ provider.modelsCount }}</p>
              </div>
            </div>
            <div class="flex items-center gap-2">
              <Activity class="h-4 w-4 text-chart-1" />
              <div>
                <p class="text-muted-foreground">Requests</p>
                <p class="font-medium">{{ formatNumber(provider.requestCount) }}</p>
              </div>
            </div>
            <div class="flex items-center gap-2">
              <CheckCircle class="h-4 w-4 text-chart-2" />
              <div>
                <p class="text-muted-foreground">Last Used</p>
                <p class="font-medium">
                  {{ new Date(provider.lastUsed).toLocaleDateString() }}
                </p>
              </div>
            </div>
          </div>
        </UiCardContent>
      </UiCard>
    </div>

    <!-- Provider Creation Wizard Modal -->
    <LazyModalsProvidersCreate v-model:open="showCreateModal" @created="onProviderCreated" />
  </div>
</template>

