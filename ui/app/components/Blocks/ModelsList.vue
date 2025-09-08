<script lang="ts">
import {
  MoreVertical,
  Cpu,
  Activity,
  AlertCircle,
  Circle,
  DollarSign,
  Eye,
  BarChart3,
  Trash2,
  Settings,
  Zap
} from "lucide-vue-next"

interface Model {
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

interface Props {
  models: Model[]
  showEmpty?: boolean
  emptyTitle?: string
  emptyDescription?: string
  showDropdownActions?: boolean
}
</script>
<script setup lang="ts">
withDefaults(defineProps<Props>(), {
  showEmpty: true,
  emptyTitle: "No models found",
  emptyDescription: "Try adjusting your search or register a new model deployment.",
  showDropdownActions: true
})

const emit = defineEmits<{
  selectModel: [model: Model]
  deleteModel: [model: Model]
}>()

const getProviderBadgeClass = (provider: string) => {
  switch (provider) {
    case "OpenAI":
      return "bg-primary/10 text-primary border border-primary/20"
    case "Azure OpenAI":
      return "bg-chart-1/10 text-chart-1 border border-chart-1/20"
    case "Anthropic":
      return "bg-chart-3/10 text-chart-3 border border-chart-3/20"
    default:
      return "bg-muted/10 text-muted-foreground border border-border"
  }
}

const handleModelClick = (model: Model) => {
  emit("selectModel", model)
}

const deleteModel = (model: Model) => {
  emit("deleteModel", model)
}
</script>

<template>
  <div class="flex flex-col gap-4">
    <div v-if="models.length === 0 && showEmpty" class="text-center py-12">
      <AlertCircle class="mx-auto h-12 w-12 text-muted-foreground" />
      <h3 class="mt-4 text-lg font-medium">{{ emptyTitle }}</h3>
      <p class="text-muted-foreground">{{ emptyDescription }}</p>
    </div>

    <UiCard v-for="model in models" :key="model.id" interactive @click="handleModelClick(model)">
      <UiCardHeader>
        <div class="flex items-start justify-between">
          <div class="space-y-1 flex-1">
            <div class="flex items-center gap-2">
              <Cpu class="h-5 w-5 text-primary" />
              <UiCardTitle class="text-lg">{{ model.name }}</UiCardTitle>
              <UiStatusBadge :status="model.status" />
            </div>
            <UiCardDescription class="text-sm"> {{ model.deployment }} • {{ model.version }} </UiCardDescription>
            <div class="flex items-center gap-4 text-xs">
              <div class="flex items-center gap-1">
                <span class="text-muted-foreground">Provider:</span>
                <UiBadge :class="getProviderBadgeClass(model.provider)" class="text-xs">
                  {{ model.provider }}
                </UiBadge>
              </div>
              <span class="text-muted-foreground">•</span>
              <div class="flex items-center gap-2">
                <span class="text-muted-foreground">Model:</span>
                <ModelBadge :model="model.modelType" size="sm" />
              </div>
              <span class="text-muted-foreground">•</span>
              <div class="flex items-center gap-1">
                <span class="text-muted-foreground">Apps:</span>
                <span class="font-medium text-foreground">{{ model.applications.length }}</span>
              </div>
            </div>
          </div>
          <UiDropdownMenu v-if="showDropdownActions">
            <UiDropdownMenuTrigger as-child>
              <UiButton variant="ghost" size="sm" @click.stop>
                <MoreVertical class="h-4 w-4" />
              </UiButton>
            </UiDropdownMenuTrigger>
            <UiDropdownMenuContent align="end" class="w-48">
              <UiDropdownMenuItem @click="navigateTo(`/models/${model.id}`)">
                <Eye class="mr-2 size-4" />
                View Details
              </UiDropdownMenuItem>
              <UiDropdownMenuItem @click="navigateTo(`/models/${model.id}/config`)">
                <Settings class="mr-2 size-4" />
                Configure Model
              </UiDropdownMenuItem>
              <UiDropdownMenuItem @click="navigateTo(`/models/${model.id}/analytics`)">
                <BarChart3 class="mr-2 size-4" />
                View Analytics
              </UiDropdownMenuItem>
              <UiDropdownMenuSeparator />
              <UiDropdownMenuItem class="text-destructive" @click="deleteModel(model)">
                <Trash2 class="mr-2 size-4" />
                Remove Model
              </UiDropdownMenuItem>
            </UiDropdownMenuContent>
          </UiDropdownMenu>
        </div>
      </UiCardHeader>
      <UiCardContent>
        <div class="grid grid-cols-1 md:grid-cols-4 gap-4 text-sm">
          <div class="flex items-center gap-2">
            <Zap class="h-4 w-4 text-chart-4" />
            <div>
              <p class="text-muted-foreground">Max Tokens</p>
              <p class="font-medium">{{ formatNumber(model.maxTokens) }}</p>
            </div>
          </div>
          <div class="flex items-center gap-2">
            <Activity class="h-4 w-4 text-chart-1" />
            <div>
              <p class="text-muted-foreground">Requests</p>
              <p class="font-medium">{{ formatNumber(model.requestCount) }}</p>
            </div>
          </div>
          <div class="flex items-center gap-2">
            <DollarSign class="h-4 w-4 text-chart-3" />
            <div>
              <p class="text-muted-foreground">Cost/Token</p>
              <p class="font-medium">{{ formatCurrency(model.costPerToken) }}</p>
            </div>
          </div>
          <div class="flex items-center gap-2">
            <Circle class="h-4 w-4 text-chart-2" />
            <div>
              <p class="text-muted-foreground">Last Used</p>
              <p class="font-medium">
                {{ new Date(model.lastUsed).toLocaleDateString() }}
              </p>
            </div>
          </div>
        </div>
      </UiCardContent>
    </UiCard>
  </div>
</template>

