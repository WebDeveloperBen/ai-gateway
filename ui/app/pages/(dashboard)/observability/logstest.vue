<script lang="ts">
interface LogEntry {
  id: string
  timestamp: string
  level: "error" | "warn" | "info" | "debug"
  message: string
  application: string
  method: string
  endpoint: string
  statusCode: number
  responseTime: number
  userId?: string
  tags: string[]
  requestId: string
  source: "gateway" | "application" | "system"
}

const logData: LogEntry[] = [
  {
    id: "log_001",
    timestamp: "2025-01-15T14:30:25.123Z",
    level: "error",
    message: "Rate limit exceeded for application Customer Service Bot",
    application: "Customer Service Bot",
    method: "POST",
    endpoint: "/v1/chat/completions",
    statusCode: 429,
    responseTime: 45,
    userId: "user_123",
    tags: ["rate-limit", "production", "customer-success"],
    requestId: "req_abc123",
    source: "gateway"
  },
  {
    id: "log_002",
    timestamp: "2025-01-15T14:29:58.456Z",
    level: "info",
    message: "Successful API request processed",
    application: "Content Generator",
    method: "POST",
    endpoint: "/v1/chat/completions",
    statusCode: 200,
    responseTime: 1250,
    userId: "user_456",
    tags: ["success", "marketing", "content"],
    requestId: "req_def456",
    source: "gateway"
  },
  {
    id: "log_003",
    timestamp: "2025-01-15T14:29:45.789Z",
    level: "warn",
    message: "Content filter triggered - potential inappropriate content detected",
    application: "Analytics Dashboard",
    method: "POST",
    endpoint: "/v1/chat/completions",
    statusCode: 200,
    responseTime: 890,
    userId: "user_789",
    tags: ["content-filter", "moderation", "analytics"],
    requestId: "req_ghi789",
    source: "system"
  },
  {
    id: "log_004",
    timestamp: "2025-01-15T14:29:30.012Z",
    level: "info",
    message: "New API key generated for application",
    application: "Customer Service Bot",
    method: "POST",
    endpoint: "/api/keys",
    statusCode: 201,
    responseTime: 124,
    userId: "admin_001",
    tags: ["api-key", "admin", "security"],
    requestId: "req_jkl012",
    source: "application"
  },
  {
    id: "log_005",
    timestamp: "2025-01-15T14:29:15.345Z",
    level: "error",
    message: "Authentication failed - invalid API key",
    application: "Unknown",
    method: "GET",
    endpoint: "/v1/models",
    statusCode: 401,
    responseTime: 23,
    tags: ["auth-failure", "security", "invalid-key"],
    requestId: "req_mno345",
    source: "gateway"
  },
  {
    id: "log_006",
    timestamp: "2025-01-15T14:29:00.678Z",
    level: "debug",
    message: "Request processed through load balancer",
    application: "Content Generator",
    method: "POST",
    endpoint: "/v1/embeddings",
    statusCode: 200,
    responseTime: 456,
    userId: "user_999",
    tags: ["load-balancer", "embeddings", "debug"],
    requestId: "req_pqr678",
    source: "system"
  },
  {
    id: "log_007",
    timestamp: "2025-01-15T14:28:45.901Z",
    level: "warn",
    message: "High response time detected - performance degradation",
    application: "Analytics Dashboard",
    method: "POST",
    endpoint: "/v1/chat/completions",
    statusCode: 200,
    responseTime: 5670,
    userId: "user_111",
    tags: ["performance", "slow-response", "monitoring"],
    requestId: "req_stu901",
    source: "system"
  },
  {
    id: "log_008",
    timestamp: "2025-01-15T14:28:30.234Z",
    level: "info",
    message: "Policy applied successfully to application request",
    application: "Customer Service Bot",
    method: "POST",
    endpoint: "/v1/chat/completions",
    statusCode: 200,
    responseTime: 789,
    userId: "user_222",
    tags: ["policy", "governance", "success"],
    requestId: "req_vwx234",
    source: "gateway"
  },
  {
    id: "log_009",
    timestamp: "2025-01-15T14:28:15.567Z",
    level: "error",
    message: "Downstream service unavailable - OpenAI API timeout",
    application: "Content Generator",
    method: "POST",
    endpoint: "/v1/chat/completions",
    statusCode: 504,
    responseTime: 30000,
    userId: "user_333",
    tags: ["timeout", "openai", "service-unavailable"],
    requestId: "req_yz1567",
    source: "gateway"
  },
  {
    id: "log_010",
    timestamp: "2025-01-15T14:28:00.890Z",
    level: "info",
    message: "User session started",
    application: "Analytics Dashboard",
    method: "POST",
    endpoint: "/auth/login",
    statusCode: 200,
    responseTime: 145,
    userId: "user_444",
    tags: ["auth", "session", "login"],
    requestId: "req_abc890",
    source: "application"
  }
]
</script>

