import { useEffect, useState } from 'react'
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
import type { User, Role } from '../types'
import { createUserSchema, updateUserSchema } from '../types'
import { useCreateUser, useUpdateUser } from '../api'

interface UserModalProps {
  isOpen: boolean
  onClose: () => void
  userToEdit?: User | null
  currentUserRole: Role
}

export function UserModal({
  isOpen,
  onClose,
  userToEdit,
  currentUserRole,
}: UserModalProps) {
  const isEditing = !!userToEdit

  const [name, setName] = useState('')
  const [username, setUsername] = useState('')
  const [password, setPassword] = useState('')
  const [role, setRole] = useState<Role>('siswa')
  const [isActive, setIsActive] = useState(true)
  const [errors, setErrors] = useState<Record<string, string>>({})
  const [generalError, setGeneralError] = useState('')

  const createUserMutation = useCreateUser()
  const updateUserMutation = useUpdateUser()

  useEffect(() => {
    if (userToEdit) {
      setName(userToEdit.name)
      setUsername(userToEdit.username)
      setPassword('')
      setRole(userToEdit.role)
      setIsActive(userToEdit.is_active)
    } else {
      setName('')
      setUsername('')
      setPassword('')
      setRole('siswa')
      setIsActive(true)
    }
    setErrors({})
    setGeneralError('')
  }, [userToEdit, isOpen])

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault()
    setErrors({})
    setGeneralError('')

    if (isEditing && userToEdit) {
      const validation = updateUserSchema.safeParse({
        name,
        password: password || undefined,
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
        await updateUserMutation.mutateAsync({
          id: userToEdit.id,
          data: {
            name: validation.data.name,
            is_active: validation.data.is_active,
            ...(validation.data.password ? { password: validation.data.password } : {}),
          },
        })
        onClose()
      } catch (err: any) {
        setGeneralError(err.message || 'Terjadi kesalahan saat memperbarui akun pengguna')
      }
    } else {
      const validation = createUserSchema.safeParse({
        name,
        username,
        password,
        role,
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
        await createUserMutation.mutateAsync(validation.data)
        onClose()
      } catch (err: any) {
        setGeneralError(err.message || 'Terjadi kesalahan saat membuat pengguna baru')
      }
    }
  }

  const isPending = createUserMutation.isPending || updateUserMutation.isPending

  return (
    <Modal
      isOpen={isOpen}
      onClose={onClose}
      title={isEditing ? 'Ubah Data Pengguna' : 'Tambah Pengguna Baru'}
      description={
        isEditing
          ? `Perbarui profil dan status akun ${userToEdit?.username}`
          : 'Lengkapi formulir untuk membuat akun staf atau siswa'
      }
    >
      <form onSubmit={handleSubmit} className="space-y-4">
        {generalError && (
          <div className="p-3 text-xs bg-red-50 border border-red-200 text-red-700 rounded-lg">
            {generalError}
          </div>
        )}

        <div className="space-y-1.5">
          <Label htmlFor="name">Nama Lengkap</Label>
          <Input
            id="name"
            placeholder="misal: Siti Nurhaliza"
            value={name}
            onChange={(e) => {
              setName(e.target.value)
              if (errors.name) setErrors((prev) => ({ ...prev, name: '' }))
            }}
            disabled={isPending}
          />
          {errors.name && <p className="text-xs text-red-500">{errors.name}</p>}
        </div>

        <div className="space-y-1.5">
          <Label htmlFor="username">Username</Label>
          <Input
            id="username"
            placeholder="misal: sitisiswa1"
            value={username}
            onChange={(e) => {
              setUsername(e.target.value)
              if (errors.username) setErrors((prev) => ({ ...prev, username: '' }))
            }}
            disabled={isEditing || isPending}
          />
          {isEditing && (
            <p className="text-[11px] text-slate-400">Username tidak dapat diubah setelah dibuat.</p>
          )}
          {errors.username && <p className="text-xs text-red-500">{errors.username}</p>}
        </div>

        <div className="space-y-1.5">
          <Label htmlFor="password">
            {isEditing ? 'Password Baru (Kosongkan jika tidak diubah)' : 'Password'}
          </Label>
          <Input
            id="password"
            type="password"
            placeholder={isEditing ? '••••••••' : 'Minimal 6 karakter'}
            value={password}
            onChange={(e) => {
              setPassword(e.target.value)
              if (errors.password) setErrors((prev) => ({ ...prev, password: '' }))
            }}
            disabled={isPending}
          />
          {errors.password && <p className="text-xs text-red-500">{errors.password}</p>}
        </div>

        {!isEditing && (
          <div className="space-y-1.5">
            <Label htmlFor="role">Role Pengguna</Label>
            <Select
              value={role}
              onValueChange={(val) => setRole(val as Role)}
              disabled={isPending}
            >
              <SelectTrigger id="role" className="w-full">
                <SelectValue placeholder="Pilih Role" />
              </SelectTrigger>
              <SelectContent position="popper">
                <SelectItem value="siswa">Siswa</SelectItem>
                {currentUserRole === 'admin' && <SelectItem value="asesor">Asesor</SelectItem>}
                {currentUserRole === 'admin' && <SelectItem value="admin">Admin</SelectItem>}
              </SelectContent>
            </Select>
            {errors.role && <p className="text-xs text-red-500">{errors.role}</p>}
          </div>
        )}

        {isEditing && (
          <div className="flex items-center space-x-2 pt-2">
            <input
              type="checkbox"
              id="is_active"
              checked={isActive}
              onChange={(e) => setIsActive(e.target.checked)}
              disabled={isPending}
              className="w-4 h-4 rounded text-indigo-600 border-slate-300 focus:ring-indigo-500 cursor-pointer"
            />
            <Label htmlFor="is_active" className="cursor-pointer text-sm font-medium">
              Akun Aktif (Dapat Login)
            </Label>
          </div>
        )}

        <div className="flex items-center justify-end space-x-2 pt-4 border-t border-slate-100">
          <Button type="button" variant="outline" onClick={onClose} disabled={isPending}>
            Batal
          </Button>
          <Button type="submit" disabled={isPending}>
            {isPending ? 'Menyimpan...' : isEditing ? 'Simpan Perubahan' : 'Buat Pengguna'}
          </Button>
        </div>
      </form>
    </Modal>
  )
}
