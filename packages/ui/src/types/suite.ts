import type { Microservice } from "@/types/microservice"

export interface TestSuite {
  id: string
  name: string
  description: string
  microservices: Microservice[]
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
