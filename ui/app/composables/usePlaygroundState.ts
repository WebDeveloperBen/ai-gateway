import { createInjectionState } from "@vueuse/core"

interface SavedPrompt {
  id: string
  name: string
  description?: string
  tags: string[]
  environments: string[]
  applications: string[]
  currentVersion: string
  createdAt: string
  updatedAt: string
  versions: PromptVersion[]
}

interface PromptVersion {
  id: string
  version: string
  name: string
  description?: string
  content: string
  systemPrompt?: string
  parameters?: {
    temperature?: number
    maxTokens?: number
    topP?: number
    frequencyPenalty?: number
    presencePenalty?: number
  }
  tags: string[]
  createdAt: string
  createdBy: string
  isPublished: boolean
  publishedAt?: string
}

interface ModelData {
  name: string
  maxTokens: number
  costPer1kTokens: {
    input: number
    output: number
  }
}

interface TestResult {
  id: string
  model: string
  timestamp: string
  response: string
  success: boolean
  tokensUsed: {
    total: number
  }
  responseTime: number
  estimatedCost: number
  error?: string
}

type PlaygroundEditorTabOptions = "system" | "user"

const [useProvidePlaygroundState, usePlaygroundState] = createInjectionState(
  (initialState?: {
    activePromptTab?: "system" | "user"
    activeTemplateTab?: "system" | "user"
    promptText?: string
    systemPrompt?: string
    isLoading?: boolean
  }) => {
    // Basic editor state
    const activePromptTab = ref<PlaygroundEditorTabOptions>(initialState?.activePromptTab ?? "user")
    const activeTemplateTab = ref<PlaygroundEditorTabOptions>(initialState?.activeTemplateTab ?? "user")
    const promptText = ref<string>(initialState?.promptText ?? "")
    const systemPrompt = ref<string>(initialState?.systemPrompt ?? "")
    const isLoading = ref(initialState?.isLoading ?? false)

    // Prompt/version state using reactive for better encapsulation
    const promptState = reactive({
      currentPrompt: null as SavedPrompt | null,
      currentVersion: null as PromptVersion | null,
      selectedVersionId: ""
    })

    // Model and test state
    const modelState = reactive({
      selectedModel: "",
      selectedModelData: null as ModelData | null | undefined,
      testResults: [] as TestResult[],
      parameters: {
        temperature: 0.7,
        maxTokens: 1000,
        topP: 1.0,
        frequencyPenalty: 0.0,
        presencePenalty: 0.0
      }
    })

    // Modal state
    const modalState = reactive({
      showParametersModal: false,
      showPromptLibrary: false,
      showReplaceWarning: false
    })

    // Library state
    const libraryState = reactive({
      searchQuery: "",
      environmentFilter: "",
      applicationFilter: "",
      formModel: {
        promptId: "",
        versionId: ""
      }
    })

    // Template refs - to be set by components
    const promptTextarea = ref<HTMLTextAreaElement | null>(null)
    const systemPromptTextarea = ref<HTMLTextAreaElement | null>(null)

    // Template insertion
    const insertTemplate = (template: PromptTemplate, templateType: PlaygroundEditorTabOptions = "user") => {
      const isSystemTemplate = templateType === "system"
      const textarea = isSystemTemplate ? systemPromptTextarea.value : promptTextarea.value

      // Switch to the appropriate tab
      activePromptTab.value = templateType

      if (!textarea) {
        // Fallback to simple append if textarea not available
        if (isSystemTemplate) {
          systemPrompt.value += (systemPrompt.value ? "\n\n" : "") + template.content
        } else {
          promptText.value += (promptText.value ? "\n\n" : "") + template.content
        }
        return
      }

      // Insert at cursor position
      const cursorPosition = textarea.selectionStart || 0
      const currentText = isSystemTemplate ? systemPrompt.value : promptText.value

      const beforeText = currentText.slice(0, cursorPosition)
      const afterText = currentText.slice(cursorPosition)
      const newText = beforeText + template.content + afterText

      if (isSystemTemplate) {
        systemPrompt.value = newText
      } else {
        promptText.value = newText
      }

      nextTick(() => {
        if (textarea) {
          textarea.focus()
          const newCursorPosition = cursorPosition + template.content.length
          textarea.setSelectionRange(newCursorPosition, newCursorPosition)
        }
      })
    }

    // Prompt editor actions
    const actions = {
      clearPrompt: () => {
        promptText.value = ""
        promptTextarea.value?.focus()
      },

      clearSystemPrompt: () => {
        systemPrompt.value = ""
      },

      copyPrompt: async () => {
        try {
          let content = ""
          if (systemPrompt.value.trim()) {
            content += `SYSTEM: ${systemPrompt.value.trim()}\n\n`
          }
          if (promptText.value.trim()) {
            content += `USER: ${promptText.value.trim()}`
          }
          await navigator.clipboard.writeText(content)
        } catch (err) {
          console.error("Failed to copy prompt:", err)
        }
      },

      savePrompt: () => {
        // TODO: Implement save logic - would call API
        console.log("Saving prompt:", {
          content: promptText.value,
          systemPrompt: systemPrompt.value,
          currentPrompt: promptState.currentPrompt
        })
      },

      copyResponse: async (response: string) => {
        try {
          await navigator.clipboard.writeText(response)
        } catch (err) {
          console.error("Failed to copy response:", err)
        }
      },

      loadVersion: () => {
        if (!promptState.currentPrompt || !promptState.selectedVersionId) {
          // Load draft (empty state)
          promptState.currentVersion = null
          promptText.value = ""
          systemPrompt.value = ""
          // Reset parameters to defaults
          modelState.parameters = {
            temperature: 0.7,
            maxTokens: 1000,
            topP: 1.0,
            frequencyPenalty: 0.0,
            presencePenalty: 0.0
          }
          return
        }

        const version = promptState.currentPrompt.versions.find((v) => v.id === promptState.selectedVersionId)
        if (version) {
          promptState.currentVersion = version
          promptText.value = version.content
          systemPrompt.value = version.systemPrompt || ""

          // Load version's parameters or use defaults
          if (version.parameters) {
            modelState.parameters = {
              temperature: version.parameters.temperature ?? 0.7,
              maxTokens: version.parameters.maxTokens ?? 1000,
              topP: version.parameters.topP ?? 1.0,
              frequencyPenalty: version.parameters.frequencyPenalty ?? 0.0,
              presencePenalty: version.parameters.presencePenalty ?? 0.0
            }
          } else {
            // No parameters saved with this version, use defaults
            modelState.parameters = {
              temperature: 0.7,
              maxTokens: 1000,
              topP: 1.0,
              frequencyPenalty: 0.0,
              presencePenalty: 0.0
            }
          }
        }
      }
    }

    // Modal actions
    const modalActions = {
      openParametersModal: () => {
        modalState.showParametersModal = true
      },
      closeParametersModal: () => {
        modalState.showParametersModal = false
      },
      openPromptLibrary: () => {
        modalState.showPromptLibrary = true
        libraryState.formModel = { promptId: "", versionId: "" }
        libraryState.searchQuery = ""
        libraryState.environmentFilter = ""
        libraryState.applicationFilter = ""
      },
      closePromptLibrary: () => {
        modalState.showPromptLibrary = false
      },
      openReplaceWarning: () => {
        modalState.showReplaceWarning = true
      },
      closeReplaceWarning: () => {
        modalState.showReplaceWarning = false
      }
    }

    // Utility functions
    const utils = {
      estimateTokens: (text: string): number => {
        if (!text) return 0
        let tokenEstimate = 0

        tokenEstimate += text.length / 4
        const punctuation = (text.match(/[.,;:!?(){}[\]"']/g) || []).length
        tokenEstimate += punctuation * 0.5
        const lineBreaks = (text.match(/\n/g) || []).length
        tokenEstimate += lineBreaks * 0.3
        const numbers = (text.match(/\d+/g) || []).length
        tokenEstimate += numbers * 0.2

        return Math.ceil(tokenEstimate)
      },

      formatTime: (dateString: string) => {
        return new Date(dateString).toLocaleTimeString("en-US", {
          hour: "2-digit",
          minute: "2-digit"
        })
      },

      formatCurrency: (amount: number) => {
        return new Intl.NumberFormat("en-US", {
          style: "currency",
          currency: "USD",
          minimumFractionDigits: 6
        }).format(amount)
      }
    }

    return {
      // Basic editor state
      activePromptTab,
      activeTemplateTab,
      promptText,
      systemPrompt,
      isLoading,

      // Reactive state objects
      promptState,
      modelState,
      modalState,
      libraryState,

      // Template refs
      promptTextarea,
      systemPromptTextarea,

      // Functions
      insertTemplate,
      actions,
      modalActions,
      utils
    }
  }
)

export { useProvidePlaygroundState }
export { usePlaygroundState }

