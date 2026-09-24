import { useState, useEffect } from 'react'
import { Modal } from '../../../components/Modal'
import { Button } from '../../../components/ui/button'
import { Input } from '../../../components/ui/input'
import { Label } from '../../../components/ui/label'
import {
  Select,
  SelectContent,
  SelectItem,
  SelectTrigger,
  SelectValue,
} from '../../../components/ui/select'
import type { Exam } from '../types'
import { examSchema } from '../types'
import { useCreateExam, useUpdateExam } from '../api'
import { useMaterials } from '../../material/api'

interface ExamModalProps {
  isOpen: boolean
  onClose: () => void
  examToEdit?: Exam | null
}

export function ExamModal({ isOpen, onClose, examToEdit }: ExamModalProps) {
  const isEditing = !!examToEdit
  const { data: materials = [] } = useMaterials()

  const [title, setTitle] = useState('')
  const [materialId, setMaterialId] = useState<string>('')
  const [argumentsPerSession, setArgumentsPerSession] = useState<number>(3)
  const [isActive, setIsActive] = useState<boolean>(true)
  const [errors, setErrors] = useState<Record<string, string>>({})
  const [generalError, setGeneralError] = useState('')

  const createMutation = useCreateExam()
  const updateMutation = useUpdateExam()

  useEffect(() => {
    if (examToEdit) {
      setTitle(examToEdit.title)
      setMaterialId(String(examToEdit.material_id))
      setArgumentsPerSession(examToEdit.arguments_per_session || 3)
      setIsActive(examToEdit.is_active)
    } else {
      setTitle('')
      setMaterialId(materials.length > 0 ? String(materials[0].id) : '')
      setArgumentsPerSession(3)
      setIsActive(true)
    }
    setErrors({})
    setGeneralError('')
  }, [examToEdit, isOpen, materials])

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault()
    setErrors({})
    setGeneralError('')

    const validation = examSchema.safeParse({
      title,
      material_id: Number(materialId),
      arguments_per_session: Number(argumentsPerSession),
      is_active: isActive,
    })

    if (!validation.success) {
      const fieldErrors: Record<string, string> = {}
      for (const issue of validation.error.issues) {
        if (issue.path[0]) {
          fieldErrors[String(issue.path[0])] = issue.message
        }
      }
      setErrors(fieldErrors)
      return
    }

    try {
      if (isEditing && examToEdit) {
        await updateMutation.mutateAsync({
          id: examToEdit.id,
          data: validation.data,
        })
      } else {
        await createMutation.mutateAsync(validation.data)
      }
      onClose()
    } catch (err: any) {
      setGeneralError(err.message || 'Gagal menyimpan ujian')
    }
  }

  const isPending = createMutation.isPending || updateMutation.isPending

  return (
    <Modal
      isOpen={isOpen}
      onClose={onClose}
      title={isEditing ? 'Ubah Data Ujian' : 'Buat Paket Ujian Baru'}
      description="Tentukan materi bacaan dan kuota argumen yang diacak untuk setiap sesi pengerjaan"
      maxWidth="md"
    >
      <form onSubmit={handleSubmit} className="space-y-4">
        {generalError && (
          <div className="p-3 text-xs bg-red-50 border border-red-200 text-red-700 rounded-lg">
            {generalError}
          </div>
        )}

        <div className="space-y-1.5">
          <Label htmlFor="exam-title">Judul Latihan / Ujian</Label>
          <Input
            id="exam-title"
            placeholder="misal: Latihan Argumentasi Toulmin 1"
            value={title}
            onChange={(e) => {
              setTitle(e.target.value)
              if (errors.title) setErrors((prev) => ({ ...prev, title: '' }))
            }}
            disabled={isPending}
            autoFocus
          />
          {errors.title && <p className="text-xs text-red-500">{errors.title}</p>}
        </div>

        <div className="space-y-1.5">
          <Label htmlFor="exam-material">Pilih Materi Bacaan</Label>
          <Select
            value={materialId}
            onValueChange={(val) => {
              setMaterialId(val)
              if (errors.material_id) setErrors((prev) => ({ ...prev, material_id: '' }))
            }}
            disabled={isPending || materials.length === 0}
          >
            <SelectTrigger id="exam-material" className="w-full">
              <SelectValue
                placeholder={
                  materials.length === 0
                    ? 'Belum ada materi dibuat'
                    : '-- Pilih Materi Bacaan --'
                }
              />
            </SelectTrigger>
            <SelectContent position="popper">
              {materials.map((m) => (
                <SelectItem key={m.id} value={String(m.id)}>
                  {m.title} ({m.arguments?.length || 0} Argumen)
                </SelectItem>
              ))}
            </SelectContent>
          </Select>
          {errors.material_id && <p className="text-xs text-red-500">{errors.material_id}</p>}
        </div>

        <div className="space-y-1.5">
          <Label htmlFor="exam-args-per-sess">Argumen per Sesi (Maksimal diacak)</Label>
          <Input
            id="exam-args-per-sess"
            type="number"
            min={1}
            max={50}
            value={argumentsPerSession}
            onChange={(e) => {
              setArgumentsPerSession(Number(e.target.value))
              if (errors.arguments_per_session) {
                setErrors((prev) => ({ ...prev, arguments_per_session: '' }))
              }
            }}
            disabled={isPending}
          />
          {errors.arguments_per_session && (
            <p className="text-xs text-red-500">{errors.arguments_per_session}</p>
          )}
          <p className="text-[11px] text-slate-400">
            Sistem akan mengacak sejumlah argumen ini setiap kali siswa memulai sesi latihan baru.
          </p>
        </div>

        <div className="flex items-center space-x-2 pt-1">
          <input
            type="checkbox"
            id="exam-is-active"
            checked={isActive}
            onChange={(e) => setIsActive(e.target.checked)}
            disabled={isPending}
            className="w-4 h-4 rounded text-indigo-600 border-slate-300 focus:ring-indigo-500 cursor-pointer"
          />
          <Label htmlFor="exam-is-active" className="cursor-pointer text-sm font-medium">
            Ujian Aktif (Dapat Dikerjakan Siswa)
          </Label>
        </div>

        <div className="flex items-center justify-end space-x-2 pt-4 border-t border-slate-100">
          <Button type="button" variant="outline" onClick={onClose} disabled={isPending}>
            Batal
          </Button>
          <Button type="submit" disabled={isPending}>
            {isPending ? 'Menyimpan...' : isEditing ? 'Simpan' : 'Buat Ujian'}
          </Button>
        </div>
      </form>
    </Modal>
  )
}
