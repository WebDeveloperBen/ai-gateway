<script setup lang="ts">
import { Activity } from "lucide-vue-next"
import type { FunctionalComponent } from "vue"

interface ActivityItem {
  timestamp: string
  action: string
  details: string
}

const props = withDefaults(
  defineProps<{
    /** Array of activity items */
    activities: ActivityItem[]
    /** Icon for empty state */
    emptyIcon?: FunctionalComponent
    /** Title for empty state */
    emptyTitle?: string
    /** Description for empty state */
    emptyDescription?: string
  }>(),
  {
    emptyIcon: () => Activity,
    emptyTitle: "No Recent Activity",
    emptyDescription: "Activity will appear here when actions are performed"
  }
)
</script>

<template>
  <div class="space-y-4">
    <div v-if="activities.length === 0" class="text-center py-8">
      <component :is="props.emptyIcon" class="mx-auto h-12 w-12 text-muted-foreground" />
      <h3 class="mt-4 text-lg font-medium">{{ props.emptyTitle }}</h3>
      <p class="text-muted-foreground">{{ props.emptyDescription }}</p>
    </div>

    <div v-for="activity in activities" :key="activity.timestamp" class="flex items-start gap-3 p-3 border rounded-lg">
      <div class="w-2 h-2 rounded-full mt-2 bg-blue-500"></div>
      <div class="flex-1">
        <div class="flex items-center gap-2 mb-1">
          <span class="font-medium text-sm">{{ activity.action }}</span>
          <span class="text-xs text-muted-foreground">
            {{ new Date(activity.timestamp).toLocaleString() }}
          </span>
        </div>
        <p class="text-sm text-muted-foreground">{{ activity.details }}</p>
      </div>
    </div>
  </div>
</template>
