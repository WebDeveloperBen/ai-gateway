<script setup lang="ts">
interface DiffSegment {
  type: "unchanged" | "added" | "removed" | "modified"
  baseText?: string
  comparisonText?: string
  context: string
}

interface PolicyPair {
  id: string
  type: "aligned" | "unmatched_base" | "unmatched_comparison"
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

interface PolicyDiffData {
  metadata: {
    runId: string
    timestamp: string
    baseFilename: string
    comparisonFilename: string
    baseVersion: string
    comparisonVersion: string
    summary: string
  }
  pairs: PolicyPair[]
}

interface Props {
  diffData: PolicyDiffData
}

const props = defineProps<Props>()

const emit = defineEmits<{
  "select-section": [pairId: string]
}>()

// Document reconstruction
const baseDocument = computed(() => {
  return props.diffData.pairs
    .filter((pair) => pair.type === "aligned" || pair.type === "unmatched_base")
    .sort((a, b) => a.section.localeCompare(b.section))
    .map((pair) => ({
      id: pair.id,
      title: pair.title,
      section: pair.section,
      content: pair.baseContent || "",
      type: pair.type
    }))
})

const comparisonDocument = computed(() => {
  return props.diffData.pairs
    .filter((pair) => pair.type === "aligned" || pair.type === "unmatched_comparison")
    .sort((a, b) => a.section.localeCompare(b.section))
    .map((pair) => ({
      id: pair.id,
      title: pair.title,
      section: pair.section,
      content: pair.comparisonContent || "",
      type: pair.type
    }))
})

// Statistics
const stats = computed(() => {
  const aligned = props.diffData.pairs.filter((p) => p.type === "aligned").length
  const removed = props.diffData.pairs.filter((p) => p.type === "unmatched_base").length
  const added = props.diffData.pairs.filter((p) => p.type === "unmatched_comparison").length
  const totalChanges = props.diffData.pairs.reduce((acc, pair) => acc + pair.keyDifferences.length, 0)

  return { aligned, removed, added, total: aligned + removed + added, totalChanges }
})

const parseContent = (content: string) => {
  return content
    .replace(/\*\*(.*?)\*\*/g, "<strong>$1</strong>")
    .replace(/\*(.*?)\*/g, "<em>$1</em>")
    .replace(/`(.*?)`/g, '<code class="bg-gray-100 px-1 rounded text-sm">$1</code>')
    .replace(/\n\n/g, '</p><p class="mb-4">')
    .replace(/\n/g, " ")
}

const getSectionTypeClass = (type: string) => {
  switch (type) {
    case "unmatched_base":
      return "border-l-4 border-red-400 hover:bg-red-50"
    case "unmatched_comparison":
      return "border-l-4 border-green-400 hover:bg-green-50"
    case "aligned":
      return "border-l-4 border-blue-400 hover:bg-blue-50"
    default:
      return "border-l-4 border-gray-400 hover:bg-gray-50"
  }
}

const handleSectionClick = (pairId: string) => {
  emit("select-section", pairId)
}

// Synchronized scrolling
const baseScrollRef = ref<HTMLElement | null>(null)
const comparisonScrollRef = ref<HTMLElement | null>(null)
const isScrollingSynced = ref(true)

const onBaseScroll = () => {
  if (!isScrollingSynced.value || !baseScrollRef.value || !comparisonScrollRef.value) return
  const scrollPercentage =
    baseScrollRef.value.scrollTop / (baseScrollRef.value.scrollHeight - baseScrollRef.value.clientHeight)
  comparisonScrollRef.value.scrollTop =
    scrollPercentage * (comparisonScrollRef.value.scrollHeight - comparisonScrollRef.value.clientHeight)
}

const onComparisonScroll = () => {
  if (!isScrollingSynced.value || !baseScrollRef.value || !comparisonScrollRef.value) return
  const scrollPercentage =
    comparisonScrollRef.value.scrollTop /
    (comparisonScrollRef.value.scrollHeight - comparisonScrollRef.value.clientHeight)
  baseScrollRef.value.scrollTop =
    scrollPercentage * (baseScrollRef.value.scrollHeight - baseScrollRef.value.clientHeight)
}

const toggleScrollSync = () => {
  isScrollingSynced.value = !isScrollingSynced.value
}
</script>

<template>
  <div class="h-full flex flex-col bg-white">
    <!-- Document Overview Header -->
    <div class="border-b border-gray-200 px-6 py-4 bg-gradient-to-r from-gray-50 to-white">
      <div class="flex items-start justify-between mb-4">
        <div class="flex-1">
          <h3 class="text-xl font-semibold text-gray-900 mb-2">Full Document Comparison</h3>
          <p class="text-sm text-gray-600">{{ props.diffData.metadata.summary }}</p>
        </div>

        <div class="flex items-center gap-6">
          <!-- Statistics -->
          <div class="flex items-center gap-4 text-sm">
            <div class="text-center">
              <div class="text-lg font-semibold text-blue-600">{{ stats.aligned }}</div>
              <div class="text-xs text-gray-500">Modified</div>
            </div>
            <div class="text-center">
              <div class="text-lg font-semibold text-red-600">{{ stats.removed }}</div>
              <div class="text-xs text-gray-500">Removed</div>
            </div>
            <div class="text-center">
              <div class="text-lg font-semibold text-green-600">{{ stats.added }}</div>
              <div class="text-xs text-gray-500">Added</div>
            </div>
            <div class="text-center">
              <div class="text-lg font-semibold text-orange-600">{{ stats.totalChanges }}</div>
              <div class="text-xs text-gray-500">Total Changes</div>
            </div>
          </div>

