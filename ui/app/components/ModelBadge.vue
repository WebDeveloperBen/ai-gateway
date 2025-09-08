<script setup lang="ts">
interface Props {
  model: string
  size?: "sm" | "default"
}

const props = withDefaults(defineProps<Props>(), {
  size: "default"
})

const getBadgeClass = (model: string) => {
  switch (model) {
    case "gpt-4":
      return "bg-primary/10 text-primary border border-primary/20"
    case "gpt-3.5-turbo":
      return "bg-chart-4/10 text-chart-4 border border-chart-4/20"
    default:
      return "bg-muted/10 text-muted-foreground border border-border"
  }
}

const sizeClass = computed(() => {
  return props.size === "sm" ? "text-xs px-2 py-0.5" : ""
})
</script>

<template>
  <UiBadge :class="[getBadgeClass(model), sizeClass]">
    {{ model }}
  </UiBadge>
</template>