import { Card, CardContent, CardHeader, CardTitle } from "@/components/ui/card"
import { Badge } from "@/components/ui/badge"
import { useApi } from "@/hooks/use-api"
import { Wifi, CheckCircle, XCircle, Clock } from "lucide-react"

interface PingData {
    url: string
    ok: boolean
    code: number
    latency_ms: number
}

export function PingStatus() {
    const { data, isLoading, isError, error } = useApi<PingData[]>("http://localhost:5551/api/ping/all")

    if (isLoading) {
        return (
            <Card>
                <CardHeader>
                    <CardTitle className="flex items-center gap-2">
                        <Wifi className="h-5 w-5" />
                        Network Status
                    </CardTitle>
                </CardHeader>
                <CardContent>
                    <div className="animate-pulse space-y-3">
                        <div className="flex items-center justify-between">
                            <div className="h-4 bg-muted rounded w-48"></div>
                            <div className="h-6 bg-muted rounded w-16"></div>
                        </div>
                    </div>
                </CardContent>
            </Card>
        )
    }

    if (isError) {
        return (
            <Card>
                <CardHeader>
                    <CardTitle className="flex items-center gap-2">
                        <Wifi className="h-5 w-5" />
                        Network Status
                    </CardTitle>
                </CardHeader>
                <CardContent>
                    <p className="text-destructive">Error: {(error as Error).message}</p>
                </CardContent>
            </Card>
        )
    }

    return (
        <Card>
            <CardHeader>
                <CardTitle className="flex items-center gap-2">
                    <Wifi className="h-5 w-5" />
                    Network Status
                </CardTitle>
            </CardHeader>
            <CardContent>
                <div className="space-y-4">
                    {data?.map((ping, index) => (
                        <div key={index} className="flex items-center justify-between p-3 rounded-lg border">
                            <div className="flex items-center gap-3">
                                {ping.ok ? (
                                    <CheckCircle className="h-5 w-5 text-green-500" />
                                ) : (
                                    <XCircle className="h-5 w-5 text-destructive" />
                                )}
                                <div className="flex flex-col">
                                    <span className="font-medium text-sm">{ping.url}</span>
                                    <span className="text-xs text-muted-foreground">Status Code: {ping.code || "N/A"}</span>
                                </div>
                            </div>
                            <div className="flex items-center gap-2">
                                <Badge variant={ping.ok ? "default" : "destructive"}>{ping.ok ? "Online" : "Offline"}</Badge>
                                {ping.ok && ping.latency_ms > 0 && (
                                    <Badge variant="outline" className="flex items-center gap-1">
                                        <Clock className="h-3 w-3" />
                                        {ping.latency_ms}ms
                                    </Badge>
                                )}
                            </div>
                        </div>
                    ))}
                </div>
            </CardContent>
        </Card>
    )
}

