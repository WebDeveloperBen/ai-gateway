<script lang="ts">
import type { DataSourceData } from "@/models/datasource"

interface SyncLog {
  id: string
  timestamp: string
  status: "success" | "error" | "warning" | "info"
  message: string
  documentsProcessed?: number
  documentsAdded?: number
  documentsUpdated?: number
  documentsRemoved?: number
  duration?: number
  errorDetails?: string
}

interface DocumentMetrics {
  totalDocuments: number
  indexedToday: number
  averageDocumentSize: number
  lastIndexed: string
  categories: { name: string; count: number; percentage: number }[]
}

interface PerformanceMetrics {
  syncDuration: { avg: number; min: number; max: number }
  throughput: { docsPerSecond: number; bytesPerSecond: number }
  reliability: { successRate: number; uptime: number }
  costs: { monthly: number; perDocument: number }
}

// Sample data - in real implementation this would come from API based on route params
const dataSource = ref<DataSourceData>({
  id: "ds_1",
  name: "Customer Documentation",
  description: "Product documentation and user guides for customer support",
  type: "documentation",
  source: "confluence",
  url: "https://company.atlassian.net/wiki",
  status: "active",
  schedule: "0 2 * * *",
  lastSync: "2025-01-15T02:00:00Z",
  nextSync: "2025-01-16T02:00:00Z",
  documentsCount: 1250,
  owner: "Alice Johnson",
  tags: ["documentation", "support", "customer"]
})

const syncLogs = ref<SyncLog[]>([
  {
    id: "log_1",
    timestamp: "2025-01-15T02:00:00Z",
    status: "success",
    message: "Sync completed successfully",
    documentsProcessed: 1250,
    documentsAdded: 15,
    documentsUpdated: 8,
    documentsRemoved: 2,
    duration: 280
  },
  {
    id: "log_2",
    timestamp: "2025-01-14T02:00:00Z",
    status: "success",
    message: "Sync completed successfully",
    documentsProcessed: 1235,
    documentsAdded: 12,
    documentsUpdated: 5,
    documentsRemoved: 1,
    duration: 245
  },
  {
    id: "log_3",
    timestamp: "2025-01-13T02:00:00Z",
    status: "warning",
    message: "Sync completed with warnings",
    documentsProcessed: 1224,
    documentsAdded: 8,
    documentsUpdated: 3,
    documentsRemoved: 0,
    duration: 320,
    errorDetails: "3 documents failed to index due to unsupported format"
  },
  {
    id: "log_4",
    timestamp: "2025-01-12T02:00:00Z",
    status: "error",
    message: "Sync failed - Authentication error",
    documentsProcessed: 0,
    duration: 15,
    errorDetails: "Invalid API credentials. Please check authentication settings."
  },
  {
    id: "log_5",
    timestamp: "2025-01-11T02:00:00Z",
    status: "success",
    message: "Sync completed successfully",
    documentsProcessed: 1216,
    documentsAdded: 6,
    documentsUpdated: 12,
    documentsRemoved: 3,
    duration: 195
  }
])

const documentMetrics = ref<DocumentMetrics>({
  totalDocuments: 1250,
  indexedToday: 15,
  averageDocumentSize: 2.4,
  lastIndexed: "2025-01-15T02:00:00Z",
  categories: [
    { name: "User Guides", count: 450, percentage: 36 },
    { name: "API Documentation", count: 325, percentage: 26 },
    { name: "Troubleshooting", count: 275, percentage: 22 },
    { name: "Release Notes", count: 125, percentage: 10 },
    { name: "FAQs", count: 75, percentage: 6 }
  ]
})

const performanceMetrics = ref<PerformanceMetrics>({
  syncDuration: { avg: 240, min: 195, max: 320 },
  throughput: { docsPerSecond: 5.2, bytesPerSecond: 12.8 },
  reliability: { successRate: 96.2, uptime: 99.8 },
  costs: { monthly: 24.5, perDocument: 0.019 }
})

// Chart data for sync performance over time
const syncPerformanceData = ref([
  { date: "Jan 11", duration: 195, documents: 1216, success: true },
  { date: "Jan 12", duration: 15, documents: 0, success: false },
  { date: "Jan 13", duration: 320, documents: 1224, success: true },
  { date: "Jan 14", duration: 245, documents: 1235, success: true },
  { date: "Jan 15", duration: 280, documents: 1250, success: true }
])
</script>

