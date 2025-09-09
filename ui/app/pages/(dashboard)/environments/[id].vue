<script lang="ts">
interface Environment {
  id: string
  name: string
  description: string
  owner: string
  teams: string[]
  status: "active" | "inactive"
  memberCount: number
  applicationCount: number
  monthlyRequests: number
  createdAt: string
  lastActivity: string
}

interface TeamMember {
  id: string
  name: string
  email: string
  role: "Owner" | "Admin" | "Developer" | "Viewer"
  team: string
  lastActive: string
}

interface Application {
  id: string
  name: string
  description: string
  status: "active" | "inactive"
  apiKeyCount: number
  monthlyRequests: number
  lastUsed: string
  models: string[]
}

interface ApiKey {
  id: string
  name: string
  type: "Production" | "Development" | "Testing"
  status: "active" | "inactive" | "expired"
  lastUsed: string
  createdAt: string
  usageLimit: number
  currentUsage: number
}

interface ActivityItem {
  id: string
  type: "key_created" | "app_deployed" | "user_added" | "config_changed" | "policy_updated"
  description: string
  user: string
  timestamp: string
}

// Sample data - would come from API in real implementation
const environment = ref<Environment>({
  id: "env_1",
  name: "Production",
  description: "Live production environment for customer-facing applications",
  owner: "Alice Johnson",
  teams: ["Engineering", "DevOps"],
  status: "active",
  memberCount: 15,
  applicationCount: 8,
  monthlyRequests: 2500000,
  createdAt: "2024-10-01T08:00:00Z",
  lastActivity: "2025-01-15T14:30:00Z"
})

const teamMembers = ref<TeamMember[]>([
  {
    id: "1",
    name: "Alice Johnson",
    email: "alice@company.com",
    role: "Owner",
    team: "Engineering",
    lastActive: "2025-01-15T14:30:00Z"
  },
  {
    id: "2",
    name: "Bob Smith",
    email: "bob@company.com",
    role: "Admin",
    team: "Engineering",
    lastActive: "2025-01-15T12:15:00Z"
  },
  {
    id: "3",
    name: "Carol Williams",
    email: "carol@company.com",
    role: "Developer",
    team: "DevOps",
    lastActive: "2025-01-15T10:45:00Z"
  },
  {
    id: "4",
    name: "David Brown",
    email: "david@company.com",
    role: "Developer",
    team: "Engineering",
    lastActive: "2025-01-14T16:20:00Z"
  },
  {
    id: "5",
    name: "Emma Davis",
    email: "emma@company.com",
    role: "Admin",
    team: "DevOps",
    lastActive: "2025-01-15T09:30:00Z"
  }
])

const applications = ref<Application[]>([
  {
    id: "app_1",
    name: "Customer Service Bot",
    description: "AI-powered customer support assistant",
    status: "active",
    apiKeyCount: 3,
    monthlyRequests: 450000,
    lastUsed: "2025-01-15T14:30:00Z",
    models: ["gpt-4", "gpt-3.5-turbo"]
  },
  {
    id: "app_2",
    name: "Content Generator",
    description: "Automated content creation system",
    status: "active",
    apiKeyCount: 2,
    monthlyRequests: 128000,
    lastUsed: "2025-01-14T16:45:00Z",
    models: ["gpt-4"]
  },
  {
    id: "app_3",
    name: "Rag Pipeline",
    description: "Automated content creation system",
    status: "active",
    apiKeyCount: 2,
    monthlyRequests: 128000,
    lastUsed: "2025-01-14T16:45:00Z",
    models: ["gpt-4"]
  }
])

const apiKeys = ref<ApiKey[]>([
  {
    id: "key_1",
    name: "Prod-CustomerBot-Primary",
    type: "Production",
    status: "active",
    lastUsed: "2025-01-15T14:30:00Z",
    createdAt: "2024-10-01T08:00:00Z",
    usageLimit: 1000000,
    currentUsage: 450000
  },
  {
    id: "key_2",
    name: "Prod-ContentGen-Primary",
    type: "Production",
    status: "active",
    lastUsed: "2025-01-14T16:45:00Z",
    createdAt: "2024-10-15T10:30:00Z",
    usageLimit: 500000,
    currentUsage: 128000
  }
])

