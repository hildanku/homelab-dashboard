
import { Card, CardContent, CardHeader, CardTitle } from "@/components/ui/card"
import { Badge } from "@/components/ui/badge"
import { useApi } from "@/hooks/use-api"
import { Monitor } from "lucide-react"

interface SystemData {
    os_name: string
    os_version: string
    kernel: string
    arch: string
    go_version: string
}

export function SystemInfo() {
    const { data, isLoading, isError, error } = useApi<SystemData>("http://localhost:5551/api/system")

    if (isLoading) {
        return (
            <Card>
                <CardHeader>
                    <CardTitle className="flex items-center gap-2">
                        <Monitor className="h-5 w-5" />
                        System Information
                    </CardTitle>
                </CardHeader>
                <CardContent>
                    <div className="animate-pulse space-y-2">
                        <div className="h-4 bg-muted rounded w-3/4"></div>
                        <div className="h-4 bg-muted rounded w-1/2"></div>
                        <div className="h-4 bg-muted rounded w-2/3"></div>
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
                        <Monitor className="h-5 w-5" />
                        System Information
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
                    <Monitor className="h-5 w-5" />
                    System Information
                </CardTitle>
            </CardHeader>
            <CardContent className="space-y-4">
                <div className="grid grid-cols-1 md:grid-cols-2 gap-4">
                    <div className="space-y-2">
                        <div className="flex items-center justify-between">
                            <span className="text-sm text-muted-foreground">OS</span>
                            <Badge variant="outline">{data?.os_name}</Badge>
                        </div>
                        <div className="flex items-center justify-between">
                            <span className="text-sm text-muted-foreground">Version</span>
                            <span className="text-sm font-mono">{data?.os_version}</span>
                        </div>
                        <div className="flex items-center justify-between">
                            <span className="text-sm text-muted-foreground">Architecture</span>
                            <Badge variant="secondary">{data?.arch}</Badge>
                        </div>
                    </div>
                    <div className="space-y-2">
                        <div className="flex items-center justify-between">
                            <span className="text-sm text-muted-foreground">Kernel</span>
                            <span className="text-sm font-mono">{data?.kernel}</span>
                        </div>
                        <div className="flex items-center justify-between">
                            <span className="text-sm text-muted-foreground">Go Version</span>
                            <Badge variant="outline">{data?.go_version}</Badge>
                        </div>
                    </div>
                </div>
            </CardContent>
        </Card>
    )
}