          <!-- Version Info -->
          <div class="flex items-center gap-4 text-sm">
            <div class="text-center">
              <div class="text-xs text-gray-500 mb-1">Base</div>
              <UiBadge variant="outline" size="sm">v{{ props.diffData.metadata.baseVersion }}</UiBadge>
            </div>
            <div class="w-6 h-px bg-gray-300"></div>
            <div class="text-center">
              <div class="text-xs text-gray-500 mb-1">Updated</div>
              <UiBadge variant="default" size="sm">v{{ props.diffData.metadata.comparisonVersion }}</UiBadge>
            </div>
          </div>
        </div>
      </div>
    </div>

    <!-- Document Content -->
    <div class="flex-1 flex overflow-hidden h-0">
      <!-- Base Document -->
      <div class="w-1/2 border-r border-gray-200 flex flex-col">
        <div class="px-6 py-3 bg-red-50 border-b border-red-100 flex items-center justify-between">
          <div class="flex items-center gap-3">
            <div class="w-3 h-3 bg-red-500 rounded-full"></div>
            <span class="text-sm font-medium text-red-800">{{ props.diffData.metadata.baseFilename }}</span>
          </div>
          <button
            @click="toggleScrollSync"
            :class="[
              'text-xs px-2 py-1 rounded transition-colors',
              isScrollingSynced
                ? 'text-red-700 bg-red-100 hover:bg-red-200'
                : 'text-gray-600 bg-gray-100 hover:bg-gray-200'
            ]"
            :title="isScrollingSynced ? 'Disable scroll sync' : 'Enable scroll sync'"
          >
            {{ isScrollingSynced ? "🔗" : "🔗❌" }} Sync
          </button>
        </div>
        <div ref="baseScrollRef" @scroll="onBaseScroll" class="flex-1 overflow-y-auto" style="height: 600px;">
          <div class="p-6 space-y-6">
            <div
              v-for="section in baseDocument"
              :key="section.id"
              :class="[
                'p-4 rounded-lg cursor-pointer transition-all hover:shadow-sm border border-gray-200',
                getSectionTypeClass(section.type)
              ]"
              @click="handleSectionClick(section.id)"
            >
              <div class="flex items-start justify-between mb-3">
                <h4 class="text-lg font-semibold text-gray-900 flex-1">{{ section.title }}</h4>
                <span class="text-xs text-gray-500 px-2 py-1 bg-gray-100 rounded-full ml-3 flex-shrink-0">{{
                  section.section
                }}</span>
              </div>
              <div
                class="prose prose-sm max-w-none text-gray-800 leading-relaxed"
                v-html="parseContent(section.content)"
              ></div>
            </div>
          </div>
        </div>
      </div>

      <!-- Comparison Document -->
      <div class="w-1/2 flex flex-col">
        <div class="px-6 py-3 bg-green-50 border-b border-green-100 flex items-center gap-3">
          <div class="w-3 h-3 bg-green-500 rounded-full"></div>
          <span class="text-sm font-medium text-green-800">{{ props.diffData.metadata.comparisonFilename }}</span>
        </div>
        <div ref="comparisonScrollRef" @scroll="onComparisonScroll" class="flex-1 overflow-y-auto" style="height: 600px;">
          <div class="p-6 space-y-6">
            <div
              v-for="section in comparisonDocument"
              :key="section.id"
              :class="[
                'p-4 rounded-lg cursor-pointer transition-all hover:shadow-sm border border-gray-200',
                getSectionTypeClass(section.type)
              ]"
              @click="handleSectionClick(section.id)"
            >
              <div class="flex items-start justify-between mb-3">
                <h4 class="text-lg font-semibold text-gray-900 flex-1">{{ section.title }}</h4>
                <span class="text-xs text-gray-500 px-2 py-1 bg-gray-100 rounded-full ml-3 flex-shrink-0">{{
                  section.section
                }}</span>
              </div>
              <div
                class="prose prose-sm max-w-none text-gray-800 leading-relaxed"
                v-html="parseContent(section.content)"
              ></div>
            </div>
          </div>
        </div>
      </div>
    </div>

    <!-- Quick Stats Footer -->
    <div class="border-t border-gray-200 bg-gray-50 px-6 py-3">
      <div class="flex items-center justify-between text-sm text-gray-600">
        <div class="flex items-center gap-6">
          <span>{{ stats.total }} sections compared</span>
          <span>{{ stats.totalChanges }} total changes identified</span>
        </div>
        <div class="flex items-center gap-4">
          <div class="flex items-center gap-2">
            <div class="w-3 h-3 border-l-4 border-blue-400 bg-blue-50"></div>
            <span>Modified sections</span>
          </div>
          <div class="flex items-center gap-2">
            <div class="w-3 h-3 border-l-4 border-red-400 bg-red-50"></div>
            <span>Removed sections</span>
          </div>
          <div class="flex items-center gap-2">
            <div class="w-3 h-3 border-l-4 border-green-400 bg-green-50"></div>
            <span>Added sections</span>
          </div>
        </div>
      </div>
    </div>
  </div>
</template>

<style scoped>
.prose {
  color: rgb(55 65 81);
  line-height: 1.75;
}

.prose p {
  margin-bottom: 1rem;
}

.prose strong {
  font-weight: 600;
  color: rgb(17 24 39);
}

.prose em {
  font-style: italic;
  color: rgb(75 85 99);
}
</style>
