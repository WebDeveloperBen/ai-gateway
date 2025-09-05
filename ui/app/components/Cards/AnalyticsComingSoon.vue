<script setup lang="ts">
import { Activity } from "lucide-vue-next"
</script>
<template>
  <UiCard>
    <UiCardHeader>
      <UiCardTitle class="flex items-center gap-2">
        <Activity class="h-5 w-5" />
        Usage Analytics
      </UiCardTitle>
      <UiCardDescription>Request volume and performance metrics over time</UiCardDescription>
    </UiCardHeader>
    <UiCardContent>
      <ChartPlaceholder
        :icon="Activity"
        title="Charts Coming Soon"
        description="Interactive usage analytics and performance charts will be available here"
      />
    </UiCardContent>
  </UiCard>
</template>
