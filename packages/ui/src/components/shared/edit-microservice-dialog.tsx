"use client"

import { useState, useEffect } from "react"
import { useForm } from "react-hook-form"
import { zodResolver } from "@hookform/resolvers/zod"
import { toast } from "sonner"
import {
  Dialog,
  DialogContent,
  DialogHeader,
  DialogTitle,
  DialogTrigger,
} from "@/components/ui/dialog"
import {
  MicroserviceForm,
  microserviceSchema,
  type MicroserviceFormValues,
} from "@/components/shared/microservice-form"
import { useUpdateMicroservice } from "@/hooks/use-microservices"
import type { Microservice } from "@/types/microservice"

interface EditMicroserviceDialogProps {
  suiteId: string
  microservice: Microservice
  trigger: React.ReactElement
}

export function EditMicroserviceDialog({
  suiteId,
  microservice,
  trigger,
}: EditMicroserviceDialogProps) {
  const [open, setOpen] = useState(false)
  const { mutateAsync, isPending } = useUpdateMicroservice(suiteId)

  const form = useForm<MicroserviceFormValues>({
    resolver: zodResolver(microserviceSchema),
    defaultValues: toFormValues(microservice),
  })

  useEffect(() => {
    if (open) form.reset(toFormValues(microservice))
  }, [open, microservice, form])

  async function onSubmit(values: MicroserviceFormValues) {
    try {
      await mutateAsync({
        id: microservice.id,
        payload: {
          name: values.name,
          description: values.description,
          language: values.language,
          mainEndpoints: values.mainEndpoints.map((e) => e.value),
          cpuLimit: values.cpuLimit || undefined,
          memoryLimit: values.memoryLimit || undefined,
          sloLatencyP99Ms: values.sloLatencyP99Ms,
          sloErrorRatePct: values.sloErrorRatePct,
        },
      })
      toast.success("Microservice updated")
      setOpen(false)
    } catch (err) {
      toast.error(err instanceof Error ? err.message : "Failed to update microservice")
    }
  }

  return (
    <Dialog open={open} onOpenChange={setOpen}>
      <DialogTrigger render={trigger} />
      <DialogContent className="sm:max-w-lg max-h-[90vh] overflow-y-auto">
        <DialogHeader>
          <DialogTitle>Edit Microservice</DialogTitle>
        </DialogHeader>
        <MicroserviceForm
          form={form}
          isPending={isPending}
          onSubmit={onSubmit}
          onCancel={() => setOpen(false)}
          submitLabel="Save Changes"
        />
      </DialogContent>
    </Dialog>
  )
}

function toFormValues(ms: Microservice): MicroserviceFormValues {
  return {
    name: ms.name,
    description: ms.description,
    language: ms.language,
    mainEndpoints: ms.mainEndpoints.map((v) => ({ value: v })),
    cpuLimit: ms.cpuLimit ?? "",
    memoryLimit: ms.memoryLimit ?? "",
    sloLatencyP99Ms: ms.sloLatencyP99Ms,
    sloErrorRatePct: ms.sloErrorRatePct,
  }
}
