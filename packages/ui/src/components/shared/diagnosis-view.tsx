"use client"

import { Skeleton } from "@/components/ui/skeleton"
import { Separator } from "@/components/ui/separator"
import { ErrorState } from "@/components/shared/error-state"
import { ErrorPlanSection } from "@/components/shared/error-plan-section"
import { BottlenecksSection } from "@/components/shared/bottlenecks-section"
import { NextStepsSection } from "@/components/shared/next-steps-section"
import { useDiagnosis } from "@/hooks/use-diagnosis"
import { formatDateTime } from "@/lib/format"

interface DiagnosisViewProps {
  runId: string
}

export function DiagnosisView({ runId }: DiagnosisViewProps) {
  const { data: diagnosis, isLoading, isError, refetch } = useDiagnosis(runId)

  if (isLoading) {
    return (
      <div className="space-y-6">
        <Skeleton className="h-6 w-48" />
        <div className="grid gap-3 sm:grid-cols-2">
          <Skeleton className="h-36 w-full" />
          <Skeleton className="h-36 w-full" />
          <Skeleton className="h-36 w-full" />
        </div>
        <Skeleton className="h-48 w-full" />
      </div>
    )
  }

  if (isError || !diagnosis) {
    return (
      <ErrorState
        title="Failed to load diagnosis"
        message="The analysis result could not be retrieved."
        onRetry={refetch}
      />
    )
  }

  return (
    <div className="space-y-8">
      <div className="flex items-center justify-between">
        <h2 className="text-lg font-semibold">AI Diagnosis</h2>
        <p className="text-xs text-muted-foreground">{formatDateTime(diagnosis.createdAt)}</p>
      </div>

      <ErrorPlanSection errorPlan={diagnosis.errorPlan} />
      <Separator />
      <BottlenecksSection bottlenecks={diagnosis.bottlenecks} />
      <Separator />
      <NextStepsSection nextSteps={diagnosis.nextSteps} />
    </div>
  )
}
