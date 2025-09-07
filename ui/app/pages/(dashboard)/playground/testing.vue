<script lang="ts">
import {
  Play,
  Clock,
  DollarSign,
  Cpu,
  FileText,
  Copy,
  RotateCcw,
  Library,
  MessageSquare,
  Code,
  PenTool,
  BarChart,
  Shield,
  Users,
  Zap,
  Loader2
} from "lucide-vue-next"
import type { FormBuilder } from "~/components/Ui/FormBuilder/FormBuilder.vue"

interface Model {
  id: string
  name: string
  provider: string
  description: string
  maxTokens: number
  costPer1kTokens: {
    input: number
    output: number
  }
}

interface PromptTemplate {
  id: string
  name: string
  content: string
  category: string
  description: string
  icon: any
}

interface TestResult {
  id: string
  timestamp: string
  prompt: string
  response: string
  model: string
  tokensUsed: {
    input: number
    output: number
    total: number
  }
  responseTime: number
  estimatedCost: number
  success: boolean
  error?: string
}

// Sample data - would come from API
const availableModels = ref<Model[]>([
  {
    id: "gpt-4",
    name: "GPT-4",
    provider: "OpenAI",
    description: "Most capable GPT-4 model",
    maxTokens: 8192,
    costPer1kTokens: { input: 0.03, output: 0.06 }
  },
  {
    id: "gpt-3.5-turbo",
    name: "GPT-3.5 Turbo",
    provider: "OpenAI",
    description: "Fast and efficient model",
    maxTokens: 4096,
    costPer1kTokens: { input: 0.0015, output: 0.002 }
  },
  {
    id: "azure-gpt-4",
    name: "GPT-4 (Azure)",
    provider: "Azure OpenAI",
    description: "GPT-4 via Azure OpenAI",
    maxTokens: 8192,
    costPer1kTokens: { input: 0.03, output: 0.06 }
  }
])

const promptTemplates = ref<PromptTemplate[]>([
  {
    id: "customer-support",
    name: "Customer Support",
    content:
      "You are a helpful customer support assistant. Always be polite and professional when helping customers with their inquiries.",
    category: "Support",
    description: "Professional customer service assistant",
    icon: MessageSquare
  },

  {
    id: "customer-support",
    name: "Customer Support",
    content:
      "You are a helpful customer support assistant. Always be polite and professional when helping customers with their inquiries.",
    category: "Support",
    description: "Professional customer service assistant",
    icon: MessageSquare
  },
  {
    id: "customer-support",
    name: "Customer Support",
    content:
      "You are a helpful customer support assistant. Always be polite and professional when helping customers with their inquiries.",
    category: "Support",
    description: "Professional customer service assistant",
    icon: MessageSquare
  },
  {
    id: "customer-support",
    name: "Customer Support",
    content:
      "You are a helpful customer support assistant. Always be polite and professional when helping customers with their inquiries.",
    category: "Support",
    description: "Professional customer service assistant",
    icon: MessageSquare
  },
  {
    id: "code-review",
    name: "Code Review",
    content: "Review the following code for best practices, security issues, and optimization opportunities:",
    category: "Development",
    description: "Technical code review and analysis",
    icon: Code
  },
  {
    id: "content-writer",
    name: "Content Writer",
    content: "Write engaging, informative content that is well-structured and matches the target audience:",
    category: "Writing",
    description: "Creative content generation",
    icon: PenTool
  },
  {
    id: "data-analyst",
    name: "Data Analyst",
    content: "Analyze the following data and provide actionable insights:",
    category: "Analytics",
    description: "Data analysis and insights",
    icon: BarChart
  },
  {
    id: "security-expert",
    name: "Security Expert",
    content: "As a cybersecurity expert, review and provide security recommendations for:",
    category: "Security",
    description: "Security analysis and recommendations",
    icon: Shield
  },
  {
    id: "hr-assistant",
    name: "HR Assistant",
    content: "As an HR professional, help with employee-related questions and policies:",
    category: "Human Resources",
    description: "HR support and guidance",
    icon: Users
  }
])
</script>
<script setup lang="ts">
useSeoMeta({ title: "Model Testing - LLM Gateway" })

