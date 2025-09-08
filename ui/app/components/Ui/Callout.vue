<script setup lang="ts">
import { computed, type Component } from "vue"
import { tv, type VariantProps } from "tailwind-variants"

/**
 * Accessible, polymorphic Callout built with tailwind-variants.
 * - Back-compat: accepts `variant` (color) like your old API
 * - New: `color` (alias of `variant`), `tone`, `size`, `withAccent`, `dismissible`, `shadow`, `as`
 * - Slots: icon, title, description, default (content), actions
 * - Emits: v-model:open (for dismiss)
 */

export const calloutStyles = tv({
  slots: {
    root: "relative rounded-lg border p-4",
    body: "flex gap-3",
    icon: "w-5 h-5 shrink-0 mt-0.5",
    content: "text-sm flex-1",
    title: "font-medium mb-1",
    description: "",
    actions: "mt-3 flex flex-wrap gap-2",
    close:
      "absolute top-2 right-2 inline-flex items-center justify-center rounded-md p-1 outline-none focus-visible:ring-2 focus-visible:ring-ring/50 hover:bg-black/5 dark:hover:bg-white/5"
  },
  variants: {
    size: {
      sm: {
        root: "p-3",
        icon: "w-4 h-4 mt-0",
        title: "text-[13px]",
        description: "text-xs",
        actions: "mt-2"
      },
      md: {
        root: "p-4",
        title: "text-sm",
        description: "text-sm"
      }
    },
    tone: {
      /** Subtle = soft surface, colored text */
      subtle: {},
      /** Solid = strong surface, white text */
      solid: {},
      /** Outline = transparent bg, colored border/text */
      outline: {}
    },
    color: {
      blue: {},
      green: {},
      yellow: {},
      red: {},
      purple: {},
      gray: {}
    },
    withAccent: {
      true: {
        root: "pl-3 before:absolute before:inset-y-0 before:left-0 before:w-1 before:rounded-l-lg"
      }
    },
    dismissible: {
      true: { root: "pr-10" }
    },
    shadow: {
      true: { root: "shadow-sm" }
    },
    bordered: {
      false: { root: "border-transparent" }
    }
  },
  defaultVariants: {
    size: "md",
    tone: "subtle",
    color: "blue",
    withAccent: true,
    shadow: false,
    bordered: true
  },
  compoundVariants: [
    // BLUE
    {
      tone: "subtle",
      color: "blue",
      class: {
        root: "bg-blue-50 dark:bg-blue-950/50 border-blue-200 dark:border-blue-800 before:bg-blue-500/60",
        icon: "text-blue-600 dark:text-blue-400",
        title: "text-blue-900 dark:text-blue-100",
        description: "text-blue-700 dark:text-blue-300"
      }
    },
    {
      tone: "solid",
      color: "blue",
      class: {
        root: "bg-blue-600 text-white dark:bg-blue-500 border-blue-700/60 dark:border-blue-600/60 before:bg-blue-400/90",
        icon: "text-white",
        title: "text-white",
        description: "text-blue-50"
      }
    },
    {
      tone: "outline",
      color: "blue",
      class: {
        root: "bg-transparent border-blue-300 dark:border-blue-800",
        icon: "text-blue-600 dark:text-blue-400",
        title: "text-blue-900 dark:text-blue-100",
        description: "text-blue-700 dark:text-blue-300"
      }
    },

    // GREEN
    {
      tone: "subtle",
      color: "green",
      class: {
        root: "bg-green-50 dark:bg-green-950/50 border-green-200 dark:border-green-800 before:bg-green-500/60",
        icon: "text-green-600 dark:text-green-400",
        title: "text-green-900 dark:text-green-100",
        description: "text-green-700 dark:text-green-300"
      }
    },
    {
      tone: "solid",
      color: "green",
      class: {
        root: "bg-green-600 text-white dark:bg-green-500 border-green-700/60 dark:border-green-600/60 before:bg-green-400/90",
        icon: "text-white",
        title: "text-white",
        description: "text-green-50"
      }
    },
    {
      tone: "outline",
      color: "green",
      class: {
        root: "bg-transparent border-green-300 dark:border-green-800",
        icon: "text-green-600 dark:text-green-400",
        title: "text-green-900 dark:text-green-100",
        description: "text-green-700 dark:text-green-300"
      }
    },

    // YELLOW
    {
      tone: "subtle",
      color: "yellow",
      class: {
        root: "bg-yellow-50 dark:bg-yellow-950/50 border-yellow-200 dark:border-yellow-800 before:bg-yellow-500/60",
        icon: "text-yellow-600 dark:text-yellow-400",
        title: "text-yellow-900 dark:text-yellow-100",
        description: "text-yellow-700 dark:text-yellow-300"
      }
    },
    {
      tone: "solid",
      color: "yellow",
      class: {
        root: "bg-yellow-500 text-yellow-950 dark:text-yellow-950 border-yellow-600/60 before:bg-yellow-400/90",
        icon: "text-yellow-950",
        title: "text-yellow-950",
        description: "text-yellow-950/80"
      }
    },
    {
      tone: "outline",
      color: "yellow",
      class: {
        root: "bg-transparent border-yellow-300 dark:border-yellow-800",
        icon: "text-yellow-600 dark:text-yellow-400",
        title: "text-yellow-900 dark:text-yellow-100",
        description: "text-yellow-700 dark:text-yellow-300"
      }
    },

    // RED
    {
      tone: "subtle",
      color: "red",
      class: {
        root: "bg-red-50 dark:bg-red-950/50 border-red-200 dark:border-red-800 before:bg-red-500/60",
        icon: "text-red-600 dark:text-red-400",
        title: "text-red-900 dark:text-red-100",
        description: "text-red-700 dark:text-red-300"
      }
    },
    {
      tone: "solid",
      color: "red",
      class: {
        root: "bg-red-600 text-white dark:bg-red-500 border-red-700/60 dark:border-red-600/60 before:bg-red-400/90",
        icon: "text-white",
        title: "text-white",
        description: "text-red-50"
      }
    },
    {
      tone: "outline",
      color: "red",
      class: {
        root: "bg-transparent border-red-300 dark:border-red-800",
        icon: "text-red-600 dark:text-red-400",
        title: "text-red-900 dark:text-red-100",
        description: "text-red-700 dark:text-red-300"
      }
    },

    // PURPLE
    {
      tone: "subtle",
      color: "purple",
      class: {
        root: "bg-purple-50 dark:bg-purple-950/50 border-purple-200 dark:border-purple-800 before:bg-purple-500/60",
        icon: "text-purple-600 dark:text-purple-400",
        title: "text-purple-900 dark:text-purple-100",
        description: "text-purple-700 dark:text-purple-300"
      }
    },
    {
      tone: "solid",
      color: "purple",
      class: {
        root: "bg-purple-600 text-white dark:bg-purple-500 border-purple-700/60 dark:border-purple-600/60 before:bg-purple-400/90",
        icon: "text-white",
        title: "text-white",
        description: "text-purple-50"
      }
    },
    {
      tone: "outline",
      color: "purple",
      class: {
        root: "bg-transparent border-purple-300 dark:border-purple-800",
        icon: "text-purple-600 dark:text-purple-400",
        title: "text-purple-900 dark:text-purple-100",
        description: "text-purple-700 dark:text-purple-300"
      }
    },

    // GRAY
    {
      tone: "subtle",
      color: "gray",
      class: {
        root: "bg-gray-50 dark:bg-gray-950/50 border-gray-200 dark:border-gray-800 before:bg-gray-400/60",
        icon: "text-gray-600 dark:text-gray-400",
        title: "text-gray-900 dark:text-gray-100",
        description: "text-gray-700 dark:text-gray-300"
      }
    },
    {
      tone: "solid",
      color: "gray",
      class: {
        root: "bg-gray-700 text-white dark:bg-gray-600 border-gray-800/60 before:bg-gray-500/90",
        icon: "text-white",
        title: "text-white",
        description: "text-gray-100"
      }
    },
    {
      tone: "outline",
      color: "gray",
      class: {
        root: "bg-transparent border-gray-300 dark:border-gray-800",
        icon: "text-gray-600 dark:text-gray-400",
        title: "text-gray-900 dark:text-gray-100",
        description: "text-gray-700 dark:text-gray-300"
      }
    }
  ]
})

