import { Crown, Eye, Shield, UserCheck, Users } from "lucide-vue-next"
import type { FunctionalComponent } from "vue"

export const availableRoles = [
  { value: "Admin", label: "Admin", icon: Shield, color: "text-blue-600" },
  { value: "Developer", label: "Developer", icon: UserCheck, color: "text-green-600" },
  { value: "Viewer", label: "Viewer", icon: Eye, color: "text-purple-600" }
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
      return "text-orange-600 bg-orange-50 border-orange-200"
    case "Admin":
      return "text-blue-600 bg-blue-50 border-blue-200"
    case "Developer":
      return "text-green-600 bg-green-50 border-green-200"
    case "Viewer":
      return "text-purple-600 bg-purple-50 border-purple-200"
    default:
      return "text-gray-600 bg-gray-50 border-gray-200"
  }
}

export function getStatusColor(status: string) {
  return status === "active"
    ? "text-green-600 bg-green-50 border-green-200"
    : "text-gray-600 bg-gray-50 border-gray-200"
}
