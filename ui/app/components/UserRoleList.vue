<script setup lang="ts">
import { MoreVertical, Edit, Trash2, UserX, UserCheck } from "lucide-vue-next"

interface RoleAssignment {
  userId: string
  userName: string
  userEmail: string
  role: "Owner" | "Admin" | "Developer" | "Viewer"
  team: string
  status: "active" | "inactive"
  lastActive: string
  avatar?: string
  assignedDate: string
}

interface Props {
  assignments: RoleAssignment[]
  compact?: boolean
  viewOnly?: boolean
}

const props = withDefaults(defineProps<Props>(), {
  compact: false,
  viewOnly: false
})

const emit = defineEmits<{
  editRole: [assignment: RoleAssignment]
  toggleStatus: [assignment: RoleAssignment]
  removeRole: [assignment: RoleAssignment]
  viewUser: [assignment: RoleAssignment]
}>()

function handleEditRole(assignment: RoleAssignment) {
  emit("editRole", assignment)
}

function handleToggleStatus(assignment: RoleAssignment) {
  emit("toggleStatus", assignment)
}

function handleRemoveRole(assignment: RoleAssignment) {
  emit("removeRole", assignment)
}

function getStatusColor(status: string) {
  switch (status) {
    case "active":
      return "text-emerald-600 dark:text-emerald-400 bg-emerald-50 dark:bg-emerald-950/30 border-emerald-200 dark:border-emerald-800/50"
    case "inactive":
      return "text-gray-600 dark:text-gray-400 bg-gray-50 dark:bg-gray-900/30 border-gray-200 dark:border-gray-700/50"
    default:
      return "text-gray-600 dark:text-gray-400 bg-gray-50 dark:bg-gray-900/30 border-gray-200 dark:border-gray-700/50"
  }
}

async function handleViewUser(assignment: RoleAssignment) {
  await navigateTo(`/users?userId=${assignment.userId}`)
}
</script>

<template>
  <div class="space-y-4">
    <div
      v-for="assignment in assignments"
      :key="assignment.userId"
      class="flex items-center justify-between border rounded-lg hover:bg-muted/50 transition-colors"
      :class="compact ? 'p-3' : 'p-4'"
    >
      <div class="flex items-center gap-4">
        <UiAvatar :class="compact ? 'size-8' : 'size-10'">
          <UiAvatarImage v-if="assignment.avatar" :src="assignment.avatar" :alt="assignment.userName" />
          <UiAvatarFallback>{{
            assignment.userName
              .split(" ")
              .map((n) => n[0])
              .join("")
          }}</UiAvatarFallback>
        </UiAvatar>

        <div class="space-y-1">
          <div class="flex items-center gap-2">
            <p :class="compact ? 'text-sm font-medium' : 'font-medium'">{{ assignment.userName }}</p>
            <UiBadge :variant="getRoleColor(assignment.role)" :class="compact ? 'text-xs' : 'text-xs'">
              <component :is="getRoleIcon(assignment.role)" class="size-3 mr-1" />
              {{ assignment.role }}
            </UiBadge>
          </div>
          <div :class="`flex items-center gap-4 ${compact ? 'text-xs' : 'text-sm'} text-muted-foreground`">
            <span>{{ assignment.userEmail }}</span>
            <span>•</span>
            <span>{{ assignment.team }} team</span>
            <span v-if="!compact">•</span>
            <span v-if="!compact">Assigned {{ new Date(assignment.assignedDate).toLocaleDateString() }}</span>
          </div>
        </div>
      </div>

      <div class="flex items-center gap-2">
        <div
          class="inline-flex items-center gap-1 px-2 py-1 rounded-md text-xs font-medium border"
          :class="getStatusColor(assignment.status)"
        >
          <div
            class="size-1.5 rounded-full"
            :class="assignment.status === 'active' ? 'bg-emerald-500' : 'bg-gray-400'"
          />
          {{ assignment.status === "active" ? "Active" : "Inactive" }}
        </div>

        <template v-if="!props.viewOnly">
          <UiDropdownMenu>
            <UiDropdownMenuTrigger as-child>
              <UiButton variant="ghost" size="sm">
                <MoreVertical class="size-4" />
              </UiButton>
            </UiDropdownMenuTrigger>
            <UiDropdownMenuContent align="end" class="w-48">
              <UiDropdownMenuItem @click="handleEditRole(assignment)">
                <Edit class="mr-2 size-4" />
                Change Role
              </UiDropdownMenuItem>
              <UiDropdownMenuItem @click="handleToggleStatus(assignment)">
                <component :is="assignment.status === 'active' ? UserX : UserCheck" class="mr-2 size-4" />
                {{ assignment.status === "active" ? "Deactivate" : "Activate" }}
              </UiDropdownMenuItem>
              <UiDropdownMenuSeparator />
              <UiDropdownMenuItem class="text-destructive" @click="handleRemoveRole(assignment)">
                <Trash2 class="mr-2 size-4" />
                Remove Role
              </UiDropdownMenuItem>
            </UiDropdownMenuContent>
          </UiDropdownMenu>
        </template>
        <template v-else>
          <UiButton variant="outline" size="sm" @click="handleViewUser(assignment)"> View </UiButton>
        </template>
      </div>
    </div>

    <div v-if="assignments.length === 0" class="text-center py-8 text-muted-foreground">
      <div class="size-12 mx-auto mb-4 text-muted-foreground/50 flex items-center justify-center">
        <component :is="getRoleIcon('Admin')" class="size-8" />
      </div>
      <h3 class="text-lg font-medium mb-2">No role assignments found</h3>
      <p class="text-sm">No users are assigned to this environment yet.</p>
    </div>
  </div>
</template>
