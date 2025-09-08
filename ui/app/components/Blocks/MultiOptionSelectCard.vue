<script setup lang="ts" generic="T">
interface SelectableCardProps<T = any> {
  /** The data object for this card */
  item: T
  /** Whether this card is selected */
  selected?: boolean
  /** Whether to show checkbox */
  showCheckbox?: boolean
  /** Whether selection is disabled */
  disabled?: boolean
  /** Additional CSS classes */
  class?: string
}

const props = withDefaults(defineProps<SelectableCardProps<T>>(), {
  selected: false,
  showCheckbox: false,
  disabled: false
})

const emit = defineEmits<{
  select: [item: T]
  toggle: [item: T]
}>()

const handleClick = () => {
  if (props.disabled) return

  if (props.showCheckbox) {
    emit("toggle", props.item)
  } else {
    emit("select", props.item)
  }
}

const handleCheckboxChange = () => {
  if (props.disabled) return
  emit("toggle", props.item)
}
</script>

<template>
  <div
    class="group relative border rounded-xl bg-card transition-all duration-200 overflow-hidden cursor-pointer"
    :class="[
      selected ? 'ring-2 ring-primary shadow-lg border-primary/50' : 'hover:border-primary/30 hover:shadow-lg',
      disabled && 'opacity-50 cursor-not-allowed hover:border-border hover:shadow-none',
      props.class
    ]"
    @click="handleClick"
  >
    <!-- Checkbox (if enabled) -->
    <div v-if="showCheckbox" class="absolute top-4 left-4 z-10" @click.stop>
      <UiCheckbox :checked="selected" :disabled="disabled" @update:checked="handleCheckboxChange" />
    </div>

    <!-- Selection Indicator (for single select mode) -->
    <div v-else class="absolute top-4 right-4 z-10">
      <div v-if="selected" class="h-6 w-6 rounded-full bg-primary flex items-center justify-center">
        <span class="text-white text-sm">✓</span>
      </div>
      <div
        v-else
        class="h-6 w-6 rounded-full border-2 border-muted group-hover:border-primary transition-colors"
        :class="disabled && 'group-hover:border-muted'"
      />
    </div>

    <!-- Card Content -->
    <div class="p-6" :class="showCheckbox ? 'pl-14' : ''">
      <!-- Header slot -->
      <div class="mb-4">
        <slot name="header" :item="item" :selected="selected" />
      </div>

      <!-- Content slot -->
      <div class="space-y-3">
        <slot name="content" :item="item" :selected="selected" />
      </div>

      <!-- Footer slot -->
      <div class="mt-4">
        <slot name="footer" :item="item" :selected="selected" />
      </div>
    </div>

    <!-- Expandable content (shown when selected) -->
    <div v-if="selected && $slots.expanded" class="border-t bg-muted/10">
      <slot name="expanded" :item="item" />
    </div>
  </div>
</template>

