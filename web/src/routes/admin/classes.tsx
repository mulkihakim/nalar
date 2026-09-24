import { ProtectedRoute } from '../../features/auth/components/ProtectedRoute'
import { useAuth } from '../../features/auth/context'
import { useClasses } from '../../features/class/api'
import { ClassTable } from '../../features/class/components/ClassTable'

export function AdminClassesPage() {
  const { user } = useAuth()
  const { data: classes = [], isLoading } = useClasses()

  return (
    <ProtectedRoute allowedRoles={['admin', 'asesor']}>
      <div className="space-y-6">
        <div>
          <h1 className="text-2xl sm:text-3xl font-bold tracking-tight text-slate-900">
            Manajemen Kelas
          </h1>
          <p className="text-sm text-slate-500 mt-1">
            Kelola data kelas dan keanggotaan siswa (Login sebagai: <span className="font-semibold text-slate-700">{user?.name}</span>, role: <span className="uppercase text-xs font-bold text-indigo-600">{user?.role}</span>)
          </p>
        </div>

        <ClassTable
          classes={classes}
          isLoading={isLoading}
          currentUserRole={user?.role || ''}
          currentUserId={user?.id || 0}
        />
      </div>
    </ProtectedRoute>
  )
}
