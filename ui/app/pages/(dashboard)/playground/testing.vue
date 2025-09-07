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
  Zap,
  Loader2,
  Save,
  Settings,
  History,
  Tag,
  AlertCircle
} from "lucide-vue-next"
import type { FormBuilder } from "~/components/Ui/FormBuilder/FormBuilder.vue"

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

const availableApplications = ref<Application[]>([
  { id: "chat-app", name: "Chat Application", description: "Main customer chat interface", environment: "prod" },
  { id: "analytics-api", name: "Analytics API", description: "Data analytics service", environment: "prod" },
  { id: "user-service", name: "User Service", description: "User management service", environment: "prod" },
  { id: "support-bot", name: "Support Bot", description: "Automated support assistant", environment: "staging" }
])

const savedPrompts = ref<SavedPrompt[]>([
  {
    id: "prompt_customer_support_v1",
    name: "Customer Support Assistant",
    description: "Standard customer support prompt for chat applications",
    tags: ["customer-support", "chat", "professional"],
    environments: ["prod", "staging"],
    applications: ["chat-app", "support-bot"],
    currentVersion: "v1.2.0",
    createdAt: "2024-01-15T10:00:00Z",
    updatedAt: "2024-02-01T14:30:00Z",
    versions: [
      {
        id: "v1.0.0",
        version: "v1.0.0",
        name: "Initial Version",
        systemPrompt:
          "You are a helpful customer support assistant for a SaaS company. Always be polite and professional.",
        content: "Please help me with: {customer_query}",
        parameters: {
          temperature: 0.3,
          maxTokens: 500,
          topP: 0.9,
          frequencyPenalty: 0.0,
          presencePenalty: 0.0
        },
        tags: ["customer-support"],
        createdAt: "2024-01-15T10:00:00Z",
        createdBy: "john.doe@company.com",
        isPublished: true,
        publishedAt: "2024-01-15T10:05:00Z"
      },
      {
        id: "v1.1.0",
        version: "v1.1.0",
        name: "Added Empathy",
        systemPrompt:
          "You are a helpful and empathetic customer support assistant for a SaaS company. Always be polite, professional, and understanding of customer concerns. Use a warm, friendly tone.",
        content:
          "Customer inquiry: {customer_query}\n\nPlease provide a helpful response that acknowledges their concern and offers practical solutions.",
        parameters: {
          temperature: 0.4,
          maxTokens: 750,
          topP: 0.95,
          frequencyPenalty: 0.1,
          presencePenalty: 0.1
        },
        tags: ["customer-support", "empathy"],
        createdAt: "2024-01-20T15:30:00Z",
        createdBy: "jane.smith@company.com",
        isPublished: true,
        publishedAt: "2024-01-21T09:00:00Z"
      },
      {
        id: "v1.2.0",
        version: "v1.2.0",
        name: "Enhanced Guidelines",
        systemPrompt:
          "You are a helpful and empathetic customer support assistant for a SaaS company. Always be polite, professional, and understanding of customer concerns. Use a warm, friendly tone. If you cannot resolve an issue, offer to escalate to a human agent. Always end with asking if there's anything else you can help with.",
        content:
          "Customer inquiry: {customer_query}\n\nPlease provide a comprehensive response that:\n1. Acknowledges their concern\n2. Offers step-by-step solutions if applicable\n3. Suggests escalation if needed\n4. Asks for additional help",
        parameters: {
          temperature: 0.5,
          maxTokens: 1000,
          topP: 0.9,
          frequencyPenalty: 0.2,
          presencePenalty: 0.1
        },
        tags: ["customer-support", "empathy", "escalation"],
        createdAt: "2024-02-01T14:30:00Z",
        createdBy: "john.doe@company.com",
        isPublished: true,
        publishedAt: "2024-02-01T16:00:00Z"
      }
    ]
  }
])

