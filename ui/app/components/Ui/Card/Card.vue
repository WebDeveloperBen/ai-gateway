<template>
  <Primitive data-slot="card" :as="as" :as-child="asChild" :class="styles({ interactive: props.interactive, padding: props.padding, class: props.class })">
    <slot>
      <slot name="header">
        <UiCardHeader>
          <slot name="title">
            <UiCardTitle v-if="title || $slots.title" :title="title" />
          </slot>
          <slot name="description">
            <UiCardDescription v-if="description || $slots.description" :description="description" />
          </slot>
        </UiCardHeader>
      </slot>
      <slot v-if="content || $slots.content" name="content">
        <UiCardContent>
          <div v-html="content" />
        </UiCardContent>
      </slot>
      <slot name="footer" />
    </slot>
  </Primitive>
</template>

<script lang="ts" setup>
import { Primitive } from "reka-ui"
import type { PrimitiveProps } from "reka-ui"
import type { HTMLAttributes } from "vue"

const props = withDefaults(
  defineProps<
    PrimitiveProps & {
      /** Title that should be displayed. Passed to the `CardTitle` component */
      title?: string
      /** Description that should be displayed. Passed to the `CardDescription` component */
      description?: string
      /** Content that should be displayed. Passed to the `CardContent` component */
      content?: string
      /** Custom class(es) to add to the element */
      class?: HTMLAttributes["class"]
      /** Whether the card should have hover effects */
      interactive?: boolean
      /** Padding variant */
      padding?: "default" | "compact"
    }
  >(),
  { as: "div", interactive: false, padding: "default" }
)

const styles = tv({
  base: "flex flex-col rounded-xl border bg-card text-card-foreground shadow-sm transition-all duration-200",
  variants: {
    interactive: {
      true: "hover:shadow-md hover:border-primary/20 hover:-translate-y-0.5 cursor-pointer"
    },
    padding: {
      default: "gap-6 py-6",
      compact: "gap-3 p-4"
    }
  },
  defaultVariants: {
    padding: "default"
  }
})
</script>