const activityFeed = ref<ActivityItem[]>([
  {
    id: "1",
    type: "app_deployed",
    description: "Deployed Customer Service Bot v2.1.0",
    user: "Alice Johnson",
    timestamp: "2025-01-15T14:30:00Z"
  },
  {
    id: "2",
    type: "user_added",
    description: "Added Emma Davis to DevOps team",
    user: "Bob Smith",
    timestamp: "2025-01-15T09:30:00Z"
  },
  {
    id: "3",
    type: "key_created",
    description: "Created new API key 'Prod-ContentGen-Backup'",
    user: "Carol Williams",
    timestamp: "2025-01-14T16:20:00Z"
  }
])
</script>

<script setup lang="ts">
import {
  Globe,
  Users,
  Key,
  Layers,
  Activity,
  Clock,
  Settings,
  Shield,
  BarChart3,
  UserPlus,
  Plus,
  Edit,
  Copy,
  RefreshCw
} from "lucide-vue-next"
import type { StatsCardProps } from "@/components/Cards/Stats.vue"

useSeoMeta({
  title: `${environment.value.name} - Environment Details`,
  description: `Manage ${environment.value.name} environment settings, teams, and applications`
})

// Stats for the environment
const environmentStats: StatsCardProps[] = [
  {
    title: "Total Requests",
    value: formatNumber(environment.value.monthlyRequests),
    icon: Activity,
    description: "This month",
    variant: "chart-1"
  },
  {
    title: "Applications",
    value: environment.value.applicationCount.toString(),
    icon: Layers,
    description: `${applications.value.filter((app) => app.status === "active").length} active`,
    variant: "chart-2"
  },
  {
    title: "API Keys",
    value: apiKeys.value.length.toString(),
    icon: Key,
    description: `${apiKeys.value.filter((key) => key.status === "active").length} active`,
    variant: "chart-3"
  },
  {
    title: "Team Members",
    value: environment.value.memberCount.toString(),
    icon: Users,
    description: `${environment.value.teams.length} teams`,
    variant: "default"
  }
]

function formatNumber(num: number): string {
  if (num >= 1000000) return (num / 1000000).toFixed(1) + "M"
  if (num >= 1000) return (num / 1000).toFixed(1) + "K"
  return num.toString()
}

function getActivityIcon(type: string) {
  switch (type) {
    case "key_created":
      return Key
    case "app_deployed":
      return Layers
    case "user_added":
      return UserPlus
    case "config_changed":
      return Settings
    case "policy_updated":
      return Shield
    default:
      return Activity
  }
}

function formatDate(dateString: string): string {
  return new Date(dateString).toLocaleString()
}

function getUsagePercentage(current: number, limit: number): number {
  return Math.round((current / limit) * 100)
}

const showCloneModal = ref(false)
const showManageAccessModal = ref(false)

// Helper functions for the new team members section
function getInitials(name: string): string {
  return name.split(' ').map(n => n[0]).join('').toUpperCase()
}

function formatTimeAgo(dateString: string): string {
  const date = new Date(dateString)
  const now = new Date()
  const diffInHours = Math.floor((now.getTime() - date.getTime()) / (1000 * 60 * 60))
  
  if (diffInHours < 1) return 'just now'
  if (diffInHours < 24) return `${diffInHours}h ago`
  const diffInDays = Math.floor(diffInHours / 24)
  if (diffInDays < 7) return `${diffInDays}d ago`
  return `${Math.floor(diffInDays / 7)}w ago`
}

function getRoleBadgeClass(role: string): string {
  switch (role) {
    case 'Owner':
      return 'bg-purple-100 text-purple-700 border-purple-200 dark:bg-purple-900/20 dark:text-purple-400'
    case 'Admin':
      return 'bg-red-100 text-red-700 border-red-200 dark:bg-red-900/20 dark:text-red-400'
    case 'Developer':
      return 'bg-blue-100 text-blue-700 border-blue-200 dark:bg-blue-900/20 dark:text-blue-400'
    case 'Viewer':
      return 'bg-gray-100 text-gray-700 border-gray-200 dark:bg-gray-900/20 dark:text-gray-400'
    default:
      return 'bg-gray-100 text-gray-700 border-gray-200 dark:bg-gray-900/20 dark:text-gray-400'
  }
}

function handleViewAnalytics() {
  navigateTo({
    path: "/analytics",
    query: {
      environment: environment.value.id,
      name: environment.value.name,
      timeRange: "30d"
    }
  })
}

