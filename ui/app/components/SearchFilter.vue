<script lang="ts">
export interface FilterConfig {
  key: string
  label: string
  options: Array<{ value: string; label: string; icon?: FunctionalComponent }>
}

export interface SearchConfig<TItem> {
  fields: Array<keyof TItem>
  placeholder: string
}

export interface DisplayConfig<TItem> {
  getItemText: (item: TItem) => string
  getItemValue: (item: TItem) => string
  getItemIcon?: (item: TItem) => FunctionalComponent
}
</script>

<script setup lang="ts" generic="TItem extends Record<string, any>">
import { onClickOutside } from "@vueuse/core"
import { X, User, CheckCircle, XCircle, Circle } from "lucide-vue-next"
import type { FunctionalComponent } from "vue"

const containerRef = ref(null)
onClickOutside(containerRef, () => {
  showFilters.value = false
})

const props = defineProps<{
  items: TItem[]
  filters: FilterConfig[]
  searchConfig: SearchConfig<TItem>
  displayConfig: DisplayConfig<TItem>
  maxSearchResults?: number
}>()

const emit = defineEmits<{
  filtersChanged: [filters: Record<string, string>]
  itemSelected: [item: TItem]
}>()

// State
const searchQuery = ref("")
const selectedFilters = ref<Record<string, string>>({})
const showFilters = ref(false)

// Initialize filter defaults
onMounted(() => {
  props.filters.forEach((filter) => {
    selectedFilters.value[filter.key] = "all"
  })
})

// Functions
function setFilters(newFilters: Record<string, string>) {
  selectedFilters.value = { ...selectedFilters.value, ...newFilters }
  emit("filtersChanged", { ...selectedFilters.value })
}

defineExpose({ setFilters })
function selectItem(item: TItem) {
  // Apply name to filters if filter exists
  const nameFilter = props.filters.find((f) => f.key === "name")
  if (nameFilter && item.name) {
    selectedFilters.value.name = item.name
    emit("filtersChanged", { ...selectedFilters.value })
  }
  emit("itemSelected", item)
  searchQuery.value = ""
  showFilters.value = false
}

function setFilter(key: string, value: string) {
  selectedFilters.value[key] = value
  showFilters.value = false
  emit("filtersChanged", { ...selectedFilters.value })
}

function clearFilter(key: string) {
  selectedFilters.value[key] = "all"
  emit("filtersChanged", { ...selectedFilters.value })
}

function clearAllFilters() {
  searchQuery.value = ""
  props.filters.forEach((filter) => {
    selectedFilters.value[filter.key] = "all"
  })
  showFilters.value = false
  emit("filtersChanged", { ...selectedFilters.value })
}

// Computed
const hasActiveFilters = computed(() => {
  return Object.values(selectedFilters.value).some((value) => value !== "all")
})

const searchResults = computed(() => {
  if (!searchQuery.value) {
    return {
      items: [],
      filterResults: props.filters.reduce(
        (acc, filter) => {
          acc[filter.key] = []
          return acc
        },
        {} as Record<string, any[]>
      )
    }
  }

  const query = searchQuery.value.toLowerCase()

  // Search items
  const matchingItems = props.items
    .filter((item) => props.searchConfig.fields.some((field) => String(item[field]).toLowerCase().includes(query)))
    .slice(0, props.maxSearchResults || 6)

  // Search filter options
  const filterResults: Record<string, any[]> = {}
  props.filters.forEach((filter) => {
    filterResults[filter.key] = filter.options.filter((option) => option.label.toLowerCase().includes(query))
  })

  return {
    items: matchingItems,
    filterResults
  }
})

function getFilterIcon(filter: FilterConfig, option: { value: string; label: string; icon?: FunctionalComponent }) {
  if (option.icon) return option.icon
  if (filter.key === "status") {
    return option.value === "active" ? CheckCircle : XCircle
  }
  return User
}

