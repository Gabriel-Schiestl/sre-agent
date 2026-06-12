import { useQuery } from "@tanstack/react-query"
import { getDiagnosis } from "@/lib/api"

export function useDiagnosis(runId: string) {
  return useQuery({
    queryKey: ["diagnosis", runId],
    queryFn: () => getDiagnosis(runId),
    enabled: !!runId,
  })
}
