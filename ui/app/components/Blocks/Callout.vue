<script setup lang="ts">
import type { Component } from "vue"

interface CalloutProps {
  /** Callout variant/color theme */
  variant?: "blue" | "green" | "yellow" | "red" | "purple" | "gray"
  /** Icon component to display */
  icon?: Component
  /** Title text */
  title?: string
  /** Description text */
  description?: string
}

const props = withDefaults(defineProps<CalloutProps>(), {
  variant: "blue"
})

const variantClasses = computed(() => {
  const variants = {
    blue: {
      container: "bg-blue-50 dark:bg-blue-950/50 border-blue-200 dark:border-blue-800",
      icon: "text-blue-600 dark:text-blue-400",
      title: "text-blue-900 dark:text-blue-100",
      description: "text-blue-700 dark:text-blue-300"
    },
    green: {
      container: "bg-green-50 dark:bg-green-950/50 border-green-200 dark:border-green-800",
      icon: "text-green-600 dark:text-green-400",
      title: "text-green-900 dark:text-green-100",
      description: "text-green-700 dark:text-green-300"
    },
    yellow: {
      container: "bg-yellow-50 dark:bg-yellow-950/50 border-yellow-200 dark:border-yellow-800",
      icon: "text-yellow-600 dark:text-yellow-400",
      title: "text-yellow-900 dark:text-yellow-100",
      description: "text-yellow-700 dark:text-yellow-300"
    },
    red: {
      container: "bg-red-50 dark:bg-red-950/50 border-red-200 dark:border-red-800",
      icon: "text-red-600 dark:text-red-400",
      title: "text-red-900 dark:text-red-100",
      description: "text-red-700 dark:text-red-300"
    },
    purple: {
      container: "bg-purple-50 dark:bg-purple-950/50 border-purple-200 dark:border-purple-800",
      icon: "text-purple-600 dark:text-purple-400",
      title: "text-purple-900 dark:text-purple-100",
      description: "text-purple-700 dark:text-purple-300"
    },
    gray: {
      container: "bg-gray-50 dark:bg-gray-950/50 border-gray-200 dark:border-gray-800",
      icon: "text-gray-600 dark:text-gray-400",
      title: "text-gray-900 dark:text-gray-100",
      description: "text-gray-700 dark:text-gray-300"
    }
  }

  return variants[props.variant]
})
</script>

<template>
  <div class="border rounded-lg p-4" :class="variantClasses.container">
    <div class="flex gap-3">
      <!-- Icon -->
      <component v-if="icon" :is="icon" class="w-5 h-5 flex-shrink-0 mt-0.5" :class="variantClasses.icon" />
      <slot name="icon" :classes="variantClasses.icon" />

      <!-- Content -->
      <div class="text-sm flex-1">
        <!-- Title -->
        <p v-if="title" class="font-medium mb-1" :class="variantClasses.title">
          {{ title }}
        </p>
        <slot name="title" :classes="variantClasses.title" />

        <!-- Description -->
        <p v-if="description" :class="variantClasses.description">
          {{ description }}
        </p>
        <slot name="description" :classes="variantClasses.description" />

        <!-- Default slot for custom content -->
        <slot :classes="variantClasses" />
      </div>
    </div>
  </div>
</template>
