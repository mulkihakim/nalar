import { useState, useEffect } from 'react'
import { GraduationCap, UserCheck } from 'lucide-react'
import { Modal } from '../../../components/Modal'
import { Button } from '../../../components/ui/button'
import type { Exam } from '../types'
import { useSetExamAccess, useExam } from '../api'
import { useClasses } from '../../class/api'
import { useUsers } from '../../user/api'

interface ExamAccessModalProps {
  isOpen: boolean
  onClose: () => void
  targetExam: Exam | null
}

export function ExamAccessModal({ isOpen, onClose, targetExam }: ExamAccessModalProps) {
  const examId = targetExam?.id || 0
  const { data: fullExam } = useExam(examId)
  const { data: classes = [] } = useClasses()
  const { data: students = [] } = useUsers('siswa')

  const [selectedClassIds, setSelectedClassIds] = useState<number[]>([])
  const [selectedStudentIds, setSelectedStudentIds] = useState<number[]>([])
  const [errorMsg, setErrorMsg] = useState('')

  const setAccessMutation = useSetExamAccess()

  useEffect(() => {
    if (fullExam) {
      setSelectedClassIds(fullExam.classes?.map((c) => c.id) || [])
      setSelectedStudentIds(fullExam.students?.map((s) => s.id) || [])
    } else if (targetExam) {
      setSelectedClassIds(targetExam.classes?.map((c) => c.id) || [])
      setSelectedStudentIds(targetExam.students?.map((s) => s.id) || [])
    }
    setErrorMsg('')
  }, [fullExam, targetExam, isOpen])

  if (!targetExam) return null

  const toggleClass = (id: number) => {
    setSelectedClassIds((prev) =>
      prev.includes(id) ? prev.filter((cid) => cid !== id) : [...prev, id]
    )
  }

  const toggleStudent = (id: number) => {
    setSelectedStudentIds((prev) =>
      prev.includes(id) ? prev.filter((sid) => sid !== id) : [...prev, id]
    )
  }

  const handleSave = async () => {
    setErrorMsg('')
    try {
      await setAccessMutation.mutateAsync({
        id: targetExam.id,
        data: {
          class_ids: selectedClassIds,
          student_ids: selectedStudentIds,
        },
      })
      onClose()
    } catch (err: any) {
      setErrorMsg(err.message || 'Gagal mengatur akses ujian')
    }
  }

  return (
    <Modal
      isOpen={isOpen}
      onClose={onClose}
      title={`Atur Hak Akses: ${targetExam.title}`}
      description="Pilih kelas atau siswa individu yang berhak melihat dan mengerjakan latihan ini"
      maxWidth="xl"
    >
      <div className="space-y-6">
        {errorMsg && (
          <div className="p-3 text-xs bg-red-50 border border-red-200 text-red-700 rounded-lg">
            {errorMsg}
          </div>
        )}

        {/* Akses Kelas */}
        <div className="space-y-2.5">
          <div className="flex items-center justify-between">
            <span className="text-xs font-bold uppercase tracking-wider text-slate-700 flex items-center gap-1.5">
              <GraduationCap className="w-4 h-4 text-indigo-600" />
              <span>Akses Berdasarkan Kelas ({selectedClassIds.length} Terpilih)</span>
            </span>
          </div>

          <div className="border border-slate-200 rounded-xl p-3 bg-slate-50/50 max-h-48 overflow-y-auto space-y-1.5">
            {classes.length === 0 ? (
              <div className="text-xs text-slate-400 text-center py-3">
                Belum ada kelas yang terdaftar.
              </div>
            ) : (
              classes.map((c) => {
                const isChecked = selectedClassIds.includes(c.id)
                return (
                  <label
                    key={c.id}
                    className={`flex items-center justify-between p-2.5 rounded-lg border cursor-pointer transition-colors ${
                      isChecked
                        ? 'bg-indigo-50/80 border-indigo-200 text-indigo-900'
                        : 'bg-white border-slate-200 hover:bg-slate-50'
                    }`}
                  >
                    <div className="flex items-center gap-2.5">
                      <input
                        type="checkbox"
                        checked={isChecked}
                        onChange={() => toggleClass(c.id)}
                        className="w-4 h-4 rounded text-indigo-600 border-slate-300 focus:ring-indigo-500"
                      />
                      <span className="text-sm font-medium">{c.name}</span>
                    </div>
                    <span className="text-xs text-slate-400">
                      {c.members?.length || 0} Siswa
                    </span>
                  </label>
                )
              })
            )}
          </div>
        </div>

        {/* Akses Siswa Individu */}
        <div className="space-y-2.5">
          <div className="flex items-center justify-between">
            <span className="text-xs font-bold uppercase tracking-wider text-slate-700 flex items-center gap-1.5">
              <UserCheck className="w-4 h-4 text-indigo-600" />
              <span>Akses Khusus Siswa Individu ({selectedStudentIds.length} Terpilih)</span>
            </span>
          </div>

          <div className="border border-slate-200 rounded-xl p-3 bg-slate-50/50 max-h-48 overflow-y-auto space-y-1.5">
            {students.length === 0 ? (
              <div className="text-xs text-slate-400 text-center py-3">
                Belum ada siswa yang terdaftar.
              </div>
            ) : (
              students.map((s) => {
                const isChecked = selectedStudentIds.includes(s.id)
                return (
                  <label
                    key={s.id}
                    className={`flex items-center justify-between p-2.5 rounded-lg border cursor-pointer transition-colors ${
                      isChecked
                        ? 'bg-indigo-50/80 border-indigo-200 text-indigo-900'
                        : 'bg-white border-slate-200 hover:bg-slate-50'
                    }`}
                  >
                    <div className="flex items-center gap-2.5">
                      <input
                        type="checkbox"
                        checked={isChecked}
                        onChange={() => toggleStudent(s.id)}
                        className="w-4 h-4 rounded text-indigo-600 border-slate-300 focus:ring-indigo-500"
                      />
                      <span className="text-sm font-medium">{s.name}</span>
                    </div>
                    <span className="text-xs font-mono text-slate-400">@{s.username}</span>
                  </label>
                )
              })
            )}
          </div>
        </div>

        <div className="flex items-center justify-end space-x-2 pt-4 border-t border-slate-100">
          <Button variant="outline" onClick={onClose} disabled={setAccessMutation.isPending}>
            Batal
          </Button>
          <Button onClick={handleSave} disabled={setAccessMutation.isPending}>
            {setAccessMutation.isPending ? 'Menyimpan...' : 'Simpan Hak Akses'}
          </Button>
        </div>
      </div>
    </Modal>
  )
}
