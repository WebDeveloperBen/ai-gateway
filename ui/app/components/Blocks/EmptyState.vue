<script setup lang="ts">
import { MessageSquare } from "lucide-vue-next"

defineProps<{ title: string; subtext: string }>()
</script>
<template>
  <div class="flex flex-col h-full items-center justify-center flex-1 text-center text-muted-foreground p-8">
    <MessageSquare class="h-12 w-12 mb-3 opacity-50" />
    <p class="text-sm mb-1">{{ title }}</p>
    <p class="text-xs opacity-70">{{ subtext }}</p>
  </div>
</template>
