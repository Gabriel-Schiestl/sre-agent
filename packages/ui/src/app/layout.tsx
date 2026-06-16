import type { Metadata } from "next"
import { Geist } from "next/font/google"
import { Activity } from "lucide-react"
import "./globals.css"
import { QueryProvider } from "@/providers/query-provider"
import { Toaster } from "@/components/ui/sonner"
import { SidebarNav } from "@/components/shared/sidebar-nav"

const geist = Geist({ subsets: ["latin"], variable: "--font-geist-sans" })

export const metadata: Metadata = {
  title: "Stress Test Intelligence Platform",
  description: "AI-powered stress test analysis",
}

export default function RootLayout({ children }: { children: React.ReactNode }) {
  return (
    <html lang="en" className={`${geist.variable} h-full antialiased`}>
      <body className="h-full flex bg-background text-foreground">
        <aside className="w-60 shrink-0 border-r border-sidebar-border flex flex-col bg-sidebar">
          <div className="h-14 flex items-center px-4 border-b border-sidebar-border gap-3">
            <div className="flex h-8 w-8 items-center justify-center rounded-lg bg-primary shadow-sm shrink-0">
              <Activity className="size-4 text-primary-foreground" />
            </div>
            <div className="min-w-0">
              <p className="font-semibold text-sm text-sidebar-foreground tracking-tight leading-none">SRE Agent</p>
              <p className="text-[10px] text-muted-foreground mt-0.5 leading-none">AI Analysis</p>
            </div>
          </div>
          <nav className="flex-1 px-2 py-3 space-y-0.5">
            <SidebarNav />
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
