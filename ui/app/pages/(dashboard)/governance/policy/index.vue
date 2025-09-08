<script lang="ts">
interface Policy {
  id: string
  name: string
  description: string
  type: "rate_limit" | "content_filter" | "access_control" | "cost_limit"
  status: "active" | "inactive"
  priority: "high" | "medium" | "low"
  applicationsCount: number
  lastModified: string
  createdBy: string
}

const policies = ref<Policy[]>([
  {
    id: "policy_1",
    name: "Rate Limiting - Production",
    description: "Enforce strict rate limits for production applications to prevent abuse",
    type: "rate_limit",
    status: "active",
    priority: "high",
    applicationsCount: 8,
    lastModified: "2025-01-15T14:30:00Z",
    createdBy: "Alice Johnson"
  },
  {
    id: "policy_2", 
    name: "Content Safety Filter",
    description: "Block harmful or inappropriate content across all AI interactions",
    type: "content_filter",
    status: "active",
    priority: "high",
    applicationsCount: 12,
    lastModified: "2025-01-14T09:15:00Z",
    createdBy: "Bob Smith"
  },
  {
    id: "policy_3",
    name: "Development Access Control",
    description: "Restrict development environment access to authorized team members only",
    type: "access_control", 
    status: "active",
    priority: "medium",
    applicationsCount: 3,
    lastModified: "2025-01-13T16:45:00Z",
    createdBy: "Carol Williams"
  },
  {
    id: "policy_4",
    name: "Monthly Cost Cap",
    description: "Prevent applications from exceeding monthly spending limits",
    type: "cost_limit",
    status: "inactive",
    priority: "medium",
    applicationsCount: 0,
    lastModified: "2025-01-10T11:20:00Z",
    createdBy: "David Brown"
  }
])
</script>
<script setup lang="ts">
import { Plus, Shield, AlertTriangle, CheckCircle, XCircle, Clock, Users, FileText, Filter, Zap } from "lucide-vue-next"
import type { StatsCardProps } from "@/components/Cards/Stats.vue"
import SearchFilter from "@/components/SearchFilter.vue"
import type { FilterConfig, SearchConfig, DisplayConfig } from "@/components/SearchFilter.vue"

const showCreateModal = ref(false)
const activeFilters = ref<Record<string, string>>({})

const appConfig = useAppConfig()

const handlePolicySelect = (policy: Policy) => {
  navigateTo(`/governance/policy/${policy.id}`)
}

const onPolicyCreated = (policyId: string) => {
  console.log("Policy created:", policyId)
  navigateTo("/governance/policy", { replace: true })
}

const route = useRoute()
onMounted(() => {
  if (route.query.create === "policy") {
    showCreateModal.value = true
  }
})

const policyTypes = [...new Set(policies.value.map(p => p.type))]
const priorities = [...new Set(policies.value.map(p => p.priority))]

const filterConfigs: FilterConfig[] = [
  {
    key: "name",
    label: "Policy",
    options: policies.value.map(policy => ({ value: policy.name, label: policy.name, icon: Shield }))
  },
  {
    key: "type",
    label: "Type",
    options: [
      { value: "rate_limit", label: "Rate Limiting", icon: Clock },
      { value: "content_filter", label: "Content Filter", icon: Filter },
      { value: "access_control", label: "Access Control", icon: Shield },
      { value: "cost_limit", label: "Cost Limit", icon: Zap }
    ]
  },
  {
    key: "status",
    label: "Status",
    options: [
      { value: "active", label: "Active", icon: CheckCircle },
      { value: "inactive", label: "Inactive", icon: XCircle }
    ]
  },
  {
    key: "priority",
    label: "Priority",
    options: [
      { value: "high", label: "High", icon: AlertTriangle },
      { value: "medium", label: "Medium", icon: FileText },
      { value: "low", label: "Low", icon: FileText }
    ]
  }
]

const searchConfig: SearchConfig<Policy> = {
  fields: ["name", "description"],
  placeholder: "Search policies, filter by type, status, priority..."
}

const displayConfig: DisplayConfig<Policy> = {
  getItemText: (policy) => `${policy.name} - ${policy.description}`,
  getItemValue: (policy) => policy.name,
  getItemIcon: () => Shield
}

function handleFiltersChanged(filters: Record<string, string>) {
  activeFilters.value = filters
}

function handleItemSelected(policy: Policy) {
  navigateTo(`/governance/policy/${policy.id}`)
}

const filteredPolicies = computed(() => {
  let filtered = policies.value

  if (activeFilters.value.name && activeFilters.value.name !== "all") {
    filtered = filtered.filter(policy => policy.name === activeFilters.value.name)
  }
  if (activeFilters.value.type && activeFilters.value.type !== "all") {
    filtered = filtered.filter(policy => policy.type === activeFilters.value.type)
  }
  if (activeFilters.value.status && activeFilters.value.status !== "all") {
    filtered = filtered.filter(policy => policy.status === activeFilters.value.status)
  }
  if (activeFilters.value.priority && activeFilters.value.priority !== "all") {
    filtered = filtered.filter(policy => policy.priority === activeFilters.value.priority)
  }

  return filtered
})

const statsIcons = [Shield, CheckCircle, AlertTriangle, Users]
const statsVariants: StatsCardProps["variant"][] = ["default", "chart-1", "chart-2", "chart-3"]

