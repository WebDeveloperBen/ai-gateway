<script setup lang="ts">
interface Metadata {
  runId: string
  timestamp: string
  baseFilename: string
  comparisonFilename: string
  baseVersion: string
  comparisonVersion: string
  summary: string
}

interface Props {
  metadata: Metadata
}

interface Emits {
  (event: "copy-run-id"): void
}

const props = defineProps<Props>()
const emit = defineEmits<Emits>()

const formattedTimestamp = computed(() => {
  return new Date(props.metadata.timestamp).toLocaleString()
})

const handleCopyRunId = () => {
  emit("copy-run-id")
}
</script>

<template>
  <header class="bg-white border-b border-gray-200 px-6 py-4">
    <div class="flex items-start justify-between">
      <div class="flex-1">
        <div class="flex items-center gap-4 mb-2">
          <h1 class="text-xl font-semibold text-gray-900">Policy Diff Comparison</h1>
          <UiBadge variant="secondary" class="text-xs">
            {{ formattedTimestamp }}
          </UiBadge>
        </div>

        <div class="flex items-center gap-6 text-sm text-gray-600 mb-3">
          <div class="flex items-center gap-2">
            <span class="font-medium">Base:</span>
            <span>{{ metadata.baseFilename }}</span>
            <UiBadge variant="outline" size="sm">v{{ metadata.baseVersion }}</UiBadge>
          </div>

          <div class="flex items-center gap-2">
            <span class="font-medium">Comparison:</span>
            <span>{{ metadata.comparisonFilename }}</span>
            <UiBadge variant="outline" size="sm">v{{ metadata.comparisonVersion }}</UiBadge>
          </div>
        </div>

        <p class="text-sm text-gray-600 max-w-2xl">
          {{ metadata.summary }}
        </p>
      </div>

      <div class="flex items-center gap-3">
        <div class="text-right">
          <div class="text-xs text-gray-500 mb-1">Run ID</div>
          <div class="font-mono text-sm text-gray-700">{{ metadata.runId }}</div>
        </div>
        <ButtonsCopy :text="metadata.runId" @copied="handleCopyRunId" size="sm" />
      </div>
    </div>
  </header>
</template>

