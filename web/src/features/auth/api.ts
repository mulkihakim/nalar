import { useMutation, useQuery } from '@tanstack/react-query'
import { apiClient } from '../../lib/api-client'
import type { LoginResponse, MeResponse } from '../../types'
import type { LoginFormData } from './types'

export async function loginUser(data: LoginFormData): Promise<LoginResponse> {
  return apiClient<LoginResponse>('/auth/login', {
    method: 'POST',
    body: JSON.stringify(data),
  })
}

export async function fetchMe(): Promise<MeResponse> {
  return apiClient<MeResponse>('/auth/me')
}

export function useLoginMutation() {
  return useMutation({
    mutationFn: loginUser,
  })
}

export function useMeQuery(enabled = true) {
  return useQuery({
    queryKey: ['auth', 'me'],
    queryFn: fetchMe,
    enabled,
    retry: false,
  })
}
