<script setup lang="ts">
import { AlertCircle } from "lucide-vue-next"

interface Props {
  open: boolean
  pendingPromptId: string
}

interface Emits {
  "update:open": [value: boolean]
  confirmReplace: []
  cancelReplace: []
}

const props = defineProps<Props>()
const emit = defineEmits<Emits>()

const isOpen = computed({
  get: () => props.open,
  set: (value) => emit("update:open", value)
})

function confirmReplace() {
  emit("confirmReplace")
}

function cancelReplace() {
  emit("cancelReplace")
}
</script>

<template>
  <UiDialog v-model:open="isOpen">
    <UiDialogContent class="min-w-md">
      <template #header>
        <div class="flex items-center gap-3">
          <div class="flex h-10 w-10 items-center justify-center rounded-full bg-orange-100 dark:bg-orange-900/20">
            <AlertCircle class="h-5 w-5 text-orange-600 dark:text-orange-400" />
          </div>
          <div>
            <h3 class="font-semibold">Replace Current Content?</h3>
            <p class="text-sm text-muted-foreground">You have unsaved changes in the prompt editor.</p>
          </div>
        </div>
      </template>

      <template #content>
        <p class="text-sm text-muted-foreground">
          {{ pendingPromptId ? "Loading a saved prompt" : "Creating a new prompt" }} will replace your current content.
          This action cannot be undone.
        </p>
      </template>

      <template #footer>
        <UiDialogFooter class="flex gap-3 justify-end">
          <UiButton variant="outline" @click="cancelReplace"> Cancel </UiButton>
          <UiButton variant="destructive" @click="confirmReplace"> Replace Content </UiButton>
        </UiDialogFooter>
      </template>
    </UiDialogContent>
  </UiDialog>
</template>

