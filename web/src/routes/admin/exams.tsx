import { ProtectedRoute } from '../../features/auth/components/ProtectedRoute'
import { useAuth } from '../../features/auth/context'
import { useExams } from '../../features/exam/api'
import { ExamTable } from '../../features/exam/components/ExamTable'

export function AdminExamsPage() {
  const { user } = useAuth()
  const { data: exams = [], isLoading } = useExams()

  return (
    <ProtectedRoute allowedRoles={['admin', 'asesor']}>
      <div className="space-y-6">
        <div>
          <h1 className="text-2xl sm:text-3xl font-bold tracking-tight text-slate-900">
            Manajemen Ujian / Latihan
          </h1>
          <p className="text-sm text-slate-500 mt-1">
            Atur paket latihan, kuota argumen per sesi, dan izin akses kelas atau siswa (Login sebagai: <span className="font-semibold text-slate-700">{user?.name}</span>, role: <span className="uppercase text-xs font-bold text-indigo-600">{user?.role}</span>)
          </p>
        </div>

        <ExamTable
          exams={exams}
          isLoading={isLoading}
          currentUserRole={user?.role || ''}
          currentUserId={user?.id || 0}
        />
      </div>
    </ProtectedRoute>
  )
}
