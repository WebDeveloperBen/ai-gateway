<script setup lang="ts">
const appConfig = useAppConfig()
</script>

<template>
  <div class="flex flex-col gap-6">
    <PageHeader title="Simple Test Page" :subtext="`Testing routing in ${appConfig.app.name}`">
    </PageHeader>

    <UiCard>
      <UiCardHeader>
        <UiCardTitle>Test Content</UiCardTitle>
        <UiCardDescription>This is a simple test to verify routing works</UiCardDescription>
      </UiCardHeader>
      <UiCardContent>
        <p>If you can see this page, the routing is working correctly!</p>
        <p>Path should be: /observability/simple-test</p>
      </UiCardContent>
    </UiCard>
  </div>
</template>