<script setup lang="ts">
interface Props {
  status: "active" | "inactive" | string
  variant?: "default" | "compact"
}

withDefaults(defineProps<Props>(), {
  variant: "default"
})

const getBadgeClass = (status: string) => {
  switch (status) {
    case "active":
      return "!bg-chart-2/10 !text-chart-2 !border-chart-2/20"
    case "inactive":
      return "!bg-muted !text-muted-foreground !border-border"
    default:
      return "!bg-muted !text-muted-foreground !border-border"
  }
}
</script>

<template>
  <UiBadge variant="outline" :class="getBadgeClass(status)">
    {{ status }}
  </UiBadge>
</template>
