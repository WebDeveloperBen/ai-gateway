<script setup lang="ts">
interface Props {
  value: number
  highThreshold?: number
  mediumThreshold?: number
  size?: "sm" | "default"
}

const props = withDefaults(defineProps<Props>(), {
  highThreshold: 20000,
  mediumThreshold: 5000,
  size: "default"
})

const getIndicatorClass = () => {
  if (props.value > props.highThreshold) {
    return 'bg-chart-2' // High activity - green
  } else if (props.value > props.mediumThreshold) {
    return 'bg-chart-1' // Medium activity - orange
  } else {
    return 'bg-muted' // Low activity - muted
  }
}

const sizeClass = computed(() => {
  return props.size === "sm" ? "w-1.5 h-1.5" : "w-2 h-2"
})
</script>

<template>
  <div :class="[getIndicatorClass(), sizeClass, 'rounded-full']"></div>
</template>