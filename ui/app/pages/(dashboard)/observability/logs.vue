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
import { Download, Filter, Code, History, Play, ChevronDown, ChevronUp } from "lucide-vue-next"
import type { Config } from "datatables.net"
import type DataTableRef from "datatables.net"

const appConfig = useAppConfig()

const tableRef = shallowRef<InstanceType<typeof DataTableRef<any[]>> | null>(null)

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

const searchQuery = ref("")
const advancedQuery = ref("")
const isExporting = ref(false)
const showAdvancedQuery = ref(false)
const queryHistory = ref<string[]>([])
// KQL-like query parsing
const parseAdvancedQuery = (query: string) => {
  if (!query.trim()) return logData

  let filteredData = [...logData]

  // Simple query parser for demo purposes
  const parts = query.split("|").map((p) => p.trim())

  for (const part of parts) {
    if (
      part.startsWith("where ") ||
      part.includes("==") ||
      part.includes("!=") ||
      part.includes(">") ||
      part.includes("<")
    ) {
      // Handle where clauses
      filteredData = applyWhereClause(filteredData, part)
    } else if (part.startsWith("order by ")) {
      // Handle ordering
      filteredData = applyOrderBy(filteredData, part)
    } else if (part.startsWith("limit ") || part.startsWith("take ")) {
      // Handle limit
      filteredData = applyLimit(filteredData, part)
    } else if (part.includes("contains")) {
      // Handle contains operations
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
  return logData
})

const handleSearch = (query: string) => {
  if (tableRef.value) {
    tableRef.value.search(query).draw()
  }
}

const handleExportCopy = () => {
  if (tableRef.value) {
    tableRef.value.buttons("copy:name").trigger()
  }
}

const handleExportCSV = () => {
  if (tableRef.value) {
    tableRef.value.buttons("csv:name").trigger()
  }
}

const handleExportExcel = () => {
  if (tableRef.value) {
    tableRef.value.buttons("excel:name").trigger()
  }
}

const handleExportJSON = () => {
  if (tableRef.value) {
    isExporting.value = true
    const data = tableRef.value.data().toArray()
    const jsonData = JSON.stringify(data, null, 2)
    const blob = new Blob([jsonData], { type: "application/json" })
    const url = URL.createObjectURL(blob)
    const a = document.createElement("a")
    a.href = url
    a.download = `logs-${new Date().toISOString().split("T")[0]}.json`
    a.click()
    URL.revokeObjectURL(url)
    isExporting.value = false
  }
}

const handleClearFilters = () => {
  searchQuery.value = ""
  advancedQuery.value = ""
  if (tableRef.value) {
    tableRef.value.search("").columns().search("").draw()
  }
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

const loadHistoryQuery = (query: string) => {
  advancedQuery.value = query
}

watch(searchQuery, (newQuery) => {
  if (!advancedQuery.value.trim()) {
    handleSearch(newQuery)
  }
})

const options: Config = {
  dom: "t<'flex flex-col lg:flex-row gap-5 lg:items-center lg:justify-between pt-3 p-5'li><''p>",
  select: true,
  autoWidth: true,
  responsive: true,
  order: [[0, "desc"]], // Order by timestamp descending (most recent first)
  pageLength: 25,
  lengthMenu: [
    [10, 25, 50, 100],
    [10, 25, 50, 100]
  ],
  searching: true,

  buttons: [
    { extend: "copy", name: "copy" },
    { extend: "csv", name: "csv" },
    { extend: "excel", name: "excel" }
  ],

  columns: [
    {
      data: "timestamp",
      title: "Timestamp",
      width: "180px",
      render(data, _, _row, ___) {
        const date = new Date(data)
        return `
          <div class="text-xs font-mono">
            <div class="font-medium">${date.toLocaleDateString()}</div>
            <div class="text-muted-foreground">${date.toLocaleTimeString()}</div>
          </div>
        `
      }
    },
    {
      data: "level",
      title: "Level",
      width: "100px",
      render(data, _, _row, ___) {
        const colorClass = getLevelColor(data)
        return `
          <div class="inline-flex items-center gap-1 px-2 py-1 rounded-md text-xs font-medium border ${colorClass}">
            <div class="size-2 rounded-full bg-current opacity-60"></div>
            ${data.toUpperCase()}
          </div>
        `
      }
    },
    {
      data: "source",
      title: "Source",
      width: "100px",
      render(data, _, _row, ___) {
        return `
          <div class="inline-flex items-center gap-1 px-2 py-1 rounded-md text-xs font-medium border bg-muted/20 border-muted">
            ${data}
          </div>
        `
      }
    },
    {
      data: "application",
      title: "Application",
      width: "150px",
      render(data, _, _row, ___) {
        return `<div class="font-medium text-sm truncate" title="${data}">${data}</div>`
      }
    },
    {
      data: "message",
      title: "Message",
      render(data, _, _row, ___) {
        return `<div class="text-sm max-w-md truncate" title="${data}">${data}</div>`
      }
    },
    {
      data: "method",
      title: "Method",
      width: "80px",
      render(data, _, _row, ___) {
        const methodColors = {
          GET: "text-blue-600 bg-blue-50 border-blue-200",
          POST: "text-green-600 bg-green-50 border-green-200",
          PUT: "text-yellow-600 bg-yellow-50 border-yellow-200",
          DELETE: "text-red-600 bg-red-50 border-red-200"
        }
        const colorClass = methodColors[data as keyof typeof methodColors] || "text-gray-600 bg-gray-50 border-gray-200"
        return `<span class="inline-flex items-center px-2 py-1 rounded text-xs font-medium border ${colorClass}">${data}</span>`
      }
    },
    {
      data: "statusCode",
      title: "Status",
      width: "80px",
      render(data, _, _row, ___) {
        const colorClass = getStatusColor(data)
        return `<span class="inline-flex items-center px-2 py-1 rounded text-xs font-medium border ${colorClass}">${data}</span>`
      }
    },
    {
      data: "responseTime",
      title: "Response Time",
      width: "120px",
      render(data, _, _row, ___) {
        const isSlowResponse = data > 2000
        const colorClass = isSlowResponse
          ? "text-red-600 bg-red-50"
          : data > 1000
            ? "text-yellow-600 bg-yellow-50"
            : "text-green-600 bg-green-50"
        return `<span class="text-xs font-mono ${colorClass} px-1 py-0.5 rounded">${data}ms</span>`
      }
    },
    {
      data: "tags",
      title: "Tags",
      width: "200px",
      render(data, _, _row, ___) {
        const tags = Array.isArray(data) ? data : []
        const displayTags = tags.slice(0, 2)
        const remainingCount = Math.max(0, tags.length - 2)

        let html = displayTags
          .map(
            (tag) =>
              `<span class="inline-flex items-center px-1.5 py-0.5 rounded text-xs bg-secondary text-secondary-foreground border mr-1">${tag}</span>`
          )
          .join("")

        if (remainingCount > 0) {
          html += `<span class="text-xs text-muted-foreground">+${remainingCount}</span>`
        }

        return `<div class="flex flex-wrap gap-1" title="${tags.join(", ")}">${html}</div>`
      }
    },
    {
      data: "requestId",
      title: "Request ID",
      width: "120px",
      render(data, _, _row, ___) {
        return `<code class="text-xs bg-muted px-1.5 py-0.5 rounded font-mono">${data}</code>`
      }
    }
  ]
}
</script>

<template>
  <div class="flex flex-col gap-6">
    <PageHeader title="Logs" :subtext="`Monitor and analyze ${appConfig.app.name} gateway logs in real-time`">
      <div class="flex items-center gap-2">
        <UiDropdownMenu>
          <UiDropdownMenuTrigger as-child>
            <UiButton variant="outline" size="sm" class="gap-2">
              <Download class="h-4 w-4" />
              Export
            </UiButton>
          </UiDropdownMenuTrigger>
          <UiDropdownMenuContent align="end">
            <UiDropdownMenuItem @click="handleExportCopy">
              <Download class="mr-2 size-4" />
              Copy to Clipboard
            </UiDropdownMenuItem>
            <UiDropdownMenuItem @click="handleExportCSV">
              <Download class="mr-2 size-4" />
              Export CSV
            </UiDropdownMenuItem>
            <UiDropdownMenuItem @click="handleExportExcel">
              <Download class="mr-2 size-4" />
              Export Excel
            </UiDropdownMenuItem>
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

    <UiDatatable :data="filteredLogData" :options="options" @ready="tableRef = $event" class="w-full" />
  </div>
</template>