const systemPromptTemplates = ref<PromptTemplate[]>([
  {
    id: "customer-support-system",
    name: "Customer Support",
    content:
      "You are a helpful customer support assistant for a SaaS company. Always be polite, professional, and empathetic. Follow these guidelines:\n\n- Listen to customer concerns carefully\n- Provide clear, step-by-step solutions\n- Escalate complex issues when appropriate\n- Use a warm, friendly tone\n- End conversations by asking if there's anything else you can help with",
    category: "Support",
    description: "Professional customer service system prompt",
    icon: MessageSquare
  },
  {
    id: "technical-writer-system",
    name: "Technical Writer",
    content:
      "You are an experienced technical writer specializing in software documentation. Your role is to:\n\n- Write clear, concise documentation\n- Use proper technical terminology\n- Structure information logically\n- Include relevant examples and code snippets\n- Ensure accessibility for different skill levels",
    category: "Documentation",
    description: "Technical documentation specialist",
    icon: FileText
  },
  {
    id: "code-reviewer-system",
    name: "Code Reviewer",
    content:
      "You are a senior software engineer conducting code reviews. Focus on:\n\n- Code quality and best practices\n- Security vulnerabilities\n- Performance optimization\n- Maintainability and readability\n- Testing coverage\n- Architecture patterns",
    category: "Development",
    description: "Senior engineer code reviewer",
    icon: Code
  },
  {
    id: "data-analyst-system",
    name: "Data Analyst",
    content:
      "You are a data analyst with expertise in statistics and business intelligence. Your responsibilities include:\n\n- Analyzing data patterns and trends\n- Creating actionable insights\n- Explaining complex findings simply\n- Recommending data-driven decisions\n- Identifying potential data quality issues",
    category: "Analytics",
    description: "Professional data analyst",
    icon: BarChart
  },
  {
    id: "security-expert-system",
    name: "Security Expert",
    content:
      "You are a cybersecurity expert with deep knowledge of:\n\n- Security vulnerabilities and threats\n- Risk assessment methodologies\n- Compliance standards (SOC2, ISO 27001, etc.)\n- Incident response procedures\n- Security architecture best practices\n\nAlways prioritize security while being practical about implementation.",
    category: "Security",
    description: "Cybersecurity specialist",
    icon: Shield
  },
  {
    id: "creative-writer-system",
    name: "Creative Writer",
    content:
      "You are a creative writer with expertise in storytelling and content creation. Your strengths include:\n\n- Engaging narrative techniques\n- Adapting tone for different audiences\n- Creative problem-solving\n- Brand voice consistency\n- Compelling calls-to-action\n\nFocus on creating memorable, impactful content.",
    category: "Writing",
    description: "Creative content specialist",
    icon: PenTool
  }
])

const userPromptTemplates = ref<PromptTemplate[]>([
  {
    id: "help-request",
    name: "Help Request",
    content:
      "I need help with {describe_your_issue}. Here are the details:\n\n- What I'm trying to do: {goal}\n- What I've tried so far: {attempts}\n- Current error/problem: {specific_issue}\n\nCan you provide step-by-step guidance?",
    category: "Support",
    description: "Structured help request template",
    icon: MessageSquare
  },
  {
    id: "code-review-request",
    name: "Code Review",
    content:
      "Please review the following code for:\n\n```{language}\n{your_code_here}\n```\n\nSpecific areas of concern:\n- {concern_1}\n- {concern_2}\n\nPlease focus on {security|performance|maintainability|best_practices}.",
    category: "Development",
    description: "Code review request template",
    icon: Code
  },
  {
    id: "data-analysis-request",
    name: "Data Analysis",
    content:
      "Please analyze the following dataset and provide insights:\n\n**Dataset:** {describe_your_data}\n**Context:** {business_context}\n**Questions to answer:**\n1. {question_1}\n2. {question_2}\n3. {question_3}\n\n**Expected outcome:** {what_you_hope_to_learn}",
    category: "Analytics",
    description: "Data analysis request template",
    icon: BarChart
  },
  {
    id: "content-brief",
    name: "Content Brief",
    content:
      "Create {content_type} about {topic}.\n\n**Target audience:** {audience_description}\n**Tone:** {professional|casual|technical|friendly}\n**Key points to cover:**\n- {point_1}\n- {point_2}\n- {point_3}\n\n**Call to action:** {desired_action}\n**Word count:** {approximate_length}",
    category: "Writing",
    description: "Content creation brief template",
    icon: PenTool
  },
  {
    id: "security-assessment",
    name: "Security Assessment",
    content:
      "Please conduct a security assessment of {system/application/process}.\n\n**Scope:**\n- {component_1}\n- {component_2}\n\n**Compliance requirements:** {standards_if_applicable}\n**Risk tolerance:** {low|medium|high}\n**Timeline:** {assessment_deadline}\n\nFocus on {authentication|authorization|data_protection|network_security}.",
    category: "Security",
    description: "Security assessment request",
    icon: Shield
  },
  {
    id: "technical-explanation",
    name: "Technical Explanation",
    content:
      "Explain {technical_concept} in simple terms.\n\n**Audience level:** {beginner|intermediate|advanced}\n**Format preference:** {step-by-step|overview|detailed_guide}\n**Include examples:** {yes/no}\n\n**Specific aspects to cover:**\n- {aspect_1}\n- {aspect_2}\n\nPlease use analogies where helpful.",
    category: "Documentation",
    description: "Technical explanation request",
    icon: FileText
  }
])
</script>
<script setup lang="ts">
useSeoMeta({ title: "Model Testing - LLM Gateway" })

// Initialize playground state provider
const playgroundState = useProvidePlaygroundState({
  activePromptTab: "user",
  activeTemplateTab: "user",
  promptText: "",
  systemPrompt: "",
  isLoading: false
})

const {
  activePromptTab,
  activeTemplateTab,
  promptText,
  systemPrompt,
  isLoading,
  promptTextarea,
  systemPromptTextarea
} = playgroundState

// State management
const selectedModel = ref<string>("")
const testResults = ref<TestResult[]>([])
const showParametersModal = ref(false)

// Model parameters
const modelParameters = ref({
  temperature: 0.7,
  maxTokens: 1000,
  topP: 1.0,
  frequencyPenalty: 0.0,
  presencePenalty: 0.0
})

