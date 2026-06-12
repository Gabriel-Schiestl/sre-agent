import { SuiteDetailView } from "@/components/shared/suite-detail-view"

export default async function SuiteDetailPage({
  params,
}: {
  params: Promise<{ id: string }>
}) {
  const { id } = await params
  return <SuiteDetailView id={id} />
}
