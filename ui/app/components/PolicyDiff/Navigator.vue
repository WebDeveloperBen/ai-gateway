<script setup lang="ts">
import { GitCompare, MinusCircle, PlusCircle, ArrowRight } from "lucide-vue-next"

interface PolicyPair {
  id: string
  type: 'aligned' | 'unmatched_base' | 'unmatched_comparison'
  title: string
  section: string
  filename: string
  baseVersion?: string
  comparisonVersion?: string
  similarities: string[]
  keyDifferences: string[]
}

interface GroupedPairs {
  aligned: PolicyPair[]
  unmatched_base: PolicyPair[]
  unmatched_comparison: PolicyPair[]
}

interface Props {
  groupedPairs: GroupedPairs
  selectedPairId: string | null
}

interface Emits {
  (event: 'select-pair', pairId: string): void
}

defineProps<Props>()
const emit = defineEmits<Emits>()

const handleSelectPair = (pairId: string) => {
  emit('select-pair', pairId)
}

const getStatusIcon = (type: PolicyPair['type']) => {
  switch (type) {
    case 'aligned':
      return GitCompare
    case 'unmatched_base':
      return MinusCircle
    case 'unmatched_comparison':
      return PlusCircle
    default:
      return GitCompare
  }
}

const getStatusColor = (type: PolicyPair['type']) => {
  switch (type) {
    case 'aligned':
      return 'text-blue-600'
    case 'unmatched_base':
      return 'text-red-600'
    case 'unmatched_comparison':
      return 'text-green-600'
    default:
      return 'text-gray-600'
  }
}
</script>

<template>
  <nav class="h-full bg-white overflow-y-auto">
    <div class="p-4">
      <h2 class="text-lg font-semibold text-gray-900 mb-4">Policy Changes</h2>

      <div class="space-y-6">
        <div v-if="groupedPairs.aligned.length > 0">
          <div class="flex items-center justify-between mb-3">
            <h3 class="text-sm font-medium text-gray-700">Aligned Changes</h3>
            <UiBadge variant="secondary" size="sm">
              {{ groupedPairs.aligned.length }}
            </UiBadge>
          </div>

          <div class="space-y-2">
            <button
              v-for="pair in groupedPairs.aligned"
              :key="pair.id"
              @click="handleSelectPair(pair.id)"
              :class="[
                'w-full text-left p-3 rounded-lg border transition-colors',
                selectedPairId === pair.id
                  ? 'border-blue-500 bg-blue-50'
                  : 'border-gray-200 hover:border-gray-300 hover:bg-gray-50'
              ]"
            >
              <div class="flex items-start gap-3">
                <component
                  :is="getStatusIcon(pair.type)"
                  :class="['w-4 h-4 mt-0.5', getStatusColor(pair.type)]"
                />
                <div class="flex-1 min-w-0">
                  <div class="font-medium text-sm text-gray-900 truncate">
                    {{ pair.title }}
                  </div>
                  <div class="text-xs text-gray-500 mb-2 truncate">
                    {{ pair.section }}
                  </div>
                  <div class="flex items-center gap-2">
                    <UiBadge variant="outline" size="sm" v-if="pair.baseVersion">
                      {{ pair.baseVersion }}
                    </UiBadge>
                    <ArrowRight class="w-3 h-3 text-gray-400" />
                    <UiBadge variant="outline" size="sm" v-if="pair.comparisonVersion">
                      {{ pair.comparisonVersion }}
                    </UiBadge>
                  </div>
                  <div class="text-xs text-gray-500 mt-1">
                    {{ pair.keyDifferences.length }} key difference{{ pair.keyDifferences.length !== 1 ? 's' : '' }}
                  </div>
                </div>
              </div>
            </button>
          </div>
        </div>

        <div v-if="groupedPairs.unmatched_base.length > 0">
          <div class="flex items-center justify-between mb-3">
            <h3 class="text-sm font-medium text-gray-700">Removed in Comparison</h3>
            <UiBadge variant="destructive" size="sm">
              {{ groupedPairs.unmatched_base.length }}
            </UiBadge>
          </div>

          <div class="space-y-2">
            <button
              v-for="pair in groupedPairs.unmatched_base"
              :key="pair.id"
              @click="handleSelectPair(pair.id)"
              :class="[
                'w-full text-left p-3 rounded-lg border transition-colors',
                selectedPairId === pair.id
                  ? 'border-red-500 bg-red-50'
                  : 'border-gray-200 hover:border-gray-300 hover:bg-gray-50'
              ]"
            >
              <div class="flex items-start gap-3">
                <component
                  :is="getStatusIcon(pair.type)"
                  :class="['w-4 h-4 mt-0.5', getStatusColor(pair.type)]"
                />
                <div class="flex-1 min-w-0">
                  <div class="font-medium text-sm text-gray-900 truncate">
                    {{ pair.title }}
                  </div>
                  <div class="text-xs text-gray-500 mb-2 truncate">
                    {{ pair.section }}
                  </div>
                  <UiBadge variant="outline" size="sm" v-if="pair.baseVersion">
                    {{ pair.baseVersion }}
                  </UiBadge>
                  <div class="text-xs text-gray-500 mt-1">
                    Only in base version
                  </div>
                </div>
              </div>
            </button>
          </div>
        </div>

        <div v-if="groupedPairs.unmatched_comparison.length > 0">
          <div class="flex items-center justify-between mb-3">
            <h3 class="text-sm font-medium text-gray-700">Added in Comparison</h3>
            <UiBadge variant="default" size="sm">
              {{ groupedPairs.unmatched_comparison.length }}
            </UiBadge>
          </div>

          <div class="space-y-2">
            <button
              v-for="pair in groupedPairs.unmatched_comparison"
              :key="pair.id"
              @click="handleSelectPair(pair.id)"
              :class="[
                'w-full text-left p-3 rounded-lg border transition-colors',
                selectedPairId === pair.id
                  ? 'border-green-500 bg-green-50'
                  : 'border-gray-200 hover:border-gray-300 hover:bg-gray-50'
              ]"
            >
              <div class="flex items-start gap-3">
                <component
                  :is="getStatusIcon(pair.type)"
                  :class="['w-4 h-4 mt-0.5', getStatusColor(pair.type)]"
                />
                <div class="flex-1 min-w-0">
                  <div class="font-medium text-sm text-gray-900 truncate">
                    {{ pair.title }}
                  </div>
                  <div class="text-xs text-gray-500 mb-2 truncate">
                    {{ pair.section }}
                  </div>
                  <UiBadge variant="outline" size="sm" v-if="pair.comparisonVersion">
                    {{ pair.comparisonVersion }}
                  </UiBadge>
                  <div class="text-xs text-gray-500 mt-1">
                    Only in comparison version
                  </div>
                </div>
              </div>
            </button>
          </div>
        </div>
      </div>
    </div>
  </nav>
</template>