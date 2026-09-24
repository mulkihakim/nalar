import { useState } from 'react'
import { Plus, Search, Edit3, Trash2, Users } from 'lucide-react'
import type { Class } from '../types'
import { Button } from '../../../components/ui/button'
import { Input } from '../../../components/ui/input'
import { ConfirmDialog } from '../../../components/ConfirmDialog'
import { ClassModal } from './ClassModal'
import { MembersModal } from './MembersModal'
import { useDeleteClass } from '../api'

interface ClassTableProps {
  classes: Class[]
  isLoading: boolean
  currentUserRole: string
  currentUserId: number
}

export function ClassTable({
  classes,
  isLoading,
  currentUserRole,
  currentUserId,
}: ClassTableProps) {
  const [searchTerm, setSearchTerm] = useState('')
  const [isClassModalOpen, setIsClassModalOpen] = useState(false)
  const [isMembersModalOpen, setIsMembersModalOpen] = useState(false)
  const [selectedClass, setSelectedClass] = useState<Class | null>(null)
  const [classToDelete, setClassToDelete] = useState<Class | null>(null)
  const [deleteError, setDeleteError] = useState<string | null>(null)

  const deleteClassMutation = useDeleteClass()

  const filteredClasses = classes.filter((c) =>
    c.name.toLowerCase().includes(searchTerm.toLowerCase())
  )

  const handleCreate = () => {
    setSelectedClass(null)
    setIsClassModalOpen(true)
  }

  const handleEdit = (c: Class) => {
    setSelectedClass(c)
    setIsClassModalOpen(true)
  }

  const handleManageMembers = (c: Class) => {
    setSelectedClass(c)
    setIsMembersModalOpen(true)
  }

  const handleDeleteConfirm = async () => {
    if (!classToDelete) return
    setDeleteError(null)

    try {
      await deleteClassMutation.mutateAsync(classToDelete.id)
      setClassToDelete(null)
    } catch (err: any) {
      setDeleteError(err.message || 'Gagal menghapus kelas')
    }
  }

  const canManageClass = (c: Class) => {
    if (currentUserRole === 'admin') return true
    if (currentUserRole === 'asesor' && c.owner_id === currentUserId) return true
    return false
  }

  return (
    <div className="space-y-4">
      {deleteError && (
        <div className="p-3 text-xs bg-red-50 border border-red-200 text-red-700 rounded-lg flex items-center justify-between">
          <span>{deleteError}</span>
          <button
            type="button"
            onClick={() => setDeleteError(null)}
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
            placeholder="Cari nama kelas..."
            value={searchTerm}
            onChange={(e) => setSearchTerm(e.target.value)}
            className="pl-9 bg-white"
          />
        </div>

        <Button onClick={handleCreate} className="gap-1.5 cursor-pointer">
          <Plus className="w-4 h-4" />
          <span>Buat Kelas Baru</span>
        </Button>
      </div>

      {/* Table Container */}
      <div className="bg-white rounded-xl border border-slate-200 shadow-xs overflow-hidden">
        <div className="overflow-x-auto">
          <table className="w-full text-left text-sm text-slate-600">
            <thead className="bg-slate-50 border-b border-slate-200 text-xs font-semibold text-slate-500 uppercase tracking-wider">
              <tr>
                <th className="py-3.5 px-4">Nama Kelas</th>
                <th className="py-3.5 px-4">Pemilik (Asesor)</th>
                <th className="py-3.5 px-4">Jumlah Anggota</th>
                <th className="py-3.5 px-4 text-right">Aksi</th>
              </tr>
            </thead>
            <tbody className="divide-y divide-slate-100">
              {isLoading ? (
                <tr>
                  <td colSpan={4} className="py-8 text-center text-slate-400">
                    Memuat daftar kelas...
                  </td>
                </tr>
              ) : filteredClasses.length === 0 ? (
                <tr>
                  <td colSpan={4} className="py-8 text-center text-slate-400">
                    Belum ada kelas yang dibuat.
                  </td>
                </tr>
              ) : (
                filteredClasses.map((c) => {
                  const memberCount = c.members ? c.members.length : 0
                  const isOwnerOrAdmin = canManageClass(c)

                  return (
                    <tr key={c.id} className="hover:bg-slate-50/70 transition-colors">
                      <td className="py-3.5 px-4 font-semibold text-slate-900">
                        {c.name}
                      </td>
                      <td className="py-3.5 px-4 text-xs text-slate-600">
                        {c.owner?.name || `User #${c.owner_id}`}
                      </td>
                      <td className="py-3.5 px-4">
                        <span className="inline-flex items-center gap-1.5 px-2.5 py-0.5 rounded-full text-xs font-medium bg-slate-100 text-slate-700">
                          <Users className="w-3 h-3 text-slate-400" />
                          <span>{memberCount} Siswa</span>
                        </span>
                      </td>
                      <td className="py-3.5 px-4 text-right space-x-1.5">
                        <Button
                          variant="outline"
                          size="xs"
                          onClick={() => handleManageMembers(c)}
                          className="gap-1 cursor-pointer text-indigo-600 hover:text-indigo-700 hover:bg-indigo-50"
                        >
                          <Users className="w-3 h-3" />
                          <span>Kelola Anggota</span>
                        </Button>

                        {isOwnerOrAdmin && (
                          <>
                            <Button
                              variant="outline"
                              size="xs"
                              onClick={() => handleEdit(c)}
                              className="gap-1 cursor-pointer"
                            >
                              <Edit3 className="w-3 h-3" />
                              <span>Ubah</span>
                            </Button>
                            <Button
                              variant="ghost"
                              size="xs"
                              onClick={() => setClassToDelete(c)}
                              className="text-red-500 hover:text-red-700 hover:bg-red-50 cursor-pointer"
                              title="Hapus Kelas"
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

      <ClassModal
        isOpen={isClassModalOpen}
        onClose={() => setIsClassModalOpen(false)}
        classToEdit={selectedClass}
      />

      <MembersModal
        isOpen={isMembersModalOpen}
        onClose={() => setIsMembersModalOpen(false)}
        targetClass={selectedClass}
      />

      <ConfirmDialog
        isOpen={!!classToDelete}
        onClose={() => setClassToDelete(null)}
        onConfirm={handleDeleteConfirm}
        title="Hapus Kelas"
        description={`Apakah Anda yakin ingin menghapus kelas "${classToDelete?.name}"? Siswa yang terdaftar tidak akan dihapus dari sistem.`}
        confirmText="Hapus Kelas"
        variant="destructive"
        isLoading={deleteClassMutation.isPending}
      />
    </div>
  )
}
