<script setup lang="ts">
import type { FunctionalComponent } from "vue"

defineProps<{ action: () => void; title: string; icon: FunctionalComponent }>()
</script>
<template>
  <UiButton variant="default" size="default" class="gap-2" @click="action">
    <component :is="icon" class="h-4 w-4" />
    {{ title }}
  </UiButton>
</template>
