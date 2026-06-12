import { Card, CardContent, CardHeader, CardTitle } from "@/components/ui/card"
import type { ErrorCategory } from "@/types/diagnosis"

const SEVERITY_STYLE: Record<string, string> = {
  critical: "bg-red-100 text-red-700 dark:bg-red-900/30 dark:text-red-400",
  high: "bg-orange-100 text-orange-700 dark:bg-orange-900/30 dark:text-orange-400",
  medium: "bg-yellow-100 text-yellow-700 dark:bg-yellow-900/30 dark:text-yellow-400",
  low: "bg-muted text-muted-foreground",
}

interface ErrorPlanSectionProps {
  errorPlan: ErrorCategory[]
}

export function ErrorPlanSection({ errorPlan }: ErrorPlanSectionProps) {
  if (errorPlan.length === 0) return null

  return (
    <section className="space-y-3">
      <h3 className="text-base font-semibold">Error Categories</h3>
      <div className="grid gap-3 sm:grid-cols-2">
        {errorPlan.map((category, i) => (
          <Card key={i}>
            <CardHeader className="pb-2">
              <div className="flex items-start justify-between gap-2">
                <CardTitle className="text-sm leading-tight">{category.category}</CardTitle>
                <span
                  className={`shrink-0 inline-block px-2 py-0.5 rounded-full text-xs font-medium capitalize ${
                    SEVERITY_STYLE[category.severity] ?? SEVERITY_STYLE.low
                  }`}
                >
                  {category.severity}
                </span>
              </div>
            </CardHeader>
            <CardContent className="space-y-3">
              <p className="text-xs text-muted-foreground">{category.description}</p>

              <div className="flex items-center gap-2">
                <span className="text-2xl font-bold tabular-nums">
                  {category.occurrences.toLocaleString()}
                </span>
                <span className="text-xs text-muted-foreground">occurrences</span>
              </div>

              {category.affectedEndpoints.length > 0 && (
                <div className="flex flex-wrap gap-1">
                  {category.affectedEndpoints.map((ep) => (
                    <span
                      key={ep}
                      className="inline-block bg-muted px-2 py-0.5 rounded text-xs font-mono truncate max-w-[200px]"
                      title={ep}
                    >
                      {ep}
                    </span>
                  ))}
                </div>
              )}
            </CardContent>
          </Card>
        ))}
      </div>
    </section>
  )
}