// State management
const selectedModel = ref<string>("")
const promptText = ref<string>("")
const isLoading = ref(false)
const testResults = ref<TestResult[]>([])

// Textarea ref for cursor position
const promptTextarea = useTemplateRef<HTMLTextAreaElement>("promptTextarea")

// Model selection
const selectedModelData = computed(() => {
  return availableModels.value.find((m) => m.id === selectedModel.value)
})

// Form fields for model configuration
const formFields: FormBuilder[] = [
  {
    variant: "Select",
    label: "Model",
    name: "model",
    placeholder: "Select a model...",
    required: true,
    options: availableModels.value.map((model) => ({
      value: model.id,
      label: `${model.name} (${model.provider})`,
      description: model.description
    })),
    wrapperClass: "col-span-full"
  }
]

// Functions
const insertTemplate = (template: PromptTemplate) => {
  if (!promptTextarea.value) return

  const textarea = promptTextarea.value
  const start = textarea.selectionStart
  const end = textarea.selectionEnd
  const currentText = promptText.value

  // Insert template at cursor position
  const newText = currentText.substring(0, start) + template.content + currentText.substring(end)
  promptText.value = newText

  // Focus back to textarea and set cursor position
  nextTick(() => {
    textarea.focus()
    const newCursorPosition = start + template.content.length
    textarea.setSelectionRange(newCursorPosition, newCursorPosition)
  })
}

const runTest = async () => {
  if (!selectedModel.value || !promptText.value.trim()) {
    return
  }

  isLoading.value = true

  try {
    const startTime = Date.now()

    // TODO: Replace with actual API call
    // const response = await $fetch('/api/playground/test', {
    //   method: 'POST',
    //   body: {
    //     model: selectedModel.value,
    //     prompt: promptText.value,
    //     maxTokens: 1000
    //   }
    // })

    // Simulate API response
    await new Promise((resolve) => setTimeout(resolve, 2000))
    const endTime = Date.now()

    const mockResponse =
      "This is a sample response from the AI model. In a real implementation, this would be the actual response from your selected model."
    const inputTokens = estimateTokens(promptText.value)
    const outputTokens = estimateTokens(mockResponse)
    const totalTokens = inputTokens + outputTokens

    const modelData = selectedModelData.value!
    const estimatedCost =
      (inputTokens * modelData.costPer1kTokens.input + outputTokens * modelData.costPer1kTokens.output) / 1000

    const result: TestResult = {
      id: `test_${Date.now()}`,
      timestamp: new Date().toISOString(),
      prompt: promptText.value,
      response: mockResponse,
      model: modelData.name,
      tokensUsed: {
        input: inputTokens,
        output: outputTokens,
        total: totalTokens
      },
      responseTime: endTime - startTime,
      estimatedCost,
      success: true
    }

    testResults.value.unshift(result)
  } catch (error) {
    console.error("Test failed:", error)
    const result: TestResult = {
      id: `test_${Date.now()}`,
      timestamp: new Date().toISOString(),
      prompt: promptText.value,
      response: "",
      model: selectedModelData.value?.name || "Unknown",
      tokensUsed: { input: 0, output: 0, total: 0 },
      responseTime: 0,
      estimatedCost: 0,
      success: false,
      error: "Failed to get response from model"
    }

    testResults.value.unshift(result)
  } finally {
    isLoading.value = false
  }
}

const clearPrompt = () => {
  promptText.value = ""
  promptTextarea.value?.focus()
}

const copyResponse = async (response: string) => {
  try {
    await navigator.clipboard.writeText(response)
    // TODO: Show success toast
  } catch (err) {
    console.error("Failed to copy:", err)
  }
}

const formatCurrency = (amount: number) => {
  return new Intl.NumberFormat("en-US", {
    style: "currency",
    currency: "USD",
    minimumFractionDigits: 6
  }).format(amount)
}

const copyPrompt = async () => {
  try {
    await navigator.clipboard.writeText(promptText.value)
    // TODO: Show success toast
    console.log("Prompt copied to clipboard")
  } catch (err) {
    console.error("Failed to copy prompt:", err)
  }
}

const formatTime = (dateString: string) => {
  return new Date(dateString).toLocaleTimeString("en-US", {
    hour: "2-digit",
    minute: "2-digit"
  })
}

