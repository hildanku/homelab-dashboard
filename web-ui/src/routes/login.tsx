import { useState, type FormEvent } from 'react'
import { createFileRoute } from '@tanstack/react-router'

import { Button } from '@/components/ui/button'
import { Card, CardContent, CardDescription, CardHeader, CardTitle } from '@/components/ui/card'
import { Input } from '@/components/ui/input'
import { loginRequest } from '@/lib/auth'

export const Route = createFileRoute('/login')({
    validateSearch: (search: Record<string, unknown>) => {
        const redirect = search.redirect
        if (typeof redirect === 'string' && redirect.startsWith('/') && !redirect.startsWith('//')) {
            return { redirect }
        }
        return { redirect: '/management' }
    },
    component: LoginPage,
})

function LoginPage() {
    const { redirect } = Route.useSearch()
    const [username, setUsername] = useState('')
    const [password, setPassword] = useState('')
    const [error, setError] = useState('')
    const [submitting, setSubmitting] = useState(false)

    async function handleSubmit(event: FormEvent<HTMLFormElement>) {
        event.preventDefault()
        setSubmitting(true)
        setError('')

        try {
            await loginRequest(username, password)
            window.location.replace(redirect)
        } catch {
            setError('Username atau password salah')
            setSubmitting(false)
        }
    }

    return (
        <div className="flex min-h-svh items-center justify-center bg-background px-4">
            <Card className="w-full max-w-sm">
                <CardHeader>
                    <CardTitle>Homelab Dashboard</CardTitle>
                    <CardDescription>Masuk untuk melihat dashboard homelab.</CardDescription>
                </CardHeader>
                <CardContent>
                    <form className="flex flex-col gap-4" onSubmit={handleSubmit}>
                        <div className="grid gap-2">
                            <label htmlFor="username" className="text-sm font-medium">Username</label>
                            <Input
                                id="username"
                                autoComplete="username"
                                value={username}
                                onChange={(event) => setUsername(event.target.value)}
                                required
                            />
                        </div>
                        <div className="grid gap-2">
                            <label htmlFor="password" className="text-sm font-medium">Password</label>
                            <Input
                                id="password"
                                type="password"
                                autoComplete="current-password"
                                value={password}
                                onChange={(event) => setPassword(event.target.value)}
                                required
                            />
                        </div>
                        {error && <p className="text-sm text-destructive">{error}</p>}
                        <Button type="submit" disabled={submitting}>
                            {submitting ? 'Memproses...' : 'Masuk'}
                        </Button>
                    </form>
                </CardContent>
            </Card>
        </div>
    )
}
