<script setup lang="ts">
import { Ban } from "lucide-vue-next"

defineProps<{ action: () => void; title: string }>()
</script>
<template>
  <UiButton variant="outline" size="default" class="gap-2 text-orange-600" @click="action">
    <Ban class="h-4 w-4" />
    Disable {{ title }}
  </UiButton>
</template>
