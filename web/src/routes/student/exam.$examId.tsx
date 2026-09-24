import { useState, useEffect } from 'react'
import { useParams, useSearch, useNavigate } from '@tanstack/react-router'
import { ProtectedRoute } from '@/features/auth/components/ProtectedRoute'
import {
  useSessionDetails,
  useRecordDrop,
  useConfirmArgument,
  useMonitoringAnalytics,
  useAnalysisAnalytics,
  useStartSession,
} from '@/features/session/api'
import { ArgumentBoard } from '@/features/session/components/ArgumentBoard'
import { MonitoringView } from '@/features/session/components/MonitoringView'
import { SessionSummary } from '@/features/session/components/SessionSummary'
import { Modal } from '@/components/Modal'
import { ConfirmDialog } from '@/components/ConfirmDialog'
import { Button } from '@/components/ui/button'
import { Badge } from '@/components/ui/badge'
import type { OptionItem } from '@/features/session/types'
import {
  ArrowLeft,
  BookOpen,
  Sparkles,
  Users,
  CheckCircle2,
  Lock,
  ChevronLeft,
  ChevronRight,
} from 'lucide-react'

export function StudentExamPage() {
  const { examId } = useParams({ strict: false }) as { examId: string }
  const search = useSearch({ strict: false }) as { sessionId?: number }
  const navigate = useNavigate()

  const parsedExamId = parseInt(examId, 10)
  const sessionId = search.sessionId ? Number(search.sessionId) : 0

  const { data: session, isLoading, refetch } = useSessionDetails(sessionId)
  const recordDropMutation = useRecordDrop()
  const confirmArgumentMutation = useConfirmArgument()
  const startSessionMutation = useStartSession()

  // State untuk navigasi antar soal (melihat soal yang sudah dijawab atau aktif)
  const [viewingIndex, setViewingIndex] = useState<number>(0)
  const [hasInitializedIndex, setHasInitializedIndex] = useState(false)

  // Dialog konfirmasi sebelum keluar ke daftar ujian
  const [isExitDialogOpen, setIsExitDialogOpen] = useState(false)

  // Material text modal
  const [isMaterialOpen, setIsMaterialOpen] = useState(false)

  // Monitoring modal for social mode
  const [isMonitoringOpen, setIsMonitoringOpen] = useState(false)
  const [isMonitoringSuccess, setIsMonitoringSuccess] = useState(false)
  const [activeArgForMonitoring, setActiveArgForMonitoring] = useState<number | undefined>()

  const { data: monitoringData, isLoading: isLoadingMonitoring } = useMonitoringAnalytics(
    sessionId,
    activeArgForMonitoring
  )

  // Analysis data for completed session
  const isSessionCompleted = session?.status === 'completed'
  const isAllArgumentsCompleted = session?.arguments.every((a) => a.completed) || false

  // State untuk berpindah antara tampilan Soal (Review) dan Halaman Hasil (Summary)
  // null = ikuti status sesi (jika sesi dari awal sudah selesai saat dibuka, tampilkan summary)
  const [isViewingSummary, setIsViewingSummary] = useState<boolean | null>(null)
  const shouldShowSummary = isViewingSummary ?? isSessionCompleted

  const { data: analysisData, isLoading: isLoadingAnalysis } = useAnalysisAnalytics(
    sessionId,
    isSessionCompleted && session?.mode === 'social'
  )

  // Inisialisasi viewingIndex ke argumen pertama yang belum selesai
  useEffect(() => {
    if (session && !hasInitializedIndex) {
      const firstUncompleted = session.arguments.findIndex((a) => !a.completed)
      if (firstUncompleted !== -1) {
        setViewingIndex(firstUncompleted)
      } else {
        setViewingIndex(0)
      }
      setHasInitializedIndex(true)
    }
  }, [session, hasInitializedIndex])

  if (isLoading) {
    return (
      <ProtectedRoute allowedRoles={['siswa']}>
        <div className="py-20 text-center text-slate-500">Memuat sesi ujian...</div>
      </ProtectedRoute>
    )
  }

  if (!session) {
    return (
      <ProtectedRoute allowedRoles={['siswa']}>
        <div className="max-w-md mx-auto py-20 text-center">
          <p className="text-slate-600 mb-4">Sesi latihan tidak ditemukan atau sudah selesai.</p>
          <Button onClick={() => navigate({ to: '/student/exams' as any })}>
            Kembali ke Daftar Ujian
          </Button>
        </div>
      </ProtectedRoute>
    )
  }

  // Cari argumen aktif pertama yang belum selesai
  const activeUncompletedIndex = session.arguments.findIndex((a) => !a.completed)
  const currentArg = session.arguments[viewingIndex] || session.arguments[0]
  const isViewingCompletedArg = currentArg ? currentArg.completed : false
  const isLastArgument = viewingIndex === session.arguments.length - 1

  const handleDropOption = (option: OptionItem, slot: 'ground' | 'warrant') => {
    if (!currentArg || isViewingCompletedArg) return
    recordDropMutation.mutate({
      sessionId: session.id,
      argumentId: currentArg.id,
      payload: {
        option_id: option.id,
        slot,
      },
    })
  }

  const handleConfirm = async (groundOptionId: number, warrantOptionId: number) => {
    if (!currentArg || isViewingCompletedArg) return null
    try {
      const res = await confirmArgumentMutation.mutateAsync({
        sessionId: session.id,
        argumentId: currentArg.id,
        payload: {
          ground_option_id: groundOptionId,
          warrant_option_id: warrantOptionId,
        },
      })
      if (res.completed) {
        setActiveArgForMonitoring(currentArg.id)
      }
      if (res.is_session_completed) {
        // Jangan langsung mengganti layar ke summary di belakang modal jika sesi selesai
        setIsViewingSummary(false)
      }
      await refetch()
      return res
    } catch (err) {
      alert(err instanceof Error ? err.message : 'Gagal mengonfirmasi jawaban')
      return null
    }
  }

  const handleNextArgument = async () => {
    const updated = await refetch()
    if (updated.data) {
      const allDone =
        updated.data.arguments.every((a) => a.completed) || updated.data.status === 'completed'
      if (allDone) {
        setIsViewingSummary(true)
        return
      }
      // Pindahkan ke soal belum selesai berikutnya
      const nextUncompleted = updated.data.arguments.findIndex((a) => !a.completed)
      if (nextUncompleted !== -1) {
        setViewingIndex(nextUncompleted)
      } else {
        // Semua sudah selesai
        setViewingIndex(0)
      }
    }
  }

  const handleRestart = async () => {
    try {
      const newSession = await startSessionMutation.mutateAsync({
        examId: parsedExamId,
        payload: {
          mode: session.mode,
        },
      })
      setHasInitializedIndex(false)
      navigate({
        to: `/student/exam/${parsedExamId}` as any,
        search: { sessionId: newSession.id } as any,
      })
    } catch (err) {
      alert(err instanceof Error ? err.message : 'Gagal memulai latihan baru')
    }
  }

  const handleBackToExamsClick = () => {
    if (isSessionCompleted) {
      navigate({ to: '/student/exams' as any })
    } else {
      setIsExitDialogOpen(true)
    }
  }

  return (
    <ProtectedRoute allowedRoles={['siswa']}>
      <div className="max-w-5xl mx-auto space-y-4 pb-12">
        {/* Top Bar Navigation */}
        <div className="flex flex-col sm:flex-row sm:items-center justify-between gap-3 bg-white p-4 sm:p-5 rounded-2xl border border-slate-200 shadow-xs">
          <div className="flex items-center gap-3">
            <Button
              variant="outline"
              size="sm"
              onClick={handleBackToExamsClick}
              className="text-slate-600 hover:text-slate-900"
            >
              <ArrowLeft className="w-4 h-4 mr-1.5" />
              Daftar Ujian
            </Button>

            <div>
              <div className="flex items-center gap-2">
                <h1 className="font-bold text-slate-900 text-base md:text-lg">
                  {session.exam_title}
                </h1>
                <Badge variant="outline" className="capitalize text-[11px] bg-slate-50">
                  {session.mode === 'help' ? (
                    <span className="flex items-center gap-1 text-amber-700">
                      <Sparkles className="w-3 h-3 text-amber-500" /> Mode Bantuan
                    </span>
                  ) : session.mode === 'social' ? (
                    <span className="flex items-center gap-1 text-indigo-700">
                      <Users className="w-3 h-3 text-indigo-500" /> Mode Sosial
                    </span>
                  ) : (
                    'Mode Standar'
                  )}
                </Badge>
              </div>
              <p className="text-xs text-slate-500 mt-0.5">Percobaan #{session.attempt_no}</p>
            </div>
          </div>

          <div className="flex items-center gap-2 self-end sm:self-center">
            <Button
              variant="outline"
              size="sm"
              onClick={() => setIsMaterialOpen(true)}
              className="text-xs text-indigo-700 bg-indigo-50 border-indigo-200 hover:bg-indigo-100"
            >
              <BookOpen className="w-3.5 h-3.5 mr-1.5" />
              Baca Materi
            </Button>
          </div>
        </div>

        {/* Stepper / Navigasi Argumen (Sebelumnya / Sedang Dikerjakan / Setelahnya) */}
        {!shouldShowSummary && (
          <div className="bg-white p-3 rounded-2xl border border-slate-200 shadow-xs flex flex-wrap items-center justify-between gap-3">
            <div className="flex items-center gap-1.5 overflow-x-auto py-0.5">
              <span className="text-xs font-semibold text-slate-400 mr-1 hidden sm:inline">
                Pilih Soal:
              </span>
              {session.arguments.map((arg, idx) => {
                const isCompleted = arg.completed
                const isCurrentActive = idx === activeUncompletedIndex
                const isSelected = idx === viewingIndex
                const isLocked =
                  !isCompleted &&
                  !isCurrentActive &&
                  (activeUncompletedIndex === -1 || idx > activeUncompletedIndex) &&
                  !isAllArgumentsCompleted &&
                  session.status !== 'completed'

                let btnStyle = 'border-slate-200 text-slate-700 hover:bg-slate-50'
                if (isSelected) {
                  btnStyle = 'border-indigo-600 bg-indigo-600 text-white shadow-xs font-bold'
                } else if (isCompleted) {
                  btnStyle = 'border-emerald-200 bg-emerald-50 text-emerald-700 hover:bg-emerald-100 font-medium'
                } else if (isLocked) {
                  btnStyle = 'border-slate-100 bg-slate-50 text-slate-400 cursor-not-allowed opacity-60'
                }

                return (
                  <button
                    key={arg.id}
                    type="button"
                    disabled={isLocked}
                    onClick={() => setViewingIndex(idx)}
                    className={`flex items-center gap-1 px-3 py-1.5 rounded-xl border text-xs transition-all ${btnStyle}`}
                    title={isLocked ? 'Selesaikan soal sebelumnya terlebih dahulu' : `Buka Soal ${idx + 1}`}
                  >
                    {isCompleted ? (
                      <CheckCircle2 className={`w-3.5 h-3.5 ${isSelected ? 'text-white' : 'text-emerald-600'}`} />
                    ) : isLocked ? (
                      <Lock className="w-3 h-3 text-slate-400" />
                    ) : (
                      <span className={`w-2 h-2 rounded-full ${isSelected ? 'bg-white' : 'bg-indigo-600'}`} />
                    )}
                    <span>Argumen {idx + 1}</span>
                    {isCompleted && <span className="text-[10px] opacity-80">(Selesai)</span>}
                  </button>
                )
              })}
            </div>

            {/* Tombol Navigasi Sebelumnya / Berikutnya / Selesaikan Ujian */}
            <div className="flex items-center gap-1.5 text-xs ml-auto">
              <Button
                variant="outline"
                size="sm"
                disabled={viewingIndex === 0}
                onClick={() => setViewingIndex((prev) => Math.max(0, prev - 1))}
                className="h-8 px-2.5 text-xs text-slate-600"
              >
                <ChevronLeft className="w-3.5 h-3.5 mr-0.5" /> Sebelumnya
              </Button>
              <Button
                variant="outline"
                size="sm"
                disabled={
                  viewingIndex >= session.arguments.length - 1 ||
                  (!session.arguments[viewingIndex].completed && viewingIndex === activeUncompletedIndex)
                }
                onClick={() =>
                  setViewingIndex((prev) => Math.min(session.arguments.length - 1, prev + 1))
                }
                className="h-8 px-2.5 text-xs text-slate-600"
              >
                Berikutnya <ChevronRight className="w-3.5 h-3.5 ml-0.5" />
              </Button>

              {(isAllArgumentsCompleted || session.status === 'completed') && (
                <Button
                  size="sm"
                  onClick={() => setIsViewingSummary(true)}
                  className="h-8 px-3 text-xs bg-emerald-600 hover:bg-emerald-700 text-white font-semibold ml-1 shadow-xs"
                >
                  Lihat Hasil
                </Button>
              )}
            </div>
          </div>
        )}

        {/* Jika Sesi Sudah Selesai dan Sedang Melihat Ringkasan Hasil */}
        {shouldShowSummary ? (
          <SessionSummary
            session={session}
            analysisData={analysisData}
            isLoadingAnalysis={isLoadingAnalysis}
            onRestart={handleRestart}
            onBackToList={() => navigate({ to: '/student/exams' as any })}
            onReviewQuestions={() => setIsViewingSummary(false)}
          />
        ) : (
          /* Layar Argumen Toulmin Aktif / Review */
          <ArgumentBoard
            key={currentArg.id}
            argument={currentArg}
            mode={session.mode}
            onDropOption={handleDropOption}
            onConfirm={handleConfirm}
            onNextArgument={handleNextArgument}
            isLastArgument={isLastArgument || isAllArgumentsCompleted || session.status === 'completed'}
            isSubmitting={confirmArgumentMutation.isPending}
            showMonitoring={session.mode === 'social'}
            onShowMonitoring={(isSuccess = false) => {
              setActiveArgForMonitoring(currentArg.id)
              setIsMonitoringSuccess(isSuccess)
              setIsMonitoringOpen(true)
            }}
            isReadOnly={isViewingCompletedArg || session.status === 'completed' || isAllArgumentsCompleted}
            isAllCompleted={isAllArgumentsCompleted || session.status === 'completed'}
            onFinishExam={() => setIsViewingSummary(true)}
          />
        )}

        {/* Modal Baca Materi */}
        <Modal
          isOpen={isMaterialOpen}
          onClose={() => setIsMaterialOpen(false)}
          title={`Materi: ${session.material_title || session.exam_title}`}
          maxWidth="2xl"
        >
          <div className="space-y-4 py-2">
            <div className="prose prose-slate max-h-[450px] overflow-y-auto text-sm text-slate-800 leading-relaxed pr-2 whitespace-pre-line">
              {session.material_content || 'Tidak ada teks materi.'}
            </div>
            <div className="flex justify-end pt-3 border-t border-slate-100">
              <Button onClick={() => setIsMaterialOpen(false)} className="bg-indigo-600 text-white">
                Selesai Membaca
              </Button>
            </div>
          </div>
        </Modal>

        {/* Modal Monitoring Pola Teman (Mode Sosial) */}
        <MonitoringView
          isOpen={isMonitoringOpen}
          onClose={() => {
            setIsMonitoringOpen(false)
            setIsMonitoringSuccess(false)
            if (isAllArgumentsCompleted || session.status === 'completed') {
              setIsViewingSummary(false)
            }
          }}
          data={monitoringData}
          claimText={currentArg?.claim_text}
          isLoading={isLoadingMonitoring}
          isSuccessDialog={isMonitoringSuccess}
          onNextArgument={() => {
            setIsMonitoringOpen(false)
            setIsMonitoringSuccess(false)
            if (isLastArgument || isAllArgumentsCompleted || session.status === 'completed') {
              setIsViewingSummary(true)
            } else {
              handleNextArgument()
            }
          }}
          isLastArgument={isLastArgument || isAllArgumentsCompleted || session.status === 'completed'}
        />

        {/* Dialog Konfirmasi Keluar ke Daftar Ujian */}
        <ConfirmDialog
          isOpen={isExitDialogOpen}
          onClose={() => setIsExitDialogOpen(false)}
          onConfirm={() => {
            setIsExitDialogOpen(false)
            navigate({ to: '/student/exams' as any })
          }}
          title="Kembali ke Daftar Ujian?"
          description="Sesi latihan Anda akan tetap berstatus berjalan dan progres Anda tidak akan hilang. Anda dapat melanjutkannya kembali kapan saja dari halaman daftar ujian."
          confirmText="Ya, Keluar"
          cancelText="Lanjutkan Latihan"
        />
      </div>
    </ProtectedRoute>
  )
}
