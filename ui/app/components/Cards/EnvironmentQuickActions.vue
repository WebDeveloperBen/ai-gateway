<script setup lang="ts">
import { BarChart3, Download, Copy } from "lucide-vue-next"
</script>

<template>
  <UiCard>
    <UiCardHeader>
      <UiCardTitle class="text-sm">Quick Actions</UiCardTitle>
    </UiCardHeader>
    <UiCardContent class="space-y-2">
      <UiButton variant="outline" size="sm" class="w-full justify-start">
        <Download class="size-4 mr-2" />
        Export Config
      </UiButton>
      <UiButton variant="outline" size="sm" class="w-full justify-start">
        <Copy class="size-4 mr-2" />
        Clone Environment
      </UiButton>
      <UiButton variant="outline" size="sm" class="w-full justify-start">
        <BarChart3 class="size-4 mr-2" />
        View Analytics
      </UiButton>
    </UiCardContent>
  </UiCard>
</template>
