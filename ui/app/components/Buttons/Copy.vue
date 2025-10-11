<template>
  <UiTooltip>
    <UiTooltipTrigger as-child>
      <UiButton
        :variant="variant"
        :size="size"
        class="disabled:opacity-100 gap-2"
        :aria-label="copied ? 'Copied' : 'Copy to clipboard'"
        :disabled="copied"
        @click="handleCopy"
      >
        <component v-if="!copied" :is="Copy" class="size-4" aria-hidden="true" />
        <component v-else :is="Check" aria-hidden="true" class="size-4 text-emerald-500" />

        <span v-if="showLabel">
          {{ copied ? "Copied!" : "Copy" }}
        </span>
      </UiButton>
    </UiTooltipTrigger>

    <UiTooltipContent align="center" class="px-2 py-1 text-xs">
      {{ copied ? "Copied" : tooltipText }}
    </UiTooltipContent>
  </UiTooltip>
</template>

<script lang="ts" setup>
import { useClipboard } from "@vueuse/core"
import { Check, Copy } from "lucide-vue-next"
import type { ButtonProps } from "../Ui/Button.vue"

interface Props {
  text: string | null | undefined
  showLabel?: boolean
  size?: ButtonProps["size"]
  variant?: "outline" | "ghost" | "default"
  tooltipText?: string
  copiedDuring?: number // ms; default 1500 like VueUse
}

const props = withDefaults(defineProps<Props>(), {
  showLabel: false,
  size: "sm",
  variant: "outline",
  tooltipText: "Click to copy",
  copiedDuring: 1500
})

const emit = defineEmits<{
  (e: "copied"): void
}>()

const { copy, copied } = useClipboard({ copiedDuring: props.copiedDuring })

const handleCopy = async () => {
  if (!props.text) return
  await copy(props.text)
  emit("copied")
}
</script>
