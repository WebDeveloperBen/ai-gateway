<script setup lang="ts">
import { ChevronLeft, ChevronRight, X, Check, ExternalLink, Sparkles, Zap } from "lucide-vue-next"

interface Props {
  open: boolean
}

interface ProviderOption {
  id: string
  name: string
  description: string
  icon: string
  popular?: boolean
}

interface CreateProviderData {
  service: string
  environment: string
  name: string
  slug: string
  // Authentication
  authType: "api-key" | "managed-identity" | "entra-id"
  apiKey?: string
  tenantId?: string
  clientId?: string
  clientSecret?: string
  // Azure specific
  resourceName?: string
  deploymentName?: string
  customDomain?: string
  apiVersion?: string
  region?: string
}

defineProps<Props>()
const emit = defineEmits<{
  "update:open": [value: boolean]
  created: [providerId: string]
}>()

const currentStep = ref(1)
const maxSteps = 4
const loading = ref(false)
const isTransitioning = ref(false)
const showSuccess = ref(false)

const formData = ref<CreateProviderData>({
  service: "",
  environment: "",
  name: "",
  slug: "",
  authType: "api-key",
  apiKey: "",
  tenantId: "",
  clientId: "",
  clientSecret: "",
  resourceName: "",
  deploymentName: "",
  customDomain: "",
  apiVersion: "2024-02-15-preview",
  region: ""
})

// Environment options
const environments = [
  { value: "production", label: "Production", description: "Live production environment", icon: "🚀" },
  { value: "staging", label: "Staging", description: "Pre-production testing", icon: "🧪" },
  { value: "development", label: "Development", description: "Development and testing", icon: "💻" },
  { value: "sandbox", label: "Sandbox", description: "Experimental environment", icon: "🏖️" }
]

// Azure authentication types
const azureAuthTypes = [
  {
    value: "api-key",
    label: "API Key",
    description: "Simple API key authentication",
    icon: "🔑",
    recommended: true
  },
  {
    value: "managed-identity",
    label: "Managed Identity",
    description: "Azure managed identity (no secrets)",
    icon: "🛡️",
    recommended: false
  },
  {
    value: "entra-id",
    label: "Entra ID (Service Principal)",
    description: "Service principal with client credentials",
    icon: "🏢",
    recommended: false
  }
]

const providers: ProviderOption[] = [
  {
    id: "openai",
    name: "OpenAI",
    description: "GPT-4, GPT-3.5 and other OpenAI models",
    icon: "🤖",
    popular: true
  },
  {
    id: "azure-openai",
    name: "Azure OpenAI",
    description: "OpenAI models hosted on Microsoft Azure",
    icon: "☁️",
    popular: true
  },
  {
    id: "anthropic",
    name: "Anthropic",
    description: "Claude 3 and other Anthropic models",
    icon: "🧠"
  },
  {
    id: "google",
    name: "Google AI",
    description: "Gemini and PaLM models from Google",
    icon: "🔍"
  },
  {
    id: "cohere",
    name: "Cohere",
    description: "Command and Embed models from Cohere",
    icon: "⚡"
  },
  {
    id: "custom",
    name: "Custom Provider",
    description: "Connect to your own API endpoint",
    icon: "🔧"
  }
]

// Form validation
const errors = ref<Record<string, string>>({})

const validateStep = (step: number): boolean => {
  errors.value = {}

  switch (step) {
    case 1:
      if (!formData.value.service) {
        errors.value.service = "Please select a provider"
        return false
      }
      break
    case 2:
      if (!formData.value.name.trim()) {
        errors.value.name = "Provider name is required"
        return false
      }
      if (!formData.value.slug.trim()) {
        errors.value.slug = "Slug is required"
        return false
      }
      if (!/^[a-z0-9-]+$/.test(formData.value.slug)) {
        errors.value.slug = "Slug can only contain lowercase letters, numbers, and hyphens"
        return false
      }
      if (!formData.value.environment) {
        errors.value.environment = "Please select an environment"
        return false
      }
      break
    case 3:
      // Authentication validation based on type
      if (formData.value.authType === "api-key") {
        if (!formData.value.apiKey?.trim()) {
          errors.value.apiKey = "API key is required"
          return false
        }
      } else if (formData.value.authType === "entra-id") {
        if (!formData.value.tenantId?.trim()) {
          errors.value.tenantId = "Tenant ID is required for Entra ID authentication"
          return false
        }
        if (!formData.value.clientId?.trim()) {
          errors.value.clientId = "Client ID is required for Entra ID authentication"
          return false
        }
        if (!formData.value.clientSecret?.trim()) {
          errors.value.clientSecret = "Client Secret is required for Entra ID authentication"
          return false
        }
      }

      // Azure OpenAI specific validation
      if (formData.value.service === "azure-openai") {
        if (!formData.value.resourceName?.trim() && !formData.value.customDomain?.trim()) {
          errors.value.resourceName = "Resource name or custom domain is required"
          return false
        }
        if (!formData.value.deploymentName?.trim()) {
          errors.value.deploymentName = "Deployment name is required for Azure OpenAI"
          return false
        }
      }

      if (formData.value.service === "custom" && !formData.value.customDomain?.trim()) {
        errors.value.customDomain = "Base URL is required for custom providers"
        return false
      }
      break
  }

  return true
}

