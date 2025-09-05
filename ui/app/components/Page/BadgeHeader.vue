<script setup lang="ts">
const props = withDefaults(
  defineProps<{
    title: string
    badgeStatus: string
    subtext: string
  }>(),
  {}
)

const getStatusBadgeClass = computed(() => {
  return props.badgeStatus === "active"
    ? "bg-green-100 text-green-800 dark:bg-green-900 dark:text-green-200"
    : "bg-gray-100 text-gray-800 dark:bg-gray-900 dark:text-gray-200"
})
</script>
<template>
  <div class="flex items-center justify-between">
    <div class="flex items-start justify-between">
      <div>
        <div class="flex items-center gap-3">
          <h1 class="text-3xl font-bold text-primary tracking-tight">{{ title }}</h1>
          <UiBadge :class="getStatusBadgeClass">
            {{ badgeStatus }}
          </UiBadge>
        </div>
        <p class="text-muted-foreground">
          {{ subtext }}
        </p>
      </div>
    </div>
    <!-- Right aligned Action bar -->
    <slot />
  </div>
</template>
