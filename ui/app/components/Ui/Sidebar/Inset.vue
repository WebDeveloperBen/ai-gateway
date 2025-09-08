<template>
  <main data-slot="sidebar-inset" :class="sideBarInsetStyles({ class: props.class })">
    <slot />
  </main>
</template>
<script lang="ts">
import type { HTMLAttributes } from "vue"

export const sideBarInsetStyles = tv({
  base: "relative flex min-w-0 flex-1 flex-col bg-background md:peer-data-[variant=inset]:m-2 md:peer-data-[variant=inset]:ml-0 md:peer-data-[variant=inset]:rounded-xl md:peer-data-[variant=inset]:shadow-sm md:peer-data-[variant=inset]:peer-data-[state=collapsed]:ml-2 md:w-[calc(100%-var(--sidebar-width))] md:peer-data-[state=collapsed]:w-[calc(100%-var(--sidebar-width-icon))]"
})
</script>
<script setup lang="ts">
const props = defineProps<{
  /**
   * Additional classes to apply to the parent element.
   */
  class?: HTMLAttributes["class"]
}>()
</script>