function getActiveFilterLabel(filterKey: string) {
  const filter = props.filters.find((f) => f.key === filterKey)
  const option = filter?.options.find((opt) => opt.value === selectedFilters.value[filterKey])
  return option?.label || selectedFilters.value[filterKey]
}
</script>

<template>
  <div class="space-y-4" ref="containerRef">
    <!-- Search & Filter Command -->
    <UiCommand class="rounded-lg border shadow-sm">
      <div class="flex items-center px-2" cmdk-input-wrapper @click="showFilters = true">
        <UiCommandInput
          v-model="searchQuery"
          :placeholder="searchConfig.placeholder"
          class="flex h-10 w-full rounded-md bg-transparent py-3 text-sm outline-none placeholder:text-muted-foreground disabled:cursor-not-allowed disabled:opacity-50"
        />
        <UiButton
          v-if="searchQuery"
          variant="ghost"
          size="sm"
          class="h-auto p-1 hover:bg-transparent"
          @click="searchQuery = ''"
        >
          <X class="size-4 text-muted-foreground hover:text-foreground" />
        </UiButton>
      </div>

      <UiCommandList v-if="searchQuery || showFilters">
        <UiCommandEmpty>No results found.</UiCommandEmpty>

        <!-- Filter Options (when not searching) -->
        <template v-if="!searchQuery">
          <UiCommandGroup v-for="filter in filters" :key="filter.key" :heading="`Filter by ${filter.label}`">
            <UiCommandItem
              v-for="option in filter.options"
              :key="option.value"
              :value="`${filter.key}:${option.value}`"
              :text="option.label"
              :icon="getFilterIcon(filter, option)"
              @select="setFilter(filter.key, option.value)"
            />
            <UiCommandItem
              :value="`${filter.key}:all`"
              :text="`All ${filter.label}`"
              :icon="Circle"
              @select="setFilter(filter.key, 'all')"
            />
          </UiCommandGroup>
          <UiCommandSeparator v-if="filters.length > 1" />
        </template>

        <!-- Search Results -->
        <template v-else>
          <!-- Items -->
          <UiCommandGroup v-if="searchResults.items.length" heading="Items">
            <UiCommandItem
              v-for="item in searchResults.items"
              :key="displayConfig.getItemValue(item)"
              :value="displayConfig.getItemValue(item)"
              :text="displayConfig.getItemText(item)"
              :icon="displayConfig.getItemIcon?.(item) || User"
              @select="selectItem(item)"
            />
          </UiCommandGroup>

          <!-- Filter Results -->
          <template v-for="filter in filters" :key="filter.key">
            <UiCommandGroup v-if="searchResults.filterResults[filter.key]?.length" :heading="filter.label">
              <UiCommandItem
                v-for="option in searchResults.filterResults[filter.key]"
                :key="option.value"
                :value="`${filter.key}:${option.value}`"
                :text="option.label"
                :icon="getFilterIcon(filter, option)"
                @select="setFilter(filter.key, option.value)"
              />
            </UiCommandGroup>
          </template>
        </template>
      </UiCommandList>
    </UiCommand>

    <!-- Active Filters -->
    <div v-if="hasActiveFilters" class="flex items-center justify-between gap-4">
      <div class="flex items-center gap-2">
        <span class="text-sm text-muted-foreground">Active filters:</span>
        <template v-for="filter in filters" :key="filter.key">
          <UiBadge v-if="selectedFilters[filter.key] !== 'all'" variant="secondary" class="gap-1">
            {{ filter.label }}: {{ getActiveFilterLabel(filter.key) }}
            <button @click="clearFilter(filter.key)" class="ml-1 hover:bg-muted-foreground/20 rounded-sm">
              <X class="h-3 w-3" />
            </button>
          </UiBadge>
        </template>
      </div>
      <UiButton variant="outline" size="sm" @click="clearAllFilters">Clear All</UiButton>
    </div>
  </div>
</template>
