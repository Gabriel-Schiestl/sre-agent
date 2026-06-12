import type { Microservice } from "@/types/microservice"
import type { TestRun } from "@/types/run"

export interface TestSuite {
  id: string
  name: string
  description: string
  microservices: Microservice[]
  runs: TestRun[]
  createdAt: string
  updatedAt: string
}

export interface CreateSuitePayload {
  name: string
  description: string
}

export interface UpdateSuitePayload {
  name: string
  description: string
}
