import { Badge } from "@/components/ui/badge"
import { Card, CardContent, CardHeader, CardTitle } from "@/components/ui/card"
import type { Bottleneck } from "@/types/diagnosis"

const CONFIDENCE_STYLE: Record<string, string> = {
  high: "bg-red-100 text-red-700 dark:bg-red-900/30 dark:text-red-400",
  medium: "bg-yellow-100 text-yellow-700 dark:bg-yellow-900/30 dark:text-yellow-400",
  low: "bg-muted text-muted-foreground",
}

interface BottlenecksSectionProps {
  bottlenecks: Bottleneck[]
}

export function BottlenecksSection({ bottlenecks }: BottlenecksSectionProps) {
  if (bottlenecks.length === 0) return null

  return (
    <section className="space-y-3">
      <h3 className="text-base font-semibold">Bottleneck Hypotheses</h3>
      <div className="space-y-3">
        {bottlenecks.map((b, i) => (
          <Card key={i}>
            <CardHeader className="pb-3">
              <div className="flex items-center justify-between gap-2">
                <CardTitle className="text-sm font-semibold">{b.microservice}</CardTitle>
                <span
                  className={`inline-block px-2 py-0.5 rounded-full text-xs font-medium capitalize ${
                    CONFIDENCE_STYLE[b.confidence] ?? CONFIDENCE_STYLE.low
                  }`}
                >
                  {b.confidence} confidence
                </span>
              </div>
            </CardHeader>
            <CardContent className="space-y-3">
              {b.hypotheses
                .sort((a, z) => a.priority - z.priority)
                .map((h, j) => (
                  <div key={j} className="flex gap-3">
                    <span className="shrink-0 text-xs font-bold text-muted-foreground w-4">
                      {h.priority}.
                    </span>
                    <div className="space-y-0.5">
                      <p className="text-sm font-medium">{h.title}</p>
                      <p className="text-xs text-muted-foreground">{h.evidence}</p>
                    </div>
                  </div>
                ))}
            </CardContent>
          </Card>
        ))}
      </div>
    </section>
  )
}
