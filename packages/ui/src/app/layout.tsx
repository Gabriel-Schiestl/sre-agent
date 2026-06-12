import type { Metadata } from "next"
import { Geist } from "next/font/google"
import Link from "next/link"
import "./globals.css"
import { QueryProvider } from "@/providers/query-provider"
import { Toaster } from "@/components/ui/sonner"

const geist = Geist({ subsets: ["latin"], variable: "--font-geist-sans" })

export const metadata: Metadata = {
  title: "Stress Test Intelligence Platform",
  description: "AI-powered stress test analysis",
}

export default function RootLayout({ children }: { children: React.ReactNode }) {
  return (
    <html lang="en" className={`${geist.variable} h-full antialiased`}>
      <body className="h-full flex bg-background text-foreground">
        <aside className="w-56 shrink-0 border-r border-border flex flex-col">
          <div className="h-14 flex items-center px-5 border-b border-border">
            <span className="font-semibold text-sm tracking-tight">SRE Agent</span>
          </div>
          <nav className="flex-1 p-3 space-y-1">
            <Link
              href="/suites"
              className="flex items-center gap-2 px-3 py-2 rounded-md text-sm text-muted-foreground hover:text-foreground hover:bg-accent transition-colors"
            >
              Test Suites
            </Link>
          </nav>
        </aside>
        <main className="flex-1 overflow-auto">
          <QueryProvider>{children}</QueryProvider>
        </main>
        <Toaster richColors closeButton />
      </body>
    </html>
  )
}
