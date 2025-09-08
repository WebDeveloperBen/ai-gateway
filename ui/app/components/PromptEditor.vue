<script setup lang="ts">
import { FileText, Copy, Save, Settings, MessageSquare, RotateCcw, Clock, Cpu, DollarSign, Zap } from "lucide-vue-next"

interface Props {
  showParametersModal: () => void
}

const props = defineProps<Props>()

// Use the playground state
const playgroundState = usePlaygroundState()
if (!playgroundState) {
  throw new Error("PromptEditor must be used within a playground state provider")
}

const {
  // Basic state
  activePromptTab,
  promptText,
  systemPrompt,
  isLoading,

  // Reactive state
  promptState,
  modelState,

  // Template refs
  promptTextarea,
  systemPromptTextarea,

  // Functions
  insertTemplate,
  actions,
  utils
} = playgroundState

// Template tab switching helper
const switchToTab = (tab: "system" | "user") => {
  activePromptTab.value = tab
}
</script>

<template>
  <div class="flex-1 flex flex-col rounded-lg border bg-background shadow-sm">
    <!-- Card Header -->
    <div class="p-4 border-b">
      <div class="flex items-center justify-between">
        <div class="flex items-center gap-3">
          <FileText class="h-5 w-5 text-primary" />
          <div>
            <h3 class="font-semibold text-base">Prompt Editor</h3>
            <div v-if="promptState.currentPrompt" class="flex items-center gap-2 text-xs text-muted-foreground mt-0.5">
              <span>{{ promptState.currentPrompt.name }}</span>
              <span>•</span>
              <span
                v-if="promptState.currentVersion"
                class="inline-flex items-center px-2 py-0.5 rounded-full bg-blue-100 text-blue-700 dark:bg-blue-900/20 dark:text-blue-400"
              >
                {{ promptState.currentVersion.version }}
              </span>
              <span
                v-else
                class="inline-flex items-center px-2 py-0.5 rounded-full bg-orange-100 text-orange-700 dark:bg-orange-900/20 dark:text-orange-400"
              >
                Draft
              </span>
            </div>
            <div v-else class="text-xs text-muted-foreground mt-0.5">New unsaved prompt</div>
          </div>
        </div>

        <!-- Controls -->
        <div class="flex items-center gap-3">
          <!-- Version Selector -->
          <div
            v-if="promptState.currentPrompt && promptState.currentPrompt.versions.length > 0"
            class="flex items-center gap-2"
          >
            <label class="text-xs text-muted-foreground">Version:</label>
            <select
              v-model="promptState.selectedVersionId"
              @change="actions.loadVersion"
              class="text-xs border border-border rounded px-2 py-1 bg-background min-w-24"
            >
              <option value="">Draft</option>
              <option
                v-for="version in promptState.currentPrompt.versions.slice().reverse()"
                :key="version.id"
                :value="version.id"
              >
                {{ version.version }}
              </option>
            </select>
          </div>

          <div class="h-6 border-l border-border"></div>

          <button
            class="inline-flex items-center gap-2 px-3 py-1.5 text-sm font-medium rounded-lg border border-border bg-background hover:bg-muted/50 transition-all duration-150 ease-in-out hover:shadow-sm disabled:opacity-50 disabled:cursor-not-allowed"
            type="button"
            @click="actions.copyPrompt"
            :disabled="!promptText && !systemPrompt"
            title="Copy current content"
          >
            <Copy class="h-3.5 w-3.5" />
            <span class="hidden sm:inline">Copy</span>
          </button>

          <button
            class="inline-flex items-center gap-2 px-3 py-1.5 text-sm font-medium rounded-lg bg-primary text-primary-foreground hover:bg-primary/90 transition-all duration-150 ease-in-out hover:shadow-sm disabled:opacity-50 disabled:cursor-not-allowed"
            type="button"
            @click="actions.savePrompt"
            :disabled="isLoading || (!promptText.trim() && !systemPrompt.trim())"
            title="Save as new version"
          >
            <Save class="h-3.5 w-3.5" />
            <span class="hidden sm:inline">Save</span>
          </button>

          <button
            @click="props.showParametersModal"
            class="inline-flex items-center gap-2 px-3 py-1.5 text-sm font-medium rounded-lg border border-border bg-background hover:bg-muted/50 transition-all duration-150 ease-in-out hover:shadow-sm"
            title="Model Parameters"
            type="button"
          >
            <Settings class="h-3.5 w-3.5" />
            <span class="hidden sm:inline">Parameters</span>
          </button>

          <div class="flex items-center gap-3 ml-3 px-3 py-2 rounded-lg bg-muted/30">
            <div
              class="h-2.5 w-2.5 rounded-full transition-colors"
              :class="{
                'bg-green-500 animate-pulse': modelState.testResults.length > 0 && modelState.testResults[0]?.success,
                'bg-red-500 animate-pulse': modelState.testResults.length > 0 && !modelState.testResults[0]?.success,
                'bg-yellow-500 animate-pulse': isLoading,
                'bg-gray-400': modelState.testResults.length === 0 && !isLoading
              }"
            />
            <span class="text-sm font-medium text-foreground">
              {{
                isLoading
                  ? "Testing..."
                  : modelState.testResults.length > 0
                    ? modelState.testResults[0]?.success
                      ? "Success"
                      : "Error"
                    : "Ready"
              }}
            </span>
          </div>
        </div>
      </div>
    </div>

    <!-- Prompt Tabs and Content -->
    <div class="flex flex-1 min-h-0">
      <!-- Left Side: Tabbed Prompt Editor -->
      <div class="flex-1 flex flex-col border-r">
        <!-- Prompt Type Tabs -->
        <div class="border-b bg-muted/20">
          <div class="flex">
            <button
              @click="switchToTab('system')"
              :class="[
                'px-4 py-3 text-sm font-medium border-b-2 transition-colors flex items-center gap-2',
                activePromptTab === 'system'
                  ? 'border-primary text-primary bg-background'
                  : 'border-transparent text-muted-foreground hover:text-foreground hover:bg-background/50'
              ]"
            >
              <Settings class="h-4 w-4" />
              System Prompt
              <span v-if="systemPrompt.trim()" class="h-2 w-2 bg-blue-500 rounded-full"></span>
            </button>
            <button
              @click="switchToTab('user')"
              :class="[
                'px-4 py-3 text-sm font-medium border-b-2 transition-colors flex items-center gap-2',
                activePromptTab === 'user'
                  ? 'border-primary text-primary bg-background'
                  : 'border-transparent text-muted-foreground hover:text-foreground hover:bg-background/50'
              ]"
            >
              <MessageSquare class="h-4 w-4" />
              User Prompt
              <span v-if="promptText.trim()" class="h-2 w-2 bg-green-500 rounded-full"></span>
            </button>
          </div>
        </div>

        <!-- System Prompt Editor -->
        <div v-show="activePromptTab === 'system'" class="flex-1 relative min-h-0">
          <UiTextarea
            ref="systemPromptTextarea"
            v-model="systemPrompt"
            placeholder="Define the AI's role, behavior, and constraints here... This sets the context for how the AI should respond."
            class="absolute inset-0 w-full h-full resize-none border-0 focus-visible:ring-0 font-mono text-sm leading-relaxed p-4 bg-background"
            :disabled="isLoading"
          />

          <!-- Character/token count overlay -->
          <div
            class="absolute bottom-3 right-3 flex items-center gap-2 rounded-md bg-background/90 backdrop-blur-sm px-3 py-1.5 text-xs text-muted-foreground border shadow-sm"
          >
            <span>{{ systemPrompt.length.toLocaleString() }} chars</span>
            <span>•</span>
            <span>~{{ utils.estimateTokens(systemPrompt).toLocaleString() }} tokens</span>
          </div>

          <!-- Clear System Prompt Button -->
          <button
            v-if="systemPrompt.trim()"
            @click="actions.clearSystemPrompt"
            class="absolute top-3 right-3 inline-flex items-center gap-1 px-2 py-1 text-xs font-medium rounded border border-border bg-background/90 hover:bg-muted/50 transition-all duration-150 ease-in-out"
            title="Clear system prompt"
          >
            <RotateCcw class="h-3 w-3" />
            Clear
          </button>
        </div>

        <!-- User Prompt Editor -->
        <div v-show="activePromptTab === 'user'" class="flex-1 relative min-h-0">
          <UiTextarea
            ref="promptTextarea"
            v-model="promptText"
            placeholder="Enter your user prompt here... Click on templates from the right sidebar to insert them at your cursor position."
            class="absolute inset-0 w-full h-full resize-none border-0 focus-visible:ring-0 font-mono text-sm leading-relaxed p-4 bg-background"
            :disabled="isLoading"
          />

          <!-- Character/token count overlay -->
          <div
            class="absolute bottom-3 right-3 flex items-center gap-2 rounded-md bg-background/90 backdrop-blur-sm px-3 py-1.5 text-xs text-muted-foreground border shadow-sm"
          >
            <span>{{ promptText.length.toLocaleString() }} chars</span>
            <span>•</span>
            <span>~{{ utils.estimateTokens(promptText).toLocaleString() }} tokens</span>
            <span
              v-if="
                modelState.selectedModelData &&
                utils.estimateTokens(promptText) + utils.estimateTokens(systemPrompt) >
                  modelState.selectedModelData.maxTokens
              "
              class="text-red-500 font-medium"
            >
              (total exceeds {{ modelState.selectedModelData.maxTokens.toLocaleString() }} limit)
            </span>
          </div>

          <!-- Clear User Prompt Button -->
          <button
            v-if="promptText.trim()"
            @click="actions.clearPrompt"
            class="absolute top-3 right-3 inline-flex items-center gap-1 px-2 py-1 text-xs font-medium rounded border border-border bg-background/90 hover:bg-muted/50 transition-all duration-150 ease-in-out"
            title="Clear user prompt"
          >
            <RotateCcw class="h-3 w-3" />
            Clear
          </button>
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
            <div v-if="modelState.testResults.length > 0" class="p-4">
              <div class="space-y-4">
                <div
                  v-for="result in modelState.testResults"
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
                          {{ utils.formatTime(result.timestamp) }}
                        </div>
                      </div>

                      <div class="flex items-center gap-2">
                        <div class="text-xs text-muted-foreground">
                          {{ result.response.length.toLocaleString() }} chars
                        </div>
                        <button
                          class="inline-flex items-center gap-1 px-2 py-1 text-xs font-medium rounded border border-border bg-background hover:bg-muted/50 transition-all duration-150 ease-in-out"
                          type="button"
                          @click="actions.copyResponse(result.response)"
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
                        <span class="text-muted-foreground">{{ result.time }}ms</span>
                      </div>

                      <div class="flex items-center gap-1">
                        <Zap class="h-3 w-3 text-purple-600 dark:text-purple-400" />
                        <span class="text-muted-foreground">{{ result.tokensUsed.total.toLocaleString() }} tokens</span>
                      </div>

                      <div class="flex items-center gap-1">
                        <DollarSign class="h-3 w-3 text-green-600 dark:text-green-400" />
                        <span class="text-muted-foreground">{{ utils.formatCurrency(result.estimatedCost) }}</span>
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
            <BlocksEmptyState title="No responses yet" subtext="Run a test to see AI responses here" v-else />
          </UiScrollArea>
        </div>
      </div>
    </div>
  </div>
</template>

