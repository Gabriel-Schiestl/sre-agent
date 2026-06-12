"use client"

import { useRouter } from "next/navigation"
import { Button } from "@/components/ui/button"
import {
  Table,
  TableBody,
  TableCell,
  TableHead,
  TableHeader,
  TableRow,
} from "@/components/ui/table"
import { CreateRunDialog } from "@/components/shared/create-run-dialog"
import { formatDateTime } from "@/lib/format"
import type { TestRun, RunStatus } from "@/types/run"

const STATUS_CONFIG: Record<RunStatus, { label: string; className: string }> = {
  pending: { label: "Pending", className: "bg-muted text-muted-foreground" },
  analyzing: {
    label: "Analyzing",
    className: "bg-yellow-100 text-yellow-800 dark:bg-yellow-900/30 dark:text-yellow-400",
  },
  done: { label: "Done", className: "bg-green-100 text-green-800 dark:bg-green-900/30 dark:text-green-400" },
  failed: { label: "Failed", className: "bg-red-100 text-red-800 dark:bg-red-900/30 dark:text-red-400" },
}

function StatusBadge({ status }: { status: RunStatus }) {
  const { label, className } = STATUS_CONFIG[status]
  return (
    <span
      className={`inline-flex items-center gap-1.5 px-2 py-0.5 rounded-full text-xs font-medium ${className}`}
    >
      {status === "analyzing" && (
        <span className="h-1.5 w-1.5 rounded-full bg-yellow-500 animate-pulse" />
      )}
      {label}
    </span>
  )
}

interface RunsSectionProps {
  suiteId: string
  runs: TestRun[]
}

export function RunsSection({ suiteId, runs }: RunsSectionProps) {
  const router = useRouter()

  return (
    <section>
      <div className="flex items-center justify-between mb-4">
        <h2 className="text-lg font-semibold">Test Runs</h2>
        <CreateRunDialog
          suiteId={suiteId}
          trigger={
            <Button variant="outline" size="sm">
              New Run
            </Button>
          }
        />
      </div>

      {runs.length === 0 ? (
        <div className="rounded-lg border border-dashed p-8 text-center text-muted-foreground">
          <p className="text-sm">No runs executed yet.</p>
          <p className="text-xs mt-1">Upload a .jtl file to start your first analysis.</p>
        </div>
      ) : (
        <Table>
          <TableHeader>
            <TableRow>
              <TableHead>Name</TableHead>
              <TableHead>Virtual Users</TableHead>
              <TableHead>Duration</TableHead>
              <TableHead>Status</TableHead>
              <TableHead>Created</TableHead>
            </TableRow>
          </TableHeader>
          <TableBody>
            {runs.map((run) => (
              <TableRow
                key={run.id}
                className="cursor-pointer"
                onClick={() => router.push(`/runs/${run.id}`)}
              >
                <TableCell className="font-medium">{run.name}</TableCell>
                <TableCell>{run.virtualUsers.toLocaleString()}</TableCell>
                <TableCell>{run.durationSeconds}s</TableCell>
                <TableCell>
                  <StatusBadge status={run.status} />
                </TableCell>
                <TableCell className="text-muted-foreground whitespace-nowrap">
                  {formatDateTime(run.createdAt)}
                </TableCell>
              </TableRow>
            ))}
          </TableBody>
        </Table>
      )}
    </section>
  )
}
