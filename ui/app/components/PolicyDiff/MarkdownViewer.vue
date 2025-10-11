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

interface Props {
  pair: PolicyPair
}

const props = defineProps<Props>()

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

// Copy functionality
const copyContent = async (content: string) => {
  try {
    await navigator.clipboard.writeText(content.replace(/<[^>]*>/g, "")) // Strip HTML tags
    // Could add toast notification here
  } catch (err) {
    console.error("Failed to copy content:", err)
  }
}

// Keyboard navigation
const emit = defineEmits<{
  "navigate-next": []
  "navigate-previous": []
}>()

onMounted(() => {
  const handleKeyboard = (event: KeyboardEvent) => {
    if (event.key === "ArrowRight" && (event.metaKey || event.ctrlKey)) {
      emit("navigate-next")
      event.preventDefault()
    } else if (event.key === "ArrowLeft" && (event.metaKey || event.ctrlKey)) {
      emit("navigate-previous")
      event.preventDefault()
    } else if (event.key === "s" && (event.metaKey || event.ctrlKey)) {
      toggleScrollSync()
      event.preventDefault()
    }
  }

  document.addEventListener("keydown", handleKeyboard)

  onUnmounted(() => {
    document.removeEventListener("keydown", handleKeyboard)
  })
})

const renderHighlightedContent = (segments: DiffSegment[], isComparison = false) => {
  if (!segments || segments.length === 0) {
    return isComparison ? props.pair.comparisonContent : props.pair.baseContent
  }

  return segments
    .filter((segment) => {
      const text = isComparison ? segment.comparisonText : segment.baseText
      return text && text.trim().length > 0
    })
    .map((segment) => {
      const text = isComparison ? segment.comparisonText : segment.baseText
      const type = segment.type

      if (type === "unchanged") {
        return text
      } else if (type === "modified" || type === "added" || type === "removed") {
        const className = isComparison
          ? type === "modified" || type === "added"
            ? "bg-green-100 text-green-900 px-1 rounded"
            : ""
          : type === "modified" || type === "removed"
            ? "bg-red-100 text-red-900 px-1 rounded"
            : ""
        return `<mark class="${className}">${text}</mark>`
      }

      return text
    })
    .join(" ")
}

const parseContent = (content: string) => {
  // Simple markdown-like parsing for better formatting
  return content
    .replace(/\*\*(.*?)\*\*/g, "<strong>$1</strong>") // Bold
    .replace(/\*(.*?)\*/g, "<em>$1</em>") // Italic
    .replace(/`(.*?)`/g, '<code class="bg-gray-100 px-1 rounded text-sm">$1</code>') // Code
    .replace(/\n\n/g, '</p><p class="mb-4">') // Paragraphs
    .replace(/\n/g, " ") // Single line breaks become spaces
}
</script>

