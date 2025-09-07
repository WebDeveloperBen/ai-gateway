export type Model = {
  id: string
  name: string
  provider: string
  description: string
  maxTokens: number
  costPer1kTokens: {
    input: number
    output: number
  }
}
