import { useQuery, useMutation, useQueryClient } from "@tanstack/react-query"
import { getSuites, getSuite, createSuite, updateSuite, deleteSuite } from "@/lib/api"
import type { CreateSuitePayload, UpdateSuitePayload } from "@/types/suite"

export function useSuites() {
  return useQuery({ queryKey: ["suites"], queryFn: getSuites })
}

export function useSuite(id: string) {
  return useQuery({ queryKey: ["suites", id], queryFn: () => getSuite(id), enabled: !!id })
}

export function useCreateSuite() {
  const qc = useQueryClient()
  return useMutation({
    mutationFn: (payload: CreateSuitePayload) => createSuite(payload),
    onSuccess: () => qc.invalidateQueries({ queryKey: ["suites"] }),
  })
}

export function useUpdateSuite(id: string) {
  const qc = useQueryClient()
  return useMutation({
    mutationFn: (payload: UpdateSuitePayload) => updateSuite(id, payload),
    onSuccess: () => {
      qc.invalidateQueries({ queryKey: ["suites"] })
      qc.invalidateQueries({ queryKey: ["suites", id] })
    },
  })
}

export function useDeleteSuite() {
  const qc = useQueryClient()
  return useMutation({
    mutationFn: (id: string) => deleteSuite(id),
    onSuccess: () => qc.invalidateQueries({ queryKey: ["suites"] }),
  })
}
