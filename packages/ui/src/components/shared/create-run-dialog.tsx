"use client"

import { useState, useRef } from "react"
import { useForm } from "react-hook-form"
import { zodResolver } from "@hookform/resolvers/zod"
import { z } from "zod"
import { toast } from "sonner"
import { useRouter } from "next/navigation"
import {
  Dialog,
  DialogContent,
  DialogHeader,
  DialogTitle,
  DialogTrigger,
} from "@/components/ui/dialog"
import {
  Form,
  FormControl,
  FormField,
  FormItem,
  FormLabel,
  FormMessage,
} from "@/components/ui/form"
import { Input } from "@/components/ui/input"
import { Textarea } from "@/components/ui/textarea"
import { Button } from "@/components/ui/button"
import { useCreateRun } from "@/hooks/use-runs"
import { formatFileSize } from "@/lib/format"

const schema = z.object({
  name: z.string().min(1, "Name is required"),
  virtualUsers: z.number().int().positive("Must be a positive integer"),
  durationSeconds: z.number().int().positive("Must be a positive integer"),
  notes: z.string().optional(),
})

type FormValues = z.infer<typeof schema>

interface CreateRunDialogProps {
  suiteId: string
  trigger: React.ReactElement
}

export function CreateRunDialog({ suiteId, trigger }: CreateRunDialogProps) {
  const [open, setOpen] = useState(false)
  const [file, setFile] = useState<File | null>(null)
  const [fileError, setFileError] = useState<string | null>(null)
  const fileRef = useRef<HTMLInputElement>(null)
  const router = useRouter()
  const { mutateAsync, isPending } = useCreateRun(suiteId)

  const form = useForm<FormValues>({
    resolver: zodResolver(schema),
    defaultValues: { name: "", notes: "" },
  })

  function handleFileChange(e: React.ChangeEvent<HTMLInputElement>) {
    const selected = e.target.files?.[0]
    setFileError(null)
    if (!selected) {
      setFile(null)
      return
    }
    if (!selected.name.endsWith(".jtl")) {
      setFileError("Only .jtl files are accepted")
      setFile(null)
      e.target.value = ""
      return
    }
    setFile(selected)
  }

  async function onSubmit(values: FormValues) {
    if (!file) {
      setFileError("A .jtl file is required")
      return
    }
    try {
      const run = await mutateAsync({
        name: values.name,
        virtualUsers: values.virtualUsers,
        durationSeconds: values.durationSeconds,
        notes: values.notes || undefined,
        file,
      })
      toast.success("Run created — analysis will start shortly")
      setOpen(false)
      form.reset()
      setFile(null)
      router.push(`/runs/${run.id}`)
    } catch (err) {
      toast.error(err instanceof Error ? err.message : "Failed to create run")
    }
  }

  return (
    <Dialog open={open} onOpenChange={setOpen}>
      <DialogTrigger render={trigger} />
      <DialogContent className="sm:max-w-md">
        <DialogHeader>
          <DialogTitle>New Test Run</DialogTitle>
        </DialogHeader>
        <Form {...form}>
          <form onSubmit={form.handleSubmit(onSubmit)} className="space-y-4">
            <FormField
              control={form.control}
              name="name"
              render={({ field }) => (
                <FormItem>
                  <FormLabel>Name</FormLabel>
                  <FormControl>
                    <Input placeholder="e.g. Peak load — 500 VU" {...field} />
                  </FormControl>
                  <FormMessage />
                </FormItem>
              )}
            />

            <div className="grid grid-cols-2 gap-4">
              <FormField
                control={form.control}
                name="virtualUsers"
                render={({ field }) => (
                  <FormItem>
                    <FormLabel>Virtual Users</FormLabel>
                    <FormControl>
                      <Input
                        type="number"
                        min={1}
                        placeholder="500"
                        value={field.value ?? ""}
                        onChange={(e) =>
                          field.onChange(
                            e.target.value === "" ? undefined : Number(e.target.value)
                          )
                        }
                      />
                    </FormControl>
                    <FormMessage />
                  </FormItem>
                )}
              />
              <FormField
                control={form.control}
                name="durationSeconds"
                render={({ field }) => (
                  <FormItem>
                    <FormLabel>Duration (s)</FormLabel>
                    <FormControl>
                      <Input
                        type="number"
                        min={1}
                        placeholder="300"
                        value={field.value ?? ""}
                        onChange={(e) =>
                          field.onChange(
                            e.target.value === "" ? undefined : Number(e.target.value)
                          )
                        }
                      />
                    </FormControl>
                    <FormMessage />
                  </FormItem>
                )}
              />
            </div>

            <FormField
              control={form.control}
              name="notes"
              render={({ field }) => (
                <FormItem>
                  <FormLabel>Notes (optional)</FormLabel>
                  <FormControl>
                    <Textarea
                      className="resize-none"
                      rows={2}
                      placeholder="Any relevant context..."
                      {...field}
                    />
                  </FormControl>
                  <FormMessage />
                </FormItem>
              )}
            />

            <div className="space-y-2">
              <FormLabel>JTL File</FormLabel>
              <div
                className="rounded-md border border-dashed p-4 text-center cursor-pointer hover:bg-accent transition-colors"
                onClick={() => fileRef.current?.click()}
              >
                {file ? (
                  <div className="text-sm">
                    <p className="font-medium">{file.name}</p>
                    <p className="text-muted-foreground">{formatFileSize(file.size)}</p>
                  </div>
                ) : (
                  <p className="text-sm text-muted-foreground">
                    Click to select a <strong>.jtl</strong> file
                  </p>
                )}
              </div>
              <input
                ref={fileRef}
                type="file"
                accept=".jtl"
                className="hidden"
                onChange={handleFileChange}
              />
              {fileError && <p className="text-sm font-medium text-destructive">{fileError}</p>}
            </div>

            <div className="flex justify-end gap-2 pt-2">
              <Button type="button" variant="ghost" onClick={() => setOpen(false)}>
                Cancel
              </Button>
              <Button type="submit" disabled={isPending}>
                {isPending ? "Uploading..." : "Start Analysis"}
              </Button>
            </div>
          </form>
        </Form>
      </DialogContent>
    </Dialog>
  )
}
