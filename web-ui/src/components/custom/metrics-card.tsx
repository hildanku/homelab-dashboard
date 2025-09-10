import { Card, CardContent, CardHeader, CardTitle } from "@/components/ui/card"
import { Progress } from "@/components/ui/progress"
import { useApi } from "@/hooks/use-api"
import { formatBytes, formatPercentage } from "@/lib/utils"
import { Activity, HardDrive, MemoryStick } from "lucide-react"

interface MetricsData {
    cpu_usage: number
    ram_total: number
    ram_free: number
    ram_used: number
    ram_usage: number
    disk_total: number
    disk_free: number
    disk_used: number
    disk_usage: number
}

export function MetricsCard() {
    const {
        data,
        isLoading,
        isError,
        error,
    } = useApi<MetricsData>("http://localhost:5551/api/metrics-v2", 3000)

    if (isLoading) {
        return (
            <div className="grid grid-cols-1 md:grid-cols-3 gap-6">
                {[1, 2, 3].map((i) => (
                    <Card key={i}>
                        <CardHeader>
                            <CardTitle className="animate-pulse bg-muted h-6 rounded w-24"></CardTitle>
                        </CardHeader>
                        <CardContent>
                            <div className="animate-pulse space-y-3">
                                <div className="h-4 bg-muted rounded"></div>
                                <div className="h-6 bg-muted rounded w-16"></div>
                            </div>
                        </CardContent>
                    </Card>
                ))}
            </div>
        )
    }

    if (isError) {
        return (
            <Card>
                <CardHeader>
                    <CardTitle>System Metrics</CardTitle>
                </CardHeader>
                <CardContent>
                    <p className="text-destructive">
                        Error loading metrics: {(error as Error).message}
                    </p>
                </CardContent>
            </Card>
        )
    }

    return (
        <div className="grid grid-cols-1 md:grid-cols-3 gap-6">
            <Card>
                <CardHeader className="pb-3">
                    <CardTitle className="flex items-center gap-2 text-lg">
                        <Activity className="h-5 w-5" />
                        CPU Usage
                    </CardTitle>
                </CardHeader>
                <CardContent>
                    <div className="space-y-3">
                        <div className="flex items-center justify-between">
                            <span className="text-2xl font-bold">
                                {formatPercentage(data?.cpu_usage || 0)}
                            </span>
                            <span className="text-sm text-muted-foreground">Current</span>
                        </div>
                        <Progress value={data?.cpu_usage || 0} className="h-2" />
                    </div>
                </CardContent>
            </Card>

            <Card>
                <CardHeader className="pb-3">
                    <CardTitle className="flex items-center gap-2 text-lg">
                        <MemoryStick className="h-5 w-5" />
                        Memory
                    </CardTitle>
                </CardHeader>
                <CardContent>
                    <div className="space-y-3">
                        <div className="flex items-center justify-between">
                            <span className="text-2xl font-bold">
                                {formatPercentage(data?.ram_usage || 0)}
                            </span>
                            <span className="text-sm text-muted-foreground">
                                {formatBytes(data?.ram_used || 0)} /{" "}
                                {formatBytes(data?.ram_total || 0)}
                            </span>
                        </div>
                        <Progress value={data?.ram_usage || 0} className="h-2" />
                    </div>
                </CardContent>
            </Card>

            <Card>
                <CardHeader className="pb-3">
                    <CardTitle className="flex items-center gap-2 text-lg">
                        <HardDrive className="h-5 w-5" />
                        Storage
                    </CardTitle>
                </CardHeader>
                <CardContent>
                    <div className="space-y-3">
                        <div className="flex items-center justify-between">
                            <span className="text-2xl font-bold">
                                {formatPercentage(data?.disk_usage || 0)}
                            </span>
                            <span className="text-sm text-muted-foreground">
                                {formatBytes(data?.disk_used || 0)} /{" "}
                                {formatBytes(data?.disk_total || 0)}
                            </span>
                        </div>
                        <Progress value={data?.disk_usage || 0} className="h-2" />
                    </div>
                </CardContent>
            </Card>
        </div>
    )
}

