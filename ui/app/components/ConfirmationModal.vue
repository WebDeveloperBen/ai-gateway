<script setup lang="ts">
import { AlertTriangle } from "lucide-vue-next"

interface Props {
  open: boolean
  title: string
  description?: string
  confirmText?: string
  cancelText?: string
  variant?: "default" | "destructive"
  loading?: boolean
}

const props = withDefaults(defineProps<Props>(), {
  confirmText: "Confirm",
  cancelText: "Cancel",
  variant: "default",
  loading: false
})

const emit = defineEmits<{
  confirm: []
  cancel: []
  "update:open": [value: boolean]
}>()

const handleConfirm = () => {
  emit("confirm")
}

const handleCancel = () => {
  emit("cancel")
  emit("update:open", false)
}

const handleOpenChange = (open: boolean) => {
  emit("update:open", open)
  if (!open) {
    emit("cancel")
  }
}
</script>

<template>
  <UiDialog :open="props.open" @update:open="handleOpenChange">
    <UiDialogContent class="sm:max-w-md">
      <UiDialogHeader>
        <div class="space-y-3">
          <UiDialogTitle class="text-left flex items-center gap-2">
            <AlertTriangle
              class="h-4 w-4"
              :class="
                props.variant === 'destructive' ? 'text-red-600 dark:text-red-400' : 'text-blue-600 dark:text-blue-400'
              "
            />
            {{ props.title }}
          </UiDialogTitle>
          <UiDialogDescription v-if="props.description" class="text-left">
            {{ props.description }}
          </UiDialogDescription>
        </div>
      </UiDialogHeader>

      <UiDialogFooter class="flex-col-reverse sm:flex-row gap-2">
        <UiButton variant="outline" @click="handleCancel" :disabled="props.loading">
          {{ props.cancelText }}
        </UiButton>
        <UiButton
          :variant="props.variant === 'destructive' ? 'destructive' : 'default'"
          @click="handleConfirm"
          :disabled="props.loading"
        >
          <span v-if="!props.loading">{{ props.confirmText }}</span>
          <span v-else class="flex items-center gap-2">
            <div class="animate-spin rounded-full h-4 w-4 border-b-2 border-current"></div>
            Processing...
          </span>
        </UiButton>
      </UiDialogFooter>
    </UiDialogContent>
  </UiDialog>
</template>

