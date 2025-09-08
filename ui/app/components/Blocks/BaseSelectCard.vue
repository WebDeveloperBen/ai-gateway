<script setup lang="ts" generic="T">
const selected = defineModel<boolean>("selected", { default: false })
const expandedModel = defineModel<boolean>("expanded")

interface SelectableCardProps<T = any> {
  /** The data object for this card */
  item: T
  /** Disable interactions */
  disabled?: boolean
  /** Selection mode controls semantics & indicator */
  mode?: "single" | "multiple"
  /** Show checkbox (alias for mode === 'multiple'; kept for back-compat) */
  showCheckbox?: boolean
  /** Expand when selected (if expanded v-model not provided) */
  expandOnSelect?: boolean
  /** Root element */
  as?: string
  /** Extra classes */
  class?: any
}

const props = withDefaults(defineProps<SelectableCardProps<T>>(), {
  disabled: false,
  mode: "single",
  showCheckbox: undefined,
  expandOnSelect: false,
  as: "div"
})

const emit = defineEmits<{
  /** Back-compat events */
  select: [item: T]
  toggle: [item: T]
}>()

/** Styling with tailwind-variants */
const cardStyles = tv({
  base: "group relative overflow-hidden rounded-2xl border bg-card text-card-foreground transition-all duration-200",
  variants: {
    interactive: { true: "cursor-pointer hover:shadow-lg hover:border-primary/30" },
    selected: {
      true: "ring-2 ring-primary shadow-lg border-primary/50",
      false: ""
    },
    disabled: { true: "opacity-50 cursor-not-allowed hover:shadow-none hover:border-border" },
    padded: { true: "p-6", false: "" }
  },
  defaultVariants: { interactive: true, padded: true }
})

const indicatorStyles = tv({
  variants: {
    mode: {
      single: "absolute top-4 right-4",
      multiple: "absolute top-4 left-4"
    },
    selected: {
      true: "",
      false: ""
    }
  }
})

const showCheckbox = computed(() => props.showCheckbox ?? props.mode === "multiple")

const role = computed(() => (showCheckbox.value ? "checkbox" : "option"))
const ariaChecked = computed(() => (showCheckbox.value ? selected.value : undefined))

const isExpanded = computed({
  get: () => expandedModel.value ?? (props.expandOnSelect ? selected.value : false),
  set: (v: boolean) => {
    expandedModel.value = v
  }
})

function shouldIgnoreClick(target: EventTarget | null) {
  if (!(target instanceof HTMLElement)) return false
  return !!target.closest("a,button,input,textarea,select,label,[role='button'],[data-prevent-toggle]")
}

function onToggle() {
  if (props.disabled) return
  if (showCheckbox.value) {
    selected.value = !selected.value
    emit("toggle", props.item)
  } else {
    // single-select: selecting this card sets selected = true
    const wasSelected = selected.value
    selected.value = true
    if (!wasSelected) emit("select", props.item)
  }
}

function onClick(e: MouseEvent) {
  if (props.disabled) return
  if (shouldIgnoreClick(e.target)) return
  onToggle()
  if (props.expandOnSelect && selected.value) {
    isExpanded.value = true
  }
}

function onKeydown(e: KeyboardEvent) {
  if (props.disabled) return
  if (e.key === "Enter" || e.key === " ") {
    e.preventDefault()
    onToggle()
  }
}
</script>

<template>
  <component
    :is="as"
    :role="role"
    :aria-selected="!showCheckbox ? selected : undefined"
    :aria-checked="ariaChecked"
    :aria-disabled="disabled || undefined"
    :tabindex="disabled ? -1 : 0"
    :class="cardStyles({ selected, disabled, class: props.class })"
    @click="onClick"
    @keydown="onKeydown"
  >
    <!-- Indicator area -->
    <div :class="indicatorStyles({ mode: showCheckbox ? 'multiple' : 'single', selected })">
      <!-- Multiple: checkbox on the left -->
      <div v-if="showCheckbox" class="z-10" @click.stop>
        <UiCheckbox v-model:checked="selected" :disabled="disabled" />
      </div>

      <!-- Single: radio-like dot on the right -->
      <div v-else class="z-10">
        <div v-if="selected" class="h-6 w-6 rounded-full bg-primary flex items-center justify-center">
          <span class="text-white text-sm">✓</span>
        </div>
        <div v-else class="h-6 w-6 rounded-full border-2 border-muted group-hover:border-primary transition-colors" />
      </div>
    </div>

    <!-- Content -->
    <div :class="[{ 'pl-14': showCheckbox }, 'relative']">
      <!-- Header -->
      <div class="mb-4">
        <slot name="header" :item="item" :selected="selected" />
      </div>

      <!-- Body -->
      <div class="space-y-3">
        <slot name="content" :item="item" :selected="selected" />
      </div>

      <!-- Footer -->
      <div class="mt-4">
        <slot name="footer" :item="item" :selected="selected" />
      </div>
    </div>

    <!-- Expanded -->
    <Transition name="fade-collapse">
      <div v-if="isExpanded && $slots.expanded" class="border-t bg-muted/10">
        <slot name="expanded" :item="item" />
      </div>
    </Transition>
  </component>
</template>

<style scoped>
.fade-collapse-enter-active,
.fade-collapse-leave-active {
  transition:
    opacity 150ms ease,
    max-height 200ms ease;
}
.fade-collapse-enter-from,
.fade-collapse-leave-to {
  opacity: 0;
}
</style>
