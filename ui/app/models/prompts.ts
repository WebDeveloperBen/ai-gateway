export type PromptTemplate = {
  id: string
  name: string
  content: string
  category: string
  description: string
  icon: any
}

export type PromptVersion = {
  id: string
  version: string
  name: string
  description?: string
  content: string
  systemPrompt?: string
  parameters?: {
    temperature?: number
    maxTokens?: number
    topP?: number
    frequencyPenalty?: number
    presencePenalty?: number
  }
  tags: string[]
  createdAt: string
  createdBy: string
  isPublished: boolean
  publishedAt?: string
}

export type SavedPrompt = {
  id: string
  name: string
  description?: string
  tags: string[]
  environments: string[]
  applications: string[]
  versions: PromptVersion[]
  currentVersion: string
  createdAt: string
  updatedAt: string
}

export type TestResult = {
  id: string
  timestamp: string
  prompt: string
  response: string
  model: string
  tokensUsed: {
    input: number
    output: number
    total: number
  }
  responseTime: number
  estimatedCost: number
  success: boolean
  error?: string
}