function handleCloneEnvironment(data: {
  name: string
  description: string
  includeTeams: boolean
  includeKeys: boolean
}) {
  console.log("Cloning environment with data:", data)

  // TODO: Implement API call to clone environment
  // This would create a new environment with:
  // - New ID and the provided name/description
  // - Copy team assignments if includeTeams is true
  // - Create new API keys with similar naming if includeKeys is true
  // - No applications initially (user needs to deploy separately)

  // For now, just show success and redirect
  // navigateTo(`/environments/${newEnvironmentId}`)
}

function handleAccessUpdated(data: { environmentId: string, access: any[] }) {
  console.log("Environment access updated:", data)
  showManageAccessModal.value = false
  
  // In a real application, you would:
  // 1. Call an API to update environment access
  // 2. Refresh the team members data
  // 3. Show a success message
  
  // For now, just update the member count (simplified)
  environment.value.memberCount = data.access.filter(a => a.type === 'user').length + 
    data.access.filter(a => a.type === 'team').reduce((sum, a) => {
      const team = data.access.filter(t => t.type === 'team')
      return sum + team.length // Simplified counting
    }, 0)
}
</script>

<template>
  <div class="flex flex-col gap-6">
    <!-- Environment Header -->
    <div class="flex items-start justify-between">
      <div class="flex items-center gap-4">
        <div class="p-3 rounded-lg bg-primary/10">
          <Globe class="size-8 text-primary" />
        </div>
        <div>
          <div class="flex items-center gap-3 mb-2">
            <h1 class="text-2xl font-bold">{{ environment.name }}</h1>
            <div
              class="inline-flex items-center gap-1 px-2 py-1 rounded-md text-xs font-medium border"
              :class="getStatusColor(environment.status)"
            >
              <div
                class="size-1.5 rounded-full"
                :class="{
                  'bg-emerald-500': environment.status === 'active',
                  'bg-gray-400': environment.status === 'inactive'
                }"
              />
              {{ environment.status.charAt(0).toUpperCase() + environment.status.slice(1) }}
            </div>
          </div>
          <p class="text-muted-foreground mb-2">{{ environment.description }}</p>
          <div class="flex items-center gap-4 text-sm text-muted-foreground">
            <span>Owner: {{ environment.owner }}</span>
            <span>Created: {{ formatDate(environment.createdAt) }}</span>
            <span>Last Activity: {{ formatDate(environment.lastActivity) }}</span>
          </div>
        </div>
      </div>

      <div class="flex items-center gap-2">
        <UiButton variant="outline" size="sm" @click="showCloneModal = true">
          <Copy class="size-4 mr-2" />
          Clone Environment
        </UiButton>
        <UiButton variant="outline" size="sm" @click="handleViewAnalytics">
          <BarChart3 class="size-4 mr-2" />
          View Analytics
        </UiButton>
        <UiButton size="sm">
          <Edit class="size-4 mr-2" />
          Edit Environment
        </UiButton>
      </div>
    </div>

    <!-- Stats Overview -->
    <div class="grid gap-4 md:grid-cols-2 lg:grid-cols-4">
      <CardsStats
        v-for="stat in environmentStats"
        :key="stat.title"
        :title="stat.title"
        :value="stat.value"
        :icon="stat.icon"
        :description="stat.description"
        :variant="stat.variant"
      />
    </div>

    <!-- Main Content Grid - 2x2 Layout (4 cards, 2 per row) -->
    <div class="grid gap-6 lg:grid-cols-2">
      <!-- Applications (Top Left) -->
      <CardsDataList title="Applications" :icon="Layers">
        <template #actions>
          <UiButton @click="navigateTo('/applications')" size="sm" variant="outline">
            <Plus class="size-4" />
          </UiButton>
        </template>

        <div class="space-y-3 max-h-60 overflow-y-auto">
          <div
            v-for="app in applications"
            :key="app.id"
            class="flex items-center justify-between p-4 border rounded-lg hover:bg-muted/50 cursor-pointer transition-colors"
            @click="navigateTo(`/applications/${app.id}`)"
          >
            <div class="flex items-center gap-3">
              <div class="p-2 rounded-md bg-primary/10">
                <Layers class="size-4 text-primary" />
              </div>
              <div>
                <div class="font-medium">{{ app.name }}</div>
                <div class="text-sm text-muted-foreground">{{ app.description }}</div>
                <div class="flex items-center gap-4 text-xs text-muted-foreground mt-1">
                  <span>{{ app.apiKeyCount }} API keys</span>
                  <span>{{ formatNumber(app.monthlyRequests) }} req/month</span>
                </div>
              </div>
            </div>
            <div
              class="inline-flex items-center gap-1 px-2 py-1 rounded-md text-xs font-medium border"
              :class="getStatusColor(app.status)"
            >
              {{ app.status.charAt(0).toUpperCase() + app.status.slice(1) }}
            </div>
          </div>
        </div>
      </CardsDataList>

      <!-- API Keys (Bottom Left) -->
      <CardsDataList title="API Keys" :icon="Key">
        <div class="space-y-3">
          <div v-for="key in apiKeys" :key="key.id" class="flex items-center justify-between p-4 border rounded-lg">
            <div class="flex items-center gap-3">
              <div class="p-2 rounded-md bg-amber-100 dark:bg-amber-900/30">
                <Key class="size-4 text-amber-600 dark:text-amber-400" />
              </div>
              <div>
                <div class="font-medium">{{ key.name }}</div>
                <div class="text-sm text-muted-foreground">
                  {{ key.type }} • Created {{ formatDate(key.createdAt) }}
                </div>
                <div class="flex items-center gap-2 mt-2">
                  <div class="w-24 bg-muted rounded-full h-2">
                    <div
                      class="bg-primary h-2 rounded-full"
                      :style="{ width: `${getUsagePercentage(key.currentUsage, key.usageLimit)}%` }"
                    />
                  </div>
                  <span class="text-xs text-muted-foreground">
                    {{ getUsagePercentage(key.currentUsage, key.usageLimit) }}%
                  </span>
                </div>
              </div>
            </div>
            <div class="flex items-center gap-2">
              <div
                class="inline-flex items-center gap-1 px-2 py-1 rounded-md text-xs font-medium border"
                :class="getStatusColor(key.status)"
              >
                {{ key.status.charAt(0).toUpperCase() + key.status.slice(1) }}
              </div>
            </div>
          </div>
        </div>
      </CardsDataList>

      <!-- Team Members (Top Right) -->
      <CardsDataList title="Team Members" :icon="Users">
        <template #actions>
          <UiButton @click="showManageAccessModal = true" size="sm" variant="outline">
            <UserPlus class="size-4" />
          </UiButton>
        </template>

        <div class="space-y-3 max-h-60 overflow-y-auto">
          <div
            v-for="member in teamMembers"
            :key="member.id"
            class="flex items-center justify-between p-4 border rounded-lg hover:bg-muted/50 transition-colors"
          >
            <div class="flex items-center gap-3">
              <UiAvatar class="size-8">
                <UiAvatarFallback>{{ getInitials(member.name) }}</UiAvatarFallback>
              </UiAvatar>
              <div>
                <div class="font-medium">{{ member.name }}</div>
                <div class="text-sm text-muted-foreground">{{ member.email }}</div>
                <div class="flex items-center gap-2 text-xs text-muted-foreground mt-1">
                  <span>{{ member.team }}</span>
                  <span>•</span>
                  <span>Active {{ formatTimeAgo(member.lastActive) }}</span>
                </div>
              </div>
            </div>
            <div class="flex items-center gap-2">
              <UiBadge 
                variant="outline" 
                :class="getRoleBadgeClass(member.role)"
              >
                {{ member.role }}
              </UiBadge>
            </div>
          </div>
        </div>
      </CardsDataList>

      <!-- Recent Activity (Bottom Right) -->
      <CardsDataList title="Recent Activity" :icon="Activity">
        <div class="space-y-3 max-h-60 overflow-y-auto">
          <div
            v-for="activity in recentActivity"
            :key="activity.id"
            class="flex items-start gap-3 p-3 border rounded-lg"
          >
            <div class="p-2 rounded-md bg-blue-100 dark:bg-blue-900/30 mt-1">
              <Activity class="size-3 text-blue-600 dark:text-blue-400" />
            </div>
            <div class="flex-1 min-w-0">
              <div class="text-sm font-medium">{{ activity.description }}</div>
              <div class="text-xs text-muted-foreground mt-1">
                {{ formatTimeAgo(activity.timestamp) }}
              </div>
            </div>
          </div>
        </div>
      </CardsDataList>
    </div>

    <!-- Analytics Section - Full Width at Bottom for Lazy Loading -->
    <CardsAnalyticsComingSoon />

    <!-- Modals -->
    <ModalsEnvironmentsClone
      v-model:open="showCloneModal"
      :source-environment="environment"
      @cloned="handleCloneEnvironment"
    />

    <ModalsEnvironmentsManageAccess
      v-model:open="showManageAccessModal"
      :environment="environment"
      @updated="handleAccessUpdated"
    />
  </div>
</template>
