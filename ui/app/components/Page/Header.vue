<script setup lang="ts">
withDefaults(defineProps<{ title: string; subtext: string }>(), {})
</script>
<template>
  <div class="flex items-center justify-between">
    <div>
      <h1 class="text-3xl font-bold tracking-tight text-primary">
        {{ title }}
      </h1>
      <p class="text-muted-foreground">{{ subtext }}</p>
    </div>
    <slot />
  </div>
</template>
