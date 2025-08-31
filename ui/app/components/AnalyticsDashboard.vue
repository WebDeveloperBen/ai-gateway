<script setup lang="ts">
import { Key, Layers, Activity, Users, ChartArea, Zap } from "lucide-vue-next";
import type { FunctionalComponent } from "vue";
import type { ButtonProps } from "./Ui/Button.vue";

const cards: {
  title: string;
  icon: FunctionalComponent;
  value: string;
  change: string;
  changeColor: string;
}[] = [
  {
    title: "Active Applications",
    icon: Layers,
    value: "12",
    change: "+2 this week",
    changeColor: "text-green-600",
  },
  {
    title: "API Keys",
    icon: Key,
    value: "47",
    change: "+5 this month",
    changeColor: "text-green-600",
  },
  {
    title: "Total Users",
    icon: Users,
    value: "1,234",
    change: "+12% from last month",
    changeColor: "text-green-600",
  },
  {
    title: "Requests Today",
    icon: ChartArea,
    value: "89.2K",
    change: "-2% from yesterday",
    changeColor: "text-red-500",
  },
];

const quickActions: {
  label: string;
  icon: FunctionalComponent;
  class: string;
  variant: ButtonProps["variant"];
}[] = [
  {
    label: "Register New Application",
    icon: Layers,
    class: "bg-gradient-to-r from-violet-500 to-purple-400 w-full text-white",
    variant: "default",
  },
  {
    label: "Generate API Key",
    icon: Key,
    variant: "ghost",
    class: "justify-start border",
  },
  {
    label: "Invite User",
    icon: Users,
    variant: "ghost",
    class: "justify-start border",
  },
  {
    label: "View Analytics",
    icon: ChartArea,
    variant: "ghost",
    class: "justify-start border",
  },
];
const resourceUsage = {
  status: "Normal",
  cpu: 45,
  memory: 62,
};
</script>

<template>
  <section class="flex flex-col gap-8">
    <!-- Header -->
    <div>
      <h1 class="text-3xl font-bold text-violet-600">Dashboard</h1>

      <p class="text-lg text-gray-500 mt-1">
        Overview of your LLM proxy service performance and activity
      </p>
    </div>
    <!-- Stats -->
    <div class="grid grid-cols-1 sm:grid-cols-2 lg:grid-cols-4 gap-4">
      <UiCard v-for="card in cards" :key="card.title">
        <template #header>
          <div class="flex items-center gap-2 justify-between px-5">
            <span class="font-semibold">{{ card.title }}</span>
            <component :is="card.icon" class="text-violet-500" />
          </div>
        </template>
        <template #content>
          <UiCardContent>
            <div class="text-2xl font-bold mt-2">{{ card.value }}</div>
            <div :class="card.changeColor + ' text-xs mt-1'">
              {{ card.change }}
            </div>
          </UiCardContent>
        </template>
      </UiCard>
    </div>
    <!-- Main Content -->
    <div class="grid grid-cols-1 md:grid-cols-2 gap-6">
      <!-- Recent Activity -->
      <UiCard>
        <template #header>
          <div class="flex items-center px-5 gap-2 mb-4">
            <Activity class="text-violet-500" />
            <h2 class="font-semibold text-lg">Recent Activity</h2>
          </div>
        </template>
        <template #content>
          <UiCardContent>
            <div class="flex flex-col gap-4">
              <div class="flex gap-3 items-start">
                <span class="mt-1 w-3 h-3 rounded-full bg-green-400"></span>
                <div>
                  <div class="font-semibold">New application registered</div>
                  <div class="text-gray-600 text-sm">
                    ChatBot Pro by ACME Corp
                  </div>
                  <div class="text-xs text-gray-400">2 minutes ago</div>
                </div>
              </div>
              <div class="flex gap-3 items-start">
                <span class="mt-1 w-3 h-3 rounded-full bg-green-400"></span>
                <div>
                  <div class="font-semibold">API key generated</div>
                  <div class="text-gray-600 text-sm">
                    Production key for WebApp Dashboard
                  </div>
                  <div class="text-xs text-gray-400">15 minutes ago</div>
                </div>
              </div>
              <div class="flex gap-3 items-start">
                <span class="mt-1 w-3 h-3 rounded-full bg-yellow-400"></span>
                <div>
                  <div class="font-semibold">Rate limit exceeded</div>
                  <div class="text-gray-600 text-sm">
                    Application 'DataProcessor' hit rate limits
                  </div>
                  <div class="text-xs text-gray-400">1 hour ago</div>
                </div>
              </div>
              <div class="flex gap-3 items-start">
                <span class="mt-1 w-3 h-3 rounded-full bg-sky-400"></span>
                <div>
                  <div class="font-semibold">User invitation sent</div>
                  <div class="text-gray-600 text-sm">
                    Invited john@example.com as Developer
                  </div>
                  <div class="text-xs text-gray-400">2 hours ago</div>
                </div>
              </div>
            </div>
          </UiCardContent>
        </template>
        <template #footer>
          <UiCardFooter class="flex mx-6 justify-between mt-6 p-0">
            <UiButton variant="outline" class="w-full"
              >View All Activity</UiButton
            >
          </UiCardFooter>
        </template>
      </UiCard>
      <!-- Quick Actions -->
      <UiCard>
        <template #header>
          <div class="flex items-center gap-2 px-5 mb-4">
            <component :is="Zap" class="text-violet-500" />
            <h2 class="font-semibold text-lg">Quick Actions</h2>
          </div>
        </template>
        <template #content>
          <UiCardContent>
            <div class="flex flex-col gap-3 mb-8">
              <UiButton
                v-for="action in quickActions"
                :key="action.label"
                :variant="action.variant"
                :class="action.class"
              >
                <component :is="action.icon" class="mr-2" />{{ action.label }}
              </UiButton>
            </div>
            <div class="mb-2 flex items-center justify-between">
              <span class="text-gray-600 text-sm">Resource Usage</span>
              <span class="text-xs text-gray-500 bg-gray-100 px-2 rounded">
                {{ resourceUsage.status }}
              </span>
            </div>
            <div>
              <div class="flex justify-between text-xs mb-1">
                <span>CPU</span>
                <span>{{ resourceUsage.cpu }}%</span>
              </div>
              <UiProgress
                v-model="resourceUsage.cpu"
                color="violet"
                class="mb-3"
              />
              <div class="flex justify-between text-xs mb-1">
                <span>Memory</span>
                <span>{{ resourceUsage.memory }}%</span>
              </div>
              <UiProgress v-model="resourceUsage.memory" color="sky" />
            </div>
          </UiCardContent>
        </template>
      </UiCard>
    </div>
  </section>
</template>
