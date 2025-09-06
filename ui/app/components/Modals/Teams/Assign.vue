<script setup lang="ts">
import { Users, Building, Plus, X, Check, UserPlus } from "lucide-vue-next"

interface Environment {
  id: string
  name: string
  description: string
  teams: string[]
}

interface Team {
  id: string
  name: string
  description: string
  memberCount: number
  owner: string
}

interface Props {
  open: boolean
  environment: Environment | null
}

// Sample teams data
const availableTeams: Team[] = [
  { id: "1", name: "Engineering", description: "Core engineering team", memberCount: 12, owner: "Alice Johnson" },
  { id: "2", name: "Product", description: "Product management team", memberCount: 8, owner: "Bob Smith" },
  { id: "3", name: "Marketing", description: "Marketing and content team", memberCount: 6, owner: "Carol Williams" },
  { id: "4", name: "Customer Success", description: "Customer support team", memberCount: 10, owner: "David Brown" },
  { id: "5", name: "DevOps", description: "DevOps and infrastructure", memberCount: 5, owner: "Emma Davis" },
  { id: "6", name: "QA", description: "Quality assurance team", memberCount: 7, owner: "Frank Wilson" },
  { id: "7", name: "Sales", description: "Sales and business development", memberCount: 9, owner: "Grace Lee" },
  { id: "8", name: "Analytics", description: "Data analysis team", memberCount: 4, owner: "Henry Chen" }
]

const props = defineProps<Props>()
const emit = defineEmits<{
  "update:open": [value: boolean]
  assigned: [data: { environmentId: string; teamIds: string[] }]
}>()

const isOpen = computed({
  get: () => props.open,
  set: (value) => emit("update:open", value)
})

// Track currently assigned teams
const assignedTeams = ref<string[]>([])
const searchQuery = ref("")
const isLoading = ref(false)

// Watch for environment changes to update assigned teams
watch(
  () => props.environment,
  (newEnv) => {
    if (newEnv) {
      // Map team names back to IDs for the current environment
      assignedTeams.value = availableTeams.filter((team) => newEnv.teams.includes(team.name)).map((team) => team.id)
    }
  },
  { immediate: true }
)

// Filter teams based on search
const filteredTeams = computed(() => {
  if (!searchQuery.value) return availableTeams

  return availableTeams.filter(
    (team) =>
      team.name.toLowerCase().includes(searchQuery.value.toLowerCase()) ||
      team.description.toLowerCase().includes(searchQuery.value.toLowerCase()) ||
      team.owner.toLowerCase().includes(searchQuery.value.toLowerCase())
  )
})

// Separate assigned and unassigned teams
const currentlyAssignedTeams = computed(() =>
  filteredTeams.value.filter((team) => assignedTeams.value.includes(team.id))
)

const availableTeamsToAssign = computed(() =>
  filteredTeams.value.filter((team) => !assignedTeams.value.includes(team.id))
)

function toggleTeamAssignment(teamId: string) {
  if (assignedTeams.value.includes(teamId)) {
    assignedTeams.value = assignedTeams.value.filter((id) => id !== teamId)
  } else {
    assignedTeams.value.push(teamId)
  }
}

async function handleSave() {
  if (!props.environment) return

  isLoading.value = true
  try {
    // Here you would make an API call to update team assignments
    console.log("Updating team assignments:", {
      environmentId: props.environment.id,
      teamIds: assignedTeams.value
    })

    emit("assigned", {
      environmentId: props.environment.id,
      teamIds: assignedTeams.value
    })

    handleClose()
  } catch (error) {
    console.error("Failed to assign teams:", error)
  } finally {
    isLoading.value = false
  }
}

function handleClose() {
  isOpen.value = false
  setTimeout(() => {
    searchQuery.value = ""
  }, 150)
}
</script>