<template>
  <div class="flex flex-col">
    <!-- Enhanced Header -->
    <div class="border-b border-gray-200 px-6 py-5 bg-gradient-to-r from-gray-50 to-white">
      <div class="flex items-start justify-between mb-4">
        <div class="flex-1">
          <h3 class="text-xl font-semibold text-gray-900 mb-2">{{ props.pair.title }}</h3>
          <div class="flex items-center gap-3 text-sm text-gray-600">
            <span class="px-2 py-1 bg-blue-100 text-blue-800 rounded-full text-xs font-medium">
              {{ props.pair.section }}
            </span>
            <span class="text-gray-400">•</span>
            <span>{{ props.pair.sourceTrace }}</span>
          </div>
        </div>

        <div v-if="props.pair.type === 'aligned'" class="flex items-center gap-4 text-sm">
          <div class="text-center">
            <div class="text-xs text-gray-500 mb-1">Base</div>
            <UiBadge variant="outline" size="sm">v{{ props.pair.baseVersion }}</UiBadge>
          </div>
          <div class="w-6 h-px bg-gray-300"></div>
          <div class="text-center">
            <div class="text-xs text-gray-500 mb-1">Updated</div>
            <UiBadge variant="default" size="sm">v{{ props.pair.comparisonVersion }}</UiBadge>
          </div>
        </div>
        <div v-else class="text-center">
          <div class="text-xs text-gray-500 mb-1">
            {{ props.pair.type === "unmatched_base" ? "Removed" : "Added" }}
          </div>
          <UiBadge :variant="props.pair.type === 'unmatched_base' ? 'destructive' : 'default'" size="sm">
            v{{ props.pair.type === "unmatched_base" ? props.pair.baseVersion : props.pair.comparisonVersion }}
          </UiBadge>
        </div>
      </div>
    </div>

    <div class="flex-1 flex flex-col overflow-hidden">
      <!-- Policy Content Display -->
      <div v-if="props.pair.type === 'aligned'" class="flex-1 grid grid-cols-2 overflow-hidden">
        <!-- Base Version -->
        <div class="border-r border-gray-200 flex flex-col">
          <div class="px-6 py-3 bg-red-50 border-b border-red-100 flex items-center justify-between">
            <div class="flex items-center gap-3">
              <div class="w-3 h-3 bg-red-500 rounded-full"></div>
              <span class="text-sm font-medium text-red-800">Original Version</span>
            </div>
            <button
              @click="copyContent(renderHighlightedContent(props.pair.diffSegments, false))"
              class="text-xs text-red-700 hover:text-red-900 px-2 py-1 rounded hover:bg-red-100 transition-colors"
              title="Copy original text"
            >
              Copy
            </button>
          </div>
          <div ref="baseScrollRef" @scroll="onBaseScroll" class="flex-1 overflow-y-auto p-6 scroll-smooth h-0">
            <div
              class="prose prose-sm max-w-none leading-relaxed text-gray-800 select-text"
              v-html="parseContent(renderHighlightedContent(props.pair.diffSegments, false))"
            ></div>
          </div>
        </div>

        <!-- Comparison Version -->
        <div class="flex flex-col">
          <div class="px-6 py-3 bg-green-50 border-b border-green-100 flex items-center justify-between">
            <div class="flex items-center gap-3">
              <div class="w-3 h-3 bg-green-500 rounded-full"></div>
              <span class="text-sm font-medium text-green-800">Updated Version</span>
            </div>
            <div class="flex items-center gap-2">
              <button
                @click="toggleScrollSync"
                :class="[
                  'text-xs px-2 py-1 rounded transition-colors',
                  isScrollingSynced
                    ? 'text-green-700 bg-green-100 hover:bg-green-200'
                    : 'text-gray-600 bg-gray-100 hover:bg-gray-200'
                ]"
                :title="isScrollingSynced ? 'Disable scroll sync' : 'Enable scroll sync'"
              >
                {{ isScrollingSynced ? "🔗" : "🔗❌" }} Sync
              </button>
              <button
                @click="copyContent(renderHighlightedContent(props.pair.diffSegments, true))"
                class="text-xs text-green-700 hover:text-green-900 px-2 py-1 rounded hover:bg-green-100 transition-colors"
                title="Copy updated text"
              >
                Copy
              </button>
            </div>
          </div>
          <div
            ref="comparisonScrollRef"
            @scroll="onComparisonScroll"
            class="flex-1 overflow-y-auto p-6 scroll-smooth h-0"
          >
            <div
              class="prose prose-sm max-w-none leading-relaxed text-gray-800 select-text"
              v-html="parseContent(renderHighlightedContent(props.pair.diffSegments, true))"
            ></div>
          </div>
        </div>
      </div>

      <!-- Single Document View -->
      <div v-else class="flex-1 flex flex-col">
        <div class="px-6 py-3 bg-gray-50 border-b border-gray-200 flex items-center gap-3">
          <div
            :class="['w-3 h-3 rounded-full', props.pair.type === 'unmatched_base' ? 'bg-red-500' : 'bg-green-500']"
          ></div>
          <span class="text-sm font-medium text-gray-800">
            {{ props.pair.type === "unmatched_base" ? "Removed Content" : "New Content" }}
          </span>
        </div>
        <div class="flex-1 overflow-y-auto p-6 h-0">
          <div class="prose prose-sm max-w-none leading-relaxed text-gray-800">
            <p
              v-html="
                parseContent(
                  props.pair.type === 'unmatched_base'
                    ? props.pair.baseContent || ''
                    : props.pair.comparisonContent || ''
                )
              "
            ></p>
          </div>
        </div>
      </div>

      <!-- Enhanced Analysis Section -->
      <div class="border-t border-gray-200 bg-gradient-to-r from-gray-50 to-white">
        <div class="p-6">
          <div class="grid grid-cols-1 lg:grid-cols-2 gap-8">
            <!-- Similarities -->
            <div v-if="props.pair.similarities?.length" class="space-y-4">
              <div class="flex items-center gap-3">
                <div class="flex items-center justify-center w-8 h-8 bg-blue-100 rounded-full">
                  <div class="w-2 h-2 bg-blue-600 rounded-full"></div>
                </div>
                <h4 class="text-sm font-semibold text-gray-900">Preserved Elements</h4>
              </div>
              <ul class="space-y-3 ml-11">
                <li
                  v-for="(similarity, index) in props.pair.similarities"
                  :key="index"
                  class="flex items-start gap-3 text-sm leading-relaxed"
                >
                  <div class="w-1.5 h-1.5 rounded-full bg-blue-400 mt-2 flex-shrink-0"></div>
                  <span class="text-gray-700">{{ similarity }}</span>
                </li>
              </ul>
            </div>

            <!-- Key Differences -->
            <div v-if="props.pair.keyDifferences?.length" class="space-y-4">
              <div class="flex items-center gap-3">
                <div class="flex items-center justify-center w-8 h-8 bg-orange-100 rounded-full">
                  <div class="w-2 h-2 bg-orange-600 rounded-full"></div>
                </div>
                <h4 class="text-sm font-semibold text-gray-900">Key Changes</h4>
              </div>
              <ul class="space-y-3 ml-11">
                <li
                  v-for="(difference, index) in props.pair.keyDifferences"
                  :key="index"
                  class="flex items-start gap-3 text-sm leading-relaxed"
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
  </div>
</template>

<style scoped>
/* Custom prose styling for better readability */
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

/* Smooth highlighting */
.prose mark {
  transition: all 0.2s ease;
  font-weight: 500;
}

.prose mark:hover {
  filter: brightness(0.95);
}
</style>

