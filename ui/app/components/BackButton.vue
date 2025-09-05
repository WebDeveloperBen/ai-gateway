<script setup lang="ts">
import { ArrowLeft } from "lucide-vue-next"
import type { ButtonProps } from "./Ui/Button.vue"

interface Props {
  to: string
  label?: string
  variant?: "ghost" | "outline" | "default"
  size?: ButtonProps["size"]
}

const props = withDefaults(defineProps<Props>(), {
  label: "Back",
  variant: "outline",
  size: "sm"
})
</script>

<template>
  <UiButton :variant="props.variant" :size="props.size" as-child>
    <NuxtLink :to="props.to" class="inline-flex items-center gap-2">
      <ArrowLeft class="h-4 w-4" />
      {{ props.label }}
    </NuxtLink>
  </UiButton>
</template>

