
import { useQuery } from "@tanstack/react-query"

interface ApiResponse<T> {
    message: string
    result: T
}

export function useApi<T>(endpoint: string, interval = 5000) {
    return useQuery<T>({
        queryKey: [endpoint],
        queryFn: async () => {
            const res = await fetch(endpoint)
            if (!res.ok) {
                throw new Error(`HTTP error! status: ${res.status}`)
            }
            const data: ApiResponse<T> = await res.json()
            return data.result
        },
        refetchInterval: interval,
    })
}
