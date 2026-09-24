import { useQuery, useMutation, useQueryClient } from '@tanstack/react-query'
import { toast } from 'sonner'
import { apiClient } from '../../lib/api-client'
import type { Material, Argument, CreateMaterialInput, CreateArgumentInput } from './types'

export const materialKeys = {
  all: ['materials'] as const,
  list: () => [...materialKeys.all, 'list'] as const,
  detail: (id: number) => [...materialKeys.all, 'detail', id] as const,
  arguments: (id: number) => [...materialKeys.all, 'arguments', id] as const,
}

export function useMaterials() {
  return useQuery({
    queryKey: materialKeys.list(),
    queryFn: () => apiClient<Material[]>('/materials'),
  })
}

export function useMaterial(id: number) {
  return useQuery({
    queryKey: materialKeys.detail(id),
    queryFn: () => apiClient<Material>(`/materials/${id}`),
    enabled: id > 0,
  })
}

export function useCreateMaterial() {
  const queryClient = useQueryClient()
  return useMutation({
    mutationFn: (data: CreateMaterialInput) =>
      apiClient<Material>('/materials', {
        method: 'POST',
        body: JSON.stringify(data),
      }),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: materialKeys.all })
      toast.success('Materi berhasil dibuat')
    },
    onError: (err: any) => {
      toast.error(err?.message || 'Gagal membuat materi')
    },
  })
}

export function useUpdateMaterial() {
  const queryClient = useQueryClient()
  return useMutation({
    mutationFn: ({ id, data }: { id: number; data: CreateMaterialInput }) =>
      apiClient<Material>(`/materials/${id}`, {
        method: 'PATCH',
        body: JSON.stringify(data),
      }),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: materialKeys.all })
      toast.success('Materi berhasil diperbarui')
    },
    onError: (err: any) => {
      toast.error(err?.message || 'Gagal memperbarui materi')
    },
  })
}

export function useDeleteMaterial() {
  const queryClient = useQueryClient()
  return useMutation({
    mutationFn: (id: number) =>
      apiClient<{ message: string }>(`/materials/${id}`, {
        method: 'DELETE',
      }),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: materialKeys.all })
      toast.success('Materi berhasil dihapus')
    },
    onError: (err: any) => {
      toast.error(err?.message || 'Gagal menghapus materi')
    },
  })
}

export function useCreateArgument() {
  const queryClient = useQueryClient()
  return useMutation({
    mutationFn: ({ materialId, data }: { materialId: number; data: CreateArgumentInput }) =>
      apiClient<Argument>(`/materials/${materialId}/arguments`, {
        method: 'POST',
        body: JSON.stringify(data),
      }),
    onSuccess: (_, { materialId }) => {
      queryClient.invalidateQueries({ queryKey: materialKeys.detail(materialId) })
      queryClient.invalidateQueries({ queryKey: materialKeys.list() })
      toast.success('Argumen beserta opsi berhasil ditambahkan')
    },
    onError: (err: any) => {
      toast.error(err?.message || 'Gagal menambahkan argumen')
    },
  })
}

export function useUpdateArgument() {
  const queryClient = useQueryClient()
  return useMutation({
    mutationFn: ({
      materialId,
      argumentId,
      data,
    }: {
      materialId: number
      argumentId: number
      data: CreateArgumentInput
    }) =>
      apiClient<Argument>(`/materials/${materialId}/arguments/${argumentId}`, {
        method: 'PUT',
        body: JSON.stringify(data),
      }),
    onSuccess: (_, { materialId }) => {
      queryClient.invalidateQueries({ queryKey: materialKeys.detail(materialId) })
      queryClient.invalidateQueries({ queryKey: materialKeys.list() })
      toast.success('Argumen berhasil diperbarui')
    },
    onError: (err: any) => {
      toast.error(err?.message || 'Gagal memperbarui argumen')
    },
  })
}

export function useDeleteArgument() {
  const queryClient = useQueryClient()
  return useMutation({
    mutationFn: ({ materialId, argumentId }: { materialId: number; argumentId: number }) =>
      apiClient<{ message: string }>(`/materials/${materialId}/arguments/${argumentId}`, {
        method: 'DELETE',
      }),
    onSuccess: (_, { materialId }) => {
      queryClient.invalidateQueries({ queryKey: materialKeys.detail(materialId) })
      queryClient.invalidateQueries({ queryKey: materialKeys.list() })
      toast.success('Argumen berhasil dihapus')
    },
    onError: (err: any) => {
      toast.error(err?.message || 'Gagal menghapus argumen')
    },
  })
}