<script setup lang="ts">
import {
  Database,
  FileText,
  Calendar,
  Clock,
  Activity,
  Settings,
  RefreshCw,
  Play,
  Pause,
  Edit,
  Copy,
  Trash2,
  AlertTriangle,
  CheckCircle,
  XCircle,
  TrendingUp,
  TrendingDown,
  Zap,
  DollarSign,
  BarChart3,
  PieChart,
  Download,
  Upload,
  Globe,
  Tag,
  User,
  ArrowUpRight,
  ArrowDownRight,
  MoreHorizontal,
  ExternalLink,
  GitBranch,
  Bookmark,
  Folder,
  HardDrive
} from "lucide-vue-next"
import type { StatsCardProps } from "@/components/Cards/Stats.vue"

// Get the route parameter
const route = useRoute()
const dataSourceId = route.params.id as string

useSeoMeta({
  title: `${dataSource.value.name} - Data Source Details`,
  description: `Manage ${dataSource.value.name} data source sync settings, logs, and performance metrics`
})

// Modal states
const showEditModal = ref(false)
const showDeleteModal = ref(false)
const showConfigModal = ref(false)
const showScheduleModal = ref(false)

// Stats for the data source
const dataSourceStats: StatsCardProps[] = [
  {
    title: "Total Documents",
    value: documentMetrics.value.totalDocuments.toLocaleString(),
    icon: FileText,
    description: `+${documentMetrics.value.indexedToday} today`,
    variant: "chart-1"
  },
  {
    title: "Sync Success Rate",
    value: `${performanceMetrics.value.reliability.successRate}%`,
    icon: CheckCircle,
    description: "Last 30 days",
    variant: "chart-2"
  },
  {
    title: "Avg Sync Duration",
    value: `${performanceMetrics.value.syncDuration.avg}s`,
    icon: Clock,
    description: "Processing time",
    variant: "chart-3"
  },
  {
    title: "Monthly Cost",
    value: `$${performanceMetrics.value.costs.monthly}`,
    icon: DollarSign,
    description: `$${performanceMetrics.value.costs.perDocument}/doc`,
    variant: "default"
  }
]

// Helper functions
function getStatusColor(status: string) {
  switch (status) {
    case "active":
      return "text-emerald-600 dark:text-emerald-400 bg-emerald-50 dark:bg-emerald-950/30 border-emerald-200 dark:border-emerald-800/50"
    case "paused":
      return "text-amber-600 dark:text-amber-400 bg-amber-50 dark:bg-amber-950/30 border-amber-200 dark:border-amber-800/50"
    case "error":
      return "text-red-600 dark:text-red-400 bg-red-50 dark:bg-red-950/30 border-red-200 dark:border-red-800/50"
    case "inactive":
      return "text-gray-600 dark:text-gray-400 bg-gray-50 dark:bg-gray-900/30 border-gray-200 dark:border-gray-700/50"
    default:
      return "text-gray-600 dark:text-gray-400 bg-gray-50 dark:bg-gray-900/30 border-gray-200 dark:border-gray-700/50"
  }
}

function getSyncStatusColor(status: string) {
  switch (status) {
    case "success":
      return "text-emerald-600 dark:text-emerald-400 bg-emerald-50 dark:bg-emerald-950/30 border-emerald-200 dark:border-emerald-800/50"
    case "warning":
      return "text-amber-600 dark:text-amber-400 bg-amber-50 dark:bg-amber-950/30 border-amber-200 dark:border-amber-800/50"
    case "error":
      return "text-red-600 dark:text-red-400 bg-red-50 dark:bg-red-950/30 border-red-200 dark:border-red-800/50"
    case "info":
      return "text-blue-600 dark:text-blue-400 bg-blue-50 dark:bg-blue-950/30 border-blue-200 dark:border-blue-800/50"
    default:
      return "text-gray-600 dark:text-gray-400 bg-gray-50 dark:bg-gray-900/30 border-gray-200 dark:border-gray-700/50"
  }
}

function getSyncStatusIcon(status: string) {
  switch (status) {
    case "success":
      return CheckCircle
    case "warning":
      return AlertTriangle
    case "error":
      return XCircle
    case "info":
      return Activity
    default:
      return Activity
  }
}

