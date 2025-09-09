<script setup lang="ts">
import { FileText, Settings, MessageSquare, Copy, Play, Code, RotateCcw } from "lucide-vue-next"

interface PromptConfig {
  temperature?: number
  maxTokens?: number
  topP?: number
  model?: string
  [key: string]: any
}

interface PromptContent {
  userPrompt: string
  systemPrompt: string
  config: PromptConfig
}

interface Props {
  content: PromptContent
  status?: 'approved' | 'review' | 'draft' | 'archived'
  isEditable?: boolean
  showPlaygroundButton?: boolean
}

const props = withDefaults(defineProps<Props>(), {
  isEditable: false,
  showPlaygroundButton: true
})

const emit = defineEmits<{
  update: [content: PromptContent]
  copy: [type: 'user' | 'system' | 'config', content: string]
  openPlayground: []
}>()

const activeTab = ref<'user' | 'system' | 'config'>('user')

const localContent = ref<PromptContent>({ ...props.content })

// Watch for external content changes
watch(() => props.content, (newContent) => {
  localContent.value = { ...newContent }
  console.log('TabbedContent received new content:', newContent)
}, { deep: true, immediate: true })

// Emit changes when local content updates
const updateContent = () => {
  if (props.isEditable) {
    emit('update', { ...localContent.value })
  }
}

const copyToClipboard = async (type: 'user' | 'system' | 'config') => {
  let content: string
  
  switch (type) {
    case 'user':
      content = localContent.value.userPrompt
      break
    case 'system':
      content = localContent.value.systemPrompt
      break
    case 'config':
      content = JSON.stringify(localContent.value.config, null, 2)
      break
  }
  
  try {
    await navigator.clipboard.writeText(content)
    emit('copy', type, content)
  } catch (err) {
    console.error('Failed to copy content:', err)
  }
}

const clearContent = (type: 'user' | 'system') => {
  if (!props.isEditable) return
  
  if (type === 'user') {
    localContent.value.userPrompt = ''
  } else {
    localContent.value.systemPrompt = ''
  }
  updateContent()
}

const getStatusBadgeClass = computed(() => {
  if (!props.status) return ''
  
  switch (props.status) {
    case 'approved':
      return 'bg-green-100 text-green-800 dark:bg-green-900/20 dark:text-green-400'
    case 'review':
      return 'bg-blue-100 text-blue-800 dark:bg-blue-900/20 dark:text-blue-400'
    case 'draft':
      return 'bg-orange-100 text-orange-800 dark:bg-orange-900/20 dark:text-orange-400'
    case 'archived':
      return 'bg-gray-100 text-gray-800 dark:bg-gray-900/20 dark:text-gray-400'
    default:
      return 'bg-gray-100 text-gray-800 dark:bg-gray-900/20 dark:text-gray-400'
  }
})

const estimateTokens = (text: string): number => {
  // Rough estimation: ~4 characters per token
  return Math.ceil(text.length / 4)
}

const configEntries = computed(() => {
  if (!localContent.value.config) return []
  return Object.entries(localContent.value.config)
    .filter(([_, value]) => value !== undefined && value !== null && value !== '')
})
</script>

