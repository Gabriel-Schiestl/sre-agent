"use client"

import Link from "next/link"
import { Skeleton } from "@/components/ui/skeleton"
import { Separator } from "@/components/ui/separator"
import { RunStatusView } from "@/components/shared/run-status-view"
import { ErrorState } from "@/components/shared/error-state"
import { useRun } from "@/hooks/use-runs"
import { useSuite } from "@/hooks/use-suites"
import { formatDateTime } from "@/lib/format"

interface RunDetailViewProps {
  id: string
}

function MetaItem({ label, value }: { label: string; value: string }) {
  return (
    <div>
      <p className="text-xs text-muted-foreground">{label}</p>
      <p className="text-sm font-medium">{value}</p>
    </div>
  )
}

export function RunDetailView({ id }: RunDetailViewProps) {
  const { data: run, isLoading, isError, refetch } = useRun(id)
  const { data: suite } = useSuite(run?.testSuiteId ?? "")

  if (isLoading) {
    return (
      <div className="p-8 max-w-5xl mx-auto space-y-6">
        <Skeleton className="h-5 w-64" />
        <Skeleton className="h-8 w-80" />
        <Skeleton className="h-20 w-full" />
      </div>
    )
  }

  if (isError || !run) {
    return (
      <div className="p-8 max-w-5xl mx-auto">
        <ErrorState title="Run not found" onRetry={refetch} />
      </div>
    )
  }

  return (
    <div className="p-8 max-w-5xl mx-auto space-y-8">
      {/* Breadcrumb */}
      <nav className="flex items-center gap-2 text-sm text-muted-foreground flex-wrap">
        <Link href="/suites" className="hover:text-foreground transition-colors">
          Test Suites
        </Link>
        <span>/</span>
        {suite ? (
          <Link
            href={`/suites/${suite.id}`}
            className="hover:text-foreground transition-colors"
          >
            {suite.name}
          </Link>
        ) : (
          <span>Suite</span>
        )}
        <span>/</span>
        <span className="text-foreground font-medium">{run.name}</span>
      </nav>

      {/* Run header */}
      <div>
        <h1 className="text-2xl font-semibold tracking-tight">{run.name}</h1>
        {run.notes && <p className="text-muted-foreground mt-1">{run.notes}</p>}
      </div>

      {/* Meta grid */}
      <div className="grid grid-cols-2 sm:grid-cols-4 gap-4 p-4 rounded-lg bg-muted/50">
        <MetaItem label="Virtual Users" value={run.virtualUsers.toLocaleString()} />
        <MetaItem label="Duration" value={`${run.durationSeconds}s`} />
        <MetaItem label="Created" value={formatDateTime(run.createdAt)} />
        <MetaItem label="Suite" value={suite?.name ?? "—"} />
      </div>

      <Separator />

      <RunStatusView runId={id} />
    </div>
  )
}
