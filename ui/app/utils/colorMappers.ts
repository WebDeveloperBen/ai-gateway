import { Crown, Eye, Shield, UserCheck, Users } from "lucide-vue-next"
import type { FunctionalComponent } from "vue"

export const availableRoles = [
  { value: "Admin", label: "Admin", icon: Shield, color: "text-blue-600 dark:text-blue-400" },
  { value: "Developer", label: "Developer", icon: UserCheck, color: "text-emerald-600 dark:text-emerald-400" },
  { value: "Viewer", label: "Viewer", icon: Eye, color: "text-purple-600 dark:text-purple-400" }
]

export function getRoleIcon(role: string): FunctionalComponent {
  switch (role) {
    case "Owner":
      return Crown
    case "Admin":
      return Shield
    case "Developer":
      return UserCheck
    case "Viewer":
      return Eye
    default:
      return Users
  }
}

export function getRoleColor(role: string) {
  switch (role) {
    case "Owner":
      return "text-amber-600 dark:text-amber-400 bg-amber-50 dark:bg-amber-950/30 border-amber-200 dark:border-amber-800/50"
    case "Admin":
      return "text-blue-600 dark:text-blue-400 bg-blue-50 dark:bg-blue-950/30 border-blue-200 dark:border-blue-800/50"
    case "Developer":
      return "text-emerald-600 dark:text-emerald-400 bg-emerald-50 dark:bg-emerald-950/30 border-emerald-200 dark:border-emerald-800/50"
    case "Viewer":
      return "text-purple-600 dark:text-purple-400 bg-purple-50 dark:bg-purple-950/30 border-purple-200 dark:border-purple-800/50"
    default:
      return "text-gray-600 dark:text-gray-400 bg-gray-50 dark:bg-gray-900/30 border-gray-200 dark:border-gray-700/50"
  }
}

export function getStatusColor(status: string) {
  return status === "active"
    ? "text-emerald-600 dark:text-emerald-400 bg-emerald-50 dark:bg-emerald-950/30 border-emerald-200 dark:border-emerald-800/50"
    : "text-gray-600 dark:text-gray-400 bg-gray-50 dark:bg-gray-900/30 border-gray-200 dark:border-gray-700/50"
}
