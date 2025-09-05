<script setup lang="ts">
import { Key, Activity, Calendar, Users, Settings, Layers, Globe } from "lucide-vue-next"

// Get the app ID from route params
const route = useRoute()
const appId = route.params.appId as string

// Sample application data - this would come from an API call
const application = ref({
  id: appId,
  name: "Customer22 Service Bot",
  description: "AI-powered customer support assistant for handling common inquiries",
  status: "active",
  apiKeyCount: 3,
  monthlyRequests: 45200,
  lastUsed: "2025-01-15T10:30:00Z",
  models: ["gpt-4", "gpt-3.5-turbo"],
  team: "Customer Success",
  created: "2024-10-15T09:00:00Z",
  owners: [
    {
      id: "user_1",
      name: "John Smith",
      email: "john.smith@company.com",
      role: "Product Owner"
    },
    {
      id: "user_2",
      name: "Sarah Johnson",
      email: "sarah.johnson@company.com",
      role: "Technical Lead"
    }
  ],
  recentActivity: [
    {
      timestamp: "2025-01-15T14:30:00Z",
      action: "API Key Created",
      details: "Production Key created by John Smith"
    },
    {
      timestamp: "2025-01-15T12:15:00Z",
      action: "High Usage Alert",
      details: "Monthly request limit at 85%"
    },
    {
      timestamp: "2025-01-14T16:45:00Z",
      action: "Model Updated",
      details: "Added GPT-4 Turbo to available models"
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
</script>

<template>
  <div class="flex flex-col gap-6">
    <!-- Header -->
    <div class="flex items-start justify-between">
      <div>
        <div class="flex items-center gap-3 mb-2">
          <h1 class="text-3xl font-bold tracking-tight">{{ application.name }}</h1>
          <UiBadge :class="getStatusBadgeClass(application.status)">
            {{ application.status }}
          </UiBadge>
        </div>
        <p class="text-muted-foreground">{{ application.description }}</p>
      </div>

      <!-- Action Buttons -->
      <div class="flex items-center gap-2">
        <UiButton variant="default" size="sm" class="gap-2" @click="navigateTo(`/applications/${appId}/keys`)">
          <Key class="h-4 w-4" />
          Manage Keys
        </UiButton>
        <UiButton variant="outline" size="sm" class="gap-2">
          <Settings class="h-4 w-4" />
          Settings
        </UiButton>
      </div>
    </div>

    <!-- Application Overview Cards -->
    <div class="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-4 gap-6">
      <UiCard>
        <UiCardHeader class="flex flex-row items-center justify-between space-y-0 pb-2">
          <UiCardTitle class="text-sm font-medium">API Keys</UiCardTitle>
          <Key class="h-4 w-4 text-muted-foreground" />
        </UiCardHeader>
        <UiCardContent>
          <div class="text-2xl font-bold">{{ application.apiKeyCount }}</div>
          <p class="text-xs text-muted-foreground">Active keys</p>
        </UiCardContent>
      </UiCard>

      <UiCard>
        <UiCardHeader class="flex flex-row items-center justify-between space-y-0 pb-2">
          <UiCardTitle class="text-sm font-medium">Monthly Requests</UiCardTitle>
          <Activity class="h-4 w-4 text-muted-foreground" />
        </UiCardHeader>
        <UiCardContent>
          <div class="text-2xl font-bold">{{ formatNumber(application.monthlyRequests) }}</div>
          <p class="text-xs text-muted-foreground">This month</p>
        </UiCardContent>
      </UiCard>

      <UiCard>
        <UiCardHeader class="flex flex-row items-center justify-between space-y-0 pb-2">
          <UiCardTitle class="text-sm font-medium">Models</UiCardTitle>
          <Layers class="h-4 w-4 text-muted-foreground" />
        </UiCardHeader>
        <UiCardContent>
          <div class="text-2xl font-bold">{{ application.models.length }}</div>
          <p class="text-xs text-muted-foreground">{{ application.models.join(", ") }}</p>
        </UiCardContent>
      </UiCard>

      <UiCard>
        <UiCardHeader class="flex flex-row items-center justify-between space-y-0 pb-2">
          <UiCardTitle class="text-sm font-medium">Team</UiCardTitle>
          <Users class="h-4 w-4 text-muted-foreground" />
        </UiCardHeader>
        <UiCardContent>
          <div class="text-2xl font-bold">{{ application.team }}</div>
          <p class="text-xs text-muted-foreground">Responsible team</p>
        </UiCardContent>
      </UiCard>
    </div>

    <!-- Application Details -->
    <UiCard>
      <UiCardHeader>
        <UiCardTitle class="flex items-center gap-2">
          <Globe class="h-5 w-5" />
          Application Details
        </UiCardTitle>
      </UiCardHeader>
      <UiCardContent class="space-y-6">
        <!-- Application Info Grid -->
        <div class="grid grid-cols-1 md:grid-cols-3 gap-6">
          <div>
            <label class="text-sm font-medium text-muted-foreground flex items-center gap-2">
              <Calendar class="h-4 w-4" />
              Created
            </label>
            <p class="text-sm font-medium mt-1">{{ new Date(application.created).toLocaleDateString() }}</p>
          </div>

          <div>
            <label class="text-sm font-medium text-muted-foreground flex items-center gap-2">
              <Activity class="h-4 w-4" />
              Last Used
            </label>
            <p class="text-sm font-medium mt-1">{{ new Date(application.lastUsed).toLocaleDateString() }}</p>
          </div>

          <div>
            <label class="text-sm font-medium text-muted-foreground flex items-center gap-2">
              <Users class="h-4 w-4" />
              Owners
            </label>
            <div class="mt-1 flex flex-wrap gap-1">
              <div v-for="owner in application.owners" :key="owner.id" class="flex items-center gap-1 text-xs">
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
        <div>
          <label class="text-sm font-medium text-muted-foreground">Description</label>
          <p class="text-sm mt-1">{{ application.description }}</p>
        </div>
      </UiCardContent>
    </UiCard>

    <!-- Recent Activity -->
    <UiCard>
      <UiCardHeader>
        <UiCardTitle class="flex items-center gap-2">
          <Activity class="h-5 w-5" />
          Recent Activity
        </UiCardTitle>
        <UiCardDescription>Latest events and actions for this application</UiCardDescription>
      </UiCardHeader>
      <UiCardContent>
        <RecentActivity :activities="application.recentActivity" />
      </UiCardContent>
    </UiCard>
  </div>
</template>
