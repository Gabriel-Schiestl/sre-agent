export interface ErrorCategory {
  category: string
  description: string
  occurrences: number
  affectedEndpoints: string[]
  severity: "low" | "medium" | "high" | "critical"
}

export interface Hypothesis {
  title: string
  evidence: string
  priority: number
}

export interface Bottleneck {
  microservice: string
  confidence: "low" | "medium" | "high"
  hypotheses: Hypothesis[]
}

export interface Diagnosis {
  id: string
  testRunId: string
  errorPlan: ErrorCategory[]
  bottlenecks: Bottleneck[]
  nextSteps: string[]
  createdAt: string
}

export interface ChatMessage {
  role: "user" | "assistant"
  content: string
}
