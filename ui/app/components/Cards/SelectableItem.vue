<script setup lang="ts">
interface Props {
  selected?: boolean
  disabled?: boolean
  showCheckbox?: boolean
}

const props = withDefaults(defineProps<Props>(), {
  selected: false,
  disabled: false,
  showCheckbox: true
})

const emit = defineEmits<{
  click: []
}>()

const handleClick = () => {
  if (!props.disabled) {
    emit('click')
  }
}
</script>

<template>
  <div
    class="group relative border rounded-xl bg-card hover:shadow-lg transition-all duration-200 overflow-hidden cursor-pointer"
    :class="[
      selected ? 'ring-2 ring-primary shadow-lg' : 'hover:border-primary/30',
      disabled ? 'opacity-50 cursor-not-allowed' : ''
    ]"
    @click="handleClick"
  >
    <div class="p-6">
      <div class="flex items-start justify-between">
        <div class="flex-1 min-w-0">
          <slot name="header" />
          
          <div v-if=$slots.content class="mt-3">
            <slot name="content" />
          </div>
          
          <div v-if=$slots.footer class="mt-4">
            <slot name="footer" />
          </div>
        </div>

        <!-- Selection Indicator -->
        <div v-if="showCheckbox" class="ml-4 flex-shrink-0">
          <div
            v-if="selected"
            class="h-6 w-6 rounded-full bg-primary flex items-center justify-center"
          >
            <span class="text-white text-sm">✓</span>
          </div>
          <div
            v-else
            class="h-6 w-6 rounded-full border-2 border-muted group-hover:border-primary transition-colors"
          ></div>
        </div>
      </div>
    </div>

    <!-- Expandable content slot -->
    <slot name="expandable" />
  </div>
</template>