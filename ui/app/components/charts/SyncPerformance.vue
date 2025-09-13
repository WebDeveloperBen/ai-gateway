<script setup lang="ts">
import { TrendingUp, TrendingDown, BarChart3 } from "lucide-vue-next"

interface SyncDataPoint {
  date: string
  duration: number
  documents: number
  success: boolean
}

const props = defineProps<{
  data: SyncDataPoint[]
  title?: string
}>()

const maxDuration = computed(() => Math.max(...props.data.map(d => d.duration)))
const maxDocuments = computed(() => Math.max(...props.data.map(d => d.documents)))
const avgDuration = computed(() => props.data.reduce((sum, d) => sum + d.duration, 0) / props.data.length)
const successRate = computed(() => (props.data.filter(d => d.success).length / props.data.length) * 100)

// Calculate trend
const trend = computed(() => {
  if (props.data.length < 2) return 0
  const recent = props.data.slice(-2)
  return recent[1].duration - recent[0].duration
})

const isImproving = computed(() => trend.value < 0)
</script>

<template>
  <div class="space-y-4">
    <!-- Trend Indicator -->
    <div class="flex items-center justify-between">
      <div class="flex items-center gap-2">
        <component :is="isImproving ? TrendingDown : TrendingUp" 
                  :class="isImproving ? 'text-green-600' : 'text-red-600'" 
                  class="size-4" />
        <span class="text-sm font-medium" 
              :class="isImproving ? 'text-green-600' : 'text-red-600'">
          {{ isImproving ? 'Performance Improving' : 'Performance Slower' }}
        </span>
      </div>
      <div class="text-sm text-muted-foreground">{{ successRate.toFixed(1) }}% success rate</div>
    </div>

    <!-- Chart Visualization -->
    <div class="h-24 relative border rounded-lg p-2">
      <div class="absolute inset-2 flex items-end justify-between gap-1">
        <div v-for="(point, index) in data" :key="index" 
             class="flex flex-col items-center gap-1 flex-1">
          <!-- Duration Bar -->
          <div class="w-full bg-muted rounded-sm overflow-hidden relative" style="height: 60px;">
            <div class="absolute bottom-0 left-0 w-full rounded-sm transition-all duration-300"
                 :class="point.success ? 'bg-emerald-500' : 'bg-red-500'"
                 :style="{ height: `${(point.duration / maxDuration) * 100}%` }">
            </div>
            <!-- Tooltip on hover -->
            <div class="absolute inset-0 opacity-0 hover:opacity-100 bg-black/10 rounded-sm flex items-center justify-center transition-opacity">
              <div class="text-xs font-bold text-white bg-black/60 rounded px-1 py-0.5">
                {{ point.duration }}s
              </div>
            </div>
          </div>
          
          <!-- Date Label -->
          <div class="text-xs text-muted-foreground">{{ point.date.slice(-2) }}</div>
        </div>
      </div>
    </div>

    <!-- Mini Legend -->
    <div class="flex items-center justify-between text-xs text-muted-foreground">
      <div class="flex items-center gap-4">
        <div class="flex items-center gap-1">
          <div class="w-2 h-2 bg-emerald-500 rounded-sm"></div>
          <span>Success</span>
        </div>
        <div class="flex items-center gap-1">
          <div class="w-2 h-2 bg-red-500 rounded-sm"></div>
          <span>Failed</span>
        </div>
      </div>
      <div>Avg: {{ Math.round(avgDuration) }}s</div>
    </div>
  </div>
</template>