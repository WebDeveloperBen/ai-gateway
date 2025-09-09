<script setup lang="ts">
import { FileText, Save, Settings, Library, Play, BookOpen } from "lucide-vue-next"

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

interface TemplateItem {
  id: string
  name: string
  description: string
  icon: any
  category: string
  content: {
    userPrompt?: string
    systemPrompt?: string
    config?: PromptConfig
  }
}

const emit = defineEmits<{
  save: [data: { name: string; description: string; category: string; tags: string[]; content: PromptContent }]
  cancel: []
}>()

// Form state
const promptName = ref('')
const promptDescription = ref('')
const promptCategory = ref('')
const promptTags = ref<string[]>([])

// Content state
const promptContent = ref<PromptContent>({
  userPrompt: '',
  systemPrompt: '',
  config: {
    model: 'gpt-4',
    temperature: 0.7,
    maxTokens: 1000,
    topP: 0.9
  }
})

const isSaving = ref(false)
const showTemplates = ref(true)

// Sample templates - in real app this would come from API
const templates = ref<TemplateItem[]>([
  {
    id: '1',
    name: 'Customer Support',
    description: 'Empathetic customer service assistant',
    icon: FileText,
    category: 'Support',
    content: {
      userPrompt: 'Help me resolve this customer inquiry: {{customer_issue}}\n\nPlease provide a helpful and empathetic response.',
      systemPrompt: 'You are a customer support assistant. Be helpful, empathetic, and professional in all responses.',
      config: { model: 'gpt-4', temperature: 0.7, maxTokens: 800 }
    }
  },
  {
    id: '2', 
    name: 'Code Review',
    description: 'Technical code review assistant',
    icon: Settings,
    category: 'Development',
    content: {
      userPrompt: 'Please review this code and provide feedback:\n\n```\n{{code}}\n```',
      systemPrompt: 'You are an expert code reviewer. Provide constructive feedback focusing on best practices, security, and performance.',
      config: { model: 'gpt-4', temperature: 0.3, maxTokens: 1200 }
    }
  },
  {
    id: '3',
    name: 'Content Writer',
    description: 'Creative content generation',
    icon: BookOpen,
    category: 'Content',
    content: {
      userPrompt: 'Write engaging content about: {{topic}}\n\nTarget audience: {{audience}}',
      systemPrompt: 'You are a skilled content writer. Create engaging, informative, and well-structured content.',
      config: { model: 'gpt-4', temperature: 0.8, maxTokens: 1500 }
    }
  }
])

const categories = computed(() => {
  const cats = new Set(templates.value.map(t => t.category))
  return Array.from(cats)
})

const filteredTemplates = computed(() => {
  return templates.value
})

const handleContentUpdate = (content: PromptContent) => {
  promptContent.value = { ...content }
}

const handleContentCopy = (type: string, content: string) => {
  console.log(`${type} content copied:`, content)
  // TODO: Show success toast
}

const handleSave = async () => {
  if (!promptName.value.trim()) {
    // TODO: Show error toast
    return
  }

  isSaving.value = true
  try {
    emit('save', {
      name: promptName.value,
      description: promptDescription.value,
      category: promptCategory.value,
      tags: promptTags.value,
      content: promptContent.value
    })
  } catch (error) {
    console.error('Failed to save prompt:', error)
    // TODO: Show error toast
  } finally {
    isSaving.value = false
  }
}

const handleCancel = () => {
  emit('cancel')
}

const useTemplate = (template: TemplateItem) => {
  if (template.content.userPrompt) {
    promptContent.value.userPrompt = template.content.userPrompt
  }
  if (template.content.systemPrompt) {
    promptContent.value.systemPrompt = template.content.systemPrompt
  }
  if (template.content.config) {
    promptContent.value.config = { ...promptContent.value.config, ...template.content.config }
  }
  
  // Pre-fill form if not already filled
  if (!promptName.value && template.name) {
    promptName.value = template.name
  }
  if (!promptDescription.value && template.description) {
    promptDescription.value = template.description
  }
  if (!promptCategory.value && template.category) {
    promptCategory.value = template.category
  }
  
  // TODO: Show success toast
  console.log('Template applied:', template.name)
}

const handleOpenPlayground = () => {
  // TODO: Open playground with current content for testing
  console.log('Opening playground with current prompt...')
}

const isValid = computed(() => {
  return promptName.value.trim() && 
         (promptContent.value.userPrompt.trim() || promptContent.value.systemPrompt.trim())
})
</script>

