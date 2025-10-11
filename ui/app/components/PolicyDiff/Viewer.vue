<script setup lang="ts">
import type * as Monaco from "monaco-editor"
import type { MonacoEditor, MonacoDiffEditor } from "#components"

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

const diffEditorOptions: Monaco.editor.IStandaloneDiffEditorConstructionOptions = {
  readOnly: true,
  enableSplitViewResizing: true,
  renderSideBySide: true,
  originalEditable: false,
  theme: 'vs',
  fontSize: 14,
  lineNumbers: 'on',
  minimap: { enabled: false },
  scrollBeyondLastLine: false,
  wordWrap: 'on',
  automaticLayout: true
}

const singleEditorOptions: Monaco.editor.IStandaloneEditorConstructionOptions = {
  readOnly: true,
  theme: 'vs',
  fontSize: 14,
  lineNumbers: 'on',
  minimap: { enabled: false },
  scrollBeyondLastLine: false,
  wordWrap: 'on',
  automaticLayout: true
}

const handleDiffEditorLoad = (editor: Monaco.editor.IStandaloneDiffEditor) => {
  console.log('Diff editor loaded for pair:', props.pair.id)
}

const handleSingleEditorLoad = (editor: Monaco.editor.IStandaloneCodeEditor) => {
  console.log('Single editor loaded for pair:', props.pair.id)
}

const getSingleEditorContent = () => {
  if (props.pair.type === 'unmatched_base') {
    return props.pair.baseContent || ''
  }
  return props.pair.comparisonContent || ''
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
      <div v-if="props.pair.type === 'aligned'" class="flex-1">
        <ClientOnly>
          <MonacoDiffEditor
            ref="diffEditorRef"
            :original="props.pair.baseContent || ''"
            :value="props.pair.comparisonContent || ''"
            :lang="'json'"
            :options="diffEditorOptions"
            @load="handleDiffEditorLoad"
            class="h-full"
          />
          <template #fallback>
            <div class="h-full flex items-center justify-center text-gray-500">
              <div class="animate-spin rounded-full h-8 w-8 border-b-2 border-blue-600"></div>
            </div>
          </template>
        </ClientOnly>
      </div>

      <div v-else class="flex-1">
        <ClientOnly>
          <MonacoEditor
            ref="singleEditorRef"
            :value="getSingleEditorContent()"
            :lang="'json'"
            :options="singleEditorOptions"
            @load="handleSingleEditorLoad"
            class="h-full"
          />
          <template #fallback>
            <div class="h-full flex items-center justify-center text-gray-500">
              <div class="animate-spin rounded-full h-8 w-8 border-b-2 border-blue-600"></div>
            </div>
          </template>
        </ClientOnly>
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