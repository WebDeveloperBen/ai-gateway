<script lang="ts">
export type ApiKey = {
  id: string
  name: string
  keyPrefix: string
  applicationId: string
  applicationName: string
  created: string
  lastUsed: string
  expires?: string
  status: string
  permissions?: string[]
  requestCount: number
  rateLimit: string
}

export type KeyProps = {
  keys: ApiKey[]
  showEmpty?: boolean
  emptyTitle?: string
  emptyDescription?: string
  showDropdownActions?: boolean
}
</script>

<script setup lang="ts">
import { MoreVertical, Key, Activity, Circle, CheckCircle, Clock } from "lucide-vue-next"

const props = withDefaults(defineProps<KeyProps>(), {
  showEmpty: true,
  emptyTitle: "No API keys found",
  emptyDescription: "Try adjusting your search or create a new API key.",
  showDropdownActions: true
})

const emit = defineEmits<{
  selectApiKey: [apiKey: ApiKey]
  regenerateKey: [apiKey: ApiKey]
  deleteKey: [apiKey: ApiKey]
}>()

const formatNumber = (num: number) => {
  return new Intl.NumberFormat().format(num)
}

const getKeyAccent = (status: string) => {
  return status === "active" ? "chart-2" : "none"
}

const handleApiKeyClick = (key: ApiKey) => {
  emit("selectApiKey", key)
}

const handleRegenerateKey = (key: ApiKey) => {
  emit("regenerateKey", key)
}

const handleDeleteKey = (key: ApiKey) => {
  emit("deleteKey", key)
}

const isExpiringSoon = (expiryDate: string) => {
  const expiry = new Date(expiryDate)
  const now = new Date()
  const daysUntilExpiry = Math.ceil((expiry.getTime() - now.getTime()) / (1000 * 60 * 60 * 24))
  return daysUntilExpiry <= 30 && daysUntilExpiry > 0
}

const isExpired = (expiryDate: string) => {
  const expiry = new Date(expiryDate)
  const now = new Date()
  return expiry < now
}

const getExpiryStatus = (expiryDate: string) => {
  if (isExpired(expiryDate)) return "expired"
  if (isExpiringSoon(expiryDate)) return "expiring-soon"
  return "valid"
}
</script>

<template>
  <div class="flex flex-col gap-4">
    <div v-if="keys.length === 0 && showEmpty" class="text-center py-12">
      <Key class="mx-auto h-12 w-12 text-muted-foreground" />
      <h3 class="mt-4 text-lg font-medium">{{ emptyTitle }}</h3>
      <p class="text-muted-foreground">{{ emptyDescription }}</p>
    </div>

    <UiCard
      v-for="key in keys"
      :key="key.id"
      interactive
      padding="compact"
      :accent="getKeyAccent(key.status)"
      class="transition-all duration-200 hover:shadow-md"
      @click="handleApiKeyClick(key)"
    >
      <div class="p-6">
        <div class="flex items-start justify-between mb-4">
          <div class="flex-1">
            <div class="flex items-center gap-3 mb-2">
              <div class="p-2 rounded-lg bg-primary/10">
                <Key class="h-4 w-4 text-primary" />
              </div>
              <div>
                <h3 class="font-semibold text-lg text-foreground">{{ key.name }}</h3>
                <p class="text-sm text-muted-foreground">{{ key.applicationName }}</p>
              </div>
              <UiStatusBadge :status="key.status" />
            </div>

            <!-- Metadata row -->
            <div class="flex items-center gap-4 text-xs text-muted-foreground">
              <template v-if="key.expires">
                <div class="flex items-center gap-1">
                  <Clock
                    class="h-3 w-3"
                    :class="{
                      'text-red-500': key.expires && getExpiryStatus(key.expires) === 'expired',
                      'text-amber-500': key.expires && getExpiryStatus(key.expires) === 'expiring-soon',
                      'text-muted-foreground': key.expires && getExpiryStatus(key.expires) === 'valid'
                    }"
                  />
                  <span
                    class="font-medium"
                    :class="{
                      'text-red-500': key.expires && getExpiryStatus(key.expires) === 'expired',
                      'text-amber-500': key.expires && getExpiryStatus(key.expires) === 'expiring-soon',
                      'text-foreground': key.expires && getExpiryStatus(key.expires) === 'valid'
                    }"
                  >
                    Expires {{ new Date(key.expires).toLocaleDateString() }}
                  </span>
                </div>
              </template>
            </div>
          </div>

          <UiDropdownMenu v-if="showDropdownActions">
            <UiDropdownMenuTrigger as-child>
              <UiButton variant="ghost" size="sm" @click.stop>
                <MoreVertical class="h-4 w-4" />
              </UiButton>
            </UiDropdownMenuTrigger>
            <UiDropdownMenuContent align="end">
              <UiDropdownMenuItem @click="handleApiKeyClick(key)"> View Details </UiDropdownMenuItem>
              <UiDropdownMenuItem @click="handleRegenerateKey(key)"> Regenerate Key </UiDropdownMenuItem>
              <UiDropdownMenuItem> Copy Key </UiDropdownMenuItem>
              <UiDropdownMenuSeparator />
              <UiDropdownMenuItem class="text-red-600" @click="handleDeleteKey(key)"> Delete Key </UiDropdownMenuItem>
            </UiDropdownMenuContent>
          </UiDropdownMenu>
        </div>

        <!-- API Key Display -->
        <ApiKeysReveal @click.stop :key-id="key.id" :key-prefix="key.keyPrefix" size="sm" />

        <!-- Stats Grid -->
        <div class="grid grid-cols-2 md:grid-cols-3 gap-4 pt-4 border-t">
          <div class="text-center">
            <div class="flex items-center justify-center gap-1 mb-1">
              <Circle class="h-3 w-3 text-muted-foreground" />
              <span class="text-xs text-muted-foreground">Created</span>
            </div>
            <p class="text-sm font-medium">{{ new Date(key.created).toLocaleDateString() }}</p>
          </div>

          <div class="text-center">
            <div class="flex items-center justify-center gap-1 mb-1">
              <Activity class="h-3 w-3 text-muted-foreground" />
              <span class="text-xs text-muted-foreground">Last Used</span>
            </div>
            <p class="text-sm font-medium">{{ new Date(key.lastUsed).toLocaleDateString() }}</p>
          </div>

          <div class="text-center">
            <div class="flex items-center justify-center gap-1 mb-1">
              <Activity class="h-3 w-3 text-blue-500" />
              <span class="text-xs text-muted-foreground">Requests</span>
            </div>
            <div class="flex items-center justify-center gap-1">
              <p class="text-sm font-medium">{{ formatNumber(key.requestCount) }}</p>
              <UiActivityIndicator :value="key.requestCount" />
            </div>
          </div>
        </div>
      </div>
    </UiCard>
  </div>
</template>

