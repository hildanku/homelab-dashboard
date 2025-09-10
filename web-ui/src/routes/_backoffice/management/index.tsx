import { createFileRoute } from '@tanstack/react-router'
import { SystemInfo } from "@/components/custom/system-info"
import { MetricsCard } from "@/components/custom/metrics-card"
import { ServicesStatus } from "@/components/custom/services-status"
import { PingStatus } from "@/components/custom/ping-status"
import { Server, RefreshCw } from "lucide-react"
export const Route = createFileRoute('/_backoffice/management/')({
    component: App,
})

function App() {
    const handleRefresh = () => {
        window.location.reload()
    }

    return (
        <div className="min-h-screen bg-background">
            <div className="container mx-auto p-6 space-y-6">
                {/* Header */}
                <div className="flex items-center justify-between">
                    <div className="flex items-center gap-3">
                        <Server className="h-8 w-8 text-primary" />
                        <div>
                            <h1 className="text-3xl font-bold text-balance">Homelab Dashboard</h1>
                            <p className="text-muted-foreground">Monitor your server status and metrics</p>
                        </div>
                    </div>
                    <button
                        onClick={handleRefresh}
                        className="flex items-center gap-2 px-4 py-2 bg-primary text-primary-foreground rounded-lg hover:bg-primary/90 transition-colors"
                    >
                        <RefreshCw className="h-4 w-4" />
                        Refresh
                    </button>
                </div>

                {/* System Information */}
                <SystemInfo />

                {/* Metrics Cards */}
                <MetricsCard />

                {/* Services and Network Status */}
                <div className="grid grid-cols-1 lg:grid-cols-2 gap-6">
                    <ServicesStatus />
                    <PingStatus />
                </div>

                {/* Footer */}
                <div className="text-center text-sm text-muted-foreground pt-6 border-t">
                    <p>Dashboard updates automatically every 3-5 seconds</p>
                </div>
            </div>
        </div>
    )
}
