<script setup lang="ts">
import { Plus, Shield, Search, Server } from "lucide-vue-next"
import { toast } from "vue-sonner"

interface Application {
  id: string
  name: string
  description: string
  team: string
  environment: string
  status: "active" | "inactive"
  hasPolicy?: boolean
}

const props = defineProps<{
  open: boolean
  policyId: string
  policyName: string
}>()

const emit = defineEmits<{
  "update:open": [value: boolean]
  assigned: [applicationId: string]
}>()

// Mock applications data - replace with real API call
const allApplications: Application[] = [
  {
    id: "app_3",
    name: "Email Assistant",
    description: "AI-powered email drafting and response assistant",
    team: "Customer Success",
    environment: "production",
    status: "active",
    hasPolicy: false
  },
  {
    id: "app_4",
    name: "Data Insights Bot",
    description: "Automated data analysis and reporting tool",
    team: "Analytics",
    environment: "production",
    status: "active",
    hasPolicy: false
  },
  {
    id: "app_6",
    name: "Marketing Copy Generator",
    description: "Content generation for marketing campaigns",
    team: "Marketing",
    environment: "staging",
    status: "active",
    hasPolicy: true
  },
  {
    id: "app_7",
    name: "Code Review Assistant",
    description: "AI-powered code review and suggestions",
    team: "Engineering",
    environment: "development",
    status: "active",
    hasPolicy: false
  },
  {
    id: "app_8",
    name: "Sales Chatbot",
    description: "Lead qualification and initial customer engagement",
    team: "Sales",
    environment: "production",
    status: "inactive",
    hasPolicy: false
  }
]

const searchQuery = ref("")
const selectedApplicationIds = ref<string[]>([])
const isAssigning = ref(false)

// Dialog state
const dialogOpen = computed({
  get: () => props.open,
  set: (value) => emit("update:open", value)
})

// Filter applications that don't already have this policy applied
const availableApplications = computed(() => {
  return allApplications.filter((app) => {
    const matchesSearch =
      searchQuery.value === "" ||
      app.name.toLowerCase().includes(searchQuery.value.toLowerCase()) ||
      app.team.toLowerCase().includes(searchQuery.value.toLowerCase()) ||
      app.description.toLowerCase().includes(searchQuery.value.toLowerCase())

    // Only show applications that don't already have a policy or have an inactive policy
    const isAvailable = !app.hasPolicy

    return matchesSearch && isAvailable
  })
})

const canAssign = computed(() => {
  return selectedApplicationIds.value.length > 0 && !isAssigning.value
})

const onAssign = async () => {
  if (!canAssign.value) return

  isAssigning.value = true

  try {
    // TODO: Implement API call to assign policy to applications
    console.log(`Assigning policy ${props.policyId} to applications:`, selectedApplicationIds.value)

    // Simulate API delay
    await new Promise((resolve) => setTimeout(resolve, 1000))

    const appNames = selectedApplicationIds.value
      .map((id) => allApplications.find((app) => app.id === id)?.name)
      .filter(Boolean)
      .join(", ")

    // Show success toast
    toast({
      title: "Policy assigned successfully",
      description: `${props.policyName} has been applied to ${appNames}.`,
      duration: 5000,
      icon: "lucide:check"
    })

    // Close modal and reset
    closeDialog(true)

    // Emit assigned events for each application
    selectedApplicationIds.value.forEach((appId) => {
      emit("assigned", appId)
    })
  } catch (error) {
    console.error("Failed to assign policy:", error)
    toast({
      title: "Assignment failed",
      description: "Failed to assign policy to applications. Please try again.",
      duration: 5000,
      icon: "lucide:x",
      variant: "destructive"
    })
  } finally {
    isAssigning.value = false
  }
}

const closeDialog = (save: boolean) => {
  if (!save && selectedApplicationIds.value.length > 0) {
    toast({
      title: "Assignment cancelled",
      description: "Policy assignment has been cancelled.",
      duration: 3000,
      icon: "lucide:x"
    })
  }
  dialogOpen.value = false
  selectedApplicationIds.value = []
  searchQuery.value = ""
}

const toggleApplicationSelection = (appId: string) => {
  const index = selectedApplicationIds.value.indexOf(appId)
  if (index > -1) {
    selectedApplicationIds.value.splice(index, 1)
  } else {
    selectedApplicationIds.value.push(appId)
  }
}

const getEnvironmentColor = (environment: string) => {
  switch (environment) {
    case "production":
      return "text-red-700 dark:text-red-400 bg-red-50 dark:bg-red-950/50 border-red-200 dark:border-red-800"
    case "staging":
      return "text-yellow-700 dark:text-yellow-400 bg-yellow-50 dark:bg-yellow-950/50 border-yellow-200 dark:border-yellow-800"
    case "development":
      return "text-blue-700 dark:text-blue-400 bg-blue-50 dark:bg-blue-950/50 border-blue-200 dark:border-blue-800"
    default:
      return "text-muted-foreground bg-muted/50 border-border"
  }
}

