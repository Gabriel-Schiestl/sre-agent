import { useMutation, useQueryClient } from "@tanstack/react-query"
import { createMicroservice, updateMicroservice, deleteMicroservice } from "@/lib/api"
import type { CreateMicroservicePayload, UpdateMicroservicePayload } from "@/types/microservice"

export function useCreateMicroservice(suiteId: string) {
  const qc = useQueryClient()
  return useMutation({
    mutationFn: (payload: CreateMicroservicePayload) => createMicroservice(suiteId, payload),
    onSuccess: () => qc.invalidateQueries({ queryKey: ["suites", suiteId] }),
  })
}

export function useUpdateMicroservice(suiteId: string) {
  const qc = useQueryClient()
  return useMutation({
    mutationFn: ({ id, payload }: { id: string; payload: UpdateMicroservicePayload }) =>
      updateMicroservice(id, payload),
    onSuccess: () => qc.invalidateQueries({ queryKey: ["suites", suiteId] }),
  })
}

export function useDeleteMicroservice(suiteId: string) {
  const qc = useQueryClient()
  return useMutation({
    mutationFn: (id: string) => deleteMicroservice(id),
    onSuccess: () => qc.invalidateQueries({ queryKey: ["suites", suiteId] }),
  })
}
