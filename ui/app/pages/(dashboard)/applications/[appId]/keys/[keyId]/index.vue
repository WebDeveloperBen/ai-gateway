<script setup lang="ts">
import {
  Key,
  RotateCw,
  Trash2,
  Activity,
  Calendar,
  AlertTriangle,
  Clock,
  MoreVertical,
  Code,
  TestTube,
  Copy,
  Users,
  Ban
} from "lucide-vue-next"
import { toast } from "vue-sonner"

// Get the app ID and key ID from route params
const route = useRoute()
const appId = route.params.appId as string
const keyId = route.params.keyId as string

// Modal state management
const showDisableModal = ref(false)
const showDeleteModal = ref(false)
const actionLoading = ref(false)

// Actions
const handleDisableKey = async () => {
  actionLoading.value = true
  try {
    // TODO: Implement API call to disable key
    await new Promise((resolve) => setTimeout(resolve, 1000)) // Simulate API call
    apiKey.value.status = "inactive"
    showDisableModal.value = false
  } catch (error) {
    console.error("Failed to disable key:", error)
  } finally {
    actionLoading.value = false
  }
}

const handleDeleteKey = async () => {
  actionLoading.value = true
  try {
    // TODO: Implement API call to delete key
    await new Promise((resolve) => setTimeout(resolve, 1000)) // Simulate API call
    // In real app, this would redirect to keys list after deletion
    showDeleteModal.value = false
    await navigateTo(`/applications/${appId}/keys`)
  } catch (error) {
    console.error("Failed to delete key:", error)
  } finally {
    actionLoading.value = false
  }
}

// Sample API key data - this would come from an API call
const apiKey = ref({
  id: keyId,
  key: "",
  name: "Production API Key",
  keyPrefix: "sk-",
  applicationId: appId,
  applicationName: "Customer Service Bot",
  description: "Main production key for customer service bot integration",
  status: "active",
  permissions: ["read", "write"],
  created: "2024-12-15T10:00:00Z",
  lastUsed: "2025-01-15T14:30:00Z",
  expiresAt: null,
  owners: [
    {
      id: "user_1",
      name: "John Smith",
      email: "john.smith@company.com",
      role: "Admin"
    },
    {
      id: "user_2",
      name: "Sarah Johnson",
      email: "sarah.johnson@company.com",
      role: "Developer"
    }
  ],
  requestCount: {
    total: 25400,
    today: 240,
    thisWeek: 1680,
    thisMonth: 7200
  },
  errorRate: 2.1,
  recentActivity: [
    {
      timestamp: "2025-01-15T14:30:00Z",
      method: "POST",
      endpoint: "/v1/chat/completions",
      status: 200,
      ip: "192.168.1.100"
    },
    {
      timestamp: "2025-01-15T14:25:00Z",
      method: "POST",
      endpoint: "/v1/chat/completions",
      status: 200,
      ip: "192.168.1.100"
    },
    {
      timestamp: "2025-01-15T14:20:00Z",
      method: "GET",
      endpoint: "/v1/models",
      status: 200,
      ip: "192.168.1.100"
    }
  ]
})

const formatNumber = (num: number) => {
  return new Intl.NumberFormat().format(num)
}

const getStatusBadgeClass = (status: string) => {
  return status === "active"
    ? "bg-green-100 text-green-800 dark:bg-green-900 dark:text-green-200"
    : "bg-gray-100 text-gray-800 dark:bg-gray-900 dark:text-gray-200"
}

const getStatusColor = (status: number) => {
  if (status >= 200 && status < 300) return "text-green-600"
  if (status >= 400 && status < 500) return "text-yellow-600"
  if (status >= 500) return "text-red-600"
  return "text-gray-600"
}

// Code examples for integration
type CodeExampleKey = "curl" | "javascript" | "python"

const codeExamples: Record<CodeExampleKey, string> = {
  curl: `curl -X POST "https://api.yourdomain.com/v1/chat/completions" \\
  -H "Authorization: Bearer YOUR_API_KEY_HERE" \\
  -H "Content-Type: application/json" \\
  -d '{
    "model": "gpt-4",
    "messages": [{"role": "user", "content": "Hello!"}]
  }'`,

  javascript: `const response = await fetch('https://api.yourdomain.com/v1/chat/completions', {
  method: 'POST',
  headers: {
    'Authorization': 'Bearer YOUR_API_KEY_HERE',
    'Content-Type': 'application/json',
  },
  body: JSON.stringify({
    model: 'gpt-4',
    messages: [{ role: 'user', content: 'Hello!' }]
  })
});`,

  python: `import requests

headers = {
    'Authorization': 'Bearer YOUR_API_KEY_HERE',
    'Content-Type': 'application/json',
}

data = {
    'model': 'gpt-4',
    'messages': [{'role': 'user', 'content': 'Hello!'}]
}

response = requests.post('https://api.yourdomain.com/v1/chat/completions', 
                        headers=headers, json=data)`
}

