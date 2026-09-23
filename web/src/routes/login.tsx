import { Navigate } from '@tanstack/react-router'
import { LoginForm } from '../features/auth/components/LoginForm'
import { useAuth } from '../features/auth/context'

export function LoginPage() {
  const { user, isAuthenticated, isLoading } = useAuth()

  if (isLoading) {
    return (
      <div className="flex items-center justify-center min-h-[50vh]">
        <div className="text-slate-500 animate-pulse">Memuat...</div>
      </div>
    )
  }

  if (isAuthenticated && user) {
    if (user.role === 'admin') {
      return <Navigate to="/admin/users" />
    }
    if (user.role === 'asesor') {
      return <Navigate to="/admin/classes" />
    }
    return <Navigate to="/student/exams" />
  }

  return <LoginForm />
}
