<script setup lang="ts">
import { Users } from "lucide-vue-next"
import type { FunctionalComponent } from "vue"

withDefaults(defineProps<{ title: string; icon: FunctionalComponent }>(), { title: "Change me", icon: () => Users })
</script>
<template>
  <!-- Team Members -->
  <UiCard>
    <UiCardHeader>
      <div class="flex items-center justify-between">
        <UiCardTitle class="flex items-center gap-2">
          <component :is="icon" class="size-5" />
          {{ title }}
        </UiCardTitle>

        <!-- action bar -->
        <slot name="actions" />
      </div>
    </UiCardHeader>
    <UiCardContent>
      <slot />
    </UiCardContent>
  </UiCard>
</template>
