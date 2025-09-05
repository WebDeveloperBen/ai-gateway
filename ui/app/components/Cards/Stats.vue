<script setup lang="ts">
interface StatsCardProps {
  title: string
  value: string | number
  icon: any
  description?: string
  textColor?: string
}

withDefaults(defineProps<StatsCardProps>(), {
  textColor: "text-foreground"
})
</script>

<template>
  <UiCard>
    <UiCardHeader class="flex flex-row items-center justify-between space-y-0 pb-2">
      <UiCardTitle class="text-sm font-medium">{{ title }}</UiCardTitle>
      <component :is="icon" class="h-4 w-4 text-muted-foreground" />
    </UiCardHeader>
    <UiCardContent>
      <div class="text-2xl font-bold" :class="textColor">
        {{ typeof value === "number" ? formatNumber(value) : value }}
      </div>
      <p v-if="description" class="text-xs text-muted-foreground">
        {{ description }}
      </p>
    </UiCardContent>
  </UiCard>
</template>
