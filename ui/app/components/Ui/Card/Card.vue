<template>
  <Primitive
    data-slot="card"
    :as="as"
    :as-child="asChild"
    :class="
      styles({ interactive: props.interactive, padding: props.padding, accent: props.accent, class: props.class })
    "
  >
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
      /** Team-based color accent variant */
      accent?: "none" | "primary" | "chart-1" | "chart-2" | "chart-3" | "chart-4"
    }
  >(),
  { as: "div", interactive: false, padding: "default", accent: "none" }
)

const styles = tv({
  base: "flex flex-col rounded-xl border border-border/40 bg-card/95 backdrop-blur-[2px] text-card-foreground shadow-sm ring-1 ring-border/20 transition-all duration-200",
  variants: {
    interactive: {
      true: "hover:shadow-md hover:ring-primary/30 hover:border-primary/30 hover:-translate-y-0.5 cursor-pointer hover:bg-card"
    },
    padding: {
      default: "gap-6 py-6",
      compact: "gap-3 p-4"
    },
    accent: {
      none: "",
      primary: "border-l-4 border-l-primary bg-gradient-to-r from-primary/5 to-transparent",
      "chart-1": "border-l-4 border-l-chart-1 bg-gradient-to-r from-chart-1/5 to-transparent",
      "chart-2": "border-l-4 border-l-chart-2 bg-gradient-to-r from-chart-2/5 to-transparent",
      "chart-3": "border-l-4 border-l-chart-3 bg-gradient-to-r from-chart-3/5 to-transparent",
      "chart-4": "border-l-4 border-l-chart-4 bg-gradient-to-r from-chart-4/5 to-transparent"
    }
  },
  defaultVariants: {
    padding: "default",
    accent: "none"
  }
})
</script>
