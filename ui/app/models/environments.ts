export type EnvironmentData = {
  id: string
  name: string
  description: string
  status: "active" | "inactive"
  memberCount: number
  applicationCount: number
  teams: string[]
  owner: string
  monthlyRequests: number
  createdAt: string
  lastActivity: string
}