<template>
  <div class="flex flex-col gap-6 h-full">
    <!-- Header -->
    <div class="flex items-center justify-between">
      <div>
        <h1 class="text-3xl font-bold text-primary tracking-tight">Create New Prompt</h1>
        <p class="text-muted-foreground">Design a new AI prompt template</p>
      </div>
      
      <div class="flex items-center gap-2">
        <UiButton 
          variant="outline" 
          @click="handleCancel"
          :disabled="isSaving"
        >
          Cancel
        </UiButton>
        <UiButton 
          variant="default" 
          class="gap-2" 
          @click="handleSave"
          :disabled="!isValid || isSaving"
        >
          <Save class="h-4 w-4" />
          {{ isSaving ? 'Saving...' : 'Save Prompt' }}
        </UiButton>
      </div>
    </div>

    <!-- Prompt Metadata Form -->
    <UiCard>
      <UiCardHeader>
        <UiCardTitle class="flex items-center gap-2">
          <Settings class="h-5 w-5" />
          Prompt Information
        </UiCardTitle>
        <UiCardDescription>Basic details about your prompt</UiCardDescription>
      </UiCardHeader>
      <UiCardContent class="space-y-4">
        <div class="grid grid-cols-1 md:grid-cols-2 gap-4">
          <div>
            <label class="text-sm font-medium mb-2 block">Prompt Name *</label>
            <UiInput
              v-model="promptName"
              placeholder="e.g. Customer Support Assistant"
              required
            />
          </div>
          <div>
            <label class="text-sm font-medium mb-2 block">Category</label>
            <UiInput
              v-model="promptCategory"
              placeholder="e.g. Customer Service, Development"
              list="categories"
            />
            <datalist id="categories">
              <option v-for="cat in categories" :key="cat" :value="cat" />
            </datalist>
          </div>
        </div>
        
        <div>
          <label class="text-sm font-medium mb-2 block">Description</label>
          <UiTextarea
            v-model="promptDescription"
            placeholder="Describe what this prompt does and when to use it..."
            rows="2"
          />
        </div>
      </UiCardContent>
    </UiCard>

    <!-- Main Content Area -->
    <div class="flex gap-6 flex-1 min-h-0">
      <!-- Template Sidebar -->
      <div v-if="showTemplates" class="w-80 rounded-xl border bg-card shadow-lg overflow-hidden">
        <div class="p-4 border-b bg-card">
          <h3 class="font-semibold text-sm flex items-center gap-2 text-foreground">
            <Library class="h-4 w-4 text-primary" />
            Prompt Templates
          </h3>
          <p class="text-xs text-muted-foreground mt-1">Start with a template or build from scratch</p>
        </div>

        <div class="flex-1 min-h-0">
          <UiScrollArea class="h-full">
            <div class="p-4">
              <div class="space-y-3">
                <div
                  v-for="template in filteredTemplates"
                  :key="template.id"
                  class="flex flex-col gap-3 p-4 rounded-xl border border-border bg-background hover:bg-primary/5 hover:border-primary/30 cursor-pointer transition-all duration-200 ease-in-out group hover:shadow-md hover:scale-[1.02]"
                  @click="useTemplate(template)"
                  :title="template.description"
                >
                  <div class="flex items-center justify-between">
                    <div
                      class="flex items-center justify-center h-10 w-10 rounded-lg bg-primary/10 group-hover:bg-primary/20 transition-colors shadow-sm"
                    >
                      <component :is="template.icon" class="h-5 w-5 text-primary" />
                    </div>
                    <UiBadge variant="outline" size="sm">
                      {{ template.category }}
                    </UiBadge>
                  </div>
                  <div>
                    <p class="text-sm font-semibold text-foreground">{{ template.name }}</p>
                    <p class="text-xs text-muted-foreground leading-tight mt-1">{{ template.description }}</p>
                  </div>
                </div>
              </div>
            </div>
          </UiScrollArea>
        </div>
      </div>

      <!-- Prompt Content Editor -->
      <div class="flex-1 min-h-96">
        <PromptsTabbedContent
          :content="promptContent"
          status="draft"
          :is-editable="true"
          :show-playground-button="true"
          @update="handleContentUpdate"
          @copy="handleContentCopy"
          @open-playground="handleOpenPlayground"
        />
      </div>
    </div>
  </div>
</template>