"use client"

import { Skeleton } from "@/components/ui/skeleton"
import { DiagnosisView } from "@/components/shared/diagnosis-view"
import { useRunStatus } from "@/hooks/use-runs"

interface RunStatusViewProps {
  runId: string
}

export function RunStatusView({ runId }: RunStatusViewProps) {
  const { data: run, isLoading } = useRunStatus(runId)

  if (isLoading || !run) {
    return (
      <div className="space-y-4">
        <Skeleton className="h-6 w-40" />
        <Skeleton className="h-4 w-72" />
      </div>
    )
  }

  if (run.status === "pending") {
    return (
      <div className="flex flex-col items-center justify-center gap-3 py-16 text-center">
        <div className="text-4xl">🕐</div>
        <p className="font-medium">Waiting to start processing...</p>
        <p className="text-sm text-muted-foreground">
          The run is queued and will be analyzed shortly.
        </p>
      </div>
    )
  }

  if (run.status === "analyzing") {
    return (
      <div className="flex flex-col items-center justify-center gap-4 py-16 text-center">
        <div className="h-10 w-10 rounded-full border-4 border-primary border-t-transparent animate-spin" />
        <p className="font-medium">Analyzing results with AI...</p>
        <p className="text-sm text-muted-foreground">
          This may take a minute. The page updates automatically.
        </p>
        <div className="w-64 h-1.5 bg-muted rounded-full overflow-hidden">
          <div className="h-full bg-primary rounded-full animate-[progress_2s_ease-in-out_infinite]" />
        </div>
      </div>
    )
  }

  if (run.status === "failed") {
    return (
      <div className="flex flex-col items-center justify-center gap-3 py-16 text-center">
        <div className="text-4xl">❌</div>
        <p className="font-semibold text-destructive">Analysis failed</p>
        <p className="text-sm text-muted-foreground max-w-sm">
          The AI could not complete the analysis. Check that the .jtl file is valid and try
          creating a new run.
        </p>
      </div>
    )
  }

  return <DiagnosisView runId={runId} />
}
