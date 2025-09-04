<script setup lang="ts">
import { Eye, EyeOff, Copy, Check } from "lucide-vue-next"
import { toast } from "vue-sonner"

interface Props {
  keyId: string
  keyPrefix: string
  size?: 'sm' | 'md' | 'lg'
  showCopyButton?: boolean
}

const props = withDefaults(defineProps<Props>(), {
  size: 'md',
  showCopyButton: true
})

// State management
const visible = ref(false)
const fetched = ref<string | null>(null)
const loading = ref(false)
const justCopied = ref(false)

const fetchKey = async () => {
  if (fetched.value || loading.value) return
  
  try {
    loading.value = true
    
    // TODO: Implement actual API call to fetch the key
    // Simulate API call
    await new Promise(resolve => setTimeout(resolve, 500))
    
    // Simulate fetched key - in real app this would come from API
    const fullKey = `sk-proj-${Math.random().toString(36).substring(2, 15)}${Math.random().toString(36).substring(2, 15)}${Math.random().toString(36).substring(2, 10)}`
    
    fetched.value = fullKey
  } catch (error) {
    console.error('Failed to fetch API key:', error)
    toast({
      title: "Error",
      description: "Failed to fetch API key. Please try again.",
      duration: 3000,
      icon: "lucide:x",
      variant: "destructive"
    })
  } finally {
    loading.value = false
  }
}

const toggleVisibility = async () => {
  if (visible.value) {
    visible.value = false
  } else {
    if (!fetched.value) {
      await fetchKey()
    }
    if (fetched.value) {
      visible.value = true
    }
  }
}

const copyToClipboard = async () => {
  try {
    if (!fetched.value) {
      await fetchKey()
    }
    
    if (!fetched.value) {
      throw new Error('Failed to retrieve API key')
    }
    
    await navigator.clipboard.writeText(fetched.value)
    justCopied.value = true
    
    toast({
      title: "Copied to clipboard",
      description: "API key has been copied to your clipboard.",
      duration: 3000,
      icon: "lucide:check"
    })
    
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

const maskKey = (keyPrefix: string) => {
  return keyPrefix + "************************"
}

const displayText = computed(() => {
  if (visible.value && fetched.value) {
    return fetched.value
  }
  if (loading.value && visible.value) {
    return "Loading..."
  }
  return maskKey(props.keyPrefix)
})

const sizeClasses = computed(() => {
  switch (props.size) {
    case 'sm':
      return {
        container: 'p-2',
        text: 'text-xs',
        button: 'h-6 w-6'
      }
    case 'lg':
      return {
        container: 'p-4',
        text: 'text-base',
        button: 'h-8 w-8'
      }
    default: // md
      return {
        container: 'p-3',
        text: 'text-sm',
        button: 'h-7 w-7'
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
      variant="ghost" 
      size="sm" 
      @click="toggleVisibility"
      :disabled="loading"
      :class="sizeClasses.button"
    >
      <div 
        v-if="loading" 
        class="border-2 border-current border-t-transparent rounded-full animate-spin"
        :class="props.size === 'sm' ? 'h-3 w-3' : 'h-4 w-4'"
      ></div>
      <Eye v-else-if="!visible" :class="props.size === 'sm' ? 'h-3 w-3' : 'h-4 w-4'" />
      <EyeOff v-else :class="props.size === 'sm' ? 'h-3 w-3' : 'h-4 w-4'" />
    </UiButton>
    
    <UiButton 
      v-if="showCopyButton"
      variant="outline" 
      size="sm" 
      @click="copyToClipboard"
      :disabled="justCopied"
      :class="[sizeClasses.button, 'gap-1']"
    >
      <Check v-if="justCopied" class="text-green-600" :class="props.size === 'sm' ? 'h-3 w-3' : 'h-4 w-4'" />
      <Copy v-else :class="props.size === 'sm' ? 'h-3 w-3' : 'h-4 w-4'" />
      <span v-if="props.size !== 'sm'" class="text-xs">{{ justCopied ? 'Copied!' : 'Copy' }}</span>
    </UiButton>
  </div>
</template>