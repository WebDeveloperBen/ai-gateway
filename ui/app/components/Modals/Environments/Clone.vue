<script setup lang="ts">
import { Copy, Globe, Users, Key, AlertCircle } from "lucide-vue-next"

interface Environment {
  id: string
  name: string
  description: string
  owner: string
  teams: string[]
  status: "active" | "inactive"
  memberCount: number
  applicationCount: number
  monthlyRequests: number
  createdAt: string
  lastActivity: string
}

interface Props {
  open: boolean
  sourceEnvironment: Environment
}

const props = defineProps<Props>()
const emit = defineEmits<{
  "update:open": [value: boolean]
  cloned: [data: { name: string; description: string; includeTeams: boolean; includeKeys: boolean }]
}>()

const formData = reactive({
  name: `${props.sourceEnvironment?.name || ""} (Copy)`,
  description: props.sourceEnvironment?.description || "",
  includeTeams: true,
  includeKeys: false
})

const cloneOptions = [
  {
    key: "includeTeams",
    title: "Team Assignments",
    description: "Copy all team member assignments and roles",
    icon: Users,
    recommended: true
  },
  {
    key: "includeKeys",
    title: "API Keys Structure",
    description: "Create new API keys with same naming patterns (keys will be regenerated)",
    icon: Key,
    recommended: false
  }
]

function handleSubmit() {
  emit("cloned", {
    name: formData.name,
    description: formData.description,
    includeTeams: formData.includeTeams,
    includeKeys: formData.includeKeys
  })
  emit("update:open", false)
}
</script>

<template>
  <UiDialog :open="open" @update:open="emit('update:open', $event)">
    <UiDialogContent class="md:min-w-xl lg:min-w-3xl">
      <UiDialogHeader>
        <UiDialogTitle class="flex items-center gap-2">
          <Copy class="size-5" />
          Clone Environment
        </UiDialogTitle>
        <UiDialogDescription>
          Create a copy of "{{ sourceEnvironment?.name }}" with your chosen configurations
        </UiDialogDescription>
      </UiDialogHeader>

      <form @submit.prevent="handleSubmit" class="space-y-6">
        <!-- Source Environment Info -->
        <div class="rounded-lg border bg-muted/20 p-4">
          <div class="flex items-center gap-3 mb-3">
            <div class="p-2 rounded-md bg-primary/10">
              <Globe class="size-4 text-primary" />
            </div>
            <div>
              <h4 class="font-medium">Source: {{ sourceEnvironment?.name }}</h4>
              <p class="text-sm text-muted-foreground">{{ sourceEnvironment?.description }}</p>
            </div>
          </div>

          <div class="grid grid-cols-3 gap-4 text-sm">
            <div class="text-center">
              <div class="font-medium">{{ sourceEnvironment?.memberCount }}</div>
              <div class="text-muted-foreground">Team Members</div>
            </div>
            <div class="text-center">
              <div class="font-medium">{{ sourceEnvironment?.applicationCount }}</div>
              <div class="text-muted-foreground">Applications</div>
            </div>
            <div class="text-center">
              <div class="font-medium">{{ sourceEnvironment?.teams?.length }}</div>
              <div class="text-muted-foreground">Teams</div>
            </div>
          </div>
        </div>

        <!-- New Environment Details -->
        <div class="space-y-4">
          <div>
            <UiLabel for="name">Environment Name</UiLabel>
            <UiInput id="name" v-model="formData.name" placeholder="Enter environment name" class="mt-1" required />
          </div>

          <div>
            <UiLabel for="description">Description</UiLabel>
            <UiTextarea
              id="description"
              v-model="formData.description"
              placeholder="Describe the purpose of this environment"
              class="mt-1"
              :rows="3"
            />
          </div>
        </div>

        <!-- Clone Options -->
        <div class="space-y-4">
          <h4 class="font-medium">What to include in the clone:</h4>

          <div class="space-y-3">
            <div v-for="option in cloneOptions" :key="option.key" class="flex items-start gap-3 p-4 border rounded-lg">
              <UiCheckbox :id="option.key" v-model="formData[option.key]" class="mt-1" />
              <div class="flex-1">
                <div class="flex items-center gap-2 mb-1">
                  <component :is="option.icon" class="size-4" />
                  <label :for="option.key" class="font-medium cursor-pointer">
                    {{ option.title }}
                  </label>
                  <UiBadge v-if="option.recommended" variant="secondary" class="text-xs"> Recommended </UiBadge>
                </div>
                <p class="text-sm text-muted-foreground">{{ option.description }}</p>
              </div>
            </div>
          </div>

          <div
            class="flex items-start gap-2 p-3 bg-amber-50 dark:bg-amber-900/20 rounded-lg border border-amber-200 dark:border-amber-800"
          >
            <AlertCircle class="size-4 text-amber-600 dark:text-amber-400 mt-0.5 flex-shrink-0" />
            <div class="text-sm text-amber-800 dark:text-amber-200">
              <strong>Note:</strong> Applications and their data will not be cloned. You'll need to deploy applications
              to the new environment separately.
            </div>
          </div>
        </div>

        <UiDialogFooter>
          <UiButton type="button" variant="outline" @click="emit('update:open', false)"> Cancel </UiButton>
          <UiButton type="submit">
            <Copy class="size-4 mr-2" />
            Clone Environment
          </UiButton>
        </UiDialogFooter>
      </form>
    </UiDialogContent>
  </UiDialog>
</template>
