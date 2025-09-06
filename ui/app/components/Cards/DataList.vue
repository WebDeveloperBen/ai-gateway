<script setup lang="ts">
import { Users } from "lucide-vue-next"
import type { FunctionalComponent } from "vue"

withDefaults(defineProps<{ title: string; icon: FunctionalComponent }>(), { title: "Change me", icon: () => Users })
</script>
<template>
  <UiCard class="w-full min-w-0">
    <UiCardHeader>
      <div class="flex items-center justify-between">
        <UiCardTitle class="flex items-center gap-2 min-w-0">
          <component :is="icon" class="size-5 flex-shrink-0" />
          <span class="truncate">{{ title }}</span>
        </UiCardTitle>

        <!-- action bar -->
        <div class="flex-shrink-0">
          <slot name="actions" />
        </div>
      </div>
    </UiCardHeader>
    <UiCardContent class="min-w-0">
      <slot />
    </UiCardContent>
  </UiCard>
</template>
