<script setup lang="ts">
import {
  ArrowLeft,
  Key,
  Edit,
  RotateCw,
  Trash2,
  Activity,
  Calendar,
  Shield,
  Zap,
  AlertTriangle,
  Clock,
  Globe,
  MoreVertical,
  Code,
  TestTube,
  Copy
} from "lucide-vue-next"

// Get the key ID from route params
const route = useRoute()
const keyId = route.params.id as string

// Sample API key data - this would come from an API call
const apiKey = ref({
  id: keyId,
  key: "",
  name: "Production API Key",
  keyPrefix: "sk-",
  applicationId: "app_1",
  applicationName: "Customer Service Bot",
  description: "Main production key for customer service bot integration",
  status: "active",
  permissions: ["read", "write"],
  rateLimit: "1000/hour",
  created: "2024-12-15T10:00:00Z",
  lastUsed: "2025-01-15T14:30:00Z",
  expiresAt: null,
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

// Actions

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
type CodeExampleKey = 'curl' | 'javascript' | 'python'

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
    // Could add toast notification here
  } catch (err) {
    console.error("Failed to copy code example: ", err)
  }
}
</script>

<template>
  <div class="flex flex-col gap-6">
    <!-- Breadcrumb Navigation -->
    <div class="flex items-center gap-2 text-sm text-muted-foreground">
      <NuxtLink to="/applications" class="hover:text-foreground">Applications</NuxtLink>
      <span>/</span>
      <NuxtLink to="/applications/keys" class="hover:text-foreground">API Keys</NuxtLink>
      <span>/</span>
      <span class="text-foreground">{{ apiKey.name }}</span>
    </div>

    <!-- Header -->
    <div class="flex items-start justify-between">
      <div class="flex items-start gap-4">
        <UiButton variant="ghost" size="sm" as-child class="gap-2">
          <NuxtLink to="/applications/keys">
            <ArrowLeft class="h-4 w-4" />
            Back to Keys
          </NuxtLink>
        </UiButton>

        <div>
          <div class="flex items-center gap-3 mb-2">
            <h1 class="text-3xl font-bold tracking-tight">{{ apiKey.name }}</h1>
            <UiBadge :class="getStatusBadgeClass(apiKey.status)">
              {{ apiKey.status }}
            </UiBadge>
          </div>
          <p class="text-muted-foreground">
            Used by <strong>{{ apiKey.applicationName }}</strong>
          </p>
        </div>
      </div>

      <!-- Action Buttons -->
      <div class="flex items-center gap-2">
        <UiButton variant="outline" size="sm" class="gap-2">
          <Edit class="h-4 w-4" />
          Edit Settings
        </UiButton>
        <UiButton variant="outline" size="sm" class="gap-2">
          <RotateCw class="h-4 w-4" />
          Regenerate
        </UiButton>
        <UiDropdownMenu>
          <UiDropdownMenuTrigger as-child>
            <UiButton variant="outline" size="sm">
              <MoreVertical class="h-4 w-4" />
            </UiButton>
          </UiDropdownMenuTrigger>
          <UiDropdownMenuContent align="end">
            <UiDropdownMenuItem>
              <TestTube class="h-4 w-4 mr-2" />
              Test API Key
            </UiDropdownMenuItem>
            <UiDropdownMenuSeparator />
            <UiDropdownMenuItem class="text-red-600">
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
            <ApiKeyDisplay :key-id="apiKey.id" :key-prefix="apiKey.keyPrefix" size="lg" />
          </div>
        </div>

        <!-- Key Details Grid -->
        <div class="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-4 gap-6">
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
              <Shield class="h-4 w-4" />
              Permissions
            </label>
            <p class="text-sm font-medium mt-1">{{ apiKey.permissions.join(", ") }}</p>
          </div>

          <div>
            <label class="text-sm font-medium text-muted-foreground flex items-center gap-2">
              <Zap class="h-4 w-4" />
              Rate Limit
            </label>
            <p class="text-sm font-medium mt-1">{{ apiKey.rateLimit }}</p>
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

    <!-- Charts Placeholder -->
    <div class="grid grid-cols-1 lg:grid-cols-2 gap-6">
      <UiCard>
        <UiCardHeader>
          <UiCardTitle>Usage Over Time</UiCardTitle>
          <UiCardDescription>API requests in the last 30 days</UiCardDescription>
        </UiCardHeader>
        <UiCardContent>
          <div class="h-[300px] flex items-center justify-center border-2 border-dashed border-muted rounded-lg">
            <div class="text-center text-muted-foreground">
              <Activity class="h-12 w-12 mx-auto mb-2 opacity-50" />
              <p>Usage Chart Coming Soon</p>
            </div>
          </div>
        </UiCardContent>
      </UiCard>

      <UiCard>
        <UiCardHeader>
          <UiCardTitle>Status Code Breakdown</UiCardTitle>
          <UiCardDescription>Response status distribution</UiCardDescription>
        </UiCardHeader>
        <UiCardContent>
          <div class="h-[300px] flex items-center justify-center border-2 border-dashed border-muted rounded-lg">
            <div class="text-center text-muted-foreground">
              <Activity class="h-12 w-12 mx-auto mb-2 opacity-50" />
              <p>Status Chart Coming Soon</p>
            </div>
          </div>
        </UiCardContent>
      </UiCard>
    </div>

    <!-- Recent Activity -->
    <UiCard>
      <UiCardHeader>
        <UiCardTitle class="flex items-center gap-2">
          <Globe class="h-5 w-5" />
          Recent Activity
        </UiCardTitle>
        <UiCardDescription>Latest API requests using this key</UiCardDescription>
      </UiCardHeader>
      <UiCardContent>
        <div class="space-y-4">
          <div
            v-for="activity in apiKey.recentActivity"
            :key="activity.timestamp"
            class="flex items-center justify-between p-3 border rounded-lg"
          >
            <div class="flex items-center gap-3">
              <div class="w-2 h-2 rounded-full bg-green-500"></div>
              <div>
                <div class="flex items-center gap-2">
                  <span class="font-medium text-sm">{{ activity.method }}</span>
                  <code class="text-sm text-muted-foreground">{{ activity.endpoint }}</code>
                </div>
                <div class="text-xs text-muted-foreground">
                  {{ new Date(activity.timestamp).toLocaleString() }} • {{ activity.ip }}
                </div>
              </div>
            </div>
            <UiBadge variant="outline" :class="getStatusColor(activity.status)">
              {{ activity.status }}
            </UiBadge>
          </div>
        </div>
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
            v-for="lang in (['curl', 'javascript', 'python'] as CodeExampleKey[])"
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
  </div>
</template>
