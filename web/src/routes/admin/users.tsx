import { ProtectedRoute } from '../../features/auth/components/ProtectedRoute'
import { useAuth } from '../../features/auth/context'
import { useUsers } from '../../features/user/api'
import { UserTable } from '../../features/user/components/UserTable'

export function AdminUsersPage() {
  const { user } = useAuth()
  const { data: users = [], isLoading } = useUsers()

  return (
    <ProtectedRoute allowedRoles={['admin', 'asesor']}>
      <div className="space-y-6">
        <div>
          <h1 className="text-2xl sm:text-3xl font-bold tracking-tight text-slate-900">
            Manajemen Pengguna
          </h1>
          <p className="text-sm text-slate-500 mt-1">
            Kelola data staf dan siswa (Login sebagai: <span className="font-semibold text-slate-700">{user?.name}</span>, role: <span className="uppercase text-xs font-bold text-indigo-600">{user?.role}</span>)
          </p>
        </div>

        <UserTable
          users={users}
          isLoading={isLoading}
          currentUserRole={user?.role || 'siswa'}
          currentUserId={user?.id || 0}
        />
      </div>
    </ProtectedRoute>
  )
}
