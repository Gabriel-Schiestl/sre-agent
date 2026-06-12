export interface Microservice {
  id: string
  testSuiteId: string
  name: string
  description: string
  language: string
  mainEndpoints: string[]
  cpuLimit?: string
  memoryLimit?: string
  sloLatencyP99Ms?: number
  sloErrorRatePct?: number
  createdAt: string
}

export interface CreateMicroservicePayload {
  name: string
  description: string
  language: string
  mainEndpoints: string[]
  cpuLimit?: string
  memoryLimit?: string
  sloLatencyP99Ms?: number
  sloErrorRatePct?: number
}

export type UpdateMicroservicePayload = CreateMicroservicePayload

export const LANGUAGE_OPTIONS = [
  { value: "Go", label: "Go" },
  { value: "Java", label: "Java" },
  { value: "Node.js", label: "Node.js" },
  { value: "Python", label: "Python" },
  { value: ".NET", label: ".NET" },
  { value: "Other", label: "Other" },
] as const
