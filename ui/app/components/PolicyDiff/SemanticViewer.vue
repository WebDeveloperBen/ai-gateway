<script setup lang="ts">
interface DiffSegment {
  type: 'unchanged' | 'added' | 'removed' | 'modified'
  baseText?: string
  comparisonText?: string
  context: string
}

interface PolicyPair {
  id: string
  type: 'aligned' | 'unmatched_base' | 'unmatched_comparison'
  title: string
  section: string
  baseContent?: string
  comparisonContent?: string
  diffSegments?: DiffSegment[]
  filename: string
  baseVersion?: string
  comparisonVersion?: string
  sourceTrace: string
  similarities: string[]
  keyDifferences: string[]
}

interface Props {
  pair: PolicyPair
}

const props = defineProps<Props>()

const renderHighlightedText = (segments: DiffSegment[], isComparison = false) => {
  if (!segments) return []

  return segments.map((segment, index) => ({
    key: index,
    text: isComparison ? segment.comparisonText : segment.baseText,
    type: segment.type,
    context: segment.context
  })).filter(item => item.text)
}
</script>

<template>
  <div class="h-full flex flex-col bg-white">
    <!-- Header -->
    <div class="border-b border-gray-200 px-6 py-4">
      <div class="flex items-center justify-between mb-3">
        <div class="flex-1">
          <h3 class="text-lg font-semibold text-gray-900">{{ props.pair.title }}</h3>
          <p class="text-sm text-gray-600 mt-1">{{ props.pair.section }}</p>
        </div>

        <div v-if="props.pair.type === 'aligned'" class="flex items-center gap-4 text-sm text-gray-600">
          <div class="flex items-center gap-2">
            <span class="font-medium">Base:</span>
            <UiBadge variant="outline" size="sm">v{{ props.pair.baseVersion }}</UiBadge>
          </div>
          <div class="flex items-center gap-2">
            <span class="font-medium">Comparison:</span>
            <UiBadge variant="outline" size="sm">v{{ props.pair.comparisonVersion }}</UiBadge>
          </div>
        </div>
        <div v-else class="text-sm text-gray-600">
          <UiBadge variant="outline" size="sm" v-if="props.pair.baseVersion">v{{ props.pair.baseVersion }}</UiBadge>
          <UiBadge variant="outline" size="sm" v-if="props.pair.comparisonVersion">v{{ props.pair.comparisonVersion }}</UiBadge>
        </div>
      </div>

      <div class="text-sm text-gray-500">
        <span class="font-medium">Source:</span> {{ props.pair.sourceTrace }}
      </div>
    </div>

    <div class="flex-1 flex flex-col overflow-hidden">
      <!-- Policy Text Comparison -->
      <div v-if="props.pair.type === 'aligned'" class="flex-1 flex overflow-hidden">
        <!-- Base Version -->
        <div class="flex-1 border-r border-gray-200 flex flex-col">
          <div class="px-4 py-3 bg-red-50 border-b border-gray-200 text-sm font-medium text-red-800">
            Base Version (v{{ props.pair.baseVersion }})
          </div>
          <div class="flex-1 p-6 overflow-y-auto">
            <div class="prose prose-sm max-w-none">
              <template v-if="props.pair.diffSegments">
                <span
                  v-for="segment in renderHighlightedText(props.pair.diffSegments, false)"
                  :key="segment.key"
                  :class="{
                    'bg-red-100 text-red-800 px-1 rounded': segment.type === 'removed' || segment.type === 'modified',
                    'bg-gray-100 text-gray-600': segment.type === 'unchanged'
                  }"
                  class="leading-relaxed"
                >{{ segment.text }}</span>
              </template>
              <template v-else>
                <p class="text-gray-800 leading-relaxed">{{ props.pair.baseContent }}</p>
              </template>
            </div>
          </div>
        </div>

        <!-- Comparison Version -->
        <div class="flex-1 flex flex-col">
          <div class="px-4 py-3 bg-green-50 border-b border-gray-200 text-sm font-medium text-green-800">
            Comparison Version (v{{ props.pair.comparisonVersion }})
          </div>
          <div class="flex-1 p-6 overflow-y-auto">
            <div class="prose prose-sm max-w-none">
              <template v-if="props.pair.diffSegments">
                <span
                  v-for="segment in renderHighlightedText(props.pair.diffSegments, true)"
                  :key="segment.key"
                  :class="{
                    'bg-green-100 text-green-800 px-1 rounded': segment.type === 'added' || segment.type === 'modified',
                    'bg-gray-100 text-gray-600': segment.type === 'unchanged'
                  }"
                  class="leading-relaxed"
                >{{ segment.text }}</span>
              </template>
              <template v-else>
                <p class="text-gray-800 leading-relaxed">{{ props.pair.comparisonContent }}</p>
              </template>
            </div>
          </div>
        </div>
      </div>

      <!-- Single Document View -->
      <div v-else class="flex-1 flex flex-col">
        <div class="px-4 py-3 bg-gray-50 border-b border-gray-200 text-sm font-medium text-gray-800">
          {{ props.pair.type === 'unmatched_base' ? 'Removed in New Version' : 'Added in New Version' }}
          (v{{ props.pair.type === 'unmatched_base' ? props.pair.baseVersion : props.pair.comparisonVersion }})
        </div>
        <div class="flex-1 p-6 overflow-y-auto">
          <div class="prose prose-sm max-w-none">
            <p class="text-gray-800 leading-relaxed">
              {{ props.pair.type === 'unmatched_base' ? props.pair.baseContent : props.pair.comparisonContent }}
            </p>
          </div>
        </div>
      </div>

      <!-- Analysis Section -->
      <div class="border-t border-gray-200 bg-gray-50">
        <div class="grid grid-cols-1 lg:grid-cols-2 gap-6 p-6">
          <!-- Similarities -->
          <div v-if="props.pair.similarities?.length">
            <h4 class="text-sm font-medium text-gray-900 mb-3 flex items-center gap-2">
              <div class="w-2 h-2 rounded-full bg-blue-500"></div>
              Similarities
            </h4>
            <ul class="space-y-2">
              <li
                v-for="(similarity, index) in props.pair.similarities"
                :key="index"
                class="flex items-start gap-3 text-sm"
              >
                <div class="w-1.5 h-1.5 rounded-full bg-blue-400 mt-2 flex-shrink-0"></div>
                <span class="text-gray-700">{{ similarity }}</span>
              </li>
            </ul>
          </div>

          <!-- Key Differences -->
          <div v-if="props.pair.keyDifferences?.length">
            <h4 class="text-sm font-medium text-gray-900 mb-3 flex items-center gap-2">
              <div class="w-2 h-2 rounded-full bg-orange-500"></div>
              Key Differences
            </h4>
            <ul class="space-y-2">
              <li
                v-for="(difference, index) in props.pair.keyDifferences"
                :key="index"
                class="flex items-start gap-3 text-sm"
              >
                <div class="w-1.5 h-1.5 rounded-full bg-orange-400 mt-2 flex-shrink-0"></div>
                <span class="text-gray-700">{{ difference }}</span>
              </li>
            </ul>
          </div>
        </div>
      </div>
    </div>
  </div>
</template>