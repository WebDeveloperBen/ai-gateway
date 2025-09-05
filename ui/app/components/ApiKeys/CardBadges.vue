<script setup lang="ts">
import { Circle, CheckCircle, Clock } from "lucide-vue-next"
import type { ApiKey } from "../ApiKeysList.vue"

interface Props {
  apiKey: ApiKey
}

defineProps<Props>()

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
  if (isExpired(expiryDate)) return 'expired'
  if (isExpiringSoon(expiryDate)) return 'expiring-soon'
  return 'valid'
}
</script>

<template>
  <UiCardContent>
    <div class="grid grid-cols-1 md:grid-cols-3 gap-4 text-sm">
      <div class="flex items-center gap-2">
        <Circle class="h-4 w-4 text-chart-4" />
        <div>
          <p class="text-muted-foreground">Created</p>
          <p class="font-medium">{{ new Date(apiKey.created).toLocaleDateString() }}</p>
        </div>
      </div>
      
      <div v-if="apiKey.expires" class="flex items-center gap-2">
        <Clock 
          class="h-4 w-4" 
          :class="{
            'text-red-500': getExpiryStatus(apiKey.expires) === 'expired',
            'text-amber-500': getExpiryStatus(apiKey.expires) === 'expiring-soon',
            'text-chart-3': getExpiryStatus(apiKey.expires) === 'valid'
          }"
        />
        <div>
          <p class="text-muted-foreground">Expires</p>
          <p 
            class="font-medium" 
            :class="{
              'text-red-500': getExpiryStatus(apiKey.expires) === 'expired',
              'text-amber-500': getExpiryStatus(apiKey.expires) === 'expiring-soon'
            }"
          >
            {{ new Date(apiKey.expires).toLocaleDateString() }}
          </p>
        </div>
      </div>

      <div class="flex items-center gap-2">
        <CheckCircle class="h-4 w-4 text-primary" />
        <div>
          <p class="text-muted-foreground">Status</p>
          <div class="flex items-center gap-2">
            <p class="font-medium capitalize">{{ apiKey.status }}</p>
            <UiBadge 
              v-if="apiKey.expires && getExpiryStatus(apiKey.expires) === 'expired'"
              variant="destructive"
              size="sm"
            >
              Expired
            </UiBadge>
            <UiBadge 
              v-else-if="apiKey.expires && getExpiryStatus(apiKey.expires) === 'expiring-soon'"
              variant="secondary"
              size="sm"
              class="bg-amber-100 text-amber-800 border-amber-200"
            >
              Expiring Soon
            </UiBadge>
          </div>
        </div>
      </div>
    </div>
  </UiCardContent>
</template>
