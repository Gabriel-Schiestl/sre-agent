"use client"

import { Button } from "@/components/ui/button"
import { CreateSuiteDialog } from "@/components/shared/create-suite-dialog"
import { SuitesTable } from "@/components/shared/suites-table"
import { ErrorState } from "@/components/shared/error-state"
import { useSuites } from "@/hooks/use-suites"

export default function SuitesPage() {
  const { data: suites, isLoading, isError, refetch } = useSuites()

  return (
    <div className="p-8 max-w-5xl mx-auto">
      <div className="flex items-center justify-between mb-6">
        <div>
          <h1 className="text-2xl font-semibold tracking-tight">Test Suites</h1>
          <p className="text-sm text-muted-foreground mt-1">
            Manage your recurring stress test scenarios
          </p>
        </div>
        <CreateSuiteDialog trigger={<Button>New Suite</Button>} />
      </div>

      {isError ? (
        <ErrorState
          title="Failed to load suites"
          message="Could not connect to the API."
          onRetry={refetch}
        />
      ) : (
        <SuitesTable suites={suites ?? []} isLoading={isLoading} />
      )}
    </div>
  )
}
