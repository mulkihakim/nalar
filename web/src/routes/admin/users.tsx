import { ProtectedRoute } from '../../features/auth/components/ProtectedRoute'
import { Card, CardContent, CardHeader, CardTitle } from '../../components/ui/card'
import { useAuth } from '../../features/auth/context'

export function AdminUsersPage() {
  const { user } = useAuth()

  return (
    <ProtectedRoute allowedRoles={['admin', 'asesor']}>
      <div className="space-y-6">
        <div>
          <h1 className="text-3xl font-bold tracking-tight">Manajemen Pengguna</h1>
          <p className="text-slate-500">
            Kelola data staf dan siswa (Login sebagai: {user?.name}, role: {user?.role})
          </p>
        </div>

        <Card>
          <CardHeader>
            <CardTitle>Daftar Pengguna</CardTitle>
          </CardHeader>
          <CardContent>
            <p className="text-sm text-slate-600">
              Modul manajemen pengguna (F-02) akan dimuat di sini.
            </p>
          </CardContent>
        </Card>
      </div>
    </ProtectedRoute>
  )
}
