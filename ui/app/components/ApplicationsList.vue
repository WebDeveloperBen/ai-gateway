<script lang="ts">
import { MoreVertical, Key, Activity, AlertCircle, Circle, Layers } from "lucide-vue-next"
interface Application {
  id: string
  name: string
  description: string
  status: "active" | "inactive"
  apiKeyCount: number
  monthlyRequests: number
  lastUsed: string
  models: string[]
  team: string
}

interface Props {
  applications: Application[]
  showEmpty?: boolean
  emptyTitle?: string
  emptyDescription?: string
  showDropdownActions?: boolean
}
</script>
<script setup lang="ts">
withDefaults(defineProps<Props>(), {
  showEmpty: true,
  emptyTitle: "No applications found",
  emptyDescription: "Try adjusting your search or create a new application.",
  showDropdownActions: true
})

const emit = defineEmits<{
  selectApplication: [application: Application]
}>()

const formatNumber = (num: number) => {
  return new Intl.NumberFormat().format(num)
}

const handleApplicationClick = (app: Application) => {
  emit("selectApplication", app)
}
</script>

<template>
  <div class="flex flex-col gap-4">
    <div v-if="applications.length === 0 && showEmpty" class="text-center py-12">
      <AlertCircle class="mx-auto h-12 w-12 text-muted-foreground" />
      <h3 class="mt-4 text-lg font-medium">{{ emptyTitle }}</h3>
      <p class="text-muted-foreground">{{ emptyDescription }}</p>
    </div>

    <UiCard v-for="app in applications" :key="app.id" interactive @click="handleApplicationClick(app)">
      <UiCardHeader>
        <div class="flex items-start justify-between">
          <div class="space-y-1 flex-1">
            <div class="flex items-center">
              <Layers class="h-5 w-5 text-primary" />
              <UiCardTitle class="text-lg px-2.5">{{ app.name }}</UiCardTitle>
              <UiStatusBadge :status="app.status" />
            </div>
            <UiCardDescription class="text-sm">
              {{ app.description }}
            </UiCardDescription>
            <div class="flex items-center gap-4 text-xs">
              <div class="flex items-center gap-1">
                <span class="text-muted-foreground">Team:</span>
                <span class="font-medium text-foreground">{{ app.team }}</span>
              </div>
              <span class="text-muted-foreground">•</span>
              <div class="flex items-center gap-2">
                <span class="text-muted-foreground">Models:</span>
                <div class="flex gap-1">
                  <UiModelBadge v-for="model in app.models" :key="model" :model="model" size="sm" />
                </div>
              </div>
            </div>
          </div>
          <UiDropdownMenu v-if="showDropdownActions">
            <UiDropdownMenuTrigger as-child>
              <UiButton variant="ghost" size="sm" @click.stop>
                <MoreVertical class="h-4 w-4" />
              </UiButton>
            </UiDropdownMenuTrigger>
            <UiDropdownMenuContent align="end">
              <UiDropdownMenuItem as-child>
                <NuxtLink :to="`/applications/${app.id}`"> View Details </NuxtLink>
              </UiDropdownMenuItem>
              <UiDropdownMenuItem as-child>
                <NuxtLink :to="`/applications/${app.id}/keys`"> Manage API Keys </NuxtLink>
              </UiDropdownMenuItem>
              <UiDropdownMenuItem> View Analytics </UiDropdownMenuItem>
              <UiDropdownMenuSeparator />
              <UiDropdownMenuItem class="text-red-600"> Delete Application </UiDropdownMenuItem>
            </UiDropdownMenuContent>
          </UiDropdownMenu>
        </div>
      </UiCardHeader>
      <UiCardContent>
        <div class="grid grid-cols-1 md:grid-cols-3 gap-4 text-sm">
          <div class="flex items-center gap-2">
            <Key class="h-4 w-4 text-chart-3" />
            <div>
              <p class="text-muted-foreground">API Keys</p>
              <p class="font-medium">{{ app.apiKeyCount }}</p>
            </div>
          </div>
          <div class="flex items-center gap-2">
            <Activity class="h-4 w-4 text-chart-1" />
            <div>
              <p class="text-muted-foreground">Monthly Requests</p>
              <div class="flex items-center gap-2">
                <p class="font-medium">{{ formatNumber(app.monthlyRequests) }}</p>
                <UiActivityIndicator :value="app.monthlyRequests" />
              </div>
            </div>
          </div>
          <div class="flex items-center gap-2">
            <Circle class="h-4 w-4 text-chart-2" />
            <div>
              <p class="text-muted-foreground">Last Used</p>
              <p class="font-medium">
                {{ new Date(app.lastUsed).toLocaleDateString() }}
              </p>
            </div>
          </div>
        </div>
      </UiCardContent>
    </UiCard>
  </div>
</template>
