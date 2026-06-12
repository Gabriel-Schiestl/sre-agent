"use client"

import { useState, useEffect } from "react"
import { useForm } from "react-hook-form"
import { zodResolver } from "@hookform/resolvers/zod"
import { z } from "zod"
import { toast } from "sonner"
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
import { useUpdateSuite } from "@/hooks/use-suites"
import type { TestSuite } from "@/types/suite"

const schema = z.object({
  name: z.string().min(3, "Name must be at least 3 characters"),
  description: z.string().min(1, "Description is required"),
})

type FormValues = z.infer<typeof schema>

interface EditSuiteDialogProps {
  suite: TestSuite
  trigger: React.ReactElement
}

export function EditSuiteDialog({ suite, trigger }: EditSuiteDialogProps) {
  const [open, setOpen] = useState(false)
  const { mutateAsync, isPending } = useUpdateSuite(suite.id)

  const form = useForm<FormValues>({
    resolver: zodResolver(schema),
    defaultValues: { name: suite.name, description: suite.description },
  })

  useEffect(() => {
    if (open) {
      form.reset({ name: suite.name, description: suite.description })
    }
  }, [open, suite, form])

  async function onSubmit(values: FormValues) {
    try {
      await mutateAsync(values)
      toast.success("Suite updated")
      setOpen(false)
    } catch (err) {
      toast.error(err instanceof Error ? err.message : "Failed to update suite")
    }
  }

  return (
    <Dialog open={open} onOpenChange={setOpen}>
      <DialogTrigger render={trigger} />
      <DialogContent className="sm:max-w-md">
        <DialogHeader>
          <DialogTitle>Edit Suite</DialogTitle>
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
                    <Input {...field} />
                  </FormControl>
                  <FormMessage />
                </FormItem>
              )}
            />
            <FormField
              control={form.control}
              name="description"
              render={({ field }) => (
                <FormItem>
                  <FormLabel>Description</FormLabel>
                  <FormControl>
                    <Textarea className="resize-none" rows={3} {...field} />
                  </FormControl>
                  <FormMessage />
                </FormItem>
              )}
            />
            <div className="flex justify-end gap-2 pt-2">
              <Button type="button" variant="ghost" onClick={() => setOpen(false)}>
                Cancel
              </Button>
              <Button type="submit" disabled={isPending}>
                {isPending ? "Saving..." : "Save Changes"}
              </Button>
            </div>
          </form>
        </Form>
      </DialogContent>
    </Dialog>
  )
}
