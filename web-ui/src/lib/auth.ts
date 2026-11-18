export interface Session {
  authenticated: boolean
  username?: string
}

interface ApiResponse<T> {
  message: string
  result: T
}

export const UNAUTHENTICATED: Session = { authenticated: false }

export const authKeys = {
  session: ["auth", "session"] as const,
}

async function parseResponse<T>(res: Response): Promise<T> {
  const data = (await res.json().catch(() => null)) as ApiResponse<T> | null
  if (!res.ok) {
    throw new Error(data?.message ?? "Request failed")
  }
  if (!data) {
    throw new Error("Invalid response")
  }
  return data.result
}

export async function getSession(): Promise<Session> {
  const res = await fetch("/api/auth/me", { credentials: "include" })
  if (res.status === 401) {
    return UNAUTHENTICATED
  }

  const result = await parseResponse<{ username?: string }>(res)
  return { authenticated: true, username: result.username }
}

export async function loginRequest(username: string, password: string): Promise<Session> {
  const res = await fetch("/api/login", {
    method: "POST",
    credentials: "include",
    headers: { "Content-Type": "application/json" },
    body: JSON.stringify({ username, password }),
  })

  const result = await parseResponse<{ username?: string }>(res)
  return { authenticated: true, username: result.username }
}

export async function logoutRequest(): Promise<void> {
  const res = await fetch("/api/logout", {
    method: "POST",
    credentials: "include",
  })
  await parseResponse<null>(res)
}
