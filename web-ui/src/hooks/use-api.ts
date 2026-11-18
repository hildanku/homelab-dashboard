
import { useQuery, useQueryClient } from "@tanstack/react-query"
import { useNavigate } from "@tanstack/react-router"

import { authKeys, UNAUTHENTICATED } from "@/lib/auth"

interface ApiResponse<T> {
    message: string
    result: T
}

export function useApi<T>(endpoint: string, interval = 5000) {
    const navigate = useNavigate()
    const queryClient = useQueryClient()

    return useQuery<T>({
        queryKey: [endpoint],
        queryFn: async () => {
            const res = await fetch(endpoint, { credentials: "include" })
            if (res.status === 401) {
                queryClient.setQueryData(authKeys.session, UNAUTHENTICATED)
                void navigate({ to: "/login", search: { redirect: "/management" }, replace: true })
            }
            if (!res.ok) {
                throw new Error(`HTTP error! status: ${res.status}`)
            }
            const data: ApiResponse<T> = await res.json()
            return data.result
        },
        refetchInterval: interval,
    })
}