// Watch for modal close to reset state
watch(dialogOpen, (isOpen) => {
  if (!isOpen) {
    selectedApplicationIds.value = []
    searchQuery.value = ""
  }
})
</script>

<template>
  <UiDialog v-model:open="dialogOpen">
    <UiDialogContent class="sm:max-w-3xl max-h-[90vh] overflow-y-auto">
      <template #header>
        <div class="flex items-center gap-3 pb-4">
          <div
            class="w-10 h-10 rounded-full bg-gradient-to-br from-blue-500 to-purple-600 flex items-center justify-center shadow-sm"
          >
            <Shield class="w-5 h-5 text-white" />
          </div>
          <div>
            <h2 class="text-lg font-semibold">Assign Policy to Applications</h2>
            <p class="text-sm text-muted-foreground">
              Select applications to apply the <strong>{{ policyName }}</strong> policy
            </p>
          </div>
        </div>
      </template>

      <template #content>
        <div class="space-y-6">
          <!-- Search Input -->
          <div class="relative">
            <Search class="absolute left-3 top-1/2 transform -translate-y-1/2 text-muted-foreground h-4 w-4" />
            <UiInput
              v-model="searchQuery"
              placeholder="Search applications by name, team, or description..."
              class="pl-10"
            />
          </div>

          <!-- Selection Summary -->
          <div v-if="selectedApplicationIds.length > 0" class="bg-muted/30 border rounded-lg p-4">
            <p class="text-sm font-medium mb-2">Selected Applications ({{ selectedApplicationIds.length }})</p>
            <div class="flex flex-wrap gap-2">
              <UiBadge
                v-for="appId in selectedApplicationIds"
                :key="appId"
                variant="secondary"
                class="text-xs cursor-pointer hover:bg-muted"
                @click="toggleApplicationSelection(appId)"
              >
                {{ allApplications.find((app) => app.id === appId)?.name }}
                ×
              </UiBadge>
            </div>
          </div>

          <!-- Available Applications List -->
          <div class="space-y-4 max-h-96 overflow-y-auto p-1">
            <div v-if="availableApplications.length === 0" class="text-center py-8 text-muted-foreground">
              <Server class="mx-auto h-12 w-12 mb-4 opacity-50" />
              <p class="font-medium mb-1">No available applications</p>
              <p class="text-sm">
                {{
                  searchQuery
                    ? "No applications match your search."
                    : "All applications already have policies assigned."
                }}
              </p>
            </div>

            <UiSelectableCard
              v-for="app in availableApplications"
              :key="app.id"
              :selected="selectedApplicationIds.includes(app.id)"
              :show-checkbox="true"
              @click="toggleApplicationSelection(app.id)"
            >
              <template #header>
                <div class="flex items-center gap-3">
                  <!-- App Icon -->
                  <div class="p-2 rounded-lg bg-primary/10">
                    <Server class="h-5 w-5 text-primary" />
                  </div>

                  <!-- App Name and Environment -->
                  <div class="flex-1">
                    <div class="flex items-center gap-3">
                      <h4 class="font-semibold text-lg text-foreground group-hover:text-primary transition-colors">
                        {{ app.name }}
                      </h4>
                      <div
                        class="inline-flex items-center gap-1 px-2 py-1 rounded-md text-xs font-medium border"
                        :class="getEnvironmentColor(app.environment)"
                      >
                        {{ app.environment }}
                      </div>
                    </div>
                  </div>
                </div>
              </template>

              <template #content>
                <p class="text-muted-foreground text-sm leading-relaxed">
                  {{ app.description }}
                </p>

                <div class="flex items-center gap-4 text-xs text-muted-foreground">
                  <span>{{ app.team }}</span>
                  <span>•</span>
                  <span class="capitalize">{{ app.status }}</span>
                </div>
              </template>
            </UiSelectableCard>
          </div>

          <BlocksCallout
            variant="blue"
            title="Policy Assignment"
            description="
            Selected applications will immediately start using this policy. You can modify or remove policy
            assignments at any time.
            "
            :icon="Shield"
          />
        </div>
      </template>

      <template #footer>
        <UiDialogFooter>
          <UiButton variant="outline" type="button" class="mt-2 sm:mt-0" @click="closeDialog(false)"> Cancel </UiButton>
          <UiButton type="submit" class="gap-2" :disabled="!canAssign" @click="onAssign">
            <Plus class="w-4 h-4" />
            {{
              isAssigning
                ? "Assigning..."
                : `Assign to ${selectedApplicationIds.length} Application${selectedApplicationIds.length === 1 ? "" : "s"}`
            }}
          </UiButton>
        </UiDialogFooter>
      </template>
    </UiDialogContent>
  </UiDialog>
</template>
