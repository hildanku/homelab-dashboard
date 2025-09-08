
import { Card, CardContent, CardHeader, CardTitle } from "@/components/ui/card"
import { Badge } from "@/components/ui/badge"
import { useApi } from "@/hooks/use-api"
import { Settings, CheckCircle, XCircle } from "lucide-react"

interface ServiceData {
    name: string
    found: boolean
    count: number
}

export function ServicesStatus() {
    const { data, isLoading, isError, error } = useApi<ServiceData[]>("http://localhost:5551/api/services")

    if (isLoading) {
        return (
            <Card>
                <CardHeader>
                    <CardTitle className="flex items-center gap-2">
                        <Settings className="h-5 w-5" />
                        Services Status
                    </CardTitle>
                </CardHeader>
                <CardContent>
                    <div className="animate-pulse space-y-3">
                        {[1, 2, 3].map((i) => (
                            <div key={i} className="flex items-center justify-between">
                                <div className="h-4 bg-muted rounded w-24"></div>
                                <div className="h-6 bg-muted rounded w-16"></div>
                            </div>
                        ))}
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
                        <Settings className="h-5 w-5" />
                        Services Status
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
                    <Settings className="h-5 w-5" />
                    Services Status
                </CardTitle>
            </CardHeader>
            <CardContent>
                <div className="space-y-4">
                    {data?.map((service) => (
                        <div key={service.name} className="flex items-center justify-between p-3 rounded-lg border">
                            <div className="flex items-center gap-3">
                                {service.found ? (
                                    <CheckCircle className="h-5 w-5 text-green-500" />
                                ) : (
                                    <XCircle className="h-5 w-5 text-destructive" />
                                )}
                                <span className="font-medium capitalize">{service.name}</span>
                            </div>
                            <div className="flex items-center gap-2">
                                <Badge variant={service.found ? "default" : "destructive"}>
                                    {service.found ? "Running" : "Stopped"}
                                </Badge>
                                {service.found && (
                                    <Badge variant="outline">
                                        {service.count} instance{service.count !== 1 ? "s" : ""}
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
