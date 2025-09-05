<script setup lang="ts">
import { Key, Layers, Activity, Users, ChartArea, Zap } from "lucide-vue-next"
import type { FunctionalComponent } from "vue"
import type { ButtonProps } from "./Ui/Button.vue"

const cards: {
  title: string
  icon: FunctionalComponent
  value: string
  change: string
  changeColor: string
}[] = [
  {
    title: "Active Applications",
    icon: Layers,
    value: "12",
    change: "+2 this week",
    changeColor: "text-green-600"
  },
  {
    title: "API Keys",
    icon: Key,
    value: "47",
    change: "+5 this month",
    changeColor: "text-green-600"
  },
  {
    title: "Total Users",
    icon: Users,
    value: "1,234",
    change: "+12% from last month",
    changeColor: "text-green-600"
  },
  {
    title: "Requests Today",
    icon: ChartArea,
    value: "89.2K",
    change: "-2% from yesterday",
    changeColor: "text-red-500"
  }
]

const quickActions: {
  label: string
  icon: FunctionalComponent
  class: string
  variant: ButtonProps["variant"]
}[] = [
  {
    label: "Register New Application",
    icon: Layers,
    class: "w-full justify-start",
    variant: "default"
  },
  {
    label: "Generate API Key",
    icon: Key,
    variant: "outline",
    class: "justify-start w-full"
  },
  {
    label: "Invite User",
    icon: Users,
    variant: "outline",
    class: "justify-start w-full"
  },
  {
    label: "View Analytics",
    icon: ChartArea,
    variant: "outline",
    class: "justify-start w-full"
  }
]
const resourceUsage = {
  status: "Normal",
  cpu: 45,
  memory: 62
}

// Quick action handlers
const handleQuickAction = (actionLabel: string) => {
  switch (actionLabel) {
    case "Register New Application":
      navigateTo("/applications?create=application")
      break
    case "Generate API Key":
      navigateTo("/applications/keys?create=apikey")
      break
    case "Invite User":
      navigateTo("/users?create=user")
      break
    case "View Analytics":
      navigateTo("/metrics")
      break
  }
}
</script>

<template>
  <section class="flex flex-col gap-8">
    <!-- Header -->
    <div>
      <h1 class="text-3xl font-bold tracking-tight text-primary">Dashboard</h1>
      <p class="text-muted-foreground">Overview of your LLM proxy service performance and activity</p>
    </div>
    <!-- Stats -->
    <div class="grid grid-cols-1 sm:grid-cols-2 lg:grid-cols-4 gap-6">
      <CardsStats
        v-for="card in cards"
        :key="card.title"
        :title="card.title"
        :value="card.value"
        :icon="card.icon"
        :description="card.change"
      />
    </div>
    <!-- Main Content -->
    <div class="grid grid-cols-1 md:grid-cols-2 gap-6">
      <!-- Recent Activity -->
      <UiCard>
        <UiCardHeader>
          <div class="flex items-center gap-2">
            <Activity class="h-5 w-5 text-primary" />
            <UiCardTitle>Recent Activity</UiCardTitle>
          </div>
        </UiCardHeader>
        <UiCardContent>
          <div class="flex flex-col gap-4">
            <div class="flex gap-3 items-start">
              <span class="mt-1 w-3 h-3 rounded-full bg-chart-2"></span>
              <div>
                <div class="font-semibold text-foreground">New application registered</div>
                <div class="text-muted-foreground text-sm">ChatBot Pro by ACME Corp</div>
                <div class="text-xs text-muted-foreground">2 minutes ago</div>
              </div>
            </div>
            <div class="flex gap-3 items-start">
              <span class="mt-1 w-3 h-3 rounded-full bg-chart-2"></span>
              <div>
                <div class="font-semibold text-foreground">API key generated</div>
                <div class="text-muted-foreground text-sm">Production key for WebApp Dashboard</div>
                <div class="text-xs text-muted-foreground">15 minutes ago</div>
              </div>
            </div>
            <div class="flex gap-3 items-start">
              <span class="mt-1 w-3 h-3 rounded-full bg-chart-1"></span>
              <div>
                <div class="font-semibold text-foreground">Rate limit exceeded</div>
                <div class="text-muted-foreground text-sm">Application 'DataProcessor' hit rate limits</div>
                <div class="text-xs text-muted-foreground">1 hour ago</div>
              </div>
            </div>
            <div class="flex gap-3 items-start">
              <span class="mt-1 w-3 h-3 rounded-full bg-chart-3"></span>
              <div>
                <div class="font-semibold text-foreground">User invitation sent</div>
                <div class="text-muted-foreground text-sm">Invited john@example.com as Developer</div>
                <div class="text-xs text-muted-foreground">2 hours ago</div>
              </div>
            </div>
          </div>
        </UiCardContent>
        <UiCardFooter>
          <UiButton variant="outline" class="w-full">View All Activity</UiButton>
        </UiCardFooter>
      </UiCard>
      <!-- Quick Actions -->
      <UiCard>
        <UiCardHeader>
          <div class="flex items-center gap-2">
            <Zap class="h-5 w-5 text-primary" />
            <UiCardTitle>Quick Actions</UiCardTitle>
          </div>
        </UiCardHeader>
        <UiCardContent>
          <div class="flex flex-col gap-3 mb-8">
            <UiButton
              v-for="action in quickActions"
              :key="action.label"
              :variant="action.variant"
              :class="action.class"
              @click="handleQuickAction(action.label)"
            >
              <component :is="action.icon" class="mr-2 h-4 w-4" />{{ action.label }}
            </UiButton>
          </div>
          <div class="mb-2 flex items-center justify-between">
            <span class="text-muted-foreground text-sm">Resource Usage</span>
            <UiBadge variant="secondary" class="text-xs">
              {{ resourceUsage.status }}
            </UiBadge>
          </div>
          <div>
            <div class="flex justify-between text-xs mb-1 text-foreground">
              <span>CPU</span>
              <span>{{ resourceUsage.cpu }}%</span>
            </div>
            <UiProgress v-model="resourceUsage.cpu" class="mb-3" />
            <div class="flex justify-between text-xs mb-1 text-foreground">
              <span>Memory</span>
              <span>{{ resourceUsage.memory }}%</span>
            </div>
            <UiProgress v-model="resourceUsage.memory" />
          </div>
        </UiCardContent>
      </UiCard>
    </div>
  </section>
</template>