<template>
  <UiDialog v-model:open="isOpen">
    <UiDialogContent class="sm:max-w-2xl max-h-[80vh]">
      <template #header>
        <UiDialogTitle class="sr-only">Assign Teams to Environment</UiDialogTitle>
        <div class="flex items-center gap-3 pb-4">
          <div
            class="w-10 h-10 rounded-full bg-gradient-to-br from-primary to-primary/80 flex items-center justify-center shadow-sm"
          >
            <UserPlus class="w-5 h-5 text-primary-foreground" />
          </div>
          <div class="flex-1">
            <h2 class="text-lg font-semibold">Manage Environment Access</h2>
            <p class="text-sm text-muted-foreground">
              {{ environment?.name ? `Configure team access for ${environment.name}` : "Loading..." }}
            </p>
          </div>
        </div>
      </template>

      <template #content>
        <div class="space-y-4">
          <!-- Search -->
          <div class="relative">
            <UiInput v-model="searchQuery" placeholder="Search teams by name, description, or owner..." class="pl-10" />
            <Users class="absolute left-3 top-1/2 transform -translate-y-1/2 text-muted-foreground size-4" />
          </div>

          <!-- Summary -->
          <div class="p-3 bg-muted/20 rounded-lg">
            <div class="flex items-center justify-between text-sm">
              <span class="text-muted-foreground">
                {{ assignedTeams.length }} of {{ availableTeams.length }} teams assigned
              </span>
              <span class="font-medium">
                ~{{
                  availableTeams.filter((t) => assignedTeams.includes(t.id)).reduce((sum, t) => sum + t.memberCount, 0)
                }}
                members will have access
              </span>
            </div>
          </div>

          <div class="flex-1 overflow-hidden">
            <!-- Currently Assigned Teams -->
            <div v-if="currentlyAssignedTeams.length > 0" class="mb-6">
              <h3 class="text-sm font-medium mb-3 text-muted-foreground">
                Assigned Teams ({{ currentlyAssignedTeams.length }})
              </h3>
              <div class="space-y-2 max-h-40 overflow-y-auto">
                <div
                  v-for="team in currentlyAssignedTeams"
                  :key="team.id"
                  class="flex items-center justify-between p-3 bg-emerald-50 dark:bg-emerald-950/30 border border-emerald-200 dark:border-emerald-800/50 rounded-lg"
                >
                  <div class="flex items-center gap-3">
                    <div class="p-2 rounded-md bg-emerald-100 dark:bg-emerald-900/50">
                      <Building class="size-4 text-emerald-600 dark:text-emerald-400" />
                    </div>
                    <div>
                      <div class="font-medium text-sm">{{ team.name }}</div>
                      <div class="text-xs text-muted-foreground">{{ team.memberCount }} members • {{ team.owner }}</div>
                    </div>
                  </div>
                  <UiButton
                    variant="outline"
                    size="sm"
                    @click="toggleTeamAssignment(team.id)"
                    class="text-destructive hover:text-destructive"
                  >
                    <X class="size-4" />
                  </UiButton>
                </div>
              </div>
            </div>

            <!-- Available Teams -->
            <div v-if="availableTeamsToAssign.length > 0">
              <h3 class="text-sm font-medium mb-3 text-muted-foreground">
                Available Teams ({{ availableTeamsToAssign.length }})
              </h3>
              <div class="space-y-2 max-h-60 overflow-y-auto">
                <div
                  v-for="team in availableTeamsToAssign"
                  :key="team.id"
                  class="flex items-center justify-between p-3 border rounded-lg hover:bg-muted/50 cursor-pointer transition-colors"
                  @click="toggleTeamAssignment(team.id)"
                >
                  <div class="flex items-center gap-3">
                    <div class="p-2 rounded-md bg-muted">
                      <Building class="size-4 text-muted-foreground" />
                    </div>
                    <div>
                      <div class="font-medium text-sm">{{ team.name }}</div>
                      <div class="text-xs text-muted-foreground">
                        {{ team.description }} • {{ team.memberCount }} members
                      </div>
                      <div class="text-xs text-muted-foreground">Owner: {{ team.owner }}</div>
                    </div>
                  </div>
                  <UiButton variant="outline" size="sm" @click.stop="toggleTeamAssignment(team.id)">
                    <Plus class="size-4" />
                  </UiButton>
                </div>
              </div>
            </div>

            <!-- No teams found -->
            <div v-if="filteredTeams.length === 0" class="text-center py-8">
              <Users class="size-12 text-muted-foreground mx-auto mb-4" />
              <p class="text-muted-foreground">No teams found matching your search.</p>
            </div>

            <!-- All teams assigned -->
            <div
              v-if="availableTeamsToAssign.length === 0 && assignedTeams.length > 0 && searchQuery === ''"
              class="text-center py-4"
            >
              <Check class="size-8 text-emerald-500 mx-auto mb-2" />
              <p class="text-sm text-muted-foreground">All teams have been assigned to this environment.</p>
            </div>
          </div>

          <!-- Actions -->
          <div class="flex justify-end space-x-2 pt-4 border-t">
            <UiButton type="button" variant="outline" @click="handleClose" :disabled="isLoading"> Cancel </UiButton>
            <UiButton type="button" @click="handleSave" :loading="isLoading" class="gap-2">
              <Check class="h-4 w-4" />
              Save Changes
            </UiButton>
          </div>
        </div>
      </template>
    </UiDialogContent>
  </UiDialog>
</template>
