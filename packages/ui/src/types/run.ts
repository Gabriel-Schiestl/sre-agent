export type RunStatus = "pending" | "analyzing" | "done" | "failed"

export interface TestRun {
  id: string
  testSuiteId: string
  name: string
  virtualUsers: number
  durationSeconds: number
  notes: string
  status: RunStatus
  createdAt: string
}

export interface CreateRunPayload {
  name: string
  virtualUsers: number
  durationSeconds: number
  notes?: string
  file: File
}