<script setup lang="ts">
import { Download, Filter, Code, History, Play, ChevronDown, ChevronUp, ChevronRight } from "lucide-vue-next"
import type { TableColumn } from "~/components/Blocks/DataTable.vue"

const appConfig = useAppConfig()

const getLevelColor = (level: LogEntry["level"]) => {
  switch (level) {
    case "error":
      return "text-red-600 bg-red-50 border-red-200"
    case "warn":
      return "text-yellow-600 bg-yellow-50 border-yellow-200"
    case "info":
      return "text-blue-600 bg-blue-50 border-blue-200"
    case "debug":
      return "text-gray-600 bg-gray-50 border-gray-200"
    default:
      return "text-gray-600 bg-gray-50 border-gray-200"
  }
}

const getStatusColor = (statusCode: number) => {
  if (statusCode >= 200 && statusCode < 300) {
    return "text-green-600 bg-green-50 border-green-200"
  } else if (statusCode >= 300 && statusCode < 400) {
    return "text-blue-600 bg-blue-50 border-blue-200"
  } else if (statusCode >= 400 && statusCode < 500) {
    return "text-yellow-600 bg-yellow-50 border-yellow-200"
  } else if (statusCode >= 500) {
    return "text-red-600 bg-red-50 border-red-200"
  }
  return "text-gray-600 bg-gray-50 border-gray-200"
}

const getMethodColor = (method: string) => {
  const methodColors = {
    GET: "text-blue-600 bg-blue-50 border-blue-200",
    POST: "text-green-600 bg-green-50 border-green-200",
    PUT: "text-yellow-600 bg-yellow-50 border-yellow-200",
    DELETE: "text-red-600 bg-red-50 border-red-200"
  }
  return methodColors[method as keyof typeof methodColors] || "text-gray-600 bg-gray-50 border-gray-200"
}

const searchQuery = ref("")
const advancedQuery = ref("")
const isExporting = ref(false)
const showAdvancedQuery = ref(false)
const queryHistory = ref<string[]>([])

// KQL-like query parsing (same as original)
const parseAdvancedQuery = (query: string) => {
  if (!query.trim()) return logData

  let filteredData = [...logData]

  const parts = query.split("|").map((p) => p.trim())

  for (const part of parts) {
    if (
      part.startsWith("where ") ||
      part.includes("==") ||
      part.includes("!=") ||
      part.includes(">") ||
      part.includes("<")
    ) {
      filteredData = applyWhereClause(filteredData, part)
    } else if (part.startsWith("order by ")) {
      filteredData = applyOrderBy(filteredData, part)
    } else if (part.startsWith("limit ") || part.startsWith("take ")) {
      filteredData = applyLimit(filteredData, part)
    } else if (part.includes("contains")) {
      filteredData = applyContains(filteredData, part)
    }
  }

  return filteredData
}

const applyWhereClause = (data: LogEntry[], clause: string) => {
  const cleanClause = clause.replace("where ", "").trim()

  if (cleanClause.includes('level == "error"')) {
    return data.filter((log) => log.level === "error")
  } else if (cleanClause.includes('level == "warn"')) {
    return data.filter((log) => log.level === "warn")
  } else if (cleanClause.includes('level == "info"')) {
    return data.filter((log) => log.level === "info")
  } else if (cleanClause.includes("statusCode == 401")) {
    return data.filter((log) => log.statusCode === 401)
  } else if (cleanClause.includes("statusCode == 403")) {
    return data.filter((log) => log.statusCode === 403)
  } else if (cleanClause.includes("statusCode >= 400")) {
    return data.filter((log) => log.statusCode >= 400)
  } else if (cleanClause.includes("statusCode >= 500")) {
    return data.filter((log) => log.statusCode >= 500)
  } else if (cleanClause.includes("responseTime > 2000")) {
    return data.filter((log) => log.responseTime > 2000)
  } else if (cleanClause.includes("responseTime > 1000")) {
    return data.filter((log) => log.responseTime > 1000)
  } else if (cleanClause.includes("ago(1h)")) {
    const oneHourAgo = new Date(Date.now() - 60 * 60 * 1000)
    return data.filter((log) => new Date(log.timestamp) > oneHourAgo)
  } else if (cleanClause.includes("ago(1d)")) {
    const oneDayAgo = new Date(Date.now() - 24 * 60 * 60 * 1000)
    return data.filter((log) => new Date(log.timestamp) > oneDayAgo)
  }

  return data
}

