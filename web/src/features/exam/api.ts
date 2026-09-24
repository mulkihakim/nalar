import { useQuery, useMutation, useQueryClient } from '@tanstack/react-query'
import { toast } from 'sonner'
import { apiClient } from '../../lib/api-client'
import type { Exam, CreateExamInput, UpdateExamInput, ExamAccessInput } from './types'

export const examKeys = {
  all: ['exams'] as const,
  list: () => [...examKeys.all, 'list'] as const,
  detail: (id: number) => [...examKeys.all, 'detail', id] as const,
}

export function useExams() {
  return useQuery({
    queryKey: examKeys.list(),
    queryFn: () => apiClient<Exam[]>('/exams'),
  })
}

export function useExam(id: number) {
  return useQuery({
    queryKey: examKeys.detail(id),
    queryFn: () => apiClient<Exam>(`/exams/${id}`),
    enabled: id > 0,
  })
}

export function useCreateExam() {
  const queryClient = useQueryClient()
  return useMutation({
    mutationFn: (data: CreateExamInput) =>
      apiClient<Exam>('/exams', {
        method: 'POST',
        body: JSON.stringify(data),
      }),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: examKeys.all })
      toast.success('Paket ujian berhasil dibuat')
    },
    onError: (err: any) => {
      toast.error(err?.message || 'Gagal membuat paket ujian')
    },
  })
}

export function useUpdateExam() {
  const queryClient = useQueryClient()
  return useMutation({
    mutationFn: ({ id, data }: { id: number; data: UpdateExamInput }) =>
      apiClient<Exam>(`/exams/${id}`, {
        method: 'PATCH',
        body: JSON.stringify(data),
      }),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: examKeys.all })
      toast.success('Paket ujian berhasil diperbarui')
    },
    onError: (err: any) => {
      toast.error(err?.message || 'Gagal memperbarui paket ujian')
    },
  })
}

export function useDeleteExam() {
  const queryClient = useQueryClient()
  return useMutation({
    mutationFn: (id: number) =>
      apiClient<{ message: string }>(`/exams/${id}`, {
        method: 'DELETE',
      }),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: examKeys.all })
      toast.success('Paket ujian berhasil dihapus')
    },
    onError: (err: any) => {
      toast.error(err?.message || 'Gagal menghapus paket ujian')
    },
  })
}

export function useSetExamAccess() {
  const queryClient = useQueryClient()
  return useMutation({
    mutationFn: ({ id, data }: { id: number; data: ExamAccessInput }) =>
      apiClient<{ message: string }>(`/exams/${id}/access`, {
        method: 'PUT',
        body: JSON.stringify(data),
      }),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: examKeys.all })
      toast.success('Akses peserta ujian berhasil diperbarui')
    },
    onError: (err: any) => {
      toast.error(err?.message || 'Gagal memperbarui akses ujian')
    },
  })
}

export function useSetExamStatus() {
  const queryClient = useQueryClient()
  return useMutation({
    mutationFn: ({ id, isActive }: { id: number; isActive: boolean }) =>
      apiClient<{ message: string }>(`/exams/${id}/status`, {
        method: 'PATCH',
        body: JSON.stringify({ is_active: isActive }),
      }),
    onSuccess: (_, { isActive }) => {
      queryClient.invalidateQueries({ queryKey: examKeys.all })
      toast.success(isActive ? 'Paket ujian berhasil diaktifkan' : 'Paket ujian berhasil dinonaktifkan')
    },
    onError: (err: any) => {
      toast.error(err?.message || 'Gagal memperbarui status ujian')
    },
  })
}
