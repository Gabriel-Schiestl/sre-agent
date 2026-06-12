interface NextStepsSectionProps {
  nextSteps: string[]
}

export function NextStepsSection({ nextSteps }: NextStepsSectionProps) {
  if (nextSteps.length === 0) return null

  return (
    <section className="space-y-3">
      <h3 className="text-base font-semibold">Recommended Next Steps</h3>
      <ol className="space-y-2">
        {nextSteps.map((step, i) => (
          <li key={i} className="flex gap-3 text-sm">
            <span className="shrink-0 flex h-5 w-5 items-center justify-center rounded-full bg-primary text-primary-foreground text-xs font-bold">
              {i + 1}
            </span>
            <span className="leading-5">{step}</span>
          </li>
        ))}
      </ol>
    </section>
  )
}
