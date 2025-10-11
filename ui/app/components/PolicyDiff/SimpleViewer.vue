<script setup lang="ts">
interface PolicyPair {
  id: string
  type: 'aligned' | 'unmatched_base' | 'unmatched_comparison'
  baseContent?: string
  comparisonContent?: string
  filename: string
  baseVersion?: string
  comparisonVersion?: string
  sourceTrace: string
  deltaDescription: string[]
}

interface Props {
  pair: PolicyPair
}

const props = defineProps<Props>()

const getSingleEditorContent = () => {
  if (props.pair.type === 'unmatched_base') {
    return props.pair.baseContent || ''
  }
  return props.pair.comparisonContent || ''
}

const formatJson = (content: string) => {
  try {
    return JSON.stringify(JSON.parse(content), null, 2)
  } catch {
    return content
  }
}
</script>

<template>
  <div class="h-full flex flex-col bg-white">
    <div class="border-b border-gray-200 px-6 py-4">
      <div class="flex items-center justify-between mb-3">
        <div class="flex items-center gap-3">
          <h3 class="text-lg font-semibold text-gray-900">{{ props.pair.filename }}</h3>
          <UiBadge
            v-if="props.pair.type === 'aligned'"
            variant="secondary"
            size="sm"
          >
            Side-by-side comparison
          </UiBadge>
          <UiBadge
            v-else-if="props.pair.type === 'unmatched_base'"
            variant="destructive"
            size="sm"
          >
            Base Version
          </UiBadge>
          <UiBadge
            v-else-if="props.pair.type === 'unmatched_comparison'"
            variant="default"
            size="sm"
          >
            Comparison Version
          </UiBadge>
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

      <div class="text-sm text-gray-500 mb-2">
        <span class="font-medium">Source:</span> {{ props.pair.sourceTrace }}
      </div>
    </div>

    <div class="flex-1 flex flex-col">
      <!-- Aligned comparison - side by side -->
      <div v-if="props.pair.type === 'aligned'" class="flex-1 flex">
        <div class="flex-1 border-r border-gray-200">
          <div class="px-4 py-2 bg-red-50 border-b border-gray-200 text-sm font-medium text-red-800">
            Base (v{{ props.pair.baseVersion }})
          </div>
          <div class="p-4 overflow-auto h-full">
            <pre class="text-sm text-gray-800 whitespace-pre-wrap font-mono">{{ formatJson(props.pair.baseContent || '') }}</pre>
          </div>
        </div>

        <div class="flex-1">
          <div class="px-4 py-2 bg-green-50 border-b border-gray-200 text-sm font-medium text-green-800">
            Comparison (v{{ props.pair.comparisonVersion }})
          </div>
          <div class="p-4 overflow-auto h-full">
            <pre class="text-sm text-gray-800 whitespace-pre-wrap font-mono">{{ formatJson(props.pair.comparisonContent || '') }}</pre>
          </div>
        </div>
      </div>

      <!-- Single view for unmatched items -->
      <div v-else class="flex-1">
        <div class="px-4 py-2 bg-gray-50 border-b border-gray-200 text-sm font-medium text-gray-800">
          {{ props.pair.type === 'unmatched_base' ? 'Base Version' : 'Comparison Version' }}
          (v{{ props.pair.type === 'unmatched_base' ? props.pair.baseVersion : props.pair.comparisonVersion }})
        </div>
        <div class="p-4 overflow-auto h-full">
          <pre class="text-sm text-gray-800 whitespace-pre-wrap font-mono">{{ formatJson(getSingleEditorContent()) }}</pre>
        </div>
      </div>

      <div class="border-t border-gray-200 bg-gray-50 px-6 py-4">
        <h4 class="text-sm font-medium text-gray-900 mb-3">Changes Summary</h4>
        <ul class="space-y-2">
          <li
            v-for="(delta, index) in props.pair.deltaDescription"
            :key="index"
            class="flex items-start gap-3 text-sm"
          >
            <div class="w-2 h-2 rounded-full bg-blue-500 mt-2 flex-shrink-0"></div>
            <span class="text-gray-700">{{ delta }}</span>
          </li>
        </ul>
      </div>
    </div>
  </div>
</template>