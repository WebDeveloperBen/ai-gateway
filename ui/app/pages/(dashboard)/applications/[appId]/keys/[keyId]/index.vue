<script setup lang="ts">
import { Key, RotateCw, Trash2, Calendar, Clock, MoreVertical, Users, Ban, TestTubeDiagonal } from "lucide-vue-next"

// Get the app ID and key ID from route params
const route = useRoute()
const appId = route.params.appId as string
const keyId = route.params.keyId as string

// Modal state management
const showDisableModal = ref(false)
const showDeleteModal = ref(false)
const actionLoading = ref(false)

// Actions
const handleDisableKey = async () => {
  actionLoading.value = true
  try {
    // TODO: Implement API call to disable key
    await new Promise((resolve) => setTimeout(resolve, 1000)) // Simulate API call
    apiKey.value.status = "inactive"
    showDisableModal.value = false
  } catch (error) {
    console.error("Failed to disable key:", error)
  } finally {
    actionLoading.value = false
  }
}

const handleDeleteKey = async () => {
  actionLoading.value = true
  try {
    // TODO: Implement API call to delete key
    await new Promise((resolve) => setTimeout(resolve, 1000)) // Simulate API call
    // In real app, this would redirect to keys list after deletion
    showDeleteModal.value = false
    await navigateTo(`/applications/${appId}/keys`)
  } catch (error) {
    console.error("Failed to delete key:", error)
  } finally {
    actionLoading.value = false
  }
}

// Sample API key data - this would come from an API call
const apiKey = ref({
  id: keyId,
  key: "",
  name: "Prwwoduction API Key",
  keyPrefix: "sk-",
  applicationId: appId,
  applicationName: "Customer Service Bot",
  description: "Main production key for customer service bot integration",
  status: "active",
  permissions: ["read", "write"],
  created: "2024-12-15T10:00:00Z",
  lastUsed: "2025-01-15T14:30:00Z",
  expiresAt: null,
  owners: [
    {
      id: "user_1",
      name: "John Smith",
      email: "john.smith@company.com",
      role: "Admin"
    },
    {
      id: "user_2",
      name: "Sarah Johnson",
      email: "sarah.johnson@company.com",
      role: "Developer"
    }
  ],
  requestCount: {
    total: 25400,
    today: 240,
    thisWeek: 1680,
    thisMonth: 7200
  },
  errorRate: 2.1,
  recentActivity: [
    {
      timestamp: "2025-01-15T14:30:00Z",
      method: "POST",
      endpoint: "/v1/chat/completions",
      status: 200,
      ip: "192.168.1.100"
    },
    {
      timestamp: "2025-01-15T14:25:00Z",
      method: "POST",
      endpoint: "/v1/chat/completions",
      status: 200,
      ip: "192.168.1.100"
    },
    {
      timestamp: "2025-01-15T14:20:00Z",
      method: "GET",
      endpoint: "/v1/models",
      status: 200,
      ip: "192.168.1.100"
    }
  ]
})
</script>

<template>
  <div class="flex flex-col gap-6">
    <PageBadgeHeader :subtext="`Used by ${apiKey.applicationName}`" :badge-status="apiKey.status" :title="apiKey.name">
      <!-- Action Buttons -->
      <ButtonsGroup>
        <ButtonsGoTo :action="() => navigateTo('/playgrouns/prompts')" title="Go Playground" :icon="TestTubeDiagonal" />
        <ButtonsDisable :action="() => (showDisableModal = true)" title="Key" />
        <ButtonsDelete :action="() => (showDisableModal = true)" />
      </ButtonsGroup>
    </PageBadgeHeader>
    <!-- Key Information Card -->
    <UiCard>
      <UiCardHeader>
        <UiCardTitle class="flex items-center px-0 gap-2">
          <Key class="h-5 w-5" />
          API Key Information
        </UiCardTitle>
      </UiCardHeader>
      <UiCardContent class="space-y-6">
        <!-- API Key Display -->
        <div>
          <label class="text-sm font-medium text-muted-foreground">API Secret Key</label>
          <div class="mt-2">
            <ApiKeysReveal :key-id="apiKey.id" :key-prefix="apiKey.keyPrefix" size="lg" show-copy-text />
          </div>
        </div>

        <!-- Key Details Grid -->
        <div class="grid grid-cols-1 md:grid-cols-3 gap-6">
          <div>
            <label class="text-sm font-medium text-muted-foreground flex items-center gap-2">
              <Calendar class="h-4 w-4" />
              Created
            </label>
            <p class="text-sm font-medium mt-1">{{ new Date(apiKey.created).toLocaleDateString() }}</p>
          </div>

          <div>
            <label class="text-sm font-medium text-muted-foreground flex items-center gap-2">
              <Clock class="h-4 w-4" />
              Last Used
            </label>
            <p class="text-sm font-medium mt-1">{{ new Date(apiKey.lastUsed).toLocaleDateString() }}</p>
          </div>

          <div>
            <label class="text-sm font-medium text-muted-foreground flex items-center gap-2">
              <Users class="h-4 w-4" />
              Owners
            </label>
            <div class="mt-1 flex flex-wrap gap-1">
              <div v-for="owner in apiKey.owners" :key="owner.id" class="flex items-center gap-1 text-xs">
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
        <div v-if="apiKey.description">
          <label class="text-sm font-medium text-muted-foreground">Description</label>
          <p class="text-sm mt-1">{{ apiKey.description }}</p>
        </div>
      </UiCardContent>
    </UiCard>

    <CardsCodeIntegrations />
    <CardsAnalyticsComingSoon />

    <LazyConfirmationModal
      v-model:open="showDisableModal"
      title="Disable API Key?"
      :description="`Are you sure you want to disable '${apiKey.name}'? This will prevent it from being used for API requests.`"
      confirm-text="Disable Key"
      variant="default"
      :loading="actionLoading"
      @confirm="handleDisableKey"
    />

    <LazyConfirmationModal
      v-model:open="showDeleteModal"
      title="Delete API Key?"
      :description="`Are you sure you want to permanently delete '${apiKey.name}'? This action cannot be undone and will break any applications currently using this key.`"
      confirm-text="Delete Key"
      variant="destructive"
      :loading="actionLoading"
      @confirm="handleDeleteKey"
    />
  </div>
</template>
