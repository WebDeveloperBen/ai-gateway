<script setup lang="ts">
import type { FunctionalComponent } from "vue"

interface Props {
  label: string
  value: string | number
  icon: FunctionalComponent
  change?: string
  iconColor?: string
}

withDefaults(defineProps<Props>(), {
  iconColor: "text-primary"
})
</script>

<template>
  <UiCard class="p-6">
    <div class="flex items-center gap-4">
      <div class="rounded-lg bg-primary/10 p-3">
        <component :is="icon" class="size-5" :class="iconColor" />
      </div>
      <div class="flex-1">
        <p class="text-2xl font-semibold">{{ value }}</p>
        <p class="text-sm text-muted-foreground">{{ label }}</p>
      </div>
    </div>
    <div v-if="change" class="mt-4">
      <p class="text-xs text-muted-foreground">{{ change }}</p>
    </div>
  </UiCard>
</template>

