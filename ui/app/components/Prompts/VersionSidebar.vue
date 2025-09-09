<script setup lang="ts">
import { GitBranch, Play, Copy, Archive, Trash2, Calendar, User, TrendingUp, CheckCircle, Eye, Edit, Activity } from "lucide-vue-next"

interface PromptVersion {
  id: string
  version: string
  status: 'approved' | 'review' | 'draft' | 'archived'
  deployments?: string[]
  content: {
    userPrompt: string
    systemPrompt: string
    config: any
  }
  createdAt: string
  createdBy: string
  usageCount: number
  isActive?: boolean
  isActiveInEnvironment?: boolean
}

interface Props {
  versions: PromptVersion[]
  selectedVersionId?: string
}

const props = defineProps<Props>()

const emit = defineEmits<{
  selectVersion: [versionId: string]
  deployVersion: [versionId: string]
  promoteVersion: [versionId: string, environment: string]
  cloneVersion: [versionId: string]
  archiveVersion: [versionId: string]
  deleteVersion: [versionId: string]
  openPlayground: [versionId: string]
}>()

const searchQuery = ref('')

// Archive modal state
const showArchiveModal = ref(false)
const selectedVersionForArchive = ref<PromptVersion | null>(null)

// Promotion modal state
const showPromoteModal = ref(false)
const selectedVersionForPromotion = ref<PromptVersion | null>(null)
const selectedEnvironment = ref('')

// Available environments - would come from API in real app
const availableEnvironments = ref(['Production', 'Staging', 'Development', 'Testing'])

const filteredVersions = computed(() => {
  if (!searchQuery.value) return props.versions
  
  const query = searchQuery.value.toLowerCase()
  return props.versions.filter(version => 
    version.version.toLowerCase().includes(query) ||
    version.createdBy.toLowerCase().includes(query)
  )
})

const openArchiveModal = (version: PromptVersion) => {
  selectedVersionForArchive.value = version
  showArchiveModal.value = true
}

const handleArchive = () => {
  if (selectedVersionForArchive.value) {
    emit('archiveVersion', selectedVersionForArchive.value.id)
    showArchiveModal.value = false
    selectedVersionForArchive.value = null
  }
}

const cancelArchive = () => {
  showArchiveModal.value = false
  selectedVersionForArchive.value = null
}

const openPromoteModal = (version: PromptVersion) => {
  selectedVersionForPromotion.value = version
  showPromoteModal.value = true
}

const handlePromote = () => {
  if (selectedVersionForPromotion.value) {
    emit('promoteVersion', selectedVersionForPromotion.value.id, 'current')
    showPromoteModal.value = false
    selectedVersionForPromotion.value = null
  }
}

const cancelPromote = () => {
  showPromoteModal.value = false
  selectedVersionForPromotion.value = null
}

const getStatusIcon = (status: string) => {
  switch (status) {
    case 'approved':
      return { icon: CheckCircle, class: 'text-green-600 dark:text-green-400' }
    case 'review':
      return { icon: Eye, class: 'text-blue-600 dark:text-blue-400' }
    case 'draft':
      return { icon: Edit, class: 'text-orange-600 dark:text-orange-400' }
    case 'archived':
      return { icon: Archive, class: 'text-gray-600 dark:text-gray-400' }
    default:
      return { icon: Edit, class: 'text-gray-600 dark:text-gray-400' }
  }
}

const formatDate = (dateString: string) => {
  return new Date(dateString).toLocaleDateString("en-US", {
    month: "short",
    day: "numeric",
    year: "numeric"
  })
}

const getInitials = (name: string) => {
  return name.split(" ").map(n => n[0]).join("")
}
</script>

