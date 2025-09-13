<script lang="ts">
const dataSources = ref<DataSourceData[]>([
  {
    id: "ds_1",
    name: "Customer Documentation",
    description: "Product documentation and user guides for customer support",
    type: "documentation",
    source: "confluence",
    url: "https://company.atlassian.net/wiki",
    status: "active",
    schedule: "0 2 * * *", // Daily at 2 AM
    lastSync: "2025-01-15T02:00:00Z",
    nextSync: "2025-01-16T02:00:00Z",
    documentsCount: 1250,
    owner: "Alice Johnson",
    tags: ["documentation", "support", "customer"]
  },
  {
    id: "ds_2",
    name: "Product Knowledge Base",
    description: "Internal product specifications and technical documentation",
    type: "knowledge-base",
    source: "notion",
    url: "https://notion.so/product-docs",
    status: "active",
    schedule: "0 3 * * 1", // Weekly on Monday at 3 AM
    lastSync: "2025-01-13T03:00:00Z",
    nextSync: "2025-01-20T03:00:00Z",
    documentsCount: 850,
    owner: "Bob Smith",
    tags: ["product", "technical", "internal"]
  },
  {
    id: "ds_3",
    name: "Company Policies",
    description: "HR policies, procedures, and compliance documentation",
    type: "policies",
    source: "sharepoint",
    url: "https://company.sharepoint.com/policies",
    status: "paused",
    schedule: "0 4 * * 0", // Weekly on Sunday at 4 AM
    lastSync: "2025-01-10T04:00:00Z",
    nextSync: "2025-01-19T04:00:00Z",
    documentsCount: 320,
    owner: "Carol Williams",
    tags: ["hr", "policies", "compliance"]
  },
  {
    id: "ds_4",
    name: "API Documentation",
    description: "REST API documentation and integration guides",
    type: "api-docs",
    source: "github",
    url: "https://github.com/company/api-docs",
    status: "error",
    schedule: "0 1 * * *", // Daily at 1 AM
    lastSync: "2025-01-14T01:00:00Z",
    nextSync: "2025-01-16T01:00:00Z",
    documentsCount: 150,
    owner: "David Brown",
    tags: ["api", "documentation", "developer"]
  }
])
</script>

<script setup lang="ts">
import {
  Database,
  MoreVertical,
  Edit,
  Trash2,
  CheckCircle,
  XCircle,
  Clock,
  AlertTriangle,
  FileText,
  Play,
  Pause,
  RefreshCw,
  Eye,
  Settings,
  Plus,
  Calendar,
  Tag,
  Globe,
  GitBranch,
  Bookmark,
  Activity
} from "lucide-vue-next"
import type { FilterConfig, SearchConfig, DisplayConfig } from "@/components/SearchFilter.vue"
import type { StatsCardProps } from "@/components/Cards/Stats.vue"
import type { DataSourceData } from "@/models/datasource"

useSeoMeta({ title: "Data Sources - LLM Gateway" })

// Modal states
const showCreateModal = ref(false)
const showEditModal = ref(false)
const showDeleteModal = ref(false)
const showScheduleModal = ref(false)
const editingDataSource = ref<DataSourceData | null>(null)
const deletingDataSource = ref<DataSourceData | null>(null)
const schedulingDataSource = ref<DataSourceData | null>(null)

// Filter state
const activeFilters = ref<Record<string, string>>({})

// Stats configuration
const stats: StatsCardProps[] = [
  {
    title: "Total Data Sources",
    value: dataSources.value.length.toString(),
    icon: Database,
    description: "All registered sources",
    variant: "default"
  },
  {
    title: "Active Sources",
    value: dataSources.value.filter((ds) => ds.status === "active").length.toString(),
    icon: CheckCircle,
    description: "Currently syncing",
    variant: "chart-2"
  },
  {
    title: "Total Documents",
    value: dataSources.value.reduce((sum, ds) => sum + ds.documentsCount, 0).toString(),
    icon: FileText,
    description: "Documents indexed",
    variant: "chart-3"
  },
  {
    title: "Sync Issues",
    value: dataSources.value.filter((ds) => ds.status === "error").length.toString(),
    icon: AlertTriangle,
    description: "Requiring attention",
    variant: "chart-1"
  }
]