// Prompt versioning state
const selectedPromptId = ref<string>("")
const selectedVersionId = ref<string>("")
const currentPrompt = ref<SavedPrompt | null>(null)
const currentVersion = ref<PromptVersion | null>(null)
const showReplaceWarning = ref(false)
const showPromptLibrary = ref(false)
const pendingPromptId = ref<string>("")
const pendingVersionId = ref<string>("")
const promptMetadata = ref({
  name: "",
  description: "",
  tags: [] as string[],
  environments: [] as string[],
  applications: [] as string[]
})

// Form model for the UiFormBuilder
const formModel = ref({
  model: ""
})

// Modal form model for prompt library selection
const libraryFormModel = ref({
  promptId: "",
  versionId: ""
})

// Library search and filter state
const librarySearchQuery = ref("")
const libraryEnvironmentFilter = ref("")
const libraryApplicationFilter = ref("")

// Model selection
const selectedModelData = computed(() => {
  return availableModels.value.find((m) => m.id === selectedModel.value)
})

// Form fields for model configuration
const formFields = computed((): FormBuilder[] => [
  {
    variant: "Select",
    label: "Model",
    name: "model",
    placeholder: "Select a model...",
    required: true,
    modelValue: formModel.value.model,
    options: availableModels.value.map((model) => ({
      value: model.id,
      label: `${model.name} (${model.provider})`,
      description: model.description
    })),
    wrapperClass: "col-span-8"
  }
])

// Filtered prompts based on search and filters
const filteredPrompts = computed(() => {
  let filtered = savedPrompts.value

  // Apply search filter
  if (librarySearchQuery.value.trim()) {
    const query = librarySearchQuery.value.toLowerCase()
    filtered = filtered.filter(
      (prompt) =>
        prompt.name.toLowerCase().includes(query) ||
        prompt.description?.toLowerCase().includes(query) ||
        prompt.tags.some((tag) => tag.toLowerCase().includes(query))
    )
  }

  // Apply environment filter
  if (libraryEnvironmentFilter.value) {
    filtered = filtered.filter((prompt) => prompt.environments.includes(libraryEnvironmentFilter.value))
  }

  // Apply application filter
  if (libraryApplicationFilter.value) {
    filtered = filtered.filter((prompt) => prompt.applications.includes(libraryApplicationFilter.value))
  }

  return filtered
})

// Selected prompt and version for cleaner template access
const selectedPromptInLibrary = computed(() => {
  return filteredPrompts.value.find((p) => p.id === libraryFormModel.value.promptId)
})

const selectedVersionInLibrary = computed(() => {
  return selectedPromptInLibrary.value?.versions.find((v) => v.id === libraryFormModel.value.versionId)
})

// Functions

