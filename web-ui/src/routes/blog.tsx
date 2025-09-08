import { createFileRoute } from '@tanstack/react-router'

export const Route = createFileRoute('/blog')({
    component: BlogPage,
})

export default function BlogPage() {
    return (
        <div className="p-4">
            <strong>data from API</strong>
            <h1 className="text-xl font-bold">Blog</h1>
        </div>
    )
}
