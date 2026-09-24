import { ProtectedRoute } from '../../features/auth/components/ProtectedRoute'
import { useAuth } from '../../features/auth/context'
import { useMaterials } from '../../features/material/api'
import { MaterialTable } from '../../features/material/components/MaterialTable'

export function AdminMaterialsPage() {
  const { user } = useAuth()
  const { data: materials = [], isLoading } = useMaterials()

  return (
    <ProtectedRoute allowedRoles={['admin', 'asesor']}>
      <div className="space-y-6">
        <div>
          <h1 className="text-2xl sm:text-3xl font-bold tracking-tight text-slate-900">
            Manajemen Materi & Argumen
          </h1>
          <p className="text-sm text-slate-500 mt-1">
            Kelola bacaan topik dan susun butir argumen model Toulmin (Login sebagai: <span className="font-semibold text-slate-700">{user?.name}</span>, role: <span className="uppercase text-xs font-bold text-indigo-600">{user?.role}</span>)
          </p>
        </div>

        <MaterialTable
          materials={materials}
          isLoading={isLoading}
          currentUserRole={user?.role || ''}
          currentUserId={user?.id || 0}
        />
      </div>
    </ProtectedRoute>
  )
}
