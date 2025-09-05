<script setup lang="ts">
import { Eye, EyeOff, RotateCcw, Loader2, Check, Copy } from "lucide-vue-next"
import { toast } from "vue-sonner"

interface Props {
  keyId: string
  keyPrefix: string
  size?: "sm" | "md" | "lg"
  showRegenerateButton?: boolean
  showRegenerateText?: boolean
}

const props = withDefaults(defineProps<Props>(), {
  size: "md",
  showRegenerateButton: true,
  showRegenerateText: false
})

const emit = defineEmits<{
  regenerate: [keyId: string]
}>()

// State management
const revealed = ref(false)
const regenerating = ref(false)

const regeneratedKey = ref<string | null>(null)
const copied = ref(false)

const fakeRegenerateApiCall = async (): Promise<string> => {
  // Stub: Replace with real backend logic
  await new Promise((resolve) => setTimeout(resolve, 700))
  return `${props.keyPrefix}${generateUniqueId(props.keyId)}${Math.random().toString(36).slice(2, 8)}`
}

const regenerateKey = async () => {
  try {
    regenerating.value = true
    toast("Regenerating...", { description: "Please wait while we generate a new API key.", duration: 2000 })

    // Simulate API call
    const newKey = await fakeRegenerateApiCall()
    regeneratedKey.value = newKey
    revealed.value = true
    toast("Key regenerated", { description: "A new API key has been generated.", duration: 3000 })
  } catch (err) {
    toast("Regeneration failed", { description: "Please try again.", duration: 3000 })
  } finally {
    regenerating.value = false
  }
}

const justCopied = shallowRef(false)

const copyToClipboard = async (text: string | null) => {
  if (!text) return
  try {
    await navigator.clipboard.writeText(text)
    justCopied.value = true

    toast({
      title: "Copied to clipboard",
      description: "API key has been copied to your clipboard.",
      duration: 3000,
      icon: "lucide:check"
    })

    // Reset the copied state after 2 seconds
    setTimeout(() => {
      justCopied.value = false
    }, 2000)
  } catch (err) {
    console.error("Failed to copy text: ", err)
    toast({
      title: "Copy failed",
      description: "Failed to copy to clipboard. Please copy manually.",
      duration: 3000,
      icon: "lucide:x",
      variant: "destructive"
    })
  }
}

const generateUniqueId = (keyId: string) => {
  const hash = keyId.split("").reduce((a, b) => {
    a = (a << 5) - a + b.charCodeAt(0)
    return a & a
  }, 0)
  return Math.abs(hash).toString(36).substring(0, 3)
}

const displayText = computed(() => {
  if (regeneratedKey.value) {
    return regeneratedKey.value
  }
  if (revealed.value) {
    return `${props.keyPrefix}${generateUniqueId(props.keyId)}`
  }
  return "sk-" + "*".repeat(48)
})

const sizeClasses = computed(() => {
  switch (props.size) {
    case "sm":
      return {
        container: "p-2",
        text: "text-xs",
        button: "h-6 w-6"
      }
    case "lg":
      return {
        container: "p-4",
        text: "text-base",
        button: "h-8 w-8"
      }
    default: // md
      return {
        container: "p-3",
        text: "text-sm",
        button: "h-7 w-7"
      }
  }
})
</script>

<template>
  <div class="flex items-center gap-3 bg-muted/50 rounded-lg border" :class="sizeClasses.container">
    <code class="font-mono flex-1 break-all" :class="sizeClasses.text">
      {{ displayText }}
    </code>

    <UiButton
      v-if="regeneratedKey"
      variant="outline"
      size="sm"
      @click="copyToClipboard(regeneratedKey)"
      class="gap-2"
      :disabled="justCopied"
    >
      <Check v-if="justCopied" class="h-4 w-4 text-green-600" />
      <Copy v-else class="h-4 w-4" />
      {{ justCopied ? "Copied!" : "Copy" }}
    </UiButton>

    <UiButton
      v-if="showRegenerateButton"
      :variant="props.size === 'sm' ? 'ghost' : 'outline'"
      size="sm"
      @click="regenerateKey"
      :disabled="regenerating"
      :class="props.showRegenerateText ? 'gap-2' : ''"
    >
      <Loader2
        v-if="regenerating"
        class="animate-spin text-blue-600"
        :class="props.size === 'sm' ? 'h-4 w-4' : 'h-4 w-4'"
      />
      <RotateCcw v-else :class="props.size === 'sm' ? 'h-4 w-4' : 'h-4 w-4'" />
      <span v-if="props.showRegenerateText">{{ regenerating ? "Regenerating..." : "Regenerate" }}</span>
    </UiButton>
  </div>
</template>
