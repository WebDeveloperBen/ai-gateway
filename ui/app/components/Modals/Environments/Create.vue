<script setup lang="ts">
import { Globe, Building, Users } from "lucide-vue-next"
import { useForm } from "vee-validate"
import type { FormBuilder } from "@/components/Ui/FormBuilder/FormBuilder.vue"

interface Props {
  open: boolean
}

// Sample teams for selection
const availableTeams = [
  { id: "1", name: "Engineering" },
  { id: "2", name: "Product" },
  { id: "3", name: "Marketing" },
  { id: "4", name: "Customer Success" },
  { id: "5", name: "DevOps" },
  { id: "6", name: "QA" },
  { id: "7", name: "Sales" },
  { id: "8", name: "Analytics" }
]

// Sample users for owner selection
const availableUsers = [
  { id: "1", name: "Alice Johnson", email: "alice@company.com" },
  { id: "2", name: "Bob Smith", email: "bob@company.com" },
  { id: "3", name: "Carol Williams", email: "carol@company.com" },
  { id: "4", name: "David Brown", email: "david@company.com" },
  { id: "5", name: "Emma Davis", email: "emma@company.com" }
]

const props = defineProps<Props>()
const emit = defineEmits<{
  "update:open": [value: boolean]
  created: [environmentData: any]
}>()

const isOpen = computed({
  get: () => props.open,
  set: (value) => emit("update:open", value)
})

// Options for form fields
const ownerOptions = availableUsers.map((user) => ({
  value: user.id,
  label: `${user.name} (${user.email})`
}))

// FormBuilder field definitions - simplified for user flexibility
const formFields: FormBuilder[] = [
  {
    variant: "Input",
    name: "name",
    label: "Environment Name",
    placeholder: "e.g., Production, Staging, Development, QA, Demo",
    hint: "Choose any name that fits your organization's needs",
    required: true
  },
  {
    variant: "Textarea",
    name: "description",
    label: "Description",
    placeholder: "Brief description of the environment's purpose and usage",
    hint: "What is this environment used for?",
    required: true,
    rows: 2
  },
  {
    variant: "Select",
    name: "owner",
    label: "Environment Owner",
    hint: "The person responsible for managing this environment",
    options: ownerOptions,
    required: true
  }
]

// Form setup using vee-validate
const { handleSubmit, resetForm, values, isSubmitting } = useForm<{
  name: string
  description: string
  owner: string
  teams: string[]
}>({
  initialValues: {
    name: "",
    description: "",
    owner: "",
    teams: []
  }
})

// Team selection
const selectedTeams = ref<string[]>([])

const isLoading = shallowRef(false)

const onSubmit = handleSubmit(async (formData) => {
  isLoading.value = true
  try {
    const owner = availableUsers.find((u) => u.id === formData.owner)
    const selectedTeamNames = selectedTeams.value
      .map((teamId) => availableTeams.find((t) => t.id === teamId)?.name)
      .filter(Boolean)

    if (!owner) {
      throw new Error("Invalid owner selection")
    }

    const environmentData = {
      id: `env_${Date.now()}`,
      name: formData.name,
      description: formData.description,
      owner: owner.name,
      teams: selectedTeamNames,
      status: "active" as const,
      memberCount: selectedTeamNames.length * 3, // Estimate
      applicationCount: 0,
      monthlyRequests: 0,
      createdAt: new Date().toISOString(),
      lastActivity: new Date().toISOString()
    }

    console.log("Creating environment:", environmentData)

    emit("created", environmentData)
    handleClose()
  } catch (error) {
    console.error("Failed to create environment:", error)
  } finally {
    isLoading.value = false
  }
})

const handleFormSubmit = () => {
  onSubmit()
}

function handleClose() {
  isOpen.value = false
  setTimeout(() => {
    resetForm()
    selectedTeams.value = []
  }, 150)
}

function addTeam(teamId: string) {
  if (!selectedTeams.value.includes(teamId)) {
    selectedTeams.value.push(teamId)
  }
}

function removeTeam(teamId: string) {
  selectedTeams.value = selectedTeams.value.filter((id) => id !== teamId)
}

// Removed type/isolation icon functions - simplified approach
</script>

<template>
  <UiDialog v-model:open="isOpen">
    <UiDialogContent class="sm:max-w-2xl">
      <template #header>
        <UiDialogTitle class="sr-only">Create Environment</UiDialogTitle>
        <div class="flex items-center gap-3 pb-4">
          <div
            class="w-10 h-10 rounded-full bg-gradient-to-br from-primary to-primary/80 flex items-center justify-center shadow-sm"
          >
            <Globe class="w-5 h-5 text-primary-foreground" />
          </div>
          <div>
            <h2 class="text-lg font-semibold">Create Environment</h2>
            <p class="text-sm text-muted-foreground">Set up a new deployment environment with team access controls</p>
          </div>
        </div>
      </template>

      <template #content>
        <form @submit.prevent="handleFormSubmit" class="space-y-6">
          <!-- Form Fields -->
          <UiFormBuilder :fields="formFields" />

          <!-- Team Assignment -->
          <div class="space-y-3">
            <div class="flex items-center justify-between">
              <div>
                <label class="text-sm font-medium">Teams</label>
                <p class="text-xs text-muted-foreground">Select teams that will have access to this environment</p>
              </div>
            </div>

            <!-- Team Selection -->
            <div class="grid grid-cols-2 gap-2">
              <div
                v-for="team in availableTeams"
                :key="team.id"
                class="flex items-center justify-between p-3 border rounded-lg cursor-pointer transition-colors"
                :class="selectedTeams.includes(team.id) ? 'bg-primary/5 border-primary/30' : 'hover:bg-muted/50'"
                @click="selectedTeams.includes(team.id) ? removeTeam(team.id) : addTeam(team.id)"
              >
                <div class="flex items-center gap-2">
                  <Building class="size-4 text-muted-foreground" />
                  <span class="text-sm font-medium">{{ team.name }}</span>
                </div>
                <UiCheckbox
                  :model-value="selectedTeams.includes(team.id)"
                  @click.stop
                  @update:model-value="(checked: any) => (checked ? addTeam(team.id) : removeTeam(team.id))"
                />
              </div>
            </div>
          </div>

          <!-- Environment Preview -->
          <div v-if="values.name" class="p-4 bg-muted/20 rounded-lg border">
            <h3 class="text-sm font-medium mb-3">Environment Preview</h3>

            <div class="space-y-3">
              <div class="flex items-center gap-3">
                <Globe class="size-5 text-primary" />
                <div>
                  <div class="font-medium text-sm">{{ values.name }}</div>
                  <div class="text-xs text-muted-foreground">{{ values.description || "Custom Environment" }}</div>
                </div>
              </div>

              <div v-if="selectedTeams.length > 0" class="flex items-center gap-2 text-xs text-muted-foreground">
                <Users class="size-4" />
                <span>{{ selectedTeams.length }} team{{ selectedTeams.length === 1 ? "" : "s" }} assigned</span>
              </div>
            </div>
          </div>

          <div class="flex justify-end space-x-2">
            <UiButton type="button" variant="outline" @click="handleClose" :disabled="isSubmitting"> Cancel </UiButton>
            <UiButton type="button" @click="handleFormSubmit" :loading="isSubmitting" class="gap-2">
              <Globe class="h-4 w-4" />
              Create Environment
            </UiButton>
          </div>
        </form>
      </template>
    </UiDialogContent>
  </UiDialog>
</template>
