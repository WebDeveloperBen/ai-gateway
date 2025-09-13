<template>
  <div class="flex items-center justify-center max-h-64 rounded-lg border-2 border-dashed border-muted">
    <div class="text-center space-y-3">
      <component :is="icon" class="h-12 w-12 text-muted-foreground mx-auto" />
      <div>
        <h3 class="font-medium text-lg">{{ title }}</h3>
        <p class="text-muted-foreground">{{ description }}</p>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { Activity } from "lucide-vue-next"

withDefaults(
  defineProps<{
    /** The icon to display */
    icon?: any
    /** The title text */
    title?: string
    /** The description text */
    description?: string
  }>(),
  {
    icon: Activity,
    title: "Charts Coming Soon",
    description: "Interactive analytics and charts will be available here"
  }
)
</script>