<template>
  <UiCard class="flex-1 flex flex-col">
    <!-- Card Header -->
    <UiCardHeader>
      <div class="flex items-center justify-between">
        <div class="flex items-center gap-3">
          <FileText class="h-5 w-5 text-primary" />
          <div>
            <UiCardTitle class="flex items-center gap-2">
              Prompt Content
              <UiBadge 
                v-if="status" 
                :class="getStatusBadgeClass"
                size="sm"
              >
                {{ status }}
              </UiBadge>
            </UiCardTitle>
            <UiCardDescription>View and edit prompt content across tabs</UiCardDescription>
          </div>
        </div>

        <!-- Action Buttons -->
        <div class="flex items-center gap-2" v-if="showPlaygroundButton">
          <UiButton 
            variant="outline" 
            size="sm" 
            class="gap-2"
            @click="emit('openPlayground')"
          >
            <Play class="h-3 w-3" />
            Playground
          </UiButton>
        </div>
      </div>
    </UiCardHeader>

    <!-- Tabs Navigation -->
    <div class="border-b bg-muted/20">
      <div class="flex">
        <!-- User Prompt Tab -->
        <button
          @click="activeTab = 'user'"
          :class="[
            'flex-1 px-4 py-3 text-sm font-medium border-b-2 transition-colors flex items-center justify-center gap-2',
            activeTab === 'user'
              ? 'border-primary text-primary bg-background'
              : 'border-transparent text-muted-foreground hover:text-foreground hover:bg-background/50'
          ]"
        >
          <MessageSquare class="h-4 w-4" />
          User Prompt
          <span v-if="localContent.userPrompt?.trim()" class="h-2 w-2 bg-green-500 rounded-full"></span>
          <UiBadge 
            v-if="localContent.userPrompt?.trim()" 
            variant="secondary" 
            class="text-xs ml-1"
          >
            ~{{ estimateTokens(localContent.userPrompt).toLocaleString() }}
          </UiBadge>
        </button>

        <!-- System Prompt Tab -->
        <button
          @click="activeTab = 'system'"
          :class="[
            'flex-1 px-4 py-3 text-sm font-medium border-b-2 transition-colors flex items-center justify-center gap-2',
            activeTab === 'system'
              ? 'border-primary text-primary bg-background'
              : 'border-transparent text-muted-foreground hover:text-foreground hover:bg-background/50'
          ]"
        >
          <Settings class="h-4 w-4" />
          System Prompt
          <span v-if="localContent.systemPrompt?.trim()" class="h-2 w-2 bg-blue-500 rounded-full"></span>
          <UiBadge 
            v-if="localContent.systemPrompt?.trim()" 
            variant="secondary" 
            class="text-xs ml-1"
          >
            ~{{ estimateTokens(localContent.systemPrompt).toLocaleString() }}
          </UiBadge>
        </button>

        <!-- Config Tab -->
        <button
          @click="activeTab = 'config'"
          :class="[
            'flex-1 px-4 py-3 text-sm font-medium border-b-2 transition-colors flex items-center justify-center gap-2',
            activeTab === 'config'
              ? 'border-primary text-primary bg-background'
              : 'border-transparent text-muted-foreground hover:text-foreground hover:bg-background/50'
          ]"
        >
          <Code class="h-4 w-4" />
          Configuration
          <span v-if="configEntries.length > 0" class="h-2 w-2 bg-orange-500 rounded-full"></span>
          <UiBadge 
            v-if="configEntries.length > 0" 
            variant="secondary" 
            class="text-xs ml-1"
          >
            {{ configEntries.length }}
          </UiBadge>
        </button>
      </div>
    </div>

    <!-- Tab Content -->
    <UiCardContent class="flex-1 p-0 min-h-0">
      <!-- User Prompt Content -->
      <div v-show="activeTab === 'user'" class="h-full relative">
        <UiTextarea
          v-if="isEditable"
          v-model="localContent.userPrompt"
          @input="updateContent"
          placeholder="Enter your user prompt here... This is the main instruction for the AI assistant."
          class="absolute inset-0 w-full h-full resize-none border-0 focus-visible:ring-0 font-mono text-sm leading-relaxed p-4 bg-background"
        />
        <div 
          v-else
          class="absolute inset-0 w-full h-full p-4 overflow-auto"
        >
          <pre 
            v-if="localContent.userPrompt?.trim()" 
            class="whitespace-pre-wrap text-sm bg-muted p-4 rounded-lg border font-mono leading-relaxed"
          >{{ localContent.userPrompt }}</pre>
          <div v-else class="flex items-center justify-center h-full text-muted-foreground">
            <div class="text-center">
              <MessageSquare class="h-8 w-8 mx-auto mb-2 opacity-50" />
              <p class="text-sm">No user prompt configured</p>
            </div>
          </div>
        </div>

        <!-- Action buttons overlay -->
        <div class="absolute top-3 right-3 flex items-center gap-2">
          <UiButton 
            v-if="localContent.userPrompt?.trim()"
            variant="ghost" 
            size="sm" 
            class="gap-2 bg-background/90 backdrop-blur-sm"
            @click="copyToClipboard('user')"
          >
            <Copy class="h-3 w-3" />
            Copy
          </UiButton>
          <UiButton 
            v-if="isEditable && localContent.userPrompt?.trim()"
            variant="ghost" 
            size="sm" 
            class="gap-2 bg-background/90 backdrop-blur-sm"
            @click="clearContent('user')"
          >
            <RotateCcw class="h-3 w-3" />
            Clear
          </UiButton>
        </div>
      </div>

      <!-- System Prompt Content -->
      <div v-show="activeTab === 'system'" class="h-full relative">
        <UiTextarea
          v-if="isEditable"
          v-model="localContent.systemPrompt"
          @input="updateContent"
          placeholder="Define the AI's role, behavior, and constraints here... This sets the context for how the AI should respond."
          class="absolute inset-0 w-full h-full resize-none border-0 focus-visible:ring-0 font-mono text-sm leading-relaxed p-4 bg-background"
        />
        <div 
          v-else
          class="absolute inset-0 w-full h-full p-4 overflow-auto"
        >
          <pre 
            v-if="localContent.systemPrompt?.trim()" 
            class="whitespace-pre-wrap text-sm bg-muted p-4 rounded-lg border font-mono leading-relaxed"
          >{{ localContent.systemPrompt }}</pre>
          <div v-else class="flex items-center justify-center h-full text-muted-foreground">
            <div class="text-center">
              <Settings class="h-8 w-8 mx-auto mb-2 opacity-50" />
              <p class="text-sm">No system prompt configured</p>
            </div>
          </div>
        </div>

        <!-- Action buttons overlay -->
        <div class="absolute top-3 right-3 flex items-center gap-2">
          <UiButton 
            v-if="localContent.systemPrompt?.trim()"
            variant="ghost" 
            size="sm" 
            class="gap-2 bg-background/90 backdrop-blur-sm"
            @click="copyToClipboard('system')"
          >
            <Copy class="h-3 w-3" />
            Copy
          </UiButton>
          <UiButton 
            v-if="isEditable && localContent.systemPrompt?.trim()"
            variant="ghost" 
            size="sm" 
            class="gap-2 bg-background/90 backdrop-blur-sm"
            @click="clearContent('system')"
          >
            <RotateCcw class="h-3 w-3" />
            Clear
          </UiButton>
        </div>
      </div>

      <!-- Configuration Content -->
      <div v-show="activeTab === 'config'" class="h-full p-4 overflow-auto">
        <div v-if="configEntries.length > 0" class="space-y-4">
          <div class="grid grid-cols-1 md:grid-cols-2 gap-4">
            <div
              v-for="[key, value] in configEntries"
              :key="key"
              class="p-3 border rounded-lg bg-muted/20"
            >
              <div class="flex items-center justify-between mb-1">
                <label class="text-sm font-medium text-foreground capitalize">
                  {{ key.replace(/([A-Z])/g, ' $1').toLowerCase() }}
                </label>
                <UiBadge variant="outline" size="sm">
                  {{ typeof value }}
                </UiBadge>
              </div>
              <div class="text-sm text-muted-foreground font-mono">
                {{ typeof value === 'object' ? JSON.stringify(value) : String(value) }}
              </div>
            </div>
          </div>

          <div class="pt-2 border-t">
            <UiButton 
              variant="outline" 
              size="sm" 
              class="gap-2"
              @click="copyToClipboard('config')"
            >
              <Copy class="h-3 w-3" />
              Copy Configuration JSON
            </UiButton>
          </div>
        </div>
        
        <div v-else class="flex items-center justify-center h-full text-muted-foreground">
          <div class="text-center">
            <Code class="h-8 w-8 mx-auto mb-2 opacity-50" />
            <p class="text-sm">No configuration parameters set</p>
          </div>
        </div>
      </div>
    </UiCardContent>
  </UiCard>
</template>