// More accurate token estimation function
const estimateTokens = (text: string): number => {
  if (!text) return 0

  // More sophisticated token estimation
  // Account for spaces, punctuation, and typical tokenization patterns
  const totalChars = text.length

  // GPT models typically have these token ratios:
  // - ~4 characters per token for English text
  // - Punctuation and special characters can be separate tokens
  // - Code and structured text can have different ratios

  let tokenEstimate = 0

  // Base estimation: ~4 chars per token
  tokenEstimate += totalChars / 4

  // Add extra tokens for punctuation
  const punctuation = (text.match(/[.,;:!?(){}[\]"']/g) || []).length
  tokenEstimate += punctuation * 0.5

  // Add extra tokens for line breaks and formatting
  const lineBreaks = (text.match(/\n/g) || []).length
  tokenEstimate += lineBreaks * 0.3

  // Add extra tokens for numbers and special characters
  const numbers = (text.match(/\d+/g) || []).length
  tokenEstimate += numbers * 0.2

  return Math.ceil(tokenEstimate)
}

// Set default model
onMounted(() => {
  if (availableModels.value.length > 0) {
    selectedModel.value = availableModels.value[0]?.id ?? ""
  }
})
</script>

<template>
  <div class="flex flex-col gap-6 h-screen overflow-hidden">
    <!-- Header -->
    <PageHeader title="Model Testing Playground" subtext="Test prompts against your registered models and providers" />

    <!-- Configuration Form -->
    <div class="flex gap-5 items-end flex-shrink-0">
      <div class="flex-1">
        <form @submit.prevent>
          <fieldset :disabled="isLoading">
            <UiFormBuilder class="grid grid-cols-12 gap-5" :fields="formFields" />
          </fieldset>
        </form>
      </div>

      <!-- Test Button -->
      <button
        class="inline-flex items-center justify-center gap-2 px-6 py-2.5 text-sm font-semibold rounded-lg bg-primary text-primary-foreground hover:bg-primary/90 transition-all duration-150 ease-in-out hover:shadow-md disabled:opacity-50 disabled:cursor-not-allowed whitespace-nowrap"
        type="button"
        @click="runTest"
        :disabled="!selectedModel || !promptText.trim() || isLoading"
      >
        <Play v-if="!isLoading" class="h-4 w-4" />
        <Loader2 v-else class="h-4 w-4 animate-spin" />
        {{ isLoading ? "Testing..." : "Test Prompt" }}
      </button>
    </div>

    <!-- Main Editor Area -->
    <div class="flex gap-6 flex-1 min-h-0">
      <!-- Prompt Editor and Chat Responses Card -->
      <div class="flex-1 flex flex-col rounded-lg border bg-background shadow-sm">
        <!-- Card Header -->
        <div class="p-4 border-b">
          <div class="flex items-center justify-between">
            <div class="flex items-center gap-2">
              <FileText class="h-5 w-5 text-primary" />
              <h3 class="font-semibold text-base">Testing Playground</h3>
            </div>

            <!-- Controls -->
            <div class="flex items-center gap-3">
              <button
                class="inline-flex items-center gap-2 px-4 py-2 text-sm font-medium rounded-lg bg-secondary text-secondary-foreground hover:bg-secondary/80 border border-border/50 transition-all duration-150 ease-in-out hover:shadow-sm hover:border-border disabled:opacity-50 disabled:cursor-not-allowed"
                type="button"
                @click="clearPrompt"
                :disabled="isLoading || !promptText"
                title="Clear prompt"
              >
                <RotateCcw class="h-4 w-4" />
                <span>Clear</span>
              </button>

              <button
                class="inline-flex items-center gap-2 px-4 py-2 text-sm font-medium rounded-lg border border-border bg-background hover:bg-muted/50 transition-all duration-150 ease-in-out hover:shadow-sm hover:border-border/60 disabled:opacity-50 disabled:cursor-not-allowed"
                type="button"
                @click="copyPrompt"
                :disabled="!promptText"
                title="Copy prompt to clipboard"
              >
                <Copy class="h-4 w-4" />
                <span>Copy</span>
              </button>

              <div class="flex items-center gap-3 ml-3 px-3 py-2 rounded-lg bg-muted/30">
                <div
                  class="h-2.5 w-2.5 rounded-full transition-colors"
                  :class="{
                    'bg-green-500 animate-pulse': testResults.length > 0 && testResults[0]?.success,
                    'bg-red-500 animate-pulse': testResults.length > 0 && !testResults[0]?.success,
                    'bg-yellow-500 animate-pulse': isLoading,
                    'bg-gray-400': testResults.length === 0 && !isLoading
                  }"
                />
                <span class="text-sm font-medium text-foreground">
                  {{
                    isLoading
                      ? "Testing..."
                      : testResults.length > 0
                        ? testResults[0]?.success
                          ? "Success"
                          : "Error"
                        : "Ready"
                  }}
                </span>
              </div>
            </div>
          </div>
        </div>

        <!-- Side by Side Content -->
        <div class="flex flex-1 min-h-0">
          <!-- Left Side: Prompt Editor -->
          <div class="flex-1 flex flex-col border-r">
            <div class="p-3 border-b bg-muted/20">
              <h4 class="text-xs font-medium text-muted-foreground">Prompt Editor</h4>
            </div>

            <div class="flex-1 relative min-h-0">
              <UiTextarea
                ref="promptTextarea"
                v-model="promptText"
                placeholder="Enter your prompt here... Click on templates from the right sidebar to insert them at your cursor position."
                class="absolute inset-0 w-full h-full resize-none border-0 focus-visible:ring-0 font-mono text-sm leading-relaxed p-4 bg-transparent"
                :disabled="isLoading"
              />

              <!-- Character/token count overlay -->
              <div
                class="absolute bottom-3 right-3 flex items-center gap-2 rounded-md bg-background/90 backdrop-blur-sm px-3 py-1.5 text-xs text-muted-foreground border shadow-sm"
              >
                <span>{{ promptText.length.toLocaleString() }} chars</span>
                <span>•</span>
                <span>~{{ estimateTokens(promptText).toLocaleString() }} tokens</span>
                <span
                  v-if="selectedModelData && estimateTokens(promptText) > selectedModelData.maxTokens"
                  class="text-red-500 font-medium"
                >
                  (exceeds {{ selectedModelData.maxTokens.toLocaleString() }} limit)
                </span>
              </div>
            </div>
          </div>

          <!-- Right Side: Chat Responses -->
          <div class="flex-1 flex flex-col min-h-0">
            <div class="p-3 border-b bg-muted/20">
              <h4 class="text-xs font-medium text-muted-foreground">AI Responses</h4>
            </div>

            <div class="flex-1 min-h-0">
              <UiScrollArea class="h-full">
                <!-- Test Results -->
                <div v-if="testResults.length > 0" class="p-4">
                  <div class="space-y-4">
                    <div
                      v-for="result in testResults"
                      :key="result.id"
                      class="border border-border rounded-lg overflow-hidden bg-card shadow-sm"
                    >
                      <!-- Header with Inline Stats -->
                      <div class="p-3 border-b bg-card">
                        <div class="flex items-center justify-between mb-2">
                          <div class="flex items-center gap-2">
                            <Cpu class="h-3 w-3 text-primary" />
                            <h5 class="font-medium text-xs">{{ result.model }}</h5>
                            <div class="text-xs text-muted-foreground">
                              {{ formatTime(result.timestamp) }}
                            </div>
                          </div>

                          <div class="flex items-center gap-2">
                            <div class="text-xs text-muted-foreground">
                              {{ result.response.length.toLocaleString() }} chars
                            </div>
                            <button
                              class="inline-flex items-center gap-1 px-2 py-1 text-xs font-medium rounded border border-border bg-background hover:bg-muted/50 transition-all duration-150 ease-in-out"
                              type="button"
                              @click="copyResponse(result.response)"
                              title="Copy response"
                            >
                              <Copy class="h-3 w-3" />
                            </button>
                          </div>
                        </div>

                        <!-- Compact Inline Stats -->
                        <div v-if="result.success" class="flex items-center gap-4 text-xs">
                          <div class="flex items-center gap-1">
                            <Clock class="h-3 w-3 text-blue-600 dark:text-blue-400" />
                            <span class="text-muted-foreground">{{ result.responseTime }}ms</span>
                          </div>

                          <div class="flex items-center gap-1">
                            <Zap class="h-3 w-3 text-purple-600 dark:text-purple-400" />
                            <span class="text-muted-foreground"
                              >{{ result.tokensUsed.total.toLocaleString() }} tokens</span
                            >
                          </div>

                          <div class="flex items-center gap-1">
                            <DollarSign class="h-3 w-3 text-green-600 dark:text-green-400" />
                            <span class="text-muted-foreground">{{ formatCurrency(result.estimatedCost) }}</span>
                          </div>
                        </div>
                      </div>

                      <!-- Response Content -->
                      <div v-if="result.success" class="bg-background">
                        <UiScrollArea class="max-h-48">
                          <div
                            class="p-3 text-sm font-mono leading-relaxed whitespace-pre-wrap border-l-4 border-l-primary/20 bg-muted/20"
                          >
                            {{ result.response }}
                          </div>
                        </UiScrollArea>
                      </div>

                      <!-- Error State -->
                      <div v-else class="p-3 bg-red-50 dark:bg-red-950/20">
                        <div class="flex items-center gap-2 mb-2">
                          <span class="text-red-600 dark:text-red-400 font-bold text-xs">✗</span>
                          <div class="font-medium text-red-900 dark:text-red-100 text-xs">Request Failed</div>
                        </div>
                        <div
                          class="bg-red-100 dark:bg-red-900/30 border border-red-200 dark:border-red-700 rounded p-2 text-xs text-red-800 dark:text-red-200"
                        >
                          {{ result.error }}
                        </div>
                      </div>
                    </div>
                  </div>
                </div>

                <!-- Empty State -->
                <div
                  v-else
                  class="flex flex-col items-center justify-center flex-1 text-center text-muted-foreground p-8"
                >
                  <MessageSquare class="h-12 w-12 mb-3 opacity-50" />
                  <p class="text-sm mb-1">No responses yet</p>
                  <p class="text-xs opacity-70">Run a test to see AI responses here</p>
                </div>
              </UiScrollArea>
            </div>
          </div>
        </div>
      </div>

      <!-- Template Sidebar -->
      <div class="w-80 flex flex-col rounded-xl border bg-card shadow-lg overflow-hidden">
        <!-- Templates Section -->
        <div class="p-4 border-b bg-card">
          <h3 class="font-semibold text-sm flex items-center gap-2 text-foreground">
            <Library class="h-4 w-4 text-primary" />
            Prompt Templates
          </h3>
          <p class="text-xs text-muted-foreground mt-1">Click to insert at cursor position</p>
        </div>

        <div class="flex-1 min-h-0">
          <UiScrollArea class="h-full">
            <div class="p-4">
              <div class="grid grid-cols-2 gap-3">
                <div
                  v-for="template in promptTemplates"
                  :key="template.id"
                  class="flex flex-col gap-3 p-4 rounded-xl border border-border bg-background hover:bg-primary/5 hover:border-primary/30 cursor-pointer transition-all duration-200 ease-in-out group hover:shadow-md hover:scale-[1.02]"
                  @click="insertTemplate(template)"
                  :title="template.description"
                >
                  <div
                    class="flex items-center justify-center h-10 w-10 rounded-lg bg-primary/10 group-hover:bg-primary/20 transition-colors shadow-sm"
                  >
                    <component :is="template.icon" class="h-5 w-5 text-primary" />
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
    </div>

    <!-- Model Info Bar -->
    <div v-if="selectedModelData" class="pt-4 border-t bg-gradient-to-r from-muted/20 to-muted/30 p-4 rounded-lg">
      <div class="text-sm text-muted-foreground text-center">
        <span class="font-medium">{{ selectedModelData.name }}</span>
        <span class="mx-2">•</span>
        <span>Max: {{ selectedModelData.maxTokens.toLocaleString() }} tokens</span>
        <span class="mx-2">•</span>
        <span
          >Cost: {{ formatCurrency(selectedModelData.costPer1kTokens.input) }}/1k in,
          {{ formatCurrency(selectedModelData.costPer1kTokens.output) }}/1k out</span
        >
      </div>
    </div>
  </div>
</template>
