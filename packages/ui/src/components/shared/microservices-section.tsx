"use client"

import { useState } from "react"
import { toast } from "sonner"
import { Badge } from "@/components/ui/badge"
import { Button } from "@/components/ui/button"
import { Card, CardContent, CardHeader, CardTitle } from "@/components/ui/card"
import {
  Dialog,
  DialogContent,
  DialogHeader,
  DialogTitle,
  DialogDescription,
} from "@/components/ui/dialog"
import { CreateMicroserviceDialog } from "@/components/shared/create-microservice-dialog"
import { EditMicroserviceDialog } from "@/components/shared/edit-microservice-dialog"
import { useDeleteMicroservice } from "@/hooks/use-microservices"
import type { Microservice } from "@/types/microservice"

interface MicroservicesSectionProps {
  suiteId: string
  microservices: Microservice[]
}

export function MicroservicesSection({ suiteId, microservices }: MicroservicesSectionProps) {
  const { mutateAsync: deleteMs, isPending: isDeleting } = useDeleteMicroservice(suiteId)
  const [deleteTarget, setDeleteTarget] = useState<Microservice | null>(null)

  async function handleDelete() {
    if (!deleteTarget) return
    try {
      await deleteMs(deleteTarget.id)
      toast.success(`"${deleteTarget.name}" removed`)
      setDeleteTarget(null)
    } catch (err) {
      toast.error(err instanceof Error ? err.message : "Failed to delete microservice")
    }
  }

  return (
    <section>
      <div className="flex items-center justify-between mb-4">
        <h2 className="text-lg font-semibold">Microservices</h2>
        <CreateMicroserviceDialog
          suiteId={suiteId}
          trigger={
            <Button variant="outline" size="sm">
              Add Microservice
            </Button>
          }
        />
      </div>

      {microservices.length === 0 ? (
        <div className="rounded-lg border border-dashed p-8 text-center text-muted-foreground">
          <p className="text-sm">No microservices registered yet.</p>
          <p className="text-xs mt-1">Add microservices to give the AI context for diagnosis.</p>
        </div>
      ) : (
        <div className="grid gap-3 sm:grid-cols-2 lg:grid-cols-3">
          {microservices.map((ms) => (
            <Card key={ms.id} className="relative">
              <CardHeader className="pb-2">
                <div className="flex items-start justify-between gap-2">
                  <CardTitle className="text-base leading-tight">{ms.name}</CardTitle>
                  <Badge variant="secondary" className="shrink-0 text-xs">
                    {ms.language}
                  </Badge>
                </div>
                <p className="text-xs text-muted-foreground line-clamp-2">{ms.description}</p>
              </CardHeader>
              <CardContent className="space-y-3">
                <div>
                  <p className="text-xs font-medium text-muted-foreground mb-1">Endpoints</p>
                  <div className="flex flex-wrap gap-1">
                    {ms.mainEndpoints.slice(0, 3).map((ep) => (
                      <span
                        key={ep}
                        className="inline-block bg-muted px-2 py-0.5 rounded text-xs font-mono truncate max-w-40"
                        title={ep}
                      >
                        {ep}
                      </span>
                    ))}
                    {ms.mainEndpoints.length > 3 && (
                      <span className="text-xs text-muted-foreground">
                        +{ms.mainEndpoints.length - 3} more
                      </span>
                    )}
                  </div>
                </div>

                {(ms.sloLatencyP99Ms || ms.sloErrorRatePct) && (
                  <div className="flex flex-wrap gap-2 text-xs text-muted-foreground">
                    {ms.sloLatencyP99Ms && <span>p99 &lt; {ms.sloLatencyP99Ms}ms</span>}
                    {ms.sloErrorRatePct && <span>errors &lt; {ms.sloErrorRatePct}%</span>}
                  </div>
                )}

                <div className="flex gap-1 pt-1">
                  <EditMicroserviceDialog
                    suiteId={suiteId}
                    microservice={ms}
                    trigger={
                      <Button variant="ghost" size="sm" className="h-7 text-xs">
                        Edit
                      </Button>
                    }
                  />
                  <Button
                    variant="ghost"
                    size="sm"
                    className="h-7 text-xs text-destructive hover:text-destructive hover:bg-destructive/10"
                    onClick={() => setDeleteTarget(ms)}
                  >
                    Delete
                  </Button>
                </div>
              </CardContent>
            </Card>
          ))}
        </div>
      )}

      <Dialog open={!!deleteTarget} onOpenChange={(o) => !o && setDeleteTarget(null)}>
        <DialogContent className="sm:max-w-sm">
          <DialogHeader>
            <DialogTitle>Remove Microservice</DialogTitle>
            <DialogDescription>
              Remove &ldquo;{deleteTarget?.name}&rdquo; from this suite?
            </DialogDescription>
          </DialogHeader>
          <div className="flex justify-end gap-2 pt-2">
            <Button variant="ghost" onClick={() => setDeleteTarget(null)}>
              Cancel
            </Button>
            <Button variant="destructive" disabled={isDeleting} onClick={handleDelete}>
              {isDeleting ? "Removing..." : "Remove"}
            </Button>
          </div>
        </DialogContent>
      </Dialog>
    </section>
  )
}
