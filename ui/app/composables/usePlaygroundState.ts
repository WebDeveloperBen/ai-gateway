import { createInjectionState } from '@vueuse/core'

const [useProvidePlaygroundState, usePlaygroundState] = createInjectionState(
  (initialState?: {
    activePromptTab?: 'system' | 'user'
    activeTemplateTab?: 'system' | 'user'
    promptText?: string
    systemPrompt?: string
    isLoading?: boolean
  }) => {
    const activePromptTab = ref<'system' | 'user'>(initialState?.activePromptTab ?? 'user')
    const activeTemplateTab = ref<'system' | 'user'>(initialState?.activeTemplateTab ?? 'user')
    const promptText = ref<string>(initialState?.promptText ?? '')
    const systemPrompt = ref<string>(initialState?.systemPrompt ?? '')
    const isLoading = ref(initialState?.isLoading ?? false)

    const promptTextarea = useTemplateRef<HTMLTextAreaElement>('promptTextarea')
    const systemPromptTextarea = useTemplateRef<HTMLTextAreaElement>('systemPromptTextarea')

    const insertTemplate = (template: PromptTemplate, templateType: 'system' | 'user' = activeTemplateTab.value) => {
      const isSystemTemplate = templateType === 'system'
      const textarea = isSystemTemplate ? systemPromptTextarea.value : promptTextarea.value

      if (!textarea) {
        activePromptTab.value = templateType
        activeTemplateTab.value = templateType
        nextTick(() => {
          insertTemplate(template, templateType)
        })
        return
      }

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

    return {
      activePromptTab,
      activeTemplateTab,
      promptText,
      systemPrompt,
      isLoading,
      promptTextarea,
      systemPromptTextarea,
      insertTemplate
    }
  }
)

export { useProvidePlaygroundState }
export { usePlaygroundState }