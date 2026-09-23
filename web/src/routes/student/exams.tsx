import { ProtectedRoute } from '../../features/auth/components/ProtectedRoute'
import { Card, CardContent, CardHeader, CardTitle } from '../../components/ui/card'
import { useAuth } from '../../features/auth/context'

export function StudentExamsPage() {
  const { user } = useAuth()

  return (
    <ProtectedRoute allowedRoles={['siswa']}>
      <div className="space-y-6">
        <div>
          <h1 className="text-3xl font-bold tracking-tight">Daftar Ujian & Latihan</h1>
          <p className="text-slate-500">
            Pilih paket latihan untuk memulai sesi (Siswa: {user?.name})
          </p>
        </div>

        <Card>
          <CardHeader>
            <CardTitle>Ujian Tersedia</CardTitle>
          </CardHeader>
          <CardContent>
            <p className="text-sm text-slate-600">
              Daftar ujian siswa (F-06/F-14) akan dimuat di sini.
            </p>
          </CardContent>
        </Card>
      </div>
    </ProtectedRoute>
  )
}
