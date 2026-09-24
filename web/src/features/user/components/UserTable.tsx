import { useState } from 'react'
import { Plus, Search, Edit3, Trash2 } from 'lucide-react'
import type { User, Role } from '../types'
import { Badge } from '../../../components/ui/badge'
import { Button } from '../../../components/ui/button'
import { Input } from '../../../components/ui/input'
import {
  Select,
  SelectContent,
  SelectItem,
  SelectTrigger,
  SelectValue,
} from '../../../components/ui/select'
import { ConfirmDialog } from '../../../components/ConfirmDialog'
import { UserModal } from './UserModal'
import { useDeleteUser } from '../api'

interface UserTableProps {
  users: User[]
  isLoading: boolean
  currentUserRole: Role
  currentUserId: number
}

export function UserTable({
  users,
  isLoading,
  currentUserRole,
  currentUserId,
}: UserTableProps) {
  const [searchTerm, setSearchTerm] = useState('')
  const [roleFilter, setRoleFilter] = useState<string>('all')
  const [isModalOpen, setIsModalOpen] = useState(false)
  const [selectedUser, setSelectedUser] = useState<User | null>(null)

  // Hard delete confirmation state
  const [userToDelete, setUserToDelete] = useState<User | null>(null)
  const [actionError, setActionError] = useState<string | null>(null)

  const deleteUserMutation = useDeleteUser()

  const filteredUsers = users.filter((u) => {
    const matchesSearch =
      u.name.toLowerCase().includes(searchTerm.toLowerCase()) ||
      u.username.toLowerCase().includes(searchTerm.toLowerCase())
    const matchesRole = roleFilter === 'all' || u.role === roleFilter
    return matchesSearch && matchesRole
  })

  const handleCreate = () => {
    setSelectedUser(null)
    setIsModalOpen(true)
  }

  const handleEdit = (u: User) => {
    setSelectedUser(u)
    setIsModalOpen(true)
  }

  // Can the current user edit this target user?
  // Admin can edit anyone.
  // Asesor can only edit students they created.
  const canEditUser = (u: User) => {
    if (currentUserRole === 'admin') return true
    if (currentUserRole === 'asesor' && u.role === 'siswa' && u.created_by === currentUserId) return true
    return false
  }

  // Hard delete: Only admin, and cannot delete own account
  const canDeleteUser = (u: User) => {
    return currentUserRole === 'admin' && u.id !== currentUserId
  }

  const handleDeleteConfirm = async () => {
    if (!userToDelete) return
    setActionError(null)

    try {
      await deleteUserMutation.mutateAsync(userToDelete.id)
      setUserToDelete(null)
    } catch (err: any) {
      setActionError(err.message || 'Gagal menghapus pengguna. Pastikan pengguna belum pernah mengerjakan ujian.')
    }
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
        <div className="flex flex-1 items-center gap-2 max-w-md">
          <div className="relative flex-1">
            <Search className="w-4 h-4 absolute left-3 top-1/2 -translate-y-1/2 text-slate-400" />
            <Input
              placeholder="Cari nama atau username..."
              value={searchTerm}
              onChange={(e) => setSearchTerm(e.target.value)}
              className="pl-9 bg-white"
            />
          </div>

          <div className="w-36">
            <Select value={roleFilter} onValueChange={setRoleFilter}>
              <SelectTrigger className="w-full bg-white h-8 text-xs">
                <SelectValue placeholder="Semua Role" />
              </SelectTrigger>
              <SelectContent position="popper">
                <SelectItem value="all">Semua Role</SelectItem>
                <SelectItem value="siswa">Siswa</SelectItem>
                <SelectItem value="asesor">Asesor</SelectItem>
                {currentUserRole === 'admin' && <SelectItem value="admin">Admin</SelectItem>}
              </SelectContent>
            </Select>
          </div>
        </div>

        <Button onClick={handleCreate} className="gap-1.5 cursor-pointer">
          <Plus className="w-4 h-4" />
          <span>Tambah Pengguna</span>
        </Button>
      </div>

      {/* Table Container */}
      <div className="bg-white rounded-xl border border-slate-200 shadow-xs overflow-hidden">
        <div className="overflow-x-auto">
          <table className="w-full text-left text-sm text-slate-600">
            <thead className="bg-slate-50 border-b border-slate-200 text-xs font-semibold text-slate-500 uppercase tracking-wider">
              <tr>
                <th className="py-3.5 px-4">Nama Lengkap</th>
                <th className="py-3.5 px-4">Username</th>
                <th className="py-3.5 px-4">Role</th>
                <th className="py-3.5 px-4">Status</th>
                <th className="py-3.5 px-4 text-right">Aksi</th>
              </tr>
            </thead>
            <tbody className="divide-y divide-slate-100">
              {isLoading ? (
                <tr>
                  <td colSpan={5} className="py-8 text-center text-slate-400">
                    Memuat data pengguna...
                  </td>
                </tr>
              ) : filteredUsers.length === 0 ? (
                <tr>
                  <td colSpan={5} className="py-8 text-center text-slate-400">
                    Tidak ada data pengguna yang sesuai.
                  </td>
                </tr>
              ) : (
                filteredUsers.map((u) => (
                  <tr key={u.id} className="hover:bg-slate-50/70 transition-colors">
                    <td className="py-3.5 px-4 font-medium text-slate-900">
                      <div className="flex items-center gap-2">
                        <div className="w-7 h-7 rounded-full bg-indigo-100 text-indigo-700 flex items-center justify-center font-bold text-xs">
                          {u.name.charAt(0).toUpperCase()}
                        </div>
                        <span>{u.name}</span>
                      </div>
                    </td>
                    <td className="py-3.5 px-4 font-mono text-xs text-slate-600">
                      @{u.username}
                    </td>
                    <td className="py-3.5 px-4">
                      <Badge
                        variant={
                          u.role === 'admin'
                            ? 'destructive'
                            : u.role === 'asesor'
                            ? 'default'
                            : 'secondary'
                        }
                      >
                        {u.role}
                      </Badge>
                    </td>
                    <td className="py-3.5 px-4">
                      {u.is_active ? (
                        <Badge variant="outline" className="bg-emerald-50 text-emerald-700 border-emerald-200">
                          Aktif
                        </Badge>
                      ) : (
                        <Badge variant="destructive">Nonaktif</Badge>
                      )}
                    </td>
                    <td className="py-3.5 px-4 text-right">
                      <div className="flex items-center justify-end gap-1.5">
                        {canEditUser(u) ? (
                          <Button
                            variant="outline"
                            size="xs"
                            onClick={() => handleEdit(u)}
                            className="gap-1 cursor-pointer"
                          >
                            <Edit3 className="w-3 h-3" />
                            <span>Ubah</span>
                          </Button>
                        ) : (
                          <span className="text-xs text-slate-300 italic">Hanya Baca</span>
                        )}

                        {canDeleteUser(u) && (
                          <Button
                            variant="outline"
                            size="xs"
                            onClick={() => setUserToDelete(u)}
                            className="gap-1 text-red-600 hover:text-red-700 hover:bg-red-50 hover:border-red-200 cursor-pointer"
                            title="Hapus Pengguna Permanen"
                          >
                            <Trash2 className="w-3 h-3" />
                            <span className="hidden sm:inline">Hapus</span>
                          </Button>
                        )}
                      </div>
                    </td>
                  </tr>
                ))
              )}
            </tbody>
          </table>
        </div>
      </div>

      <UserModal
        isOpen={isModalOpen}
        onClose={() => setIsModalOpen(false)}
        userToEdit={selectedUser}
        currentUserRole={currentUserRole}
      />

      <ConfirmDialog
        isOpen={!!userToDelete}
        onClose={() => setUserToDelete(null)}
        onConfirm={handleDeleteConfirm}
        title="Hapus Pengguna Permanen"
        description={`Apakah Anda yakin ingin menghapus akun @${userToDelete?.username} (${userToDelete?.name}) secara permanen? Pengguna hanya dapat dihapus jika belum pernah mengerjakan ujian apa pun.`}
        confirmText="Hapus Permanen"
        variant="destructive"
        isLoading={deleteUserMutation.isPending}
      />
    </div>
  )
}