export type CalloutStyleProps = VariantProps<typeof calloutStyles>

interface CalloutProps {
  /** Visual color (back-compat: alias of `variant`) */
  color?: CalloutStyleProps["color"]
  /** Old API alias for `color` */
  variant?: CalloutStyleProps["color"]
  /** Surface intensity */
  tone?: CalloutStyleProps["tone"]
  /** Size (spacing & typography) */
  size?: CalloutStyleProps["size"]
  /** Accent bar at left */
  withAccent?: boolean
  /** Add subtle shadow */
  shadow?: boolean
  /** Show border */
  bordered?: boolean
  /** Optional icon component */
  icon?: Component
  /** Heading */
  title?: string
  /** Supporting text */
  description?: string
  /** Render as */
  as?: string | Component
  /** Extra classes for the root */
  class?: any
  /** Accessible role */
  role?: "status" | "alert" | "note" | "region"
  /** When true, shows close button */
  dismissible?: boolean
}

// v-model for open state
const open = defineModel<boolean>("open", { default: true })

const props = withDefaults(defineProps<CalloutProps>(), {
  tone: "subtle",
  withAccent: true,
  shadow: false,
  bordered: true,
  role: "status",
  dismissible: false
})

const emit = defineEmits<{ (e: "dismiss"): void }>()

// Back-compat: map `variant` -> `color` if provided
const color = computed(() => props.color ?? props.variant ?? "blue")

const elementType = computed(() => props.as ?? "div")

const live = computed(() => (props.role === "alert" ? "assertive" : "polite"))

const classes = computed(() =>
  calloutStyles({
    size: props.size,
    tone: props.tone,
    color: color.value,
    withAccent: props.withAccent,
    shadow: props.shadow,
    bordered: props.bordered,
    dismissible: props.dismissible
  })
)

function close() {
  open.value = false
  emit("dismiss")
}
</script>

<template>
  <component :is="elementType" v-if="open" :class="classes.root({ class: props.class })" :role="role" :aria-live="live">
    <button v-if="dismissible" type="button" :class="classes.close()" @click="close" aria-label="Dismiss">
      <!-- simple ×, you can swap to an icon slot later -->
      <span aria-hidden="true">&times;</span>
    </button>

    <div :class="classes.body()">
      <!-- Icon: prop OR named slot -->
      <slot name="icon" :classes="classes.icon">
        <component v-if="icon" :is="icon" :class="classes.icon()" />
      </slot>

      <div :class="classes.content()">
        <!-- Title -->
        <slot name="title" :classes="classes.title">
          <p v-if="title" :class="classes.title()">{{ title }}</p>
        </slot>

        <!-- Description -->
        <slot name="description" :classes="classes.description">
          <p v-if="description" :class="classes.description()">{{ description }}</p>
        </slot>

        <!-- Default slot -->
        <slot :classes="classes" />

        <!-- Actions -->
        <div v-if="$slots.actions" :class="classes.actions()">
          <slot name="actions" />
        </div>
      </div>
    </div>
  </component>
</template>

