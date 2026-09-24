import { useQuery, useMutation, useQueryClient } from '@tanstack/react-query'
import { toast } from 'sonner'
import { apiClient } from '../../lib/api-client'
import type { User, CreateUserInput, UpdateUserInput } from './types'

export const userKeys = {
  all: ['users'] as const,
  list: (role?: string) => [...userKeys.all, 'list', role] as const,
  detail: (id: number) => [...userKeys.all, 'detail', id] as const,
}

export function useUsers(role?: string) {
  return useQuery({
    queryKey: userKeys.list(role),
    queryFn: () => {
      const params = role ? { role } : undefined
      return apiClient<User[]>('/users', { params })
    },
  })
}

export function useCreateUser() {
  const queryClient = useQueryClient()
  return useMutation({
    mutationFn: (data: CreateUserInput) =>
      apiClient<User>('/users', {
        method: 'POST',
        body: JSON.stringify(data),
      }),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: userKeys.all })
      toast.success('Pengguna berhasil dibuat')
    },
    onError: (err: any) => {
      toast.error(err?.message || 'Gagal membuat pengguna')
    },
  })
}

export function useUpdateUser() {
  const queryClient = useQueryClient()
  return useMutation({
    mutationFn: ({ id, data }: { id: number; data: UpdateUserInput }) =>
      apiClient<User>(`/users/${id}`, {
        method: 'PATCH',
        body: JSON.stringify(data),
      }),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: userKeys.all })
      toast.success('Data pengguna berhasil diperbarui')
    },
    onError: (err: any) => {
      toast.error(err?.message || 'Gagal memperbarui data pengguna')
    },
  })
}

export function useDeleteUser() {
  const queryClient = useQueryClient()
  return useMutation({
    mutationFn: (id: number) =>
      apiClient<{ message: string }>(`/users/${id}`, {
        method: 'DELETE',
      }),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: userKeys.all })
      toast.success('Pengguna berhasil dihapus')
    },
    onError: (err: any) => {
      toast.error(err?.message || 'Gagal menghapus pengguna')
    },
  })
}
