<script setup lang="ts">
import { Settings, Save } from "lucide-vue-next"

interface Application {
  id: string
  name: string
  description: string
  status: string
  team: string
}

interface Props {
  open: boolean
  application: Application
  loading?: boolean
}

const props = withDefaults(defineProps<Props>(), {
  loading: false
})

const emit = defineEmits<{
  "update:open": [value: boolean]
  save: [application: Application]
  cancel: []
}>()

const form = ref({
  name: props.application.name,
  description: props.application.description,
  team: props.application.team,
  status: props.application.status
})

watch(() => props.application, (newApp) => {
  form.value = {
    name: newApp.name,
    description: newApp.description,
    team: newApp.team,
    status: newApp.status
  }
}, { deep: true })

const handleSave = () => {
  emit("save", {
    ...props.application,
    ...form.value
  })
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
    <UiDialogContent class="sm:max-w-2xl">
      <UiDialogHeader>
        <UiDialogTitle class="flex items-center gap-2">
          <Settings class="h-5 w-5" />
          Application Settings
        </UiDialogTitle>
        <UiDialogDescription>
          Update the configuration and details for this application.
        </UiDialogDescription>
      </UiDialogHeader>

      <div class="space-y-6 py-4">
        <div class="space-y-4">
          <div class="grid grid-cols-2 gap-4">
            <div class="space-y-2">
              <label class="text-sm font-medium">Application Name</label>
              <UiInput
                v-model="form.name"
                placeholder="Enter application name"
                :disabled="props.loading"
              />
            </div>

            <div class="space-y-2">
              <label class="text-sm font-medium">Team</label>
              <UiInput
                v-model="form.team"
                placeholder="Enter team name"
                :disabled="props.loading"
              />
            </div>
          </div>

          <div class="space-y-2">
            <label class="text-sm font-medium">Description</label>
            <UiTextarea
              v-model="form.description"
              placeholder="Enter application description"
              rows="3"
              :disabled="props.loading"
            />
          </div>

          <div class="space-y-2">
            <label class="text-sm font-medium">Status</label>
            <select
              v-model="form.status"
              :disabled="props.loading"
              class="flex h-10 w-full rounded-md border border-input bg-background px-3 py-2 text-sm ring-offset-background file:border-0 file:bg-transparent file:text-sm file:font-medium placeholder:text-muted-foreground focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-ring focus-visible:ring-offset-2 disabled:cursor-not-allowed disabled:opacity-50"
            >
              <option value="active">Active</option>
              <option value="inactive">Inactive</option>
              <option value="maintenance">Maintenance</option>
            </select>
          </div>
        </div>
      </div>

      <UiDialogFooter>
        <UiButton variant="outline" @click="handleCancel" :disabled="props.loading">
          Cancel
        </UiButton>
        <UiButton @click="handleSave" :disabled="props.loading">
          <span v-if="!props.loading" class="flex items-center gap-2">
            <Save class="h-4 w-4" />
            Save Changes
          </span>
          <span v-else class="flex items-center gap-2">
            <div class="animate-spin rounded-full h-4 w-4 border-b-2 border-current"></div>
            Saving...
          </span>
        </UiButton>
      </UiDialogFooter>
    </UiDialogContent>
  </UiDialog>
</template>