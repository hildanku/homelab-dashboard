import { createFileRoute } from '@tanstack/react-router'

export const Route = createFileRoute('/')({
    component: RouteComponent,
})

function RouteComponent() {
    return (
        <div className="relative flex min-h-svh flex-col items-center justify-center overflow-hidden bg-background">
            <div className="pointer-events-none absolute inset-0 bg-[radial-gradient(ellipse_at_center,var(--tw-gradient-stops))] from-primary/10 via-background to-background" />

            <div className="relative z-10 flex flex-col items-center gap-6 px-6 text-center">
                <span className="rounded-full border border-border px-4 py-1 text-xs font-medium uppercase tracking-widest text-muted-foreground">
                    In Development
                </span>

                <h1 className="text-4xl font-bold tracking-tight text-balance sm:text-6xl">
                    Coming Soon
                </h1>

                <p className="max-w-md text-muted-foreground text-pretty">
                    Something new is on the way. We are building a better homelab
                    experience. Stay tuned.
                </p>

                <div className="mt-4 font-mono text-xs uppercase tracking-[0.3em] text-muted-foreground/60">
                    by Svega Labs
                </div>
            </div>
        </div>
    )
}