const selectedCodeExample = ref<CodeExampleKey>("curl")

const copyCodeExample = async (example: CodeExampleKey) => {
  try {
    await navigator.clipboard.writeText(codeExamples[example])
    toast.success("Code copied successfully!", {})
    // Could add toast notification here
  } catch (err) {
    console.error("Failed to copy code example: ", err)
  }
}
</script>

<template>
  <div class="flex flex-col gap-6">
    <!-- Header -->
    <div class="flex items-start justify-between">
      <div>
        <div class="flex items-center gap-3 mb-2">
          <h1 class="text-3xl font-bold tracking-tight">{{ apiKey.name }}</h1>
          <UiBadge :class="getStatusBadgeClass(apiKey.status)">
            {{ apiKey.status }}
          </UiBadge>
        </div>
        <p class="text-muted-foreground">
          Used by
          <NuxtLink :to="`/applications/${apiKey.applicationId}`" class="font-medium text-foreground hover:underline">{{
            apiKey.applicationName
          }}</NuxtLink>
        </p>
      </div>

      <!-- Action Buttons -->
      <div class="flex items-center gap-2">
        <UiButton variant="default" size="sm" class="gap-2" @click="$router.push('/playground/prompts')">
          <TestTube class="h-4 w-4" />
          Test API Key
        </UiButton>
        <UiButton variant="outline" size="sm" class="gap-2 text-orange-600" @click="showDisableModal = true">
          <Ban class="h-4 w-4" />
          Disable Key
        </UiButton>
        <UiDropdownMenu>
          <UiDropdownMenuTrigger as-child>
            <UiButton variant="outline" size="sm">
              <MoreVertical class="h-4 w-4" />
            </UiButton>
          </UiDropdownMenuTrigger>
          <UiDropdownMenuContent align="end">
            <UiDropdownMenuItem>
              <RotateCw class="h-4 w-4 mr-2" />
              Regenerate Key
            </UiDropdownMenuItem>
            <UiDropdownMenuSeparator />
            <UiDropdownMenuItem class="text-red-600" @click="showDeleteModal = true">
              <Trash2 class="h-4 w-4 mr-2" />
              Delete Key
            </UiDropdownMenuItem>
          </UiDropdownMenuContent>
        </UiDropdownMenu>
      </div>
    </div>

    <!-- Key Information Card -->
    <UiCard>
      <UiCardHeader>
        <UiCardTitle class="flex items-center gap-2">
          <Key class="h-5 w-5" />
          API Key Information
        </UiCardTitle>
      </UiCardHeader>
      <UiCardContent class="space-y-6">
        <!-- API Key Display -->
        <div>
          <label class="text-sm font-medium text-muted-foreground">API Key</label>
          <div class="mt-2">
            <ApiKeyDisplay :key-id="apiKey.id" :key-prefix="apiKey.keyPrefix" size="lg" show-copy-text />
          </div>
        </div>

        <!-- Key Details Grid -->
        <div class="grid grid-cols-1 md:grid-cols-3 gap-6">
          <div>
            <label class="text-sm font-medium text-muted-foreground flex items-center gap-2">
              <Calendar class="h-4 w-4" />
              Created
            </label>
            <p class="text-sm font-medium mt-1">{{ new Date(apiKey.created).toLocaleDateString() }}</p>
          </div>

          <div>
            <label class="text-sm font-medium text-muted-foreground flex items-center gap-2">
              <Clock class="h-4 w-4" />
              Last Used
            </label>
            <p class="text-sm font-medium mt-1">{{ new Date(apiKey.lastUsed).toLocaleDateString() }}</p>
          </div>

          <div>
            <label class="text-sm font-medium text-muted-foreground flex items-center gap-2">
              <Users class="h-4 w-4" />
              Owners
            </label>
            <div class="mt-1 flex flex-wrap gap-1">
              <div v-for="owner in apiKey.owners" :key="owner.id" class="flex items-center gap-1 text-xs">
                <div class="w-5 h-5 bg-primary/10 rounded-full flex items-center justify-center">
                  <span class="text-[10px] font-medium text-primary">{{
                    owner.name
                      .split(" ")
                      .map((n) => n[0])
                      .join("")
                  }}</span>
                </div>
                <span class="text-sm font-medium">{{ owner.name }}</span>
              </div>
            </div>
          </div>
        </div>

        <!-- Description -->
        <div v-if="apiKey.description">
          <label class="text-sm font-medium text-muted-foreground">Description</label>
          <p class="text-sm mt-1">{{ apiKey.description }}</p>
        </div>
      </UiCardContent>
    </UiCard>

    <!-- Usage Analytics Dashboard -->
    <div class="grid grid-cols-1 lg:grid-cols-4 gap-6">
      <!-- Usage Stats -->
      <UiCard>
        <UiCardHeader class="flex flex-row items-center justify-between space-y-0 pb-2">
          <UiCardTitle class="text-sm font-medium">Total Requests</UiCardTitle>
          <Activity class="h-4 w-4 text-muted-foreground" />
        </UiCardHeader>
        <UiCardContent>
          <div class="text-2xl font-bold">{{ formatNumber(apiKey.requestCount.total) }}</div>
          <p class="text-xs text-muted-foreground">+{{ apiKey.requestCount.today }} today</p>
        </UiCardContent>
      </UiCard>

      <UiCard>
        <UiCardHeader class="flex flex-row items-center justify-between space-y-0 pb-2">
          <UiCardTitle class="text-sm font-medium">This Week</UiCardTitle>
          <Calendar class="h-4 w-4 text-muted-foreground" />
        </UiCardHeader>
        <UiCardContent>
          <div class="text-2xl font-bold">{{ formatNumber(apiKey.requestCount.thisWeek) }}</div>
          <p class="text-xs text-muted-foreground">Weekly usage</p>
        </UiCardContent>
      </UiCard>

      <UiCard>
        <UiCardHeader class="flex flex-row items-center justify-between space-y-0 pb-2">
          <UiCardTitle class="text-sm font-medium">This Month</UiCardTitle>
          <Calendar class="h-4 w-4 text-muted-foreground" />
        </UiCardHeader>
        <UiCardContent>
          <div class="text-2xl font-bold">{{ formatNumber(apiKey.requestCount.thisMonth) }}</div>
          <p class="text-xs text-muted-foreground">Monthly usage</p>
        </UiCardContent>
      </UiCard>

      <UiCard>
        <UiCardHeader class="flex flex-row items-center justify-between space-y-0 pb-2">
          <UiCardTitle class="text-sm font-medium">Error Rate</UiCardTitle>
          <AlertTriangle class="h-4 w-4 text-muted-foreground" />
        </UiCardHeader>
        <UiCardContent>
          <div class="text-2xl font-bold">{{ apiKey.errorRate }}%</div>
          <p class="text-xs text-muted-foreground">Last 24 hours</p>
        </UiCardContent>
      </UiCard>
    </div>

    <!-- Usage Analytics Chart -->
    <UiCard>
      <UiCardHeader>
        <UiCardTitle class="flex items-center gap-2">
          <Activity class="h-5 w-5" />
          Usage Analytics
        </UiCardTitle>
        <UiCardDescription>Request volume and performance metrics over time</UiCardDescription>
      </UiCardHeader>
      <UiCardContent>
        <ChartPlaceholder
          :icon="Activity"
          title="Charts Coming Soon"
          description="Interactive usage analytics and performance charts will be available here"
        />
      </UiCardContent>
    </UiCard>

    <!-- Integration Helper -->
    <UiCard>
      <UiCardHeader>
        <UiCardTitle class="flex items-center gap-2">
          <Code class="h-5 w-5" />
          Integration Examples
        </UiCardTitle>
        <UiCardDescription>Sample code to get you started with this API key</UiCardDescription>
      </UiCardHeader>
      <UiCardContent>
        <!-- Language Tabs -->
        <div class="flex items-center gap-2 mb-4">
          <UiButton
            v-for="lang in ['curl', 'javascript', 'python'] as CodeExampleKey[]"
            :key="lang"
            variant="outline"
            size="sm"
            @click="selectedCodeExample = lang"
            :class="selectedCodeExample === lang ? 'bg-primary text-primary-foreground' : ''"
          >
            {{ lang.charAt(0).toUpperCase() + lang.slice(1) }}
          </UiButton>
        </div>

        <!-- Code Block -->
        <div class="relative">
          <pre
            class="bg-muted p-4 rounded-lg text-sm overflow-x-auto"
          ><code>{{ codeExamples[selectedCodeExample]}}</code></pre>
          <UiButton
            variant="outline"
            size="sm"
            class="absolute top-2 right-2 gap-2"
            @click="copyCodeExample(selectedCodeExample)"
          >
            <Copy class="h-3 w-3" />
            Copy
          </UiButton>
        </div>
      </UiCardContent>
    </UiCard>

    <!-- Confirmation Modals -->
    <ConfirmationModal
      v-model:open="showDisableModal"
      title="Disable API Key?"
      :description="`Are you sure you want to disable '${apiKey.name}'? This will prevent it from being used for API requests.`"
      confirm-text="Disable Key"
      variant="default"
      :loading="actionLoading"
      @confirm="handleDisableKey"
    />

    <ConfirmationModal
      v-model:open="showDeleteModal"
      title="Delete API Key?"
      :description="`Are you sure you want to permanently delete '${apiKey.name}'? This action cannot be undone and will break any applications currently using this key.`"
      confirm-text="Delete Key"
      variant="destructive"
      :loading="actionLoading"
      @confirm="handleDeleteKey"
    />
  </div>
</template>