function getTypeIcon(type: string) {
  switch (type) {
    case "documentation":
      return FileText
    case "knowledge-base":
      return Bookmark
    case "policies":
      return Settings
    case "api-docs":
      return GitBranch
    default:
      return FileText
  }
}

function getSourceIcon(source: string) {
  switch (source) {
    case "confluence":
      return Globe
    case "notion":
      return Bookmark
    case "sharepoint":
      return Folder
    case "github":
      return GitBranch
    case "google-drive":
      return HardDrive
    default:
      return Globe
  }
}

function formatDate(dateString: string): string {
  return new Date(dateString).toLocaleString()
}

function formatTimeAgo(dateString: string): string {
  const date = new Date(dateString)
  const now = new Date()
  const diffInHours = Math.floor((now.getTime() - date.getTime()) / (1000 * 60 * 60))

  if (diffInHours < 1) return "just now"
  if (diffInHours < 24) return `${diffInHours}h ago`
  const diffInDays = Math.floor(diffInHours / 24)
  if (diffInDays < 7) return `${diffInDays}d ago`
  return `${Math.floor(diffInDays / 7)}w ago`
}

function formatDuration(seconds: number): string {
  if (seconds < 60) return `${seconds}s`
  const minutes = Math.floor(seconds / 60)
  const remainingSeconds = seconds % 60
  return `${minutes}m ${remainingSeconds}s`
}

function formatSchedule(cronExpression: string) {
  const parts = cronExpression.split(" ")
  if (parts.length !== 5) return cronExpression

  const [minute, hour, day, month, dayOfWeek] = parts

  if (day === "*" && month === "*" && dayOfWeek === "*") {
    return `Daily at ${hour}:${minute.padStart(2, "0")}`
  } else if (day === "*" && month === "*" && dayOfWeek !== "*") {
    const days = ["Sunday", "Monday", "Tuesday", "Wednesday", "Thursday", "Friday", "Saturday"]
    return `Weekly on ${days[parseInt(dayOfWeek)]} at ${hour}:${minute.padStart(2, "0")}`
  }

  return cronExpression
}

// Action handlers
const syncNow = async () => {
  await triggerManualSync()
}

const pauseSync = async () => {
  console.log("Pausing sync for:", dataSource.value.name)
  dataSource.value.status = "paused"
}

const resumeSync = async () => {
  console.log("Resuming sync for:", dataSource.value.name)
  dataSource.value.status = "active"
}

const openSourceUrl = () => {
  window.open(dataSource.value.url, "_blank")
}

const downloadLogs = () => {
  console.log("Downloading sync logs for:", dataSource.value.name)
  // TODO: Implement log download functionality
}

// Real-time updates simulation
const lastSyncTime = computed(() => formatTimeAgo(dataSource.value.lastSync))
const nextSyncTime = computed(() => formatTimeAgo(dataSource.value.nextSync))

// Auto-refresh logic
const autoRefresh = ref(true)
const refreshInterval = ref<NodeJS.Timeout | null>(null)

onMounted(() => {
  if (autoRefresh.value) {
    refreshInterval.value = setInterval(() => {
      // In real implementation, this would fetch fresh data from API
      console.log("Auto-refreshing data source metrics...")
    }, 30000) // Refresh every 30 seconds
  }
})

onUnmounted(() => {
  if (refreshInterval.value) {
    clearInterval(refreshInterval.value)
  }
})

// Sync in progress simulation
const syncInProgress = ref(false)

const triggerManualSync = async () => {
  syncInProgress.value = true
  console.log("Starting manual sync for:", dataSource.value.name)

  try {
    // Simulate API call
    await new Promise((resolve) => setTimeout(resolve, 2000))

    // Update last sync time
    dataSource.value.lastSync = new Date().toISOString()

    // Add a new sync log entry
    const newLog: SyncLog = {
      id: `log_${Date.now()}`,
      timestamp: new Date().toISOString(),
      status: "success",
      message: "Manual sync completed successfully",
      documentsProcessed: dataSource.value.documentsCount,
      documentsAdded: Math.floor(Math.random() * 10),
      documentsUpdated: Math.floor(Math.random() * 5),
      documentsRemoved: Math.floor(Math.random() * 2),
      duration: Math.floor(Math.random() * 200) + 100
    }

    syncLogs.value.unshift(newLog)

    // Show success toast
    // toast.success("Sync completed successfully")
  } catch (error) {
    console.error("Sync failed:", error)
    // toast.error("Sync failed. Please try again.")
  } finally {
    syncInProgress.value = false
  }
}
</script>

