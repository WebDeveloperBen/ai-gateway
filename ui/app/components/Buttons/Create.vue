<script setup lang="ts">
import { Plus } from 'lucide-vue-next';

defineProps<{
  action: () => void
  title: string
}>()
</script>
<template>
  <UiButton @click="action" class="gap-2">
    <Plus class="h-4 w-4" />
    {{ title }}
  </UiButton>
</template>
