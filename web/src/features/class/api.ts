import { useQuery, useMutation, useQueryClient } from '@tanstack/react-query'
import { toast } from 'sonner'
import { apiClient } from '../../lib/api-client'
import type { Class, User, CreateClassInput, UpdateClassInput, AddMemberInput } from './types'

export const classKeys = {
  all: ['classes'] as const,
  list: () => [...classKeys.all, 'list'] as const,
  detail: (id: number) => [...classKeys.all, 'detail', id] as const,
  members: (id: number) => [...classKeys.all, 'members', id] as const,
}

export function useClasses() {
  return useQuery({
    queryKey: classKeys.list(),
    queryFn: () => apiClient<Class[]>('/classes'),
  })
}

export function useClass(id: number) {
  return useQuery({
    queryKey: classKeys.detail(id),
    queryFn: () => apiClient<Class>(`/classes/${id}`),
    enabled: id > 0,
  })
}

export function useClassMembers(id: number) {
  return useQuery({
    queryKey: classKeys.members(id),
    queryFn: () => apiClient<User[]>(`/classes/${id}/members`),
    enabled: id > 0,
  })
}

export function useCreateClass() {
  const queryClient = useQueryClient()
  return useMutation({
    mutationFn: (data: CreateClassInput) =>
      apiClient<Class>('/classes', {
        method: 'POST',
        body: JSON.stringify(data),
      }),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: classKeys.all })
      toast.success('Kelas berhasil dibuat')
    },
    onError: (err: any) => {
      toast.error(err?.message || 'Gagal membuat kelas')
    },
  })
}

export function useUpdateClass() {
  const queryClient = useQueryClient()
  return useMutation({
    mutationFn: ({ id, data }: { id: number; data: UpdateClassInput }) =>
      apiClient<Class>(`/classes/${id}`, {
        method: 'PATCH',
        body: JSON.stringify(data),
      }),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: classKeys.all })
      toast.success('Kelas berhasil diperbarui')
    },
    onError: (err: any) => {
      toast.error(err?.message || 'Gagal memperbarui kelas')
    },
  })
}

export function useDeleteClass() {
  const queryClient = useQueryClient()
  return useMutation({
    mutationFn: (id: number) =>
      apiClient<{ message: string }>(`/classes/${id}`, {
        method: 'DELETE',
      }),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: classKeys.all })
      toast.success('Kelas berhasil dihapus')
    },
    onError: (err: any) => {
      toast.error(err?.message || 'Gagal menghapus kelas')
    },
  })
}

export function useAddClassMember() {
  const queryClient = useQueryClient()
  return useMutation({
    mutationFn: ({ classId, data }: { classId: number; data: AddMemberInput }) =>
      apiClient<{ message: string }>(`/classes/${classId}/members`, {
        method: 'POST',
        body: JSON.stringify(data),
      }),
    onSuccess: (_, { classId }) => {
      queryClient.invalidateQueries({ queryKey: classKeys.members(classId) })
      queryClient.invalidateQueries({ queryKey: classKeys.list() })
      toast.success('Anggota berhasil ditambahkan ke kelas')
    },
    onError: (err: any) => {
      toast.error(err?.message || 'Gagal menambahkan anggota')
    },
  })
}

export function useRemoveClassMember() {
  const queryClient = useQueryClient()
  return useMutation({
    mutationFn: ({ classId, userId }: { classId: number; userId: number }) =>
      apiClient<{ message: string }>(`/classes/${classId}/members/${userId}`, {
        method: 'DELETE',
      }),
    onSuccess: (_, { classId }) => {
      queryClient.invalidateQueries({ queryKey: classKeys.members(classId) })
      queryClient.invalidateQueries({ queryKey: classKeys.list() })
      toast.success('Anggota berhasil dikeluarkan dari kelas')
    },
    onError: (err: any) => {
      toast.error(err?.message || 'Gagal mengeluarkan anggota')
    },
  })
}
