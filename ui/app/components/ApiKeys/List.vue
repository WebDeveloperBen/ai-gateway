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
import { MoreVertical, Key, Activity, Circle, Clock, Eye, Trash2 } from "lucide-vue-next"

withDefaults(defineProps<KeyProps>(), {
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

    <UiCard v-for="key in keys" :key="key.id" interactive @click="handleApiKeyClick(key)">
      <UiCardHeader>
        <div class="flex items-start justify-between">
          <div class="space-y-1 flex-1">
            <div class="flex items-center gap-3">
              <Key class="h-5 w-5 text-primary" />
              <UiCardTitle class="text-lg">{{ key.name }}</UiCardTitle>
              <UiStatusBadge :status="key.status" />
            </div>
            <UiCardDescription class="text-sm">
              {{ key.applicationName }}
            </UiCardDescription>
            <div class="flex items-center gap-4 text-xs">
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
              <UiDropdownMenuItem @click="handleApiKeyClick(key)">
                <Eye class="mr-2 size-4" />
                View Details
              </UiDropdownMenuItem>
              <UiDropdownMenuSeparator />
              <UiDropdownMenuItem class="text-red-600" @click="handleDeleteKey(key)">
                <Trash2 class="mr-2 size-4" />
                Delete Key
              </UiDropdownMenuItem>
            </UiDropdownMenuContent>
          </UiDropdownMenu>
        </div>
      </UiCardHeader>
      <UiCardContent>
        <!-- API Key Display -->
        <div class="mb-4">
          <ApiKeysReveal @click.stop :key-id="key.id" :key-prefix="key.keyPrefix" size="lg" />
        </div>

        <!-- Stats Grid -->
        <div class="grid grid-cols-1 md:grid-cols-3 gap-4 text-sm">
          <div class="flex items-center gap-2 justify-center">
            <Circle class="h-4 w-4 text-chart-2" />
            <div>
              <p class="text-muted-foreground">Created</p>
              <p class="font-medium">{{ new Date(key.created).toLocaleDateString() }}</p>
            </div>
          </div>
          <div class="flex items-center gap-2 justify-center">
            <Activity class="h-4 w-4 text-chart-1" />
            <div>
              <p class="text-muted-foreground">Last Used</p>
              <p class="font-medium">{{ new Date(key.lastUsed).toLocaleDateString() }}</p>
            </div>
          </div>
          <div class="flex items-center gap-2 justify-center">
            <Activity class="h-4 w-4 text-chart-3" />
            <div>
              <p class="text-muted-foreground">Requests</p>
              <div class="flex items-center gap-2">
                <p class="font-medium">{{ formatNumber(key.requestCount) }}</p>
                <UiActivityIndicator :value="key.requestCount" />
              </div>
            </div>
          </div>
        </div>
      </UiCardContent>
    </UiCard>
  </div>
</template>
