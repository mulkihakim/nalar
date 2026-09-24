import { useState } from 'react'
import { Plus, Trash2, UserPlus } from 'lucide-react'
import { Modal } from '../../../components/Modal'
import { Button } from '../../../components/ui/button'
import {
  Select,
  SelectContent,
  SelectItem,
  SelectTrigger,
  SelectValue,
} from '../../../components/ui/select'
import { ConfirmDialog } from '../../../components/ConfirmDialog'
import type { Class, User } from '../types'
import { addMemberSchema } from '../types'
import { useClassMembers, useAddClassMember, useRemoveClassMember } from '../api'
import { useUsers } from '../../user/api'

interface MembersModalProps {
  isOpen: boolean
  onClose: () => void
  targetClass: Class | null
}

export function MembersModal({ isOpen, onClose, targetClass }: MembersModalProps) {
  const classId = targetClass?.id || 0
  const { data: members = [], isLoading: loadingMembers } = useClassMembers(classId)
  const { data: allStudents = [] } = useUsers('siswa')

  const [selectedStudentId, setSelectedStudentId] = useState<string>('')
  const [errorMsg, setErrorMsg] = useState('')
  const [memberToRemove, setMemberToRemove] = useState<User | null>(null)

  const addMemberMutation = useAddClassMember()
  const removeMemberMutation = useRemoveClassMember()

  if (!targetClass) return null

  // Filter students who are not yet members
  const memberIds = new Set(members.map((m) => m.id))
  const availableStudents = allStudents.filter((s) => !memberIds.has(s.id))

  const handleAddMember = async (e: React.FormEvent) => {
    e.preventDefault()
    setErrorMsg('')

    const val = addMemberSchema.safeParse({ user_id: Number(selectedStudentId) })
    if (!val.success) {
      setErrorMsg(val.error.issues[0]?.message || 'Pilih siswa terlebih dahulu')
      return
    }

    try {
      await addMemberMutation.mutateAsync({
        classId: targetClass.id,
        data: { user_id: val.data.user_id },
      })
      setSelectedStudentId('')
    } catch (err: any) {
      setErrorMsg(err.message || 'Gagal menambahkan anggota')
    }
  }

  const handleRemoveConfirm = async () => {
    if (!memberToRemove) return
    try {
      await removeMemberMutation.mutateAsync({
        classId: targetClass.id,
        userId: memberToRemove.id,
      })
      setMemberToRemove(null)
    } catch (err: any) {
      setErrorMsg(err.message || 'Gagal menghapus anggota')
    }
  }

  return (
    <Modal
      isOpen={isOpen}
      onClose={onClose}
      title={`Anggota: ${targetClass.name}`}
      description="Kelola siswa yang terdaftar di kelas ini"
      maxWidth="lg"
    >
      <div className="space-y-6">
        {/* Form Tambah Siswa */}
        <form onSubmit={handleAddMember} className="space-y-3 bg-slate-50 p-4 rounded-xl border border-slate-200">
          <div className="text-xs font-semibold text-slate-700 flex items-center gap-1.5">
            <UserPlus className="w-3.5 h-3.5 text-indigo-600" />
            <span>Daftarkan Siswa ke Kelas Ini</span>
          </div>

          {errorMsg && (
            <div className="p-2 text-xs bg-red-50 border border-red-200 text-red-700 rounded-md">
              {errorMsg}
            </div>
          )}

          <div className="flex gap-2 items-center">
            <div className="flex-1">
              <Select
                value={selectedStudentId}
                onValueChange={setSelectedStudentId}
                disabled={availableStudents.length === 0 || addMemberMutation.isPending}
              >
                <SelectTrigger className="w-full bg-white text-xs">
                  <SelectValue
                    placeholder={
                      availableStudents.length === 0
                        ? 'Tidak ada siswa lain yang tersedia'
                        : 'Pilih Siswa untuk Ditambahkan'
                    }
                  />
                </SelectTrigger>
                <SelectContent position="popper">
                  {availableStudents.map((s) => (
                    <SelectItem key={s.id} value={String(s.id)}>
                      {s.name} (@{s.username})
                    </SelectItem>
                  ))}
                </SelectContent>
              </Select>
            </div>

            <Button
              type="submit"
              disabled={!selectedStudentId || addMemberMutation.isPending}
              className="gap-1 cursor-pointer shrink-0"
              size="sm"
            >
              <Plus className="w-3.5 h-3.5" />
              <span>{addMemberMutation.isPending ? 'Menambah...' : 'Tambah'}</span>
            </Button>
          </div>
        </form>

        {/* Daftar Anggota Terdaftar */}
        <div>
          <div className="text-xs font-semibold uppercase tracking-wider text-slate-400 mb-2">
            Daftar Anggota Saat Ini ({members.length} Siswa)
          </div>

          <div className="border border-slate-200 rounded-xl overflow-hidden divide-y divide-slate-100 max-h-64 overflow-y-auto">
            {loadingMembers ? (
              <div className="p-6 text-center text-xs text-slate-400">Memuat anggota kelas...</div>
            ) : members.length === 0 ? (
              <div className="p-6 text-center text-xs text-slate-400">
                Belum ada siswa di kelas ini. Tambahkan siswa menggunakan formulir di atas.
              </div>
            ) : (
              members.map((m) => (
                <div key={m.id} className="p-3 flex items-center justify-between hover:bg-slate-50 transition-colors">
                  <div className="flex items-center gap-2.5">
                    <div className="w-7 h-7 rounded-full bg-indigo-50 text-indigo-700 flex items-center justify-center font-bold text-xs">
                      {m.name.charAt(0).toUpperCase()}
                    </div>
                    <div>
                      <div className="text-sm font-medium text-slate-900">{m.name}</div>
                      <div className="text-[11px] font-mono text-slate-400">@{m.username}</div>
                    </div>
                  </div>

                  <Button
                    variant="ghost"
                    size="xs"
                    onClick={() => setMemberToRemove(m)}
                    disabled={removeMemberMutation.isPending}
                    className="text-red-500 hover:text-red-700 hover:bg-red-50 cursor-pointer p-1"
                    title="Keluarkan dari kelas"
                  >
                    <Trash2 className="w-3.5 h-3.5" />
                  </Button>
                </div>
              ))
            )}
          </div>
        </div>

        <div className="flex justify-end pt-3 border-t border-slate-100">
          <Button variant="outline" onClick={onClose}>
            Selesai
          </Button>
        </div>
      </div>

      <ConfirmDialog
        isOpen={!!memberToRemove}
        onClose={() => setMemberToRemove(null)}
        onConfirm={handleRemoveConfirm}
        title="Keluarkan Siswa dari Kelas"
        description={`Apakah Anda yakin ingin mengeluarkan ${memberToRemove?.name} dari kelas ini?`}
        confirmText="Keluarkan"
        variant="destructive"
        isLoading={removeMemberMutation.isPending}
      />
    </Modal>
  )
}
