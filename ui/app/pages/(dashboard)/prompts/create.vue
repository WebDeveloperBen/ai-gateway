<script setup lang="ts">
useSeoMeta({ title: "Create Prompt - LLM Gateway" })

const handleSave = async (promptData: any) => {
  try {
    // TODO: Replace with actual API call
    console.log("Saving new prompt:", promptData)
    // await $fetch('/api/prompts', { method: 'POST', body: promptData })
    
    // Navigate to the new prompt page
    // const newPromptId = result.id
    // navigateTo(`/prompts/${newPromptId}`)
    
    // For now, navigate to prompts list
    navigateTo('/prompts')
    
    // TODO: Show success toast
  } catch (error) {
    console.error("Failed to save prompt:", error)
    // TODO: Show error toast
  }
}

const handleCancel = () => {
  navigateTo('/prompts')
}
</script>

<template>
  <PromptsCreateEditor 
    @save="handleSave"
    @cancel="handleCancel"
  />
</template>