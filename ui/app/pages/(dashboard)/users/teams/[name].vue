<script lang="ts">
// Mock team data - replace with API call
const teams: TeamData[] = [
  {
    id: "team_1",
    name: "Engineering",
    description: "Core engineering team for product development",
    status: "active",
    memberCount: 12,
    owner: "Alice Johnson",
    adminCount: 2,
    developerCount: 8,
    viewerCount: 2,
    policies: ["Development Policy", "Security Policy"],
    costCenter: "ENG-001",
    createdAt: "2024-12-01T08:00:00Z",
    lastActivity: "2025-01-15T14:30:00Z"
  },
  {
    id: "team_2",
    name: "Marketing",
    description: "Content creation and marketing campaigns",
    status: "active",
    memberCount: 6,
    owner: "Carol Williams",
    adminCount: 1,
    developerCount: 2,
    viewerCount: 3,
    policies: ["Marketing Policy"],
    costCenter: "MKT-002",
    createdAt: "2024-11-15T10:30:00Z",
    lastActivity: "2025-01-14T16:45:00Z"
  },
  {
    id: "team_3",
    name: "Customer Success",
    description: "Customer support and success operations",
    status: "active",
    memberCount: 8,
    owner: "David Brown",
    adminCount: 1,
    developerCount: 3,
    viewerCount: 4,
    policies: ["Customer Data Policy", "Security Policy"],
    costCenter: "CS-003",
    createdAt: "2024-10-20T09:15:00Z",
    lastActivity: "2025-01-15T11:20:00Z"
  },
  {
    id: "team_4",
    name: "Analytics",
    description: "Data analysis and business intelligence",
    status: "inactive",
    memberCount: 3,
    owner: "Emma Davis",
    adminCount: 1,
    developerCount: 1,
    viewerCount: 1,
    policies: ["Data Policy"],
    costCenter: "ANL-004",
    createdAt: "2024-09-10T11:00:00Z",
    lastActivity: "2024-12-22T10:30:00Z"
  }
]

// Mock team members - replace with API call
const teamMembers: TeamMember[] = [
  {
    id: "1",
    name: "Alice Johnson",
    email: "alice@company.com",
    role: "Owner",
    status: "active",
    lastActive: "2 hours ago",
    avatar: "https://images.unsplash.com/photo-1494790108755-2616b95b0c97?w=40&h=40&fit=crop&crop=face"
  },
  {
    id: "2",
    name: "Bob Smith",
    email: "bob@company.com",
    role: "Admin",
    status: "active",
    lastActive: "1 day ago"
  },
  {
    id: "3",
    name: "Carol Williams",
    email: "carol@company.com",
    role: "Admin",
    status: "active",
    lastActive: "3 hours ago",
    avatar: "https://images.unsplash.com/photo-1438761681033-6461ffad8d80?w=40&h=40&fit=crop&crop=face"
  },
  {
    id: "4",
    name: "David Brown",
    email: "david@company.com",
    role: "Developer",
    status: "active",
    lastActive: "30 minutes ago"
  },
  {
    id: "5",
    name: "Emma Davis",
    email: "emma@company.com",
    role: "Developer",
    status: "active",
    lastActive: "1 hour ago",
    avatar: "https://images.unsplash.com/photo-1472099645785-5658abf4ff4e?w=40&h=40&fit=crop&crop=face"
  },
  {
    id: "6",
    name: "Frank Wilson",
    email: "frank@company.com",
    role: "Developer",
    status: "active",
    lastActive: "2 days ago"
  },
  {
    id: "7",
    name: "Grace Lee",
    email: "grace@company.com",
    role: "Viewer",
    status: "active",
    lastActive: "4 hours ago"
  },
  {
    id: "8",
    name: "Henry Chen",
    email: "henry@company.com",
    role: "Viewer",
    status: "inactive",
    lastActive: "1 week ago"
  }
]
</script>

<script setup lang="ts">
import {
  Users,
  Crown,
  Shield,
  UserCheck,
  Eye,
  Building,
  FileText,
  ArrowLeft,
  Edit,
  Plus,
  MoreVertical,
  Trash2,
  UserX
} from "lucide-vue-next"

const route = useRoute()
const teamId = route.params.name as string

// Find team data - replace with API call
const team = computed(() => teams.find((t) => t.name === teamId))

// Filter members for this team (in real app, this would be an API call with team filter)
const members = computed(() => teamMembers)

useSeoMeta({
  title: computed(() => (team.value ? `${team.value.name} - Teams` : "Team Not Found"))
})

// 404 handling
if (!team.value) {
  throw createError({
    statusCode: 404,
    statusMessage: "Team not found"
  })
}

