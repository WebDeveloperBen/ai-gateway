<script lang="ts">
import { Play, FileText, Library, MessageSquare, Code, PenTool, BarChart, Shield, Loader2 } from "lucide-vue-next"
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

const { promptText, systemPrompt, isLoading, promptState, modelState, utils } = playgroundState

// State management
const selectedModel = ref<string>("")
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
    const systemTokens = utils.estimateTokens(systemPrompt.value)
    const userTokens = utils.estimateTokens(promptText.value)
    const inputTokens = systemTokens + userTokens
    const outputTokens = Math.min(utils.estimateTokens(mockResponse), modelParameters.value.maxTokens)
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

    modelState.testResults.unshift(result)
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

    modelState.testResults.unshift(result)
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

const createNewPrompt = () => {
  promptState.currentPrompt = null
  promptState.currentVersion = null
  selectedPromptId.value = ""
  promptState.selectedVersionId = ""
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

// Clear selection when filters change
watch([librarySearchQuery, libraryApplicationFilter], () => {
  libraryFormModel.value.promptId = ""
  libraryFormModel.value.versionId = ""
})

// Watch for changes to prompt content and model parameters to switch to draft mode
watch(
  [promptText, systemPrompt, modelParameters],
  () => {
    if (promptState.currentVersion) {
      const contentChanged = promptText.value !== promptState.currentVersion.content
      const systemPromptChanged = systemPrompt.value !== (promptState.currentVersion.systemPrompt || "")

      // Check if model parameters have changed
      const parametersChanged = promptState.currentVersion.parameters
        ? modelParameters.value.temperature !== (promptState.currentVersion.parameters.temperature ?? 0.7) ||
          modelParameters.value.maxTokens !== (promptState.currentVersion.parameters.maxTokens ?? 1000) ||
          modelParameters.value.topP !== (promptState.currentVersion.parameters.topP ?? 1.0) ||
          modelParameters.value.frequencyPenalty !== (promptState.currentVersion.parameters.frequencyPenalty ?? 0.0) ||
          modelParameters.value.presencePenalty !== (promptState.currentVersion.parameters.presencePenalty ?? 0.0)
        : // If no parameters saved, check if current parameters differ from defaults
          modelParameters.value.temperature !== 0.7 ||
          modelParameters.value.maxTokens !== 1000 ||
          modelParameters.value.topP !== 1.0 ||
          modelParameters.value.frequencyPenalty !== 0.0 ||
          modelParameters.value.presencePenalty !== 0.0

      if (contentChanged || systemPromptChanged || parametersChanged) {
        // Content or parameters have been modified, switch to draft
        promptState.currentVersion = null
        promptState.selectedVersionId = ""
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
    promptState.currentPrompt = prompt as any
    promptState.currentVersion = version as any
    selectedPromptId.value = promptId
    promptState.selectedVersionId = versionId
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
      <!-- Prompt Editor -->
      <PromptEditor :showParametersModal="() => (showParametersModal = true)" />

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
          >Cost: {{ utils.formatCurrency(selectedModelData.costPer1kTokens.input) }}/1k in,
          {{ utils.formatCurrency(selectedModelData.costPer1kTokens.output) }}/1k out</span
        >
      </div>
    </div>

    <!-- Prompt Library Modal -->
    <ModalsPlaygroundPromptLibrary
      :open="showPromptLibrary"
      @update:open="showPromptLibrary = $event"
      :saved-prompts="savedPrompts"
      :available-applications="availableApplications"
      :library-search-query="librarySearchQuery"
      @update:library-search-query="librarySearchQuery = $event"
      :library-application-filter="libraryApplicationFilter"
      @update:library-application-filter="libraryApplicationFilter = $event"
      :library-form-model="libraryFormModel"
      @update:library-form-model="libraryFormModel = $event"
      @select-prompt="selectPrompt"
      @select-version="selectVersion"
      @confirm-selection="confirmPromptSelection"
    />

    <ModalsPlaygroundReplaceContentWarning
      :open="showReplaceWarning"
      :pending-prompt-id="pendingPromptId"
      @update:open="showReplaceWarning = $event"
      @confirm-replace="confirmReplaceContent"
      @cancel-replace="cancelReplaceContent"
    />

    <!-- Model Parameters Modal -->
    <ModalsPlaygroundModelParameters
      :open="showParametersModal"
      :model-parameters="modelParameters"
      :selected-model-data="selectedModelData"
      @update:open="showParametersModal = $event"
      @update:model-parameters="modelParameters = $event"
      @reset-to-defaults="
        modelParameters = {
          temperature: 0.7,
          maxTokens: 1000,
          topP: 1.0,
          frequencyPenalty: 0.0,
          presencePenalty: 0.0
        }
      "
      @apply-settings="showParametersModal = false"
    />
  </div>
</template>