const runTest = async () => {
  if (!selectedModel.value || (!promptText.value.trim() && !systemPrompt.value.trim())) {
    return
  }

  isLoading.value = true

  try {
    const startTime = Date.now()

    // Build the complete prompt with system and user parts
    const fullPrompt = buildFullPrompt()

    // TODO: Replace with actual API call
    // const response = await $fetch('/api/playground/test', {
    //   method: 'POST',
    //   body: {
    //     model: selectedModel.value,
    //     systemPrompt: systemPrompt.value,
    //     prompt: promptText.value,
    //     parameters: modelParameters.value
    //   }
    // })

    // Simulate API response with more realistic variations based on parameters
    const simulationDelay = Math.max(500, 2000 - modelParameters.value.temperature * 500)
    await new Promise((resolve) => setTimeout(resolve, simulationDelay))
    const endTime = Date.now()

    // Generate mock response that varies based on temperature
    const mockResponses = [
      "This is a focused and deterministic response from the AI model. The system has processed your request according to the specified parameters.",
      "Here's a creative response that demonstrates how the AI adapts to different temperature settings and system constraints.",
      "An analytical response showcasing the model's ability to follow system instructions while incorporating the specified behavioral parameters.",
      "A comprehensive answer that balances creativity and precision based on your configured model parameters and system prompt."
    ]

    const mockResponse =
      mockResponses[Math.floor(Math.random() * mockResponses.length)] +
      (modelParameters.value.temperature > 1.0
        ? " With higher creativity settings, responses become more varied and exploratory!"
        : "")

    // Calculate tokens including system prompt
    const systemTokens = estimateTokens(systemPrompt.value)
    const userTokens = estimateTokens(promptText.value)
    const inputTokens = systemTokens + userTokens
    const outputTokens = Math.min(estimateTokens(mockResponse), modelParameters.value.maxTokens)
    const totalTokens = inputTokens + outputTokens

    const modelData = selectedModelData.value!
    const estimatedCost =
      (inputTokens * modelData.costPer1kTokens.input + outputTokens * modelData.costPer1kTokens.output) / 1000

    const result: TestResult = {
      id: `test_${Date.now()}`,
      timestamp: new Date().toISOString(),
      prompt: fullPrompt,
      response: mockResponse,
      model: `${modelData.name} (T:${modelParameters.value.temperature})`,
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
      prompt: buildFullPrompt(),
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

const buildFullPrompt = () => {
  let fullPrompt = ""
  if (systemPrompt.value.trim()) {
    fullPrompt += `SYSTEM: ${systemPrompt.value.trim()}\n\n`
  }
  if (promptText.value.trim()) {
    fullPrompt += `USER: ${promptText.value.trim()}`
  }
  return fullPrompt
}

const clearPrompt = () => {
  promptText.value = ""
  promptTextarea.value?.focus()
}

const clearSystemPrompt = () => {
  systemPrompt.value = ""
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
    const content = buildFullPrompt()
    await navigator.clipboard.writeText(content)
    // TODO: Show success toast
    console.log("Prompt content copied to clipboard")
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

const loadVersion = () => {
  if (!currentPrompt.value || !selectedVersionId.value) {
    // Load draft (empty state)
    currentVersion.value = null
    promptText.value = ""
    systemPrompt.value = ""
    // Reset parameters to defaults
    modelParameters.value = {
      temperature: 0.7,
      maxTokens: 1000,
      topP: 1.0,
      frequencyPenalty: 0.0,
      presencePenalty: 0.0
    }
    return
  }

  const version = currentPrompt.value.versions.find((v) => v.id === selectedVersionId.value)
  if (version) {
    currentVersion.value = version
    promptText.value = version.content
    systemPrompt.value = version.systemPrompt || ""

    // Load version's parameters or use defaults
    if (version.parameters) {
      modelParameters.value = {
        temperature: version.parameters.temperature ?? 0.7,
        maxTokens: version.parameters.maxTokens ?? 1000,
        topP: version.parameters.topP ?? 1.0,
        frequencyPenalty: version.parameters.frequencyPenalty ?? 0.0,
        presencePenalty: version.parameters.presencePenalty ?? 0.0
      }
    } else {
      // No parameters saved with this version, use defaults
      modelParameters.value = {
        temperature: 0.7,
        maxTokens: 1000,
        topP: 1.0,
        frequencyPenalty: 0.0,
        presencePenalty: 0.0
      }
    }
  }
}

const savePrompt = async () => {
  // TODO: Implement save logic - this would call the API
  console.log("Saving prompt:", {
    promptId: selectedPromptId.value,
    content: promptText.value,
    metadata: promptMetadata.value
  })

  // For now, just show a placeholder
  alert("Save functionality will be implemented with API integration")
}

const createNewPrompt = () => {
  currentPrompt.value = null
  currentVersion.value = null
  selectedPromptId.value = ""
  selectedVersionId.value = ""
  promptText.value = ""
  systemPrompt.value = ""
  modelParameters.value = {
    temperature: 0.7,
    maxTokens: 1000,
    topP: 1.0,
    frequencyPenalty: 0.0,
    presencePenalty: 0.0
  }
  promptMetadata.value = {
    name: "",
    description: "",
    tags: [],
    environments: [],
    applications: []
  }
}

// Form change handlers
const handleModelChange = (value: string) => {
  selectedModel.value = value
  formModel.value.model = value
}

// Watch for library form changes
watch(
  () => libraryFormModel.value.promptId,
  (newPromptId) => {
    if (newPromptId) {
      // Reset version when prompt changes
      libraryFormModel.value.versionId = ""
    }
  },
  { immediate: false }
)

// Modal handlers
const openPromptLibrary = () => {
  showPromptLibrary.value = true
  libraryFormModel.value = { promptId: "", versionId: "" }
  librarySearchQuery.value = ""
  libraryEnvironmentFilter.value = ""
  libraryApplicationFilter.value = ""
}

const closePromptLibrary = () => {
  showPromptLibrary.value = false
  libraryFormModel.value = { promptId: "", versionId: "" }
  librarySearchQuery.value = ""
  libraryEnvironmentFilter.value = ""
  libraryApplicationFilter.value = ""
}

// Clear selection when filters change
watch([librarySearchQuery, libraryApplicationFilter], () => {
  libraryFormModel.value.promptId = ""
  libraryFormModel.value.versionId = ""
})

// Watch for changes to prompt content and model parameters to switch to draft mode
watch(
  [promptText, systemPrompt, modelParameters],
  () => {
    if (currentVersion.value) {
      const contentChanged = promptText.value !== currentVersion.value.content
      const systemPromptChanged = systemPrompt.value !== (currentVersion.value.systemPrompt || "")

      // Check if model parameters have changed
      const parametersChanged = currentVersion.value.parameters
        ? modelParameters.value.temperature !== (currentVersion.value.parameters.temperature ?? 0.7) ||
          modelParameters.value.maxTokens !== (currentVersion.value.parameters.maxTokens ?? 1000) ||
          modelParameters.value.topP !== (currentVersion.value.parameters.topP ?? 1.0) ||
          modelParameters.value.frequencyPenalty !== (currentVersion.value.parameters.frequencyPenalty ?? 0.0) ||
          modelParameters.value.presencePenalty !== (currentVersion.value.parameters.presencePenalty ?? 0.0)
        : // If no parameters saved, check if current parameters differ from defaults
          modelParameters.value.temperature !== 0.7 ||
          modelParameters.value.maxTokens !== 1000 ||
          modelParameters.value.topP !== 1.0 ||
          modelParameters.value.frequencyPenalty !== 0.0 ||
          modelParameters.value.presencePenalty !== 0.0

      if (contentChanged || systemPromptChanged || parametersChanged) {
        // Content or parameters have been modified, switch to draft
        currentVersion.value = null
        selectedVersionId.value = ""
      }
    }
  },
  { deep: true }
)

// Enhanced selection handlers for card-based interface
const selectPrompt = (prompt: SavedPrompt) => {
  if (libraryFormModel.value.promptId === prompt.id) {
    // Toggle off if already selected
    libraryFormModel.value.promptId = ""
    libraryFormModel.value.versionId = ""
  } else {
    // Select new prompt
    libraryFormModel.value.promptId = prompt.id
    libraryFormModel.value.versionId = "" // Reset version selection
  }
}

const selectVersion = (version: PromptVersion) => {
  libraryFormModel.value.versionId = version.id
}

const confirmPromptSelection = () => {
  if (!libraryFormModel.value.promptId || !libraryFormModel.value.versionId) {
    return // Validation - both fields required
  }

  const promptId = libraryFormModel.value.promptId
  const versionId = libraryFormModel.value.versionId

  // Check if we need to warn about replacing content
  if (promptText.value.trim()) {
    pendingPromptId.value = promptId
    pendingVersionId.value = versionId
    showPromptLibrary.value = false
    showReplaceWarning.value = true
  } else {
    // Load directly
    loadPromptVersion(promptId, versionId)
    showPromptLibrary.value = false
  }
}

const loadPromptVersion = (promptId: string, versionId: string) => {
  const prompt = savedPrompts.value.find((p) => p.id === promptId)
  const version = prompt?.versions.find((v) => v.id === versionId)

  if (prompt && version) {
    currentPrompt.value = prompt
    currentVersion.value = version
    selectedPromptId.value = promptId
    selectedVersionId.value = versionId
    promptText.value = version.content
    systemPrompt.value = version.systemPrompt || ""

    // Load parameters if they exist
    if (version.parameters) {
      modelParameters.value = {
        temperature: version.parameters.temperature ?? 0.7,
        maxTokens: version.parameters.maxTokens ?? 1000,
        topP: version.parameters.topP ?? 1.0,
        frequencyPenalty: version.parameters.frequencyPenalty ?? 0.0,
        presencePenalty: version.parameters.presencePenalty ?? 0.0
      }
    }

    // Update metadata
    promptMetadata.value = {
      name: prompt.name,
      description: prompt.description || "",
      tags: [...prompt.tags],
      environments: [...prompt.environments],
      applications: [...prompt.applications]
    }
  }
}

// Warning modal handlers
const confirmReplaceContent = () => {
  if (pendingPromptId.value && pendingVersionId.value) {
    loadPromptVersion(pendingPromptId.value, pendingVersionId.value)
  } else {
    createNewPrompt()
  }
  showReplaceWarning.value = false
  pendingPromptId.value = ""
  pendingVersionId.value = ""
}

const cancelReplaceContent = () => {
  showReplaceWarning.value = false
  pendingPromptId.value = ""
  pendingVersionId.value = ""
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
    formModel.value.model = availableModels.value[0]?.id ?? ""
  }
})
</script>

<template>
  <div class="flex flex-col gap-6 h-screen overflow-hidden">
    <!-- Header -->
    <PageHeader title="Model Testing Playground" subtext="Test prompts against your registered models and providers" />

    <!-- Configuration Form -->
    <div class="flex-shrink-0">
      <form @submit.prevent class="grid grid-cols-12 gap-5 items-end">
        <fieldset :disabled="isLoading" class="contents">
          <UiFormBuilder class="contents" :fields="formFields" v-model="formModel" @update:model="handleModelChange" />
        </fieldset>

        <!-- Browse Library Button -->
        <div class="col-span-2">
          <UiButton variant="outline" size="lg" class="w-full" @click="openPromptLibrary" :disabled="isLoading">
            <Library class="h-4 w-4 mr-2" />
            Browse Library
          </UiButton>
        </div>

        <!-- Test Button -->
        <div class="col-span-2">
          <UiButton
            variant="default"
            size="lg"
            class="w-full"
            @click="runTest"
            :disabled="!selectedModel || (!promptText.trim() && !systemPrompt.trim()) || isLoading"
          >
            <Play v-if="!isLoading" class="h-4 w-4 mr-2" />
            <Loader2 v-else class="h-4 w-4 animate-spin mr-2" />
            {{ isLoading ? "Testing..." : "Run Test" }}
          </UiButton>
        </div>
      </form>
    </div>

    <!-- Main Editor Area -->
    <div class="flex gap-6 flex-1 min-h-0">
      <!-- Left: Prompt Editor and Chat Responses -->
      <div class="flex-1 flex flex-col rounded-lg border bg-background shadow-sm">
        <!-- Card Header -->
        <div class="p-4 border-b">
          <div class="flex items-center justify-between">
            <div class="flex items-center gap-3">
              <FileText class="h-5 w-5 text-primary" />
              <div>
                <h3 class="font-semibold text-base">Prompt Editor</h3>
                <div v-if="currentPrompt" class="flex items-center gap-2 text-xs text-muted-foreground mt-0.5">
                  <span>{{ currentPrompt.name }}</span>
                  <span>•</span>
                  <span
                    v-if="currentVersion"
                    class="inline-flex items-center px-2 py-0.5 rounded-full bg-blue-100 text-blue-700 dark:bg-blue-900/20 dark:text-blue-400"
                  >
                    {{ currentVersion.version }}
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
              <div v-if="currentPrompt && currentPrompt.versions.length > 0" class="flex items-center gap-2">
                <label class="text-xs text-muted-foreground">Version:</label>
                <select
                  v-model="selectedVersionId"
                  @change="loadVersion"
                  class="text-xs border border-border rounded px-2 py-1 bg-background min-w-24"
                >
                  <option value="">Draft</option>
                  <option
                    v-for="version in currentPrompt.versions.slice().reverse()"
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
                @click="copyPrompt"
                :disabled="!promptText && !systemPrompt"
                title="Copy current content"
              >
                <Copy class="h-3.5 w-3.5" />
                <span class="hidden sm:inline">Copy</span>
              </button>

              <button
                class="inline-flex items-center gap-2 px-3 py-1.5 text-sm font-medium rounded-lg bg-primary text-primary-foreground hover:bg-primary/90 transition-all duration-150 ease-in-out hover:shadow-sm disabled:opacity-50 disabled:cursor-not-allowed"
                type="button"
                @click="savePrompt"
                :disabled="isLoading || (!promptText.trim() && !systemPrompt.trim())"
                title="Save as new version"
              >
                <Save class="h-3.5 w-3.5" />
                <span class="hidden sm:inline">Save</span>
              </button>

              <button
                @click="showParametersModal = true"
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

        <!-- Prompt Tabs and Content -->
        <div class="flex flex-1 min-h-0">
          <!-- Left Side: Tabbed Prompt Editor -->
          <div class="flex-1 flex flex-col border-r">
            <!-- Prompt Type Tabs -->
            <div class="border-b bg-muted/20">
              <div class="flex">
                <button
                  @click="
                    () => {
                      activePromptTab = 'system'
                      activeTemplateTab = 'system'
                    }
                  "
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
                  @click="
                    () => {
                      activePromptTab = 'user'
                      activeTemplateTab = 'user'
                    }
                  "
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
                <span>~{{ estimateTokens(systemPrompt).toLocaleString() }} tokens</span>
              </div>

              <!-- Clear System Prompt Button -->
              <button
                v-if="systemPrompt.trim()"
                @click="clearSystemPrompt"
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
                <span>~{{ estimateTokens(promptText).toLocaleString() }} tokens</span>
                <span
                  v-if="
                    selectedModelData &&
                    estimateTokens(promptText) + estimateTokens(systemPrompt) > selectedModelData.maxTokens
                  "
                  class="text-red-500 font-medium"
                >
                  (total exceeds {{ selectedModelData.maxTokens.toLocaleString() }} limit)
                </span>
              </div>

              <!-- Clear User Prompt Button -->
              <button
                v-if="promptText.trim()"
                @click="clearPrompt"
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
                <BlocksEmptyState title="No responses yet" subtext="Run a test to see AI responses here" v-else />
              </UiScrollArea>
            </div>
          </div>
        </div>
      </div>

      <!-- Right: Template Sidebar -->
      <PlaygroundSidebar
        :system-prompt-templates="systemPromptTemplates"
        :user-prompt-templates="userPromptTemplates"
      />
    </div>

    <!-- Model Info Bar -->
    <div v-if="selectedModelData">
      <UiGradientDivider />
      <div class="text-sm pt-4 text-muted-foreground text-center">
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

    <!-- Prompt Library Modal -->
    <div v-if="showPromptLibrary" class="fixed inset-0 z-50 grid place-items-center bg-black/50 backdrop-blur-sm p-4">
      <div class="w-full max-w-6xl rounded-2xl border bg-card shadow-2xl max-h-[95vh] flex flex-col overflow-hidden">
        <!-- Enhanced Header -->
        <div class="border-b bg-gradient-to-r from-card to-muted/20 px-8 py-6 flex-shrink-0">
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
            <div class="flex items-center gap-3">
              <UiButton variant="ghost" size="sm" @click="closePromptLibrary" class="h-10 w-10 p-0">
                <span class="sr-only">Close</span>
                <span class="text-xl text-muted-foreground hover:text-foreground">×</span>
              </UiButton>
            </div>
          </div>
        </div>

        <div class="flex-1 flex min-h-0 overflow-hidden">
          <!-- Enhanced Search and Filters Sidebar -->
          <div class="w-80 border-r bg-muted/10 flex flex-col">
            <!-- Search Section -->
            <div class="p-6 border-b flex-shrink-0">
              <div class="space-y-4">
                <div>
                  <label class="text-sm font-medium text-foreground mb-2 block">Search Prompts</label>
                  <div class="relative">
                    <UiInput
                      v-model="librarySearchQuery"
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
                  <span v-if="librarySearchQuery || libraryApplicationFilter" class="text-primary"> (filtered) </span>
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
                    <input
                      type="radio"
                      v-model="libraryApplicationFilter"
                      value=""
                      class="text-primary focus:ring-primary"
                    />
                    <span class="text-sm">All Applications</span>
                  </label>
                  <label
                    v-for="app in availableApplications"
                    :key="app.id"
                    class="flex items-center gap-2 cursor-pointer"
                  >
                    <input
                      type="radio"
                      v-model="libraryApplicationFilter"
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
              <div v-if="librarySearchQuery || libraryApplicationFilter" class="pt-4 border-t">
                <UiButton
                  variant="outline"
                  size="sm"
                  class="w-full"
                  @click="
                    () => {
                      librarySearchQuery = ''
                      libraryApplicationFilter = ''
                    }
                  "
                >
                  Clear All Filters
                </UiButton>
              </div>
            </div>
          </div>

          <!-- Enhanced Prompt Gallery -->
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
                      libraryFormModel.promptId === prompt.id
                        ? 'ring-2 ring-primary shadow-lg'
                        : 'hover:border-primary/30'
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
                            v-if="libraryFormModel.promptId === prompt.id"
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
                    <div v-if="libraryFormModel.promptId === prompt.id" class="border-b bg-muted/20">
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
                                libraryFormModel.versionId === version.id
                                  ? 'border-primary bg-primary/5'
                                  : 'border-border'
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
                                  v-if="libraryFormModel.versionId === version.id"
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
                    <div v-if="libraryFormModel.promptId === prompt.id && selectedVersionInLibrary" class="bg-muted/10">
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
                      librarySearchQuery || libraryApplicationFilter
                        ? "Try adjusting your search or filters to find more prompts."
                        : "Create your first prompt to get started with the library."
                    }}
                  </p>
                  <UiButton
                    variant="outline"
                    @click="
                      () => {
                        librarySearchQuery = ''
                        libraryApplicationFilter = ''
                      }
                    "
                  >
                    Clear Filters
                  </UiButton>
                </div>
              </div>
            </UiScrollArea>
          </div>
        </div>

        <!-- Actions -->
        <div class="flex gap-3 justify-end p-6 border-t flex-shrink-0">
          <UiButton variant="outline" @click="closePromptLibrary"> Cancel </UiButton>
          <UiButton
            @click="confirmPromptSelection"
            :disabled="!libraryFormModel.promptId || !libraryFormModel.versionId"
          >
            Load Prompt
          </UiButton>
        </div>
      </div>
    </div>

    <!-- Replace Content Warning Modal -->
    <div v-if="showReplaceWarning" class="fixed inset-0 z-50 grid place-items-center bg-black/40 p-4">
      <div class="w-full max-w-md rounded-xl border bg-card shadow-xl">
        <div class="p-6">
          <div class="flex items-center gap-3 mb-4">
            <div class="flex h-10 w-10 items-center justify-center rounded-full bg-orange-100 dark:bg-orange-900/20">
              <AlertCircle class="h-5 w-5 text-orange-600 dark:text-orange-400" />
            </div>
            <div>
              <h3 class="font-semibold">Replace Current Content?</h3>
              <p class="text-sm text-muted-foreground">You have unsaved changes in the prompt editor.</p>
            </div>
          </div>

          <p class="text-sm text-muted-foreground mb-6">
            {{ pendingPromptId ? "Loading a saved prompt" : "Creating a new prompt" }} will replace your current
            content. This action cannot be undone.
          </p>

          <div class="flex gap-3 justify-end">
            <button
              class="px-4 py-2 text-sm font-medium rounded-lg border border-border bg-background hover:bg-muted/50 transition-colors"
              @click="cancelReplaceContent"
            >
              Cancel
            </button>
            <button
              class="px-4 py-2 text-sm font-medium rounded-lg bg-orange-600 text-white hover:bg-orange-700 transition-colors"
              @click="confirmReplaceContent"
            >
              Replace Content
            </button>
          </div>
        </div>
      </div>
    </div>

    <!-- Model Parameters Modal -->
    <div v-if="showParametersModal" class="fixed inset-0 z-50 grid place-items-center bg-black/40 p-4">
      <div class="w-full max-w-2xl rounded-xl border bg-card shadow-xl">
        <div class="flex items-center justify-between border-b px-6 py-4">
          <div>
            <h3 class="font-semibold text-lg flex items-center gap-2">
              <Settings class="h-5 w-5 text-primary" />
              Model Parameters
            </h3>
            <p class="text-sm text-muted-foreground mt-1">Fine-tune AI behavior and output</p>
          </div>
          <UiButton variant="ghost" size="sm" @click="showParametersModal = false">
            <span class="sr-only">Close</span>
            <span class="text-lg">×</span>
          </UiButton>
        </div>

        <div class="p-6">
          <div class="grid grid-cols-1 md:grid-cols-2 gap-6">
            <!-- Temperature -->
            <div class="space-y-3">
              <div class="flex justify-between items-center">
                <label class="text-sm font-medium text-foreground">Temperature</label>
                <span class="text-sm text-muted-foreground font-mono bg-muted px-2 py-1 rounded">{{
                  modelParameters.temperature
                }}</span>
              </div>
              <input
                type="range"
                v-model.number="modelParameters.temperature"
                min="0"
                max="2"
                step="0.1"
                class="w-full h-2 bg-muted rounded-lg appearance-none cursor-pointer slider"
              />
              <p class="text-xs text-muted-foreground">
                Controls randomness. Lower = more focused, Higher = more creative
              </p>
            </div>

            <!-- Max Tokens -->
            <div class="space-y-3">
              <div class="flex justify-between items-center">
                <label class="text-sm font-medium text-foreground">Max Tokens</label>
                <span class="text-sm text-muted-foreground font-mono bg-muted px-2 py-1 rounded">{{
                  modelParameters.maxTokens
                }}</span>
              </div>
              <input
                type="range"
                v-model.number="modelParameters.maxTokens"
                min="50"
                :max="selectedModelData?.maxTokens || 4000"
                step="50"
                class="w-full h-2 bg-muted rounded-lg appearance-none cursor-pointer slider"
              />
              <p class="text-xs text-muted-foreground">Maximum response length</p>
            </div>

            <!-- Top P -->
            <div class="space-y-3">
              <div class="flex justify-between items-center">
                <label class="text-sm font-medium text-foreground">Top P</label>
                <span class="text-sm text-muted-foreground font-mono bg-muted px-2 py-1 rounded">{{
                  modelParameters.topP
                }}</span>
              </div>
              <input
                type="range"
                v-model.number="modelParameters.topP"
                min="0.1"
                max="1"
                step="0.05"
                class="w-full h-2 bg-muted rounded-lg appearance-none cursor-pointer slider"
              />
              <p class="text-xs text-muted-foreground">Controls diversity via nucleus sampling</p>
            </div>

            <!-- Frequency Penalty -->
            <div class="space-y-3">
              <div class="flex justify-between items-center">
                <label class="text-sm font-medium text-foreground">Frequency Penalty</label>
                <span class="text-sm text-muted-foreground font-mono bg-muted px-2 py-1 rounded">{{
                  modelParameters.frequencyPenalty
                }}</span>
              </div>
              <input
                type="range"
                v-model.number="modelParameters.frequencyPenalty"
                min="-2"
                max="2"
                step="0.1"
                class="w-full h-2 bg-muted rounded-lg appearance-none cursor-pointer slider"
              />
              <p class="text-xs text-muted-foreground">Reduces repetition of frequent tokens</p>
            </div>

            <!-- Presence Penalty -->
            <div class="space-y-3">
              <div class="flex justify-between items-center">
                <label class="text-sm font-medium text-foreground">Presence Penalty</label>
                <span class="text-sm text-muted-foreground font-mono bg-muted px-2 py-1 rounded">{{
                  modelParameters.presencePenalty
                }}</span>
              </div>
              <input
                type="range"
                v-model.number="modelParameters.presencePenalty"
                min="-2"
                max="2"
                step="0.1"
                class="w-full h-2 bg-muted rounded-lg appearance-none cursor-pointer slider"
              />
              <p class="text-xs text-muted-foreground">Encourages talking about new topics</p>
            </div>

            <!-- Current Model Info -->
            <div class="md:col-span-2 p-4 bg-muted/20 rounded-lg border">
              <div v-if="selectedModelData" class="text-sm">
                <h4 class="font-medium mb-2">Current Model: {{ selectedModelData.name }}</h4>
                <div class="grid grid-cols-1 md:grid-cols-3 gap-4 text-xs text-muted-foreground">
                  <div>
                    <span class="font-medium">Max Tokens:</span> {{ selectedModelData.maxTokens.toLocaleString() }}
                  </div>
                  <div>
                    <span class="font-medium">Input Cost:</span>
                    {{ formatCurrency(selectedModelData.costPer1kTokens.input) }}/1k
                  </div>
                  <div>
                    <span class="font-medium">Output Cost:</span>
                    {{ formatCurrency(selectedModelData.costPer1kTokens.output) }}/1k
                  </div>
                </div>
              </div>
            </div>
          </div>

          <!-- Actions -->
          <div class="flex gap-3 justify-between pt-6 border-t mt-6">
            <UiButton
              variant="outline"
              @click="
                modelParameters = {
                  temperature: 0.7,
                  maxTokens: 1000,
                  topP: 1.0,
                  frequencyPenalty: 0.0,
                  presencePenalty: 0.0
                }
              "
            >
              Reset to Defaults
            </UiButton>

            <div class="flex gap-3">
              <UiButton variant="outline" @click="showParametersModal = false"> Cancel </UiButton>
              <UiButton @click="showParametersModal = false"> Apply Settings </UiButton>
            </div>
          </div>
        </div>
      </div>
    </div>
  </div>
</template>