const selectedProvider = computed(() => providers.find((p) => p.id === formData.value.service))

const canProceed = computed(() => {
  return validateStep(currentStep.value)
})

const getStepTitle = computed(() => {
  const titles = ["Choose your AI provider", "Configure integration", "Set up authentication", "Review & create"]
  return titles[currentStep.value - 1]
})

const getStepDescription = computed(() => {
  const descriptions = [
    "Select the AI service you want to integrate with your gateway",
    "Set up the integration details and scope",
    "Provide your API credentials for secure access",
    "Review all settings before creating the provider"
  ]
  return descriptions[currentStep.value - 1]
})

const generateSlug = (name: string) => {
  return name
    .toLowerCase()
    .replace(/[^a-z0-9\s]/g, "")
    .replace(/\s+/g, "-")
    .trim()
}

const onNameChange = () => {
  if (formData.value.name && !formData.value.slug) {
    formData.value.slug = generateSlug(formData.value.name)
  }
}

const selectProvider = (providerId: string) => {
  formData.value.service = providerId
  nextStep()
}

const nextStep = () => {
  if (!validateStep(currentStep.value)) {
    return
  }

  if (currentStep.value < maxSteps) {
    isTransitioning.value = true
    setTimeout(() => {
      currentStep.value++
      isTransitioning.value = false
    }, 150)
  }
}

const prevStep = () => {
  if (currentStep.value > 1) {
    isTransitioning.value = true
    setTimeout(() => {
      currentStep.value--
      isTransitioning.value = false
      errors.value = {} // Clear errors when going back
    }, 150)
  }
}

const closeModal = () => {
  emit("update:open", false)
  // Reset form after animation
  setTimeout(() => {
    currentStep.value = 1
    formData.value = {
      service: "",
      environment: "",
      name: "",
      slug: "",
      authType: "api-key",
      apiKey: "",
      tenantId: "",
      clientId: "",
      clientSecret: "",
      resourceName: "",
      deploymentName: "",
      customDomain: "",
      apiVersion: "2024-02-15-preview",
      region: ""
    }
  }, 300)
}

const createProvider = async () => {
  if (!validateStep(currentStep.value)) {
    return
  }

  loading.value = true
  try {
    // TODO: Replace with actual API call
    // const response = await $fetch('/api/providers', {
    //   method: 'POST',
    //   body: formData.value
    // })

    // Simulate API call
    await new Promise((resolve) => setTimeout(resolve, 2000))

    showSuccess.value = true

    // Wait for success animation then close
    setTimeout(() => {
      const providerId = `provider_${Date.now()}`
      emit("created", providerId)
      closeModal()
    }, 2000)
  } catch (error) {
    console.error("Failed to create provider:", error)
    errors.value.general = "Failed to create provider. Please try again."
  } finally {
    loading.value = false
  }
}

const grabApiKey = () => {
  const urls: Record<string, string> = {
    openai: "https://platform.openai.com/api-keys",
    "azure-openai": "https://portal.azure.com",
    anthropic: "https://console.anthropic.com",
    google: "https://makersuite.google.com/app/apikey",
    cohere: "https://dashboard.cohere.ai/api-keys"
  }

  const url = urls[formData.value.service]
  if (url) {
    window.open(url, "_blank")
  }
}
</script>

