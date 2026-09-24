import { useState, useEffect } from 'react'
import { Modal } from '../../../components/Modal'
import { Button } from '../../../components/ui/button'
import { Input } from '../../../components/ui/input'
import { Label } from '../../../components/ui/label'
import { Textarea } from '../../../components/ui/textarea'
import type { Material } from '../types'
import { materialSchema } from '../types'
import { useCreateMaterial, useUpdateMaterial } from '../api'

interface MaterialModalProps {
  isOpen: boolean
  onClose: () => void
  materialToEdit?: Material | null
}

export function MaterialModal({ isOpen, onClose, materialToEdit }: MaterialModalProps) {
  const isEditing = !!materialToEdit
  const [title, setTitle] = useState('')
  const [content, setContent] = useState('')
  const [errors, setErrors] = useState<Record<string, string>>({})
  const [errorMsg, setErrorMsg] = useState('')

  const createMaterialMutation = useCreateMaterial()
  const updateMaterialMutation = useUpdateMaterial()

  useEffect(() => {
    if (materialToEdit) {
      setTitle(materialToEdit.title)
      setContent(materialToEdit.content)
    } else {
      setTitle('')
      setContent('')
    }
    setErrors({})
    setErrorMsg('')
  }, [materialToEdit, isOpen])

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault()
    setErrors({})
    setErrorMsg('')

    const validation = materialSchema.safeParse({ title, content })
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
      if (isEditing && materialToEdit) {
        await updateMaterialMutation.mutateAsync({
          id: materialToEdit.id,
          data: validation.data,
        })
      } else {
        await createMaterialMutation.mutateAsync(validation.data)
      }
      onClose()
    } catch (err: any) {
      setErrorMsg(err.message || 'Gagal menyimpan materi')
    }
  }

  const isPending = createMaterialMutation.isPending || updateMaterialMutation.isPending

  return (
    <Modal
      isOpen={isOpen}
      onClose={onClose}
      title={isEditing ? 'Ubah Materi' : 'Buat Materi Baru'}
      description="Tuliskan bacaan atau topik argumentasi yang akan dipelajari siswa"
      maxWidth="xl"
    >
      <form onSubmit={handleSubmit} className="space-y-4">
        {errorMsg && (
          <div className="p-3 text-xs bg-red-50 border border-red-200 text-red-700 rounded-lg">
            {errorMsg}
          </div>
        )}

        <div className="space-y-1.5">
          <Label htmlFor="material-title">Judul Materi</Label>
          <Input
            id="material-title"
            placeholder="misal: Pentingnya Berpikir Kritis dalam Era Digital"
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
          <Label htmlFor="material-content">Konten / Teks Bacaan</Label>
          <Textarea
            id="material-content"
            rows={8}
            placeholder="Tuliskan artikel atau bacaan pendukung di sini..."
            value={content}
            onChange={(e) => {
              setContent(e.target.value)
              if (errors.content) setErrors((prev) => ({ ...prev, content: '' }))
            }}
            disabled={isPending}
          />
          {errors.content && <p className="text-xs text-red-500">{errors.content}</p>}
        </div>

        <div className="flex items-center justify-end space-x-2 pt-4 border-t border-slate-100">
          <Button type="button" variant="outline" onClick={onClose} disabled={isPending}>
            Batal
          </Button>
          <Button type="submit" disabled={isPending}>
            {isPending ? 'Menyimpan...' : isEditing ? 'Simpan' : 'Buat Materi'}
          </Button>
        </div>
      </form>
    </Modal>
  )
}
