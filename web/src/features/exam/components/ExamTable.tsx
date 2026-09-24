import { useState } from 'react'
import { Plus, Search, Edit3, Trash2, KeyRound, BookOpen } from 'lucide-react'
import type { Exam } from '../types'
import { Badge } from '../../../components/ui/badge'
import { Button } from '../../../components/ui/button'
import { Input } from '../../../components/ui/input'
import { ConfirmDialog } from '../../../components/ConfirmDialog'
import { ExamModal } from './ExamModal'
import { ExamAccessModal } from './ExamAccessModal'
import { useDeleteExam, useSetExamStatus } from '../api'

interface ExamTableProps {
  exams: Exam[]
  isLoading: boolean
  currentUserRole: string
  currentUserId: number
}

export function ExamTable({
  exams,
  isLoading,
  currentUserRole,
  currentUserId,
}: ExamTableProps) {
  const [searchTerm, setSearchTerm] = useState('')
  const [isExamModalOpen, setIsExamModalOpen] = useState(false)
  const [selectedExam, setSelectedExam] = useState<Exam | null>(null)

  const [isAccessModalOpen, setIsAccessModalOpen] = useState(false)
  const [accessExam, setAccessExam] = useState<Exam | null>(null)

  // Confirm delete state
  const [examToDelete, setExamToDelete] = useState<Exam | null>(null)
  const [actionError, setActionError] = useState<string | null>(null)

  const deleteExamMutation = useDeleteExam()
  const setStatusMutation = useSetExamStatus()

  const filteredExams = exams.filter((e) =>
    e.title.toLowerCase().includes(searchTerm.toLowerCase())
  )

  const handleCreate = () => {
    setSelectedExam(null)
    setIsExamModalOpen(true)
  }

  const handleEdit = (e: Exam) => {
    setSelectedExam(e)
    setIsExamModalOpen(true)
  }

  const handleManageAccess = (e: Exam) => {
    setAccessExam(e)
    setIsAccessModalOpen(true)
  }

  const handleToggleStatus = async (e: Exam) => {
    try {
      await setStatusMutation.mutateAsync({
        id: e.id,
        isActive: !e.is_active,
      })
    } catch (err: any) {
      setActionError(err.message || 'Gagal mengubah status ujian')
    }
  }

  const handleDeleteConfirm = async () => {
    if (!examToDelete) return
    setActionError(null)
    try {
      await deleteExamMutation.mutateAsync(examToDelete.id)
      setExamToDelete(null)
    } catch (err: any) {
      setActionError(err.message || 'Gagal menghapus ujian')
    }
  }

  const canManageExam = (e: Exam) => {
    if (currentUserRole === 'admin') return true
    if (currentUserRole === 'asesor' && e.owner_id === currentUserId) return true
    return false
  }

  return (
    <div className="space-y-4">
      {actionError && (
        <div className="p-3 text-xs bg-red-50 border border-red-200 text-red-700 rounded-lg flex items-center justify-between">
          <span>{actionError}</span>
          <button
            type="button"
            onClick={() => setActionError(null)}
            className="text-red-500 hover:text-red-800 font-bold ml-2 cursor-pointer"
          >
            &times;
          </button>
        </div>
      )}

      {/* Action Toolbar */}
      <div className="flex flex-col sm:flex-row gap-3 items-stretch sm:items-center justify-between">
        <div className="relative flex-1 max-w-md">
          <Search className="w-4 h-4 absolute left-3 top-1/2 -translate-y-1/2 text-slate-400" />
          <Input
            placeholder="Cari judul ujian..."
            value={searchTerm}
            onChange={(e) => setSearchTerm(e.target.value)}
            className="pl-9 bg-white"
          />
        </div>

        <Button onClick={handleCreate} className="gap-1.5 cursor-pointer">
          <Plus className="w-4 h-4" />
          <span>Buat Ujian Baru</span>
        </Button>
      </div>

      {/* Table Container */}
      <div className="bg-white rounded-xl border border-slate-200 shadow-xs overflow-hidden">
        <div className="overflow-x-auto">
          <table className="w-full text-left text-sm text-slate-600">
            <thead className="bg-slate-50 border-b border-slate-200 text-xs font-semibold text-slate-500 uppercase tracking-wider">
              <tr>
                <th className="py-3.5 px-4">Judul Ujian</th>
                <th className="py-3.5 px-4">Materi Terkait</th>
                <th className="py-3.5 px-4">Argumen / Sesi</th>
                <th className="py-3.5 px-4">Status</th>
                <th className="py-3.5 px-4">Akses Diberikan</th>
                <th className="py-3.5 px-4 text-right">Aksi</th>
              </tr>
            </thead>
            <tbody className="divide-y divide-slate-100">
              {isLoading ? (
                <tr>
                  <td colSpan={6} className="py-8 text-center text-slate-400">
                    Memuat daftar ujian...
                  </td>
                </tr>
              ) : filteredExams.length === 0 ? (
                <tr>
                  <td colSpan={6} className="py-8 text-center text-slate-400">
                    Belum ada paket ujian yang dibuat.
                  </td>
                </tr>
              ) : (
                filteredExams.map((e) => {
                  const isOwnerOrAdmin = canManageExam(e)
                  const classCount = e.classes?.length || 0
                  const studentCount = e.students?.length || 0

                  return (
                    <tr key={e.id} className="hover:bg-slate-50/70 transition-colors">
                      <td className="py-3.5 px-4 font-semibold text-slate-900">
                        {e.title}
                      </td>
                      <td className="py-3.5 px-4 text-xs text-slate-600">
                        <div className="flex items-center gap-1.5">
                          <BookOpen className="w-3.5 h-3.5 text-slate-400 shrink-0" />
                          <span className="truncate max-w-[180px]">
                            {e.material?.title || `Materi #${e.material_id}`}
                          </span>
                        </div>
                      </td>
                      <td className="py-3.5 px-4 text-xs font-mono text-slate-600">
                        {e.arguments_per_session} per sesi
                      </td>
                      <td className="py-3.5 px-4">
                        <button
                          type="button"
                          onClick={() => isOwnerOrAdmin && handleToggleStatus(e)}
                          className={`cursor-pointer transition-opacity ${
                            !isOwnerOrAdmin ? 'cursor-default' : 'hover:opacity-80'
                          }`}
                          title={isOwnerOrAdmin ? 'Klik untuk toggle aktif/nonaktif' : ''}
                        >
                          {e.is_active ? (
                            <Badge variant="outline" className="bg-emerald-50 text-emerald-700 border-emerald-200">
                              Aktif
                            </Badge>
                          ) : (
                            <Badge variant="destructive">Nonaktif</Badge>
                          )}
                        </button>
                      </td>
                      <td className="py-3.5 px-4 text-xs">
                        <span className="text-slate-700 font-medium">
                          {classCount} Kelas, {studentCount} Siswa
                        </span>
                      </td>
                      <td className="py-3.5 px-4 text-right space-x-1.5">
                        <Button
                          variant="outline"
                          size="xs"
                          onClick={() => handleManageAccess(e)}
                          className="gap-1 cursor-pointer text-indigo-600 hover:text-indigo-700 hover:bg-indigo-50"
                        >
                          <KeyRound className="w-3 h-3" />
                          <span>Hak Akses</span>
                        </Button>

                        {isOwnerOrAdmin && (
                          <>
                            <Button
                              variant="outline"
                              size="xs"
                              onClick={() => handleEdit(e)}
                              className="cursor-pointer"
                            >
                              <Edit3 className="w-3 h-3" />
                            </Button>
                            <Button
                              variant="ghost"
                              size="xs"
                              onClick={() => setExamToDelete(e)}
                              className="text-red-500 hover:text-red-700 hover:bg-red-50 cursor-pointer"
                              title="Hapus Ujian"
                            >
                              <Trash2 className="w-3.5 h-3.5" />
                            </Button>
                          </>
                        )}
                      </td>
                    </tr>
                  )
                })
              )}
            </tbody>
          </table>
        </div>
      </div>

      <ExamModal
        isOpen={isExamModalOpen}
        onClose={() => setIsExamModalOpen(false)}
        examToEdit={selectedExam}
      />

      <ExamAccessModal
        isOpen={isAccessModalOpen}
        onClose={() => setIsAccessModalOpen(false)}
        targetExam={accessExam}
      />

      {/* Confirm Dialog for Exam Deletion */}
      <ConfirmDialog
        isOpen={!!examToDelete}
        onClose={() => setExamToDelete(null)}
        onConfirm={handleDeleteConfirm}
        title="Hapus Paket Ujian"
        description={`Apakah Anda yakin ingin menghapus paket ujian "${examToDelete?.title}"?`}
        confirmText="Hapus Ujian"
        variant="destructive"
        isLoading={deleteExamMutation.isPending}
      />
    </div>
  )
}
