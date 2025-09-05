<script setup lang="ts">
import { Key, Activity, Calendar, Users, Settings, Layers, Globe } from "lucide-vue-next"
import type { StatsCardProps } from "~/components/Cards/Stats.vue"

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

const getStatusBadgeClass = (status: string) => {
  return status === "active"
    ? "bg-green-100 text-green-800 dark:bg-green-900 dark:text-green-200"
    : "bg-gray-100 text-gray-800 dark:bg-gray-900 dark:text-gray-200"
}

// Icon and variant arrays for index matching with API data
const statsIcons = [Key, Activity, Layers, Users]
const statsVariants: StatsCardProps["variant"][] = ["chart-3", "chart-1", "default", "chart-2"]

// Stats configuration - ready for API replacement
const statsCards = computed(() => {
  // This can be replaced with: const { data: statsData } = await useFetch(`/api/applications/${appId}/stats`)
  // Then: statsData.map((stat, index) => ({ ...stat, icon: statsIcons[index], variant: statsVariants[index] }))

  return [
    {
      title: "API Keys",
      value: application.value.apiKeyCount,
      description: "Active keys"
    },
    {
      title: "Monthly Requests",
      value: application.value.monthlyRequests,
      description: "This month"
    },
    {
      title: "Models",
      value: application.value.models.length,
      description: application.value.models.join(", ")
    },
    {
      title: "Team",
      value: application.value.team,
      description: "Responsible team"
    }
  ].map((stat, index) => ({
    ...stat,
    icon: statsIcons[index],
    variant: statsVariants[index]
  }))
})
</script>

<template>
  <div class="flex flex-col gap-6">
    <!-- Header -->

    <PageBadgeHeader :badge-status="application.status" :title="application.name" :subtext="application.description">
      <!-- Action Buttons -->

      <div class="flex items-center gap-2">
        <UiButton variant="default" class="gap-2" @click="navigateTo(`/applications/${appId}/keys`)">
          <Key class="h-4 w-4" />
          Manage Keys
        </UiButton>
        <UiButton variant="outline" class="gap-2">
          <Settings class="h-4 w-4" />
          Settings
        </UiButton>
      </div>
    </PageBadgeHeader>

    <!-- Application Overview Cards -->
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