const statsCards = computed(() => {
  return [
    {
      title: "Total Policies",
      value: policies.value.length,
      description: "All governance policies"
    },
    {
      title: "Active Policies", 
      value: policies.value.filter(p => p.status === "active").length,
      description: "Currently enforced policies"
    },
    {
      title: "High Priority",
      value: policies.value.filter(p => p.priority === "high").length,
      description: "Critical governance rules"
    },
    {
      title: "Applications Covered",
      value: policies.value.reduce((sum, policy) => sum + policy.applicationsCount, 0),
      description: "Total application assignments"
    }
  ].map((stat, index) => ({
    ...stat,
    icon: statsIcons[index],
    variant: statsVariants[index]
  }))
})

const getPolicyTypeIcon = (type: Policy["type"]) => {
  switch (type) {
    case "rate_limit": return Clock
    case "content_filter": return Filter
    case "access_control": return Shield
    case "cost_limit": return Zap
    default: return Shield
  }
}

const getPolicyTypeColor = (type: Policy["type"]) => {
  switch (type) {
    case "rate_limit": return "text-blue-600 bg-blue-50 border-blue-200"
    case "content_filter": return "text-purple-600 bg-purple-50 border-purple-200"
    case "access_control": return "text-green-600 bg-green-50 border-green-200"
    case "cost_limit": return "text-orange-600 bg-orange-50 border-orange-200"
    default: return "text-gray-600 bg-gray-50 border-gray-200"
  }
}

const getPriorityColor = (priority: Policy["priority"]) => {
  switch (priority) {
    case "high": return "text-red-600 bg-red-50 border-red-200"
    case "medium": return "text-yellow-600 bg-yellow-50 border-yellow-200"
    case "low": return "text-green-600 bg-green-50 border-green-200"
    default: return "text-gray-600 bg-gray-50 border-gray-200"
  }
}

const getStatusColor = (status: Policy["status"]) => {
  switch (status) {
    case "active": return "text-green-600 bg-green-50 border-green-200"
    case "inactive": return "text-gray-600 bg-gray-50 border-gray-200"
    default: return "text-gray-600 bg-gray-50 border-gray-200"
  }
}
</script>

<template>
  <div class="flex flex-col gap-6">
    <PageHeader 
      title="Governance Policies"
      :subtext="`Manage governance rules and policies for ${appConfig.app.name} applications`"
    >
      <UiButton @click="showCreateModal = true" class="gap-2">
        <Plus class="h-4 w-4" />
        New Policy
      </UiButton>
    </PageHeader>

    <div class="grid grid-cols-1 sm:grid-cols-2 lg:grid-cols-4 gap-6">
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

    <SearchFilter
      :items="policies"
      :filters="filterConfigs"
      :search-config="searchConfig"
      :display-config="displayConfig"
      @filters-changed="handleFiltersChanged"
      @item-selected="handleItemSelected"
    />

    <CardsDataList title="All Policies" :icon="Shield">
      <template #actions></template>

      <div class="space-y-4">
        <div 
          v-for="policy in filteredPolicies"
          :key="policy.id"
          class="flex items-center justify-between p-4 border rounded-lg hover:bg-muted/50 cursor-pointer"
          @click="handlePolicySelect(policy)"
        >
          <div class="flex items-center gap-4">
            <div class="p-2 rounded-lg bg-primary/10">
              <component :is="getPolicyTypeIcon(policy.type)" class="size-5 text-primary" />
            </div>

            <div class="space-y-1">
              <div class="flex items-center gap-2">
                <p class="font-medium">{{ policy.name }}</p>
                <div 
                  class="inline-flex items-center gap-1 px-2 py-1 rounded-md text-xs font-medium border"
                  :class="getPolicyTypeColor(policy.type)"
                >
                  <component :is="getPolicyTypeIcon(policy.type)" class="size-3" />
                  {{ policy.type.replace('_', ' ').replace(/\b\w/g, l => l.toUpperCase()) }}
                </div>
              </div>
              <p class="text-sm text-muted-foreground">{{ policy.description }}</p>
              <div class="flex items-center gap-4 text-xs text-muted-foreground">
                <span>{{ policy.applicationsCount }} applications</span>
                <span>•</span>
                <span>Modified {{ new Date(policy.lastModified).toLocaleDateString() }}</span>
                <span>•</span>
                <span>Created by {{ policy.createdBy }}</span>
              </div>
            </div>
          </div>

          <div class="flex items-center gap-2">
            <div
              class="inline-flex items-center gap-1 px-2 py-1 rounded-md text-xs font-medium border"
              :class="getPriorityColor(policy.priority)"
            >
              <component :is="policy.priority === 'high' ? AlertTriangle : FileText" class="size-3" />
              {{ policy.priority.charAt(0).toUpperCase() + policy.priority.slice(1) }}
            </div>
            <div
              class="inline-flex items-center gap-1 px-2 py-1 rounded-md text-xs font-medium border"
              :class="getStatusColor(policy.status)"
            >
              <div class="size-1.5 rounded-full" :class="policy.status === 'active' ? 'bg-green-500' : 'bg-gray-400'" />
              {{ policy.status === "active" ? "Active" : "Inactive" }}
            </div>
          </div>
        </div>
      </div>
    </CardsDataList>
  </div>
</template>
