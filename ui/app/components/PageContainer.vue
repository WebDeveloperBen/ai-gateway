<template>
  <div class="p-10 h-full bg-gradient-to-br from-background via-muted/30 to-muted relative overflow-hidden">
    <div
      class="absolute inset-0 bg-gradient-to-tr from-primary/[0.03] via-transparent to-accent/[0.03] pointer-events-none"
    ></div>
    <slot />
  </div>
</template>
