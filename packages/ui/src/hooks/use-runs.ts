import { useQuery, useMutation, useQueryClient } from "@tanstack/react-query"
import { getRun, createRun } from "@/lib/api"

export function useRun(id: string) {
  return useQuery({
    queryKey: ["runs", id],
    queryFn: () => getRun(id),
    enabled: !!id,
  })
}

export function useRunStatus(id: string) {
  return useQuery({
    queryKey: ["runs", id],
    queryFn: () => getRun(id),
    enabled: !!id,
    refetchInterval: (query) => {
      const status = query.state.data?.status
      if (status === "done" || status === "failed") return false
      return 3_000
    },
  })
}

export function useCreateRun(suiteId: string) {
  const qc = useQueryClient()
  return useMutation({
    mutationFn: (payload: {
      name: string
      virtualUsers: number
      durationSeconds: number
      notes?: string
      file: File
    }) => createRun(suiteId, payload),
    onSuccess: () => qc.invalidateQueries({ queryKey: ["suites", suiteId] }),
  })
}
