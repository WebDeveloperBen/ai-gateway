<script lang="ts">
export type StatsCardProps = {
  title: string
  value: string | number
  icon: any
  description?: string
  textColor?: string
  variant?: "default" | "chart-1" | "chart-2" | "chart-3" | "chart-4"
}
</script>
<script setup lang="ts">
const props = withDefaults(defineProps<StatsCardProps>(), {
  textColor: "text-foreground",
  variant: "default"
})

const getVariantClasses = (variant: string) => {
  switch (variant) {
    case "chart-1":
      return {
        background: "bg-chart-1/10 border-chart-1/20",
        icon: "text-chart-1",
        value: "text-chart-1"
      }
    case "chart-2":
      return {
        background: "bg-chart-2/10 border-chart-2/20",
        icon: "text-chart-2",
        value: "text-chart-2"
      }
    case "chart-3":
      return {
        background: "bg-chart-3/10 border-chart-3/20",
        icon: "text-chart-3",
        value: "text-chart-3"
      }
    case "chart-4":
      return {
        background: "bg-chart-4/10 border-chart-4/20",
        icon: "text-chart-4",
        value: "text-chart-4"
      }
    default:
      return {
        background: "bg-primary/10 border-primary/20",
        icon: "text-primary",
        value: props.textColor
      }
  }
}
</script>

<template>
  <UiCard>
    <UiCardHeader class="flex flex-row items-center justify-between space-y-0 pb-2">
      <UiCardTitle class="text-sm px-0 font-medium">{{ title }}</UiCardTitle>
      <div class="rounded-lg bg-primary/10 p-3">
        <component :is="icon" class="size-5 text-primary" />
      </div>
    </UiCardHeader>
    <UiCardContent>
      <div class="text-2xl font-bold" :class="getVariantClasses(variant).value">
        {{ typeof value === "number" ? formatNumber(value) : value }}
      </div>
      <p v-if="description" class="text-xs text-muted-foreground">
        {{ description }}
      </p>
    </UiCardContent>
  </UiCard>
</template>