<template>
  <UiDialog :open="open" @update:open="$emit('update:open', $event)">
    <UiDialogContent class="min-w-5xl h-[85vh] p-0 overflow-y-auto">
      <!-- Success State -->
      <div
        v-if="showSuccess"
        class="absolute inset-0 bg-gradient-to-br from-chart-2/20 via-primary/10 to-chart-3/20 backdrop-blur-md z-50 flex items-center justify-center"
      >
        <div class="bg-background/95 backdrop-blur-sm rounded-2xl p-8 border shadow-2xl max-w-md mx-4">
          <div class="text-center space-y-6">
            <div class="relative">
              <div class="w-20 h-20 bg-chart-2/20 rounded-full flex items-center justify-center mx-auto">
                <div class="w-12 h-12 bg-chart-2 rounded-full flex items-center justify-center">
                  <Check class="w-6 h-6 text-white animate-in zoom-in-50 duration-700" />
                </div>
              </div>
              <div class="absolute -top-1 -right-1 w-6 h-6 bg-primary rounded-full flex items-center justify-center">
                <Sparkles class="w-3 h-3 text-white" />
              </div>
            </div>
            <div>
              <h3 class="text-xl font-bold mb-2">{{ selectedProvider?.name }} Connected!</h3>
              <p class="text-muted-foreground">Your provider is configured and ready to route AI requests.</p>
            </div>
          </div>
        </div>
      </div>

      <div class="flex h-full">
        <!-- Sidebar with Steps -->
        <div class="w-80 border-r bg-muted/30">
          <div class="p-6 border-b">
            <div class="flex items-center justify-between mb-2">
              <h1 class="text-xl font-bold">Add Provider</h1>
              <UiButton variant="ghost" size="sm" @click="closeModal" :disabled="loading">
                <X class="w-4 h-4" />
              </UiButton>
            </div>
            <p class="text-sm text-muted-foreground">Connect your AI service to the gateway</p>
          </div>

          <!-- Step Navigation -->
          <div class="p-6 space-y-4">
            <div v-for="(step, index) in maxSteps" :key="step" class="relative">
              <div class="flex items-center gap-4">
                <!-- Step Indicator -->
                <div
                  class="w-10 h-10 rounded-full flex items-center justify-center text-sm font-semibold transition-all duration-300 relative z-10"
                  :class="{
                    'bg-primary text-white shadow-lg': step === currentStep,
                    'bg-chart-2 text-white': step < currentStep,
                    'bg-muted-foreground/20 text-muted-foreground': step > currentStep
                  }"
                >
                  <Check v-if="step < currentStep" class="w-5 h-5" />
                  <div v-else-if="step === currentStep" class="w-2 h-2 bg-white rounded-full animate-pulse" />
                  <template v-else>{{ step }}</template>
                </div>

                <!-- Step Content -->
                <div class="flex-1">
                  <p
                    class="font-medium transition-colors duration-300"
                    :class="{
                      'text-foreground': step <= currentStep,
                      'text-muted-foreground': step > currentStep
                    }"
                  >
                    {{ ["Choose Service", "Configure", "Authenticate", "Review"][index] }}
                  </p>
                  <p
                    class="text-xs transition-colors duration-300"
                    :class="{
                      'text-muted-foreground': step <= currentStep,
                      'text-muted-foreground/60': step > currentStep
                    }"
                  >
                    {{ ["Select AI provider", "Set up integration", "Add credentials", "Confirm settings"][index] }}
                  </p>
                </div>
              </div>

              <!-- Connecting Line -->
              <div
                v-if="step < maxSteps"
                class="absolute left-5 top-10 w-0.5 h-8 transition-colors duration-300"
                :class="{
                  'bg-chart-2': step < currentStep,
                  'bg-muted-foreground/20': step >= currentStep
                }"
              />
            </div>
          </div>

          <!-- Selected Provider Preview -->
          <div v-if="selectedProvider" class="p-6 border-t bg-background/50">
            <div class="flex items-center gap-3">
              <div class="text-2xl">{{ selectedProvider.icon }}</div>
              <div>
                <p class="font-medium text-sm">{{ selectedProvider.name }}</p>
                <p class="text-xs text-muted-foreground">Selected provider</p>
              </div>
            </div>
          </div>
        </div>

        <!-- Main Content -->
        <div class="flex-1 flex flex-col h-full min-h-0">
          <!-- Header -->
          <div class="p-8 border-b bg-gradient-to-r from-background to-muted/20">
            <div class="max-w-2xl">
              <h2 class="text-2xl font-bold mb-2">{{ getStepTitle }}</h2>
              <p class="text-muted-foreground">{{ getStepDescription }}</p>
            </div>
          </div>

          <!-- Content -->
          <div class="flex-1 overflow-hidden p-8">
            <div class="max-w-2xl mx-auto">
              <!-- Global Error -->
              <div
                v-if="errors.general"
                class="mb-6 p-4 bg-destructive/10 border border-destructive/20 rounded-xl flex items-start gap-3"
              >
                <div class="w-5 h-5 rounded-full bg-destructive/20 flex items-center justify-center mt-0.5">
                  <X class="w-3 h-3 text-destructive" />
                </div>
                <p class="text-destructive text-sm">{{ errors.general }}</p>
              </div>

              <!-- Step Content -->
              <div :class="{ 'opacity-50 pointer-events-none': isTransitioning }" class="transition-all duration-300">
                <!-- Step 1: Provider Selection -->
                <div v-if="currentStep === 1" class="space-y-4">
                  <div class="grid gap-3">
                    <button
                      v-for="provider in providers"
                      :key="provider.id"
                      @click="selectProvider(provider.id)"
                      class="group relative p-6 rounded-2xl border-2 text-left transition-all duration-200 hover:shadow-lg hover:scale-[1.01]"
                      :class="{
                        'border-primary bg-primary/5 shadow-md': formData.service === provider.id,
                        'border-border bg-card hover:border-primary/30': formData.service !== provider.id
                      }"
                    >
                      <div class="flex items-center gap-4">
                        <div class="text-3xl">{{ provider.icon }}</div>
                        <div class="flex-1 min-w-0">
                          <div class="flex items-center gap-3 mb-2">
                            <h3 class="font-bold text-lg">{{ provider.name }}</h3>
                            <UiBadge
                              v-if="provider.popular"
                              class="bg-gradient-to-r from-amber-500 to-orange-500 text-white border-0 text-xs"
                            >
                              <Sparkles class="w-3 h-3 mr-1" />
                              Popular
                            </UiBadge>
                          </div>
                          <p class="text-muted-foreground text-sm">{{ provider.description }}</p>
                        </div>
                        <div class="flex items-center gap-2">
                          <div
                            v-if="formData.service === provider.id"
                            class="w-5 h-5 bg-primary rounded-full flex items-center justify-center"
                          >
                            <Check class="w-3 h-3 text-white" />
                          </div>
                          <ChevronRight
                            class="w-5 h-5 text-muted-foreground group-hover:text-primary transition-colors"
                            :class="{ 'text-primary': formData.service === provider.id }"
                          />
                        </div>
                      </div>
                    </button>
                  </div>
                </div>

                <!-- Step 2: Configuration -->
                <div v-if="currentStep === 2" class="space-y-8">
                  <!-- Environment Selection -->
                  <div>
                    <h3 class="font-semibold mb-4 flex items-center gap-2">
                      <Zap class="w-5 h-5 text-primary" />
                      Environment
                    </h3>
                    <div class="grid grid-cols-2 gap-4">
                      <button
                        v-for="env in environments"
                        :key="env.value"
                        @click="formData.environment = env.value"
                        class="p-6 rounded-xl border-2 text-left transition-all duration-200 hover:shadow-sm"
                        :class="{
                          'border-primary bg-primary/5': formData.environment === env.value,
                          'border-border hover:border-primary/30': formData.environment !== env.value
                        }"
                      >
                        <div class="flex items-center gap-3 mb-3">
                          <div class="w-10 h-10 rounded-lg bg-chart-1/20 flex items-center justify-center text-lg">
                            {{ env.icon }}
                          </div>
                          <div>
                            <h4 class="font-semibold">{{ env.label }}</h4>
                            <p class="text-xs text-muted-foreground">{{ env.description }}</p>
                          </div>
                        </div>
                      </button>
                    </div>
                    <div v-if="errors.environment" class="mt-2">
                      <p class="text-destructive text-xs flex items-center gap-2">
                        <X class="w-3 h-3" />
                        {{ errors.environment }}
                      </p>
                    </div>
                  </div>

                  <!-- Basic Configuration -->
                  <div class="space-y-6">
                    <div>
                      <label class="block text-sm font-semibold mb-3">Provider Name</label>
                      <UiInput
                        v-model="formData.name"
                        @input="onNameChange"
                        placeholder="e.g. OpenAI Production, Azure Development..."
                        class="h-12 text-base"
                        :class="{ 'border-destructive': errors.name }"
                      />
                      <div class="mt-2 space-y-1">
                        <p v-if="errors.name" class="text-destructive text-xs flex items-center gap-2">
                          <X class="w-3 h-3" />
                          {{ errors.name }}
                        </p>
                        <p class="text-xs text-muted-foreground">Choose a descriptive name to identify this provider</p>
                      </div>
                    </div>

                    <div>
                      <label class="block text-sm font-semibold mb-3">URL Slug</label>
                      <UiInput
                        v-model="formData.slug"
                        placeholder="provider-slug"
                        class="h-12 font-mono text-base"
                        :class="{ 'border-destructive': errors.slug }"
                      />
                      <div class="mt-2 space-y-1">
                        <p v-if="errors.slug" class="text-destructive text-xs flex items-center gap-2">
                          <X class="w-3 h-3" />
                          {{ errors.slug }}
                        </p>
                        <p class="text-xs text-muted-foreground">Used in API routing and configuration references</p>
                      </div>
                    </div>
                  </div>
                </div>

                <!-- Step 3: Authentication -->
                <div v-if="currentStep === 3" class="flex-1 flex flex-col min-h-0 overflow-hidden">
                  <UiScrollArea class="flex-1 min-h-0">
                    <div class="space-y-8 h-full pb-6">
                      <!-- Azure OpenAI Enhanced Configuration -->
                      <div v-if="formData.service === 'azure-openai'">
                        <!-- Authentication Type Selection -->
                        <div class="space-y-4">
                          <h3 class="font-semibold flex items-center gap-2">
                            <div class="w-8 h-8 rounded-lg bg-primary/20 flex items-center justify-center">🔐</div>
                            Authentication Method
                          </h3>

                          <div class="space-y-3">
                            <button
                              v-for="authType in azureAuthTypes"
                              :key="authType.value"
                              @click="formData.authType = authType.value"
                              class="w-full p-4 rounded-lg border-2 text-left transition-all duration-200 hover:shadow-sm"
                              :class="{
                                'border-primary bg-primary/5': formData.authType === authType.value,
                                'border-border hover:border-primary/30': formData.authType !== authType.value
                              }"
                            >
                              <div class="flex items-center justify-between">
                                <div class="flex items-center gap-3">
                                  <div class="text-xl">{{ authType.icon }}</div>
                                  <div>
                                    <div class="flex items-center gap-2">
                                      <h4 class="font-semibold">{{ authType.label }}</h4>
                                      <UiBadge v-if="authType.recommended" variant="secondary" class="text-xs">
                                        Recommended
                                      </UiBadge>
                                    </div>
                                    <p class="text-sm text-muted-foreground">{{ authType.description }}</p>
                                  </div>
                                </div>
                                <div
                                  v-if="formData.authType === authType.value"
                                  class="w-5 h-5 bg-primary rounded-full flex items-center justify-center"
                                >
                                  <Check class="w-3 h-3 text-white" />
                                </div>
                              </div>
                            </button>
                          </div>
                        </div>

                        <!-- API Key Fields -->
                        <div v-if="formData.authType === 'api-key'" class="space-y-4">
                          <div class="flex items-center justify-between">
                            <h4 class="font-semibold">API Key</h4>
                            <UiButton variant="outline" @click="grabApiKey" class="gap-2 h-9">
                              <ExternalLink class="w-4 h-4" />
                              Get API Key
                            </UiButton>
                          </div>
                          <UiInput
                            v-model="formData.apiKey"
                            type="password"
                            placeholder="Paste your Azure OpenAI API key"
                            class="h-12 font-mono"
                            :class="{ 'border-destructive': errors.apiKey }"
                          />
                          <div class="space-y-1">
                            <p v-if="errors.apiKey" class="text-destructive text-xs flex items-center gap-2">
                              <X class="w-3 h-3" />
                              {{ errors.apiKey }}
                            </p>
                            <p class="text-xs text-muted-foreground">
                              Found in Azure Portal → Your OpenAI Resource → Keys and Endpoint
                            </p>
                          </div>
                        </div>

                        <!-- Entra ID Fields -->
                        <div v-if="formData.authType === 'entra-id'" class="space-y-6">
                          <h4 class="font-semibold">Entra ID Service Principal</h4>
                          <div class="grid grid-cols-1 gap-4">
                            <div>
                              <label class="block text-sm font-semibold mb-2">Tenant ID</label>
                              <UiInput
                                v-model="formData.tenantId"
                                placeholder="12345678-1234-1234-1234-123456789012"
                                class="h-12 font-mono"
                                :class="{ 'border-destructive': errors.tenantId }"
                              />
                              <p v-if="errors.tenantId" class="text-destructive text-xs mt-1 flex items-center gap-2">
                                <X class="w-3 h-3" />
                                {{ errors.tenantId }}
                              </p>
                            </div>
                            <div>
                              <label class="block text-sm font-semibold mb-2">Client ID (Application ID)</label>
                              <UiInput
                                v-model="formData.clientId"
                                placeholder="12345678-1234-1234-1234-123456789012"
                                class="h-12 font-mono"
                                :class="{ 'border-destructive': errors.clientId }"
                              />
                              <p v-if="errors.clientId" class="text-destructive text-xs mt-1 flex items-center gap-2">
                                <X class="w-3 h-3" />
                                {{ errors.clientId }}
                              </p>
                            </div>
                            <div>
                              <label class="block text-sm font-semibold mb-2">Client Secret</label>
                              <UiInput
                                v-model="formData.clientSecret"
                                type="password"
                                placeholder="Your client secret value"
                                class="h-12 font-mono"
                                :class="{ 'border-destructive': errors.clientSecret }"
                              />
                              <p
                                v-if="errors.clientSecret"
                                class="text-destructive text-xs mt-1 flex items-center gap-2"
                              >
                                <X class="w-3 h-3" />
                                {{ errors.clientSecret }}
                              </p>
                            </div>
                          </div>
                        </div>

                        <!-- Managed Identity Info -->
                        <div
                          v-if="formData.authType === 'managed-identity'"
                          class="p-4 bg-blue-50 border border-blue-200 rounded-lg"
                        >
                          <div class="flex items-start gap-3">
                            <div class="w-6 h-6 bg-blue-100 rounded-full flex items-center justify-center mt-0.5">
                              <Check class="w-4 h-4 text-blue-600" />
                            </div>
                            <div>
                              <h4 class="font-semibold text-blue-900">Managed Identity Configuration</h4>
                              <p class="text-sm text-blue-700 mt-1">
                                No additional credentials needed. The system will use the managed identity assigned to
                                the deployment environment. Ensure your managed identity has the "Cognitive Services
                                OpenAI User" role on your Azure OpenAI resource.
                              </p>
                            </div>
                          </div>
                        </div>

                        <!-- Azure Resource Configuration -->
                        <div class="space-y-6">
                          <h4 class="font-semibold">Azure Resource Configuration</h4>
                          <div class="grid grid-cols-1 gap-4">
                            <div>
                              <label class="block text-sm font-semibold mb-2">Resource Name or Custom Domain</label>
                              <div class="space-y-3">
                                <UiInput
                                  v-model="formData.resourceName"
                                  placeholder="my-openai-resource"
                                  class="h-12"
                                  :class="{ 'border-destructive': errors.resourceName }"
                                />
                                <div class="flex items-center gap-2">
                                  <div class="flex-1 h-px bg-border"></div>
                                  <span class="text-xs text-muted-foreground">OR</span>
                                  <div class="flex-1 h-px bg-border"></div>
                                </div>
                                <UiInput
                                  v-model="formData.customDomain"
                                  placeholder="https://my-custom-domain.openai.azure.com"
                                  class="h-12"
                                />
                              </div>
                              <div class="mt-2 space-y-1">
                                <p v-if="errors.resourceName" class="text-destructive text-xs flex items-center gap-2">
                                  <X class="w-3 h-3" />
                                  {{ errors.resourceName }}
                                </p>
                                <p class="text-xs text-muted-foreground">
                                  Resource name will create: https://[resource-name].openai.azure.com
                                </p>
                              </div>
                            </div>

                            <div>
                              <label class="block text-sm font-semibold mb-2">Deployment Name</label>
                              <UiInput
                                v-model="formData.deploymentName"
                                placeholder="gpt-4-deployment"
                                class="h-12"
                                :class="{ 'border-destructive': errors.deploymentName }"
                              />
                              <div class="mt-2 space-y-1">
                                <p
                                  v-if="errors.deploymentName"
                                  class="text-destructive text-xs flex items-center gap-2"
                                >
                                  <X class="w-3 h-3" />
                                  {{ errors.deploymentName }}
                                </p>
                                <p class="text-xs text-muted-foreground">
                                  The name of your model deployment in Azure OpenAI Studio
                                </p>
                              </div>
                            </div>

                            <div class="grid grid-cols-2 gap-4">
                              <div>
                                <label class="block text-sm font-semibold mb-2">API Version</label>
                                <UiInput v-model="formData.apiVersion" placeholder="2024-02-15-preview" class="h-12" />
                                <p class="text-xs text-muted-foreground mt-1">Latest stable API version recommended</p>
                              </div>
                              <div>
                                <label class="block text-sm font-semibold mb-2">Region (Optional)</label>
                                <UiInput v-model="formData.region" placeholder="eastus" class="h-12" />
                              </div>
                            </div>
                          </div>
                        </div>
                      </div>

                      <!-- Other Providers Authentication -->
                      <div v-else>
                        <div class="flex items-center justify-between mb-4">
                          <h3 class="font-semibold flex items-center gap-2">
                            <div class="w-8 h-8 rounded-lg bg-primary/20 flex items-center justify-center">🔐</div>
                            API Authentication
                          </h3>
                          <UiButton variant="outline" @click="grabApiKey" class="gap-2 h-9">
                            <ExternalLink class="w-4 h-4" />
                            Get API Key
                          </UiButton>
                        </div>

                        <div class="space-y-3">
                          <UiInput
                            v-model="formData.apiKey"
                            type="password"
                            placeholder="sk-... or your API key"
                            class="h-12 font-mono text-base"
                            :class="{ 'border-destructive': errors.apiKey }"
                          />
                          <div class="space-y-1">
                            <p v-if="errors.apiKey" class="text-destructive text-xs flex items-center gap-2">
                              <X class="w-3 h-3" />
                              {{ errors.apiKey }}
                            </p>
                            <p class="text-xs text-muted-foreground">Your API key is encrypted and stored securely</p>
                          </div>
                        </div>

                        <!-- Custom Provider Configuration -->
                        <div v-if="formData.service === 'custom'" class="space-y-4 mt-6">
                          <h4 class="font-semibold">Custom Endpoint</h4>
                          <UiInput
                            v-model="formData.customDomain"
                            placeholder="https://api.example.com/v1"
                            class="h-12"
                            :class="{ 'border-destructive': errors.customDomain }"
                          />
                          <p v-if="errors.customDomain" class="text-destructive text-xs">{{ errors.customDomain }}</p>
                        </div>
                      </div>
                    </div>
                  </UiScrollArea>
                </div>

                <!-- Step 4: Review -->
                <div v-if="currentStep === 4" class="space-y-6">
                  <!-- Provider Summary Card -->
                  <div class="relative overflow-hidden">
                    <!-- Background Pattern -->
                    <div class="absolute inset-0 bg-gradient-to-br from-primary/5 via-chart-1/5 to-chart-2/5"></div>
                    <div class="absolute inset-0 bg-grid-pattern opacity-[0.03]"></div>

                    <!-- Main Card -->
                    <div
                      class="relative bg-white/50 backdrop-blur-sm border-2 border-primary/10 rounded-3xl p-8 shadow-xl"
                    >
                      <!-- Header Section -->
                      <div class="flex items-start justify-between mb-8">
                        <div class="flex items-center gap-4">
                          <div class="relative">
                            <div class="w-16 h-16 bg-primary/10 rounded-2xl flex items-center justify-center text-3xl">
                              {{ selectedProvider?.icon }}
                            </div>
                            <div
                              class="absolute -bottom-1 -right-1 w-6 h-6 bg-chart-2 rounded-full flex items-center justify-center"
                            >
                              <Check class="w-3 h-3 text-white" />
                            </div>
                          </div>
                          <div>
                            <h3 class="text-2xl font-bold mb-1">{{ formData.name }}</h3>
                            <p class="text-muted-foreground text-lg">{{ selectedProvider?.name }} Provider</p>
                            <div class="flex items-center gap-2 mt-2">
                              <div class="w-2 h-2 bg-chart-2 rounded-full animate-pulse"></div>
                              <span class="text-sm text-muted-foreground">Ready for deployment</span>
                            </div>
                          </div>
                        </div>

                        <!-- Environment Badge -->
                        <div class="flex items-center gap-2">
                          <div class="text-lg">
                            {{ environments.find((e) => e.value === formData.environment)?.icon }}
                          </div>
                          <UiBadge class="bg-gradient-to-r from-primary to-primary/80 text-white border-0 px-3 py-1">
                            {{ formData.environment }} Environment
                          </UiBadge>
                        </div>
                      </div>

                      <!-- Configuration Grid -->
                      <div class="grid grid-cols-1 md:grid-cols-3 gap-6 mb-8">
                        <!-- Basic Info -->
                        <div class="space-y-4">
                          <h4 class="text-sm font-bold text-muted-foreground uppercase tracking-wide mb-3">
                            Provider Details
                          </h4>
                          <div class="space-y-3">
                            <div>
                              <p class="text-xs text-muted-foreground mb-1">Type</p>
                              <div class="flex items-center gap-2">
                                <div class="w-2 h-2 bg-primary rounded-full"></div>
                                <span class="font-semibold">{{ selectedProvider?.name }}</span>
                              </div>
                            </div>
                            <div>
                              <p class="text-xs text-muted-foreground mb-1">URL Slug</p>
                              <code class="px-3 py-2 bg-muted/50 border rounded-lg text-sm font-mono">{{
                                formData.slug
                              }}</code>
                            </div>
                          </div>
                        </div>

                        <!-- Authentication -->
                        <div class="space-y-4">
                          <h4 class="text-sm font-bold text-muted-foreground uppercase tracking-wide mb-3">
                            Authentication
                          </h4>
                          <div class="space-y-3">
                            <div>
                              <p class="text-xs text-muted-foreground mb-1">Method</p>
                              <div class="flex items-center gap-2">
                                <div class="text-sm">
                                  {{ azureAuthTypes.find((a) => a.value === formData.authType)?.icon }}
                                </div>
                                <UiBadge variant="outline" class="font-medium">
                                  {{ azureAuthTypes.find((a) => a.value === formData.authType)?.label }}
                                </UiBadge>
                              </div>
                            </div>
                            <div v-if="formData.apiKey">
                              <p class="text-xs text-muted-foreground mb-1">API Key</p>
                              <div class="flex items-center gap-2">
                                <code class="px-3 py-2 bg-muted/50 border rounded-lg text-sm font-mono"
                                  >{{ formData.apiKey.slice(0, 8) }}••••••••</code
                                >
                                <div class="w-2 h-2 bg-chart-2 rounded-full"></div>
                              </div>
                            </div>
                            <div v-if="formData.tenantId">
                              <p class="text-xs text-muted-foreground mb-1">Tenant ID</p>
                              <code class="text-xs font-mono text-muted-foreground"
                                >{{ formData.tenantId.slice(0, 8) }}...</code
                              >
                            </div>
                          </div>
                        </div>

                        <!-- Azure Configuration -->
                        <div class="space-y-4" v-if="formData.service === 'azure-openai'">
                          <h4 class="text-sm font-bold text-muted-foreground uppercase tracking-wide mb-3">
                            Azure Config
                          </h4>
                          <div class="space-y-3">
                            <div v-if="formData.resourceName">
                              <p class="text-xs text-muted-foreground mb-1">Resource</p>
                              <p class="font-mono text-sm">{{ formData.resourceName }}.openai.azure.com</p>
                            </div>
                            <div v-if="formData.customDomain">
                              <p class="text-xs text-muted-foreground mb-1">Custom Domain</p>
                              <p class="font-mono text-sm">{{ formData.customDomain }}</p>
                            </div>
                            <div v-if="formData.deploymentName">
                              <p class="text-xs text-muted-foreground mb-1">Deployment</p>
                              <div class="flex items-center gap-2">
                                <div class="w-2 h-2 bg-chart-3 rounded-full"></div>
                                <span class="font-semibold">{{ formData.deploymentName }}</span>
                              </div>
                            </div>
                            <div v-if="formData.apiVersion">
                              <p class="text-xs text-muted-foreground mb-1">API Version</p>
                              <UiBadge variant="secondary" class="text-xs">{{ formData.apiVersion }}</UiBadge>
                            </div>
                            <div v-if="formData.region">
                              <p class="text-xs text-muted-foreground mb-1">Region</p>
                              <span class="font-medium">{{ formData.region }}</span>
                            </div>
                          </div>
                        </div>
                      </div>

                      <!-- Status Banner -->
                      <div
                        class="bg-gradient-to-r from-chart-2/10 via-chart-2/5 to-chart-2/10 border border-chart-2/20 rounded-2xl p-6"
                      >
                        <div class="flex items-center justify-between">
                          <div class="flex items-center gap-4">
                            <div class="relative">
                              <div class="w-12 h-12 bg-chart-2 rounded-xl flex items-center justify-center">
                                <Sparkles class="w-6 h-6 text-white" />
                              </div>
                              <div
                                class="absolute -top-1 -right-1 w-4 h-4 bg-green-500 rounded-full border-2 border-white"
                              ></div>
                            </div>
                            <div>
                              <h4 class="font-bold text-lg">Configuration Complete</h4>
                              <p class="text-muted-foreground">
                                Your {{ selectedProvider?.name }} provider is ready for immediate deployment
                              </p>
                            </div>
                          </div>
                          <div class="text-right">
                            <div class="text-2xl font-bold text-chart-2">100%</div>
                            <p class="text-xs text-muted-foreground">Setup Complete</p>
                          </div>
                        </div>
                      </div>
                    </div>
                  </div>
                </div>
              </div>
            </div>
          </div>

          <!-- Footer -->
          <div class="border-t bg-background p-6">
            <div class="max-w-2xl mx-auto flex items-center justify-between">
              <UiButton
                v-if="currentStep > 1"
                variant="ghost"
                @click="prevStep"
                :disabled="loading || isTransitioning"
                class="gap-2"
              >
                <ChevronLeft class="w-4 h-4" />
                Back
              </UiButton>
              <div v-else></div>

              <div class="flex items-center gap-4">
                <div class="text-xs text-muted-foreground font-medium">Step {{ currentStep }} of {{ maxSteps }}</div>

                <UiButton
                  v-if="currentStep < maxSteps"
                  @click="nextStep"
                  :disabled="!canProceed || loading || isTransitioning"
                  size="lg"
                  class="gap-2 min-w-24"
                >
                  Continue
                  <ChevronRight class="w-4 h-4" />
                </UiButton>

                <UiButton
                  v-else
                  @click="createProvider"
                  :disabled="!canProceed || loading"
                  size="lg"
                  class="gap-2 min-w-32 bg-gradient-to-r from-primary to-primary/80"
                >
                  <template v-if="loading">
                    <div class="w-4 h-4 border-2 border-white/30 border-t-white rounded-full animate-spin" />
                    Creating...
                  </template>
                  <template v-else>
                    <Sparkles class="w-4 h-4" />
                    Create Provider
                  </template>
                </UiButton>
              </div>
            </div>
          </div>
        </div>
      </div>
    </UiDialogContent>
  </UiDialog>
</template>
