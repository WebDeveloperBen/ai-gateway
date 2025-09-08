<script lang="ts">
interface PolicyRule {
  id: string
  condition: string
  action: string
  value: string | number
  enabled: boolean
}

interface PolicyApplication {
  id: string
  name: string
  status: "active" | "inactive"
  team: string
  appliedDate: string
}

interface PolicyAuditLog {
  timestamp: string
  action: string
  details: string
  user: string
}

interface PolicyDetail {
  id: string
  name: string
  description: string
  type: "rate_limit" | "content_filter" | "access_control" | "cost_limit"
  status: "active" | "inactive"
  priority: "high" | "medium" | "low"
  applicationsCount: number
  lastModified: string
  createdBy: string
  createdDate: string
  version: string
  rules: PolicyRule[]
  applications: PolicyApplication[]
  auditLog: PolicyAuditLog[]
  tags: string[]
  enforcement: "strict" | "warn" | "log"
}

const policies: Record<string, PolicyDetail> = {
  policy_1: {
    id: "policy_1",
    name: "Rate Limiting - Production",
    description:
      "Enforce strict rate limits for production applications to prevent abuse and ensure fair resource allocation across all applications.",
    type: "rate_limit",
    status: "active",
    priority: "high",
    applicationsCount: 8,
    lastModified: "2025-01-15T14:30:00Z",
    createdBy: "Alice Johnson",
    createdDate: "2024-12-01T10:00:00Z",
    version: "2.1",
    enforcement: "strict",
    tags: ["production", "rate-limiting", "security"],
    rules: [
      {
        id: "rule_1",
        condition: "requests_per_minute > 1000",
        action: "block",
        value: 1000,
        enabled: true
      },
      {
        id: "rule_2",
        condition: "requests_per_hour > 50000",
        action: "throttle",
        value: 50000,
        enabled: true
      },
      {
        id: "rule_3",
        condition: "concurrent_requests > 100",
        action: "queue",
        value: 100,
        enabled: false
      }
    ],
    applications: [
      {
        id: "app_1",
        name: "Customer Service Bot",
        status: "active",
        team: "Customer Success",
        appliedDate: "2024-12-15T09:00:00Z"
      },
      {
        id: "app_2",
        name: "Content Generator",
        status: "active",
        team: "Marketing",
        appliedDate: "2024-12-20T14:30:00Z"
      },
      {
        id: "app_5",
        name: "Analytics Dashboard",
        status: "active",
        team: "Product",
        appliedDate: "2025-01-10T11:15:00Z"
      }
    ],
    auditLog: [
      {
        timestamp: "2025-01-15T14:30:00Z",
        action: "Policy Updated",
        details: "Increased rate limit threshold from 800 to 1000 requests/minute",
        user: "Alice Johnson"
      },
      {
        timestamp: "2025-01-10T11:15:00Z",
        action: "Application Added",
        details: "Applied policy to Analytics Dashboard",
        user: "Bob Smith"
      },
      {
        timestamp: "2025-01-05T16:45:00Z",
        action: "Policy Triggered",
        details: "Rate limit exceeded for Customer Service Bot - 1,234 requests blocked",
        user: "System"
      },
      {
        timestamp: "2024-12-20T14:30:00Z",
        action: "Application Added",
        details: "Applied policy to Content Generator",
        user: "Alice Johnson"
      }
    ]
  },
  policy_2: {
    id: "policy_2",
    name: "Content Safety Filter",
    description:
      "Block harmful or inappropriate content across all AI interactions to maintain platform safety and compliance.",
    type: "content_filter",
    status: "active",
    priority: "high",
    applicationsCount: 12,
    lastModified: "2025-01-14T09:15:00Z",
    createdBy: "Bob Smith",
    createdDate: "2024-11-15T14:00:00Z",
    version: "1.8",
    enforcement: "strict",
    tags: ["content-safety", "compliance", "moderation"],
    rules: [
      {
        id: "rule_1",
        condition: "contains_profanity = true",
        action: "block",
        value: "enabled",
        enabled: true
      },
      {
        id: "rule_2",
        condition: "hate_speech_score > 0.8",
        action: "block",
        value: 0.8,
        enabled: true
      },
      {
        id: "rule_3",
        condition: "adult_content_detected = true",
        action: "warn",
        value: "enabled",
        enabled: true
      }
    ],
    applications: [],
    auditLog: []
  }
}
</script>
<script setup lang="ts">
import {
  Shield,
  Clock,
  Filter,
  Zap,
  AlertTriangle,
  CheckCircle,
  XCircle,
  Settings,
  Edit,
  Play,
  Pause,
  Eye,
  Activity,
  FileText,
  Users,
  Calendar,
  User,
  Tag,
  GitBranch,
  MoreHorizontal,
  Plus,
  Trash2,
  RotateCcw,
  Download,
  Upload
} from "lucide-vue-next"
import type { StatsCardProps } from "@/components/Cards/Stats.vue"

