import { useState, useEffect } from 'react'
import { Modal } from '../../../components/Modal'
import { Button } from '../../../components/ui/button'
import { Input } from '../../../components/ui/input'
import { Label } from '../../../components/ui/label'
import type { Class } from '../types'
import { classSchema } from '../types'
import { useCreateClass, useUpdateClass } from '../api'

interface ClassModalProps {
  isOpen: boolean
  onClose: () => void
  classToEdit?: Class | null
}

export function ClassModal({ isOpen, onClose, classToEdit }: ClassModalProps) {
  const isEditing = !!classToEdit
  const [name, setName] = useState('')
  const [fieldError, setFieldError] = useState('')
  const [errorMsg, setErrorMsg] = useState('')

  const createClassMutation = useCreateClass()
  const updateClassMutation = useUpdateClass()

  useEffect(() => {
    if (classToEdit) {
      setName(classToEdit.name)
    } else {
      setName('')
    }
    setFieldError('')
    setErrorMsg('')
  }, [classToEdit, isOpen])

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault()
    setFieldError('')
    setErrorMsg('')

    const validation = classSchema.safeParse({ name })
    if (!validation.success) {
      setFieldError(validation.error.issues[0]?.message || 'Input tidak valid')
      return
    }

    try {
      if (isEditing && classToEdit) {
        await updateClassMutation.mutateAsync({
          id: classToEdit.id,
          data: { name: validation.data.name },
        })
      } else {
        await createClassMutation.mutateAsync({ name: validation.data.name })
      }
      onClose()
    } catch (err: any) {
      setErrorMsg(err.message || 'Gagal menyimpan kelas')
    }
  }

  const isPending = createClassMutation.isPending || updateClassMutation.isPending

  return (
    <Modal
      isOpen={isOpen}
      onClose={onClose}
      title={isEditing ? 'Ubah Nama Kelas' : 'Buat Kelas Baru'}
      description="Kelas digunakan untuk mengelompokkan siswa dan memberikan akses ujian"
      maxWidth="md"
    >
      <form onSubmit={handleSubmit} className="space-y-4">
        {errorMsg && (
          <div className="p-3 text-xs bg-red-50 border border-red-200 text-red-700 rounded-lg">
            {errorMsg}
          </div>
        )}

        <div className="space-y-1.5">
          <Label htmlFor="class-name">Nama Kelas</Label>
          <Input
            id="class-name"
            placeholder="misal: Kelas Argumentasi X-A"
            value={name}
            onChange={(e) => {
              setName(e.target.value)
              if (fieldError) setFieldError('')
            }}
            disabled={isPending}
            autoFocus
          />
          {fieldError && <p className="text-xs text-red-500">{fieldError}</p>}
        </div>

        <div className="flex items-center justify-end space-x-2 pt-4 border-t border-slate-100">
          <Button type="button" variant="outline" onClick={onClose} disabled={isPending}>
            Batal
          </Button>
          <Button type="submit" disabled={isPending}>
            {isPending ? 'Menyimpan...' : isEditing ? 'Simpan' : 'Buat Kelas'}
          </Button>
        </div>
      </form>
    </Modal>
  )
}
