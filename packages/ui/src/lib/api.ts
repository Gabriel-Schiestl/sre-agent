import axios from "axios"
import type { TestSuite, CreateSuitePayload, UpdateSuitePayload } from "@/types/suite"
import type { Microservice, CreateMicroservicePayload, UpdateMicroservicePayload } from "@/types/microservice"
import type { TestRun } from "@/types/run"
import type { Diagnosis } from "@/types/diagnosis"

const api = axios.create({
  baseURL: process.env.NEXT_PUBLIC_API_URL ?? "http://localhost:8080",
})

api.interceptors.response.use(
  (response) => response,
  (error) => {
    const message =
      error.response?.data?.message ??
      error.response?.data?.error ??
      error.message ??
      "Unexpected error"
    return Promise.reject(new Error(message))
  }
)

// ─── Suites ──────────────────────────────────────────────────────────────────

export async function getSuites(): Promise<TestSuite[]> {
  const { data } = await api.get<TestSuite[]>("/suites")
  return data
}

export async function getSuite(id: string): Promise<TestSuite> {
  const { data } = await api.get<TestSuite>(`/suites/${id}`)
  return data
}

export async function createSuite(payload: CreateSuitePayload): Promise<TestSuite> {
  const { data } = await api.post<TestSuite>("/suites", payload)
  return data
}

export async function updateSuite(id: string, payload: UpdateSuitePayload): Promise<TestSuite> {
  const { data } = await api.put<TestSuite>(`/suites/${id}`, payload)
  return data
}

export async function deleteSuite(id: string): Promise<void> {
  await api.delete(`/suites/${id}`)
}

// ─── Microservices ───────────────────────────────────────────────────────────

export async function createMicroservice(
  suiteId: string,
  payload: CreateMicroservicePayload
): Promise<Microservice> {
  const { data } = await api.post<Microservice>(`/suites/${suiteId}/microservices`, payload)
  return data
}

export async function updateMicroservice(
  id: string,
  payload: UpdateMicroservicePayload
): Promise<Microservice> {
  const { data } = await api.put<Microservice>(`/microservices/${id}`, payload)
  return data
}

export async function deleteMicroservice(id: string): Promise<void> {
  await api.delete(`/microservices/${id}`)
}

// ─── Runs ─────────────────────────────────────────────────────────────────────

export async function getRuns(suiteId: string): Promise<TestRun[]> {
  const { data } = await api.get<TestRun[]>(`/suites/${suiteId}/runs`)
  return data
}

export async function getRun(id: string): Promise<TestRun> {
  const { data } = await api.get<TestRun>(`/runs/${id}`)
  return data
}

export async function createRun(
  suiteId: string,
  payload: { name: string; virtualUsers: number; durationSeconds: number; notes?: string; file: File }
): Promise<TestRun> {
  const form = new FormData()
  form.append("name", payload.name)
  form.append("virtualUsers", String(payload.virtualUsers))
  form.append("durationSeconds", String(payload.durationSeconds))
  if (payload.notes) form.append("notes", payload.notes)
  form.append("jtlFile", payload.file)

  const { data } = await api.post<TestRun>(`/suites/${suiteId}/runs`, form, {
    headers: { "Content-Type": "multipart/form-data" },
  })
  return data
}

export async function getDiagnosis(runId: string): Promise<Diagnosis> {
  const { data } = await api.get<Diagnosis>(`/runs/${runId}/diagnosis`)
  return data
}
