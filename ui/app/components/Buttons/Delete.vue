<script setup lang="ts">
import { Trash2 } from "lucide-vue-next"

defineProps<{ action: () => void }>()
</script>
<template>
  <UiButton effect="shine" variant="destructive" size="icon" class="" @click="action">
    <Trash2 class="h-4 w-4" />
  </UiButton>
</template>