// Get unique values for filtering
const allTypes = [...new Set(dataSources.value.map((ds) => ds.type))]
const allSources = [...new Set(dataSources.value.map((ds) => ds.source))]
const allTags = [...new Set(dataSources.value.flatMap((ds) => ds.tags))]

const filterConfigs: FilterConfig[] = [
  {
    key: "type",
    label: "Type",
    options: allTypes.map((type) => ({ value: type, label: type.charAt(0).toUpperCase() + type.slice(1), icon: FileText }))
  },
  {
    key: "status",
    label: "Status",
    options: [
      { value: "active", label: "Active", icon: CheckCircle },
      { value: "paused", label: "Paused", icon: Pause },
      { value: "error", label: "Error", icon: AlertTriangle },
      { value: "inactive", label: "Inactive", icon: XCircle }
    ]
  },
  {
    key: "source",
    label: "Source",
    options: allSources.map((source) => ({ value: source, label: source.charAt(0).toUpperCase() + source.slice(1), icon: Globe }))
  },
  {
    key: "tag",
    label: "Tags",
    options: allTags.map((tag) => ({ value: tag, label: tag, icon: Tag }))
  }
]

const searchConfig: SearchConfig<DataSourceData> = {
  fields: ["name", "description", "type", "source", "owner"],
  placeholder: "Search data sources, filter by type, status..."
}

const displayConfig: DisplayConfig<DataSourceData> = {
  getItemText: (ds) => `${ds.name} - ${ds.description}`,
  getItemValue: (ds) => ds.name,
  getItemIcon: () => Database
}

// Event handlers
function handleFiltersChanged(filters: Record<string, string>) {
  activeFilters.value = filters
}

function handleItemSelected(dataSource: DataSourceData) {
  navigateTo(`/storage/${dataSource.id}`)
}

// Filtering logic
const filteredDataSources = computed(() => {
  let filtered = dataSources.value

  if (activeFilters.value.type && activeFilters.value.type !== "all") {
    filtered = filtered.filter((ds) => ds.type === activeFilters.value.type)
  }
  if (activeFilters.value.status && activeFilters.value.status !== "all") {
    filtered = filtered.filter((ds) => ds.status === activeFilters.value.status)
  }
  if (activeFilters.value.source && activeFilters.value.source !== "all") {
    filtered = filtered.filter((ds) => ds.source === activeFilters.value.source)
  }
  if (activeFilters.value.tag && activeFilters.value.tag !== "all") {
    filtered = filtered.filter((ds) => ds.tags.includes(activeFilters.value.tag!))
  }

  return filtered
})

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

function getStatusIcon(status: string) {
  switch (status) {
    case "active": return CheckCircle
    case "paused": return Pause
    case "error": return AlertTriangle
    case "inactive": return XCircle
    default: return XCircle
  }
}

function getTypeIcon(type: string) {
  switch (type) {
    case "documentation": return FileText
    case "knowledge-base": return Bookmark
    case "policies": return Settings
    case "api-docs": return GitBranch
    default: return FileText
  }
}

// Modal handlers
const openCreateModal = () => {
  showCreateModal.value = true
}

const openEditModal = (dataSource: DataSourceData) => {
  editingDataSource.value = dataSource
  showEditModal.value = true
}

const openDeleteModal = (dataSource: DataSourceData) => {
  deletingDataSource.value = dataSource
  showDeleteModal.value = true
}

const openScheduleModal = (dataSource: DataSourceData) => {
  schedulingDataSource.value = dataSource
  showScheduleModal.value = true
}

// Action handlers
const syncDataSource = async (dataSource: DataSourceData) => {
  console.log("Syncing data source:", dataSource.name)
  // TODO: Implement sync logic
}

const pauseDataSource = async (dataSource: DataSourceData) => {
  console.log("Pausing data source:", dataSource.name)
  dataSource.status = "paused"
}

const resumeDataSource = async (dataSource: DataSourceData) => {
  console.log("Resuming data source:", dataSource.name)
  dataSource.status = "active"
}