// Helper functions

// Modal states
const showEditModal = ref(false)
const showAddMemberModal = ref(false)
const showDeleteModal = ref(false)
const deletingMember = ref<TeamMember | null>(null)

// Handlers
const openEditModal = () => {
  showEditModal.value = true
}

const openAddMemberModal = () => {
  showAddMemberModal.value = true
}

const openDeleteMemberModal = (member: TeamMember) => {
  deletingMember.value = member
  showDeleteModal.value = true
}
</script>

<template>
  <div v-if="team" class="flex flex-col gap-6">
    <PageHeader :title="team.name" :subtext="team.description">
      <div class="flex items-center gap-2">
        <UiButton variant="outline" @click="openAddMemberModal">
          <Plus class="mr-2 size-4" />
          Add Member
        </UiButton>
        <UiButton @click="openEditModal">
          <Edit class="mr-2 size-4" />
          Edit Team
        </UiButton>
      </div>
    </PageHeader>

    <!-- Team Overview Cards -->
    <div class="grid gap-6 md:grid-cols-2 lg:grid-cols-3">
      <!-- Basic Info Card -->
      <UiCard>
        <UiCardHeader>
          <UiCardTitle class="flex items-center gap-2">
            <Building class="size-5" />
            Team Information
          </UiCardTitle>
        </UiCardHeader>
        <UiCardContent class="space-y-4">
          <div class="flex items-center justify-between">
            <span class="text-sm text-muted-foreground">Status</span>
            <div
              class="inline-flex items-center gap-1 px-2 py-1 rounded-md text-xs font-medium border"
              :class="getStatusColor(team.status)"
            >
              <div class="size-1.5 rounded-full" :class="team.status === 'active' ? 'bg-green-500' : 'bg-gray-400'" />
              {{ team.status === "active" ? "Active" : "Inactive" }}
            </div>
          </div>
          <div class="flex items-center justify-between">
            <span class="text-sm text-muted-foreground">Owner</span>
            <span class="text-sm font-medium">{{ team.owner }}</span>
          </div>
          <div class="flex items-center justify-between">
            <span class="text-sm text-muted-foreground">Created</span>
            <span class="text-sm font-medium">{{ new Date(team.createdAt).toLocaleDateString() }}</span>
          </div>
          <div class="flex items-center justify-between">
            <span class="text-sm text-muted-foreground">Last Activity</span>
            <span class="text-sm font-medium">{{ new Date(team.lastActivity).toLocaleDateString() }}</span>
          </div>
        </UiCardContent>
      </UiCard>

      <!-- Members Overview Card -->
      <UiCard>
        <UiCardHeader>
          <UiCardTitle class="flex items-center gap-2">
            <Users class="size-5" />
            Members ({{ team.memberCount }})
          </UiCardTitle>
        </UiCardHeader>
        <UiCardContent class="space-y-3">
          <div class="flex items-center justify-between">
            <div class="flex items-center gap-2">
              <Crown class="size-4 text-orange-600" />
              <span class="text-sm text-muted-foreground">Owner</span>
            </div>
            <span class="text-sm font-medium">1</span>
          </div>
          <div class="flex items-center justify-between">
            <div class="flex items-center gap-2">
              <Shield class="size-4 text-blue-600" />
              <span class="text-sm text-muted-foreground">Admins</span>
            </div>
            <span class="text-sm font-medium">{{ team.adminCount }}</span>
          </div>
          <div class="flex items-center justify-between">
            <div class="flex items-center gap-2">
              <UserCheck class="size-4 text-green-600" />
              <span class="text-sm text-muted-foreground">Developers</span>
            </div>
            <span class="text-sm font-medium">{{ team.developerCount }}</span>
          </div>
          <div class="flex items-center justify-between">
            <div class="flex items-center gap-2">
              <Eye class="size-4 text-purple-600" />
              <span class="text-sm text-muted-foreground">Viewers</span>
            </div>
            <span class="text-sm font-medium">{{ team.viewerCount }}</span>
          </div>
        </UiCardContent>
      </UiCard>

      <!-- Policies Card -->
      <UiCard>
        <UiCardHeader>
          <div class="flex items-center justify-between">
            <UiCardTitle class="flex items-center gap-2">
              <FileText class="size-5" />
              Policies ({{ team.policies.length }})
            </UiCardTitle>
            <UiButton variant="outline" size="sm" @click="navigateTo('/governance/policies')"> View All </UiButton>
          </div>
        </UiCardHeader>
        <UiCardContent>
          <div class="space-y-3 max-h-48 overflow-y-auto">
            <div
              v-for="(policy, index) in team.policies"
              :key="policy"
              class="group flex items-center gap-3 p-3 border rounded-lg hover:bg-muted/50 cursor-pointer transition-colors"
              @click="navigateTo(`/governance/policies/policy_${index + 1}`)"
            >
              <div class="p-2 rounded-lg bg-blue-50 text-blue-600">
                <FileText class="size-4" />
              </div>
              <div class="flex-1 min-w-0">
                <div class="font-medium text-sm">{{ policy }}</div>
                <div class="text-xs text-muted-foreground">Click to view policy details</div>
              </div>
              <div class="opacity-0 group-hover:opacity-100 transition-opacity">
                <ArrowLeft class="size-4 rotate-180 text-muted-foreground" />
              </div>
            </div>
            <div v-if="team.policies.length === 0" class="text-center py-4 text-muted-foreground text-sm">
              No policies assigned to this team
            </div>
          </div>
        </UiCardContent>
      </UiCard>
    </div>

    <!-- Team Members -->
    <CardsDataList :title="`${team.name} Members`" :icon="Users">
      <template #actions>
        <UiButton variant="outline" size="sm" @click="openAddMemberModal">
          <Plus class="mr-2 size-4" />
          Add Member
        </UiButton>
      </template>

      <div class="space-y-4">
        <div
          v-for="member in members"
          :key="member.id"
          class="flex items-center justify-between p-4 border rounded-lg hover:bg-muted/50"
        >
          <div class="flex items-center gap-4">
            <UiAvatar class="size-10">
              <UiAvatarImage v-if="member.avatar" :src="member.avatar" :alt="member.name" />
              <UiAvatarFallback>{{
                member.name
                  .split(" ")
                  .map((n) => n[0])
                  .join("")
              }}</UiAvatarFallback>
            </UiAvatar>

            <div class="space-y-1">
              <div class="flex items-center gap-3">
                <p class="font-medium">{{ member.name }}</p>
                <div
                  class="inline-flex items-center gap-1 px-2 py-1 rounded-md text-xs font-medium border"
                  :class="getRoleColor(member.role)"
                >
                  <component :is="getRoleIcon(member.role)" class="size-3" />
                  {{ member.role }}
                </div>
              </div>
              <div class="flex items-center gap-4 text-sm text-muted-foreground">
                <span>{{ member.email }}</span>
                <span>•</span>
                <span>Last active {{ member.lastActive }}</span>
              </div>
            </div>
          </div>

          <div class="flex items-center gap-2">
            <div
              class="inline-flex items-center gap-1 px-2 py-1 rounded-md text-xs font-medium border"
              :class="getStatusColor(member.status)"
            >
              <div class="size-1.5 rounded-full" :class="member.status === 'active' ? 'bg-green-500' : 'bg-gray-400'" />
              {{ member.status === "active" ? "Active" : "Inactive" }}
            </div>

            <UiDropdownMenu v-if="member.role !== 'Owner'">
              <UiDropdownMenuTrigger as-child>
                <UiButton variant="ghost" size="sm">
                  <MoreVertical class="size-4" />
                </UiButton>
              </UiDropdownMenuTrigger>
              <UiDropdownMenuContent align="end" class="w-48">
                <UiDropdownMenuItem @click="() => {}">
                  <Edit class="mr-2 size-4" />
                  Edit Member
                </UiDropdownMenuItem>
                <UiDropdownMenuItem @click="() => {}">
                  <component :is="member.status === 'active' ? UserX : UserCheck" class="mr-2 size-4" />
                  {{ member.status === "active" ? "Deactivate" : "Activate" }}
                </UiDropdownMenuItem>
                <UiDropdownMenuSeparator />
                <UiDropdownMenuItem class="text-destructive" @click="openDeleteMemberModal(member)">
                  <Trash2 class="mr-2 size-4" />
                  Remove from Team
                </UiDropdownMenuItem>
              </UiDropdownMenuContent>
            </UiDropdownMenu>
          </div>
        </div>
      </div>
    </CardsDataList>

    <!-- Confirmation Modal for removing members -->
    <ConfirmationModal
      v-model:open="showDeleteModal"
      title="Remove Team Member"
      :description="`Are you sure you want to remove ${deletingMember?.name} from this team? They will lose access to team resources.`"
      confirm-text="Remove Member"
      variant="destructive"
      :loading="false"
      @confirm="
        () => {
          showDeleteModal = false
          deletingMember = null
        }
      "
      @cancel="
        () => {
          showDeleteModal = false
          deletingMember = null
        }
      "
    />
  </div>
</template>