const applyContains = (data: LogEntry[], clause: string) => {
  if (clause.includes('tags contains "rate-limit"')) {
    return data.filter((log) => log.tags.includes("rate-limit"))
  } else if (clause.includes('tags contains "error"')) {
    return data.filter((log) => log.tags.some((tag) => tag.includes("error")))
  } else if (clause.includes("message contains")) {
    const searchTerm = clause.match(/"([^"]+)"/)?.[1]
    if (searchTerm) {
      return data.filter((log) => log.message.toLowerCase().includes(searchTerm.toLowerCase()))
    }
  }

  return data
}

const applyOrderBy = (data: LogEntry[], clause: string) => {
  const orderClause = clause.replace("order by ", "").trim()

  if (orderClause.includes("timestamp desc")) {
    return data.sort((a, b) => new Date(b.timestamp).getTime() - new Date(a.timestamp).getTime())
  } else if (orderClause.includes("timestamp asc")) {
    return data.sort((a, b) => new Date(a.timestamp).getTime() - new Date(b.timestamp).getTime())
  } else if (orderClause.includes("responseTime desc")) {
    return data.sort((a, b) => b.responseTime - a.responseTime)
  } else if (orderClause.includes("responseTime asc")) {
    return data.sort((a, b) => a.responseTime - b.responseTime)
  }

  return data
}

const applyLimit = (data: LogEntry[], clause: string) => {
  const limitMatch = clause.match(/(?:limit|take)\s+(\d+)/)
  if (limitMatch && limitMatch.length > 1) {
    const limit = parseInt(limitMatch[1]!)
    return data.slice(0, limit)
  }
  return data
}

const filteredLogData = computed(() => {
  if (showAdvancedQuery.value && advancedQuery.value.trim()) {
    return parseAdvancedQuery(advancedQuery.value)
  }

  // Simple search functionality
  if (searchQuery.value.trim()) {
    return logData.filter(
      (log) =>
        log.message.toLowerCase().includes(searchQuery.value.toLowerCase()) ||
        log.application.toLowerCase().includes(searchQuery.value.toLowerCase()) ||
        log.requestId.toLowerCase().includes(searchQuery.value.toLowerCase()) ||
        log.tags.some((tag) => tag.toLowerCase().includes(searchQuery.value.toLowerCase()))
    )
  }

  return logData
})

const loadHistoryQuery = (query: string) => {
  advancedQuery.value = query
}

const handleRunQuery = () => {
  if (advancedQuery.value.trim()) {
    if (!queryHistory.value.includes(advancedQuery.value)) {
      queryHistory.value.unshift(advancedQuery.value)
      if (queryHistory.value.length > 10) {
        queryHistory.value = queryHistory.value.slice(0, 10)
      }
    }
  }
}

const handleClearFilters = () => {
  searchQuery.value = ""
  advancedQuery.value = ""
}

// Export functions
const handleExportJSON = () => {
  isExporting.value = true
  const jsonData = JSON.stringify(filteredLogData.value, null, 2)
  const blob = new Blob([jsonData], { type: "application/json" })
  const url = URL.createObjectURL(blob)
  const a = document.createElement("a")
  a.href = url
  a.download = `logs2-${new Date().toISOString().split("T")[0]}.json`
  a.click()
  URL.revokeObjectURL(url)
  isExporting.value = false
}

