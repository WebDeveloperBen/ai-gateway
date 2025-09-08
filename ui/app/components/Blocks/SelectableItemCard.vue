<script setup lang="ts">
import { computed } from "vue"
import { Tag, History, Settings } from "lucide-vue-next"

interface PromptVersion {
  id: string
  version: string | number
  name: string
  createdBy: string
  createdAt: string
  content?: string
}

interface PromptItem {
  id: string
  name: string
  description?: string
  tags: string[]
  currentVersion: string | number
  versions: PromptVersion[]
  applications: any[]
}

const selected = defineModel<boolean>("selected")
const selectedVersionId = defineModel<string | null>("selectedVersionId")

const props = withDefaults(
  defineProps<{
    prompt: PromptItem
    /** Auto expand when selected if v-model:expanded is not supplied */
    expandOnSelect?: boolean
    /** Provide a token estimator; defaults to naive char/4 */
    estimateTokens?: (text: string) => number
    /** Pass extra classes to the root card */
    class?: any
  }>(),
  {
    expandOnSelect: true
  }
)

const emit = defineEmits<{
  (e: "select", prompt: PromptItem): void
  (e: "select-version", version: PromptVersion): void
}>()

const selectedVersion = computed(() => props.prompt.versions.find((v) => v.id === selectedVersionId.value) || null)

const tokenEstimator = computed(() => props.estimateTokens || ((t: string) => Math.ceil((t?.length || 0) / 4)))

function selectVersion(version: PromptVersion) {
  selectedVersionId.value = version.id
  emit("select-version", version)
}
</script>

<template>
  <BlocksSelectableCard
    :item="prompt"
    v-model:selected="selected"
    :class="props.class"
    :expand-on-select="expandOnSelect"
  >
    <!-- HEADER -->
    <template #header="{ item }">
      <div class="flex items-start justify-between">
        <div class="flex-1">
          <div class="flex items-center gap-3 mb-2">
            <h4 class="font-semibold text-lg text-foreground group-hover:text-primary transition-colors">
              {{ item.name }}
            </h4>
            <span
              class="inline-flex items-center px-2 py-1 text-xs font-medium rounded-full bg-blue-100 text-blue-700 dark:bg-blue-900/20 dark:text-blue-400"
            >
              Live: v{{ item.currentVersion }}
            </span>
            <!-- Optional extra badges slot -->
            <slot name="badges" :item="item" />
          </div>

          <p v-if="item.description" class="text-muted-foreground text-sm mb-3 leading-relaxed">
            {{ item.description }}
          </p>

          <!-- Tags -->
          <div class="flex flex-wrap gap-1 mb-3">
            <span
              v-for="tag in item.tags.slice(0, 3)"
              :key="tag"
              class="inline-flex items-center gap-1 px-2 py-1 text-xs rounded-md bg-primary/10 text-primary"
            >
              <Tag class="h-3 w-3" />
              {{ tag }}
            </span>
            <span
              v-if="item.tags.length > 3"
              class="inline-flex items-center px-2 py-1 text-xs rounded-md bg-muted text-muted-foreground"
            >
              +{{ item.tags.length - 3 }} more
            </span>
          </div>

          <!-- Metadata Row -->
          <div class="flex items-center gap-4 text-xs text-muted-foreground">
            <div class="flex items-center gap-1">
              <History class="h-3 w-3" />
              <span>{{ item.versions.length }} versions</span>
            </div>
            <div class="flex items-center gap-1">
              <Settings class="h-3 w-3" />
              <span>{{ item.applications.length }} apps</span>
            </div>
            <slot name="meta-extra" :item="item" />
          </div>
        </div>

        <!-- Selection indicator mirrors BlocksSelectableCard single-select -->
        <div class="ml-4">
          <div v-if="selected" class="h-6 w-6 rounded-full bg-primary grid place-items-center">
            <span class="text-white text-sm leading-none">✓</span>
          </div>
          <div v-else class="h-6 w-6 rounded-full border-2 border-muted group-hover:border-primary transition-colors" />
        </div>
      </div>
    </template>

    <!-- CONTENT (keep empty; header holds primary text) -->
    <template #content="{ item }">
      <slot name="content" :item="item" />
    </template>

    <!-- EXPANDED AREA: version picker + preview -->
    <template #expanded="{ item }">
      <!-- Version Selection -->
      <div class="border-b bg-muted/20">
        <div class="p-4">
          <h5 class="text-sm font-medium mb-3 flex items-center gap-2">
            <span>Select Version</span>
            <span class="text-xs text-muted-foreground">({{ item.versions.length }} available)</span>
          </h5>
          <div class="max-h-40 overflow-y-auto">
            <div class="space-y-2 pr-1">
              <div
                v-for="version in item.versions.slice().reverse()"
                :key="version.id"
                class="flex items-center justify-between p-3 border rounded-lg cursor-pointer hover:bg-muted/50 transition-colors"
                :class="selectedVersionId === version.id ? 'border-primary bg-primary/5' : 'border-border'"
                @click.stop="selectVersion(version)"
              >
                <div class="flex-1">
                  <div class="flex items-center gap-2 mb-1">
                    <span class="font-mono text-sm font-medium">{{ version.version }}</span>
                    <span class="text-sm text-foreground">{{ version.name }}</span>
                    <span
                      v-if="version.version === item.currentVersion"
                      class="inline-flex items-center px-1.5 py-0.5 text-xs font-medium rounded bg-blue-100 text-blue-700 dark:bg-blue-900/20 dark:text-blue-400"
                    >
                      Live
                    </span>
                  </div>
                  <div class="text-xs text-muted-foreground">
                    by {{ version.createdBy.split("@")[0] }} • {{ new Date(version.createdAt).toLocaleDateString() }}
                  </div>
                </div>
                <div class="ml-3">
                  <div
                    v-if="selectedVersionId === version.id"
                    class="h-4 w-4 rounded-full bg-primary grid place-items-center"
                  >
                    <span class="text-white text-[10px] leading-none">✓</span>
                  </div>
                  <div v-else class="h-4 w-4 rounded-full border border-muted"></div>
                </div>
              </div>
            </div>
          </div>
        </div>
      </div>

      <!-- Content Preview -->
      <div v-if="selectedVersion" class="bg-muted/10">
        <div class="p-4">
          <h5 class="text-sm font-medium mb-3 flex items-center justify-between">
            <span>Content Preview</span>
            <span class="text-xs text-muted-foreground">
              {{ (selectedVersion.content || "").length.toLocaleString() }} chars • ~{{
                tokenEstimator(selectedVersion.content || "").toLocaleString()
              }}
              tokens
            </span>
          </h5>
          <div class="rounded-lg border bg-background/50 p-4 max-h-32 overflow-y-auto">
            <pre class="text-xs font-mono leading-relaxed whitespace-pre-wrap text-foreground">{{
              selectedVersion.content || "No content available"
            }}</pre>
          </div>
        </div>
      </div>

      <slot name="expanded-extra" :item="item" :version="selectedVersion" />
    </template>
  </BlocksSelectableCard>
</template>
