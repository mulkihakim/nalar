import { type ReactNode } from 'react'
import { Navigate } from '@tanstack/react-router'
import { useAuth } from '../context'
import type { Role } from '../../../types'

interface ProtectedRouteProps {
  children?: ReactNode
  allowedRoles?: Role[]
}

export function ProtectedRoute({ children, allowedRoles }: ProtectedRouteProps) {
  const { user, isAuthenticated, isLoading, logout } = useAuth()

  if (isLoading) {
    return (
      <div className="flex items-center justify-center min-h-[60vh]">
        <div className="text-muted-foreground animate-pulse">Memuat data pengguna...</div>
      </div>
    )
  }

  if (!isAuthenticated || !user) {
    return <Navigate to="/login" />
  }

  if (allowedRoles && allowedRoles.length > 0 && !allowedRoles.includes(user.role)) {
    return (
      <div className="flex flex-col items-center justify-center min-h-[60vh] space-y-4">
        <h2 className="text-2xl font-bold text-red-600">403 — Akses Ditolak</h2>
        <p className="text-muted-foreground">
          Akun Anda ({user.username}, role: {user.role}) tidak memiliki izin untuk mengakses halaman ini.
        </p>
        <button
          type="button"
          onClick={logout}
          className="text-sm underline text-blue-600 hover:text-blue-800 cursor-pointer"
        >
          Keluar dan masuk dengan akun lain
        </button>
      </div>
    )
  }

  return <>{children}</>
}