const route = useRoute()
const policyId = route.params.id as string

const policy = computed(() => policies[policyId])

useSeoMeta({
  title: computed(() => (policy.value ? `${policy.value.name} - Governance Policy` : "Policy Not Found"))
})

if (!policy.value) {
  throw createError({
    statusCode: 404,
    statusMessage: "Policy not found"
  })
}

const showSettingsModal = ref(false)
const isSettingsLoading = ref(false)
const showAssignApplicationModal = ref(false)

const handlePolicyToggle = async () => {
  console.log(`Toggling policy ${policyId}`)
}

const openSettingsModal = () => {
  showSettingsModal.value = true
}

const handleSettingsSave = async (updatedPolicy: PolicyDetail) => {
  isSettingsLoading.value = true
  try {
    console.log("Policy updated:", updatedPolicy)
    isSettingsLoading.value = false
    showSettingsModal.value = false
  } catch (error) {
    console.error("Failed to update policy:", error)
  } finally {
    isSettingsLoading.value = false
  }
}

const handleSettingsCancel = () => {
  showSettingsModal.value = false
}

const openAssignApplicationModal = () => {
  showAssignApplicationModal.value = true
}

const handleApplicationAssigned = (applicationId: string) => {
  console.log(`Application ${applicationId} assigned to policy ${policyId}`)
  // TODO: Refresh policy data or update applications list
}

const getPolicyTypeIcon = (type: PolicyDetail["type"]) => {
  switch (type) {
    case "rate_limit":
      return Clock
    case "content_filter":
      return Filter
    case "access_control":
      return Shield
    case "cost_limit":
      return Zap
    default:
      return Shield
  }
}

const getPolicyTypeColor = (type: PolicyDetail["type"]) => {
  switch (type) {
    case "rate_limit":
      return "text-blue-600 bg-blue-50 border-blue-200"
    case "content_filter":
      return "text-purple-600 bg-purple-50 border-purple-200"
    case "access_control":
      return "text-green-600 bg-green-50 border-green-200"
    case "cost_limit":
      return "text-orange-600 bg-orange-50 border-orange-200"
    default:
      return "text-gray-600 bg-gray-50 border-gray-200"
  }
}

const getPriorityColor = (priority: PolicyDetail["priority"]) => {
  switch (priority) {
    case "high":
      return "text-red-600 bg-red-50 border-red-200"
    case "medium":
      return "text-yellow-600 bg-yellow-50 border-yellow-200"
    case "low":
      return "text-green-600 bg-green-50 border-green-200"
    default:
      return "text-gray-600 bg-gray-50 border-gray-200"
  }
}

const getStatusColor = (status: PolicyDetail["status"]) => {
  switch (status) {
    case "active":
      return "text-green-600 bg-green-50 border-green-200"
    case "inactive":
      return "text-gray-600 bg-gray-50 border-gray-200"
    default:
      return "text-gray-600 bg-gray-50 border-gray-200"
  }
}

const getEnforcementColor = (enforcement: PolicyDetail["enforcement"]) => {
  switch (enforcement) {
    case "strict":
      return "text-red-600 bg-red-50 border-red-200"
    case "warn":
      return "text-yellow-600 bg-yellow-50 border-yellow-200"
    case "log":
      return "text-blue-600 bg-blue-50 border-blue-200"
    default:
      return "text-gray-600 bg-gray-50 border-gray-200"
  }
}

const getApplicationStatusColor = (status: PolicyApplication["status"]) => {
  return status === "active"
    ? "text-green-600 bg-green-50 border-green-200"
    : "text-gray-600 bg-gray-50 border-gray-200"
}

const statsIcons = [Users, getPolicyTypeIcon(policy.value.type), AlertTriangle, Activity]
const statsVariants: StatsCardProps["variant"][] = ["chart-1", "chart-2", "chart-3", "default"]

