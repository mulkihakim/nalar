import { useState } from 'react'
import { useNavigate } from '@tanstack/react-router'
import { ProtectedRoute } from '@/features/auth/components/ProtectedRoute'
import { useAuth } from '@/features/auth/context'
import { useMyExams, useStartSession, useMyExamSessions } from '@/features/session/api'
import { ModeSelectDialog } from '@/features/session/components/ModeSelectDialog'
import { Modal } from '@/components/Modal'
import { Button } from '@/components/ui/button'
import { Badge } from '@/components/ui/badge'
import type { Exam } from '@/types'
import type { SessionMode } from '@/features/session/types'
import { ModeConflictApiError } from '@/features/session/api'
import { BookOpen, Play, History, CheckCircle2, Clock } from 'lucide-react'

export function StudentExamsPage() {
  const { user } = useAuth()
  const navigate = useNavigate()
  const { data: exams, isLoading, error } = useMyExams()
  const startSessionMutation = useStartSession()

  const [selectedExam, setSelectedExam] = useState<Exam | null>(null)
  const [isModeDialogOpen, setIsModeDialogOpen] = useState(false)
  const [conflictMode, setConflictMode] = useState<SessionMode | undefined>()

  // History modal
  const [historyExam, setHistoryExam] = useState<Exam | null>(null)
  const { data: historySessions, isLoading: isLoadingHistory } = useMyExamSessions(historyExam?.id || 0)

  const handleOpenStart = (exam: Exam) => {
    setSelectedExam(exam)
    setConflictMode(undefined)
    setIsModeDialogOpen(true)
  }

  const handleSelectMode = async (mode: SessionMode, confirmModeChange = false) => {
    if (!selectedExam) return

    try {
      const session = await startSessionMutation.mutateAsync({
        examId: selectedExam.id,
        payload: {
          mode,
          confirm_mode_change: confirmModeChange,
        },
      })
      setIsModeDialogOpen(false)
      navigate({
        to: `/student/exam/${selectedExam.id}` as any,
        search: { sessionId: session.id } as any,
      })
    } catch (err) {
      if (err instanceof ModeConflictApiError) {
        setConflictMode(err.currentMode)
      } else {
        alert(err instanceof Error ? err.message : 'Gagal memulai sesi')
      }
    }
  }

  return (
    <ProtectedRoute allowedRoles={['siswa']}>
      <div className="space-y-6 max-w-5xl mx-auto">
        <div>
          <h1 className="text-2xl font-bold tracking-tight text-slate-900">Daftar Latihan & Ujian</h1>
          <p className="text-sm text-slate-500">
            Selamat datang, <span className="font-semibold text-slate-700">{user?.name}</span>. Silakan pilih paket latihan di bawah ini untuk memulai latihan penalaran Toulmin.
          </p>
        </div>

        {isLoading ? (
          <div className="py-12 text-center text-slate-500">Memuat daftar ujian...</div>
        ) : error ? (
          <div className="p-4 bg-red-50 text-red-700 rounded-xl border border-red-200 text-sm">
            Gagal memuat ujian: {(error as Error).message}
          </div>
        ) : !exams || exams.length === 0 ? (
          <div className="py-16 text-center bg-white rounded-2xl border border-slate-200 p-8">
            <BookOpen className="w-12 h-12 text-slate-300 mx-auto mb-3" />
            <h3 className="font-bold text-slate-800 text-base">Belum Ada Ujian yang Ditugaskan</h3>
            <p className="text-xs text-slate-500 max-w-sm mx-auto mt-1">
              Anda belum terdaftar dalam kelas atau penugasan ujian aktif. Hubungi pengajar atau asesor Anda untuk mendapatkan akses.
            </p>
          </div>
        ) : (
          <div className="grid grid-cols-1 md:grid-cols-2 gap-5">
            {exams.map((exam) => (
              <div
                key={exam.id}
                className="bg-white p-6 rounded-2xl border border-slate-200 shadow-sm hover:border-slate-300 transition-all flex flex-col justify-between"
              >
                <div>
                  <div className="flex items-center justify-between gap-2 mb-2">
                    <span className="text-[11px] font-semibold text-indigo-600 bg-indigo-50 px-2.5 py-0.5 rounded-full border border-indigo-100">
                      Materi: {exam.material?.title || 'Bacaan Argumen'}
                    </span>
                    <Badge variant="outline" className="bg-emerald-50 text-emerald-700 border-emerald-200 text-xs">
                      Aktif
                    </Badge>
                  </div>

                  <h3 className="font-bold text-slate-900 text-lg leading-snug">{exam.title}</h3>
                  <p className="text-xs text-slate-500 mt-2 line-clamp-2">
                    {exam.material?.content || 'Baca teks materi dan susun komponen argumen Toulmin.'}
                  </p>
                </div>

                <div className="pt-6 mt-6 border-t border-slate-100 flex items-center justify-between">
                  <div className="text-xs text-slate-500">
                    <span className="font-semibold text-slate-700">{exam.arguments_per_session}</span> Argumen per sesi
                  </div>

                  <div className="flex items-center gap-2">
                    <Button
                      variant="outline"
                      size="sm"
                      onClick={() => setHistoryExam(exam)}
                      className="text-xs text-slate-600"
                    >
                      <History className="w-3.5 h-3.5 mr-1" />
                      Riwayat
                    </Button>
                    <Button
                      size="sm"
                      onClick={() => handleOpenStart(exam)}
                      className="bg-indigo-600 hover:bg-indigo-700 text-white text-xs font-semibold px-4"
                    >
                      <Play className="w-3.5 h-3.5 mr-1 fill-white" />
                      Mulai Latihan
                    </Button>
                  </div>
                </div>
              </div>
            ))}
          </div>
        )}

        {/* Modal Dialog Pilih Mode */}
        {selectedExam && (
          <ModeSelectDialog
            isOpen={isModeDialogOpen}
            onClose={() => setIsModeDialogOpen(false)}
            onSelectMode={handleSelectMode}
            isLoading={startSessionMutation.isPending}
            conflictMode={conflictMode}
            onCancelConflict={() => {
              setConflictMode(undefined)
              setIsModeDialogOpen(false)
            }}
          />
        )}

        {/* Modal Riwayat Percobaan Siswa */}
        <Modal
          isOpen={!!historyExam}
          onClose={() => setHistoryExam(null)}
          title={`Riwayat Latihan: ${historyExam?.title}`}
          maxWidth="2xl"
        >
          <div className="space-y-4 py-2">
            {isLoadingHistory ? (
              <div className="py-8 text-center text-slate-400">Memuat riwayat...</div>
            ) : !historySessions || historySessions.length === 0 ? (
              <div className="py-8 text-center text-slate-400 text-sm">
                Belum ada riwayat pengerjaan pada ujian ini.
              </div>
            ) : (
              <div className="space-y-3 max-h-[400px] overflow-y-auto pr-1">
                {historySessions.map((s) => (
                  <div
                    key={s.id}
                    className="p-4 rounded-xl border border-slate-200 bg-white flex items-center justify-between gap-4"
                  >
                    <div>
                      <div className="flex items-center gap-2">
                        <span className="font-bold text-sm text-slate-900">
                          Percobaan #{s.attempt_no}
                        </span>
                        <Badge variant="outline" className="capitalize text-[10px]">
                          Mode {s.mode}
                        </Badge>
                        {s.status === 'completed' ? (
                          <Badge variant="outline" className="bg-emerald-50 text-emerald-700 border-emerald-200 text-[10px]">
                            <CheckCircle2 className="w-3 h-3 mr-0.5" /> Selesai
                          </Badge>
                        ) : (
                          <Badge variant="outline" className="bg-amber-50 text-amber-700 border-amber-200 text-[10px]">
                            <Clock className="w-3 h-3 mr-0.5" /> Sedang Berjalan
                          </Badge>
                        )}
                      </div>
                      <p className="text-xs text-slate-500 mt-1">
                        Selesai {s.completed_arguments} dari {s.total_arguments} argumen · Total {s.total_attempts} percobaan
                      </p>
                    </div>

                    {s.status === 'in_progress' && (
                      <Button
                        size="sm"
                        onClick={() => {
                          setHistoryExam(null)
                          navigate({
                            to: `/student/exam/${historyExam?.id}` as any,
                            search: { sessionId: s.id } as any,
                          })
                        }}
                        className="bg-indigo-600 hover:bg-indigo-700 text-white text-xs"
                      >
                        Lanjutkan
                      </Button>
                    )}
                  </div>
                ))}
              </div>
            )}

            <div className="flex justify-end pt-3 border-t border-slate-100">
              <Button variant="outline" onClick={() => setHistoryExam(null)}>
                Tutup
              </Button>
            </div>
          </div>
        </Modal>
      </div>
    </ProtectedRoute>
  )
}
