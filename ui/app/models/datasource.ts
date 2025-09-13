export type DataSourceData = {
  id: string
  name: string
  description: string
  type: string
  source: string
  url: string
  status: "active" | "paused" | "error" | "inactive"
  schedule: string
  lastSync: string
  nextSync: string
  documentsCount: number
  owner: string
  tags: string[]
}

export type DataSourceType = 
  | "documentation" 
  | "knowledge-base" 
  | "policies" 
  | "api-docs"

export type DataSourceSource = 
  | "confluence" 
  | "notion" 
  | "sharepoint" 
  | "github" 
  | "google-drive"

export type DataSourceStatus = 
  | "active" 
  | "paused" 
  | "error" 
  | "inactive"