<template>
  <div class="flex flex-col gap-6">
    <!-- Data Source Header -->
    <div class="flex items-start justify-between">
      <div class="flex items-center gap-4">
        <div class="p-3 rounded-lg bg-primary/10">
          <component :is="getTypeIcon(dataSource.type)" class="size-8 text-primary" />
        </div>
        <div>
          <div class="flex items-center gap-3 mb-2">
            <h1 class="text-2xl font-bold">{{ dataSource.name }}</h1>
            <div
              class="inline-flex items-center gap-1 px-2 py-1 rounded-md text-xs font-medium border"
              :class="getStatusColor(dataSource.status)"
            >
              <div
                class="size-1.5 rounded-full"
                :class="{
                  'bg-emerald-500 animate-pulse': dataSource.status === 'active' && syncInProgress,
                  'bg-emerald-500': dataSource.status === 'active' && !syncInProgress,
                  'bg-amber-500': dataSource.status === 'paused',
                  'bg-red-500': dataSource.status === 'error',
                  'bg-gray-400': dataSource.status === 'inactive'
                }"
              />
              {{ syncInProgress ? "Syncing" : dataSource.status.charAt(0).toUpperCase() + dataSource.status.slice(1) }}
            </div>
          </div>
          <p class="text-muted-foreground mb-2">{{ dataSource.description }}</p>
          <div class="flex items-center gap-6 text-sm text-muted-foreground">
            <div class="flex items-center gap-1">
              <component :is="getSourceIcon(dataSource.source)" class="size-4" />
              <span>{{ dataSource.source.charAt(0).toUpperCase() + dataSource.source.slice(1) }}</span>
            </div>
            <div class="flex items-center gap-1">
              <User class="size-4" />
              <span>{{ dataSource.owner }}</span>
            </div>
            <div class="flex items-center gap-1">
              <Calendar class="size-4" />
              <span>{{ formatSchedule(dataSource.schedule) }}</span>
            </div>
          </div>
          <div class="flex items-center gap-2 mt-2">
            <span class="text-xs text-muted-foreground">Tags:</span>
            <div class="flex flex-wrap gap-1">
              <UiBadge v-for="tag in dataSource.tags" :key="tag" variant="secondary" class="text-xs">
                {{ tag }}
              </UiBadge>
            </div>
          </div>
        </div>
      </div>

      <div class="flex items-center gap-2">
        <UiButton variant="outline" size="sm" @click="openSourceUrl">
          <ExternalLink class="size-4 mr-2" />
          Open Source
        </UiButton>
        <UiButton variant="outline" size="sm" @click="downloadLogs">
          <Download class="size-4 mr-2" />
          Export Logs
        </UiButton>
        <UiButton
          variant="outline"
          size="sm"
          @click="syncNow"
          :disabled="dataSource.status === 'error' || syncInProgress"
          :loading="syncInProgress"
        >
          <RefreshCw class="size-4 mr-2" :class="{ 'animate-spin': syncInProgress }" />
          {{ syncInProgress ? "Syncing..." : "Sync Now" }}
        </UiButton>
        <UiDropdownMenu>
          <UiDropdownMenuTrigger as-child>
            <UiButton variant="outline" size="sm">
              <MoreHorizontal class="size-4" />
            </UiButton>
          </UiDropdownMenuTrigger>
          <UiDropdownMenuContent align="end" class="w-56">
            <UiDropdownMenuItem @click="showEditModal = true">
              <Edit class="mr-2 size-4" />
              Edit Settings
            </UiDropdownMenuItem>
            <UiDropdownMenuItem @click="showScheduleModal = true">
              <Calendar class="mr-2 size-4" />
              Update Schedule
            </UiDropdownMenuItem>
            <UiDropdownMenuItem @click="dataSource.status === 'active' ? pauseSync() : resumeSync()">
              <component :is="dataSource.status === 'active' ? Pause : Play" class="mr-2 size-4" />
              {{ dataSource.status === "active" ? "Pause" : "Resume" }} Sync
            </UiDropdownMenuItem>
            <UiDropdownMenuSeparator />
            <UiDropdownMenuItem>
              <Copy class="mr-2 size-4" />
              Clone Data Source
            </UiDropdownMenuItem>
            <UiDropdownMenuItem class="text-destructive" @click="showDeleteModal = true">
              <Trash2 class="mr-2 size-4" />
              Delete Data Source
            </UiDropdownMenuItem>
          </UiDropdownMenuContent>
        </UiDropdownMenu>
      </div>
    </div>

    <!-- Stats Overview -->
    <div class="grid gap-4 md:grid-cols-2 lg:grid-cols-4">
      <CardsStats
        v-for="stat in dataSourceStats"
        :key="stat.title"
        :title="stat.title"
        :value="stat.value"
        :icon="stat.icon"
        :description="stat.description"
        :variant="stat.variant"
      />
    </div>

    <!-- Sync Status Card - Full Width -->
    <UiCard>
      <UiCardHeader>
        <UiCardTitle class="flex items-center gap-2">
          <RefreshCw class="size-5" />
          Sync Status
        </UiCardTitle>
        <UiCardDescription>Current sync information and next scheduled run</UiCardDescription>
      </UiCardHeader>
      <UiCardContent>
        <div class="grid md:grid-cols-2 gap-6">
          <div class="space-y-4">
            <div class="flex items-center justify-between p-4 border rounded-lg">
              <div>
                <div class="text-sm font-medium">Last Sync</div>
                <div class="text-2xl font-bold text-muted-foreground">{{ lastSyncTime }}</div>
                <div class="text-xs text-muted-foreground">{{ formatDate(dataSource.lastSync) }}</div>
              </div>
              <div class="p-2 rounded-lg bg-emerald-50 dark:bg-emerald-950/30">
                <CheckCircle class="size-6 text-emerald-600 dark:text-emerald-400" />
              </div>
            </div>

            <div class="flex items-center justify-between p-4 border rounded-lg">
              <div>
                <div class="text-sm font-medium">Next Sync</div>
                <div class="text-2xl font-bold text-muted-foreground">{{ nextSyncTime }}</div>
                <div class="text-xs text-muted-foreground">{{ formatDate(dataSource.nextSync) }}</div>
              </div>
              <div class="p-2 rounded-lg bg-blue-50 dark:bg-blue-950/30">
                <Clock class="size-6 text-blue-600 dark:text-blue-400" />
              </div>
            </div>
          </div>

          <!-- Performance Chart Area -->
          <ChartPlaceholder
            :icon="BarChart3"
            title="Sync Performance Chart"
            description="Performance analytics coming soon"
          />
        </div>
      </UiCardContent>
    </UiCard>

    <!-- Main Content Grid - 1x2 Layout -->
    <div class="grid gap-6 lg:grid-cols-2">
      <!-- Recent Sync Logs (Left) -->
      <CardsDataList title="Recent Sync Logs" :icon="Activity">
        <template #actions>
          <UiButton size="sm" variant="outline" @click="downloadLogs">
            <Download class="size-4" />
          </UiButton>
        </template>

        <UiScrollArea class="h-[440px]">
          <div class="space-y-3 pr-4">
            <div
              v-for="log in syncLogs"
              :key="log.id"
              class="flex items-start gap-3 p-4 border rounded-lg hover:bg-muted/50 transition-colors"
            >
              <div
                class="p-2 rounded-md mt-1"
                :class="{
                  'bg-emerald-50 dark:bg-emerald-950/30': log.status === 'success',
                  'bg-amber-50 dark:bg-amber-950/30': log.status === 'warning',
                  'bg-red-50 dark:bg-red-950/30': log.status === 'error',
                  'bg-blue-50 dark:bg-blue-950/30': log.status === 'info'
                }"
              >
                <component
                  :is="getSyncStatusIcon(log.status)"
                  class="size-4"
                  :class="{
                    'text-emerald-600 dark:text-emerald-400': log.status === 'success',
                    'text-amber-600 dark:text-amber-400': log.status === 'warning',
                    'text-red-600 dark:text-red-400': log.status === 'error',
                    'text-blue-600 dark:text-blue-400': log.status === 'info'
                  }"
                />
              </div>

              <div class="flex-1 min-w-0">
                <div class="flex items-center gap-2 mb-1">
                  <div class="font-medium text-sm">{{ log.message }}</div>
                  <div
                    class="inline-flex items-center gap-1 px-2 py-0.5 rounded text-xs font-medium border"
                    :class="getSyncStatusColor(log.status)"
                  >
                    {{ log.status }}
                  </div>
                </div>

                <div class="text-xs text-muted-foreground mb-2">
                  {{ formatDate(log.timestamp) }} • {{ formatDuration(log.duration || 0) }}
                </div>

                <div v-if="log.documentsProcessed" class="flex items-center gap-4 text-xs text-muted-foreground">
                  <span>{{ log.documentsProcessed }} processed</span>
                  <span v-if="log.documentsAdded">+{{ log.documentsAdded }} added</span>
                  <span v-if="log.documentsUpdated">~{{ log.documentsUpdated }} updated</span>
                  <span v-if="log.documentsRemoved">-{{ log.documentsRemoved }} removed</span>
                </div>

                <div
                  v-if="log.errorDetails"
                  class="mt-2 p-2 bg-red-50 dark:bg-red-950/30 border border-red-200 dark:border-red-800 rounded text-xs text-red-700 dark:text-red-400"
                >
                  {{ log.errorDetails }}
                </div>
              </div>
            </div>
          </div>
        </UiScrollArea>
      </CardsDataList>

      <!-- Configuration Overview (Right) -->
      <CardsDataList title="Configuration" :icon="Settings">
        <template #actions>
          <UiButton size="sm" variant="outline" @click="showConfigModal = true">
            <Edit class="size-4" />
          </UiButton>
        </template>

        <div class="flex flex-col">
          <div class="space-y-3 flex-1">
            <div class="flex items-center justify-between py-2 border-b">
              <span class="text-sm font-medium">Source URL</span>
              <div class="flex items-center gap-2">
                <span class="text-sm text-muted-foreground truncate max-w-48">{{ dataSource.url }}</span>
                <UiButton size="sm" variant="ghost" @click="openSourceUrl">
                  <ExternalLink class="size-3" />
                </UiButton>
              </div>
            </div>

            <div class="flex items-center justify-between py-2 border-b">
              <span class="text-sm font-medium">Sync Schedule</span>
              <span class="text-sm text-muted-foreground">{{ formatSchedule(dataSource.schedule) }}</span>
            </div>

            <div class="flex items-center justify-between py-2 border-b">
              <span class="text-sm font-medium">Document Type</span>
              <UiBadge variant="secondary">{{ dataSource.type.replace("-", " ") }}</UiBadge>
            </div>

            <div class="flex items-center justify-between py-2 border-b">
              <span class="text-sm font-medium">Total Storage</span>
              <span class="text-sm text-muted-foreground"
                >{{ (documentMetrics.totalDocuments * documentMetrics.averageDocumentSize).toFixed(1) }}MB</span
              >
            </div>

            <div class="flex items-center justify-between py-2 border-b">
              <span class="text-sm font-medium">Owner</span>
              <span class="text-sm text-muted-foreground">{{ dataSource.owner }}</span>
            </div>

            <div class="flex items-center justify-between py-2 border-b">
              <span class="text-sm font-medium">Last Sync</span>
              <span class="text-sm text-muted-foreground">{{ lastSyncTime }}</span>
            </div>

            <div class="flex items-center justify-between py-2 border-b">
              <span class="text-sm font-medium">Next Sync</span>
              <span class="text-sm text-muted-foreground">{{ nextSyncTime }}</span>
            </div>

            <div class="flex items-center justify-between py-2">
              <span class="text-sm font-medium">Documents</span>
              <span class="text-sm text-muted-foreground">{{ dataSource.documentsCount.toLocaleString() }}</span>
            </div>
          </div>

          <div class="pt-4 border-t mt-auto">
            <div class="grid grid-cols-2 gap-3">
              <UiButton size="sm" variant="outline" @click="showScheduleModal = true">
                <Calendar class="size-4 mr-2" />
                Schedule
              </UiButton>
              <UiButton size="sm" variant="outline" @click="showConfigModal = true">
                <Settings class="size-4 mr-2" />
                Settings
              </UiButton>
            </div>
          </div>
        </div>
      </CardsDataList>
    </div>

    <!-- Modals would go here -->
    <!-- TODO: Add edit, delete, schedule, and configuration modals -->
  </div>
</template>