// Helper functions
const formatSchedule = (cronExpression: string) => {
  // Simple cron to human readable conversion
  const parts = cronExpression.split(' ')
  if (parts.length !== 5) return cronExpression
  
  const [minute, hour, day, month, dayOfWeek] = parts
  
  if (day === '*' && month === '*' && dayOfWeek === '*') {
    return `Daily at ${hour}:${minute.padStart(2, '0')}`
  } else if (day === '*' && month === '*' && dayOfWeek !== '*') {
    const days = ['Sunday', 'Monday', 'Tuesday', 'Wednesday', 'Thursday', 'Friday', 'Saturday']
    return `Weekly on ${days[parseInt(dayOfWeek)]} at ${hour}:${minute.padStart(2, '0')}`
  }
  
  return cronExpression
}

const formatDate = (dateString: string) => {
  return new Date(dateString).toLocaleString()
}

// Handle data source creation
const onDataSourceCreated = (newDataSource: DataSourceData) => {
  console.log("Data source created:", newDataSource)
  // Add new data source to the list
  dataSources.value.push(newDataSource)
  showCreateModal.value = false
}
</script>

<template>
  <div class="flex flex-col gap-6">
    <!-- Header -->
    <PageHeader title="Data Sources" subtext="Manage data sources for RAG ingestion with automated scheduling">
      <ButtonsCreate title="Add Data Source" :action="openCreateModal" />
    </PageHeader>

    <!-- Stats Cards -->
    <div class="grid gap-4 md:grid-cols-2 lg:grid-cols-4">
      <CardsStats
        v-for="stat in stats"
        :key="stat.title"
        :title="stat.title"
        :value="stat.value"
        :icon="stat.icon"
        :description="stat.description"
        :variant="stat.variant"
      />
    </div>

    <!-- Search & Filter Component -->
    <SearchFilter
      :items="dataSources"
      :filters="filterConfigs"
      :search-config="searchConfig"
      :display-config="displayConfig"
      @filters-changed="handleFiltersChanged"
      @item-selected="handleItemSelected"
    />

    <!-- Data Sources List -->
    <CardsDataList title="All Data Sources" :icon="Database">
      <template #actions>
        <ButtonsCreate title="Add Data Source" :action="openCreateModal" />
      </template>

      <div class="space-y-4">
        <div
          v-for="dataSource in filteredDataSources"
          :key="dataSource.id"
          class="flex items-center justify-between p-6 border rounded-lg hover:bg-muted/50 cursor-pointer transition-colors"
          @click="navigateTo(`/storage/${dataSource.id}`)"
        >
          <div class="flex items-center gap-6">
            <div class="p-3 rounded-lg bg-primary/10">
              <component :is="getTypeIcon(dataSource.type)" class="size-6 text-primary" />
            </div>

            <div class="space-y-2">
              <div class="flex items-center gap-3">
                <h3 class="font-semibold text-lg">{{ dataSource.name }}</h3>
                <div
                  class="inline-flex items-center gap-1 px-2 py-1 rounded-md text-xs font-medium border"
                  :class="getStatusColor(dataSource.status)"
                >
                  <component :is="getStatusIcon(dataSource.status)" class="size-3" />
                  {{ dataSource.status.charAt(0).toUpperCase() + dataSource.status.slice(1) }}
                </div>
              </div>

              <p class="text-muted-foreground">{{ dataSource.description }}</p>

              <div class="flex items-center gap-6 text-sm text-muted-foreground">
                <div class="flex items-center gap-1">
                  <FileText class="size-4" />
                  <span>{{ formatNumberPretty(dataSource.documentsCount) }} docs</span>
                </div>
                <div class="flex items-center gap-1">
                  <Calendar class="size-4" />
                  <span>{{ formatSchedule(dataSource.schedule) }}</span>
                </div>
                <div class="flex items-center gap-1">
                  <Globe class="size-4" />
                  <span>{{ dataSource.source }}</span>
                </div>
              </div>

              <!-- Tags -->
              <div class="flex items-center gap-2">
                <span class="text-xs text-muted-foreground">Tags:</span>
                <div class="flex flex-wrap gap-1">
                  <UiBadge v-for="tag in dataSource.tags" :key="tag" variant="secondary" class="text-xs">
                    {{ tag }}
                  </UiBadge>
                </div>
              </div>

              <!-- Sync Info -->
              <div class="flex items-center gap-4 text-xs text-muted-foreground">
                <span>Owner: {{ dataSource.owner }}</span>
                <span>Last sync: {{ formatDate(dataSource.lastSync) }}</span>
                <span v-if="dataSource.status === 'active'">Next sync: {{ formatDate(dataSource.nextSync) }}</span>
              </div>
            </div>
          </div>

          <div class="flex items-center gap-2" @click.stop>
            <!-- Quick Actions -->
            <UiTooltip>
              <UiTooltipTrigger as-child>
                <UiButton 
                  variant="ghost" 
                  size="sm"
                  @click="syncDataSource(dataSource)"
                  :disabled="dataSource.status === 'error'"
                >
                  <RefreshCw class="size-4" />
                </UiButton>
              </UiTooltipTrigger>
              <UiTooltipContent>Sync Now</UiTooltipContent>
            </UiTooltip>

            <UiTooltip>
              <UiTooltipTrigger as-child>
                <UiButton 
                  variant="ghost" 
                  size="sm"
                  @click="dataSource.status === 'active' ? pauseDataSource(dataSource) : resumeDataSource(dataSource)"
                >
                  <component :is="dataSource.status === 'active' ? Pause : Play" class="size-4" />
                </UiButton>
              </UiTooltipTrigger>
              <UiTooltipContent>{{ dataSource.status === 'active' ? 'Pause' : 'Resume' }}</UiTooltipContent>
            </UiTooltip>

            <!-- More Actions Dropdown -->
            <UiDropdownMenu>
              <UiDropdownMenuTrigger as-child>
                <UiButton variant="ghost" size="sm">
                  <MoreVertical class="size-4" />
                </UiButton>
              </UiDropdownMenuTrigger>
              <UiDropdownMenuContent align="end" class="w-56">
                <UiDropdownMenuItem @click="navigateTo(`/storage/${dataSource.id}`)">
                  <Eye class="mr-2 size-4" />
                  View Details
                </UiDropdownMenuItem>
                <UiDropdownMenuItem @click="openScheduleModal(dataSource)">
                  <Calendar class="mr-2 size-4" />
                  Edit Schedule
                </UiDropdownMenuItem>
                <UiDropdownMenuItem @click="syncDataSource(dataSource)">
                  <RefreshCw class="mr-2 size-4" />
                  Sync Now
                </UiDropdownMenuItem>
                <UiDropdownMenuSeparator />
                <UiDropdownMenuItem @click="openEditModal(dataSource)">
                  <Edit class="mr-2 size-4" />
                  Edit Source
                </UiDropdownMenuItem>
                <UiDropdownMenuItem>
                  <Activity class="mr-2 size-4" />
                  View Logs
                </UiDropdownMenuItem>
                <UiDropdownMenuSeparator />
                <UiDropdownMenuItem class="text-destructive" @click="openDeleteModal(dataSource)">
                  <Trash2 class="mr-2 size-4" />
                  Delete Source
                </UiDropdownMenuItem>
              </UiDropdownMenuContent>
            </UiDropdownMenu>
          </div>
        </div>
      </div>
    </CardsDataList>

    <!-- Modals -->
    <!-- Create Data Source Modal -->
    <LazyModalsStorageCreate 
      v-model:open="showCreateModal" 
      @created="onDataSourceCreated" 
    />

    <!-- Delete Confirmation Modal -->
    <ConfirmationModal
      v-model:open="showDeleteModal"
      title="Delete Data Source"
      :description="`Are you sure you want to delete ${deletingDataSource?.name}? This action cannot be undone and will remove all indexed documents.`"
      confirm-text="Delete Data Source"
      variant="destructive"
      @confirm="() => {}"
      @cancel="
        () => {
          showDeleteModal = false
          deletingDataSource = null
        }
      "
    />
  </div>
</template>