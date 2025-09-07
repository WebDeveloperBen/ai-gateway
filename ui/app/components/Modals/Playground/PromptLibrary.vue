<script setup lang="ts">
import { Library, Settings, Tag, History } from "lucide-vue-next"

const props = defineProps<{ savedPrompts: any[]; availableApplications: any[] }>()

// Use the playground state
const playgroundState = usePlaygroundState()
if (!playgroundState) {
  throw new Error("PromptLibrary modal must be used within a playground state provider")
}

const { modalState, modalActions, libraryState } = playgroundState

const isOpen = computed({
  get: () => modalState.showPromptLibrary,
  set: (value) => {
    if (value) {
      modalActions.openPromptLibrary()
    } else {
      modalActions.closePromptLibrary()
    }
  }
})

// These still need to be emitted since testing.vue handles the business logic
const emit = defineEmits<{
  selectPrompt: [prompt: any]
  selectVersion: [version: any]
  confirmSelection: []
}>()

const searchQuery = computed({
  get: () => libraryState.searchQuery,
  set: (value) => (libraryState.searchQuery = value)
})

const applicationFilter = computed({
  get: () => libraryState.applicationFilter,
  set: (value) => (libraryState.applicationFilter = value)
})

const formModel = computed({
  get: () => libraryState.formModel,
  set: (value) => (libraryState.formModel = value)
})

const filteredPrompts = computed(() => {
  let filtered = props.savedPrompts

  if (searchQuery.value) {
    const query = searchQuery.value.toLowerCase()
    filtered = filtered.filter(
      (prompt) =>
        prompt.name.toLowerCase().includes(query) ||
        prompt.description?.toLowerCase().includes(query) ||
        prompt.tags.some((tag: string) => tag.toLowerCase().includes(query))
    )
  }

  if (applicationFilter.value) {
    filtered = filtered.filter((prompt) => prompt.applications.includes(applicationFilter.value))
  }

  return filtered
})

const selectedVersionInLibrary = computed(() => {
  if (!libraryState.formModel.promptId || !libraryState.formModel.versionId) return null

  const prompt = props.savedPrompts.find((p) => p.id === libraryState.formModel.promptId)
  if (!prompt) return null

  return prompt.versions.find((v: any) => v.id === libraryState.formModel.versionId)
})

function selectPrompt(prompt: any) {
  emit("selectPrompt", prompt)
}

function selectVersion(version: any) {
  emit("selectVersion", version)
}

function clearFilters() {
  searchQuery.value = ""
  applicationFilter.value = ""
}

function confirmPromptSelection() {
  emit("confirmSelection")
}

function closePromptLibrary() {
  isOpen.value = false
}

function estimateTokens(text: string): number {
  return Math.ceil(text.length / 4)
}
</script>

