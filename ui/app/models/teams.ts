export type TeamData = {
  id: string
  name: string
  description: string
  status: Status
  memberCount: number
  owner: string
  adminCount: number
  developerCount: number
  viewerCount: number
  policies: string[]
  costCenter: string
  createdAt: string
  lastActivity: string
}

export type TeamMember = {
  id: string
  name: string
  email: string
  role: AvailableRoles
  status: Status
  lastActive: string
  avatar?: string
}
