"use client"

import Link from "next/link"
import { Button } from "@/components/ui/button"
import { Skeleton } from "@/components/ui/skeleton"
import { Separator } from "@/components/ui/separator"
import { EditSuiteDialog } from "@/components/shared/edit-suite-dialog"
import { MicroservicesSection } from "@/components/shared/microservices-section"
import { RunsSection } from "@/components/shared/runs-section"
import { ErrorState } from "@/components/shared/error-state"
import { useSuite } from "@/hooks/use-suites"

interface SuiteDetailViewProps {
  id: string
}

export function SuiteDetailView({ id }: SuiteDetailViewProps) {
  const { data: suite, isLoading, isError, refetch } = useSuite(id)

  if (isLoading) {
    return (
      <div className="p-8 max-w-5xl mx-auto space-y-8">
        <div className="space-y-2">
          <Skeleton className="h-6 w-48" />
          <Skeleton className="h-8 w-72" />
          <Skeleton className="h-4 w-96" />
        </div>
        <Skeleton className="h-48 w-full" />
        <Skeleton className="h-64 w-full" />
      </div>
    )
  }

  if (isError || !suite) {
    return (
      <div className="p-8 max-w-5xl mx-auto">
        <ErrorState
          title="Suite not found"
          message="Could not load the suite. It may have been deleted or the API is unavailable."
          onRetry={refetch}
        />
      </div>
    )
  }

  return (
    <div className="p-8 max-w-5xl mx-auto space-y-8">
      {/* Breadcrumb */}
      <nav className="flex items-center gap-2 text-sm text-muted-foreground">
        <Link href="/suites" className="hover:text-foreground transition-colors">
          Test Suites
        </Link>
        <span>/</span>
        <span className="text-foreground font-medium">{suite.name}</span>
      </nav>

      {/* Suite Header */}
      <div className="flex items-start justify-between gap-4">
        <div>
          <h1 className="text-2xl font-semibold tracking-tight">{suite.name}</h1>
          <p className="text-muted-foreground mt-1">{suite.description}</p>
        </div>
        <EditSuiteDialog
          suite={suite}
          trigger={
            <Button variant="outline" size="sm">
              Edit Suite
            </Button>
          }
        />
      </div>

      <Separator />

      <MicroservicesSection suiteId={id} microservices={suite.microservices} />

      <Separator />

      <RunsSection suiteId={id} runs={suite.runs} />
    </div>
  )
}
