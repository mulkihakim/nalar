import { useState } from 'react'
import { ProtectedRoute } from '@/features/auth/components/ProtectedRoute'
import { useExams } from '@/features/exam/api'
import { useStaffAnalysis, useExamResults } from '@/features/session/api'
import { AnalysisView } from '@/features/session/components/AnalysisView'
import { ResultsTable } from '@/features/session/components/ResultsTable'
import { Modal } from '@/components/Modal'
import { Button } from '@/components/ui/button'
import { Badge } from '@/components/ui/badge'
import { Card, CardContent } from '@/components/ui/card'
import type { Exam } from '@/types'
import type { AnalysisOptionStat } from '@/features/session/types'
import {
  ArrowLeft,
  BarChart3,
  ListChecks,
  BookOpen,
  Users,
  CheckCircle2,
  XCircle,
} from 'lucide-react'

export function AdminResultsPage() {
  const { data: exams, isLoading: isLoadingExams } = useExams()
  const [selectedExam, setSelectedExam] = useState<Exam | null>(null)
  const [activeTab, setActiveTab] = useState<'analysis' | 'results'>('analysis')
  const [selectedOptionForModal, setSelectedOptionForModal] = useState<AnalysisOptionStat | null>(
    null
  )

  const examIdNum = selectedExam ? selectedExam.id : 0

  const { data: staffAnalysis, isLoading: isLoadingAnalysis } = useStaffAnalysis(examIdNum)
  const { data: results, isLoading: isLoadingResults } = useExamResults(examIdNum)

  return (
    <ProtectedRoute allowedRoles={['admin', 'asesor']}>
      <div className="space-y-6 max-w-6xl mx-auto pb-12">
        {!selectedExam ? (
          /* View 1: List Paket Ujian */
          <div className="space-y-5">
            <div>
              <h1 className="text-2xl font-bold tracking-tight text-slate-900">Hasil Ujian</h1>
              <p className="text-sm text-slate-500 mt-1">
                Pilih paket latihan/ujian di bawah ini untuk melihat hasil analitik kelompok dan distribusi pilihan siswa.
              </p>
            </div>

            {isLoadingExams ? (
              <div className="py-16 text-center text-slate-400 text-sm">Memuat daftar ujian...</div>
            ) : !exams || exams.length === 0 ? (
              <Card>
                <CardContent className="py-16 text-center text-slate-500 text-sm">
                  Belum ada paket ujian yang tersedia.
                </CardContent>
              </Card>
            ) : (
              <div className="overflow-x-auto rounded-2xl border border-slate-200 bg-white shadow-xs">
                <table className="w-full text-left border-collapse text-sm">
                  <thead>
                    <tr className="bg-slate-50/80 border-b border-slate-200 text-xs font-semibold text-slate-600">
                      <th className="py-4 px-5">Judul Paket Ujian</th>
                      <th className="py-4 px-5">Materi Terkait</th>
                      <th className="py-4 px-5 text-center">Argumen / Sesi</th>
                      <th className="py-4 px-5 text-center">Status</th>
                      <th className="py-4 px-5 text-right">Aksi</th>
                    </tr>
                  </thead>
                  <tbody className="divide-y divide-slate-100">
                    {exams.map((exam) => (
                      <tr
                        key={exam.id}
                        className="hover:bg-indigo-50/30 transition-colors group cursor-pointer"
                        onClick={() => {
                          setSelectedExam(exam)
                          setActiveTab('analysis')
                        }}
                      >
                        <td className="py-4 px-5 font-semibold text-slate-900">
                          <div className="flex items-center gap-2">
                            <span className="group-hover:text-indigo-600 transition-colors">
                              {exam.title}
                            </span>
                          </div>
                        </td>
                        <td className="py-4 px-5 text-slate-600">
                          <div className="flex items-center gap-1.5 text-xs">
                            <BookOpen className="w-3.5 h-3.5 text-slate-400" />
                            {exam.material?.title || 'Materi Umum'}
                          </div>
                        </td>
                        <td className="py-4 px-5 text-center font-medium text-slate-700">
                          {exam.arguments_per_session}
                        </td>
                        <td className="py-4 px-5 text-center">
                          {exam.is_active ? (
                            <Badge
                              variant="outline"
                              className="bg-emerald-50 text-emerald-700 border-emerald-200 text-[11px]"
                            >
                              <CheckCircle2 className="w-3 h-3 mr-1" /> Aktif
                            </Badge>
                          ) : (
                            <Badge
                              variant="outline"
                              className="bg-slate-50 text-slate-500 border-slate-200 text-[11px]"
                            >
                              <XCircle className="w-3 h-3 mr-1" /> Nonaktif
                            </Badge>
                          )}
                        </td>
                        <td className="py-4 px-5 text-right">
                          <Button
                            size="sm"
                            variant="default"
                            className="bg-indigo-600 hover:bg-indigo-700 text-white text-xs shadow-xs"
                            onClick={(e) => {
                              e.stopPropagation()
                              setSelectedExam(exam)
                              setActiveTab('analysis')
                            }}
                          >
                            <BarChart3 className="w-3.5 h-3.5 mr-1.5" />
                            Lihat Hasil Analisis
                          </Button>
                        </td>
                      </tr>
                    ))}
                  </tbody>
                </table>
              </div>
            )}
          </div>
        ) : (
          /* View 2: Detail Analitik & Hasil Ujian */
          <div className="space-y-6">
            {/* Header Detail */}
            <div className="flex flex-col sm:flex-row sm:items-center justify-between gap-4 bg-white p-5 rounded-2xl border border-slate-200 shadow-xs">
              <div className="flex items-start sm:items-center gap-3">
                <Button
                  variant="outline"
                  size="sm"
                  onClick={() => {
                    setSelectedExam(null)
                    setSelectedOptionForModal(null)
                  }}
                  className="text-slate-600 hover:text-slate-900 shrink-0"
                >
                  <ArrowLeft className="w-4 h-4 mr-1.5" />
                  Daftar Ujian
                </Button>
                <div>
                  <h1 className="text-lg md:text-xl font-bold text-slate-900">
                    {selectedExam.title}
                  </h1>
                  <p className="text-xs text-slate-500 flex items-center gap-2 mt-0.5">
                    <span>Materi: {selectedExam.material?.title || 'Umum'}</span>
                    <span>•</span>
                    <span>{selectedExam.arguments_per_session} Argumen per Sesi</span>
                  </p>
                </div>
              </div>

              {/* Tab Switcher: Analitik Kelompok vs Ringkasan Hasil Siswa */}
              <div className="flex items-center gap-1.5 bg-slate-100 p-1 rounded-xl self-start sm:self-center border border-slate-200/60">
                <button
                  type="button"
                  onClick={() => setActiveTab('analysis')}
                  className={`flex items-center gap-1.5 px-3 py-1.5 rounded-lg text-xs font-semibold transition-all ${
                    activeTab === 'analysis'
                      ? 'bg-white text-indigo-700 shadow-xs'
                      : 'text-slate-600 hover:text-slate-900'
                  }`}
                >
                  <BarChart3 className="w-3.5 h-3.5" />
                  Analitik Sosial Kelompok
                </button>
                <button
                  type="button"
                  onClick={() => setActiveTab('results')}
                  className={`flex items-center gap-1.5 px-3 py-1.5 rounded-lg text-xs font-semibold transition-all ${
                    activeTab === 'results'
                      ? 'bg-white text-indigo-700 shadow-xs'
                      : 'text-slate-600 hover:text-slate-900'
                  }`}
                >
                  <ListChecks className="w-3.5 h-3.5" />
                  Hasil Siswa
                </button>
              </div>
            </div>

            {/* Content Tab */}
            {activeTab === 'analysis' ? (
              <AnalysisView
                data={staffAnalysis}
                isLoading={isLoadingAnalysis}
                isStaffView={true}
                onSelectOption={(option) => setSelectedOptionForModal(option)}
              />
            ) : (
              <ResultsTable results={results || []} isLoading={isLoadingResults} />
            )}
          </div>
        )}

        {/* Modal Daftar Siswa Pemilih Opsi (Staff View) */}
        <Modal
          isOpen={selectedOptionForModal !== null}
          onClose={() => setSelectedOptionForModal(null)}
          title="Daftar Siswa Pemilih Opsi"
          maxWidth="2xl"
        >
          {selectedOptionForModal && (
            <div className="space-y-4 py-2">
              {/* Opsi Box */}
              <div className="p-3.5 bg-slate-50 rounded-xl border border-slate-200 space-y-2">
                <div className="flex items-center justify-between">
                  <Badge
                    variant="outline"
                    className={`uppercase text-[10px] font-bold ${
                      selectedOptionForModal.type === 'ground'
                        ? 'bg-blue-50 text-blue-700 border-blue-200'
                        : 'bg-purple-50 text-purple-700 border-purple-200'
                    }`}
                  >
                    {selectedOptionForModal.type === 'ground' ? 'Ground' : 'Warrant'}
                  </Badge>
                  <span className="text-xs font-bold text-slate-700">
                    Rasio Pilihan: {selectedOptionForModal.ratio}
                  </span>
                </div>
                <p className="text-xs font-medium text-slate-900">
                  {selectedOptionForModal.text}
                </p>
              </div>

              {/* Tabel Siswa Pemilih */}
              <div>
                <h5 className="text-xs font-bold text-slate-700 mb-2 flex items-center gap-1.5">
                  <Users className="w-3.5 h-3.5 text-indigo-600" />
                  Siswa yang Pernah Memilih Opsi Ini:
                </h5>

                {!selectedOptionForModal.students || selectedOptionForModal.students.length === 0 ? (
                  <div className="p-6 text-center text-slate-400 text-xs border border-dashed rounded-xl">
                    Belum ada siswa yang memilih opsi ini.
                  </div>
                ) : (
                  <div className="overflow-x-auto rounded-xl border border-slate-200">
                    <table className="w-full text-left border-collapse text-xs">
                      <thead>
                        <tr className="bg-slate-50 border-b border-slate-200 text-slate-600 font-semibold">
                          <th className="py-2.5 px-3">Nama Siswa</th>
                          <th className="py-2.5 px-3">Username</th>
                          <th className="py-2.5 px-3 text-center">Jumlah Dipilih</th>
                          <th className="py-2.5 px-3 text-right">Terakhir Dipilih</th>
                        </tr>
                      </thead>
                      <tbody className="divide-y divide-slate-100 bg-white">
                        {selectedOptionForModal.students.map((st) => (
                          <tr key={st.student_id} className="hover:bg-slate-50/70">
                            <td className="py-2.5 px-3 font-semibold text-slate-800">
                              {st.student_name}
                            </td>
                            <td className="py-2.5 px-3 text-slate-500">{st.username}</td>
                            <td className="py-2.5 px-3 text-center font-bold text-indigo-600">
                              {st.attempts} kali
                            </td>
                            <td className="py-2.5 px-3 text-right text-slate-500">
                              {st.last_attempt
                                ? new Date(st.last_attempt).toLocaleString('id-ID', {
                                    dateStyle: 'medium',
                                    timeStyle: 'short',
                                  })
                                : '-'}
                            </td>
                          </tr>
                        ))}
                      </tbody>
                    </table>
                  </div>
                )}
              </div>

              <div className="flex justify-end pt-3 border-t border-slate-100">
                <Button
                  size="sm"
                  variant="outline"
                  onClick={() => setSelectedOptionForModal(null)}
                  className="text-xs"
                >
                  Tutup
                </Button>
              </div>
            </div>
          )}
        </Modal>
      </div>
    </ProtectedRoute>
  )
}