// Define table columns using TanStack Table format
const columns: TableColumn<LogEntry>[] = [
  {
    id: "_expander",
    header: "",
    size: 36,
    enableSorting: false,
    enableHiding: false,
    cell: ({ row }) =>
      h(
        "button",
        {
          type: "button",
          class:
            "inline-flex items-center justify-center size-6 rounded-md border bg-background hover:bg-muted transition",
          "aria-label": row.getIsExpanded() ? "Collapse row" : "Expand row",
          onClick: (e: Event) => {
            e.stopPropagation()
            row.toggleExpanded()
          }
        },
        [
          h(ChevronRight, {
            class: ["h-4 w-4 transition-transform", row.getIsExpanded() ? "rotate-90" : "rotate-0"]
              .filter(Boolean)
              .join(" ")
          })
        ]
      ),
    meta: {
      class: {
        th: "w-9",
        td: "w-9"
      }
    }
  },
  {
    accessorKey: "timestamp",
    header: "Timestamp",
    cell: ({ getValue }) => {
      const date = new Date(getValue<string>())
      return h("div", { class: "text-xs font-mono leading-tight" }, [
        h("div", { class: "font-medium" }, date.toLocaleDateString()),
        h("div", { class: "text-muted-foreground" }, date.toLocaleTimeString())
      ])
    },
    size: 180
  },
  {
    accessorKey: "level",
    header: "Level",
    cell: ({ getValue }) => {
      const level = getValue<LogEntry["level"]>()
      const colorClass = getLevelColor(level)
      return h(
        "div",
        {
          class: `inline-flex items-center gap-1 px-1.5 py-0.5 rounded text-xs font-medium border ${colorClass}`
        },
        [h("div", { class: "size-1.5 rounded-full bg-current opacity-70" }), level.toUpperCase()]
      )
    },
    size: 100
  },
  {
    accessorKey: "source",
    header: "Source",
    cell: ({ getValue }) => {
      const source = getValue<string>()
      return h(
        "div",
        {
          class: "inline-flex items-center px-1.5 py-0.5 rounded text-xs font-medium border bg-muted/20 border-muted"
        },
        source
      )
    },
    size: 100
  },
  {
    accessorKey: "application",
    header: "Application",
    cell: ({ getValue }) => {
      const app = getValue<string>()
      return h(
        "div",
        {
          class: "font-medium text-xs truncate",
          title: app
        },
        app
      )
    },
    size: 150
  },
  {
    accessorKey: "message",
    header: "Message",
    cell: ({ getValue }) => {
      const message = getValue<string>()
      return h(
        "div",
        {
          class: "text-xs max-w-md truncate leading-tight",
          title: message
        },
        message
      )
    }
  },
  {
    accessorKey: "method",
    header: "Method",
    cell: ({ getValue }) => {
      const method = getValue<string>()
      const colorClass = getMethodColor(method)
      return h(
        "span",
        {
          class: `inline-flex items-center px-1.5 py-0.5 rounded text-xs font-medium border ${colorClass}`
        },
        method
      )
    },
    size: 80
  },
  {
    accessorKey: "statusCode",
    header: "Status",
    cell: ({ getValue }) => {
      const status = getValue<number>()
      const colorClass = getStatusColor(status)
      return h(
        "span",
        {
          class: `inline-flex items-center px-1.5 py-0.5 rounded text-xs font-medium border ${colorClass}`
        },
        status.toString()
      )
    },
    size: 80
  },
  {
    accessorKey: "responseTime",
    header: "Response Time",
    cell: ({ getValue }) => {
      const time = getValue<number>()
      const colorClass =
        time > 2000
          ? "text-red-600 bg-red-50"
          : time > 1000
            ? "text-yellow-600 bg-yellow-50"
            : "text-green-600 bg-green-50"
      return h(
        "span",
        {
          class: `text-xs font-mono ${colorClass} px-1 py-0.5 rounded leading-none`
        },
        `${time}ms`
      )
    },
    size: 120
  },
  {
    accessorKey: "tags",
    header: "Tags",
    cell: ({ getValue }) => {
      const tags = getValue<string[]>()
      const displayTags = tags.slice(0, 2)
      const remainingCount = Math.max(0, tags.length - 2)

      const tagElements = displayTags.map((tag) =>
        h(
          "span",
          {
            class:
              "inline-flex items-center px-1 py-0.5 rounded text-xs bg-secondary text-secondary-foreground border mr-1"
          },
          tag
        )
      )

      if (remainingCount > 0) {
        tagElements.push(h("span", { class: "text-xs text-muted-foreground" }, `+${remainingCount}`))
      }

      return h(
        "div",
        {
          class: "flex flex-wrap gap-0.5",
          title: tags.join(", ")
        },
        tagElements
      )
    },
    size: 200
  },
  {
    accessorKey: "requestId",
    header: "Request ID",
    cell: ({ getValue }) => {
      const id = getValue<string>()
      return h(
        "code",
        {
          class: "text-xs bg-muted px-1 py-0.5 rounded font-mono leading-none"
        },
        id
      )
    },
    size: 120
  }
]

