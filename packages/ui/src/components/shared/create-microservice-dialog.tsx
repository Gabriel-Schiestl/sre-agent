"use client"

import { useState } from "react"
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
import { useCreateMicroservice } from "@/hooks/use-microservices"

interface CreateMicroserviceDialogProps {
  suiteId: string
  trigger: React.ReactElement
}

export function CreateMicroserviceDialog({ suiteId, trigger }: CreateMicroserviceDialogProps) {
  const [open, setOpen] = useState(false)
  const { mutateAsync, isPending } = useCreateMicroservice(suiteId)

  const form = useForm<MicroserviceFormValues>({
    resolver: zodResolver(microserviceSchema),
    defaultValues: {
      name: "",
      description: "",
      language: "",
      mainEndpoints: [{ value: "" }],
      cpuLimit: "",
      memoryLimit: "",
    },
  })

  async function onSubmit(values: MicroserviceFormValues) {
    try {
      await mutateAsync({
        name: values.name,
        description: values.description,
        language: values.language,
        mainEndpoints: values.mainEndpoints.map((e) => e.value),
        cpuLimit: values.cpuLimit || undefined,
        memoryLimit: values.memoryLimit || undefined,
        sloLatencyP99Ms: values.sloLatencyP99Ms,
        sloErrorRatePct: values.sloErrorRatePct,
      })
      toast.success("Microservice added")
      setOpen(false)
      form.reset()
    } catch (err) {
      toast.error(err instanceof Error ? err.message : "Failed to add microservice")
    }
  }

  return (
    <Dialog open={open} onOpenChange={setOpen}>
      <DialogTrigger render={trigger} />
      <DialogContent className="sm:max-w-lg max-h-[90vh] overflow-y-auto">
        <DialogHeader>
          <DialogTitle>Add Microservice</DialogTitle>
        </DialogHeader>
        <MicroserviceForm
          form={form}
          isPending={isPending}
          onSubmit={onSubmit}
          onCancel={() => setOpen(false)}
          submitLabel="Add Microservice"
        />
      </DialogContent>
    </Dialog>
  )
}