<template>
  <UiDialog v-model:open="isOpen">
    <UiDialogContent class="min-w-6xl min-h-0 p-0 flex-1 flex-col">
      <template #header>
        <div class="border-b bg-gradient-to-r from-card to-muted/20 px-8 py-6">
          <div class="flex items-start justify-between">
            <div class="flex items-center gap-4">
              <div class="p-3 rounded-xl bg-primary/10 border">
                <Library class="h-6 w-6 text-primary" />
              </div>
              <div>
                <h3 class="font-bold text-2xl text-foreground">Prompt Library</h3>
                <p class="text-muted-foreground mt-1 text-base">Discover and load professional prompt templates</p>
              </div>
            </div>
          </div>
        </div>
      </template>

      <template #content>
        <div class="flex-1 flex min-h-0 overflow-hidden">
          <!-- Search and Filters Sidebar -->
          <div class="w-80 border-r bg-muted/10 flex flex-col">
            <!-- Search Section -->
            <div class="p-6 border-b flex-shrink-0">
              <div class="space-y-4">
                <div>
                  <label class="text-sm font-medium text-foreground mb-2 block">Search Prompts</label>
                  <div class="relative">
                    <UiInput
                      v-model="searchQuery"
                      placeholder="Search by name, description, or tags..."
                      class="w-full pl-10"
                    />
                    <div class="absolute left-3 top-1/2 transform -translate-y-1/2">
                      <span class="text-muted-foreground">🔍</span>
                    </div>
                  </div>
                </div>

                <!-- Filter Stats -->
                <div class="flex items-center gap-2 text-xs text-muted-foreground">
                  <span>{{ filteredPrompts.length }} of {{ savedPrompts.length }} prompts</span>
                  <span v-if="searchQuery || applicationFilter" class="text-primary"> (filtered) </span>
                </div>
              </div>
            </div>

            <!-- Filters Section -->
            <div class="p-6 space-y-6 flex-1 overflow-y-auto min-h-0">
              <!-- Application Filter -->
              <div>
                <label class="text-sm font-medium text-foreground mb-3 flex items-center gap-2">
                  <Settings class="h-4 w-4" />
                  Application
                </label>
                <div class="space-y-2">
                  <label class="flex items-center gap-2 cursor-pointer">
                    <input type="radio" v-model="applicationFilter" value="" class="text-primary focus:ring-primary" />
                    <span class="text-sm">All Applications</span>
                  </label>
                  <label
                    v-for="app in availableApplications"
                    :key="app.id"
                    class="flex items-center gap-2 cursor-pointer"
                  >
                    <input
                      type="radio"
                      v-model="applicationFilter"
                      :value="app.id"
                      class="text-primary focus:ring-primary"
                    />
                    <span class="text-sm">{{ app.name }}</span>
                    <span class="text-xs text-muted-foreground">
                      ({{ savedPrompts.filter((p) => p.applications.includes(app.id)).length }})
                    </span>
                  </label>
                </div>
              </div>

              <!-- Clear Filters -->
              <div v-if="searchQuery || applicationFilter" class="pt-4 border-t">
                <UiButton variant="outline" size="sm" class="w-full" @click="clearFilters">
                  Clear All Filters
                </UiButton>
              </div>
            </div>
          </div>

          <!-- Prompt Gallery -->
          <div class="flex-1 flex flex-col min-h-0 overflow-hidden">
            <UiScrollArea class="flex-1 min-h-0">
              <div class="p-6">
                <div v-if="filteredPrompts.length > 0" class="space-y-4">
                  <!-- Prompt Cards Grid -->
                  <div
                    v-for="prompt in filteredPrompts"
                    :key="prompt.id"
                    class="group relative border rounded-xl bg-card hover:shadow-lg transition-all duration-200 overflow-hidden"
                    :class="
                      formModel.promptId === prompt.id ? 'ring-2 ring-primary shadow-lg' : 'hover:border-primary/30'
                    "
                  >
                    <!-- Card Header -->
                    <div class="p-6 border-b cursor-pointer" @click="selectPrompt(prompt)">
                      <div class="flex items-start justify-between">
                        <div class="flex-1">
                          <div class="flex items-center gap-3 mb-2">
                            <h4
                              class="font-semibold text-lg text-foreground group-hover:text-primary transition-colors"
                            >
                              {{ prompt.name }}
                            </h4>
                            <span
                              class="inline-flex items-center px-2 py-1 text-xs font-medium rounded-full bg-blue-100 text-blue-700 dark:bg-blue-900/20 dark:text-blue-400"
                            >
                              Live: v{{ prompt.currentVersion }}
                            </span>
                          </div>

                          <p v-if="prompt.description" class="text-muted-foreground text-sm mb-3 leading-relaxed">
                            {{ prompt.description }}
                          </p>

                          <!-- Tags -->
                          <div class="flex flex-wrap gap-1 mb-3">
                            <span
                              v-for="tag in prompt.tags.slice(0, 3)"
                              :key="tag"
                              class="inline-flex items-center gap-1 px-2 py-1 text-xs rounded-md bg-primary/10 text-primary"
                            >
                              <Tag class="h-3 w-3" />
                              {{ tag }}
                            </span>
                            <span
                              v-if="prompt.tags.length > 3"
                              class="inline-flex items-center px-2 py-1 text-xs rounded-md bg-muted text-muted-foreground"
                            >
                              +{{ prompt.tags.length - 3 }} more
                            </span>
                          </div>

                          <!-- Metadata Row -->
                          <div class="flex items-center gap-4 text-xs text-muted-foreground">
                            <div class="flex items-center gap-1">
                              <History class="h-3 w-3" />
                              <span>{{ prompt.versions.length }} versions</span>
                            </div>
                            <div class="flex items-center gap-1">
                              <Settings class="h-3 w-3" />
                              <span>{{ prompt.applications.length }} apps</span>
                            </div>
                          </div>
                        </div>

                        <!-- Selection Indicator -->
                        <div class="ml-4">
                          <div
                            v-if="formModel.promptId === prompt.id"
                            class="h-6 w-6 rounded-full bg-primary flex items-center justify-center"
                          >
                            <span class="text-white text-sm">✓</span>
                          </div>
                          <div
                            v-else
                            class="h-6 w-6 rounded-full border-2 border-muted group-hover:border-primary transition-colors"
                          ></div>
                        </div>
                      </div>
                    </div>

                    <!-- Version Selection (Shown when prompt is selected) -->
                    <div v-if="formModel.promptId === prompt.id" class="border-b bg-muted/20">
                      <div class="p-4">
                        <h5 class="text-sm font-medium mb-3 flex items-center gap-2">
                          <span>Select Version</span>
                          <span class="text-xs text-muted-foreground">({{ prompt.versions.length }} available)</span>
                        </h5>
                        <div class="max-h-40 overflow-y-auto">
                          <div class="space-y-2 pr-1">
                            <div
                              v-for="version in prompt.versions.slice().reverse()"
                              :key="version.id"
                              class="flex items-center justify-between p-3 border rounded-lg cursor-pointer hover:bg-muted/50 transition-colors"
                              :class="
                                formModel.versionId === version.id ? 'border-primary bg-primary/5' : 'border-border'
                              "
                              @click="selectVersion(version)"
                            >
                              <div class="flex-1">
                                <div class="flex items-center gap-2 mb-1">
                                  <span class="font-mono text-sm font-medium">{{ version.version }}</span>
                                  <span class="text-sm text-foreground">{{ version.name }}</span>
                                  <span
                                    v-if="version.version === prompt.currentVersion"
                                    class="inline-flex items-center px-1.5 py-0.5 text-xs font-medium rounded bg-blue-100 text-blue-700 dark:bg-blue-900/20 dark:text-blue-400"
                                  >
                                    Live
                                  </span>
                                </div>
                                <div class="text-xs text-muted-foreground">
                                  by {{ version.createdBy.split("@")[0] }} •
                                  {{ new Date(version.createdAt).toLocaleDateString() }}
                                </div>
                              </div>
                              <div class="ml-3">
                                <div
                                  v-if="formModel.versionId === version.id"
                                  class="h-4 w-4 rounded-full bg-primary flex items-center justify-center"
                                >
                                  <span class="text-white text-xs">✓</span>
                                </div>
                                <div v-else class="h-4 w-4 rounded-full border border-muted"></div>
                              </div>
                            </div>
                          </div>
                        </div>
                      </div>
                    </div>

                    <!-- Content Preview (Shown when version is selected) -->
                    <div v-if="formModel.promptId === prompt.id && selectedVersionInLibrary" class="bg-muted/10">
                      <div class="p-4">
                        <h5 class="text-sm font-medium mb-3 flex items-center justify-between">
                          <span>Content Preview</span>
                          <span class="text-xs text-muted-foreground">
                            {{ (selectedVersionInLibrary.content || "").length.toLocaleString() }} chars • ~{{
                              estimateTokens(selectedVersionInLibrary.content || "").toLocaleString()
                            }}
                            tokens
                          </span>
                        </h5>
                        <div class="rounded-lg border bg-background/50 p-4 max-h-32 overflow-y-auto">
                          <pre class="text-xs font-mono leading-relaxed whitespace-pre-wrap text-foreground">{{
                            selectedVersionInLibrary.content || "No content available"
                          }}</pre>
                        </div>
                      </div>
                    </div>
                  </div>
                </div>

                <!-- Empty State -->
                <div v-else class="flex flex-col items-center justify-center h-full text-center py-12">
                  <div class="w-24 h-24 rounded-full bg-muted/20 flex items-center justify-center mb-4">
                    <Library class="h-12 w-12 text-muted-foreground" />
                  </div>
                  <h4 class="font-medium text-lg mb-2">No prompts found</h4>
                  <p class="text-muted-foreground text-sm mb-6 max-w-sm">
                    {{
                      searchQuery || applicationFilter
                        ? "Try adjusting your search or filters to find more prompts."
                        : "Create your first prompt to get started with the library."
                    }}
                  </p>
                  <UiButton variant="outline" @click="clearFilters"> Clear Filters </UiButton>
                </div>
              </div>
            </UiScrollArea>
          </div>
        </div>
      </template>

      <template #footer>
        <UiDialogFooter class="flex gap-3 justify-end p-6 border-t flex-shrink-0">
          <UiButton variant="outline" @click="closePromptLibrary"> Cancel </UiButton>
          <UiButton @click="confirmPromptSelection" :disabled="!formModel.promptId || !formModel.versionId">
            Load Prompt
          </UiButton>
        </UiDialogFooter>
      </template>
    </UiDialogContent>
  </UiDialog>
</template>