// Table sorting - default to timestamp descending
const sorting = ref([{ id: "timestamp", desc: true }])
</script>

<template>
  <div class="flex flex-col gap-6 w-full max-w-full min-w-0">
    <PageHeader
      title="Logs (TanStack Table)"
      :subtext="`Monitor and analyze ${appConfig.app.name} gateway logs with TanStack Table`"
    >
      <div class="flex items-center gap-2">
        <UiDropdownMenu>
          <UiDropdownMenuTrigger as-child>
            <UiButton variant="outline" size="sm" class="gap-2">
              <Download class="h-4 w-4" />
              Export
            </UiButton>
          </UiDropdownMenuTrigger>
          <UiDropdownMenuContent align="end">
            <UiDropdownMenuItem @click="handleExportJSON" :disabled="isExporting">
              <Download class="mr-2 size-4" />
              {{ isExporting ? "Exporting..." : "Export JSON" }}
            </UiDropdownMenuItem>
          </UiDropdownMenuContent>
        </UiDropdownMenu>
      </div>
    </PageHeader>

    <UiCard>
      <UiCardHeader>
        <UiCardTitle class="flex items-center gap-2">
          <Filter class="h-5 w-5" />
          Search & Filters
        </UiCardTitle>
        <UiCardDescription> Search logs and apply filters to find specific events </UiCardDescription>
      </UiCardHeader>
      <UiCardContent>
        <div class="flex flex-col gap-4">
          <!-- Query Mode Toggle -->
          <div class="flex items-center gap-2">
            <UiButton variant="outline" size="sm" class="gap-2" @click="showAdvancedQuery = !showAdvancedQuery">
              <Code class="h-4 w-4" />
              {{ showAdvancedQuery ? "Simple Search" : "Advanced Query" }}
              <component :is="showAdvancedQuery ? ChevronUp : ChevronDown" class="h-4 w-4" />
            </UiButton>
          </div>

          <!-- Simple Search (default) -->
          <div v-if="!showAdvancedQuery" class="flex items-center gap-2">
            <UiInput
              v-model="searchQuery"
              placeholder="Search logs by message, application, tags, request ID..."
              class="flex-1"
            />
            <UiButton v-if="searchQuery" variant="outline" size="sm" @click="handleClearFilters"> Clear </UiButton>
          </div>

          <!-- Advanced Query Interface -->
          <div v-else class="space-y-4">
            <!-- Query Input -->
            <div class="flex items-start gap-2">
              <UiTextarea
                v-model="advancedQuery"
                placeholder='Enter KQL-like query, e.g.: level == "error" | where responseTime > 1000 | order by timestamp desc | limit 50'
                class="flex-1 font-mono text-sm min-h-[80px]"
              />
              <UiButton class="gap-2 mt-1" @click="handleRunQuery" :disabled="!advancedQuery.trim()">
                <Play class="h-4 w-4" />
                Run
              </UiButton>
            </div>

            <!-- Query History -->
            <div v-if="queryHistory.length > 0">
              <div class="flex items-center gap-2 mb-2">
                <History class="h-4 w-4" />
                <span class="text-sm font-medium">Recent Queries</span>
              </div>
              <div class="space-y-1">
                <UiButton
                  v-for="historyQuery in queryHistory.slice(0, 5)"
                  :key="historyQuery"
                  variant="ghost"
                  size="sm"
                  class="w-full justify-start text-left h-auto py-1"
                  @click="loadHistoryQuery(historyQuery)"
                >
                  <span class="text-xs font-mono truncate">{{ historyQuery }}</span>
                </UiButton>
              </div>
            </div>

            <!-- Query Syntax Help -->
            <div class="text-xs text-muted-foreground bg-muted/20 p-3 rounded-lg">
              <div class="font-medium mb-1">Query Syntax Examples:</div>
              <div class="space-y-1 font-mono">
                <div><strong>Filter by level:</strong> level == "error"</div>
                <div><strong>Time range:</strong> timestamp > ago(1h)</div>
                <div><strong>Response time:</strong> responseTime > 2000</div>
                <div><strong>Contains text:</strong> tags contains "rate-limit"</div>
                <div><strong>Status codes:</strong> statusCode >= 400</div>
                <div><strong>Sorting:</strong> order by responseTime desc</div>
                <div><strong>Limit results:</strong> limit 50</div>
                <div>
                  <strong>Combine:</strong> level == "error" | where responseTime > 1000 | order by timestamp desc
                </div>
              </div>
            </div>
          </div>
        </div>
      </UiCardContent>
    </UiCard>

    <!-- Logs Table using TanStack Table -->
    <BlocksDataTable :data="filteredLogData" :columns="columns" v-model:sorting="sorting" class="border rounded-lg">
      <template #expanded="{ row }">
        <div class="bg-slate-50 p-4 rounded-lg space-y-4 text-sm border">
          <div class="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-4">
            <div class="space-y-3">
              <h4 class="font-semibold text-slate-900 border-b border-slate-200 pb-1">Request Details</h4>
              <div class="space-y-2 text-xs">
                <div>
                  <span class="font-medium text-slate-700">Request ID:</span>
                  <code class="bg-white px-2 py-1 rounded text-slate-800 border">{{ row.original.requestId }}</code>
                </div>
                <div>
                  <span class="font-medium text-slate-700">Method:</span>
                  <span
                    :class="`inline-flex items-center px-2 py-1 rounded text-xs font-medium border ${getMethodColor(row.original.method)}`"
                    >{{ row.original.method }}</span
                  >
                </div>
                <div>
                  <span class="font-medium text-slate-700">Endpoint:</span>
                  <code class="bg-white px-2 py-1 rounded text-slate-800 border">{{ row.original.endpoint }}</code>
                </div>
                <div>
                  <span class="font-medium text-slate-700">Status Code:</span>
                  <span
                    :class="`inline-flex items-center px-2 py-1 rounded text-xs font-medium border ${getStatusColor(row.original.statusCode)}`"
                    >{{ row.original.statusCode }}</span
                  >
                </div>
              </div>
            </div>

            <div class="space-y-3">
              <h4 class="font-semibold text-slate-900 border-b border-slate-200 pb-1">Application & User</h4>
              <div class="space-y-2 text-xs">
                <div>
                  <span class="font-medium text-slate-700">Application:</span>
                  <span class="text-slate-900 font-medium">{{ row.original.application }}</span>
                </div>
                <div>
                  <span class="font-medium text-slate-700">User ID:</span>
                  <span class="text-slate-800">{{ row.original.userId || "N/A" }}</span>
                </div>
                <div>
                  <span class="font-medium text-slate-700">Source:</span>
                  <span
                    class="inline-flex items-center px-2 py-1 rounded text-xs font-medium border bg-slate-100 border-slate-300 text-slate-700"
                    >{{ row.original.source }}</span
                  >
                </div>
                <div>
                  <span class="font-medium text-slate-700">Response Time:</span>
                  <span
                    :class="`font-mono ${row.original.responseTime > 2000 ? 'text-red-600' : row.original.responseTime > 1000 ? 'text-yellow-600' : 'text-green-600'}`"
                    >{{ row.original.responseTime }}ms</span
                  >
                </div>
              </div>
            </div>

            <div class="space-y-3">
              <h4 class="font-semibold text-slate-900 border-b border-slate-200 pb-1">Metadata</h4>
              <div class="space-y-2 text-xs">
                <div>
                  <span class="font-medium text-slate-700">Timestamp:</span>
                  <span class="font-mono text-slate-800">{{ new Date(row.original.timestamp).toLocaleString() }}</span>
                </div>
                <div>
                  <span class="font-medium text-slate-700">Level:</span>
                  <span
                    :class="`inline-flex items-center gap-1 px-2 py-1 rounded text-xs font-medium border ${getLevelColor(row.original.level)}`"
                    ><div class="w-1.5 h-1.5 rounded-full bg-current opacity-70"></div>
                    {{ row.original.level.toUpperCase() }}</span
                  >
                </div>
                <div>
                  <span class="font-medium text-slate-700">Tags:</span>
                  <div class="flex flex-wrap gap-1 mt-1">
                    <span
                      v-for="tag in row.original.tags"
                      :key="tag"
                      class="inline-flex items-center px-2 py-1 rounded text-xs bg-slate-200 text-slate-700 border border-slate-300"
                    >
                      {{ tag }}
                    </span>
                  </div>
                </div>
              </div>
            </div>
          </div>

          <div class="space-y-3 pt-2 border-t border-slate-200">
            <h4 class="font-semibold text-slate-900">Message</h4>
            <div class="bg-white p-3 rounded border border-slate-200 font-mono text-xs leading-relaxed text-slate-800">
              {{ row.original.message }}
            </div>
          </div>
        </div>
      </template>
    </BlocksDataTable>
  </div>
</template>