<template>
  <div class="w-80 flex flex-col rounded-xl border bg-card shadow-lg overflow-hidden">
    <!-- Header -->
    <div class="p-4 border-b bg-card">
      <h3 class="font-semibold text-sm flex items-center gap-2 text-foreground">
        <GitBranch class="h-4 w-4 text-primary" />
        Prompt Versions
      </h3>
      <p class="text-xs text-muted-foreground mt-1">Select a version to view or edit</p>
    </div>

    <!-- Search -->
    <div class="p-3 border-b bg-muted/20">
      <UiInput
        v-model="searchQuery"
        placeholder="Search versions..."
        class="text-sm"
      />
    </div>

    <!-- Versions List -->
    <div class="flex-1 min-h-0">
      <UiScrollArea class="h-full">
        <div class="p-2">
          <div class="space-y-2">
            <div
              v-for="version in filteredVersions"
              :key="version.id"
              :class="[
                'p-3 rounded-lg border cursor-pointer transition-all duration-200 ease-in-out hover:shadow-sm',
                selectedVersionId === version.id
                  ? 'border-primary bg-primary/5 ring-1 ring-primary/20'
                  : 'border-border bg-background hover:bg-muted/30 hover:border-border/60'
              ]"
              @click="emit('selectVersion', version.id)"
            >
              <!-- Version Header -->
              <div class="flex items-start justify-between mb-3">
                <div class="flex items-center gap-2">
                  <span class="font-semibold text-sm text-foreground">{{ version.version }}</span>
                  <div class="flex items-center gap-1">
                    <component 
                      :is="getStatusIcon(version.status).icon" 
                      :class="['h-4 w-4', getStatusIcon(version.status).class]"
                      :title="version.status.charAt(0).toUpperCase() + version.status.slice(1)"
                    />
                    <!-- Active Status Indicator -->
                    <div 
                      v-if="version.isActiveInEnvironment" 
                      class="flex items-center gap-1 ml-1"
                    >
                      <div class="h-2 w-2 bg-green-500 rounded-full animate-pulse"></div>
                      <span class="text-xs text-green-600 dark:text-green-400 font-medium">ACTIVE</span>
                    </div>
                  </div>
                </div>
                
                <!-- Quick Actions -->
                <div class="flex items-center gap-1">
                  <button
                    @click.stop="emit('openPlayground', version.id)"
                    class="p-1 rounded hover:bg-muted/50 transition-colors"
                    title="Test in Playground"
                  >
                    <Play class="h-3 w-3 text-muted-foreground hover:text-primary" />
                  </button>
                  <button
                    @click.stop="emit('cloneVersion', version.id)"
                    class="p-1 rounded hover:bg-muted/50 transition-colors"
                    title="Clone Version"
                  >
                    <Copy class="h-3 w-3 text-muted-foreground hover:text-primary" />
                  </button>
                  <button
                    v-if="version.status === 'approved' && !version.isActiveInEnvironment"
                    @click.stop="openPromoteModal(version)"
                    class="p-1 rounded hover:bg-muted/50 transition-colors"
                    title="Deploy to Current Environment"
                  >
                    <TrendingUp class="h-3 w-3 text-muted-foreground hover:text-green-600" />
                  </button>
                </div>
              </div>

              <!-- Version Meta -->
              <div class="space-y-2">
                <div class="flex items-center justify-between text-xs">
                  <div class="flex items-center gap-2 text-muted-foreground">
                    <Calendar class="h-3 w-3" />
                    <span>{{ formatDate(version.createdAt) }}</span>
                  </div>
                  <div class="flex items-center gap-1">
                    <div class="w-4 h-4 bg-primary/10 rounded-full flex items-center justify-center">
                      <span class="text-xs font-medium text-primary">{{ getInitials(version.createdBy) }}</span>
                    </div>
                    <span class="text-muted-foreground">{{ version.createdBy }}</span>
                  </div>
                </div>
              </div>


              <!-- Archive Action -->
              <div v-if="version.status !== 'archived'" class="flex justify-end mt-2">
                <button
                  @click.stop="openArchiveModal(version)"
                  class="p-1 rounded hover:bg-muted/50 transition-colors"
                  title="Archive Version"
                >
                  <Archive class="h-3 w-3 text-muted-foreground hover:text-orange-600" />
                </button>
              </div>
            </div>
          </div>
        </div>
      </UiScrollArea>
    </div>


    <!-- Archive Confirmation Modal -->
    <UiDialog v-model:open="showArchiveModal">
      <UiDialogContent class="max-w-md">
        <UiDialogHeader>
          <UiDialogTitle class="flex items-center gap-2">
            <Archive class="h-5 w-5 text-orange-600" />
            Archive Version
          </UiDialogTitle>
          <UiDialogDescription>
            Are you sure you want to archive version {{ selectedVersionForArchive?.version }}?
          </UiDialogDescription>
        </UiDialogHeader>
        
        <div class="py-4">
          <div class="p-3 rounded-lg bg-orange-50 dark:bg-orange-900/20 border border-orange-200 dark:border-orange-700">
            <p class="text-sm text-orange-800 dark:text-orange-200">
              <strong>Note:</strong> Archived versions can still be viewed but cannot be deployed or edited.
            </p>
          </div>
        </div>
        
        <UiDialogFooter>
          <UiButton variant="outline" @click="cancelArchive">
            Cancel
          </UiButton>
          <UiButton 
            variant="destructive" 
            @click="handleArchive"
          >
            <Archive class="h-4 w-4 mr-2" />
            Archive
          </UiButton>
        </UiDialogFooter>
      </UiDialogContent>
    </UiDialog>

    <!-- Deploy to Environment Modal -->
    <UiDialog v-model:open="showPromoteModal">
      <UiDialogContent class="max-w-md">
        <UiDialogHeader>
          <UiDialogTitle class="flex items-center gap-2">
            <TrendingUp class="h-5 w-5 text-green-600" />
            Deploy Version
          </UiDialogTitle>
          <UiDialogDescription>
            Deploy version {{ selectedVersionForPromotion?.version }} to the current environment
          </UiDialogDescription>
        </UiDialogHeader>
        
        <div class="space-y-4 py-4">
          <div class="p-3 rounded-lg bg-blue-50 dark:bg-blue-900/20 border border-blue-200 dark:border-blue-700">
            <p class="text-sm text-blue-800 dark:text-blue-200">
              <strong>Target Environment:</strong> Current environment (shown in sidebar)
            </p>
          </div>
          
          <div class="p-3 rounded-lg bg-yellow-50 dark:bg-yellow-900/20 border border-yellow-200 dark:border-yellow-700">
            <p class="text-sm text-yellow-800 dark:text-yellow-200">
              <strong>Warning:</strong> This will replace the currently active version in this environment.
            </p>
          </div>
        </div>
        
        <UiDialogFooter>
          <UiButton variant="outline" @click="cancelPromote">
            Cancel
          </UiButton>
          <UiButton 
            variant="default" 
            @click="handlePromote"
            class="bg-green-600 hover:bg-green-700"
          >
            <TrendingUp class="h-4 w-4 mr-2" />
            Deploy
          </UiButton>
        </UiDialogFooter>
      </UiDialogContent>
    </UiDialog>
  </div>
</template>