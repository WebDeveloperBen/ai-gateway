<script lang="ts" setup>
const props = withDefaults(
  defineProps<{
    statusCode?: number
    fatal?: boolean
    unhandled?: boolean
    statusMessage?: string
    message?: string
    data?: unknown
    cause?: unknown
  }>(),
  {
    statusCode: 404,
    fatal: false,
    unhandled: false,
    statusMessage: "",
    message: "We can't find this page",
    data: undefined,
    cause: undefined
  }
)

const title = computed(() => {
  if (!props.message) return "Error"
  return props.message
})

useSeoMeta({ title })
</script>
<template>
  <div
    class="relative flex min-h-[100dvh] flex-col md:min-h-screen md:flex-row bg-background pt-[env(safe-area-inset-top)] pb-[env(safe-area-inset-bottom)]"
  >
    <!-- Image: on top for mobile, side-by-side from md+ -->
    <div class="order-1 w-full md:order-2 md:w-1/2">
      <img src="/error-background.avif" alt="error page" class="h-48 w-full object-cover md:h-full" />
    </div>

    <!-- Content -->
    <div class="order-2 flex w-full items-center justify-center px-4 py-10 md:order-1 md:w-1/2 md:px-8 md:py-0">
      <div class="max-w-md text-center md:text-left">
        <p class="mb-3 text-sm font-semibold tracking-tight text-primary">{{ statusCode }} error</p>
        <h1 class="text-3xl font-bold tracking-tight sm:text-4xl lg:text-5xl">
          {{ title }}
        </h1>
        <UiButton class="mt-6 w-full sm:w-auto" @click="clearError({ redirect: '#' })"> Take me home </UiButton>
      </div>
    </div>
  </div>
</template>