const statsCards = computed(() => {
  return [
    {
      title: "Applied Applications",
      value: policy.value.applicationsCount,
      description: "Currently using this policy"
    },
    {
      title: "Policy Type",
      value: policy.value.type.replace("_", " ").replace(/\b\w/g, (l) => l.toUpperCase()),
      description: "Governance category"
    },
    {
      title: "Priority Level",
      value: policy.value.priority.charAt(0).toUpperCase() + policy.value.priority.slice(1),
      description: "Enforcement priority"
    },
    {
      title: "Active Rules",
      value: policy.value.rules.filter((rule) => rule.enabled).length,
      description: `of ${policy.value.rules.length} total rules`
    }
  ].map((stat, index) => ({
    ...stat,
    icon: statsIcons[index],
    variant: statsVariants[index]
  }))
})
</script>

<template>
  <div v-if="policy" class="flex flex-col gap-6">
    <PageBadgeHeader :badge-status="policy.status" :title="policy.name" :subtext="policy.description">
      <div class="flex items-center gap-2">
        <UiButton variant="outline" size="sm" class="gap-2">
          <Download class="h-4 w-4" />
          Export
        </UiButton>
        <UiButton
          :variant="policy.status === 'active' ? 'destructive' : 'default'"
          size="sm"
          class="gap-2"
          @click="handlePolicyToggle"
        >
          <component :is="policy.status === 'active' ? Pause : Play" class="h-4 w-4" />
          {{ policy.status === "active" ? "Disable" : "Enable" }}
        </UiButton>
        <UiButton variant="default" class="gap-2" @click="openSettingsModal">
          <Settings class="h-4 w-4" />
          Settings
        </UiButton>
      </div>
    </PageBadgeHeader>

    <div class="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-4 gap-6">
      <CardsStats
        v-for="card in statsCards"
        :key="card.title"
        :title="card.title"
        :value="card.value"
        :icon="card.icon"
        :description="card.description"
        :variant="card.variant"
      />
    </div>

    <div class="grid grid-cols-1 lg:grid-cols-2 gap-6">
      <UiCard>
        <UiCardHeader>
          <UiCardTitle class="flex items-center gap-2">
            <FileText class="h-5 w-5" />
            Policy Details
          </UiCardTitle>
        </UiCardHeader>
        <UiCardContent class="space-y-4">
          <div class="grid grid-cols-2 gap-4">
            <div>
              <label class="text-sm font-medium text-muted-foreground flex items-center gap-2">
                <component :is="getPolicyTypeIcon(policy.type)" class="h-4 w-4" />
                Type
              </label>
              <div class="mt-1">
                <div
                  class="inline-flex items-center gap-1 px-2 py-1 rounded-md text-xs font-medium border"
                  :class="getPolicyTypeColor(policy.type)"
                >
                  <component :is="getPolicyTypeIcon(policy.type)" class="size-3" />
                  {{ policy.type.replace("_", " ").replace(/\b\w/g, (l) => l.toUpperCase()) }}
                </div>
              </div>
            </div>

            <div>
              <label class="text-sm font-medium text-muted-foreground flex items-center gap-2">
                <AlertTriangle class="h-4 w-4" />
                Priority
              </label>
              <div class="mt-1">
                <div
                  class="inline-flex items-center gap-1 px-2 py-1 rounded-md text-xs font-medium border"
                  :class="getPriorityColor(policy.priority)"
                >
                  <component :is="policy.priority === 'high' ? AlertTriangle : FileText" class="size-3" />
                  {{ policy.priority.charAt(0).toUpperCase() + policy.priority.slice(1) }}
                </div>
              </div>
            </div>

            <div>
              <label class="text-sm font-medium text-muted-foreground flex items-center gap-2">
                <Shield class="h-4 w-4" />
                Enforcement
              </label>
              <div class="mt-1">
                <div
                  class="inline-flex items-center gap-1 px-2 py-1 rounded-md text-xs font-medium border"
                  :class="getEnforcementColor(policy.enforcement)"
                >
                  <Shield class="size-3" />
                  {{ policy.enforcement.charAt(0).toUpperCase() + policy.enforcement.slice(1) }}
                </div>
              </div>
            </div>

            <div>
              <label class="text-sm font-medium text-muted-foreground flex items-center gap-2">
                <GitBranch class="h-4 w-4" />
                Version
              </label>
              <p class="text-sm font-medium mt-1">{{ policy.version }}</p>
            </div>
          </div>

          <div class="grid grid-cols-2 gap-4">
            <div>
              <label class="text-sm font-medium text-muted-foreground flex items-center gap-2">
                <Calendar class="h-4 w-4" />
                Created
              </label>
              <p class="text-sm font-medium mt-1">{{ new Date(policy.createdDate).toLocaleDateString() }}</p>
            </div>

            <div>
              <label class="text-sm font-medium text-muted-foreground flex items-center gap-2">
                <User class="h-4 w-4" />
                Created By
              </label>
              <p class="text-sm font-medium mt-1">{{ policy.createdBy }}</p>
            </div>
          </div>

          <div>
            <label class="text-sm font-medium text-muted-foreground flex items-center gap-2">
              <Tag class="h-4 w-4" />
              Tags
            </label>
            <div class="mt-1 flex flex-wrap gap-1">
              <UiBadge v-for="tag in policy.tags" :key="tag" variant="secondary" class="text-xs">
                {{ tag }}
              </UiBadge>
            </div>
          </div>
        </UiCardContent>
      </UiCard>

      <UiCard>
        <UiCardHeader>
          <UiCardTitle class="flex items-center gap-2">
            <Settings class="h-5 w-5" />
            Policy Rules
          </UiCardTitle>
          <UiCardDescription>{{ policy.rules.length }} rules configured</UiCardDescription>
        </UiCardHeader>
        <UiCardContent>
          <div class="space-y-3">
            <div
              v-for="rule in policy.rules"
              :key="rule.id"
              class="p-3 border rounded-lg"
              :class="rule.enabled ? 'bg-muted/20 border-muted' : 'bg-muted/10 border-muted/50'"
            >
              <div class="flex items-center justify-between mb-2">
                <div class="flex items-center gap-2">
                  <div class="size-2 rounded-full" :class="rule.enabled ? 'bg-primary' : 'bg-muted-foreground/50'" />
                  <span class="text-sm font-medium">{{ rule.condition }}</span>
                </div>
                <UiBadge :variant="rule.enabled ? 'default' : 'secondary'" class="text-xs">
                  {{ rule.enabled ? "Active" : "Disabled" }}
                </UiBadge>
              </div>
              <div class="text-xs text-muted-foreground">
                <strong>Action:</strong> {{ rule.action }} <strong class="ml-2">Value:</strong> {{ rule.value }}
              </div>
            </div>
          </div>
        </UiCardContent>
      </UiCard>
    </div>

    <CardsDataList title="Applied Applications" :icon="Users">
      <template #actions>
        <UiButton variant="outline" size="sm" class="gap-2" @click="openAssignApplicationModal">
          <Plus class="h-4 w-4" />
          Add Application
        </UiButton>
      </template>

      <div class="space-y-4">
        <div
          v-for="app in policy.applications"
          :key="app.id"
          class="flex items-center justify-between p-4 border rounded-lg hover:bg-muted/50 cursor-pointer"
          @click="navigateTo(`/applications/${app.id}`)"
        >
          <div class="flex items-center gap-4">
            <div class="p-2 rounded-lg bg-primary/10">
              <Shield class="size-5 text-primary" />
            </div>

            <div class="space-y-1">
              <p class="font-medium">{{ app.name }}</p>
              <div class="flex items-center gap-4 text-sm text-muted-foreground">
                <span>{{ app.team }}</span>
                <span>•</span>
                <span>Applied {{ new Date(app.appliedDate).toLocaleDateString() }}</span>
              </div>
            </div>
          </div>

          <div class="flex items-center gap-2">
            <div
              class="inline-flex items-center gap-1 px-2 py-1 rounded-md text-xs font-medium border"
              :class="getApplicationStatusColor(app.status)"
            >
              <div class="size-1.5 rounded-full" :class="app.status === 'active' ? 'bg-green-500' : 'bg-gray-400'" />
              {{ app.status === "active" ? "Active" : "Inactive" }}
            </div>

            <UiDropdownMenu @click.stop>
              <UiDropdownMenuTrigger as-child>
                <UiButton variant="ghost" size="sm">
                  <MoreHorizontal class="size-4" />
                </UiButton>
              </UiDropdownMenuTrigger>
              <UiDropdownMenuContent align="end">
                <UiDropdownMenuItem>
                  <Eye class="mr-2 size-4" />
                  View Application
                </UiDropdownMenuItem>
                <UiDropdownMenuSeparator />
                <UiDropdownMenuItem class="text-destructive">
                  <Trash2 class="mr-2 size-4" />
                  Remove Policy
                </UiDropdownMenuItem>
              </UiDropdownMenuContent>
            </UiDropdownMenu>
          </div>
        </div>

        <div v-if="policy.applications.length === 0" class="text-center py-8 text-muted-foreground">
          No applications currently using this policy
        </div>
      </div>
    </CardsDataList>

    <!-- Assign Application Modal -->
    <LazyModalsPolicyAssignApplication
      v-model:open="showAssignApplicationModal"
      :policy-id="policyId"
      :policy-name="policy.name"
      @assigned="handleApplicationAssigned"
    />
  </div>
</template>
